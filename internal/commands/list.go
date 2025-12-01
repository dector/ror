package commands

import (
	"fmt"

	"github.com/dector/ror/internal/task"
)

func CmdListTasks(config task.RunnerConfig) {
	if len(config.Tasks) == 0 {
		fmt.Println("No tasks found in ror.kdl")
		return
	}

	fmt.Println("Available tasks:")
	for name, task := range config.Tasks {
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
