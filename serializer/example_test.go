package serializer_test

import (
	"log"

	"github.com/go-packagist/go-kratos-components/serializer/json"
)

var j = json.Serializer

func ExampleJSON() {
	bytes, err := j.Serialize(map[string]string{
		"key": "value",
	})
	if err != nil {
		log.Fatal(err)
	}

	var dest map[string]string
	err = j.Unserialize(bytes, &dest)
	if err != nil {
		log.Fatal(err)
	}
}
