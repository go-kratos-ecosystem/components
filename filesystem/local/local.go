package local

import (
	"os"

	"github.com/go-kratos-ecosystem/components/v2/filesystem"
)

type local struct {
	root     string
	prefixer *filesystem.PathPrefixer
}

type Option func(*local)

func WithPrefixer(prefixer *filesystem.PathPrefixer) Option {
	return func(l *local) {
		l.prefixer = prefixer
	}
}

func New(root string, opts ...Option) filesystem.Filesystem {
	l := &local{
		root: root,
	}

	for _, opt := range opts {
		opt(l)
	}

	if l.prefixer == nil {
		l.prefixer = filesystem.NewPathPrefixer(root)
	}

	return l
}

func (l *local) Exists(path string) (bool, error) {
	path = l.prefixer.PrefixPath(path)

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (l *local) Get(path string) ([]byte, error) {
	return os.ReadFile(
		l.prefixer.PrefixPath(path),
	)
}

func (l *local) Put(path string, value []byte) error {
	path = l.prefixer.PrefixPath(path)
	return os.WriteFile(path, value, 0644) //nolint:gofumpt
}

func (l *local) Prepend(path string, value []byte) error {
	// TODO implement me
	panic("implement me")
}

func (l *local) Append(path string, value []byte) error {
	// TODO implement me
	panic("implement me")
}

func (l *local) Delete(path string) error {
	// TODO implement me
	panic("implement me")
}

func (l *local) Copy(src, dst string) error {
	// TODO implement me
	panic("implement me")
}

func (l *local) Move(src, dst string) error {
	return os.Rename(
		l.prefixer.PrefixPath(src),
		l.prefixer.PrefixPath(dst),
	)
}

func (l *local) Size(path string) (int64, error) {
	// TODO implement me
	panic("implement me")
}
