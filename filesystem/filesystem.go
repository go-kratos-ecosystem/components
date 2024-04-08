package filesystem

type Filesystem interface {
	// Read reads the file named by filename and returns the contents.
	Read(filename string) ([]byte, error)

	// Write writes data to a file named by filename.
	Write(filename string, data []byte) error

	// Exists reports whether the named file or directory exists.
	Exists(filename string) bool

	// Remove removes the named file or directory.
	Remove(filename string) error

	// MkdirAll creates a directory named path, along with any necessary parents.
	MkdirAll(path string) error

	// List returns the names of the files inside the directory.
	List(path string) ([]string, error)

	// Size returns the size in bytes of the named file.
	Size(filename string) (int64, error)

	// MimeType returns the MIME type of the file.
	MimeType(filename string) (string, error)

	// Rename renames (moves) oldpath to newpath.
	Rename(oldpath, newpath string) error

	// Close closes the filesystem.
	Close() error
}
