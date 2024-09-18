package tracing

import (
	"fmt"
	"runtime"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

type options struct {
	tracerName   string
	tp           trace.TracerProvider
	spanHandlers []SpanHandler
}

type Option func(*options)

func newOptions(opts ...Option) *options {
	o := &options{
		tracerName: "gin",
		tp:         otel.GetTracerProvider(),
		spanHandlers: []SpanHandler{
			baseHandler{},
		},
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func WithTracerName(tracerName string) Option {
	return func(o *options) {
		o.tracerName = tracerName
	}
}

func WithTracerProvider(tp trace.TracerProvider) Option {
	return func(o *options) {
		o.tp = tp
	}
}

func WithSpanHandlers(handlers ...SpanHandler) Option {
	return func(o *options) {
		o.spanHandlers = append(o.spanHandlers, handlers...)
	}
}

// New returns a new tracing middleware.
func New(opts ...Option) gin.HandlerFunc {
	o := newOptions(opts...)

	tracer := o.tp.Tracer(o.tracerName)

	return func(c *gin.Context) {
		ctx, span := tracer.Start(c.Request.Context(), c.Request.URL.Path)
		defer span.End()

		for _, h := range o.spanHandlers {
			h.StartSpan(ctx, span, c)
		}

		defer func() {
			if r := recover(); r != nil {
				// stack trace
				stackTrace := make([]byte, 2048)
				n := runtime.Stack(stackTrace, false)

				span.SetAttributes(
					semconv.CodeStacktrace(string(stackTrace[:n])),
				)
				span.SetStatus(codes.Error, fmt.Sprintf("%v", r))
				span.RecordError(fmt.Errorf("%v", r))

				span.End() // if panic, end the span

				panic(r)
			}

			span.SetStatus(codes.Ok, "OK")
		}()

		defer func() {
			for _, h := range o.spanHandlers {
				h.EndSpan(ctx, span, c)
			}
		}()

		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
