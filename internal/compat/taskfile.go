package compat

import (
	"fmt"
	"os/exec"

	"github.com/dector/ror/internal/io"
)

func RunTaskfile(io io.IO, args []string) {
	cmd := exec.Command("task", args...)
	cmd.Stdout = io.Std().Stdout()
	cmd.Stderr = io.Std().Stderr()
	cmd.Stdin = io.Std().Stdin()

	if err := io.Shell().RunCommand(cmd); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			io.Std().Exit(exitErr.ExitCode())
		}
		fmt.Fprintf(io.Std().Stderr(), "Error running task: %v\n", err)
		io.Std().Exit(1)
	}
}