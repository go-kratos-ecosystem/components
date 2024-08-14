package jet

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createServer(t *testing.T, formatter Formatter, packer Packer) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		_, _ = buf.ReadFrom(r.Body)

		// A. parse and check request
		request, err := formatter.ParseRequest(buf.Bytes())
		assert.NoError(t, err)
		assert.Equal(t, "/example/_user/_money/balance", request.Path)

		var params []string
		err = packer.Unpack(request.Params, &params)
		assert.NoError(t, err)
		assert.Equal(t, []string{"flc"}, params)

		// B. response
		var balance float64 = 100.0
		result, err := packer.Pack(balance)
		assert.NoError(t, err)

		response, err := formatter.FormatResponse(&RPCResponse{
			ID:     request.ID,
			Result: result,
		}, nil)
		assert.NoError(t, err)

		_, _ = w.Write(response)
	}))
}

func TestClient_Invoke(t *testing.T) {
	// init
	formatter := DefaultFormatter
	packer := DefaultPacker

	// create srv
	srv := createServer(t, formatter, packer)
	defer srv.Close()

	// create transport
	transport, err := NewHTTPTransporter(
		WithHTTPTransporterAddr(srv.URL),
	)
	assert.NoError(t, err)

	// create client
	client, err := NewClient(
		WithService("Example/User/MoneyService"),
		WithTransporter(transport),
		WithFormatter(formatter),
		WithPacker(packer),
		WithMiddleware(func(next Handler) Handler {
			return func(ctx context.Context, name string, request any) (response any, err error) {
				assert.Equal(t, "balance", name)
				assert.Equal(t, []any{"flc"}, request)
				return next(ctx, name, request)
			}
		}),
	)
	assert.NoError(t, err)

	client.Use(func(next Handler) Handler {
		return func(ctx context.Context, name string, request any) (response any, err error) {
			assert.Equal(t, "balance", name)
			assert.Equal(t, []any{"flc"}, request)
			return next(ctx, name, request)
		}
	})

	// call service
	var balance float64
	err = client.Invoke(context.Background(), "balance", []any{"flc"}, &balance)
	assert.NoError(t, err)
	assert.Equal(t, 100.0, balance)
}

func TestClient_InvalidErrs(t *testing.T) {
	_, err := NewClient()
	assert.Error(t, err)
	assert.Equal(t, ErrClientServiceIsRequired, err)

	_, err = NewClient(WithService("test"), WithTransporter(nil))
	assert.Error(t, err)
	assert.Equal(t, ErrClientTransporterIsRequired, err)
}
