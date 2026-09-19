package internal

import "github.com/dector/ror/internal/task"

type ParsedArgs struct {
	RorArgs  []string
	Command  string // reserved command selected via '+alias' (e.g. "init")
	TaskName string
	TaskArgs []string
}

type Env struct {
	VerbosityLevel int
}

type Context struct {
	Env     Env
	Project task.Project
	Args    ParsedArgs
}
