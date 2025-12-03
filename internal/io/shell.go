package io

import (
	"os"
	"os/exec"
)

type Shell interface {
	ExecuteCommand(cmd *exec.Cmd) ([]byte, error)
	RunCommand(cmd *exec.Cmd) error
	Environ() []string
}

type RealShell struct{}

func (self *RealShell) ExecuteCommand(cmd *exec.Cmd) ([]byte, error) {
	return cmd.Output()
}

func (self *RealShell) RunCommand(cmd *exec.Cmd) error {
	return cmd.Run()
}

func (self *RealShell) Environ() []string {
	return os.Environ()
}
