package storage

import (
	"context"
	"errors"
)

var ErrNotSupported = errors.New("[storage] not supported")

type Storage interface {
	// Get retrieves the value for the given path.
	Get(ctx context.Context, path string) ([]byte, error)

	// Set sets the value for the given path.
	Set(ctx context.Context, path string, value []byte) error

	// Delete deletes the value for the given path.
	Delete(ctx context.Context, path string) error

	// Has checks if the path exists.
	Has(ctx context.Context, path string) (bool, error)

	// Move moves the value from the old path to the new path.
	Move(ctx context.Context, oldPath, newPath string) error

	// Link creates a hard link from the old path to the new path.
	Link(ctx context.Context, oldPath, newPath string) error

	// Symlink creates a symbolic link from the old path to the new path.
	Symlink(ctx context.Context, oldPath, newPath string) error

	// Files lists the files in the given path.
	Files(ctx context.Context, path string) ([]string, error)

	// AllFiles lists all the files in the given path.(including subdirectories)
	AllFiles(ctx context.Context, path string) ([]string, error)

	// Directories lists the directories in the given path.
	Directories(ctx context.Context, path string) ([]string, error)

	// AllDirectories lists all the directories in the given path.(including subdirectories)
	AllDirectories(ctx context.Context, path string) ([]string, error)

	// IsFile checks if the path is a file.
	IsFile(ctx context.Context, path string) (bool, error)

	// IsDir checks if the path is a directory.
	IsDir(ctx context.Context, path string) (bool, error)

	// Name(ctx context.Context, path string) (string, error)
	//
	// Basename(ctx context.Context, path string) (string, error)
}

type Copyable interface {
	Copy(ctx context.Context, oldPath, newPath string) error
}

type noopStorage struct{}

var NoopStorage Storage = noopStorage{}

func (noopStorage) Get(context.Context, string) ([]byte, error)              { return nil, nil }
func (noopStorage) Set(context.Context, string, []byte) error                { return nil }
func (noopStorage) Delete(context.Context, string) error                     { return nil }
func (noopStorage) Has(context.Context, string) (bool, error)                { return true, nil }
func (noopStorage) Move(context.Context, string, string) error               { return nil }
func (noopStorage) Link(context.Context, string, string) error               { return nil }
func (noopStorage) Symlink(context.Context, string, string) error            { return nil }
func (noopStorage) Files(context.Context, string) ([]string, error)          { return nil, nil }
func (noopStorage) AllFiles(context.Context, string) ([]string, error)       { return nil, nil }
func (noopStorage) Directories(context.Context, string) ([]string, error)    { return nil, nil }
func (noopStorage) AllDirectories(context.Context, string) ([]string, error) { return nil, nil }
func (noopStorage) IsFile(context.Context, string) (bool, error)             { return false, nil }
func (noopStorage) IsDir(context.Context, string) (bool, error)              { return false, nil }
