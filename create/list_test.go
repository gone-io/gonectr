package create

import (
	"github.com/stretchr/testify/assert"
	"reflect"
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

func Test_getReadmeDesc(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"success",
			args{
				filename: "testdata/test.md",
			},
			"Apollo Configuration Center Example",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, getReadmeDesc(tt.args.filename), "getReadmeDesc(%v)", tt.args.filename)
		})
	}
}

func Test_listDirModule1(t *testing.T) {
	module, err := listDirModule("testdata/projs")
	assert.Nil(t, err)
	assert.True(t, reflect.DeepEqual(module, []string{"proj", "tree/proj2", "tree/proj3"}))
}
