package task

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
	Tasks map[string]Task
}
