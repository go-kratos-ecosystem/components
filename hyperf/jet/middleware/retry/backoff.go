package retry

import (
	"time"
)

var DefaultBackoff = LinearBackoff(100 * time.Millisecond)

type BackoffFunc func(attempt int) time.Duration

func NoBackoff() BackoffFunc {
	return func(_ int) time.Duration {
		return 0
	}
}

// LinearBackoff returns a backoff function that increases the delay linearly.
func LinearBackoff(delay time.Duration) BackoffFunc {
	return func(attempt int) time.Duration {
		return delay * time.Duration(attempt)
	}
}

// ExponentialBackoff returns a backoff function that increases the delay exponentially.
func ExponentialBackoff(delay time.Duration) BackoffFunc {
	return func(attempt int) time.Duration {
		return delay * time.Duration(1<<uint(attempt))
	}
}

// ConstantBackoff returns a backoff function that always returns the same delay.
func ConstantBackoff(delay time.Duration) BackoffFunc {
	return func(_ int) time.Duration {
		return delay
	}
}
