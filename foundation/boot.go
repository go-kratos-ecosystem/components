package foundation

import (
	"context"
	"time"
)

const BootName = "foundation.boot"

type BootEvent struct {
	Time time.Time
}

func NewBoot(opts ...Option) func(context.Context) error {
	opt := options{}

	for _, o := range opts {
		o(&opt)
	}

	return func(context.Context) error {
		if opt.m != nil {
			opt.m.Close(BootName)
		}

		if opt.d != nil {
			opt.d.Dispatch(BootName, &BootEvent{
				Time: time.Now(),
			})
		}

		return nil
	}
}
