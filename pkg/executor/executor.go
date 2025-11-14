package executor

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"snip/pkg/snippet"

	"github.com/manifoldco/promptui"
)

// Executor handles command execution with placeholder replacement
type Executor struct{}

// NewExecutor creates a new Executor instance
func NewExecutor() *Executor {
	return &Executor{}
}

// Execute executes a snippet, prompting for placeholders if needed
func (e *Executor) Execute(snip snippet.Snippet) error {
	placeholders := snip.ExtractPlaceholders()

	// If there are placeholders, prompt for values
	values := make(map[string]string)
	if len(placeholders) > 0 {
		fmt.Println("\nPlease provide values for the following placeholders:")
		for _, placeholder := range placeholders {
			prompt := promptui.Prompt{
				Label: placeholder,
			}

			value, err := prompt.Run()
			if err != nil {
				return fmt.Errorf("failed to get placeholder value: %w", err)
			}

			values[placeholder] = value
		}
	}

	// Replace placeholders with actual values
	finalCommand := snip.ReplacePlaceholders(values)

	fmt.Printf("\nExecuting: %s\n\n", finalCommand)

	// Execute the command
	return e.executeShellCommand(finalCommand)
}

// executeShellCommand executes a shell command
func (e *Executor) executeShellCommand(command string) error {
	var cmd *exec.Cmd

	// Determine the shell based on OS
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}

	// Set up command to use current stdin/stdout/stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command execution failed: %w", err)
	}

	return nil
}
