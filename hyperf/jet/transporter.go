package jet

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
)

var (
	ErrHTTPTransporterAddrIsRequired     = errors.New("jet/transporter: addr is required")
	ErrorHTTPTransporterClientIsRequired = errors.New("jet/transporter: client is required")
)

type Transporter interface {
	Send(ctx context.Context, data []byte) ([]byte, error)
}

// --------------------------------------------------------
// HTTPTransporter implementation
// --------------------------------------------------------

// HTTPTransporter is a http transporter
type HTTPTransporter struct {
	addr string
	*http.Client
}

type HTTPTransporterOption func(*HTTPTransporter)

func WithHTTPTransporterAddr(addr string) HTTPTransporterOption {
	return func(t *HTTPTransporter) {
		t.addr = addr
	}
}

func WithHTTPTransporterClient(client *http.Client) HTTPTransporterOption {
	return func(t *HTTPTransporter) {
		t.Client = client
	}
}

func NewHTTPTransporter(opts ...HTTPTransporterOption) (*HTTPTransporter, error) {
	transport := &HTTPTransporter{
		Client: http.DefaultClient,
	}
	for _, opt := range opts {
		opt(transport)
	}

	// validate
	if transport.addr == "" {
		return nil, ErrHTTPTransporterAddrIsRequired
	}
	if transport.Client == nil {
		return nil, ErrorHTTPTransporterClientIsRequired
	}

	return transport, nil
}

func (t *HTTPTransporter) Send(ctx context.Context, data []byte) ([]byte, error) {
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
