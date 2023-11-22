package serializer

type Serializable interface {
	// Serialize the given data into bytes.
	Serialize(data interface{}) ([]byte, error)

	// Unserialize the given bytes into the given data.
	Unserialize(src []byte, dest interface{}) error
}
