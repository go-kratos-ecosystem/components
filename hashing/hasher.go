package hashing

type Hasher interface {
	// Make a hash value from the given value.
	Make(value string) (string, error)

	// MustMake a hash value from the given value.
	// If an error occurs, it will panic.
	MustMake(value string) string

	// Check the given value matches the given hashed value.
	// If you Make() is error, it will return false.
	Check(value, hashedValue string) bool
}
