package internal

import "github.com/dector/ror/internal/task"

type ParsedArgs struct {
	RorArgs  []string
	TaskName string
	TaskArgs []string
}

type Env struct {
	VerboseOutput     bool
	VeryVerboseOutput bool
}

type Context struct {
	Env     Env
	Project task.Project
	Args    ParsedArgs
}
