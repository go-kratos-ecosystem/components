package sms

import (
	"context"
	"errors"
)

var (
	ErrInvalidPhone   = errors.New("sms: invalid phone")
	ErrInvalidMessage = errors.New("sms: invalid message")
)

type Phone struct {
	IDDCode string
	Number  string
}

type Content struct {
	Text      string
	Template  string
	Variables map[string]string
}

type Message struct {
	Phone   *Phone
	Content *Content
}

type Sms struct {
	provider Provider
}

func New(provider Provider) *Sms {
	return &Sms{
		provider: provider,
	}
}

func (s *Sms) Send(ctx context.Context, message *Message) error {
	return s.provider.Send(ctx, message)
}
