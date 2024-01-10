package sms

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSms(t *testing.T) {
	sms := New(NewNullProvider())

	assert.NoError(t, sms.Send(nil, nil, nil))
	assert.NoError(t, sms.SendText(nil, nil, ""))
	assert.NoError(t, sms.SendTemplate(nil, nil, "", nil))
	assert.NoError(t, sms.SendTextWithNumber(nil, "", ""))
	assert.NoError(t, sms.SendTemplateWithNumber(nil, "", "", nil))
}
