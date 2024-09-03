package json

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSON(t *testing.T) {
	j1, j2 := Codec, Codec

	assert.Same(t, j1, j2)

	data := map[string]any{
		"foo": "bar",
	}

	// marshal
	bytes1, err := json.Marshal(data)
	assert.NoError(t, err)

	bytes2, err := j1.Marshal(data)
	assert.NoError(t, err)

	assert.Equal(t, bytes1, bytes2)

	// unmarshal
	dest1, dest2 := make(map[string]any), make(map[string]any)
	assert.NoError(t, json.Unmarshal(bytes1, &dest1))
	assert.NoError(t, j1.Unmarshal(bytes1, &dest2))

	assert.Equal(t, dest1, dest2)
}
