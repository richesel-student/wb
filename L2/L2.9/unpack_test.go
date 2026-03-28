package main

import "testing"

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		hasError bool
	}{
		{"a4bc2d5e", "aaaabccddddde", false},
		{"abcd", "abcd", false},
		{"45", "", true},
		{"", "", false},
	}

	for _, tt := range tests {
		res, err := unpack(tt.input)

		if tt.hasError {
			if err == nil {
				t.Errorf("expected error for input %q", tt.input)
			}
			continue
		}

		if err != nil {
			t.Errorf("unexpected error for input %q: %v", tt.input, err)
			continue
		}

		if res != tt.expected {
			t.Errorf("input %q: expected %q, got %q",
				tt.input, tt.expected, res)
		}
	}
}
