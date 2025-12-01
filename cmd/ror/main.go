package main

import (
	"fmt"
	"os"

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

	// Validation
	if len(config.Tasks) == 0 {
		fmt.Println("No tasks found in ror.kdl")
		return
	}

	// Print tasks
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
