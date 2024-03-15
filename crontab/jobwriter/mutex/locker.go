package mutex

import "time"

type Locker interface {
	Lock(name string, expired time.Duration) error
	Unlock(name string) error
}
