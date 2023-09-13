package serializer

type Serializable interface {
	Serialize(data interface{}) ([]byte, error)
	Unserialize(src []byte, dest interface{}) error
}
