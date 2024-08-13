package jet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneralPathGenerator_Generate(t *testing.T) {
	generator := NewGeneralPathGenerator()

	tests := []struct {
		service string
		name    string
		want    string
	}{
		{`CN12345/User/MoneyService`, "getBalance", "/c_n12345/_user/_money/getBalance"},
		{`USA888/User/OrderService`, "create", "/u_s_a888/_user/_order/create"},
		{`Usa888/Money/LinkService`, "Create", "/usa888/_money/_link/Create"},
		{`/CN12345/User/MoneyService`, "getBalance", "/_c_n12345/_user/_money/getBalance"},
		{`/USA888/User/OrderService`, "create", "/_u_s_a888/_user/_order/create"},
		{`/Usa888/Money/LinkService`, "Create", "/_usa888/_money/_link/Create"},
		{`Foo\UserService`, "create", "/user/create"},
		{`\Foo\UserService`, "Create", "/user/Create"},
	}

	for _, tt := range tests {
		t.Run(tt.service+"@"+tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, generator.Generate(tt.service, tt.name))
		})
	}
}

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
