package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	err1 = errors.New("err1")
	err2 = errors.New("err2")
	err3 = errors.New("err3")
)

func TestGroup(t *testing.T) {
	g := NewGroup()

	assert.Equal(t, 0, g.Len())
	assert.Equal(t, g, g.Add(nil, err1, err2))

	assert.Equal(t, 2, g.Len())
	assert.Equal(t, multipleErrors, g.Error())
	assert.Equal(t, []error{err1, err2}, g.Errors())
	assert.True(t, g.Has(err1))
	assert.True(t, g.Has(err2))
	assert.False(t, g.Has(err3))
	assert.Equal(t, err1, g.First())
}
