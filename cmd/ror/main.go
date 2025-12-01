package main

import (
	"fmt"
	"os"

	"github.com/dector/ror/internal/commands"
	"github.com/dector/ror/internal/config"
	"github.com/dector/ror/internal/task"
	"github.com/dector/ror/internal/utils"
)

var configFile = "ror.kdl"

func main() {
	config := config.ParseConfig(configFile)

	execute(config, utils.SubSlice(os.Args, 1))
}

func execute(config task.RunnerConfig, args []string) {
	if len(args) == 0 {
		commands.CmdListTasks(config)
		return
	}

	taskName := args[0]
	if err := runTaskWithDependencies(taskName, config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
