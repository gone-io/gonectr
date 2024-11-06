package generate

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func Test_scanDirGenCode(t *testing.T) {
	rel, err := filepath.Rel("abc/d1/d2", "abc/d1/d2/d3")
	assert.Nil(t, err)
	println(rel)
}
