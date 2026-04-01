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

func TestNewFilter_Error(t *testing.T) {
	_, err := NewFilter([]string{"["}) // invalid regex
	if err == nil {
		t.Error("expected error for invalid regex")
	}
}

func TestNewFilter_EmptyPatterns(t *testing.T) {
	f, err := NewFilter([]string{})
	if err != nil {
		t.Fatalf("NewFilter with empty patterns should not error: %v", err)
	}
	if f.ContainsSensitive("anything") {
		t.Error("empty filter should not detect anything")
	}
	if got := f.Redact("password = secret"); got != "password = secret" {
		t.Errorf("empty filter should not modify text, got %q", got)
	}
}

func TestFilter_ZeroValue(t *testing.T) {
	var f Filter
	// Zero value should not panic
	_ = f.Redact("test")
	if f.ContainsSensitive("test") {
		t.Error("zero value filter should not detect anything (no patterns)")
	}
}

func TestDefaultFilter_Patterns(t *testing.T) {
	f := DefaultFilter()
	tests := []struct {
		input        string
		shouldRedact bool
	}{
		{"password = secret", true},
		{"api_key = abcdefghijklmnopqrstuvwxyz123456", true},
		{"token = abcdefghijklmnopqrstuvwxyz123456", true},
		{"AWS key AKIA0123456789ABCDEF", true},
		{"-----BEGIN RSA PRIVATE KEY-----", true},
		{"secret = mysecret123", true},
		{"nothing sensitive", false},
	}
	for _, tt := range tests {
		redacted := f.Redact(tt.input)
		didRedact := redacted != tt.input
		if didRedact != tt.shouldRedact {
			t.Errorf("pattern mismatch for %q: redacted=%v, expected redact=%v", tt.input, didRedact, tt.shouldRedact)
		}
	}
}
