# Tracing

This middleware provides tracing capabilities for your application. 

It is based on the OpenTelmetry specification and can be used with any OpenTelemetry compatible tracing backend.

## Usage Example

```go
package main

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/go-kratos-ecosystem/components/v2/gin/middleware/tracing"
)

func main() {
	router := gin.Default()
	router.Use(tracing.New())
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	if err := router.Run(":8000"); err != nil {
		panic(err)
	}
}
```