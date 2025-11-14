package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"snip/pkg/snippet"

	"gopkg.in/yaml.v3"
)

// Storage handles reading and writing snippets to YAML file
type Storage struct {
	filePath string
}

// SnippetsFile represents the structure of the YAML file
type SnippetsFile struct {
	Snippets []snippet.Snippet `yaml:"snippets"`
}

// NewStorage creates a new Storage instance
func NewStorage() (*Storage, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	filePath := filepath.Join(homeDir, ".snip.yaml")
	return &Storage{filePath: filePath}, nil
}

// Load loads all snippets from the YAML file
func (s *Storage) Load() ([]snippet.Snippet, error) {
	// Check if file exists
	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		// File doesn't exist, return empty slice
		return []snippet.Snippet{}, nil
	}

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read snippets file: %w", err)
	}

	var snippetsFile SnippetsFile
	if err := yaml.Unmarshal(data, &snippetsFile); err != nil {
		return nil, fmt.Errorf("failed to parse snippets file: %w", err)
	}

	return snippetsFile.Snippets, nil
}

// Save saves all snippets to the YAML file
func (s *Storage) Save(snippets []snippet.Snippet) error {
	snippetsFile := SnippetsFile{
		Snippets: snippets,
	}

	data, err := yaml.Marshal(&snippetsFile)
	if err != nil {
		return fmt.Errorf("failed to marshal snippets: %w", err)
	}

	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write snippets file: %w", err)
	}

	return nil
}

// Add adds a new snippet
func (s *Storage) Add(newSnippet snippet.Snippet) error {
	snippets, err := s.Load()
	if err != nil {
		return err
	}

	snippets = append(snippets, newSnippet)
	return s.Save(snippets)
}

// Search performs fuzzy search on snippets
func (s *Storage) Search(query string) ([]snippet.Snippet, error) {
	snippets, err := s.Load()
	if err != nil {
		return nil, err
	}

	if query == "" {
		return snippets, nil
	}

	var results []snippet.Snippet
	for _, snip := range snippets {
		if snip.MatchesSearch(query) {
			results = append(results, snip)
		}
	}

	return results, nil
}

// GetFilePath returns the path to the snippets file
func (s *Storage) GetFilePath() string {
	return s.filePath
}
