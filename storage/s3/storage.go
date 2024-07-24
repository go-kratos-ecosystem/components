package s3

import (
	"context"
	"io"

	"github.com/go-kratos-ecosystem/components/v2/helpers"
	"github.com/go-kratos-ecosystem/components/v2/storage"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Storage struct {
	s3     *s3.Client
	bucket string

	root     string
	prefixer *storage.PathPrefixer
}

var _ storage.Storage = (*Storage)(nil)

type Option func(*Storage)

func WithRoot(root string) Option {
	return func(s *Storage) {
		s.root = root
	}
}

func New(s3 *s3.Client, bucket string, opts ...Option) *Storage {
	s := &Storage{
		s3:     s3,
		bucket: bucket,
		root:   "",
	}
	for _, opt := range opts {
		opt(s)
	}

	s.prefixer = storage.NewPathPrefixer(s.root)

	return s
}

func (s *Storage) Get(ctx context.Context, path string) ([]byte, error) {
	output, err := s.s3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: helpers.Ptr(s.bucket),
		Key:    helpers.Ptr(s.prefixer.Prefix(path)),
	})
	if err != nil {
		return nil, err
	}

	bytes, err := io.ReadAll(output.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (s *Storage) Set(ctx context.Context, path string, value []byte) error {
	// TODO implement me
	panic("implement me")
}

func (s *Storage) Delete(ctx context.Context, path string) error {
	// TODO implement me
	panic("implement me")
}

func (s *Storage) Has(ctx context.Context, path string) (bool, error) {
	// TODO implement me
	panic("implement me")
}

func (s *Storage) Prepend(ctx context.Context, path string, value []byte) error {
	// TODO implement me
	panic("implement me")
}

func (s *Storage) Append(ctx context.Context, path string, value []byte) error {
	// TODO implement me
	panic("implement me")
}

func (s *Storage) Move(ctx context.Context, oldPath, newPath string) error {
	// TODO implement me
	panic("implement me")
}

func (s *Storage) Copy(ctx context.Context, oldPath, newPath string) error {
	// TODO implement me
	panic("implement me")
}

func (s *Storage) Link(ctx context.Context, oldPath, newPath string) error {
	// TODO implement me
	panic("implement me")
}

func (s *Storage) Symlink(ctx context.Context, oldPath, newPath string) error {
	// TODO implement me
	panic("implement me")
}

func (s *Storage) Files(ctx context.Context, path string) ([]string, error) {
	// TODO implement me
	panic("implement me")
}

func (s *Storage) AllFiles(ctx context.Context, path string) ([]string, error) {
	// TODO implement me
	panic("implement me")
}

func (s *Storage) Directories(ctx context.Context, path string) ([]string, error) {
	// TODO implement me
	panic("implement me")
}

func (s *Storage) AllDirectories(ctx context.Context, path string) ([]string, error) {
	// TODO implement me
	panic("implement me")
}

func (s *Storage) IsFile(ctx context.Context, path string) (bool, error) {
	// TODO implement me
	panic("implement me")
}

func (s *Storage) IsDir(ctx context.Context, path string) (bool, error) {
	// TODO implement me
	panic("implement me")
}
