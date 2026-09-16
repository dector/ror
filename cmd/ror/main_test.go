package main

import (
	"reflect"
	"testing"

	"github.com/dector/ror/internal"
)

func TestParseArgumentsCommandAliases(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want internal.ParsedArgs
	}{
		{
			name: "+init becomes the init command",
			args: []string{"+init"},
			want: internal.ParsedArgs{TaskName: "init", TaskArgs: []string{}},
		},
		{
			name: "+version becomes the version command",
			args: []string{"+version"},
			want: internal.ParsedArgs{TaskName: "version", TaskArgs: []string{}},
		},
		{
			name: "flags before the alias stay as ror flags",
			args: []string{"-v", "+version"},
			want: internal.ParsedArgs{RorArgs: []string{"-v"}, TaskName: "version", TaskArgs: []string{}},
		},
		{
			name: "arguments after the alias go to the command",
			args: []string{"+version", "--short"},
			want: internal.ParsedArgs{TaskName: "version", TaskArgs: []string{"--short"}},
		},
		{
			name: "+list is still a flag, not an alias",
			args: []string{"+list"},
			want: internal.ParsedArgs{RorArgs: []string{"+list"}, TaskArgs: []string{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseArguments(tt.args)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseArguments(%v) = %+v, want %+v", tt.args, got, tt.want)
			}
		})
	}
}
