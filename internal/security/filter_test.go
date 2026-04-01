package security

import (
	"testing"
)

func TestFilter_Redact(t *testing.T) {
	patterns := []string{
		`password\s*=\s*\w+`,
		`token: [a-f0-9]{32}`,
	}
	f, err := NewFilter(patterns)
	if err != nil {
		t.Fatalf("NewFilter failed: %v", err)
	}

	tests := []struct {
		input    string
		expected string
	}{
		{"password = secret123", "[REDACTED]"},
		{"token: abcdef1234567890abcdef1234567890", "[REDACTED]"},
		{"no sensitive data", "no sensitive data"},
		{"multiple: password = foo token: abcdef1234567890abcdef1234567890",
			"multiple: [REDACTED] [REDACTED]"},
	}

	for _, tt := range tests {
		got := f.Redact(tt.input)
		if got != tt.expected {
			t.Errorf("Redact(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestFilter_ContainsSensitive(t *testing.T) {
	f, _ := NewFilter([]string{`password\s*=\s*\w+`})
	if !f.ContainsSensitive("password = hello") {
		t.Error("ContainsSensitive should detect password")
	}
	if f.ContainsSensitive("no password here") {
		t.Error("ContainsSensitive false positive")
	}
}

func TestDefaultFilter(t *testing.T) {
	f := DefaultFilter()
	if f == nil {
		t.Fatal("DefaultFilter returned nil")
	}
	// Ensure it redacts a known pattern
	text := "password: supersecret"
	redacted := f.Redact(text)
	if redacted == text {
		t.Errorf("DefaultFilter did not redact password")
	}
}
