package commands

import (
	"fmt"

	"github.com/dector/ror/internal"
	"github.com/dector/ror/internal/config"
	"github.com/dector/ror/internal/io"
	"github.com/dector/ror/internal/task"
	"github.com/fatih/color"
)

// CmdValidate checks that the given file is a valid ror configuration.
func CmdValidate(io io.IO, env internal.Env, args []string) {
	red := color.New(color.FgRed, color.Bold).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	if len(args) != 1 {
		fmt.Fprintf(io.Std().Stderr(), "%s usage: ror +validate <file>\n", red("Error:"))
		io.Std().Exit(1)
		return
	}

	filename := args[0]
	if err := config.ValidateProject(io, filename); err != nil {
		fmt.Fprintf(io.Std().Stderr(), "%s %s is invalid: %v\n", red("Error:"), filename, err)
		io.Std().Exit(1)
		return
	}

	if env.VerbosityLevel >= task.VerbosityNormal {
		io.Std().Printf("%s %s is valid\n", green("OK:"), filename)
	}
}
