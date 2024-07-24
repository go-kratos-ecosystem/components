package local

import (
	"context"
	"os"

	"github.com/go-kratos-ecosystem/components/v2/storage"
)

type Storage struct {
	root     string
	prefixer *storage.PathPrefixer
}

var _ storage.Storage = (*Storage)(nil)

func New(root string) *Storage {
	return &Storage{
		root:     root,
		prefixer: storage.NewPathPrefixer(root),
	}
}

func (s *Storage) Get(_ context.Context, path string) ([]byte, error) {
	return os.ReadFile(s.prefixer.Prefix(path))
}

func (s *Storage) Set(_ context.Context, path string, value []byte) error {
	return os.WriteFile(s.prefixer.Prefix(path), value, 0o644) //nolint:mnd
}

func (s *Storage) Delete(_ context.Context, path string) error {
	return os.Remove(s.prefixer.Prefix(path))
}

func (s *Storage) Has(_ context.Context, path string) (bool, error) {
	_, err := os.Stat(s.prefixer.Prefix(path))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (s *Storage) Move(_ context.Context, oldPath, newPath string) error {
	return os.Rename(s.prefixer.Prefix(oldPath), s.prefixer.Prefix(newPath))
}

func (s *Storage) Link(_ context.Context, oldPath, newPath string) error {
	return os.Link(s.prefixer.Prefix(oldPath), s.prefixer.Prefix(newPath))
}

func (s *Storage) Symlink(_ context.Context, oldPath, newPath string) error {
	return os.Symlink(s.prefixer.Prefix(oldPath), s.prefixer.Prefix(newPath))
}

func (s *Storage) Files(_ context.Context, path string) ([]string, error) {
	f, err := os.ReadDir(s.prefixer.Prefix(path))
	if err != nil {
		return nil, err
	}

	var files []string //nolint:prealloc
	for _, file := range f {
		if !file.IsDir() {
			files = append(files, file.Name())
		}
	}
	return files, nil
}

func (s *Storage) AllFiles(ctx context.Context, path string) ([]string, error) {
	f, err := os.ReadDir(s.prefixer.Prefix(path))
	if err != nil {
		return nil, err
	}

	var files []string //molint:prealloc
	for _, file := range f {
		if !file.IsDir() {
			files = append(files, file.Name())
		} else {
			subFiles, err := s.AllFiles(ctx, file.Name())
			if err != nil {
				return nil, err
			}
			files = append(files, subFiles...)
		}
	}

	return files, nil
}

func (s *Storage) Directories(_ context.Context, path string) ([]string, error) {
	f, err := os.ReadDir(s.prefixer.Prefix(path))
	if err != nil {
		return nil, err
	}

	var dirs []string //nolint:prealloc
	for _, file := range f {
		if !file.IsDir() {
			continue
		}
		dirs = append(dirs, file.Name())
	}
	return dirs, nil
}

func (s *Storage) AllDirectories(ctx context.Context, path string) ([]string, error) {
	f, err := os.ReadDir(s.prefixer.Prefix(path))
	if err != nil {
		return nil, err
	}

	var dirs []string
	for _, file := range f {
		if !file.IsDir() {
			continue
		}
		subDirs, err := s.AllDirectories(ctx, file.Name())
		if err != nil {
			return nil, err
		}
		dirs = append(dirs, subDirs...)
	}
	return dirs, nil
}

func (s *Storage) IsFile(_ context.Context, path string) (bool, error) {
	info, err := os.Stat(s.prefixer.Prefix(path))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return !info.IsDir(), nil
}

func (s *Storage) IsDir(_ context.Context, path string) (bool, error) {
	info, err := os.Stat(s.prefixer.Prefix(path))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return info.IsDir(), nil
}
