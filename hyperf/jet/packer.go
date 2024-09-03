package jet

import (
	"encoding/json"
)

var DefaultPacker Packer = NewJSONPacker()

type Packer interface {
	Pack(any) ([]byte, error)
	Unpack([]byte, any) error
}

// JSONPacker is a json packer
type JSONPacker struct{}

func NewJSONPacker() *JSONPacker {
	return &JSONPacker{}
}

func (p *JSONPacker) Pack(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (p *JSONPacker) Unpack(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
