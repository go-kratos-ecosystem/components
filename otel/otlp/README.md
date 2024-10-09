# Otlp Configuration

This package provides a configuration for the OpenTelemetry Protocol (OTLP) exporter.

## Usage Example

```go
package main

import (
	"context"

	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	"github.com/go-kratos-ecosystem/components/v2/otel/otlp"
)

func main() {
	ctx := context.TODO()

	// resource
	res, err := sdkresource.New(ctx,
		sdkresource.WithHost(),
		sdkresource.WithTelemetrySDK(),
		sdkresource.WithContainer(),

		sdkresource.WithAttributes(
			semconv.ServiceName("service-name"),
			semconv.DeploymentEnvironment("prod"),
		),
	)
	if err != nil {
		panic(err)
	}

	// transport
	transport := otlp.NewGRPCTransport("localhost:4317", otlp.WithGRPCTransportInsecure(true))

	// client
	client := otlp.NewClient(
		otlp.WithResource(res),
		otlp.WithTransport(transport),
	)

	if err := client.Configure(ctx); err != nil {
		panic(err)
	}

	defer client.Shutdown(ctx)

	// do something
}
```