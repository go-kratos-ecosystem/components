package helper

import (
	"encoding/json"
)

func Tap[T any](value T, callbacks ...func(T)) T {
	for _, callback := range callbacks {
		if callback != nil {
			callback(value)
		}
	}

	return value
}

func With[T any](value T, callbacks ...func(T) T) T {
	for _, callback := range callbacks {
		if callback != nil {
			value = callback(value)
		}
	}

	return value
}

func Chain[T any](fns ...func(T) T) func(T) T {
	return func(v T) T {
		for _, fn := range fns {
			v = fn(v)
		}
		return v
	}
}

func When[T any](value T, condition bool, callbacks ...func(T) T) T {
	if condition {
		return With(value, callbacks...)
	}

	return value
}

func Scan(src any, dest any) error {
	bytes, err := json.Marshal(src)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, dest)
}
