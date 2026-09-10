package debounce

type config struct {
	onError func(error)
}

// Option configures a [Debouncer].
type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (f optionFunc) apply(c *config) {
	f(c)
}

// WithOnError sets a callback invoked whenever the handler returns an error.
// A nil callback disables error notifications, which is the default.
//
// The callback runs synchronously on the worker, including for calls initiated
// by [Debouncer.Flush]. Flush and [Debouncer.Close] wait for it to return.
// It may call [Debouncer.Trigger] but must not synchronously call Flush or Close
// on its own debouncer. Panics are not recovered.
func WithOnError(fn func(error)) Option {
	return optionFunc(func(c *config) { c.onError = fn })
}

func newConfig(options ...Option) config {
	var c config
	for _, option := range options {
		if option != nil {
			option.apply(&c)
		}
	}

	return c
}
