package task

type options struct {
	size       int // task queue size
	goroutines int // max goroutines
}

func newOptions(opts ...Option) *options {
	o := &options{}
	o.apply(opts...)
	o.init()
	return o
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

func (o *options) init() {
	if o.size <= 0 {
		o.size = 100
	}

	if o.goroutines <= 0 {
		o.goroutines = 100
	}
}

type Option func(*options)

func Size(size int) Option {
	return func(o *options) {
		o.size = size
	}
}

func Goroutines(goroutines int) Option {
	return func(o *options) {
		o.goroutines = goroutines
	}
}
