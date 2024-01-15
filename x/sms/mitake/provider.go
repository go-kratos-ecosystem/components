// Package mitake is a sms provider for mitake.
// See: https://sms.mitake.com.tw/
package mitake

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/text/encoding/traditionalchinese"

	"github.com/go-kratos-ecosystem/components/v2/x/sms"
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

func (p *provider) Send(ctx context.Context, message *sms.Message) (err error) {
	if err = p.verify(message); err != nil {
		return
	}

	// Convert to Big5
	text, err := traditionalchinese.Big5.NewEncoder().String(message.Content.Text)
	if err != nil {
		return
	}

	// Combine params
	params := url.Values{}
	params.Set("username", p.username)
	params.Set("password", p.password)
	params.Set("type", "now")
	params.Set("encoding", "big5")
	params.Set("dstaddr", message.Phone.Number)
	params.Set("smbody", text)

	// new message
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.api+"?"+params.Encode(), nil)
	if err != nil {
		return
	}

	// send message
	rep, err := p.httpClient.Do(req)
	if err != nil {
		return
	}
	defer rep.Body.Close()

	// check rep
	if rep.StatusCode != http.StatusOK {
		err = fmt.Errorf("sms mitake: http status code: %d", rep.StatusCode)
		return
	}

	return nil
}

func (p *provider) verify(message *sms.Message) error {
	if message.Phone == nil || message.Phone.Number == "" {
		return sms.ErrInvalidPhone
	}

	if message.Content == nil || message.Content.Text == "" {
		return sms.ErrInvalidMessage
	}

	return nil
}
