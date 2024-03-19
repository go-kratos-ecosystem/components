package values

import "encoding/json"

// If returns trueVal if condition is true, otherwise falseVal
//
//	If(true, "foo", "bar") // "foo"
//	If(false, "foo", "bar") // "bar"
func If[T any](condition bool, trueVal T, falseVal T) T {
	if condition {
		return trueVal
	}

	return falseVal
}

func NilIf[T any](condition bool, trueVal T) *T {
	if condition {
		return &trueVal
	}

	return nil
}

// Tap calls the given callback with the given value then returns the value.
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
//	Transform(1, strconv.Itoa) // "1"
//	Transform("foo", func(s string) *foo {
//		return &foo{Name: s}
//	}) // &foo{Name: "foo"}
func Transform[T, R any](value T, callback func(T) R) R {
	return callback(value)
}

// When calls the given callbacks with the given value if condition is true then return the value.
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

// Default returns defaultValue if value is zero, otherwise value.
//
//	Default("", "foo") // "foo"
//	Default("bar", "foo") // "bar"
func Default[T comparable](value T, defaultValue T) T {
	var zero T
	if value == zero {
		return defaultValue
	}
	return value
}

// DefaultWithFunc returns defaultValue if value is zero, otherwise value.
//
//	DefaultWithFunc("", func() string { return "foo" }) // "foo"
//	DefaultWithFunc("bar", func() string { return "foo" }) // "bar"
func DefaultWithFunc[T comparable](value T, defaultValue func() T) T {
	var zero T
	if value == zero {
		return defaultValue()
	}
	return value
}

// Ptr returns a pointer to the value.
//
//	Ptr("foo") // *string("foo")
//	Ptr(1) // *int(1)
func Ptr[T any](value T) *T {
	return &value
}

// Val returns the value of the pointer.
// If the pointer is nil, return the zero value.
//
//	Val((*string)(nil)) // ""
//	Val(Ptr("foo")) // "foo"
func Val[T any](value *T) T {
	if value != nil {
		return *value
	}
	var zero T
	return zero
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
