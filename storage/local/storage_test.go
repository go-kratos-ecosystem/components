package local

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ctx = context.Background()

func TestStorage(t *testing.T) {
	local := New("./testfile/dir1")

	// set
	assert.NoError(t, local.Set(ctx, "test", []byte("test")))

	// get
	data, err := local.Get(ctx, "test")
	assert.NoError(t, err)
	assert.Equal(t, []byte("test"), data)

	data, err = local.Get(ctx, "missing")
	assert.Error(t, err)
	assert.Nil(t, data)

	// has
	has, err := local.Has(ctx, "test")
	assert.NoError(t, err)
	assert.True(t, has)

	missing, err := local.Has(ctx, "missing")
	assert.NoError(t, err)
	assert.False(t, missing)

	// move
	assert.NoError(t, local.Move(ctx, "test", "test2"))
	has, err = local.Has(ctx, "test")
	assert.NoError(t, err)
	assert.False(t, has)
	has, err = local.Has(ctx, "test2")
	assert.NoError(t, err)
	assert.True(t, has)

	// link
	assert.NoError(t, local.Link(ctx, "test3", "test4"))
	has, err = local.Has(ctx, "test4")
	assert.NoError(t, err)
	assert.True(t, has)

	// symlink
	assert.NoError(t, local.Symlink(ctx, "test3", "test5"))
	has, err = local.Has(ctx, "test5")
	assert.NoError(t, err)
	assert.True(t, has)

	// delete
	// assert.NoError(t, local.Delete("test"))
	// has, err = local.Has("test")
	// assert.NoError(t, err)
	// assert.False(t, has)
}
