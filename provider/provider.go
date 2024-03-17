package provider

import "context"

type Provider func(ctx context.Context) (context.Context, error)
