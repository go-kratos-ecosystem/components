package helpers

import (
	"encoding/json"
	"time"
)

// If returns trueVal if condition is true, otherwise falseVal
//
// Example:
//
//	If(true, "foo", "bar") // "foo"
//	If(false, "foo", "bar") // "bar"
func If[T any](condition bool, trueVal T, falseVal T) T {
	if condition {
		return trueVal
	}

	return falseVal
}

// Tap calls the given callback with the given value then returns the value.
//
// Example:
//
//	Tap("foo", func(s string) {
//		fmt.Println(s) // "foo" and os.Stdout will print "foo"
//	}, func(s string) {
//		// more callbacks
//	}...)
func Tap[T any](value T, callbacks ...func(T)) T {
	for _, callback := range callbacks {
		if callback != nil {
			callback(value)
		}
	}

	return value
}

// With calls the given callbacks with the given value then return the value.
//
// Example:
//
//	With("foo", func(s string) string {
//		return s + "bar"
//	}, func(s string) string {
//		return s + "baz"
//	}) // "foobarbaz"
func With[T any](value T, callbacks ...func(T) T) T {
	for _, callback := range callbacks {
		if callback != nil {
			value = callback(value)
		}
	}

	return value
}

// Transform calls the given callback with the given value then return the result.
//
// Example:
//
//	Transform(1, strconv.Itoa) // "1"
//	Transform("foo", func(s string) *foo {
//		return &foo{Name: s}
//	}) // &foo{Name: "foo"}
func Transform[T, R any](value T, callback func(T) R) R {
	return callback(value)
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

// When calls the given callbacks with the given value if condition is true then return the value.
//
// Example:
//
//	When("foo", true, func(s string) string {
//		return s + "bar"
//	}, func(s string) string {
//		return s + "baz"
//	}) // "foobarbaz"
func When[T any](value T, condition bool, callbacks ...func(T) T) T {
	if condition {
		return With(value, callbacks...)
	}

	return value
}

// Scan sets the value of dest to the value of src.
//
// Example:
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

// Default returns defaultValue if value is zero, otherwise value.
func Default[T comparable](value T, defaultValue T) T {
	var zero T
	if value == zero {
		return defaultValue
	}
	return value
}

// DefaultWith returns defaultValue if value is zero, otherwise value.
func DefaultWith[T comparable](value T, defaultValue func() T) T {
	var zero T
	if value == zero {
		return defaultValue()
	}
	return value
}
