package cmd

import (
	"fmt"
	"time"

	"snip/pkg/snippet"
	"snip/pkg/storage"

	"github.com/google/uuid"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Add a new snippet",
	Long:  `Add a new command snippet with optional placeholders and description.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runNew()
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}

func runNew() error {
	// Prompt for command
	commandPrompt := promptui.Prompt{
		Label: "Command",
	}

	command, err := commandPrompt.Run()
	if err != nil {
		return fmt.Errorf("failed to get command: %w", err)
	}

	if command == "" {
		return fmt.Errorf("command cannot be empty")
	}

	// Prompt for description (optional)
	descriptionPrompt := promptui.Prompt{
		Label:   "Description (optional)",
		Default: "",
	}

	description, err := descriptionPrompt.Run()
	if err != nil {
		return fmt.Errorf("failed to get description: %w", err)
	}

	// Create snippet
	newSnippet := snippet.Snippet{
		ID:          uuid.New().String(),
		Command:     command,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Save snippet
	store, err := storage.NewStorage()
	if err != nil {
		return err
	}

	if err := store.Add(newSnippet); err != nil {
		return fmt.Errorf("failed to save snippet: %w", err)
	}

	// Show success message with placeholders if any
	placeholders := newSnippet.ExtractPlaceholders()
	if len(placeholders) > 0 {
		fmt.Printf("\n✓ Snippet added successfully with placeholders: %v\n", placeholders)
	} else {
		fmt.Println("\n✓ Snippet added successfully")
	}

	return nil
}
