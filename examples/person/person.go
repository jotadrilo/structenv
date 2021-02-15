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
