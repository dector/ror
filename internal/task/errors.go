package task

import "fmt"

// TaskNotFoundError reports a requested task that is not defined in the project.
type TaskNotFoundError struct {
	Name string
}

func (e *TaskNotFoundError) Error() string {
	return fmt.Sprintf("task '%s' not found", e.Name)
}
