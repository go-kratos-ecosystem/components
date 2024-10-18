package locker

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
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

type Owner interface {
	// Name returns the name of the owner.
	Name() string

	// Release releases the lock.
	Release(ctx context.Context) (bool, error)
}

type owner struct {
	name   string
	locker Locker
}

type OwnerOption func(*owner)

func WithName(name string) OwnerOption {
	return func(o *owner) {
		o.name = name
	}
}

func NewOwner(locker Locker, opts ...OwnerOption) Owner {
	o := &owner{
		name:   uuid.New().String(),
		locker: locker,
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func (o *owner) Name() string {
	return o.name
}

func (o *owner) Release(ctx context.Context) (bool, error) {
	return o.locker.Release(ctx, o)
}
