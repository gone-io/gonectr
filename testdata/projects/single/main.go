package main

import "github.com/gone-io/gone"

func main() {
	gone.Default.Run(func(dep struct {
		point Point `gone:"*"`
	}) {
		println(dep.point.Echo())
	})
}

type Point struct {
	gone.Flag
}

func (p *Point) Echo() string {
	return "I am a point"
}
