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

	// path
	path := local.Path(ctx, "1.jpg")
	assert.Equal(t, "./testfile/dir1/1.jpg", path)

	// delete
	// assert.NoError(t, local.Delete("test"))
	// has, err = local.Has("test")
	// assert.NoError(t, err)
	// assert.False(t, has)
}

func TestStorage_Path(t *testing.T) {
	local := New("./testfile/path")

	tests := []struct {
		path string
		want string
	}{
		{"1.jpg", "./testfile/path/1.jpg"},
		{"2.jpg", "./testfile/path/2.jpg"},
		{"2", "./testfile/path/2"},
		{"2/3", "./testfile/path/2/3"},
		{"/4/3.jpg", "./testfile/path/4/3.jpg"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			assert.Equal(t, tt.want, local.Path(ctx, tt.path))
		})
	}
}

func TestStorage_Name(t *testing.T) {
	local := New("./testfile/path")

	tests := []struct {
		path string
		want string
	}{
		{".jpg", ""},
		{"", "."},
		{"1.jpg", "1"},
		{"2.jpg", "2"},
		{"2", "2"},
		{"2/3", "3"},
		{"/4/3.jpg", "3"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			assert.Equal(t, tt.want, local.Name(ctx, tt.path))
		})
	}
}

func TestStorage_Basename(t *testing.T) {
	local := New("./testfile/path")

	tests := []struct {
		path string
		want string
	}{
		{".jpg", ".jpg"},
		{"", "."},
		{"1.jpg", "1.jpg"},
		{"2.jpg", "2.jpg"},
		{"2", "2"},
		{"2/3", "3"},
		{"/4/3.jpg", "3.jpg"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			assert.Equal(t, tt.want, local.Basename(ctx, tt.path))
		})
	}
}

func TestStorage_Dirname(t *testing.T) {
	local := New("./testfile/path")

	tests := []struct {
		path string
		want string
	}{
		{".jpg", "testfile/path"},
		{"", "testfile/path"},
		{"1.jpg", "testfile/path"},
		{"2.jpg", "testfile/path"},
		{"2", "testfile/path"},
		{"2/3", "testfile/path/2"},
		{"/4/3.jpg", "testfile/path/4"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			assert.Equal(t, tt.want, local.Dirname(ctx, tt.path))
		})
	}
}

func TestStorage_Extension(t *testing.T) {
	local := New("./testfile/path")

	tests := []struct {
		path string
		want string
	}{
		{".jpg", ".jpg"},
		{"", ""},
		{"1.jpg", ".jpg"},
		{"2.jpg", ".jpg"},
		{"2", ""},
		{"2/3", ""},
		{"/4/3.jpg", ".jpg"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			assert.Equal(t, tt.want, local.Extension(ctx, tt.path))
		})
	}
}
