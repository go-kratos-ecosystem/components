# Prints

## Usage

```go
package main

import (
	"fmt"

	"github.com/go-kratos-ecosystem/components/v2/prints"
)

func main() {
	// basic
	fmt.Println("Hello World!", "Hello Go!")
	prints.Line("Line: Hello World!", "Hello Go!")              //nolint:errcheck
	prints.Linef("Linef: %s %s\n", "Hello World!", "Hello Go!") //nolint:errcheck
	prints.NewLine()                                            //nolint:errcheck
	println("----------------------------------------")
	prints.NewLine(0) //nolint:gomnd,errcheck
	println("----------------------------------------")
	prints.NewLine(1) //nolint:gomnd,errcheck
	println("----------------------------------------")
	prints.NewLine(2) //nolint:gomnd,errcheck
	println("----------------------------------------")
	prints.NewLine(3)                                                 //nolint:gomnd,errcheck
	prints.Info("Info: Hello World!", "Hello Go!")                    //nolint:errcheck
	prints.Infof("Infof: %s %s\n", "Hello World!", "Hello Go!")       //nolint:errcheck
	prints.Comment("Comment: Hello World!", "Hello Go!")              //nolint:errcheck
	prints.Commentf("Commentf: %s %s\n", "Hello World!", "Hello Go!") //nolint:errcheck
	prints.Error("Error: Hello World!", "Hello Go!")                  //nolint:errcheck
	prints.Errorf("Errorf: %s %s\n", "Hello World!", "Hello Go!")     //nolint:errcheck
	prints.Warn("Warn: Hello World!", "Hello Go!")                    //nolint:errcheck
	prints.Warnf("Warnf: %s %s\n", "Hello World!", "Hello Go!")       //nolint:errcheck
	prints.Alert("Alert: Hello World!", "Hello Go!")                  //nolint:errcheck
	prints.Alertf("Alertf: %s %s\n", "Hello World!", "Hello Go!")     //nolint:errcheck

	// prompt
	r1, _ := prints.Ask("What is your name?", "FLC") //nolint:errcheck
	if r1 == "FLC" {
		prints.Warn("You are FLC!") //nolint:errcheck
	} else {
		prints.Infof("You are %s!\n", r1) //nolint:errcheck
	}
	r2, _ := prints.Ask("What is your name?") //nolint:errcheck
	prints.Info(r2)                           //nolint:errcheck
}
```

Output:

![](output.jpg)