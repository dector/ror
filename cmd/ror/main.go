package main

import (
	"fmt"
	"os"
	"slices"

	"github.com/dector/ror/internal"
	"github.com/dector/ror/internal/commands"
	"github.com/dector/ror/internal/compat"
	"github.com/dector/ror/internal/config"
	"github.com/dector/ror/internal/env"
	"github.com/dector/ror/internal/task"
	"github.com/dector/ror/internal/utils"
)

var configFile = "ror.kdl"

func main() {
	env, args := buildEnvAndArgs()

	ctx := createContext(env, args)
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

	parsed := internal.ParsedArgs{
		RorArgs:  rorArgs,
		TaskName: taskName,
		TaskArgs: taskArgs,
	}

	return parsed
}

func buildEnvAndArgs() (internal.Env, internal.ParsedArgs) {
	env := internal.Env{}

	osArgs := utils.SubSlice(os.Args, 1)
	args := parseArguments(osArgs)

	// TODO use count
	if slices.Contains(args.RorArgs, "-v") {
		env.VerboseOutput = true
	}
	if slices.Contains(args.RorArgs, "-vvv") {
		env.VeryVerboseOutput = true
	}

	if env.VeryVerboseOutput {
		fmt.Printf("[DEBUG] Raw arguments: %v\n", osArgs)
		fmt.Printf("[DEBUG] Parsed ror args: %v\n", args.RorArgs)
		fmt.Printf("[DEBUG] Parsed task name: %s\n", args.TaskName)
		fmt.Printf("[DEBUG] Parsed task args: %v\n", args.TaskArgs)
	}

	return env, args
}

func createContext(env internal.Env, args internal.ParsedArgs) internal.Context {
	project := buildProject(env)
	ctx := internal.Context{
		Env:     env,
		Project: project,
		Args:    args,
	}

	return ctx
}

func buildProject(env internal.Env) task.Project {
	// TODO use it
	_ = env

	// Check if ror.kdl exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// ror.kdl not found, check for Taskfile.yml
		if config.CheckTaskfileExists() {
			compat.RunTaskfile(utils.SubSlice(os.Args, 1))
			os.Exit(0)
		}
		fmt.Println("ror.kdl not found")
		os.Exit(1)
	}

	return config.ParseProject(configFile)
}

func execute(ctx internal.Context) {
	//fmt.Printf("Execute: %+v\n", ctx)

	if ctx.Args.TaskName == "" {
		commands.CmdListTasks(ctx)
		return
	}

	switch ctx.Args.TaskName {
	case "version":
		if len(ctx.Args.TaskArgs) > 0 && ctx.Args.TaskArgs[0] == "--short" {
			fmt.Printf("%s\n", env.GetShortVersion())
		} else if len(ctx.Args.TaskArgs) > 0 && ctx.Args.TaskArgs[0] == "--verbose" {
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
