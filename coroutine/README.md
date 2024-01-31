# Coroutine

## Usage

```go
package main

import (
	"log"

	"github.com/go-kratos-ecosystem/components/v2/coroutine"
)

func main() {
	funcs := []func(){
		func() {
			log.Println("1")
		},
		func() {
			log.Println("2")
		},
		func() {
			log.Println("3")
		},
	}

	// Concurrent Example1
	c := coroutine.NewConcurrent(2)
	defer c.Close()
	c.Add(funcs...)

	c.Wait()

	// Concurrent Example2
	coroutine.RunConcurrent(2, funcs...)

	// Parallel Example1
	p := coroutine.NewParallel()
	p.Add(funcs...)
	p.Wait()

	// Parallel Example2
	coroutine.RunParallel(funcs...)
	
	// Wait Example
	coroutine.Wait(funcs...)
}
```