package codec

type Codec interface {
	// Marshal the given data into bytes.
	Marshal(data any) ([]byte, error)

	// Unmarshal the given bytes into dest.
	Unmarshal(src []byte, dest any) error
}
