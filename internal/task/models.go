package task

type Task struct {
	Name        string
	Command     string
	Description string
	DependsOn   []string
}

type RunnerConfig struct {
	Tasks map[string]Task
}
