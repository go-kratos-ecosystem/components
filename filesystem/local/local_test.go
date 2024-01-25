package local

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var l = New("../../.testdata")

func TestLocal_Exists(t *testing.T) {
	exists, err := l.Exists("local-go-dark.jpeg")
	assert.NoError(t, err)
	assert.True(t, exists)

	notExists, err := l.Exists("not-exists.jpg")
	assert.NoError(t, err)
	assert.False(t, notExists)
}

func TestLocal_Get(t *testing.T) {
	data, err := l.Get("local-get.txt")
	assert.NoError(t, err)
	assert.Equal(t, "hello world", string(data))
}
