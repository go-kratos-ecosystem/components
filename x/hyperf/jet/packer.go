package jet

import (
	"encoding/json"
)

var DefaultPacker Packer = NewJsonPacker()

type Packer interface {
	Pack(interface{}) ([]byte, error)
	Unpack([]byte, interface{}) error
}

// JsonPacker is a json packer
type JsonPacker struct{}

func NewJsonPacker() *JsonPacker {
	return &JsonPacker{}
}

func (p *JsonPacker) Pack(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (p *JsonPacker) Unpack(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
