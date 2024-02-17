package foundation

import (
	"context"
	"time"
)

const BootstrapName = "foundation.bootstrap"

type BootstrapEvent struct {
	Time time.Time
}

func NewBootstrap(opts ...Option) func(context.Context) error {
	o := &options{}

	for _, opt := range opts {
		opt(o)
	}

	return func(context.Context) error {
		if o.m != nil {
			o.m.Close(BootstrapName)
		}

		if o.d != nil {
			o.d.Dispatch(BootstrapName, &BootstrapEvent{
				Time: time.Now(),
			})
		}

		return nil
	}
}
