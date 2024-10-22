package locker

import (
	"context"
	"errors"
	"time"
)

var ErrNotLocked = errors.New("locker: the locker is not locked")

type Locker interface {
	// Try tries to get the lock.
	// the return value indicates whether the lock is locked successfully.
	Try(ctx context.Context, fn func()) (bool, error)

	// Until tries to get the lock until the timeout.
	// the return value indicates whether the lock is locked successfully.
	Until(ctx context.Context, timeout time.Duration, fn func()) (bool, error)

	// Get tries to get the lock.
	// the return the owner of the lock, whether the lock is locked, and an error.
	Get(ctx context.Context) (Owner, bool, error)

	// Release releases the locked owner.
	// the return value indicates whether the lock is released successfully.
	Release(ctx context.Context, owner Owner) (bool, error)

	// ForceRelease releases the lock forcibly.
	ForceRelease(ctx context.Context) error

	// LockedOwner returns the owner of the lock.
	LockedOwner(ctx context.Context) (Owner, error)
}

type NoopLocker struct{}

var _ Locker = NoopLocker{}

func (l NoopLocker) Try(context.Context, func()) (bool, error)                  { return true, nil }
func (l NoopLocker) Until(context.Context, time.Duration, func()) (bool, error) { return true, nil }
func (l NoopLocker) Get(context.Context) (Owner, bool, error)                   { return nil, true, nil }
func (l NoopLocker) Release(context.Context, Owner) (bool, error)               { return true, nil }
func (l NoopLocker) ForceRelease(context.Context) error                         { return nil }
func (l NoopLocker) LockedOwner(context.Context) (Owner, error)                 { return nil, nil }
