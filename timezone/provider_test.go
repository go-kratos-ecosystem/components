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
	assert.NoError(t, Provider()(context.Background()))
	assert.Equal(t, "UTC", time.Local.String())

	// after
	assert.NoError(t, Provider(Local("Asia/Shanghai"))(context.Background()))
	assert.Equal(t, "Asia/Shanghai", time.Local.String())

	// err
	assert.Error(t, Provider(Local("Asia/Beijing"))(context.Background()))
}
