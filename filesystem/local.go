package filesystem

import "os"

type Local struct {
	root       string
	pathPrefix *PathPrefix
}

func NewLocal(root string) *Local {
	return &Local{
		root:       root,
		pathPrefix: NewPathPrefix(root),
	}
}

func (l *Local) Read(filename string) ([]byte, error) {
	return os.ReadFile(l.pathPrefix.ApplePrefixPath(filename))
}

func (l *Local) Write(filename string, data []byte) error {
	return os.WriteFile(l.pathPrefix.ApplePrefixPath(filename), data, 0o644) //nolint:gomnd
}

func (l *Local) Exists(filename string) bool {
	_, err := os.Stat(l.pathPrefix.ApplePrefixPath(filename))
	return err == nil
}

func (l *Local) Remove(filename string) error {
	return os.Remove(l.pathPrefix.ApplePrefixPath(filename))
}

func (l *Local) MkdirAll(path string) error {
	return os.MkdirAll(l.pathPrefix.ApplePrefixPath(path), 0o755) //nolint:gomnd
}

func (l *Local) List(path string) ([]string, error) {
	// TODO implement me
	panic("implement me")
}

func (l *Local) Size(filename string) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (l *Local) MimeType(filename string) (string, error) {
	// TODO implement me
	panic("implement me")
}

func (l *Local) Rename(oldpath, newpath string) error {
	// TODO implement me
	panic("implement me")
}

func (l *Local) Close() error {
}
