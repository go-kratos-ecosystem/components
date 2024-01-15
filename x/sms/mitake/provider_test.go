package mitake

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-kratos-ecosystem/components/v2/x/sms"
)

func TestProvider(t *testing.T) {
	var (
		username = "test"
		password = "test"
		number   = "123456789"
		text     = "Hello, world"
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, username, r.URL.Query().Get("username"))
		assert.Equal(t, password, r.URL.Query().Get("password"))
		assert.Equal(t, number, r.URL.Query().Get("dstaddr"))
		assert.Equal(t, text, r.URL.Query().Get("smbody"))

		w.Write([]byte("hello")) //nolint:errcheck
	}))
	defer srv.Close()

	p := New(username, password,
		WithAPI(srv.URL),
		WithHTTPClient(http.DefaultClient),
	)

	err := p.Send(context.Background(), &sms.Message{
		Phone: &sms.Phone{
			Number: number,
		},
		Content: &sms.Content{
			Text: text,
		},
	})
	assert.NoError(t, err)
}
