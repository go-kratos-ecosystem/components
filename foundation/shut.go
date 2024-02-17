package foundation

import (
	"context"
	"time"
)

const ShutName = "foundation.shut"

type ShutEvent struct {
	Time time.Time
}

func NewShut(opts ...Option) func(context.Context) error {
	opt := options{}

	for _, o := range opts {
		o(&opt)
	}

	return func(context.Context) error {
		if opt.m != nil {
			opt.m.Close(ShutName)
		}

		if opt.d != nil {
			opt.d.Dispatch(ShutName, &ShutEvent{
				Time: time.Now(),
			})
		}

		return nil
	}
}
