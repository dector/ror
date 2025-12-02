package commands

import (
	"fmt"

	"github.com/dector/ror/internal"
)

func CmdListTasks(ctx internal.Context) {
	if len(ctx.Config.Tasks) == 0 {
		fmt.Println("No tasks found in ror.kdl")
		return
	}

	fmt.Println("Available tasks:")
	for name, task := range ctx.Config.Tasks {
		desc := task.Description
		if desc == "" {
			desc = "(no description)"
		}
		fmt.Printf("  %s: %s\n", name, desc)
		if len(task.DependsOn) > 0 {
			fmt.Printf("    Depends on: %v\n", task.DependsOn)
		}
	}
}
