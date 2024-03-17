package features

type Retriable interface {
	Retries() int
}

type RetryFeature struct {
	retries int
}

func NewRetryFeature(retries int) *RetryFeature {
	return &RetryFeature{retries: retries}
}

func (r *RetryFeature) Retries() int {
	return r.retries
}
