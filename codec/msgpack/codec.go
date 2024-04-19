package msgpack

import (
	"github.com/vmihailenco/msgpack/v5"

	"github.com/go-kratos-ecosystem/components/v2/codec"
)

var Codec codec.Codec = &msgPackCodec{}

type msgPackCodec struct{}

func (j *msgPackCodec) Marshal(data any) ([]byte, error) {
	return msgpack.Marshal(data)
}

func (j *msgPackCodec) Unmarshal(src []byte, dest any) error {
	return msgpack.Unmarshal(src, dest)
}
