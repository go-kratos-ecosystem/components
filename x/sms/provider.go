package sms

import "context"

type Provider interface {
	Send(ctx context.Context, message *Message) error
}
