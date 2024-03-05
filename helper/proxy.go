package helper

type Proxy[T any] struct {
	target T
}

func NewProxy[T any](target T) *Proxy[T] {
	return &Proxy[T]{
		target: target,
	}
}

func (p *Proxy[T]) Tap(callbacks ...func(T)) T {
	for _, callback := range callbacks {
		if callback != nil {
			callback(p.target)
		}
	}

	return p.target
}

func (p *Proxy[T]) With(callbacks ...func(T) T) T {
	for _, callback := range callbacks {
		if callback != nil {
			p.target = callback(p.target)
		}
	}

	return p.target
}

func (p *Proxy[T]) When(condition bool, callbacks ...func(T) T) T {
	if condition {
		return p.With(callbacks...)
	}

	return p.target
}

func (p *Proxy[T]) Unless(condition bool, callbacks ...func(T) T) T {
	return p.When(!condition, callbacks...)
}

func (p *Proxy[T]) Target() T {
	return p.target
}
