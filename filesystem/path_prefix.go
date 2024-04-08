package filesystem

import "strings"

type PathPrefix struct {
	prefix string
}

func NewPathPrefix(prefix string) *PathPrefix {
	prefix = strings.TrimRight(prefix, "/") + "/"
	return &PathPrefix{
		prefix: prefix,
	}
}

func (p *PathPrefix) GetPrefix() string {
	return p.prefix
}

func (p *PathPrefix) ApplePrefixPath(path string) string {
	return p.prefix + path
}
