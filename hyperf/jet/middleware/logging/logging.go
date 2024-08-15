package logging

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/go-kratos-ecosystem/components/v2/hyperf/jet"
)

type options struct {
	logger log.Logger
}

type Option func(*options)

func Logger(logger log.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

func New(opts ...Option) jet.Middleware {
	o := options{
		logger: log.DefaultLogger,
	}
	for _, opt := range opts {
		opt(&o)
	}
	return func(next jet.Handler) jet.Handler {
		return func(ctx context.Context, client *jet.Client, name string, request any) (response any, err error) {
			defer func(starting time.Time) {
				level := log.LevelInfo
				if err != nil {
					level = log.LevelError
				}

				service := "unknown"
				if client != nil {
					service = client.GetService()
				}

				_ = log.WithContext(ctx, o.logger).Log(level,
					"kind", "jet",
					"service", service,
					"name", name,
					"request", request,
					"response", response,
					"error", err,
					"latency", time.Since(starting),
				)
			}(time.Now())
			return next(ctx, client, name, request)
		}
	}
}
