package encrypter

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

type Encrypter struct {
	key    []byte
	cipher cipher.Block
}

func New(key string) *Encrypter {
	return &Encrypter{
		key: []byte(key),
	}
}

func (e *Encrypter) Encrypt(value string) (string, error) {
	plaintext := []byte(value)

	block, err := e.getBlock()
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func (e *Encrypter) Decrypt(value string) (string, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(value)
	if err != nil {
		return "", err
	}

	block, err := e.getBlock()
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

func (e *Encrypter) getBlock() (cipher.Block, error) {
	if e.cipher != nil {
		return e.cipher, nil
	}

	c, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}

	e.cipher = c

	return c, nil
}
