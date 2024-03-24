package event

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProvider(t *testing.T) {
	var d *Dispatcher
	p := Provider(d)
	ctx, err := p(context.Background())
	assert.NoError(t, err)
	d1, ok1 := FromContext(ctx)
	assert.True(t, ok1)
	assert.Equal(t, d, d1)
}
