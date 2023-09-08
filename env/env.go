package env

type Env string

var (
	Dev   Env = "dev"
	Prod  Env = "prod"
	Debug Env = "debug"
	Stage Env = "stage"
)

var currentEnv = Prod

func SetEnv(env Env) {
	currentEnv = env
}

func Is(envs ...Env) bool {
	for _, env := range envs {
		if currentEnv == env {
			return true
		}
	}

	return false
}

func IsDev() bool {
	return Is(Dev)
}

func IsProd() bool {
	return Is(Prod)
}

func IsDebug() bool {
	return Is(Debug)
}

func IsStage() bool {
	return Is(Stage)
}

func IsUseString(envs ...string) bool {
	for _, env := range envs {
		if currentEnv == Env(env) {
			return true
		}
	}

	return false
}
