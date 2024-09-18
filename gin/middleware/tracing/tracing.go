package tracing

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

type options struct {
	tracerName string

	tp trace.TracerProvider
}

type Option func(*options)

func newOptions(opts ...Option) *options {
	o := &options{
		tp: otel.GetTracerProvider(),
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// New returns a new tracing middleware.
func New(opts ...Option) gin.HandlerFunc {
	o := newOptions(opts...)

	tracer := o.tp.Tracer(o.tracerName)

	return func(c *gin.Context) {
		ctx, span := tracer.Start(c.Request.Context(), c.Request.URL.Path)
		defer span.End()

		c.Request = c.Request.WithContext(ctx)

		span.SetAttributes(
			// see https://opentelemetry.io/docs/specs/semconv/http/http-spans/
			semconv.HTTPRequestMethodKey.String(c.Request.Method),
			semconv.URLPath(c.Request.URL.Path),
			semconv.URLScheme(c.Request.URL.Scheme),
			semconv.URLFragment(c.Request.URL.Fragment),
			semconv.HTTPRoute(c.FullPath()),
			semconv.URLQuery(c.Request.URL.RawQuery),
			semconv.ServerAddress(c.Request.Host),
			semconv.UserAgentOriginal(c.Request.UserAgent()),
		)

		defer func() {
			if r := recover(); r != nil {
				span.RecordError(r.(error))
				panic(r)
			}

			span.SetAttributes(
				semconv.HTTPResponseStatusCode(c.Writer.Status()),
				semconv.HTTPResponseSize(c.Writer.Size()),
			)
		}()

		c.Next()
	}
}
