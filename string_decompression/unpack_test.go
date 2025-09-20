package main

import (
	"testing"
)

func TestUnpackString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		err      error
	}{
		// Simple testcase
		{"a", "a", nil},
		{"abc", "abc", nil},
		{"a3", "aaa", nil},
		{"a0", "", nil},
		{"a10", "aaaaaaaaaa", nil},

		// With escape
		{"\\\\", "\\", nil},
		{"a\\3", "a3", nil},
		{"a\\\\3", "a\\\\\\", nil},

		// Errors
		{"3abc", "", ErrInvalidStartDigit},
		{"a1b2c3", "abbccc", nil},
		{"a\\", "", ErrInvalidEscape},
		{"a1000001", "", ErrRepeatTooLarge},
	}

	for _, test := range tests {
		result, err := unpackString(test.input)

		if result != test.expected {
			t.Errorf("unpackString(%q) = %q; want %q", test.input, result, test.expected)
		}

		if err != test.err {
			t.Errorf("unpackString(%q) error = %v; want %v", test.input, err, test.err)
		}
	}
}
