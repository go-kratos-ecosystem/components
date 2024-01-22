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
		root:     root,
		prefixer: filesystem.NewPathPrefixer(root),
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

func (l *local) Exists(path string) (bool, error) {
	path = l.prefixer.PrefixPath(path)

	if _, err := os.Stat(path); err != nil {
		return false, err
	} else if os.IsNotExist(err) {
		return false, nil
	}

	return true, nil
}

func (l *local) Get(path string) ([]byte, error) {
	path = l.prefixer.PrefixPath(path)

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	data := make([]byte, stat.Size())
	_, err = file.Read(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (l *local) Put(path string, value []byte) error {
	path = l.prefixer.PrefixPath(path)

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(value)
	if err != nil {
		return err
	}

	return nil
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
	// TODO implement me
	panic("implement me")
}

func (l *local) Size(path string) (int64, error) {
	// TODO implement me
	panic("implement me")
}
