package task

import (
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

// Verbosity levels
const (
	VerbositySilent      = -2
	VerbosityQuiet       = -1
	VerbosityNormal      = 0
	VerbosityVerbose     = 1
	VerbosityVeryVerbose = 2
	VerbosityDebug       = 3
)

// WhereType represents the type of a where variable
type WhereType string

const (
	WhereTypeStr WhereType = "str"
	WhereTypeCmd WhereType = "cmd"
)

// WhereVariable represents a single variable definition in a where block
type WhereVariable struct {
	Name     string          // Variable name (e.g., "greeting")
	Value    string          // Value or command to execute
	Type     WhereType       // Type: str or cmd
	Children []WhereVariable // Nested where blocks
}

// CommandTemplate represents a command with its variable definitions
type CommandTemplate struct {
	Template  string          // Command template with %%variable%% placeholders
	Variables []WhereVariable // Top-level where variables
}

type Task struct {
	Name            string
	Command         string           // Keep for backward compatibility (simple commands)
	CommandTemplate *CommandTemplate // New: for commands with where blocks
	Description     string
	DependsOn       []string
	EnvVars         map[string]string // Environment variables for the task
}

type Project struct {
	Tasks       *orderedmap.OrderedMap[string, Task]
	DefaultTask string   // Task to run when no task is specified
	DefaultArgs []string // Arguments to pass to the default task
}
