package sms

import "context"

type NullProvider struct{}

func NewNullProvider() Provider {
	return &NullProvider{}
}

func (p *NullProvider) Send(_ context.Context, _ *Message) error {
	return nil
}
