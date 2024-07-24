package storage

type Repository interface {
	Storage

	// Put sets the value for the given path.
	Put(path string, value []byte) error

	// Destroy deletes the value for the given path.
	Destroy(path string) error

	// Exists checks if the path exists.
	Exists(path string) (bool, error)

	// Missing checks if the path does not exist.
	Missing(path string) (bool, error)

	// Rename renames the value from the old path to the new path.
	Rename(oldPath, newPath string) error
}

type repository struct {
	Storage
}

func NewRepository(store Storage) Repository {
	return &repository{
		Storage: store,
	}
}

func (r *repository) Put(path string, value []byte) error {
	return r.Set(path, value)
}

func (r *repository) Destroy(path string) error {
	return r.Delete(path)
}

func (r *repository) Exists(path string) (bool, error) {
	return r.Has(path)
}

func (r *repository) Missing(path string) (bool, error) {
	had, err := r.Has(path)
	if err != nil {
		return false, err
	}

	return !had, nil
}

func (r *repository) Rename(oldPath, newPath string) error {
	return r.Move(oldPath, newPath)
}
