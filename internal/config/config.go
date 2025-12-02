package config

import (
	"fmt"
	"os"

	"github.com/dector/kdly"
	"github.com/dector/ror/internal/task"
)

var compatibilityTaskfile = true

func ParseConfig(filename string) task.RunnerConfig {
	doc, err := readConfig(filename)
	if err != nil {
		panic(err)
	}

	config := task.RunnerConfig{
		Tasks: make(map[string]task.Task),
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

		task := task.Task{
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

func readConfig(filename string) (*kdly.Document, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return kdly.Parse(string(content))
}

func CheckTaskfileExists() bool {
	if !compatibilityTaskfile {
		return false
	}

	_, err := os.Stat("Taskfile.yml")
	return err == nil
}
