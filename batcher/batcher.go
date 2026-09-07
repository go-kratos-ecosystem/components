package batcher

import (
	"context"
	"sync"
)

// Handler processes a non-empty batch of items in acceptance order.
//
// Calls are serial and receive the processing context supplied to [New]. A
// handler should stop promptly when that context is canceled. Each call receives
// a distinct slice that the batcher will not reuse. Referenced data within T is
// shared with the producer, which must synchronize any mutations.
//
// A handler must not synchronously call blocking methods on its own [Batcher].
// Handler panics follow normal Go semantics and are not recovered.
type Handler[T any] func(context.Context, []T) error

// Result describes one completed handler call.
type Result struct {
	// Size is the number of items passed to the handler.
	Size int
	// Err is the error returned by the handler, or nil on success.
	Err error
}

type request[T any] struct {
	item  T
	flush *completion
}

// Batcher combines items into bounded batches processed by one worker.
//
// A Batcher is safe for concurrent use. It must be created by [New], must not be
// copied after first use, and must be shut down when no longer needed. Its zero
// value is not valid. Pending items are stored in memory until processed.
//
// Handler errors do not stop subsequent batches. [Batcher.Flush] and
// [Batcher.Shutdown] report the first handler error up to their completion
// boundary, including errors already reported by earlier calls. Use
// [WithOnResult] to observe every batch result.
type Batcher[T any] struct {
	requests chan request[T]
	closing  chan struct{}
	finished *completion

	mu         sync.Mutex
	closed     bool
	stopOnce   sync.Once
	submitters sync.WaitGroup
}

// New starts a batcher that processes items with handler.
//
// The ctx controls processing independently of the contexts passed to
// [Batcher.Add], [Batcher.Flush], and [Batcher.Shutdown]. Canceling ctx initiates
// shutdown asynchronously and drains accepted items; the handler is still called
// for every batch, with the canceled context.
//
// New returns [ErrInvalidContext] for a nil ctx, its cancellation cause if ctx
// is already canceled, [ErrNilHandler] for a nil handler, or an error wrapping
// [ErrInvalidOption] for invalid options. No worker is started on failure.
func New[T any](ctx context.Context, handler Handler[T], options ...Option) (*Batcher[T], error) {
	if err := validateContext(ctx); err != nil {
		return nil, err
	}
	if handler == nil {
		return nil, ErrNilHandler
	}
	c, err := newConfig(options...)
	if err != nil {
		return nil, err
	}

	b := &Batcher[T]{
		requests: make(chan request[T], c.queueSize),
		closing:  make(chan struct{}),
		finished: newCompletion(),
	}
	stopCancellation := context.AfterFunc(ctx, b.stop)
	go func() {
		defer stopCancellation()
		w := worker[T]{ctx: ctx, handler: handler, config: c}
		b.finished.complete(w.run(b.requests))
	}()

	return b, nil
}

// Add waits until item is accepted, ctx is canceled, or the batcher is closed.
//
// A nil error means the item was accepted, not that processing succeeded. The
// ctx controls only admission; canceling it after acceptance does not remove
// the item or cancel processing. Items rejected with an error are never processed.
// Concurrent calls are ordered by their actual acceptance into the queue.
//
// Add returns [ErrInvalidContext] for a nil ctx, [context.Cause] of ctx when
// admission is canceled, or [ErrClosed] when admission has closed.
func (b *Batcher[T]) Add(ctx context.Context, item T) error {
	return b.submit(ctx, request[T]{item: item})
}

// TryAdd attempts to accept item without waiting for queue capacity.
//
// TryAdd returns [ErrFull] when immediate acceptance is unavailable, or
// [ErrClosed] when the batcher is closed. With an unbuffered queue, the worker
// must be ready to receive. Rejected items are never processed.
func (b *Batcher[T]) TryAdd(item T) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return ErrClosed
	}

	select {
	case b.requests <- request[T]{item: item}:
		return nil
	default:
		return ErrFull
	}
}

// Flush waits for items accepted before its queue barrier to finish processing.
// This includes all items accepted before Flush was called. A partial batch is
// processed without waiting for the flush interval.
//
// Flush returns the first handler error since construction through this barrier,
// or nil if all those batches succeeded. Errors are not cleared by a flush.
// Items accepted after the barrier do not delay its completion.
//
// Canceling ctx stops admission of the barrier or waiting for it. An accepted
// barrier continues in the background. Flush returns [ErrClosed] if admission
// has closed; use [Batcher.Shutdown] to wait for the final result. A nil ctx
// returns [ErrInvalidContext]. Canceled admission or an unfinished wait returns
// [context.Cause] of ctx; an already completed barrier's result takes precedence
// over cancellation while waiting.
func (b *Batcher[T]) Flush(ctx context.Context) error {
	result := newCompletion()
	if err := b.submit(ctx, request[T]{flush: result}); err != nil {
		return err
	}

	return result.wait(ctx)
}

// Shutdown stops admission and waits for all accepted items and result callbacks.
// It processes the final partial batch and returns the first handler error since
// construction, or nil if every batch succeeded. Repeated calls share this result.
//
// Canceling ctx stops only the wait; draining continues in the background using
// the processing context supplied to [New]. Calls already submitting when
// shutdown begins may still succeed, and their items are included in the drain.
// A nil ctx returns [ErrInvalidContext] without initiating shutdown. If draining
// has completed, its result takes precedence over cancellation of ctx.
func (b *Batcher[T]) Shutdown(ctx context.Context) error {
	if ctx == nil {
		return ErrInvalidContext
	}
	b.stop()

	return b.finished.wait(ctx)
}

func (b *Batcher[T]) submit(ctx context.Context, req request[T]) error {
	if err := validateContext(ctx); err != nil {
		return err
	}
	b.mu.Lock()
	if b.closed {
		b.mu.Unlock()

		return ErrClosed
	}
	b.submitters.Add(1)
	b.mu.Unlock()
	defer b.submitters.Done()

	select {
	case b.requests <- req:
		return nil
	case <-ctx.Done():
		return context.Cause(ctx)
	case <-b.closing:
		return ErrClosed
	}
}

func (b *Batcher[T]) stop() {
	b.stopOnce.Do(func() {
		b.mu.Lock()
		b.closed = true
		close(b.closing)
		b.mu.Unlock()
		go func() {
			b.submitters.Wait()
			close(b.requests)
		}()
	})
}

func validateContext(ctx context.Context) error {
	if ctx == nil {
		return ErrInvalidContext
	}

	return context.Cause(ctx)
}

type completion struct {
	done chan struct{}
	err  error
}

func newCompletion() *completion {
	return &completion{done: make(chan struct{})}
}

func (c *completion) complete(err error) {
	c.err = err
	close(c.done)
}

func (c *completion) wait(ctx context.Context) error {
	select {
	case <-c.done:
		return c.err
	default:
	}
	select {
	case <-c.done:
		return c.err
	case <-ctx.Done():
		select {
		case <-c.done:
			return c.err
		default:
			return context.Cause(ctx)
		}
	}
}
