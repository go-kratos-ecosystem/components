package storage

type Storage interface {
	// Get retrieves the value for the given path.
	Get(path string) ([]byte, error)

	// Set sets the value for the given path.
	Set(path string, value []byte) error

	// Delete deletes the value for the given path.
	Delete(path string) error

	// Has checks if the path exists.
	Has(path string) (bool, error)

	// Prepend prepends the value to the existing value for the given path.
	Prepend(path string, value []byte) error

	// Append appends the value to the existing value for the given path.
	Append(path string, value []byte) error

	// Move moves the value from the old path to the new path.
	Move(oldPath, newPath string) error

	// Copy copies the value from the old path to the new path.
	Copy(oldPath, newPath string) error

	// Link creates a hard link from the old path to the new path.
	Link(oldPath, newPath string) error

	// Symlink creates a symbolic link from the old path to the new path.
	Symlink(oldPath, newPath string) error

	// Files lists the files in the given path.
	Files(path string) ([]string, error)

	// AllFiles lists all the files in the given path.(including subdirectories)
	AllFiles(path string) ([]string, error)

	// Directories lists the directories in the given path.
	Directories(path string) ([]string, error)

	// AllDirectories lists all the directories in the given path.(including subdirectories)
	AllDirectories(path string) ([]string, error)

	// IsFile checks if the path is a file.
	IsFile(path string) (bool, error)

	// IsDir checks if the path is a directory.
	IsDir(path string) (bool, error)
}
