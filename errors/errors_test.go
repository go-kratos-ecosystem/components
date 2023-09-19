package errors

import (
	"testing"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/stretchr/testify/assert"
)

func TestErrors(t *testing.T) {
	assert.True(t, errors.IsBadRequest(BadRequest("BadRequest")))
	assert.True(t, errors.IsUnauthorized(Unauthorized("Unauthorized")))
	assert.True(t, errors.IsForbidden(Forbidden("Forbidden")))
	assert.True(t, errors.IsNotFound(NotFound("NotFound")))
	assert.True(t, errors.IsConflict(Conflict("Conflict")))
	assert.True(t, errors.IsInternalServer(InternalServer("InternalServer")))
	assert.True(t, errors.IsServiceUnavailable(ServiceUnavailable("ServiceUnavailable")))
	assert.True(t, errors.IsGatewayTimeout(GatewayTimeout("GatewayTimeout")))
	assert.True(t, errors.IsClientClosed(ClientClosed("ClientClosed")))
	assert.True(t, errors.IsClientClosed(ClientClosed("ClientClosed")))
}
