package jet

import (
	"strings"

	"github.com/google/uuid"
)

var DefaultPathGenerator PathGenerator = NewPathGenerator()

type PathGenerator interface {
	Generate(service string, name string) string
}

type pathGenerator struct{}

func NewPathGenerator() PathGenerator {
	return &pathGenerator{}
}

func (s *pathGenerator) Generate(service string, name string) string {
	services := strings.Split(service, "\\")
	service = strings.ReplaceAll(services[len(services)-1], "Service", "")
	service = strings.ToLower(service)
	if service[0] != '/' {
		service = "/" + service
	}

	return service + "/" + name
}

type DotPathGenerator struct{}

func (d *DotPathGenerator) Generate(service string, name string) string {
	return service + "." + name
}

type FullPathGenerator struct{}

func (f *FullPathGenerator) Generate(service string, name string) string {
	return service + "/" + name
}

// --------------------------------------------------------------------------------
// IDGenerator
// --------------------------------------------------------------------------------

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
