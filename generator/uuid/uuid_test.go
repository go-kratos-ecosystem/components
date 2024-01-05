package uuid

import (
	"testing"

	"github.com/go-packagist/go-kratos-components/strable"
	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	uuid := Generate()

	assert.NotEmpty(t, uuid)
	assert.Equal(t, 36, len(uuid))
	assert.True(t, strable.IsUuid(uuid))
}
