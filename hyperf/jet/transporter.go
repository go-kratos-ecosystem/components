package jet

import (
	"bytes"
	"context"
	"errors"
	"fmt"
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

	// check response status code
	if isHTTPTransporterServerFailed(response) {
		return nil, &HTTPTransporterServerError{
			StatusCode: response.StatusCode,
			Message:    response.Status,
			Err:        fmt.Errorf("jet/transporter: failed to send request, status code: %d", response.StatusCode),
		}
	}

	return io.ReadAll(response.Body)
}

func isHTTPTransporterServerFailed(resp *http.Response) bool {
	return resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest
}

func IsHTTPTransporterServerError(err error) bool {
	var target *HTTPTransporterServerError
	return errors.As(err, &target)
}

type HTTPTransporterServerError struct {
	StatusCode int
	Message    string
	Err        error
}

func (e *HTTPTransporterServerError) Error() string {
	return e.Err.Error()
}

func (e *HTTPTransporterServerError) Unwrap() error {
	return e.Err
}
