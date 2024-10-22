package locker

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOwner_owner(t *testing.T) {
	o := NewOwner(NoopLocker{}, WithName("test"))
	assert.Equal(t, "test", o.Name())

	ok, err := o.Release(context.Background())
	assert.NoError(t, err)
	assert.True(t, ok)
}
