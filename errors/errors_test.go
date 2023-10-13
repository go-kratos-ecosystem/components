package errors

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrors(t *testing.T) {
	assert.True(t, IsBadRequest(BadRequest("BadRequest")))
	assert.True(t, IsUnauthorized(Unauthorized("Unauthorized")))
	assert.True(t, IsForbidden(Forbidden("Forbidden")))
	assert.True(t, IsNotFound(NotFound("NotFound")))
	assert.True(t, IsConflict(Conflict("Conflict")))
	assert.True(t, IsInternalServer(InternalServer("InternalServer")))
	assert.True(t, IsServiceUnavailable(ServiceUnavailable("ServiceUnavailable")))
	assert.True(t, IsGatewayTimeout(GatewayTimeout("GatewayTimeout")))
	assert.True(t, IsClientClosed(ClientClosed("ClientClosed")))
	assert.True(t, IsClientClosed(ClientClosed("ClientClosed")))
}

func TestErrors_New(t *testing.T) {
	assert.True(t, IsForbidden(New(http.StatusForbidden, "Forbidden")))
	assert.Equal(t, New(http.StatusForbidden, "Forbidden"), Forbidden("Forbidden"))
}
