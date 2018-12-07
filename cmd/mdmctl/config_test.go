package main

import "testing"

func TestValidateServerURL(t *testing.T) {

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "http",
			input:    "http://localhost:8000",
			expected: "http://localhost:8000/",
		},
		{
			name:     "https",
			input:    "https://localhost:8000",
			expected: "https://localhost:8000/",
		},
		{
			name:     "trailing_slash",
			input:    "https://localhost:8000/",
			expected: "https://localhost:8000/",
		},
		{
			name:     "no_prefix",
			input:    "localhost:8000",
			expected: "https://localhost:8000/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualURL, _ := validateServerURL(tt.input)
			if have, want := actualURL, tt.expected; have != want {
				t.Errorf("have %s, want %s", have, want)
			}
		})
	}
}
