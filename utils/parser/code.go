package parser

const codeTpl = `%s
package %s
import (
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
)

//added LoadFunc
var loaders = []gone.LoadFunc{
}

func ThirdGonersLoad() gone.LoadFunc {
	var ops []*g.LoadOp
	for _, f := range loaders {
		ops = append(ops, g.F(f))
	}
	return g.BuildOnceLoadFunc(ops...)
}
`
