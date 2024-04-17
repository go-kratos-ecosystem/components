package chi

import (
	"context"
	"io"
	"net/http"
	"sync"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	var (
		srv = NewServer(
			chi.NewRouter(),
			Addr(":8001"),
		)
		ch = make(chan string, 1)
		wg sync.WaitGroup
	)
	wg.Add(1)

	srv.Get("/ping", func(w http.ResponseWriter, _ *http.Request) {
		defer wg.Done()
		ch <- "pong"
		_, _ = w.Write([]byte("pong"))
	})

	go func() {
		srv.Start(context.Background()) // nolint: errcheck
	}()

	resp, err := http.Get("http://localhost:8001/ping")

	wg.Wait()

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	defer resp.Body.Close() // nolint: errcheck
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "pong", string(body))

	assert.Equal(t, "pong", <-ch)

	err = srv.Stop(context.Background())
	assert.NoError(t, err)
}
