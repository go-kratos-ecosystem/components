package jet

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/davecgh/go-spew/spew"
)

type Transporter interface {
	Send(ctx context.Context, data []byte) ([]byte, error)
}

// --------------------------------------------------------
// HttpTransporter implementation
// --------------------------------------------------------

// HttpTransporter is a http transporter
type HttpTransporter struct {
	addr string
	*http.Client
}

type HttpTransporterOption func(*HttpTransporter)

func WithHttpTransporterAddr(addr string) HttpTransporterOption {
	return func(t *HttpTransporter) {
		t.addr = addr
	}
}

func WithHttpTransporterClient(client *http.Client) HttpTransporterOption {
	return func(t *HttpTransporter) {
		t.Client = client
	}
}

func NewHttpTransporter(opts ...HttpTransporterOption) (*HttpTransporter, error) {
	transport := &HttpTransporter{
		Client: http.DefaultClient,
	}
	for _, opt := range opts {
		opt(transport)
	}

	// validate
	if transport.addr == "" {
		return nil, errors.New("jet/transporter: addr is required")
	}
	if transport.Client == nil {
		return nil, errors.New("jet/transporter: client is required")
	}

	return transport, nil
}

func (t *HttpTransporter) Send(ctx context.Context, data []byte) ([]byte, error) {
	spew.Dump(data)

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, t.addr, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	response, err := t.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close() // nolint:errcheck
	return io.ReadAll(response.Body)
}
