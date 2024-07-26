package snap

import (
	"sync"
	"time"
)

type Snap[T any] struct {
	value T
	mu    sync.RWMutex
	err   error

	refresh  func() (T, error)
	expired  time.Time     // expiration time
	interval time.Duration // refresh interval
}

type Option[T any] func(*Snap[T])

func Interval[T any](interval time.Duration) Option[T] {
	return func(s *Snap[T]) {
		s.interval = interval
	}
}

func New[T any](refresh func() (T, error), opts ...Option[T]) *Snap[T] {
	s := &Snap[T]{
		refresh:  refresh,
		interval: time.Minute,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Snap[T]) Get() (T, error) {
	if s.expired.IsZero() || time.Now().After(s.expired) {
		return s.Refresh()
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.value, s.err
}

func (s *Snap[T]) Refresh() (T, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.value, s.err = s.refresh()
	s.expired = time.Now().Add(s.interval)

	return s.value, s.err
}
