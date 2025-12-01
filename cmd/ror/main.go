package main

import (
	"fmt"
	"os"

	"github.com/dector/ror/internal/commands"
	"github.com/dector/ror/internal/config"
	"github.com/dector/ror/internal/env"
	"github.com/dector/ror/internal/task"
	"github.com/dector/ror/internal/utils"
)

var configFile = "ror.kdl"

func main() {
	config := config.ParseConfig(configFile)

	// fmt.Printf("%+v\n", config)

	execute(config, utils.SubSlice(os.Args, 1))
}

func execute(config task.RunnerConfig, args []string) {
	if len(args) == 0 {
		commands.CmdListTasks(config)
		return
	}

	command := args[0]
	switch command {
	case "version":
		if len(args) > 1 && args[1] == "--short" {
			fmt.Printf("%s\n", env.GetShortVersion())
		} else if len(args) > 1 && args[1] == "--long" {
			fmt.Printf("%s\n", env.GetLongVersion())
		} else {
			fmt.Printf("%s\n", env.GetDefaultVersion())
		}
	case "help":
		printUsage()
	default:
		taskName := args[0]
		if err := runTaskWithDependencies(taskName, config); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  ror [command]")
	fmt.Println("")
	fmt.Println("Use `ror version` to get version")
}
