package json

import (
	"encoding/json"

	"github.com/go-kratos-ecosystem/components/v2/serializer"
)

var Serializer serializer.Serializable = &jsonSerializer{}

type jsonSerializer struct{}

func (j *jsonSerializer) Serialize(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func (j *jsonSerializer) Unserialize(src []byte, dest interface{}) error {
	return json.Unmarshal(src, dest)
}
