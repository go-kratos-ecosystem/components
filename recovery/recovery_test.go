package recovery

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecovery_Wrap(t *testing.T) {
	fn := func() {
		panic("test")
	}

	// default
	r := New()
	assert.Panics(t, func() {
		r.Wrap(fn)
	})

	// with handler
	r = New(WithHandler(func(err interface{}) {
		assert.Equal(t, "test", err)
	}))

	assert.NotPanics(t, func() {
		r.Wrap(fn)
	})
}
