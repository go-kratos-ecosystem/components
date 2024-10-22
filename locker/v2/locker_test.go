package locker

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLocker_NoopLocker(t *testing.T) {
	l := NoopLocker{}
	ctx := context.Background()

	ok, err := l.Try(ctx, func() {})
	assert.NoError(t, err)
	assert.True(t, ok)

	ok, err = l.Until(ctx, time.Second, func() {})
	assert.NoError(t, err)
	assert.True(t, ok)

	o, ok, err := l.Get(ctx)
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Nil(t, o)

	ok, err = l.Release(ctx, nil)
	assert.NoError(t, err)
	assert.True(t, ok)

	assert.NoError(t, l.ForceRelease(ctx))

	o, err = l.LockedOwner(ctx)
	assert.NoError(t, err)
	assert.Nil(t, o)
}
