# OTLP Configuration

This package provides a configuration for the OpenTelemetry Protocol (OTLP) exporter.

## Usage Example

```go
package main

import (
	"context"

	"go.opentelemetry.io/otel/attribute"

	"github.com/go-kratos-ecosystem/components/v2/otel/otlp"
)

func main() {
	ctx := context.TODO()

	// transport
	transport := otlp.NewGRPCTransport("localhost:4317", otlp.WithGRPCTransportInsecure(true))

	// client
	client := otlp.NewClient(
		otlp.WithServiceName("service-name"),
		otlp.WithDeploymentEnvironment("development"),
		otlp.WithAttributes(
			attribute.String("key", "value"),
			// ...
		),
		otlp.WithTransport(transport),
	)

	if err := client.Configure(ctx); err != nil {
		panic(err)
	}

	defer client.Shutdown(ctx)

	// do something
}

```