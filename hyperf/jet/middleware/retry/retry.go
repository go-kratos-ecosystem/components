package retry

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos-ecosystem/components/v2/hyperf/jet"
)

var defaultCanRetry = func(err error) bool {
	if jet.IsHTTPTransporterServerError(err) {
		return true
	}
	return false
}

type options struct {
	attempts int
	backoff  BackoffFunc
	canRetry func(err error) bool
}

type Option func(o *options)

func Attempts(attempts int) Option {
	return func(o *options) {
		o.attempts = attempts
	}
}

func CanRetry(f func(err error) bool) Option {
	return func(o *options) {
		o.canRetry = f
	}
}

func Backoff(f BackoffFunc) Option {
	return func(o *options) {
		o.backoff = f
	}
}

func NewRetry(opts ...Option) jet.Middleware {
	o := options{
		attempts: 3,
		backoff:  LinearBackoff(100 * time.Millisecond),
		canRetry: defaultCanRetry,
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

				if !o.canRetry(err) {
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
