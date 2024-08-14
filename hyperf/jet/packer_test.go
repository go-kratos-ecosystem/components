package jet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPacker_JSONPacker(t *testing.T) {
	packer := NewJSONPacker()

	// pack
	data, err := packer.Pack(1)
	assert.NoError(t, err)

	// unpack
	var v int
	err = packer.Unpack(data, &v)
	assert.NoError(t, err)
	assert.Equal(t, 1, v)
}
