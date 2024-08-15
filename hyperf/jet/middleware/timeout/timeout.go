package timeout

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos-ecosystem/components/v2/hyperf/jet"
)

var (
	ErrTimeout     = fmt.Errorf("jet/timeout: request timeout")
	defaultTimeout = time.Second * 5
)

type options struct {
	timeout time.Duration
}

type Option func(*options)

func Timeout(timeout time.Duration) Option {
	return func(o *options) {
		o.timeout = timeout
	}
}

func New(opts ...Option) jet.Middleware {
	o := options{
		timeout: defaultTimeout,
	}
	for _, opt := range opts {
		opt(&o)
	}
	return func(next jet.Handler) jet.Handler {
		return func(ctx context.Context, name string, request interface{}) (response interface{}, err error) {
			newCtx, cancel := context.WithTimeoutCause(ctx, o.timeout, ErrTimeout)
			defer cancel()

			finished := make(chan struct{}, 1)

			go func() {
				defer close(finished)
				response, err = next(newCtx, name, request)
			}()

			select {
			case <-newCtx.Done():
				if cause := context.Cause(newCtx); cause != nil {
					return nil, cause
				}
				return nil, newCtx.Err()
			case <-finished:
				return response, err
			}
		}
	}
}
