package uuid

import "github.com/google/uuid"

// Generate uuid
func Generate() string {
	return uuid.New().String()
}
