# helper

## Example

```go
package main

import (
	"fmt"

	"github.com/go-kratos-ecosystem/components/v2/helper"
)

type User struct {
	Name string
	Age  int
}

func main() {
	user := &User{Name: "foo"}

	// Tap
	user = helper.Tap(user, func(u *User) {
		u.Name = "bar"
		u.Age = 18
	})
	fmt.Println(user)
	// output:
	// &{bar 18}

	// With
	user = helper.With(user, func(u *User) *User {
		u.Name = "baz"
		u.Age = 19
		return u
	})
	fmt.Println(user)
	// output:
	// &{baz 19}

	// When
	user = helper.When(user, true, func(u *User) *User {
		u.Name = "Flc"
		u.Age = 20
		return u
	})

	fmt.Println(user)
	// output:
	// &{Flc 20}
}
```