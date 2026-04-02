package llm

import (
	"strings"
	"unicode/utf8"
)

// EstimateTokens approximates token count for GPT models.
func EstimateTokens(text string) int {
	// Rough approximation: 1 token ≈ 4 characters for English
	chars := utf8.RuneCountInString(strings.TrimSpace(text))
	return (chars + 3) / 4
}
