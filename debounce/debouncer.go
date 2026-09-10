package debounce

import (
	"context"
	"sync"
	"time"
)

// Handler processes the latest value selected by a [Debouncer].
//
// Calls are serial and receive the context derived from the one supplied to
// [New]. A handler should return promptly when its context is canceled. It may
// call [Debouncer.Trigger], but must not synchronously call [Debouncer.Flush] or
// [Debouncer.Close] on its own debouncer. Panics follow normal Go semantics.
type Handler[T any] func(context.Context, T) error

// Debouncer delays processing until triggers have been quiet for a full interval.
//
// A Debouncer is safe for concurrent use. It must be created by [New], must not
// be copied after first use, and must be closed when no longer needed. Its zero
// value is not valid. Values are not deep-copied; callers must synchronize
// mutations to referenced data that might be read by a handler.
//
// One worker processes invocations serially. Internally, at most one value can
// be running, one frozen by a flush, and one pending replacement. Handler errors
// do not stop subsequent invocations; use [WithOnError] to observe background
// failures.
type Debouncer[T any] struct {
	ctx     context.Context
	cancel  context.CancelFunc
	delay   time.Duration
	handler Handler[T]
	config  config
	wake    chan struct{}
	done    chan struct{}

	mu      sync.Mutex
	closed  bool
	pending *invocation[T]
	ready   *invocation[T]
	running *invocation[T]
}

// New starts a trailing-edge debouncer with the given delay and handler.
//
// Each [Debouncer.Trigger] restarts the delay for the replaceable pending value.
// Canceling ctx closes the debouncer, discards values not yet selected for
// execution, and cancels the running handler's context. Use [Debouncer.Close]
// to wait for cleanup.
//
// New returns [ErrInvalidContext] for a nil ctx, its cancellation cause if ctx
// is already canceled, [ErrInvalidDelay] for a non-positive delay, or
// [ErrNilHandler] for a nil handler. No worker starts on failure. Nil options
// are ignored.
func New[T any](ctx context.Context, delay time.Duration, handler Handler[T], options ...Option) (*Debouncer[T], error) {
	if err := validateContext(ctx); err != nil {
		return nil, err
	}
	if delay <= 0 {
		return nil, ErrInvalidDelay
	}
	if handler == nil {
		return nil, ErrNilHandler
	}

	processingContext, cancel := context.WithCancel(ctx)
	d := &Debouncer[T]{
		ctx:     processingContext,
		cancel:  cancel,
		delay:   delay,
		handler: handler,
		config:  newConfig(options...),
		wake:    make(chan struct{}, 1),
		done:    make(chan struct{}),
	}
	stopCancellation := context.AfterFunc(processingContext, d.stop)
	go func() {
		defer stopCancellation()
		d.run()
		d.stop()
		close(d.done)
	}()

	return d, nil
}

// Trigger replaces the pending value and restarts its quiet interval.
//
// Trigger returns once the value is stored, without waiting for the handler.
// A later trigger can replace it before execution. Triggers during a running
// handler form the next invocation and never start an overlapping handler.
// A value already frozen by [Debouncer.Flush] is not replaced.
//
// Trigger returns [ErrClosed] if the debouncer is closed or its processing
// context has been canceled. Concurrent triggers are ordered by acceptance.
func (d *Debouncer[T]) Trigger(value T) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.closed || context.Cause(d.ctx) != nil {
		return ErrClosed
	}
	if d.pending == nil {
		d.pending = &invocation[T]{done: make(chan struct{})}
	}
	d.pending.value = value
	d.pending.due = time.Now().Add(d.delay)
	d.notify()

	return nil
}

// Flush freezes the latest pending value for execution without its remaining delay.
// It waits for that invocation and its error callback, returning the handler error.
//
// Later triggers cannot replace a frozen value or postpone its execution. If a
// previous flush already occupies the frozen slot and a newer value is pending,
// Flush waits for that invocation to finish before capturing the latest pending
// value. Pending values can still be replaced during this admission wait.
//
// With no pending value, Flush waits for the existing frozen or running
// invocation. If idle, it returns nil. Concurrent flushes may share a result.
// Completed errors are not replayed; use [WithOnError] for background failures.
//
// Canceling ctx stops only admission or waiting, returning [context.Cause] of
// ctx. An already frozen invocation continues using the processing context.
// A completed result takes precedence over cancellation while waiting. A nil
// ctx returns [ErrInvalidContext]. A closed debouncer or an invocation discarded
// by [Debouncer.Close] returns [ErrClosed].
func (d *Debouncer[T]) Flush(ctx context.Context) error {
	if err := validateContext(ctx); err != nil {
		return err
	}
	for {
		call, retry, err := d.selectFlush()
		if err != nil || call == nil {
			return err
		}
		if err := wait(ctx, call.done); err != nil {
			return err
		}
		if !retry {
			return call.err
		}
		if err := context.Cause(ctx); err != nil {
			return err
		}
	}
}

// Close stops accepting triggers, discards pending and frozen values, and waits
// for the running invocation and worker to exit.
//
// Close cancels the processing context but still waits for a running handler and
// its error callback to return. Values already selected by the worker count as
// running. Call Flush before Close when the final pending value must be processed.
// Close is safe to call repeatedly or concurrently.
func (d *Debouncer[T]) Close() {
	d.stop()
	<-d.done
}

func (d *Debouncer[T]) selectFlush() (*invocation[T], bool, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.closed || context.Cause(d.ctx) != nil {
		return nil, false, ErrClosed
	}
	if d.ready != nil {
		return d.ready, d.pending != nil, nil
	}
	if d.pending != nil {
		d.ready, d.pending = d.pending, nil
		d.notify()

		return d.ready, false, nil
	}

	return d.running, false, nil
}

func (d *Debouncer[T]) stop() {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.closed {
		return
	}
	d.closed = true
	d.cancel()
	if d.pending != nil {
		d.pending.complete(ErrClosed)
		d.pending = nil
	}
	if d.ready != nil {
		d.ready.complete(ErrClosed)
		d.ready = nil
	}
	d.notify()
}

func (d *Debouncer[T]) notify() {
	select {
	case d.wake <- struct{}{}:
	default:
	}
}

func validateContext(ctx context.Context) error {
	if ctx == nil {
		return ErrInvalidContext
	}

	return context.Cause(ctx)
}

type invocation[T any] struct {
	value T
	due   time.Time
	done  chan struct{}
	err   error
}

func (c *invocation[T]) complete(err error) {
	var zero T
	c.value = zero
	c.err = err
	close(c.done)
}

func wait(ctx context.Context, done <-chan struct{}) error {
	select {
	case <-done:
		return nil
	default:
	}
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		select {
		case <-done:
			return nil
		default:
			return context.Cause(ctx)
		}
	}
}
