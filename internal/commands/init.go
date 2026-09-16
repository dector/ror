package commands

import (
	_ "embed"
	"fmt"

	"github.com/dector/ror/internal"
	"github.com/dector/ror/internal/io"
	"github.com/dector/ror/internal/task"
	"github.com/fatih/color"
)

//go:embed init_template.kdl
var initTemplate string

// CmdInit creates a new ror.kdl template in the current directory.
func CmdInit(io io.IO, env internal.Env, filename string) {
	red := color.New(color.FgRed, color.Bold).SprintFunc()

	if _, err := io.Files().Stat(filename); err == nil {
		fmt.Fprintf(io.Std().Stderr(), "%s %s already exists\n", red("Error:"), filename)
		io.Std().Exit(1)
		return
	}

	if err := io.Files().WriteFile(filename, []byte(initTemplate), 0o644); err != nil {
		fmt.Fprintf(io.Std().Stderr(), "%s %v\n", red("Error:"), err)
		io.Std().Exit(1)
		return
	}

	if env.VerbosityLevel >= task.VerbosityNormal {
		green := color.New(color.FgGreen).SprintFunc()
		io.Std().Printf("%s %s\n", green("Created"), filename)
	}
}
