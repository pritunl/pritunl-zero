package utils

import (
	"strings"
	"testing"
)

func TestFilterName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "allowed characters",
			input:    "Alpha-123",
			expected: "Alpha-123",
		},
		{
			name:     "unsupported characters",
			input:    "alpha beta/_@#",
			expected: "alphabeta",
		},
		{
			name:     "length limit",
			input:    strings.Repeat("a", nameSafeLimit+1),
			expected: strings.Repeat("a", nameSafeLimit),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := FilterName(test.input)
			if result != test.expected {
				t.Fatalf("expected %q, got %q", test.expected, result)
			}
		})
	}
}
