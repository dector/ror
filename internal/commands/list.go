package commands

import (
	"github.com/dector/ror/internal"
	"github.com/dector/ror/internal/io"
)

func CmdListTasks(io io.IO, ctx internal.Context) {
	if len(ctx.Project.Tasks) == 0 {
		io.Std().Println("No tasks found in ror.kdl")
		return
	}

	io.Std().Println("Available tasks:")
	for name, task := range ctx.Project.Tasks {
		desc := task.Description
		if desc == "" {
			desc = "(no description)"
		}

		io.Std().Printf("  %s: %s\n", name, desc)
		if len(task.DependsOn) > 0 {
			io.Std().Printf("    Depends on: %v\n", task.DependsOn)
		}
	}
}
