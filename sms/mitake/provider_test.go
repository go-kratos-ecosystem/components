package mitake

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-kratos-ecosystem/components/v2/sms"
)

func TestProvider(t *testing.T) {
	var (
		username = "test"
		password = "test"
		number   = "0939474570"
		text     = "Hello, world"
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatal("method error")
		}

		if r.URL.Query().Get("username") != username {
			t.Fatal("username error")
		}

		if r.URL.Query().Get("password") != password {
			t.Fatal("password error")
		}

		if r.URL.Query().Get("dstaddr") != number {
			t.Fatal("dstaddr error")
		}

		if r.URL.Query().Get("smbody") != text {
			t.Fatal("smbody error")
		}

		w.Write([]byte("")) //nolint:errcheck
	}))
	defer srv.Close()

	p := New(username, password,
		WithAPI(srv.URL),
	)

	err := p.Send(context.Background(), &sms.Phone{
		Number: number,
	}, &sms.Message{
		Text: text,
	})

	if err != nil {
		t.Fatal(err)
	}
}
