package errors

import (
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
