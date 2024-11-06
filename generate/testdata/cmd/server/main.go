package main

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner"
)

func main() {
	gone.
		Default.
		LoadPriest(func(cemetery gone.Cemetery) error {
			return goner.GinPriest(cemetery) // register gin priest
		}).
		Serve()
}
