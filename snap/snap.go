package snap

import (
	"sync"
	"time"
)

type Snap[T any] struct {
	value T
	mu    sync.RWMutex

	refresh  func() T
	expired  time.Time     // expiration time
	interval time.Duration // refresh interval
}

type Option[T any] func(*Snap[T])

func Interval[T any](interval time.Duration) Option[T] {
	return func(s *Snap[T]) {
		s.interval = interval
	}
}

func New[T any](refresh func() T, opts ...Option[T]) *Snap[T] {
	s := &Snap[T]{
		refresh:  refresh,
		interval: time.Minute,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Snap[T]) Get() T {
	if s.expired.IsZero() || time.Now().After(s.expired) {
		s.Refresh()
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.value
}

func (s *Snap[T]) Refresh() {
	value := s.refresh()

	s.mu.Lock()
	defer s.mu.Unlock()
	s.value = value
	s.expired = time.Now().Add(s.interval)
}
