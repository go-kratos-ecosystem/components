package recovery

import (
	"context"
	"fmt"

	"github.com/go-kratos-ecosystem/components/v2/hyperf/jet"
)

var DefaultHandler = func(_ context.Context, client *jet.Client, name string, request any, err any) error {
	service := "unknown"
	if client != nil {
		service = client.GetService()
	}
	return fmt.Errorf("service: %s, name: %s, request: %v, error: %v", service, name, request, err)
}

type HandlerFunc func(ctx context.Context, client *jet.Client, name string, request any, err any) error

type options struct {
	handler HandlerFunc
}

type Option func(*options)

func Handler(h HandlerFunc) Option {
	return func(o *options) {
		o.handler = h
	}
}

func New(opts ...Option) jet.Middleware {
	return func(next jet.Handler) jet.Handler {
		o := options{
			handler: DefaultHandler,
		}
		for _, opt := range opts {
			opt(&o)
		}

		return func(ctx context.Context, client *jet.Client, name string, request any) (response any, err error) {
			defer func() {
				if rerr := recover(); rerr != nil {
					err = o.handler(ctx, client, name, request, rerr)
				}
			}()
			return next(ctx, client, name, request)
		}
	}
}
