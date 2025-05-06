package parser

const codeTpl = `%s
package %s

import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
)

// load installed gone module LoadFunc
var loaders = []gone.LoadFunc{
}

func GoneModuleLoad(loader gone.Loader) error {
	var ops []*g.LoadOp
	for _, f := range loaders {
		ops = append(ops, g.F(f))
	}
	return g.BuildOnceLoadFunc(ops...)(loader)
}
`
