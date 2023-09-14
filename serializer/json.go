package serializer

import (
	"encoding/json"

	"github.com/go-packagist/go-kratos-components/contract"
)

var JsonSerializer = newJsonSerializer()

type jsonSerializer struct{}

func newJsonSerializer() contract.Serializable {
	return &jsonSerializer{}
}

func (j *jsonSerializer) Serialize(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func (j *jsonSerializer) Unserialize(src []byte, dest interface{}) error {
	return json.Unmarshal(src, dest)
}
