package sms

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSms(t *testing.T) {
	sms := New(NewNullProvider())
	ctx := context.Background()

	assert.NoError(t, sms.Send(ctx, nil, nil))
	assert.NoError(t, sms.SendText(ctx, nil, ""))
	assert.NoError(t, sms.SendTemplate(ctx, nil, "", nil))
	assert.NoError(t, sms.SendTextWithNumber(ctx, "", ""))
	assert.NoError(t, sms.SendTemplateWithNumber(ctx, "", "", nil))
}
