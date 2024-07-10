package memory

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"

	"github.com/go-kratos-ecosystem/components/v2/locker"
)

type Locker struct {
	locked     atomic.Bool
	lockedTime time.Time  // locked time
	mu         sync.Mutex // guard lockedTime

	seconds time.Duration
	owner   string
	sleep   time.Duration
}

type Option func(*Locker)

func WithOwner(owner string) Option {
	return func(l *Locker) {
		l.owner = owner
	}
}

func WithSleep(sleep time.Duration) Option {
	return func(l *Locker) {
		l.sleep = sleep
	}
}

var _ locker.Locker = (*Locker)(nil)

func NewLocker(seconds time.Duration, opts ...Option) *Locker {
	l := &Locker{
		seconds: seconds,
		owner:   uuid.New().String(),
		sleep:   time.Millisecond * 100, //nolint:gomnd
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

func (m *Locker) acquire() bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	acquired := m.locked.CompareAndSwap(false, true)
	if acquired {
		m.lockedTime = time.Now()
	} else {
		if time.Since(m.lockedTime) > m.seconds {
			m.locked.Store(true)
			m.lockedTime = time.Now()
			return true
		}
	}
	return acquired
}

func (m *Locker) Try(ctx context.Context, fn func() error) error {
	if !m.acquire() {
		return locker.ErrLocked
	}

	defer m.Release(ctx) //nolint:errcheck
	return fn()
}

func (m *Locker) Until(ctx context.Context, timeout time.Duration, fn func() error) error {
	starting := time.Now()

	for m.acquire() {
		if time.Since(starting) > timeout {
			return locker.ErrTimeout
		}

		time.Sleep(m.sleep)
	}

	defer m.Release(ctx) //nolint:errcheck
	return fn()
}

func (m *Locker) Release(context.Context) (bool, error) {
	if m.locked.CompareAndSwap(true, false) {
		return false, locker.ErrNotLocked
	}
	return true, nil
}

func (m *Locker) ForceRelease(context.Context) error {
	m.locked.Store(false)
	return nil
}

func (m *Locker) Owner() string {
	return m.owner
}

func (m *Locker) LockedOwner(context.Context) (string, error) {
	if m.locked.Load() {
		return m.owner, nil
	}
	return "", locker.ErrNotLocked
}
