package debounce

import "errors"

var (
	// ErrInvalidDelay indicates that [New] received a non-positive delay.
	ErrInvalidDelay = errors.New("debounce: delay must be positive")
	// ErrNilHandler indicates that [New] received a nil handler.
	ErrNilHandler = errors.New("debounce: handler must not be nil")
	// ErrInvalidContext indicates that a nil context was supplied.
	ErrInvalidContext = errors.New("debounce: context must not be nil")
	// ErrClosed indicates that a debouncer is closed or an invocation was
	// discarded before starting because the debouncer closed.
	ErrClosed = errors.New("debounce: closed")
)
