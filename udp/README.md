# UDP

## Server

```go
package main

import (
	"log"
	"net"

	"github.com/go-kratos/kratos/v2"

	"github.com/go-packagist/go-kratos-components/udp"
)

func main() {
	err := kratos.New(
		kratos.Server(
			udp.NewServer(":12190", udp.WithHandler(func(conn net.PacketConn, buf []byte, addr net.Addr) {
				log.Println(string(buf))
			}), udp.WithRecoveryHandler(func(conn net.PacketConn, buf []byte, addr net.Addr, err interface{}) {
				log.Println(err)
			})),
		),
	).Run()

	if err != nil {
		log.Fatal(err)
	}
}

```