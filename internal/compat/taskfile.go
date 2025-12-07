package compat

import (
	"bytes"
	"fmt"
	"os/exec"
	"sort"

	"github.com/dector/ror/internal/io"
	"github.com/dector/ror/internal/task"
)

func RunTaskfile(io io.IO, args []string) {
	cmd := exec.Command("task", args...)
	cmd.Stdout = io.Std().Stdout()
	cmd.Stderr = io.Std().Stderr()
	cmd.Stdin = io.Std().Stdin()

	if err := io.Shell().RunCommand(cmd); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			io.Std().Exit(exitErr.ExitCode())
		}
		fmt.Fprintf(io.Std().Stderr(), "Error running task: %v\n", err)
		io.Std().Exit(1)
	}
}

func ExportTaskfile(io io.IO, project task.Project) {
	var buf bytes.Buffer

	fmt.Fprintln(&buf, "version: '3'")
	fmt.Fprintln(&buf, "")
	fmt.Fprintln(&buf, "tasks:")

	// Iterate in order as defined in ror.kdl
	for pair := project.Tasks.Oldest(); pair != nil; pair = pair.Next() {
		name := pair.Key
		t := pair.Value
		fmt.Fprintf(&buf, "  %s:\n", name)
		if t.Description != "" {
			fmt.Fprintf(&buf, "    desc: %s\n", t.Description)
		}

		if len(t.DependsOn) > 0 {
			fmt.Fprintln(&buf, "    deps:")
			for _, dep := range t.DependsOn {
				fmt.Fprintf(&buf, "      - %s\n", dep)
			}
		}

		if len(t.EnvVars) > 0 {
			fmt.Fprintln(&buf, "    env:")
			// Sort env vars for deterministic output
			var envKeys []string
			for k := range t.EnvVars {
				envKeys = append(envKeys, k)
			}
			sort.Strings(envKeys)
			for _, k := range envKeys {
				fmt.Fprintf(&buf, "      %s: %s\n", k, t.EnvVars[k])
			}
		}

		cmd := t.Command
		vars := make(map[string]task.WhereVariable)

		if t.CommandTemplate != nil {
			cmd = convertRorToTaskfileVar(t.CommandTemplate.Template)
			collectVariables(t.CommandTemplate.Variables, vars)
		}

		if len(vars) > 0 {
			fmt.Fprintln(&buf, "    vars:")
			// Sort vars for deterministic output
			var varKeys []string
			for k := range vars {
				varKeys = append(varKeys, k)
			}
			sort.Strings(varKeys)

			for _, k := range varKeys {
				v := vars[k]
				val := convertRorToTaskfileVar(v.Value)
				if v.Type == task.WhereTypeCmd {
					fmt.Fprintf(&buf, "      %s:\n", k)
					fmt.Fprintf(&buf, "        sh: %s\n", val)
				} else {
					fmt.Fprintf(&buf, "      %s: %s\n", k, fmt.Sprintf("%q", val))
				}
			}
		}

		if cmd != "" {
			fmt.Fprintln(&buf, "    cmds:")
			fmt.Fprintf(&buf, "      - %s\n", cmd)
		}
		fmt.Fprintln(&buf, "")
	}

	err := io.Files().WriteFile("Taskfile.yml", buf.Bytes(), 0644)
	if err != nil {
		fmt.Fprintf(io.Std().Stderr(), "Error writing Taskfile.yml: %v\n", err)
		io.Std().Exit(1)
	}
}

func convertRorToTaskfileVar(input string) string {
	return task.VariableRegexp.ReplaceAllString(input, `{{.$1}}`)
}

func collectVariables(vars []task.WhereVariable, acc map[string]task.WhereVariable) {
	for _, v := range vars {
		// Add current variable
		// Note: Ror allows nested scopes, but Taskfile vars are flat per task.
		// We overwrite existing keys, assuming the last one in traversal (or depth-first) is intended,
		// though Ror's specific scoping rules might be more complex.
		// For export, a flat map is a reasonable approximation.
		acc[v.Name] = v
		collectVariables(v.Children, acc)
	}
}
