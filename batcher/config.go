package batcher

import (
	"fmt"
	"time"
)

type config struct {
	batchSize     int
	flushInterval time.Duration
	queueSize     int
	onResult      func(Result)
}

// Option configures a [Batcher].
type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (f optionFunc) apply(c *config) {
	f(c)
}

// WithBatchSize sets the maximum number of items passed to a handler at once.
// The default is 100. [New] returns [ErrInvalidOption] if size is not positive.
func WithBatchSize(size int) Option {
	return optionFunc(func(c *config) { c.batchSize = size })
}

// WithFlushInterval sets how long a partial batch may accumulate after the
// worker receives its first item. The default is one second.
//
// Time spent in the queue or executing a handler is additional to this interval.
// [New] returns [ErrInvalidOption] if interval is not positive.
func WithFlushInterval(interval time.Duration) Option {
	return optionFunc(func(c *config) { c.flushInterval = interval })
}

// WithQueueSize sets the number of pending items or flush barriers the queue
// can hold, excluding the current batch. The default is 1000.
//
// Zero creates an unbuffered queue. [New] returns [ErrInvalidOption] for a
// negative size.
func WithQueueSize(size int) Option {
	return optionFunc(func(c *config) { c.queueSize = size })
}

// WithOnResult sets a callback invoked after each handler returns, including
// successful calls. A nil callback disables result notifications.
//
// The callback runs synchronously on the worker. It must return promptly and
// must not synchronously call blocking methods on the same [Batcher].
// [Batcher.Flush] and [Batcher.Shutdown] wait for the callback to return.
func WithOnResult(fn func(Result)) Option {
	return optionFunc(func(c *config) { c.onResult = fn })
}

func newConfig(options ...Option) (config, error) {
	c := config{batchSize: 100, flushInterval: time.Second, queueSize: 1000}
	for _, option := range options {
		if option != nil {
			option.apply(&c)
		}
	}
	if c.batchSize <= 0 {
		return config{}, fmt.Errorf("%w: batch size must be positive", ErrInvalidOption)
	}
	if c.flushInterval <= 0 {
		return config{}, fmt.Errorf("%w: flush interval must be positive", ErrInvalidOption)
	}
	if c.queueSize < 0 {
		return config{}, fmt.Errorf("%w: queue size must not be negative", ErrInvalidOption)
	}

	return c, nil
}
