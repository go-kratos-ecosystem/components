package foundation

type Provider interface {
	Bootstrap(Application) error
	Terminate(Application) error
}
