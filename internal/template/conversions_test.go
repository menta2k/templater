package template

import (
	"strings"
	"testing"
)

func TestToYAML(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{
			name:     "simple map",
			input:    map[string]any{"name": "test", "version": "1.0"},
			expected: "name: test\nversion: \"1.0\"",
		},
		{
			name:     "nested map",
			input:    map[string]any{"app": map[string]any{"name": "myapp", "port": 8080}},
			expected: "app:\n  name: myapp\n  port: 8080",
		},
		{
			name:     "array",
			input:    []string{"item1", "item2", "item3"},
			expected: "- item1\n- item2\n- item3",
		},
		{
			name:     "simple string",
			input:    "hello world",
			expected: "hello world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toYAML(tt.input)
			if !strings.Contains(result, strings.Split(tt.expected, "\n")[0]) {
				t.Errorf("Expected result to contain %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestMustToYAML(t *testing.T) {
	tests := []struct {
		name        string
		input       any
		shouldPanic bool
	}{
		{
			name:        "valid input",
			input:       map[string]any{"name": "test"},
			shouldPanic: false,
		},
		{
			name:        "invalid input - function",
			input:       func() {},
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if tt.shouldPanic && r == nil {
					t.Error("Expected panic but didn't get one")
				}
				if !tt.shouldPanic && r != nil {
					t.Errorf("Unexpected panic: %v", r)
				}
			}()

			result := mustToYAML(tt.input)
			if !tt.shouldPanic && result == "" {
				t.Error("Expected non-empty result for valid input")
			}
		})
	}
}

func TestFromYAML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		hasError bool
	}{
		{
			name:     "valid yaml",
			input:    "name: test\nversion: 1.0",
			hasError: false,
		},
		{
			name:     "invalid yaml",
			input:    "invalid: yaml: content: [",
			hasError: true,
		},
		{
			name:     "empty string",
			input:    "",
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fromYAML(tt.input)
			if tt.hasError {
				if _, exists := result["Error"]; !exists {
					t.Error("Expected Error key in result for invalid YAML")
				}
			} else {
				if _, exists := result["Error"]; exists {
					t.Errorf("Unexpected error in result: %v", result["Error"])
				}
			}
		})
	}
}

func TestFromYAMLArray(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		hasError bool
	}{
		{
			name:     "valid array",
			input:    "- item1\n- item2\n- item3",
			hasError: false,
		},
		{
			name:     "invalid yaml",
			input:    "- invalid: [",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fromYAMLArray(tt.input)
			if tt.hasError {
				if len(result) == 0 || result[0] == nil {
					t.Error("Expected error message in result for invalid YAML")
				}
			} else {
				if len(result) == 0 {
					t.Error("Expected non-empty array for valid YAML")
				}
			}
		})
	}
}

func TestToJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{
			name:     "simple map",
			input:    map[string]any{"name": "test", "version": "1.0"},
			expected: `{"name":"test","version":"1.0"}`,
		},
		{
			name:     "array",
			input:    []string{"item1", "item2", "item3"},
			expected: `["item1","item2","item3"]`,
		},
		{
			name:     "string",
			input:    "hello",
			expected: `"hello"`,
		},
		{
			name:     "number",
			input:    42,
			expected: "42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toJSON(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestMustToJSON(t *testing.T) {
	tests := []struct {
		name        string
		input       any
		shouldPanic bool
	}{
		{
			name:        "valid input",
			input:       map[string]any{"name": "test"},
			shouldPanic: false,
		},
		{
			name:        "invalid input - function",
			input:       func() {},
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if tt.shouldPanic && r == nil {
					t.Error("Expected panic but didn't get one")
				}
				if !tt.shouldPanic && r != nil {
					t.Errorf("Unexpected panic: %v", r)
				}
			}()

			result := mustToJSON(tt.input)
			if !tt.shouldPanic && result == "" {
				t.Error("Expected non-empty result for valid input")
			}
		})
	}
}

func TestFromJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		hasError bool
	}{
		{
			name:     "valid json",
			input:    `{"name":"test","version":"1.0"}`,
			hasError: false,
		},
		{
			name:     "invalid json",
			input:    `{"invalid": json}`,
			hasError: true,
		},
		{
			name:     "empty string",
			input:    "",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fromJSON(tt.input)
			if tt.hasError {
				if _, exists := result["Error"]; !exists {
					t.Error("Expected Error key in result for invalid JSON")
				}
			} else {
				if _, exists := result["Error"]; exists {
					t.Errorf("Unexpected error in result: %v", result["Error"])
				}
			}
		})
	}
}

func TestFromJSONArray(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		hasError bool
	}{
		{
			name:     "valid array",
			input:    `["item1","item2","item3"]`,
			hasError: false,
		},
		{
			name:     "invalid json",
			input:    `["invalid": json]`,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fromJSONArray(tt.input)
			if tt.hasError {
				if len(result) == 0 || result[0] == nil {
					t.Error("Expected error message in result for invalid JSON")
				}
			} else {
				if len(result) != 3 {
					t.Errorf("Expected 3 items, got %d", len(result))
				}
			}
		})
	}
}

func TestToTOML(t *testing.T) {
	tests := []struct {
		name  string
		input any
	}{
		{
			name:  "simple map",
			input: map[string]any{"name": "test", "version": "1.0"},
		},
		{
			name:  "nested map",
			input: map[string]any{"database": map[string]any{"host": "localhost", "port": 5432}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toTOML(tt.input)
			if result == "" {
				t.Error("Expected non-empty TOML result")
			}
			// Basic check that it contains expected content
			if tt.name == "simple map" && (!strings.Contains(result, "name") || !strings.Contains(result, "test")) {
				t.Errorf("Expected TOML to contain name and test, got: %s", result)
			}
		})
	}
}

func TestFromTOML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		hasError bool
	}{
		{
			name:     "valid toml",
			input:    "name = \"test\"\nversion = \"1.0\"",
			hasError: false,
		},
		{
			name:     "invalid toml",
			input:    `invalid toml [[[`,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fromTOML(tt.input)
			if tt.hasError {
				if _, exists := result["Error"]; !exists {
					t.Error("Expected Error key in result for invalid TOML")
				}
			} else {
				if _, exists := result["Error"]; exists {
					t.Errorf("Unexpected error in result: %v", result["Error"])
				}
				if tt.name == "valid toml" {
					if result["name"] != "test" {
						t.Errorf("Expected name=test, got name=%v", result["name"])
					}
				}
			}
		})
	}
}