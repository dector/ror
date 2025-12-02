package task

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// ExpandCommand expands a CommandTemplate by resolving all variables
func ExpandCommand(template *CommandTemplate, veryVerbose bool) (string, error) {
	if template == nil {
		return "", fmt.Errorf("nil command template")
	}

	// Build variable scope by processing where blocks bottom-up
	scope := make(map[string]string)

	for _, variable := range template.Variables {
		if err := expandVariable(variable, scope, veryVerbose); err != nil {
			return "", err
		}
	}

	if veryVerbose {
		fmt.Printf("[DEBUG] Variable expansion complete. Scope:\n")
		for k, v := range scope {
			fmt.Printf("[DEBUG]   %%%s%% = %s\n", k, v)
		}
	}

	// Substitute variables in the template
	result := substituteVariables(template.Template, scope, veryVerbose)

	if veryVerbose {
		fmt.Printf("[DEBUG] After variable substitution: %s\n", result)
	}

	return result, nil
}

// expandVariable recursively expands a variable and adds it to scope
func expandVariable(variable WhereVariable, scope map[string]string, veryVerbose bool) error {
	if veryVerbose {
		fmt.Printf("[DEBUG] Expanding variable: %s (initial value: %s, type: %s)\n", variable.Name, variable.Value, variable.Type)
	}

	// First, expand all children (bottom-up)
	childScope := make(map[string]string)
	for _, child := range variable.Children {
		if err := expandVariable(child, childScope, veryVerbose); err != nil {
			return err
		}
	}

	// Substitute child variables in this variable's value
	value := substituteVariables(variable.Value, childScope, veryVerbose)

	if veryVerbose && len(childScope) > 0 {
		fmt.Printf("[DEBUG] Variable '%s' after child substitution: %s\n", variable.Name, value)
	}

	// If type is cmd, execute it
	if variable.Type == WhereTypeCmd {
		if veryVerbose {
			fmt.Printf("[DEBUG] Executing command for variable '%s': %s\n", variable.Name, value)
		}
		output, err := executeCommand(value)
		if err != nil {
			return fmt.Errorf("failed to execute command for variable '%s': %w", variable.Name, err)
		}
		value = strings.TrimSpace(output)
		if veryVerbose {
			fmt.Printf("[DEBUG] Command output for '%s': %s\n", variable.Name, value)
		}
	}

	// Check for undefined variables in the expanded value
	undefinedVars := findUndefinedVariables(value)
	for _, undefinedVar := range undefinedVars {
		fmt.Fprintf(os.Stderr, "Warning: undefined variable '%%%s%%' in variable '%s'\n", undefinedVar, variable.Name)
	}

	// Add to scope
	scope[variable.Name] = value

	if veryVerbose {
		fmt.Printf("[DEBUG] Variable '%s' final value: %s\n", variable.Name, value)
	}

	return nil
}

// executeCommand executes a shell command and returns stdout
func executeCommand(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// substituteVariables replaces %%variable%% placeholders in text
func substituteVariables(text string, scope map[string]string, veryVerbose bool) string {
	// Regex to find %%variable%% patterns
	re := regexp.MustCompile(`%%([a-zA-Z0-9_-]+)%%`)

	result := re.ReplaceAllStringFunc(text, func(match string) string {
		// Extract variable name (remove %% markers)
		varName := match[2 : len(match)-2]

		// Look up in scope
		if value, ok := scope[varName]; ok {
			if veryVerbose {
				fmt.Printf("[DEBUG] Substituting %%%s%% with: %s\n", varName, value)
			}
			return value
		}

		// Variable not defined - warn and leave as-is
		fmt.Fprintf(os.Stderr, "Warning: undefined variable '%%%s%%', leaving as-is\n", varName)
		return match
	})

	return result
}

// findUndefinedVariables finds all %%variable%% patterns in text
func findUndefinedVariables(text string) []string {
	re := regexp.MustCompile(`%%([a-zA-Z0-9_-]+)%%`)
	matches := re.FindAllStringSubmatch(text, -1)

	var variables []string
	for _, match := range matches {
		if len(match) > 1 {
			variables = append(variables, match[1])
		}
	}

	return variables
}
