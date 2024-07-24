package local

import (
	"os"

	"github.com/go-kratos-ecosystem/components/v2/storage"
)

type Storage struct {
	root     string
	prefixer *storage.PathPrefixer
}

var _ storage.Storage = (*Storage)(nil)

func NewStorage(root string) *Storage {
	return &Storage{
		root:     root,
		prefixer: storage.NewPathPrefixer(root),
	}
}

func (s *Storage) Get(path string) ([]byte, error) {
	return os.ReadFile(s.prefixer.Prefix(path))
}

func (s *Storage) Set(path string, value []byte) error {
	return os.WriteFile(s.prefixer.Prefix(path), value, 0o644)
}

func (s *Storage) Delete(path string) error {
	return os.Remove(s.prefixer.Prefix(path))
}

func (s *Storage) Has(path string) (bool, error) {
	_, err := os.Stat(s.prefixer.Prefix(path))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (s *Storage) Prepend(path string, value []byte) error {
	old, err := s.Get(path)
	if err != nil {
		return err
	}
	return s.Set(path, append(value, old...))
}

func (s *Storage) Append(path string, value []byte) error {
	old, err := s.Get(path)
	if err != nil {
		return err
	}
	return s.Set(path, append(old, value...))
}

func (s *Storage) Move(oldPath, newPath string) error {
	return os.Rename(s.prefixer.Prefix(oldPath), s.prefixer.Prefix(newPath))
}

func (s *Storage) Copy(oldPath, newPath string) error {
	old, err := s.Get(oldPath)
	if err != nil {
		return err
	}
	return s.Set(newPath, old)
}

func (s *Storage) Link(oldPath, newPath string) error {
	return os.Link(s.prefixer.Prefix(oldPath), s.prefixer.Prefix(newPath))
}

func (s *Storage) Symlink(oldPath, newPath string) error {
	return os.Symlink(s.prefixer.Prefix(oldPath), s.prefixer.Prefix(newPath))
}

func (s *Storage) Files(path string) ([]string, error) {
	f, err := os.ReadDir(s.prefixer.Prefix(path))
	if err != nil {
		return nil, err
	}

	var files []string
	for _, file := range f {
		if file.IsDir() {
			continue
		}
		files = append(files, file.Name())
	}
	return files, nil
}

func (s *Storage) AllFiles(path string) ([]string, error) {
	f, err := os.ReadDir(s.prefixer.Prefix(path))
	if err != nil {
		return nil, err
	}

	var files []string
	for _, file := range f {
		if file.IsDir() {
			subFiles, err := s.AllFiles(file.Name())
			if err != nil {
				return nil, err
			}
			files = append(files, subFiles...)
			continue
		}
		files = append(files, file.Name())
	}

	return files, nil
}

func (s *Storage) Directories(path string) ([]string, error) {
	f, err := os.ReadDir(s.prefixer.Prefix(path))
	if err != nil {
		return nil, err
	}

	var dirs []string
	for _, file := range f {
		if !file.IsDir() {
			continue
		}
		dirs = append(dirs, file.Name())
	}
	return dirs, nil
}

func (s *Storage) AllDirectories(path string) ([]string, error) {
	f, err := os.ReadDir(s.prefixer.Prefix(path))
	if err != nil {
		return nil, err
	}

	var dirs []string
	for _, file := range f {
		if !file.IsDir() {
			continue
		}
		subDirs, err := s.AllDirectories(file.Name())
		if err != nil {
			return nil, err
		}
		dirs = append(dirs, subDirs...)
	}
	return dirs, nil
}

func (s *Storage) IsFile(path string) (bool, error) {
	info, err := os.Stat(s.prefixer.Prefix(path))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return !info.IsDir(), nil
}

func (s *Storage) IsDir(path string) (bool, error) {
	info, err := os.Stat(s.prefixer.Prefix(path))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return info.IsDir(), nil
}
