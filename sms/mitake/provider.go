// Package mitake is a sms provider for mitake.
// See: https://sms.mitake.com.tw/
package mitake

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/text/encoding/traditionalchinese"

	"github.com/go-kratos-ecosystem/components/v2/debug"
	"github.com/go-kratos-ecosystem/components/v2/sms"
)

type provider struct {
	api      string
	username string
	password string

	httpClient *http.Client
}

type Option func(t *provider)

func WithAPI(api string) Option {
	return func(t *provider) {
		t.api = api
	}
}

func WithHTTPClient(httpClient *http.Client) Option {
	return func(t *provider) {
		t.httpClient = httpClient
	}
}

func New(username, password string, opts ...Option) sms.Provider {
	p := &provider{
		api:        "https://smsapi.mitake.com.tw",
		username:   username,
		password:   password,
		httpClient: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (p *provider) Send(ctx context.Context, phone *sms.Phone, message *sms.Message) error {
	if err := p.verify(phone, message); err != nil {
		return err
	}

	// Convert to Big5
	text, err := traditionalchinese.Big5.NewEncoder().String(message.Text)
	if err != nil {
		return err
	}

	// Combine params
	params := url.Values{}
	params.Set("username", p.username)
	params.Set("password", p.password)
	params.Set("type", "now")
	params.Set("encoding", "big5")
	params.Set("dstaddr", phone.Number)
	params.Set("smbody", text)

	// new request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.api+"?"+params.Encode(), nil)
	if err != nil {
		return err
	}

	// send request
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	debug.Dump(resp.Body)

	// check response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("sms mitake: http status code: %d", resp.StatusCode)
	}

	return nil
}

func (p *provider) verify(phone *sms.Phone, message *sms.Message) error {
	if phone.Number == "" {
		return sms.ErrInvalidPhone
	}

	if message.Text == "" {
		return sms.ErrInvalidMessage
	}

	return nil
}
