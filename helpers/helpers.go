package helpers

import (
	"encoding/json"
	"fmt"
	"time"
)

// Retry retries the given function until it returns nil or the attempts are exhausted.
// `sleeps` is the time to sleep between each attempt.
// If `sleeps` is not provided, it will not sleep.
//
//	Retry(func() error { return nil }, 3) => nil
//	Retry(func() error { return nil }, 3, time.Second) => nil
//	Retry(func() error { return fmt.Errorf("error") }, 3) => error
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

// Until retries the given function until it returns true.
// `sleeps` is the time to sleep between each attempt.
// If `sleeps` is not provided, it will not sleep.
//
//	Until(func() bool { return true }) => true
//	Until(func() bool { return true }, time.Second) => true
func Until(fn func() bool, sleeps ...time.Duration) bool {
	var sleep time.Duration
	if len(sleeps) > 0 {
		sleep = sleeps[0]
	}

	for !fn() {
		if sleep > 0 {
			time.Sleep(sleep)
		}
	}
	return true
}

func UntilTimeout(fn func() bool, timeout time.Duration, sleeps ...time.Duration) error {
	ch := make(chan error, 1)
	go func() {
		Until(fn, sleeps...)
		ch <- nil
	}()
	select {
	case err := <-ch:
		defer close(ch)
		return err
	case <-time.After(timeout):
		return fmt.Errorf("helpers: timeout after %s", timeout.String())
	}
}

// Timeout runs the given function with a timeout.
// If the function does not return before the timeout, it returns an error.
//
//	Timeout(func() error { return nil }, time.Second) => nil
//	Timeout(func() error { time.Sleep(2 * time.Second); return nil }, time.Second) => error
func Timeout(fn func() error, timeout time.Duration) error {
	ch := make(chan error, 1)
	go func() {
		ch <- fn()
	}()
	select {
	case err := <-ch:
		defer close(ch)
		return err
	case <-time.After(timeout):
		return fmt.Errorf("helpers: timeout after %s", timeout.String())
	}
}

// Repeat runs the given function `times` times or until an error is returned.
//
//	Repeat(func() error { fmt.Println("hello"); return nil }, 3) => prints hello 3 times and returns nil
//	Repeat(func() error { return fmt.Errorf("error") }, 3) => returns error
func Repeat(fn func() error, times int) error {
	for i := 0; i < times; i++ {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}

// ErrorIf returns an error if the condition is true.
//
//	ErrorIf(true, "error") => error
//	ErrorIf(false, "error") => nil
//	ErrorIf(true, "error %s", "with value") => error with value
func ErrorIf(condition bool, format string, a ...any) error {
	if condition {
		return fmt.Errorf(format, a...)
	}
	return nil
}

// PanicIf panics if the condition is true.
//
//	PanicIf(true, "error") => panic("error")
//	PanicIf(false, "error") => nil
//	PanicIf(true, "error %s", "with value") => panic("error with value")
func PanicIf(condition bool, format string, a ...any) {
	if condition {
		panic(fmt.Sprintf(format, a...))
	}
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

// Scan sets the value of dest to the value of src.
//
//	var foo string
//	Scan("bar", &foo) // foo == "bar"
//
//	var bar struct {A string}
//	Scan(struct{A string}{"foo"}, &bar) // bar == struct{A string}{"foo"}
func Scan(src any, dest any) error {
	bytes, err := json.Marshal(src)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, dest)
}
