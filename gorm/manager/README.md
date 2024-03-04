# Manager

## Example

```go
package main

import (
	"gorm.io/gorm"

	"github.com/go-kratos-ecosystem/components/v2/gorm/manager"
)

func main() {
	var db1, db2, db3, db4 *gorm.DB

	m := manager.New(db1)  // register db1 as default connection
	m.Register("db2", db2) // register db2 as named connection
	m.Register("db3", db3) // register db3 as named connection
	m.Register("db4", db4) // register db4 as named connection

	_ = m.DB          // get default connection
	_ = m.Conn()      // get default connection
	_ = m.Conn("db2") // get db2 connection
	_ = m.Conn("db3") // get db3 connection
	_ = m.Conn("db4") // get db4 connection
	_ = m.Conn("db5") // panic: The connection [db5] is not registered.
}
```