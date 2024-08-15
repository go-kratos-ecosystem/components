package jet

import "context"

type Handler func(ctx context.Context, client *Client, name string, request any) (response any, err error)

type Middleware func(Handler) Handler

// Chain chains the middlewares.
//
//	Chain(m1, m2, m3)(xxx) => m1(m2(m3(xxx))
func Chain(m ...Middleware) Middleware {
	return func(next Handler) Handler {
		for i := len(m) - 1; i >= 0; i-- {
			next = m[i](next)
		}
		return next
	}
}
