package recovery

type options struct {
	handler func(err interface{})
}

type Option func(o *options)

func WithHandler(handler func(err interface{})) Option {
	return func(o *options) {
		o.handler = handler
	}
}

type Recovery struct {
	opt *options
}

func New(opts ...Option) *Recovery {
	o := &options{}

	for _, opt := range opts {
		opt(o)
	}

	if o.handler == nil {
		o.handler = func(err interface{}) {
			panic(err)
		}
	}

	return &Recovery{
		opt: o,
	}
}

func (r *Recovery) Wrap(f func() error) error {
	defer func() {
		if err := recover(); err != nil {
			if r.opt.handler != nil {
				r.opt.handler(err)
			}
		}
	}()

	return f()
}
