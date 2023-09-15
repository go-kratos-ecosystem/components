package env

import "context"

type currentEnvKey struct{}

func NewContext(ctx context.Context, env Env) context.Context {
	return context.WithValue(ctx, currentEnvKey{}, env)
}

func FromContext(ctx context.Context) (Env, bool) {
	if env, ok := ctx.Value(currentEnvKey{}).(Env); ok {
		return env, true
	}

	return "", false
}
