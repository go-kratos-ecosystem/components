package s3

import (
	"bytes"
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

var (
	_ storage.Storage  = (*Storage)(nil)
	_ storage.Copyable = (*Storage)(nil)
)

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

	body, err := io.ReadAll(output.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (s *Storage) Set(ctx context.Context, path string, value []byte) error {
	_, err := s.s3.PutObject(ctx, &s3.PutObjectInput{
		Bucket: helpers.Ptr(s.bucket),
		Key:    helpers.Ptr(s.prefixer.Prefix(path)),
		Body:   io.NopCloser(bytes.NewBuffer(value)),
	})
	return err
}

func (s *Storage) Delete(ctx context.Context, path string) error {
	_, err := s.s3.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: helpers.Ptr(s.bucket),
		Key:    helpers.Ptr(s.prefixer.Prefix(path)),
	})
	return err
}

func (s *Storage) Has(ctx context.Context, path string) (bool, error) {
	_, err := s.s3.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: helpers.Ptr(s.bucket),
		Key:    helpers.Ptr(s.prefixer.Prefix(path)),
	})
	if err != nil {
		// if _, ok := err.(*s3.HeadObjectNotFound); ok {
		// 	return false, nil
		// }
		return false, err
	}

	return true, nil
}

func (s *Storage) Move(ctx context.Context, oldPath, newPath string) error {
	if err := s.Copy(ctx, oldPath, newPath); err != nil {
		return err
	}
	return s.Delete(ctx, oldPath)
}

func (s *Storage) Copy(ctx context.Context, oldPath, newPath string) error {
	_, err := s.s3.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     helpers.Ptr(s.bucket),
		CopySource: helpers.Ptr(s.bucket + "/" + s.prefixer.Prefix(oldPath)),
		Key:        helpers.Ptr(s.prefixer.Prefix(newPath)),
	})
	return err
}

func (s *Storage) Link(context.Context, string, string) error {
	return storage.ErrNotSupported
}

func (s *Storage) Symlink(context.Context, string, string) error {
	return storage.ErrNotSupported
}

func (s *Storage) Files(ctx context.Context, path string) ([]string, error) {
	output, err := s.s3.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: helpers.Ptr(s.bucket),
		Prefix: helpers.Ptr(s.prefixer.Prefix(path)),
	})
	if err != nil {
		return nil, err
	}

	var files []string
	for _, obj := range output.Contents {
		files = append(files, *obj.Key)
	}
	return files, nil
}

func (s *Storage) AllFiles(ctx context.Context, path string) ([]string, error) {
	output, err := s.s3.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: helpers.Ptr(s.bucket),
		Prefix: helpers.Ptr(s.prefixer.Prefix(path)),
	})
	if err != nil {
		return nil, err
	}

	var files []string
	for _, obj := range output.Contents {
		files = append(files, *obj.Key)
	}
	return files, nil
}

func (s *Storage) Directories(ctx context.Context, path string) ([]string, error) {
	output, err := s.s3.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: helpers.Ptr(s.bucket),
		Prefix: helpers.Ptr(s.prefixer.Prefix(path)),
	})
	if err != nil {
		return nil, err
	}

	var dirs []string
	for _, obj := range output.Contents {
		dirs = append(dirs, *obj.Key)
	}
	return dirs, nil
}

func (s *Storage) AllDirectories(ctx context.Context, path string) ([]string, error) {
	output, err := s.s3.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: helpers.Ptr(s.bucket),
		Prefix: helpers.Ptr(s.prefixer.Prefix(path)),
	})
	if err != nil {
		return nil, err
	}

	var dirs []string
	for _, obj := range output.Contents {
		dirs = append(dirs, *obj.Key)
	}
	return dirs, nil
}

func (s *Storage) IsFile(ctx context.Context, path string) (bool, error) {
	// TODO implement me
	panic("implement me")
}

func (s *Storage) IsDir(ctx context.Context, path string) (bool, error) {
	// TODO implement me
	panic("implement me")
}
