package create

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_listDirModule(t *testing.T) {
	module, err := listDirModule("/Users/jim/works/gone-io/goner/examples")
	if err != nil {
		t.Error(err)
	}
	t.Logf("%#v", module)
}

func Test_listExamples(t *testing.T) {
	err := listExamples()
	assert.Nil(t, err)
}
