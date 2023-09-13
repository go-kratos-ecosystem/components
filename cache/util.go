package cache

import (
	"fmt"
	"reflect"
)

var (
	ErrDestMustBePointer = fmt.Errorf("cache: dest must be a pointer")
	ErrDestMustNotBeNil  = fmt.Errorf("cache: dest must not be nil")
)

func valueOf(src interface{}, dest interface{}) error {
	rv := reflect.ValueOf(dest)

	if rv.Kind() != reflect.Ptr {
		return ErrDestMustBePointer
	}

	if rv.IsNil() {
		return ErrDestMustNotBeNil
	}

	rv.Elem().Set(reflect.ValueOf(src))

	return nil
}
