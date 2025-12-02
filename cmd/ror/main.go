package main

import (
	"fmt"
	"slices"

	"github.com/dector/ror/internal"
	"github.com/dector/ror/internal/commands"
	"github.com/dector/ror/internal/compat"
	"github.com/dector/ror/internal/config"
	"github.com/dector/ror/internal/env"
	"github.com/dector/ror/internal/io"
	"github.com/dector/ror/internal/task"
	"github.com/dector/ror/internal/utils"
)

var configFile = "ror.kdl"

func main() {
	io := io.NewIO()
	execMain(io)
}

func execMain(io io.IO) {
	env, args := buildEnvAndArgs(io)

	ctx := createContext(io, env, args)

	execute(io, ctx)
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

func buildEnvAndArgs(io io.IO) (internal.Env, internal.ParsedArgs) {
	env := internal.Env{}

	osArgs := utils.SubSlice(io.Std().Args(), 1)
	args := parseArguments(osArgs)

	// TODO use count
	if slices.Contains(args.RorArgs, "-v") {
		env.VerboseOutput = true
	}
	if slices.Contains(args.RorArgs, "-vvv") {
		env.VeryVerboseOutput = true
	}

	if env.VeryVerboseOutput {
		io.Std().Printf("[DEBUG] Raw arguments: %v\n", osArgs)
		io.Std().Printf("[DEBUG] Parsed ror args: %v\n", args.RorArgs)
		io.Std().Printf("[DEBUG] Parsed task name: %s\n", args.TaskName)
		io.Std().Printf("[DEBUG] Parsed task args: %v\n", args.TaskArgs)
	}

	return env, args
}

func createContext(io io.IO, env internal.Env, args internal.ParsedArgs) internal.Context {
	project := buildProject(io, env)
	ctx := internal.Context{
		Env:     env,
		Project: project,
		Args:    args,
	}

	return ctx
}

func buildProject(io io.IO, env internal.Env) task.Project {
	// TODO use it
	_ = env

	// Check if ror.kdl exists
	if _, err := io.Files().Stat(configFile); io.Files().IsNotExist(err) {
		// ror.kdl not found, check for Taskfile.yml
		if config.CheckTaskfileExists(io) {
			compat.RunTaskfile(io, utils.SubSlice(io.Std().Args(), 1))
			io.Std().Exit(0)
		}
		io.Std().Println("ror.kdl not found")
		io.Std().Exit(1)
	}

	return config.ParseProject(io, configFile)
}

func execute(io io.IO, ctx internal.Context) {
	//fmt.Printf("Execute: %+v\n", ctx)

	if slices.Contains(ctx.Args.RorArgs, "--export-taskfile") {
		if config.CheckTaskfileExists(io) {
			fmt.Fprintln(io.Std().Stderr(), "Error: Taskfile.yml already exists")
			io.Std().Exit(1)
		}
		compat.ExportTaskfile(io, ctx.Project)
		return
	}

	if ctx.Args.TaskName == "" {
		commands.CmdListTasks(io, ctx)
		return
	}

	switch ctx.Args.TaskName {
	case "version":
		if len(ctx.Args.TaskArgs) > 0 && ctx.Args.TaskArgs[0] == "--short" {
			io.Std().Printf("%s\n", env.GetShortVersion())
		} else if len(ctx.Args.TaskArgs) > 0 && ctx.Args.TaskArgs[0] == "--verbose" {
			io.Std().Printf("%s\n", env.GetLongVersion())
		} else {
			io.Std().Printf("%s\n", env.GetDefaultVersion())
		}
	case "help":
		printUsage(io)
	default:
		if err := runTaskWithDependencies(io, ctx); err != nil {
			fmt.Fprintf(io.Std().Stderr(), "Error: %v\n", err)
			io.Std().Exit(1)
		}
	}
}

func printUsage(io io.IO) {
	io.Std().Println("Usage:")
	io.Std().Println("  ror [flags] [task] [task args]")
	io.Std().Println("")
	io.Std().Println("Flags:")
	io.Std().Println("  --export-taskfile   Export ror.kdl to Taskfile.yml")
	io.Std().Println("  -v                  Verbose output")
	io.Std().Println("  -vvv                Very verbose output")
	io.Std().Println("")
	io.Std().Println("Commands:")
	io.Std().Println("  version             Print version")
	io.Std().Println("  help                Print this help message")
}
