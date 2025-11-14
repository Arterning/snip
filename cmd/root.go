package cmd

import (
	"fmt"
	"os"

	"snip/pkg/executor"
	"snip/pkg/storage"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "snip",
	Short: "A CLI snippet manager",
	Long: `Snip is a CLI tool for managing and executing command snippets.

Use 'snip new' to add a new snippet.
Use 'snip' to search and execute snippets.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSearch()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runSearch() error {
	store, err := storage.NewStorage()
	if err != nil {
		return err
	}

	// Interactive search
	searcher := &snippetSearcher{
		storage: store,
	}

	return searcher.Search()
}

type snippetSearcher struct {
	storage *storage.Storage
}

func (s *snippetSearcher) Search() error {
	// Custom search template
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "▸ {{ .Command | cyan }} {{ if .Description }}({{ .Description | faint }}){{ end }}",
		Inactive: "  {{ .Command }} {{ if .Description }}({{ .Description | faint }}){{ end }}",
		Selected: "▸ {{ .Command | green }}",
		Details: `
--------- Snippet Details ----------
{{ "Command:" | faint }}	{{ .Command }}
{{ if .Description }}{{ "Description:" | faint }}	{{ .Description }}{{ end }}
{{ if .ExtractPlaceholders }}{{ "Placeholders:" | faint }}	{{ .ExtractPlaceholders }}{{ end }}`,
	}

	// Custom search function
	searcher := func(input string, index int) bool {
		snippets, err := s.storage.Search(input)
		if err != nil {
			return false
		}
		if index >= len(snippets) {
			return false
		}
		return true
	}

	// Load initial snippets
	snippets, err := s.storage.Load()
	if err != nil {
		return err
	}

	if len(snippets) == 0 {
		fmt.Println("No snippets found. Use 'snip new' to add your first snippet.")
		return nil
	}

	// Create a prompt with search functionality
	prompt := promptui.Select{
		Label:             "Search snippets",
		Items:             snippets,
		Templates:         templates,
		Size:              10,
		Searcher:          searcher,
		StartInSearchMode: true,
	}

	// Custom searcher that filters snippets
	prompt.Searcher = func(input string, index int) bool {
		snip := snippets[index]
		return snip.MatchesSearch(input)
	}

	idx, _, err := prompt.Run()
	if err != nil {
		if err == promptui.ErrInterrupt {
			return nil
		}
		return fmt.Errorf("search failed: %w", err)
	}

	// Execute selected snippet
	selectedSnippet := snippets[idx]
	exec := executor.NewExecutor()
	return exec.Execute(selectedSnippet)
}
