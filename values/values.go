package values

import (
	"encoding/json"
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

// Default returns defaultValue if value is zero, otherwise value.
func Default[T comparable](value T, defaultValue T) T {
	var zero T
	if value == zero {
		return defaultValue
	}
	return value
}

// DefaultWithFunc returns defaultValue if value is zero, otherwise value.
func DefaultWithFunc[T comparable](value T, defaultValue func() T) T {
	var zero T
	if value == zero {
		return defaultValue()
	}
	return value
}
