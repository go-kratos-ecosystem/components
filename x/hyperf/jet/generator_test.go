package jet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFullPathGenerator_Generate(t *testing.T) {
	generator := NewFullPathGenerator()

	tests := []struct {
		service string
		name    string
		want    string
	}{
		{`Foo\UserService`, "query", "/Foo/UserService/query"},
		{`Foo\UserService`, "Query", "/Foo/UserService/Query"},
		{`user`, "query", "/user/query"},
	}

	for _, tt := range tests {
		t.Run(tt.service+"_"+tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, generator.Generate(tt.service, tt.name))
		})
	}
}
