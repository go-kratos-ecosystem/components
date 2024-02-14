package v2 //nolint:typecheck

type Job interface {
	Exp() string // Expression
	Run()        // Handler
}
