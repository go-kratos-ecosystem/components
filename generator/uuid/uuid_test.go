package uuid

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-kratos-ecosystem/components/v2/strable"
)

func TestGenerate(t *testing.T) {
	uuid := Generate()

	assert.NotEmpty(t, uuid)
	assert.Equal(t, 36, len(uuid))
	assert.True(t, strable.IsUuid(uuid))
}
