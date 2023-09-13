package serializer

import (
	"encoding/json"

	"github.com/go-packagist/go-kratos-components/contracts/serializer"
)

var JsonSerializer = newJsonSerializer()

type jsonSerializer struct{}

func newJsonSerializer() serializer.Serializable {
	return &jsonSerializer{}
}

func (j *jsonSerializer) Serialize(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func (j *jsonSerializer) Unserialize(src []byte, dest interface{}) error {
	return json.Unmarshal(src, dest)
}
