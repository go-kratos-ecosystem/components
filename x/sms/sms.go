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
	gw Provider
}

func New(gw Provider) *Sms {
	return &Sms{
		gw: gw,
	}
}

func (s *Sms) Send(ctx context.Context, message *Message) error {
	return s.gw.Send(ctx, message)
}
