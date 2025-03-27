package utils

import "testing"

func TestGetGoneVersionFromModuleFile(t *testing.T) {
	type args struct {
		sanDir   []string
		scanFile []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "only v1",
			args: args{
				sanDir:   []string{"testdata/only_v1"},
				scanFile: nil,
			},
			want: "v1",
		},
		{
			name: "only v2",
			args: args{
				sanDir:   []string{"testdata/only_v2"},
				scanFile: nil,
			},
			want: "v2",
		},
		{
			name: "v1 and v2",
			args: args{
				sanDir:   []string{"testdata/v1_and_v2"},
				scanFile: nil,
			},
			want: "v2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetGoneVersionFromModuleFile(tt.args.sanDir, tt.args.scanFile); got != tt.want {
				t.Errorf("GetGoneVersionFromModuleFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
