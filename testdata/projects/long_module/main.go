package main

import (
	"time"

	"github.com/gone-io/gone"
)

func main() {
	gone.Default.Run(func(dep struct {
		point Point `gone:"*"`
	}) {
		println(dep.point.Echo())
	})

	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)
		println("hello", i)
	}
}

type Point struct {
	gone.Flag
}

func (p *Point) Echo() string {
	return "I am a point"
}
