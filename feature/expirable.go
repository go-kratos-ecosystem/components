package feature

import "time"

type Expirable interface {
	Expiration() time.Duration
}
