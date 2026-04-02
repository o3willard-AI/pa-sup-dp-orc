package llm

import (
	"testing"
)

func TestEstimateTokens(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "empty string",
			input:    "",
			expected: 0,
		},
		{
			name:     "whitespace only",
			input:    "   \t\n  ",
			expected: 0,
		},
		{
			name:     "single word",
			input:    "hello",
			expected: 2, // 5 chars -> (5+3)/4 = 2
		},
		{
			name:     "exactly 4 characters",
			input:    "abcd",
			expected: 1, // (4+3)/4 = 7/4 = 1
		},
		{
			name:     "exactly 8 characters",
			input:    "abcdefgh",
			expected: 2, // (8+3)/4 = 11/4 = 2
		},
		{
			name:     "leading/trailing spaces",
			input:    "  hello world  ",
			expected: 3, // "hello world" length 11, (11+3)/4 = 14/4 = 3
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EstimateTokens(tt.input)
			if got != tt.expected {
				t.Errorf("EstimateTokens(%q) = %d, want %d", tt.input, got, tt.expected)
			}
		})
	}
}
