package locker

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOwner_owner(t *testing.T) {
	o := NewOwner(NoopLocker{}, WithOwnerName("test"))
	assert.Equal(t, "test", o.Name())
	assert.NoError(t, o.Release(context.Background()))
}
