package features

type Slug interface {
	Slug() string
}

type SlugFeature struct {
	slug string
}

func NewSlugFeature(slug string) *SlugFeature {
	return &SlugFeature{slug: slug}
}

func (n *SlugFeature) Slug() string {
	return n.slug
}
