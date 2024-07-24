package storage

import "context"

type Repository interface {
	Storage
	Copyable

	// Put sets the value for the given path.
	Put(ctx context.Context, path string, value []byte) error

	// Destroy deletes the value for the given path.
	Destroy(ctx context.Context, path string) error

	// Exists checks if the path exists.
	Exists(ctx context.Context, path string) (bool, error)

	// Missing checks if the path does not exist.
	Missing(ctx context.Context, path string) (bool, error)

	// Rename renames the value from the old path to the new path.
	Rename(ctx context.Context, oldPath, newPath string) error

	// Prepend prepends the value to the existing value for the given path.
	Prepend(ctx context.Context, path string, value []byte) error

	// Append appends the value to the existing value for the given path.
	Append(ctx context.Context, path string, value []byte) error
}

type repository struct {
	Storage
}

func NewRepository(store Storage) Repository {
	return &repository{
		Storage: store,
	}
}

func (r *repository) Put(ctx context.Context, path string, value []byte) error {
	return r.Set(ctx, path, value)
}

func (r *repository) Destroy(ctx context.Context, path string) error {
	return r.Delete(ctx, path)
}

func (r *repository) Exists(ctx context.Context, path string) (bool, error) {
	return r.Has(ctx, path)
}

func (r *repository) Missing(ctx context.Context, path string) (bool, error) {
	had, err := r.Has(ctx, path)
	if err != nil {
		return false, err
	}

	return !had, nil
}

func (r *repository) Rename(ctx context.Context, oldPath, newPath string) error {
	return r.Move(ctx, oldPath, newPath)
}

func (r *repository) Prepend(ctx context.Context, path string, value []byte) error {
	old, err := r.Get(ctx, path)
	if err != nil {
		return err
	}
	return r.Set(ctx, path, append(value, old...))
}

func (r *repository) Append(ctx context.Context, path string, value []byte) error {
	old, err := r.Get(ctx, path)
	if err != nil {
		return err
	}
	return r.Set(ctx, path, append(old, value...))
}

func (r *repository) Copy(ctx context.Context, oldPath, newPath string) error {
	if copier, ok := r.Storage.(Copyable); ok {
		return copier.Copy(ctx, oldPath, newPath)
	}

	old, err := r.Get(ctx, oldPath)
	if err != nil {
		return err
	}
	return r.Set(ctx, newPath, old)
}
