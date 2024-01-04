# UDP

## Server

```go
package main

import (
	"log"

	"github.com/go-kratos/kratos/v2"

	"github.com/go-kratos-ecosystem/components/v2/udp"
)

func main() {
	err := kratos.New(
		kratos.Server(
			udp.NewServer(":12190", udp.WithHandler(func(msg *udp.Message) {
				log.Printf("receive message: %s", msg.Body)
			}), udp.WithRecoveryHandler(func(msg *udp.Message, err interface{}) {
				log.Println(err)
			}), udp.WithReadChanSize(10240)),
		),
	).Run()

	if err != nil {
		log.Fatal(err)
	}
}

```