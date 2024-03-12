package timezone

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimezone(t *testing.T) {
	// before
	assert.Equal(t, "Local", time.Local.String())
	ctx1, err1 := Provider()(context.Background())
	assert.NoError(t, err1)
	assert.Equal(t, "UTC", time.Local.String())
	local1, ok1 := FromContext(ctx1)
	assert.True(t, ok1)
	assert.Equal(t, "UTC", local1.String())

	// after
	ctx2, err2 := Provider(Local("Asia/Shanghai"))(context.Background())
	assert.NoError(t, err2)
	assert.Equal(t, "Asia/Shanghai", time.Local.String())
	local2, ok2 := FromContext(ctx2)
	assert.True(t, ok2)
	assert.Equal(t, "Asia/Shanghai", local2.String())

	// err
	_, err3 := Provider(Local("Asia/Beijing"))(context.Background())
	assert.Error(t, err3)
}
