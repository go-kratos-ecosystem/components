package env

import "context"

func Provider(env Env) func(ctx context.Context) error {
	return func(context.Context) error {
		SetEnv(env)

		return nil
	}
}
