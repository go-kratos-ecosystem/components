package request

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware"

	"github.com/go-kratos-ecosystem/components/v2/event"
)

const (
	BeforeName = "request.before"
	AfterName  = "request.after"
)

type From string

const (
	FromClient From = "client"
	FromServer From = "server"
)

type BeforeEvent struct {
	Ctx  context.Context
	Req  interface{}
	From From
}

type AfterEvent struct {
	Ctx   context.Context
	Req   interface{}
	Reply interface{}
	Err   error
	From  From
}

func Server(d *event.Dispatcher) middleware.Middleware {
	return newMiddleware(d, FromServer)
}

func Client(d *event.Dispatcher) middleware.Middleware {
	return newMiddleware(d, FromClient)
}

func newMiddleware(d *event.Dispatcher, from From) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if d != nil {
				d.Dispatch(BeforeName, &BeforeEvent{
					Ctx: ctx, Req: req, From: from,
				})
			}

			reply, err = handler(ctx, req)

			if d != nil {
				d.Dispatch(AfterName, &AfterEvent{
					Ctx: ctx, Req: req, Reply: reply, Err: err, From: from,
				})
			}

			return
		}
	}
}
