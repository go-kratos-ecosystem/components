package retry

import "time"

type BackoffFunc func(attempt int) time.Duration

func NoBackoff(_ int) time.Duration {
	return 0
}

// LinearBackoff returns a backoff function that increases the delay linearly.
// todo: 待确定逻辑
func LinearBackoff(delay time.Duration) BackoffFunc {
	return func(attempt int) time.Duration {
		return delay
	}
}

// todo: 待确定逻辑
func ExponentialBackoff(delay time.Duration) BackoffFunc {
	return func(attempt int) time.Duration {
		return delay * time.Duration(attempt)
	}
}
