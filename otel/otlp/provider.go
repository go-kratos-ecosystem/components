package otlp

import (
	"context"

	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	"github.com/go-kratos-ecosystem/components/v2/foundation"
)

type Provider struct {
	env     string
	service string

	client *Client
}

type ProviderOption func(*Provider)

func WithProviderEnv(env string) ProviderOption {
	return func(p *Provider) {
		p.env = env
	}
}

func WithProviderService(service string) ProviderOption {
	return func(p *Provider) {
		p.service = service
	}
}

var _ foundation.Provider = (*Provider)(nil)

func NewProvider(client *Client, opts ...ProviderOption) *Provider {
	provider := &Provider{
		client: client,
	}

	for _, opt := range opts {
		opt(provider)
	}

	return provider
}

func (p *Provider) Bootstrap(ctx context.Context) (context.Context, error) {
	// resource
	res, err := sdkresource.New(ctx,
		sdkresource.WithHost(),
		sdkresource.WithTelemetrySDK(),
		sdkresource.WithContainer(),
		sdkresource.WithAttributes(
			semconv.ServiceName(p.service),
			semconv.DeploymentEnvironment(p.env),
		),
	)
	if err != nil {
		return nil, err
	}

	p.client.RegisterResource(res)

	return ctx, p.client.Configure(ctx)
}

func (p *Provider) Terminate(ctx context.Context) (context.Context, error) {
	return ctx, p.client.Shutdown(ctx)
}
