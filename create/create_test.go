package create

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_processGoMod(t *testing.T) {
	type args struct {
		content string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "single line not local",
			args: args{
				content: "replace github.com/jim-minter/gonectl v0.0.0 => github.com/jim-minter/gonectl v0.0.0",
			},
			want: "replace github.com/jim-minter/gonectl v0.0.0 => github.com/jim-minter/gonectl v0.0.0\n",
		},
		{
			name: "single line local",
			args: args{
				content: "replace github.com/jim-minter/gonectl => ./gonectl",
			},
			want: "",
		},
		{
			name: "multi line",
			args: args{
				content: "replace (\ngithub.com/jim-minter/gonectl => ./gonectl\nreplace github.com/jim-minter/gonectl v0.0.0 => ../x\n)",
			},
			want: "",
		},
		{
			name: "multi line has local and remote",
			args: args{
				content: "replace (\ngithub.com/jim-minter/gonectl => ./gonectl\nreplace github.com/jim-minter/gonectl v0.0.0 => ../x\ngithub.com/jim-minter/gonectl v0.0.0 => github.com/jim-minter/gonectl v0.0.0\n)",
			},
			want: "replace (\ngithub.com/jim-minter/gonectl v0.0.0 => github.com/jim-minter/gonectl v0.0.0\n)\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, processGoMod(tt.args.content), "processGoMod(%v)", tt.args.content)
		})
	}
}
