package locker

import (
	"context"
	"errors"
	"time"
)

var (
	ErrLocked    = errors.New("locker: the locker is locked")
	ErrTimeout   = errors.New("locker: the locker is timeout")
	ErrNotLocked = errors.New("locker: the locker is not locked")
)

type Locker interface {
	// Try tries to get the lock.
	// if the lock is locked, it will return an error.
	Try(ctx context.Context, fn func()) error

	// Until tries to get the lock until the timeout.
	// if the lock is locked, it will wait until the lock is released or the timeout is reached.
	Until(ctx context.Context, timeout time.Duration, fn func()) error

	// Get tries to get the lock.
	// if the lock is locked, it will return the owner of the lock or an error.
	Get(ctx context.Context) (Owner, error)

	// Release releases the locked owner.
	// if the owner is not the locked owner, it will return an error.
	Release(ctx context.Context, owner Owner) error

	// ForceRelease releases the lock forcibly.
	ForceRelease(ctx context.Context) error

	// LockedOwner returns the owner of the lock.
	LockedOwner(ctx context.Context) (Owner, error)
}

type NoopLocker struct{}

var _ Locker = NoopLocker{}

func (l NoopLocker) Try(context.Context, func()) error                  { return nil }
func (l NoopLocker) Until(context.Context, time.Duration, func()) error { return nil }
func (l NoopLocker) Get(context.Context) (Owner, error)                 { return nil, nil }
func (l NoopLocker) Release(context.Context, Owner) error               { return nil }
func (l NoopLocker) ForceRelease(context.Context) error                 { return nil }
func (l NoopLocker) LockedOwner(context.Context) (Owner, error)         { return nil, nil }
