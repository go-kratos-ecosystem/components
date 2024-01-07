# Gin Server

## Usage

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2"

	ginS "github.com/go-kratos-ecosystem/components/v2/gin"
)

func main() {
	gs := ginS.NewServer(
		gin.Default(),
		ginS.WithAddr(":8080"),
	)

	gs.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	app := kratos.New(
		kratos.Server(gs),
	)

	err := app.Run()
	if err != nil {
		panic(err)
	}
}
```