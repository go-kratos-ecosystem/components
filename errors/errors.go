package errors

import (
	"net/http"

	"github.com/go-kratos/kratos/v2/errors"
)

// alias
var (
	IsBadRequest         = errors.IsBadRequest
	IsUnauthorized       = errors.IsUnauthorized
	IsForbidden          = errors.IsForbidden
	IsNotFound           = errors.IsNotFound
	IsConflict           = errors.IsConflict
	IsInternalServer     = errors.IsInternalServer
	IsServiceUnavailable = errors.IsServiceUnavailable
	IsGatewayTimeout     = errors.IsGatewayTimeout
	IsClientClosed       = errors.IsClientClosed
)

func BadRequest(message string) *errors.Error {
	return errors.BadRequest(http.StatusText(http.StatusBadRequest), message)
}

func Unauthorized(message string) *errors.Error {
	return errors.Unauthorized(http.StatusText(http.StatusUnauthorized), message)
}

func Forbidden(message string) *errors.Error {
	return errors.Forbidden(http.StatusText(http.StatusForbidden), message)
}

func NotFound(message string) *errors.Error {
	return errors.NotFound(http.StatusText(http.StatusNotFound), message)
}

func Conflict(message string) *errors.Error {
	return errors.Conflict(http.StatusText(http.StatusConflict), message)
}

func InternalServer(message string) *errors.Error {
	return errors.InternalServer(http.StatusText(http.StatusInternalServerError), message)
}
func ServiceUnavailable(message string) *errors.Error {
	return errors.ServiceUnavailable(http.StatusText(http.StatusServiceUnavailable), message)
}

func GatewayTimeout(message string) *errors.Error {
	return errors.GatewayTimeout(http.StatusText(http.StatusGatewayTimeout), message)
}

func ClientClosed(message string) *errors.Error {
	return errors.ClientClosed("ClientClosed", message)
}
