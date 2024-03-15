package mutex

import (
	"errors"
	"time"
)

var ErrLocked = errors.New("crontab/jobwrapper/mutex: job is locked")

type Locker interface {
	Lock(name string, expiration time.Duration) error
	Unlock(name string) error
}
