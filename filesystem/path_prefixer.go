package filesystem

import "strings"

type PathPrefixer struct {
	prefix    string
	separator string
}

type PathPrefixerOption func(*PathPrefixer)

func WithPathPrefixerSeparator(separator string) PathPrefixerOption {
	return func(p *PathPrefixer) {
		p.separator = separator
	}
}

func NewPathPrefixer(prefix string, opts ...PathPrefixerOption) *PathPrefixer {
	fixer := &PathPrefixer{
		prefix:    strings.TrimRight(prefix, "/"),
		separator: "/",
	}

	for _, opt := range opts {
		opt(fixer)
	}

	return fixer
}

func (p *PathPrefixer) PrefixPath(path string) string {
	return p.prefix + p.separator + strings.TrimLeft(path, "/")
}
