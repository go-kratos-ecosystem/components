package hashing_test

import (
	"testing"

	"github.com/go-packagist/go-kratos-components/hashing"
	_ "github.com/go-packagist/go-kratos-components/hashing/md5"
	"github.com/stretchr/testify/assert"
)

func TestHasher_New(t *testing.T) {
	hashedValue, err := hashing.MD5.New().Make("123456")
	assert.NoError(t, err)

	assert.True(t, hashing.MD5.New().Check("123456", hashedValue))

	// unknown
	assert.Panics(t, func() {
		hashing.Hash(999999).New()
	})
}

func TestHasher_Instance(t *testing.T) {
	h1 := hashing.MD5.New()
	h2 := hashing.MD5.New()

	assert.Same(t, h1, h2)
}
