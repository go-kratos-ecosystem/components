package recovery

type Recovery struct {
	handler func(err any)
}

type Option func(o *Recovery)

func WithHandler(handler func(err any)) Option {
	return func(r *Recovery) {
		r.handler = handler
	}
}

func New(opts ...Option) *Recovery {
	r := &Recovery{
		handler: func(err any) {
			panic(err)
		},
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (r *Recovery) Wrap(f func()) {
	if r.handler != nil {
		defer func() {
			if err := recover(); err != nil {
				r.handler(err)
			}
		}()
	}

	f()
}
