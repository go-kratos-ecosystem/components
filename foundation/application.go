package foundation

type Application interface {
	Register(providers ...Provider)
	Boot() error
	Run() error
}
