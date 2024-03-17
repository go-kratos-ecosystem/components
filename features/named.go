package features

type Named interface {
	Name() string
}

type NamedFeature struct {
	name string
}

func NewNamedFeature(name string) *NamedFeature {
	return &NamedFeature{name: name}
}

func (n *NamedFeature) Name() string {
	return n.name
}
