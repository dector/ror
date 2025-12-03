package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/dector/ror/internal"
	"github.com/dector/ror/internal/io"
	taskpkg "github.com/dector/ror/internal/task"
)

type executionState struct {
	executed map[string]bool
	visiting map[string]bool
}

func newExecutionState() *executionState {
	return &executionState{
		executed: make(map[string]bool),
		visiting: make(map[string]bool),
	}
}

func runTaskWithDependencies(io io.IO, ctx internal.Context) error {
	state := newExecutionState()
	if err := runTask(io, ctx.Args.TaskName, ctx.Project, ctx.Args.TaskArgs, state, ctx.Env); err != nil {
		return err
	}

	return nil
}

func runTask(io io.IO, name string, config taskpkg.Project, args []string, state *executionState, env internal.Env) error {
	if state.visiting[name] {
		return fmt.Errorf("circular dependency detected: %s", name)
	}
	if state.executed[name] {
		return nil
	}

	state.visiting[name] = true
	defer func() { delete(state.visiting, name) }()

	task, ok := config.Tasks[name]
	if !ok {
		return fmt.Errorf("task '%s' not found", name)
	}

	if env.VeryVerboseOutput && len(task.DependsOn) > 0 {
		io.Std().Printf("[DEBUG] Task '%s' has dependencies: %v\n", name, task.DependsOn)
	}

	for _, dep := range task.DependsOn {
		if err := runTask(io, dep, config, nil, state, env); err != nil {
			return err
		}
	}

	if task.Command != "" || task.CommandTemplate != nil {
		if env.VerboseOutput {
			io.Std().Printf("Running task: %s\n", name)
		}

		// Determine the final command string
		var finalCommand string
		var err error

		if task.CommandTemplate != nil {
			if env.VeryVerboseOutput {
				io.Std().Printf("[DEBUG] Command command before expansion: %s\n", task.CommandTemplate.Template)
				if len(task.CommandTemplate.Variables) > 0 {
					io.Std().Printf("[DEBUG] Variables to expand:\n")
					for _, v := range task.CommandTemplate.Variables {
						io.Std().Printf("[DEBUG]   - %s = %s (type: %s)\n", v.Name, v.Value, v.Type)
					}
				}
			}

			// Expand variables in template
			finalCommand, err = taskpkg.ExpandCommand(io, task.CommandTemplate, env.VeryVerboseOutput)
			if err != nil {
				return fmt.Errorf("failed to expand command template for task '%s': %w", name, err)
			}

			if env.VerboseOutput {
				io.Std().Printf("Expanded command: %s\n", finalCommand)
			}
		} else {
			// Simple command
			finalCommand = task.Command
			if env.VeryVerboseOutput {
				io.Std().Printf("[DEBUG] Simple command (no expansion needed): %s\n", finalCommand)
			}
		}

		// Build command args array from final command and args
		commandArgs := []string{finalCommand}
		commandArgs = append(commandArgs, args...)

		fullCommand := strings.Join(commandArgs, " ")
		if env.VeryVerboseOutput {
			io.Std().Printf("[DEBUG] Final shell command: sh -c \"%s\"\n", fullCommand)
			if len(args) > 0 {
				io.Std().Printf("[DEBUG] Task arguments: %v\n", args)
			}
		}

		cmd := exec.Command("sh", "-c", fullCommand)
		cmd.Stdout = io.Std().Stdout()
		cmd.Stderr = io.Std().Stderr()
		cmd.Stdin = io.Std().Stdin()

		// Set environment variables from the task
		if taskEnv := buildTaskEnv(io, task, env); taskEnv != nil {
			cmd.Env = taskEnv
		}

		if err := io.Shell().RunCommand(cmd); err != nil {
			return fmt.Errorf("failed to run task '%s': %w", name, err)
		}
	}

	state.executed[name] = true
	return nil
}

// buildTaskEnv builds the environment variables for task execution
// by combining the current environment with task-specific variables
func buildTaskEnv(io io.IO, task taskpkg.Task, env internal.Env) []string {
	if len(task.EnvVars) == 0 {
		return nil
	}

	if env.VeryVerboseOutput {
		io.Std().Printf("[DEBUG] Setting environment variables:\n")
		for k, v := range task.EnvVars {
			io.Std().Printf("[DEBUG]   %s=%s\n", k, v)
		}
	}

	// Get current environment and append task-specific vars
	currentEnv := io.Shell().Environ()
	taskEnvVars := formatEnvVars(task.EnvVars)
	return append(currentEnv, taskEnvVars...)
}

// formatEnvVars converts a map of environment variables to a slice of "KEY=VALUE" strings
func formatEnvVars(envVars map[string]string) []string {
	result := make([]string, 0, len(envVars))
	for k, v := range envVars {
		result = append(result, fmt.Sprintf("%s=%s", k, v))
	}
	return result
}
