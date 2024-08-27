package logging

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"time"
)

// Redacter defines how to log an object
type Redacter interface {
	Redact() string
}

// Server returns a middleware that logs slow requests from the server.
func Server(logger log.Logger, threshold time.Duration) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			var (
				code      int32
				reason    string
				kind      string
				operation string
			)
			startTime := time.Now()
			duration := time.Since(startTime)
			if duration > threshold {
				if info, ok := transport.FromServerContext(ctx); ok {
					kind = info.Kind().String()
					operation = info.Operation()

				}
				reply, err = handler(ctx, req)
				if se := errors.FromError(err); se != nil {
					code = se.Code
					reason = se.Reason
				}
				level, stack := extractError(err)
				log.NewHelper(log.WithContext(ctx, logger)).Log(level,
					"kind", "server",
					"component", kind,
					"operation", operation,
					"args", extractArgs(req),
					"code", code,
					"reason", reason,
					"stack", stack,
					"latency", time.Since(startTime).Seconds(),
				)
			}
			return
		}
	}
}

// Client returns a middleware that logs slow requests from the client.
func Client(logger log.Logger, threshold time.Duration) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			var (
				code      int32
				reason    string
				kind      string
				operation string
			)
			startTime := time.Now()
			duration := time.Since(startTime)
			if duration > threshold {
				if info, ok := transport.FromClientContext(ctx); ok {
					kind = info.Kind().String()
					operation = info.Operation()
				}
				reply, err = handler(ctx, req)
				if se := errors.FromError(err); se != nil {
					code = se.Code
					reason = se.Reason
				}
				level, stack := extractError(err)
				log.NewHelper(log.WithContext(ctx, logger)).Log(level,
					"kind", "client",
					"component", kind,
					"operation", operation,
					"args", extractArgs(req),
					"code", code,
					"reason", reason,
					"stack", stack,
					"latency", time.Since(startTime).Seconds(),
				)
			}
			return
		}
	}
}

// extractError returns the string of the error
func extractError(err error) (log.Level, string) {
	if err != nil {
		return log.LevelError, fmt.Sprintf("%+v", err)
	}
	return log.LevelInfo, ""
}

// extractArgs returns the string of the req
func extractArgs(req interface{}) string {
	if redacter, ok := req.(Redacter); ok {
		return redacter.Redact()
	}
	if stringer, ok := req.(fmt.Stringer); ok {
		return stringer.String()
	}
	return fmt.Sprintf("%+v", req)
}
