package jet

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	// transport
	transport, err := NewHttpTransporter(
		WithHttpTransporterAddr("http://localhost:8080/"),
	)
	assert.NoError(t, err)

	// client
	client, err := NewClient(
		WithService("TW8591/Money/MoneyService"),
		WithTransporter(transport),
	)
	assert.NoError(t, err)
	assert.NoError(t, client.Invoke(context.Background(), "getBalance", "request", "response"))
}
