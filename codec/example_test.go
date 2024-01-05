package codec_test

import (
	"log"
	"testing"

	"github.com/go-kratos-ecosystem/components/v2/codec/json"
)

var j = json.Codec

func TestJson(_ *testing.T) {
	bytes, err := j.Marshal(map[string]string{
		"key": "value",
	})
	if err != nil {
		log.Fatal(err)
	}

	var dest map[string]string
	err = j.Unmarshal(bytes, &dest)
	if err != nil {
		log.Fatal(err)
	}
}
