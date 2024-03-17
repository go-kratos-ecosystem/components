package helpers

import "time"

// Retry retries the given function until it returns nil or the attempts are exhausted.
func Retry(fn func() error, attempts int, sleeps ...time.Duration) (err error) {
	var sleep time.Duration
	if len(sleeps) > 0 {
		sleep = sleeps[0]
	}

	for i := 0; i < attempts; i++ {
		if err = fn(); err == nil {
			return nil
		}
		if sleep > 0 {
			time.Sleep(sleep)
		}
	}
	return
}

// Pipe is a function that takes a value and returns a value
//
//	Pipe(m1, m2, m3)(value) => m3(m2(m1(value)))
func Pipe[T any](fns ...func(T) T) func(T) T {
	return func(v T) T {
		for _, fn := range fns {
			v = fn(v)
		}
		return v
	}
}

// PipeWithErr is a function that takes a value and returns a value and an error
//
//	PipeWithErr(m1, m2, m3)(value) => m3(m2(m1(value)))
func PipeWithErr[T any](fns ...func(T) (T, error)) func(T) (T, error) {
	var err error
	return func(v T) (T, error) {
		for _, fn := range fns {
			if v, err = fn(v); err != nil {
				return v, err
			}
		}
		return v, nil
	}
}

// Chain is a reverse Pipe
//
//	Chain(m1, m2, m3)(value) => m1(m2(m3(value)))
func Chain[T any](fns ...func(T) T) func(T) T {
	return func(v T) T {
		for i := len(fns) - 1; i >= 0; i-- {
			if fns[i] != nil {
				v = fns[i](v)
			}
		}
		return v
	}
}

// ChainWithErr is a reverse PipeWithErr
//
//	ChainWithErr(m1, m2, m3)(value) => m1(m2(m3(value)))
func ChainWithErr[T any](fns ...func(T) (T, error)) func(T) (T, error) {
	var err error
	return func(v T) (T, error) {
		for i := len(fns) - 1; i >= 0; i-- {
			if fns[i] != nil {
				if v, err = fns[i](v); err != nil {
					return v, err
				}
			}
		}
		return v, nil
	}
}
