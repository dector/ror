package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/dector/kdly"
)

var configFile = "ror.kdl"

type Task struct {
	Name        string
	Command     string
	Description string
	DependsOn   []string
}

type RunnerConfig struct {
	Tasks map[string]Task
}

func main() {
	config := buildConfig()

	if len(os.Args) < 2 {
		listTasks(config)
		return
	}

	taskName := os.Args[1]
	executed := make(map[string]bool)
	visiting := make(map[string]bool)
	if err := runTask(taskName, config, executed, visiting); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runTask(name string, config RunnerConfig, executed map[string]bool, visiting map[string]bool) error {
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

func listTasks(config RunnerConfig) {
	if len(config.Tasks) == 0 {
		fmt.Println("No tasks found in ror.kdl")
		return
	}

	fmt.Println("Available tasks:")
	for name, task := range config.Tasks {
		desc := task.Description
		if desc == "" {
			desc = "(no description)"
		}
		fmt.Printf("  %s: %s\n", name, desc)
		if len(task.DependsOn) > 0 {
			fmt.Printf("    Depends on: %v\n", task.DependsOn)
		}
	}
}

func buildConfig() RunnerConfig {
	doc, err := parseConfig()
	if err != nil {
		panic(err)
	}

	config := RunnerConfig{
		Tasks: make(map[string]Task),
	}

	// Iterate over top-level nodes
	for _, node := range doc.Nodes {
		if node.Name != "task" {
			panic(fmt.Sprintf("unknown top-level node: %s", node.Name))
		}

		if len(node.Arguments) < 1 {
			panic("task node requires a name argument")
		}

		taskName := node.Arguments[0].Value

		task := Task{
			Name:      taskName,
			DependsOn: []string{},
		}

		if node.Children != nil {
			for _, child := range node.Children {
				switch child.Name {
				case "cmd":
					if len(child.Arguments) > 0 {
						task.Command = child.Arguments[0].Value
					}
				case "description":
					if len(child.Arguments) > 0 {
						task.Description = child.Arguments[0].Value
					}
				case "depends":
					if child.Children != nil {
						for _, dep := range child.Children {
							if dep.Name == "on" && len(dep.Arguments) > 0 {
								task.DependsOn = append(task.DependsOn, dep.Arguments[0].Value)
							}
						}
					}
				default:
					panic(fmt.Sprintf("unknown property in task '%s': %s", taskName, child.Name))
				}
			}
		}

		config.Tasks[taskName] = task
	}

	return config
}

func parseConfig() (*kdly.Document, error) {
	content, err := os.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("ror.kdl not found")
			os.Exit(1)
		}
		return nil, err
	}

	return kdly.Parse(string(content))
}
