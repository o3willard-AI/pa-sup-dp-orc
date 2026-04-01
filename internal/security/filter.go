package security

import (
	"fmt"
	"regexp"
)

// Filter redacts sensitive patterns from text.
type Filter struct {
	patterns []*regexp.Regexp
	replace  string
}

// NewFilter creates a filter with compiled regex patterns.
func NewFilter(rawPatterns []string) (*Filter, error) {
	f := &Filter{
		replace: "[REDACTED]",
	}
	for _, p := range rawPatterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("invalid regex pattern %q: %w", p, err)
		}
		f.patterns = append(f.patterns, re)
	}
	return f, nil
}

// Redact replaces matches of any pattern with the replacement string.
func (f *Filter) Redact(text string) string {
	result := text
	for _, re := range f.patterns {
		result = re.ReplaceAllString(result, f.replace)
	}
	return result
}

// ContainsSensitive returns true if text matches any pattern.
func (f *Filter) ContainsSensitive(text string) bool {
	for _, re := range f.patterns {
		if re.MatchString(text) {
			return true
		}
	}
	return false
}

// DefaultFilter returns a filter with common sensitive patterns.
func DefaultFilter() *Filter {
	patterns := []string{
		`(?i)password\s*[:=]\s*['\"]?[^'\"]+`,
		`(?i)(api[_-]?key|token)[\s:=]+['\"]?[a-zA-Z0-9_\-]{20,}['\"]?`,
		`AKIA[0-9A-Z]{16}`,
		`-----BEGIN (RSA|DSA|EC|OPENSSH) PRIVATE KEY-----`,
		`(?i)secret[\s:=]+['\"]?[a-zA-Z0-9_\-]{10,}['\"]?`,
	}
	// Ignore compilation errors for known‑good patterns
	f, err := NewFilter(patterns)
	if err != nil {
		// Should never happen for the hardcoded patterns; fallback to empty filter
		f = &Filter{replace: "[REDACTED]"}
	}
	return f
}
