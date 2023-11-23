package json

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSerializer(t *testing.T) {
	var j1, j2 = Serializer, Serializer

	assert.Same(t, j1, j2)

	var data = map[string]interface{}{
		"foo": "bar",
	}

	// marshal
	bytes1, err := json.Marshal(data)
	assert.NoError(t, err)

	bytes2, err := j1.Serialize(data)
	assert.NoError(t, err)

	assert.Equal(t, bytes1, bytes2)

	// unmarshal
	var dest1, dest2 = make(map[string]interface{}), make(map[string]interface{})
	assert.NoError(t, json.Unmarshal(bytes1, &dest1))
	assert.NoError(t, j1.Unserialize(bytes1, &dest2))

	assert.Equal(t, dest1, dest2)

}
