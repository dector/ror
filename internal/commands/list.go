package commands

import (
	"github.com/dector/ror/internal"
	"github.com/dector/ror/internal/io"
	"github.com/dector/ror/internal/task"
	"github.com/fatih/color"
)

func CmdListTasks(io io.IO, ctx internal.Context) {
	// Silent and quiet modes: no output
	if ctx.Env.VerbosityLevel < task.VerbosityNormal {
		return
	}

	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	dim := color.New(color.Faint).SprintFunc()

	if ctx.Project.Tasks.Len() == 0 {
		io.Std().Println("No tasks found in ror.kdl")
		return
	}

	io.Std().Println(cyan("Project tasks:"))
	for pair := ctx.Project.Tasks.Oldest(); pair != nil; pair = pair.Next() {
		name := pair.Key
		task := pair.Value
		desc := task.Description
		if desc == "" {
			desc = dim("(no description)")
		}

		io.Std().Printf("  %-25s%s\n", yellow(name), desc)
		if len(task.DependsOn) > 0 {
			io.Std().Printf("    %s %v\n", green("Depends on:"), task.DependsOn)
		}
	}

	io.Std().Println("")
	io.Std().Println(cyan("Reserved commands:"))
	io.Std().Printf("  %-25sPrint version information\n", yellow("version"))
	io.Std().Printf("  %-25sPrint help message\n", yellow("help"))
	io.Std().Println("")
	io.Std().Println(cyan("Useful flags:"))
	io.Std().Printf("  %-25sList all available tasks\n", yellow("+list"))
	io.Std().Printf("  %-25sSilent mode (no output)\n", yellow("--silent"))
	io.Std().Printf("  %-25sQuiet mode (errors only)\n", yellow("-q, --quiet"))
	io.Std().Printf("  %-25sVerbose output (level 1)\n", yellow("-v"))
	io.Std().Printf("  %-25sVery verbose output (level 2)\n", yellow("-vv"))
	io.Std().Printf("  %-25sDebug output (level 3)\n", yellow("-vvv"))
}
