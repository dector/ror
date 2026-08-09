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
	"github.com/fatih/color"
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
	env := internal.Env{
		VerbosityLevel: task.VerbosityNormal,
	}

	osArgs := utils.SubSlice(io.Std().Args(), 1)
	args := parseArguments(osArgs)

	// Find out verbosity
	vLevel := 0
	for _, arg := range args.RorArgs {
		switch arg {
		case "-v":
			vLevel++
		case "-vv":
			vLevel += 2
		case "-vvv":
			vLevel += 3
		case "-q", "--quiet":
			vLevel -= 1
		case "-s", "--silent":
			vLevel -= 2
		}
	}
	if vLevel >= 3 {
		env.VerbosityLevel = task.VerbosityDebug
	} else if vLevel == 2 {
		env.VerbosityLevel = task.VerbosityVeryVerbose
	} else if vLevel == 1 {
		env.VerbosityLevel = task.VerbosityVerbose
	} else if vLevel == 0 {
		env.VerbosityLevel = task.VerbosityNormal
	} else if vLevel == -1 {
		env.VerbosityLevel = task.VerbosityQuiet
	} else if vLevel <= -2 {
		env.VerbosityLevel = task.VerbositySilent
	}

	if env.VerbosityLevel >= task.VerbosityDebug {
		io.Std().Printf("[DEBUG] Raw arguments: %v\n", osArgs)
		io.Std().Printf("[DEBUG] Parsed ror args: %v\n", args.RorArgs)
		io.Std().Printf("[DEBUG] Parsed task name: %s\n", args.TaskName)
		io.Std().Printf("[DEBUG] Parsed task args: %v\n", args.TaskArgs)
		io.Std().Printf("[DEBUG] Verbosity level: %d\n", env.VerbosityLevel)
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
		red := color.New(color.FgRed, color.Bold).SprintFunc()
		fmt.Fprintf(io.Std().Stderr(), "%s ror.kdl not found\n", red("Error:"))
		io.Std().Exit(1)
	}

	return config.ParseProject(io, configFile)
}

func execute(io io.IO, ctx internal.Context) {
	//fmt.Printf("Execute: %+v\n", ctx)

	if slices.Contains(ctx.Args.RorArgs, "--export-taskfile") {
		if config.CheckTaskfileExists(io) {
			red := color.New(color.FgRed, color.Bold).SprintFunc()
			fmt.Fprintf(io.Std().Stderr(), "%s Taskfile.yml already exists\n", red("Error:"))
			io.Std().Exit(1)
		}
		compat.ExportTaskfile(io, ctx.Project)
		return
	}

	// Handle +list flag explicitly
	if slices.Contains(ctx.Args.RorArgs, "+list") {
		commands.CmdListTasks(io, ctx)
		return
	}

	// If no task is specified, use default task if available
	if ctx.Args.TaskName == "" {
		if ctx.Env.VerbosityLevel >= task.VerbosityDebug {
			fmt.Fprintf(io.Std().Stderr(), "[DEBUG] Default task: '%s'\n", ctx.Project.DefaultTask)
		}
		if ctx.Project.DefaultTask != "" {
			// Use the default task
			ctx.Args.TaskName = ctx.Project.DefaultTask
			// If no task args were provided and default args are defined, use them
			if len(ctx.Args.TaskArgs) == 0 && len(ctx.Project.DefaultArgs) > 0 {
				ctx.Args.TaskArgs = ctx.Project.DefaultArgs
				if ctx.Env.VerbosityLevel >= task.VerbosityDebug {
					fmt.Fprintf(io.Std().Stderr(), "[DEBUG] Using default args: %v\n", ctx.Args.TaskArgs)
				}
			}
			if ctx.Env.VerbosityLevel >= task.VerbosityDebug {
				fmt.Fprintf(io.Std().Stderr(), "[DEBUG] Using default task: '%s'\n", ctx.Args.TaskName)
			}
		} else {
			// No default task, list tasks
			commands.CmdListTasks(io, ctx)
			return
		}
	}

	if ctx.Env.VerbosityLevel >= task.VerbosityDebug {
		fmt.Fprintf(io.Std().Stderr(), "[DEBUG] Choosen task: '%s'\n", ctx.Args.TaskName)
	}
	switch ctx.Args.TaskName {
	case "version":
		// Version command respects verbosity (silent and quiet suppress output)
		if ctx.Env.VerbosityLevel >= task.VerbosityNormal {
			if len(ctx.Args.TaskArgs) > 0 && ctx.Args.TaskArgs[0] == "--short" {
				io.Std().Printf("%s\n", env.GetShortVersion())
			} else if len(ctx.Args.TaskArgs) > 0 && ctx.Args.TaskArgs[0] == "--verbose" {
				io.Std().Printf("%s\n", env.GetLongVersion())
			} else {
				io.Std().Printf("%s\n", env.GetDefaultVersion())
			}
		}
	case "help":
		printUsage(io, ctx.Env)
	default:
		if err := runTaskWithDependencies(io, ctx); err != nil {
			red := color.New(color.FgRed, color.Bold).SprintFunc()
			fmt.Fprintf(io.Std().Stderr(), "%s %v\n", red("Error:"), err)
			io.Std().Exit(1)
		}
	}
}

func printUsage(io io.IO, env internal.Env) {
	// Silent and quiet modes: no output
	if env.VerbosityLevel < task.VerbosityNormal {
		return
	}

	bold := color.New(color.Bold).SprintFunc()
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	dim := color.New(color.Faint).SprintFunc()

	io.Std().Println(cyan("Usage:"))
	io.Std().Printf("  %s %s %s %s\n",
		bold("ror"),
		green("[flags]"),
		green("[task]"),
		green("[task args]"))
	io.Std().Println("")

	io.Std().Println(cyan("Flags:"))
	io.Std().Printf("  %-25s%s\n", yellow("+list"), "List all available tasks")
	io.Std().Printf("  %-25s%s\n", yellow("--export-taskfile"), "Export ror.kdl to Taskfile.yml")
	io.Std().Printf("  %-25s%s\n", yellow("-s, --silent"), "Silent mode (no output)")
	io.Std().Printf("  %-25s%s\n", yellow("-q, --quiet"), "Quiet mode (errors only)")
	io.Std().Printf("  %-25s%s\n", yellow("-v"), "Verbose output (level 1)")
	io.Std().Printf("  %-25s%s\n", yellow("-vv"), "Very verbose output (level 2)")
	io.Std().Printf("  %-25s%s\n", yellow("-vvv"), "Debug output (level 3)")
	io.Std().Println("")

	io.Std().Println(cyan("Commands:"))
	io.Std().Printf("  %-25s%s\n", yellow("version"), "Print version")
	io.Std().Printf("    %-23s%s\n", green("--short"), "Short version format")
	io.Std().Printf("    %-23s%s\n", green("--verbose"), "Verbose version format")
	io.Std().Printf("  %-25s%s\n", yellow("help"), "Print this help message")
	io.Std().Println("")

	io.Std().Println(cyan("Examples:"))
	io.Std().Printf("  %-25s%s\n", dim("ror"), dim("List all available tasks"))
	io.Std().Printf("  %-25s%s\n", dim("ror build"), dim("Run the 'build' task"))
	io.Std().Printf("  %-25s%s\n", dim("ror -v test"), dim("Run with verbose output"))
}
