package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorage(t *testing.T) {
	ctx := context.Background()

	_, err := NoopStorage.Get(ctx, "test")
	assert.NoError(t, err)

	assert.NoError(t, NoopStorage.Set(ctx, "noop", []byte("noop")))
	assert.NoError(t, NoopStorage.Delete(ctx, "noop"))

	has, err := NoopStorage.Has(ctx, "noop")
	assert.NoError(t, err)
	assert.True(t, has)

	assert.NoError(t, NoopStorage.Move(ctx, "noop", "noop"))
	assert.NoError(t, NoopStorage.Link(ctx, "noop", "noop"))
	assert.NoError(t, NoopStorage.Symlink(ctx, "noop", "noop"))

	files, err := NoopStorage.Files(ctx, "noop")
	assert.NoError(t, err)
	assert.Len(t, files, 0)

	allFiles, err := NoopStorage.AllFiles(ctx, "noop")
	assert.NoError(t, err)
	assert.Len(t, allFiles, 0)

	directories, err := NoopStorage.Directories(ctx, "noop")
	assert.NoError(t, err)
	assert.Len(t, directories, 0)

	allDirectories, err := NoopStorage.AllDirectories(ctx, "noop")
	assert.NoError(t, err)
	assert.Len(t, allDirectories, 0)

	isFile, err := NoopStorage.IsFile(ctx, "noop")
	assert.NoError(t, err)
	assert.False(t, isFile)

	isDir, err := NoopStorage.IsDir(ctx, "noop")
	assert.NoError(t, err)
	assert.False(t, isDir)
}
