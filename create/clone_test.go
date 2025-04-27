package create

import (
	"testing"
)

func Test_cloneOrUpdateReop(t *testing.T) {
	_ = cloneOrUpdateReop("testdata/test", "https://github.com/gone-io/gonectr.git")
}
