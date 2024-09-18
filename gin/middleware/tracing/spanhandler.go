package tracing

import (
	"context"

	"github.com/gin-gonic/gin"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

type SpanHandler interface {
	StartSpan(ctx context.Context, span trace.Span, c *gin.Context)
	EndSpan(ctx context.Context, span trace.Span, c *gin.Context)
}

type baseHandler struct{}

func (baseHandler) StartSpan(_ context.Context, span trace.Span, c *gin.Context) {
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
		semconv.ClientAddress(c.ClientIP()),
		semconv.NetworkProtocolName(c.Request.Proto),
	)
}

func (baseHandler) EndSpan(_ context.Context, span trace.Span, c *gin.Context) {
	span.SetAttributes(
		semconv.HTTPResponseStatusCode(c.Writer.Status()),
		semconv.HTTPResponseSize(c.Writer.Size()),
	)
}
