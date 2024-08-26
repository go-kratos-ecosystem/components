package errors

import (
	"errors"
	"fmt"
	"time"
)

type TimeoutError struct {
	timeout time.Duration
	err     error
}

func NewTimeoutError(timeout time.Duration, err error) *TimeoutError {
	return &TimeoutError{
		timeout: timeout,
		err:     err,
	}
}

func (e *TimeoutError) Error() string {
	return fmt.Sprintf("timeout after %s: %v", e.timeout, e.err)
}

func (e *TimeoutError) Unwrap() error {
	return e.err
}

func IsTimeoutError(err error) bool {
	var target *TimeoutError
	return errors.As(err, &target)
}
