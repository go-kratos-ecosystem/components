package sms

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSms(t *testing.T) {
	sms := New(NewNullProvider())
	ctx := context.Background()

	err := sms.Send(ctx, &Message{})
	assert.NoError(t, err)
}
