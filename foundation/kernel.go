package foundation

type Kernel struct {
	app       Application
	providers []Provider
}

func NewKernel(app Application) *Kernel {
	return &Kernel{
		app: app,
	}
}

func (k *Kernel) Register(providers ...Provider) {
	k.providers = append(k.providers, providers...)
}

func (k *Kernel) bootstrap() error {
	for _, provider := range k.providers {
		if err := provider.Bootstrap(k.app); err != nil {
			return err
		}
	}

	return nil
}

func (k *Kernel) Run() error {
	if err := k.bootstrap(); err != nil {
		return err
	}
	defer k.terminate() // nolint:errcheck

	return k.app.Run()
}

func (k *Kernel) terminate() error {
	for _, provider := range k.providers {
		if err := provider.Terminate(k.app); err != nil {
			return err
		}
	}

	return nil
}
