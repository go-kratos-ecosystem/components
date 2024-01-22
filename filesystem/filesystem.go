package filesystem

type Filesystem interface {
	Exists(path string) (bool, error)
	Get(path string) ([]byte, error)
	Put(path string, value []byte) error
	Prepend(path string, value []byte) error
	Append(path string, value []byte) error
	Delete(path string) error
	Copy(src, dst string) error
	Move(src, dst string) error
	Size(path string) (int64, error)
}
