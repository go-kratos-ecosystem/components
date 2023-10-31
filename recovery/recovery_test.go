package recovery

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecovery_Wrap(t *testing.T) {
	fn := func() error {
		panic("test")

		return nil
	}

	// default
	r := New()
	assert.Panics(t, func() {
		_ = r.Wrap(fn)
	})

	// with handler
	r = New(WithHandler(func(err interface{}) {
		assert.Equal(t, "test", err)
	}))

	assert.NoError(t, r.Wrap(fn))
}
