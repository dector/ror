package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/dector/ror/internal/task"
)

func runTaskWithDependencies(name string, config task.RunnerConfig) error {
	executed := make(map[string]bool)
	visiting := make(map[string]bool)
	if err := runTask(name, config, executed, visiting); err != nil {
		return err
	}

	return nil
}

func runTask(name string, config task.RunnerConfig, executed map[string]bool, visiting map[string]bool) error {
	if visiting[name] {
		return fmt.Errorf("circular dependency detected: %s", name)
	}
	if executed[name] {
		return nil
	}

	visiting[name] = true
	defer func() { delete(visiting, name) }()

	task, ok := config.Tasks[name]
	if !ok {
		return fmt.Errorf("task '%s' not found", name)
	}

	for _, dep := range task.DependsOn {
		if err := runTask(dep, config, executed, visiting); err != nil {
			return err
		}
	}

	if task.Command != "" {
		fmt.Printf("Running task: %s\n", name)
		// TODO: Support other shells or direct execution
		cmd := exec.Command("sh", "-c", task.Command)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to run task '%s': %w", name, err)
		}
	}

	executed[name] = true
	return nil
}
