package sms

import (
	"context"
	"errors"
	"time"
)

var (
	ErrInvalidPhone   = errors.New("sms: invalid phone")
	ErrInvalidMessage = errors.New("sms: invalid message")
)

type Phone struct {
	IDDCode string
	Number  string
}

type Message struct {
	Text      string
	Template  string
	Variables map[string]string
	Schedule  time.Time
}

type Sms struct {
	gw Provider
}

func New(gw Provider) *Sms {
	return &Sms{
		gw: gw,
	}
}

func (s *Sms) Send(ctx context.Context, phone *Phone, message *Message) error {
	return s.gw.Send(ctx, phone, message)
}

func (s *Sms) SendText(ctx context.Context, phone *Phone, text string) error {
	return s.Send(ctx, phone, &Message{
		Text: text,
	})
}

func (s *Sms) SendTemplate(ctx context.Context, phone *Phone, template string, variables map[string]string) error {
	return s.Send(ctx, phone, &Message{
		Template:  template,
		Variables: variables,
	})
}

func (s *Sms) SendTextWithNumber(ctx context.Context, phoneNumber, text string) error {
	return s.SendText(ctx, &Phone{
		Number: phoneNumber,
	}, text)
}

func (s *Sms) SendTemplateWithNumber(ctx context.Context, phoneNumber, template string, variables map[string]string) error { //nolint:lll
	return s.SendTemplate(ctx, &Phone{
		Number: phoneNumber,
	}, template, variables)
}
