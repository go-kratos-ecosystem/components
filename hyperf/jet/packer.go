package jet

import (
	"encoding/json"
)

var DefaultPacker Packer = NewJSONPacker()

type Packer interface {
	Pack(interface{}) ([]byte, error)
	Unpack([]byte, interface{}) error
}

// JSONPacker is a json packer
type JSONPacker struct{}

func NewJSONPacker() *JSONPacker {
	return &JSONPacker{}
}

func (p *JSONPacker) Pack(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (p *JSONPacker) Unpack(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
