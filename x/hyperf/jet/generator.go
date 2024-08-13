package jet

import (
	"strings"

	"github.com/google/uuid"
)

// ================================================================================================
// PathGenerator generates the path of the service method
// ================================================================================================

var DefaultPathGenerator PathGenerator = NewFullPathGenerator()

type PathGenerator interface {
	Generate(service string, name string) string
}

// FullPathGenerator generates the full path of the service method
type FullPathGenerator struct{}

func NewFullPathGenerator() *FullPathGenerator {
	return &FullPathGenerator{}
}

func (f *FullPathGenerator) Generate(service string, name string) string {
	path := strings.ReplaceAll(service, "\\", "/")
	if len(path) > 0 && path[0] != '/' {
		path = "/" + path
	}
	return path + "/" + name
}

// ================================================================================================
// IDGenerator generates the id of the request
// ================================================================================================

var DefaultIDGenerator IDGenerator = NewUUIDGenerator()

type IDGenerator interface {
	Generate() string
}

type UUIDGenerator struct{}

func NewUUIDGenerator() *UUIDGenerator {
	return &UUIDGenerator{}
}

func (u *UUIDGenerator) Generate() string {
	return uuid.New().String()
}
