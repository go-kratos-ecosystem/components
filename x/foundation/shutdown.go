package foundation

import (
	"context"
	"time"
)

const ShutdownName = "foundation.shutdown"

type ShutdownEvent struct {
	Time time.Time
}

func NewShutdown(opts ...Option) func(context.Context) error {
	o := &options{}

	for _, opt := range opts {
		opt(o)
	}

	return func(context.Context) error {
		if o.m != nil {
			o.m.Close(ShutdownName)
		}

		if o.d != nil {
			o.d.Dispatch(ShutdownName, &ShutdownEvent{
				Time: time.Now(),
			})
		}

		return nil
	}
}
