package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/dector/ror/internal"
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

func runTaskWithDependencies(ctx internal.Context) error {
	state := newExecutionState()
	if err := runTask(ctx.Args.TaskName, ctx.Project, ctx.Args.TaskArgs, state, ctx.Env); err != nil {
		return err
	}

	return nil
}

func runTask(name string, config taskpkg.Project, args []string, state *executionState, env internal.Env) error {
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
		fmt.Printf("[DEBUG] Task '%s' has dependencies: %v\n", name, task.DependsOn)
	}

	for _, dep := range task.DependsOn {
		if err := runTask(dep, config, nil, state, env); err != nil {
			return err
		}
	}

	if task.Command != "" || task.CommandTemplate != nil {
		if env.VerboseOutput {
			fmt.Printf("Running task: %s\n", name)
		}

		// Determine the final command string
		var finalCommand string
		var err error

		if task.CommandTemplate != nil {
			if env.VeryVerboseOutput {
				fmt.Printf("[DEBUG] Command command before expansion: %s\n", task.CommandTemplate.Template)
				if len(task.CommandTemplate.Variables) > 0 {
					fmt.Printf("[DEBUG] Variables to expand:\n")
					for _, v := range task.CommandTemplate.Variables {
						fmt.Printf("[DEBUG]   - %s = %s (type: %s)\n", v.Name, v.Value, v.Type)
					}
				}
			}

			// Expand variables in template
			finalCommand, err = taskpkg.ExpandCommand(task.CommandTemplate, env.VeryVerboseOutput)
			if err != nil {
				return fmt.Errorf("failed to expand command template for task '%s': %w", name, err)
			}

			if env.VerboseOutput {
				fmt.Printf("Expanded command: %s\n", finalCommand)
			}
		} else {
			// Simple command
			finalCommand = task.Command
			if env.VeryVerboseOutput {
				fmt.Printf("[DEBUG] Simple command (no expansion needed): %s\n", finalCommand)
			}
		}

		// Build command args array from final command and args
		commandArgs := []string{finalCommand}
		commandArgs = append(commandArgs, args...)

		fullCommand := strings.Join(commandArgs, " ")
		if env.VeryVerboseOutput {
			fmt.Printf("[DEBUG] Final shell command: sh -c \"%s\"\n", fullCommand)
			if len(args) > 0 {
				fmt.Printf("[DEBUG] Task arguments: %v\n", args)
			}
		}

		cmd := exec.Command("sh", "-c", fullCommand)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to run task '%s': %w", name, err)
		}
	}

	state.executed[name] = true
	return nil
}
