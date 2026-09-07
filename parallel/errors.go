package parallel

import "errors"

var (
	// ErrInvalidLimit indicates that a concurrency limit is not positive.
	ErrInvalidLimit = errors.New("parallel: concurrency limit must be greater than zero")
	// ErrNilTask indicates that a task passed to [Run], [RunLimit], [Pool.Submit],
	// [Pool.TrySubmit], or [Pool.Execute] is nil.
	ErrNilTask = errors.New("parallel: task must not be nil")
	// ErrNilFunc indicates that a callback passed to [ForEach], [Map], [MapResults],
	// [Filter], [SubmitValue], or [MapStream] is nil.
	ErrNilFunc = errors.New("parallel: function must not be nil")
	// ErrPoolClosed indicates that a pool no longer accepts tasks.
	ErrPoolClosed = errors.New("parallel: pool is closed")
	// ErrPoolFull indicates that [Pool.TrySubmit] could not immediately accept a task.
	ErrPoolFull = errors.New("parallel: pool is full")
	// ErrNilInput indicates that [MapStream] received a nil input channel.
	ErrNilInput = errors.New("parallel: input channel must not be nil")
)
