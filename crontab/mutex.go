package crontab

import (
	"context"
	"errors"
)

var ErrAnotherServerRunning = errors.New("crontab: another server running")

type Mutex interface {
	Lock(ctx context.Context, name string) error
	Unlock(ctx context.Context, name string) error
}
