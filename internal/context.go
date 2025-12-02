package internal

import "github.com/dector/ror/internal/task"

type ParsedArgs struct {
	RorArgs  []string
	TaskName string
	TaskArgs []string
}

type Context struct {
	VerboseOutput bool
	Config        task.RunnerConfig
	Args          ParsedArgs
}
