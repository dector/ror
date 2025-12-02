package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/dector/ror/internal"
	"github.com/dector/ror/internal/commands"
	"github.com/dector/ror/internal/config"
	"github.com/dector/ror/internal/env"
	"github.com/dector/ror/internal/task"
	"github.com/dector/ror/internal/utils"
)

var configFile = "ror.kdl"

func main() {
	// Check if ror.kdl exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// ror.kdl not found, check for Taskfile.yml
		if config.CheckTaskfileExists() {
			runTaskfile(os.Args[1:])
			return
		}
		fmt.Println("ror.kdl not found")
		os.Exit(1)
	}

	cfg := config.ParseConfig(configFile)

	// fmt.Printf("%+v\n", cfg)

	args := parseArguments(utils.SubSlice(os.Args, 1))
	ctx := createContext(cfg, args)
	execute(ctx)
}

func parseArguments(args []string) internal.ParsedArgs {
	var rorArgs []string
	var taskName string
	var taskArgs []string

	i := 0
	// Collect all arguments starting with '-', '+', or '--' before any task name
	for i < len(args) {
		arg := args[i]
		if len(arg) > 0 && (arg[0] == '-' || arg[0] == '+') {
			rorArgs = append(rorArgs, arg)
			i++
		} else {
			// Found the task name
			taskName = arg
			i++
			break
		}
	}

	// Collect remaining arguments as task arguments
	taskArgs = utils.SubSlice(args, i)

	return internal.ParsedArgs{
		RorArgs:  rorArgs,
		TaskName: taskName,
		TaskArgs: taskArgs,
	}
}

func createContext(config task.RunnerConfig, args internal.ParsedArgs) internal.Context {
	ctx := internal.Context{
		VerboseOutput: false,

		Config: config,
		Args:   args,
	}

	for _, arg := range args.RorArgs {
		if arg == "-v" {
			ctx.VerboseOutput = true
		}
	}

	return ctx
}

func runTaskfile(args []string) {
	cmd := exec.Command("task", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintf(os.Stderr, "Error running task: %v\n", err)
		os.Exit(1)
	}
}

func execute(ctx internal.Context) {
	if ctx.Args.TaskName == "" {
		commands.CmdListTasks(ctx)
		return
	}

	switch ctx.Args.TaskName {
	case "version":
		if len(ctx.Args.TaskArgs) > 0 && ctx.Args.TaskArgs[0] == "--short" {
			fmt.Printf("%s\n", env.GetShortVersion())
		} else if len(ctx.Args.TaskArgs) > 0 && ctx.Args.TaskArgs[0] == "--long" {
			fmt.Printf("%s\n", env.GetLongVersion())
		} else {
			fmt.Printf("%s\n", env.GetDefaultVersion())
		}
	case "help":
		printUsage()
	default:
		if err := runTaskWithDependencies(ctx); err != nil {
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
