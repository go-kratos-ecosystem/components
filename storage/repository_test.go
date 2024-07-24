package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoopRepository(t *testing.T) {
	repo := NewRepository(NoopStorage)
	ctx := context.Background()

	assert.NoError(t, repo.Put(ctx, "noop", []byte("noop")))
	assert.NoError(t, repo.Destroy(ctx, "noop"))

	exists, err := repo.Exists(ctx, "noop")
	assert.NoError(t, err)
	assert.True(t, exists)

	missing, err := repo.Missing(ctx, "noop")
	assert.NoError(t, err)
	assert.False(t, missing)

	assert.NoError(t, repo.Rename(ctx, "noop", "noop2"))
}
