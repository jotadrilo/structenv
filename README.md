# structenv

Set structs from environment variables in Golang.

## Example

```golang
package main

import (
	"fmt"
	"log"

	"github.com/jotadrilo/structenv"
)
type Person struct {
	Name   string  `env:"NAME"`
	Age    int     `env:"AGE"`
	Height float64 `env:"HEIGHT"`
	Weight float64 `env:"WEIGHT"`
}

func main() {
    p := &Person{Name: "Bob", Age: 45, Height: 1.75, Weight: 76.8}
	fmt.Printf("#1 %s (%.1fkg, %.1fm) is %d years old\n", p.Name, p.Weight, p.Height, p.Age)

	if err := structenv.Parse(p); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("#2 %s (%.1fkg, %.1fm) is %d years old\n", p.Name, p.Weight, p.Height, p.Age)
}
```

```shell
$ go run ./person.go
#1 Bob (76.800000kg, 1.750000m) is 45 years old
#2 Bob (76.800000kg, 1.750000m) is 45 years old

$ NAME=Alice AGE=45 HEIGHT=1.67 WEIGHT=62 go run ./person.go
#1 Bob (76.8kg, 1.8m) is 45 years old
#2 Alice (62.0kg, 1.7m) is 45 years old
```
