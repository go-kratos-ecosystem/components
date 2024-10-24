package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/go-kratos-ecosystem/components/v2/hyperf/jet"
)

const instrumentation = "github.com/go-kratos-ecosystem/components/v2/hyperf/jet/middleware/tracing"

type options struct {
	mp    propagation.TextMapPropagator
	tp    trace.TracerProvider
	attrs []attribute.KeyValue
}

type Option func(*options)

func WithPropagator(mp propagation.TextMapPropagator) Option {
	return func(o *options) {
		o.mp = mp
	}
}

func WithTracerProvider(tp trace.TracerProvider) Option {
	return func(o *options) {
		o.tp = tp
	}
}

func WithAttributes(attrs ...attribute.KeyValue) Option {
	return func(o *options) {
		o.attrs = append(o.attrs, attrs...)
	}
}

func newOptions(opts ...Option) options {
	o := options{
		mp: otel.GetTextMapPropagator(),
		tp: otel.GetTracerProvider(),
	}
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

func New(opts ...Option) jet.Middleware {
	o := newOptions(opts...)

	tracer := o.tp.Tracer(instrumentation)
	return func(next jet.Handler) jet.Handler {
		return func(ctx context.Context, service, method string, request any) (response any, err error) {
			ctx, span := tracer.Start(ctx, service+"/"+method)
			defer span.End()

			// 10 is the number of attributes in the following code
			attrs := make([]attribute.KeyValue, 0, len(o.attrs)+10)

			attrs = []attribute.KeyValue{
				semconv.RPCSystemKey.String("jsonrpc"),
				semconv.RPCService(service),
				semconv.RPCMethod(method),
				semconv.RPCJsonrpcErrorCode(0),     // todo
				semconv.RPCJsonrpcErrorMessage(""), // todo
				semconv.RPCJsonrpcVersion(""),      // todo
				semconv.RPCJsonrpcRequestID(""),    // todo
				semconv.ServerAddress(""),          // todo
				semconv.ServerPort(0),              // todo
				semconv.NetworkPeerAddress(""),     // todo
				semconv.NetworkPeerPort(0),         // todo
			}
			attrs = append(attrs, o.attrs...)
			span.SetAttributes(attrs...)

			response, err = next(ctx, service, method, request)
			if err != nil {
				span.SetStatus(codes.Error, err.Error())
				span.RecordError(err)
			}

			return
		}
	}
}
