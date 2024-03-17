package features

import "time"

type Expirable interface {
	Expiration() time.Duration
}

type ExpireFeature struct {
	expiration time.Duration
}

func NewExpireFeature(expiration time.Duration) *ExpireFeature {
	return &ExpireFeature{expiration: expiration}
}

func (e *ExpireFeature) Expiration() time.Duration {
	return e.expiration
}
