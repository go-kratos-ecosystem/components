package locker

import (
	"context"
	"errors"
	"time"
)

var (
	ErrLocked    = errors.New("locker: locked")
	ErrTimeout   = errors.New("locker: timeout")
	ErrNotLocked = errors.New("locker: not locked")
)

type Locker interface {
	// Try tries to acquire the lock and execute the function.
	// If the lock is already acquired, it will return ErrLocked.
	// Otherwise, it will release the lock after the function is executed, and return the error of the function.
	Try(ctx context.Context, fn func() error) error

	// Until tries to acquire the lock and execute the function.
	// If the lock is already acquired, it will wait until the lock is released or the timeout is reached.
	// If the timeout is reached, it will return ErrTimeout.
	// Otherwise, it will release the lock after the function is executed, and return the error of the function.
	Until(ctx context.Context, timeout time.Duration, fn func() error) error

	// Release releases the lock for the current owner.
	// It returns true if the lock is released successfully, otherwise it returns false.
	// If the lock acquired by another owner, it will return false.
	Release(ctx context.Context) (bool, error)

	// ForceRelease releases the lock forcefully.
	ForceRelease(ctx context.Context) error

	// Owner return the owner of the lock.
	Owner() string

	// LockedOwner returns the current owner of the lock.
	// If the lock is not acquired, it will return an empty string and ErrNotLocked.
	// Unlike Owner, which returns the owner of the current node in the cluster,
	// LockedOwner returns the owner of the lock acquired by the Cluster.
	LockedOwner(ctx context.Context) (string, error)
}
