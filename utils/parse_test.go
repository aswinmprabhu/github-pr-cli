package utils

import "testing"

func Test_parse(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ssh",
			args: args{
				line: "origin  git@github.com:aswinmprabhu/github-pr-cli.git (fetch)",
			},
			want: "aswinmprabhu/github-pr-cli",
		},
		{
			name: "https",
			args: args{
				line: "origin    https://github.com/aswinmprabhu/github-pr-cli (fetch)",
			},
			want: "aswinmprabhu/github-pr-cli",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parse(tt.args.line); got != tt.want {
				t.Errorf("parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
