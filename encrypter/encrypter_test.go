package encrypter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert.Equal(t, &Encrypter{
		key: []byte("test"),
	}, New("test"))
}

func TestEncrypter(t *testing.T) {
	e := New("EAFBSPAXDCIOGRUVNERQGXPYGPNKYATM")

	ciphertext, _ := e.Encrypt("test")
	plaintext, _ := e.Decrypt(ciphertext)

	println(ciphertext, plaintext)

	assert.NotNil(t, ciphertext)
	assert.NotNil(t, plaintext)
	assert.Equal(t, "test", plaintext)
}

func TestEncrypter_Error(t *testing.T) {
	e := New("test")

	_, err1 := e.Encrypt("test")
	assert.Error(t, err1)

	_, err2 := e.Decrypt("test")
	assert.Error(t, err2)
}

func TestEncrypter_Decrypt_Error(t *testing.T) {
	e := New("EAFBSPAXDCIOGRUVNERQGXPYGPNKYATM")

	_, err1 := e.Decrypt("j9_mcZXKVlInk8bbpBqJOpmDp")
	assert.Error(t, err1)

	_, err2 := e.Decrypt("MQ==")
	assert.Error(t, err2)
}
