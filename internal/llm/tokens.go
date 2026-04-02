package llm

import (
	"strings"
)

// EstimateTokens approximates token count for GPT models.
func EstimateTokens(text string) int {
	// Rough approximation: 1 token ≈ 4 characters for English
	chars := len(strings.TrimSpace(text))
	return (chars + 3) / 4
}