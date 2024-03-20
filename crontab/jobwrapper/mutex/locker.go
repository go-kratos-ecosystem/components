package mutex

import (
	"errors"
	"time"
)

var ErrLocked = errors.New("crontab/jobwrapper/mutex: job is locked")

type Locker interface {
	Lock(slug string, expiration time.Duration) error
	Unlock(slug string) error
}
