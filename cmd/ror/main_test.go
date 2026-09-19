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
			want: internal.ParsedArgs{Command: "init", TaskArgs: []string{}},
		},
		{
			name: "+help becomes the help command",
			args: []string{"+help"},
			want: internal.ParsedArgs{Command: "help", TaskArgs: []string{}},
		},
		{
			name: "+version becomes the version command",
			args: []string{"+version"},
			want: internal.ParsedArgs{Command: "version", TaskArgs: []string{}},
		},
		{
			name: "+validate becomes the validate command with a file argument",
			args: []string{"+validate", "ror.kdl"},
			want: internal.ParsedArgs{Command: "validate", TaskArgs: []string{"ror.kdl"}},
		},
		{
			name: "flags before the alias stay as ror flags",
			args: []string{"-v", "+version"},
			want: internal.ParsedArgs{RorArgs: []string{"-v"}, Command: "version", TaskArgs: []string{}},
		},
		{
			name: "arguments after the alias go to the command",
			args: []string{"+version", "--short"},
			want: internal.ParsedArgs{Command: "version", TaskArgs: []string{"--short"}},
		},
		{
			name: "+list is still a flag, not an alias",
			args: []string{"+list"},
			want: internal.ParsedArgs{RorArgs: []string{"+list"}, TaskArgs: []string{}},
		},
		{
			name: "bare reserved names are regular tasks",
			args: []string{"init"},
			want: internal.ParsedArgs{TaskName: "init", TaskArgs: []string{}},
		},
		{
			name: "bare help is a regular task",
			args: []string{"help", "extra"},
			want: internal.ParsedArgs{TaskName: "help", TaskArgs: []string{"extra"}},
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
