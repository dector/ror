package config

import (
	"fmt"

	"github.com/dector/kdly"
	"github.com/dector/ror/internal/io"
	"github.com/dector/ror/internal/task"
)

var compatibilityTaskfile = true

func ParseProject(io io.IO, filename string) task.Project {
	doc, err := readProject(io, filename)
	if err != nil {
		panic(err)
	}

	config := task.Project{
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
					cmd, cmdTemplate := parseCmdNode(child)
					if cmdTemplate != nil {
						task.CommandTemplate = cmdTemplate
					} else {
						task.Command = cmd
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

func readProject(io io.IO, filename string) (*kdly.Document, error) {
	content, err := io.Files().ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return kdly.Parse(string(content))
}

func CheckTaskfileExists(io io.IO) bool {
	if !compatibilityTaskfile {
		return false
	}

	_, err := io.Files().Stat("Taskfile.yml")
	return err == nil
}

// parseCmdNode parses a cmd node and returns either a simple command or a CommandTemplate
func parseCmdNode(node kdly.Node) (string, *task.CommandTemplate) {
	if len(node.Arguments) == 0 {
		panic("cmd node requires a command string")
	}

	cmdString := node.Arguments[0].Value

	// Check if there are where blocks (children)
	if len(node.Children) == 0 {
		// Simple command without where blocks
		return cmdString, nil
	}

	// Parse where blocks
	variables := parseWhereBlocks(node.Children)

	return "", &task.CommandTemplate{
		Template:  cmdString,
		Variables: variables,
	}
}

// parseWhereBlocks recursively parses where blocks
func parseWhereBlocks(nodes []kdly.Node) []task.WhereVariable {
	var variables []task.WhereVariable

	for _, node := range nodes {
		if node.Name != "where" {
			panic(fmt.Sprintf("unexpected child in cmd block: %s (expected 'where')", node.Name))
		}

		variable := parseWhereVariable(node)
		variables = append(variables, variable)
	}

	return variables
}

// parseWhereVariable parses a single where variable
func parseWhereVariable(node kdly.Node) task.WhereVariable {
	if len(node.Properties) == 0 {
		panic("where node requires at least one property (the variable definition)")
	}

	// Find the main variable property (the one that's not 'as')
	var varName string
	var varValue string
	found := false

	for _, prop := range node.Properties {
		if prop.Key != "as" {
			varName = prop.Key
			varValue = prop.Value.Value
			found = true
			break
		}
	}

	if !found {
		panic("where node missing variable property (only found 'as')")
	}

	// Check for 'execute' child node to determine type
	varType := task.WhereTypeStr // default
	var children []task.WhereVariable

	if len(node.Children) > 0 {
		// Check if first child is 'execute' node
		if len(node.Children) == 1 && node.Children[0].Name == "execute" {
			varType = task.WhereTypeCmd
		} else {
			// Parse nested where blocks recursively
			children = parseWhereBlocks(node.Children)
		}
	}

	return task.WhereVariable{
		Name:     varName,
		Value:    varValue,
		Type:     varType,
		Children: children,
	}
}
