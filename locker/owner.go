package locker

import (
	"context"

	"github.com/google/uuid"
)

type Owner interface {
	// Name returns the name of the owner.
	Name() string

	// Release releases the lock.
	Release(ctx context.Context) error
}

type owner struct {
	name   string
	locker Locker
}

type OwnerOption func(*owner)

func WithOwnerName(name string) OwnerOption {
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

func (o *owner) Release(ctx context.Context) error {
	return o.locker.Release(ctx, o)
}
