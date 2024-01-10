package sms

import "context"

type Provider interface {
	Send(ctx context.Context, phone *Phone, message *Message) error
}
