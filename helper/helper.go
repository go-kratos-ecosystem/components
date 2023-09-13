package helper

import (
	"fmt"
	"reflect"
)

// Tap calls the given callback with the given value then returns the value.
func Tap(value interface{}, callbacks ...func(interface{})) interface{} {
	for _, callback := range callbacks {
		callback(value)
	}

	return value
}

// With calls the given callbacks with the given value then return the value.
func With(value interface{}, callbacks ...func(interface{}) interface{}) interface{} {
	for _, callback := range callbacks {
		value = callback(value)
	}

	return value
}

// ValueOf sets the value of dest to the value of src.
//
// Example:
//
//	var foo string
//	ValueOf("bar", &foo)
func ValueOf(src interface{}, dest interface{}) error {
	rv := reflect.ValueOf(dest)

	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("dest must be a pointer")
	}

	if rv.IsNil() {
		return fmt.Errorf("dest must not be nil")
	}

	rv.Elem().Set(reflect.ValueOf(src))

	return nil
}

// Call calls the given function with the given params.
//
// Example:
//
//	Call(func(name string) string {
//		return "Hello " + name
//	}, "world")
func Call(fn interface{}, params ...interface{}) interface{} {
	return CallWithCtx(nil, fn, params...)
}

// CallWithCtx calls the given function with the given context and params.
//
// Example:
//
//	type Foo struct {
//		Name string
//	}
//
//	result := CallWithCtx(&Foo{Name: "Hello"}, func(ts *Foo, name string) string {
//		return ts.Name + name
//	}, "world")
func CallWithCtx(ctx interface{}, fn interface{}, params ...interface{}) interface{} {
	fv := reflect.ValueOf(fn)

	if fv.Kind() != reflect.Func {
		panic("fn must be a function")
	}

	var args []reflect.Value

	if ctx != nil {
		args = append(args, reflect.ValueOf(ctx))
	}

	for _, param := range params {
		args = append(args, reflect.ValueOf(param))
	}

	return fv.Call(args)[0].Interface()
}

// Retry retries the given function until it returns nil or max is reached.
//
// Example:
//
//	err := Retry(func() error {
//		return nil
//	}, 3)
func Retry(fn func() error, max int) error {
	var err error

	for i := 0; i < max; i++ {
		err = fn()

		if err == nil {
			break
		}
	}

	return err
}
