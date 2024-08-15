package retry

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-kratos-ecosystem/components/v2/hyperf/jet"
)

func NewRetry(opts ...Option) jet.Middleware {
	o := options{
		attempts: 3,
		backoff:  DefaultBackoff,
		allow:    DefaultAllow,
	}
	for _, opt := range opts {
		opt(&o)
	}
	return func(next jet.Handler) jet.Handler {
		return func(ctx context.Context, name string, request any) (response any, err error) {
			starting := time.Now()
			for i := 1; i <= o.attempts; i++ {
				response, err = next(ctx, name, request)
				if err == nil {
					return
				}

				if !o.allow(err) {
					return
				}

				if i < o.attempts {
					time.Sleep(o.backoff(i))
				}
			}
			return response, &Error{
				Attempts: o.attempts,
				Time:     time.Since(starting),
				Err:      err,
			}
		}
	}
}

type Error struct {
	Attempts int
	Time     time.Duration
	Err      error
}

func (e *Error) Error() string {
	return fmt.Sprintf("jet/middleware/retry: retry failed after %d attempts, time: %v, last error: %v", e.Attempts, e.Time, e.Err)
}

func (e *Error) Unwrap() error {
	return e.Err
}

func IsError(err error) bool {
	var e *Error
	ok := errors.As(err, &e)
	return ok
}
