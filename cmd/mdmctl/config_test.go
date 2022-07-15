package main

import "testing"

func TestValidateServerURL(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		errExpected bool
	}{
		{
			name:        "empty",
			input:       "",
			expected:    "",
			errExpected: true,
		},
		{
			name:        "whitespace",
			input:       " ",
			expected:    "",
			errExpected: true,
		},
		{
			name:     "http",
			input:    "http://localhost:8000",
			expected: "http://localhost:8000",
		},
		{
			name:     "https",
			input:    "https://localhost:8000",
			expected: "https://localhost:8000",
		},
		{
			name:     "trailing_slash",
			input:    "https://localhost:8000/",
			expected: "https://localhost:8000/",
		},
		{
			name:     "no_prefix",
			input:    "localhost:8000",
			expected: "https://localhost:8000",
		},
		{
			name:     "http_path",
			input:    "http://localhost:8000/path",
			expected: "http://localhost:8000/path",
		},
		{
			name:     "https_path",
			input:    "https://localhost:8000/path",
			expected: "https://localhost:8000/path",
		},
		{
			name:     "trailing_slash_path",
			input:    "https://localhost:8000/path/",
			expected: "https://localhost:8000/path/",
		},
		{
			name:     "no_prefix_path",
			input:    "localhost:8000/path",
			expected: "https://localhost:8000/path",
		},
		{
			name:     "multipart_path",
			input:    "https://localhost:8000/path1/path2/path3",
			expected: "https://localhost:8000/path1/path2/path3",
		},
		{
			name:     "multipart_path_trailing_slash",
			input:    "https://localhost:8000/longerpath1/longerpath2/longerpath3/",
			expected: "https://localhost:8000/longerpath1/longerpath2/longerpath3/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualURL, err := validateServerURL(tt.input)
			if haveErr, wantErr := err != nil, tt.errExpected; haveErr != wantErr {
				t.Errorf("have error - %t, want error - %t", haveErr, wantErr)
			}
			if have, want := actualURL, tt.expected; have != want {
				t.Errorf("have %s, want %s", have, want)
			}
		})
	}
}
