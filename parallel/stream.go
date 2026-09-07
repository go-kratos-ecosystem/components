package parallel

import (
	"context"
	"fmt"
	"sync"
)

// StreamResult contains the outcome of processing one stream input.
type StreamResult[T any] struct {
	// Index is the zero-based position in the stream's input receive order.
	Index int
	// Value is the callback result, including a value returned alongside an error.
	Value T
	// Err is the callback error for this input, or nil on success.
	Err error
}

// Stream delivers results from concurrent processing of channel inputs.
//
// A Stream must be created by [MapStream] and must not be copied after first use;
// its zero value is not valid. Call [Stream.Close] when stopping consumption
// early. The input channel remains owned by its producer.
//
// Stream methods are safe for concurrent use. Multiple receivers of Results
// share the same channel and divide results between them; results are not
// broadcast. Values are not deep-copied, so callers must synchronize mutations
// to referenced data.
type Stream[T any] struct {
	ctx     context.Context
	cancel  context.CancelFunc
	results chan StreamResult[T]
	done    chan struct{}
	err     error
}

// MapStream starts at most limit concurrent callbacks consuming input.
// Results are delivered as workers make them available, without preserving input
// order. Each result's index identifies its input receive position.
//
// A callback error is included in its [StreamResult] and does not cancel other
// callbacks or become the stream's terminal error. Results are unbuffered: each
// worker receives another input only after delivering its previous result. The
// stream holds at most limit in-flight inputs or results, excluding caller-owned
// buffers and data retained by callbacks.
//
// Closing input ends the stream after all received inputs have been processed
// and their results delivered. Canceling ctx or calling [Stream.Close] stops
// receiving and delivering; some received inputs may have no delivered result.
// Callbacks must observe their context and return promptly on cancellation.
// They must not synchronously call Close or Wait on their own stream. Panics
// follow normal Go semantics and are not recovered.
//
// MapStream returns an error wrapping [ErrInvalidLimit] if limit is not positive,
// an error wrapping [ErrNilFunc] if fn is nil, [ErrNilInput] if input is nil, or
// the cancellation cause if ctx is already canceled. The ctx must not be nil.
func MapStream[T, R any](
	ctx context.Context,
	limit int,
	input <-chan T,
	fn func(context.Context, T) (R, error),
) (*Stream[R], error) {
	if err := validateLimit(limit); err != nil {
		return nil, err
	}
	if fn == nil {
		return nil, fmt.Errorf("%w: MapStream", ErrNilFunc)
	}
	if input == nil {
		return nil, ErrNilInput
	}
	if err := context.Cause(ctx); err != nil {
		return nil, err
	}

	streamContext, cancel := context.WithCancel(ctx)
	stream := &Stream[R]{
		ctx:     streamContext,
		cancel:  cancel,
		results: make(chan StreamResult[R]),
		done:    make(chan struct{}),
	}
	source := &streamInput[T]{values: input, gate: make(chan struct{}, 1)}
	var workers sync.WaitGroup
	for range limit {
		workers.Go(func() {
			for {
				value, index, ok := source.receive(streamContext)
				if !ok || context.Cause(streamContext) != nil {
					return
				}
				result, err := fn(streamContext, value)
				select {
				case stream.results <- StreamResult[R]{Index: index, Value: result, Err: err}:
				case <-streamContext.Done():
					return
				}
			}
		})
	}
	go func() {
		workers.Wait()
		stream.err = context.Cause(streamContext)
		cancel()
		close(stream.results)
		close(stream.done)
	}()

	return stream, nil
}

// Context returns the context used by callbacks.
//
// It is canceled by [Stream.Close], parent cancellation, or normal completion.
// Producers can select on its Done channel to stop sending when consumption
// ends. Use [Stream.Wait] to distinguish normal completion from cancellation.
func (s *Stream[T]) Context() context.Context {
	return s.ctx
}

// Results returns the unbuffered result channel, closed when all workers exit.
// Callback errors appear in individual results. Use [Stream.Wait] after consuming
// results to obtain the terminal status.
func (s *Stream[T]) Results() <-chan StreamResult[T] {
	return s.results
}

// Close cancels the stream and waits for its workers to return.
//
// Close is safe to call repeatedly or concurrently. It unblocks workers waiting
// for input or result delivery, but still waits for running callbacks. It does
// not close the input channel or wait for caller-owned producers. Those producers
// can use [Stream.Context] to observe cancellation.
func (s *Stream[T]) Close() {
	s.cancel()
	<-s.done
}

// Wait waits for the stream to finish or ctx to be canceled.
//
// On normal completion, Wait returns nil even if individual callbacks failed.
// If the stream was canceled, Wait returns its cancellation cause. Closing an
// unfinished stream with [Stream.Close] causes Wait to return [context.Canceled].
// The terminal result is fixed at completion and can be read multiple times.
//
// Canceling ctx stops only this wait, returning [context.Cause] of ctx. An already
// completed stream's result takes precedence over cancellation of the wait.
// The ctx must not be nil. Results must be consumed, or the stream canceled,
// before Wait can complete; Wait does not drain results on behalf of the caller.
func (s *Stream[T]) Wait(ctx context.Context) error {
	select {
	case <-s.done:
		return s.err
	default:
	}
	select {
	case <-s.done:
		return s.err
	case <-ctx.Done():
		select {
		case <-s.done:
			return s.err
		default:
			return context.Cause(ctx)
		}
	}
}

type streamInput[T any] struct {
	gate   chan struct{}
	values <-chan T
	next   int
}

func (s *streamInput[T]) receive(ctx context.Context) (T, int, bool) {
	// Receiving and assigning the index must be one operation. Assigning an
	// atomic index after an independent receive could reorder input positions.
	var zero T
	select {
	case s.gate <- struct{}{}:
		defer func() { <-s.gate }()
	case <-ctx.Done():
		return zero, 0, false
	}
	if context.Cause(ctx) != nil {
		return zero, 0, false
	}
	select {
	case value, ok := <-s.values:
		if !ok {
			return zero, 0, false
		}
		index := s.next
		s.next++

		return value, index, true
	case <-ctx.Done():
		return zero, 0, false
	}
}
