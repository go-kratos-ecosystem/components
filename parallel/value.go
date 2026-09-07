package parallel

import (
	"context"
	"fmt"
)

// ValueFuture represents the eventual value and error of a task accepted by a [Pool].
//
// A ValueFuture supports concurrent and repeated waits. Values are shared, not
// copied deeply; callers must synchronize mutations to referenced data.
//
// ValueFuture values are created by [SubmitValue]; the zero value is not valid.
type ValueFuture[T any] struct {
	future *Future
	value  T
}

// Done returns a channel that is closed when the task finishes.
func (f *ValueFuture[T]) Done() <-chan struct{} {
	return f.future.Done()
}

// Wait waits for the task to finish or ctx to be canceled.
//
// When the task finishes, Wait returns its value and error, including a value
// returned alongside an error. A completed task's result takes precedence over
// cancellation of ctx.
//
// If ctx is canceled while the task is still running, Wait returns the zero
// value of T and [context.Cause] of ctx. Canceling ctx stops only this wait;
// it does not cancel the task.
func (f *ValueFuture[T]) Wait(ctx context.Context) (T, error) {
	err := f.future.Wait(ctx)
	select {
	case <-f.Done():
		return f.value, f.future.err
	default:
		var zero T

		return zero, err
	}
}

// SubmitValue submits fn to pool and returns a future once the task is accepted.
//
// Submission applies the same backpressure as [Pool.Submit]. The supplied ctx
// controls the task's lifetime, including time spent waiting in the queue.
// If the task is canceled before starting, its future returns the zero value of
// T and the cancellation cause.
//
// SubmitValue returns an error wrapping [ErrNilFunc] if fn is nil. Otherwise,
// submission errors are returned as documented by [Pool.Submit]. No future is
// returned when submission fails.
//
// The pool must be non-nil and created by [NewPool].
func SubmitValue[T any](ctx context.Context, pool *Pool, fn func(context.Context) (T, error)) (*ValueFuture[T], error) {
	if fn == nil {
		return nil, fmt.Errorf("%w: SubmitValue", ErrNilFunc)
	}

	result := &ValueFuture[T]{}
	future, err := pool.Submit(ctx, func(ctx context.Context) error {
		value, err := fn(ctx)
		result.value = value

		return err
	})
	if err != nil {
		return nil, err
	}
	result.future = future

	return result, nil
}
