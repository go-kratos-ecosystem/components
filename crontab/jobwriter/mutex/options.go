package mutex

import (
	"time"

	"github.com/robfig/cron/v3"
)

type options struct {
	prefix     string
	locker     Locker
	logger     cron.Logger
	expiration time.Duration
}

func newOptions(opts ...Option) *options {
	o := &options{}
	o.apply(opts...)
	o.applyDefault()
	return o
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

func (o *options) applyDefault() {
	if o.prefix == "" {
		o.prefix = "crontab:"
	}

	if o.locker == nil {
		panic("locker is required")
	}

	if o.logger == nil {
		o.logger = cron.DefaultLogger
	}

	if o.expiration == 0 {
		o.expiration = time.Minute * 60 //nolint:gomnd
	}
}

type Option func(*options)

func WithPrefix(prefix string) Option {
	return func(o *options) {
		o.prefix = prefix
	}
}

func WithLocker(locker Locker) Option {
	return func(o *options) {
		o.locker = locker
	}
}

func WithLogger(logger cron.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

func WithExpiration(expiration time.Duration) Option {
	return func(o *options) {
		o.expiration = expiration
	}
}
