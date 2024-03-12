package env

import "context"

func Provider(env Env) func(ctx context.Context) (context.Context, error) {
	return func(ctx context.Context) (context.Context, error) {
		SetEnv(env)

		return NewContext(ctx, env), nil
	}
}
