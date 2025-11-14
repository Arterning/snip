package snippet

import (
	"regexp"
	"strings"
	"time"
)

// Snippet represents a command snippet with placeholders
type Snippet struct {
	ID          string    `yaml:"id"`
	Command     string    `yaml:"command"`
	Description string    `yaml:"description,omitempty"`
	CreatedAt   time.Time `yaml:"created_at"`
	UpdatedAt   time.Time `yaml:"updated_at"`
}

// ExtractPlaceholders extracts all placeholders from the command
// Placeholders are in the format <placeholder_name>
func (s *Snippet) ExtractPlaceholders() []string {
	re := regexp.MustCompile(`<([^>]+)>`)
	matches := re.FindAllStringSubmatch(s.Command, -1)

	// Use map to deduplicate placeholders
	placeholderMap := make(map[string]bool)
	var placeholders []string

	for _, match := range matches {
		if len(match) > 1 {
			placeholder := match[1]
			if !placeholderMap[placeholder] {
				placeholderMap[placeholder] = true
				placeholders = append(placeholders, placeholder)
			}
		}
	}

	return placeholders
}

// ReplacePlaceholders replaces placeholders with actual values
func (s *Snippet) ReplacePlaceholders(values map[string]string) string {
	result := s.Command

	for placeholder, value := range values {
		result = strings.ReplaceAll(result, "<"+placeholder+">", value)
	}

	return result
}

// MatchesSearch checks if the snippet matches the search query (fuzzy search)
func (s *Snippet) MatchesSearch(query string) bool {
	query = strings.ToLower(query)
	command := strings.ToLower(s.Command)
	description := strings.ToLower(s.Description)

	// Check if query is contained in command or description
	return strings.Contains(command, query) || strings.Contains(description, query)
}
