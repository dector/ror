package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/dector/ror/internal/task"
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

func runTaskWithDependencies(name string, config task.RunnerConfig, args []string) error {
	state := newExecutionState()
	if err := runTask(name, config, args, state); err != nil {
		return err
	}

	return nil
}

func runTask(name string, config task.RunnerConfig, args []string, state *executionState) error {
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

	for _, dep := range task.DependsOn {
		if err := runTask(dep, config, nil, state); err != nil {
			return err
		}
	}

	if task.Command != "" {
		fmt.Printf("Running task: %s\n", name)

		// Build command args array from task command and args
		commandArgs := []string{task.Command}
		commandArgs = append(commandArgs, args...)

		cmd := exec.Command("sh", "-c", strings.Join(commandArgs, " "))
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
