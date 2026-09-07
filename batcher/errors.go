package batcher

import "errors"

var (
	// ErrInvalidOption indicates an invalid batch size, interval, or queue size.
	ErrInvalidOption = errors.New("batcher: invalid option")
	// ErrNilHandler indicates that [New] received a nil handler.
	ErrNilHandler = errors.New("batcher: handler must not be nil")
	// ErrInvalidContext indicates that a nil context was supplied.
	ErrInvalidContext = errors.New("batcher: context must not be nil")
	// ErrClosed indicates that a batcher no longer accepts items or flush barriers.
	ErrClosed = errors.New("batcher: closed")
	// ErrFull indicates that [Batcher.TryAdd] could not immediately accept an item.
	ErrFull = errors.New("batcher: queue is full")
)
