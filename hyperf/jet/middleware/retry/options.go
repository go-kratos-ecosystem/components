package retry

import "github.com/go-kratos-ecosystem/components/v2/hyperf/jet"

var DefaultAllow = jet.IsHTTPTransporterServerError

type AllowFunc func(err error) bool

func AllowChain(fs ...AllowFunc) AllowFunc {
	return func(err error) bool {
		for _, f := range fs {
			if f(err) {
				return true
			}
		}
		return false
	}
}

type options struct {
	attempts int
	backoff  BackoffFunc
	allow    AllowFunc // allow retry
}

type Option func(o *options)

func Attempts(attempts int) Option {
	return func(o *options) {
		o.attempts = attempts
	}
}

func Allow(f AllowFunc) Option {
	return func(o *options) {
		o.allow = f
	}
}

func Backoff(f BackoffFunc) Option {
	return func(o *options) {
		o.backoff = f
	}
}
