# Encrypter

## Usage

```go
package main

import "github.com/go-kratos-ecosystem/components/v2/encrypter"

func main() {
	e := encrypter.New("EAFBSPAXDCIOGRUVNERQGXPYGPNKYATM")

	ciphertext, _ := e.Encrypt("test")
	plaintext, _ := e.Decrypt(ciphertext)

	println(ciphertext, plaintext) // I-UVz6tds3MlRX-1afR36cLcmMw= test
}
```