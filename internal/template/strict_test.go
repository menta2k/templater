package template

import (
	"strings"
	"testing"
)

func TestStrictModeError(t *testing.T) {
	err := &StrictModeError{
		Variable: "undefined",
		Template: "test.tpl",
	}

	expected := "undefined variable 'undefined' in template 'test.tpl' (strict mode enabled)"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
}

func TestNewStrictTemplate(t *testing.T) {
	tests := []struct {
		name       string
		strictMode bool
	}{
		{"strict mode enabled", true},
		{"strict mode disabled", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl := NewStrictTemplate("test", tt.strictMode)
			if tmpl == nil {
				t.Error("NewStrictTemplate returned nil")
			}
			if tmpl.StrictMode != tt.strictMode {
				t.Errorf("Expected StrictMode %t, got %t", tt.strictMode, tmpl.StrictMode)
			}
			if tmpl.Template == nil {
				t.Error("Template field is nil")
			}
		})
	}
}

func TestStrictModeExecution(t *testing.T) {
	tests := []struct {
		name        string
		template    string
		data        map[string]any
		strictMode  bool
		expectError bool
		errorType   string
	}{
		{
			name:        "defined variable - strict mode",
			template:    "Hello {{.name}}!",
			data:        map[string]any{"name": "world"},
			strictMode:  true,
			expectError: false,
		},
		{
			name:        "defined variable - normal mode",
			template:    "Hello {{.name}}!",
			data:        map[string]any{"name": "world"},
			strictMode:  false,
			expectError: false,
		},
		{
			name:        "undefined variable - strict mode",
			template:    "Hello {{.undefined}}!",
			data:        map[string]any{"name": "world"},
			strictMode:  true,
			expectError: true,
			errorType:   "StrictModeError",
		},
		{
			name:        "undefined variable - normal mode",
			template:    "Hello {{.undefined}}!",
			data:        map[string]any{"name": "world"},
			strictMode:  false,
			expectError: false,
		},
		{
			name:        "nested undefined variable - strict mode",
			template:    "App: {{.app.undefined}}",
			data:        map[string]any{"app": map[string]any{"name": "myapp"}},
			strictMode:  true,
			expectError: true,
			errorType:   "StrictModeError",
		},
		{
			name:        "nested undefined variable - normal mode",
			template:    "App: {{.app.undefined}}",
			data:        map[string]any{"app": map[string]any{"name": "myapp"}},
			strictMode:  false,
			expectError: false,
		},
		{
			name:        "undefined root object - strict mode",
			template:    "Database: {{.database.host}}",
			data:        map[string]any{"app": map[string]any{"name": "myapp"}},
			strictMode:  true,
			expectError: true,
			errorType:   "StrictModeError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl := NewStrictTemplate("test", tt.strictMode)
			parsedTemplate, err := tmpl.ParseTemplate(tt.template)
			if err != nil {
				t.Fatalf("Failed to parse template: %v", err)
			}

			result, err := parsedTemplate.ExecuteTemplate(tt.data)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
					return
				}

				if tt.errorType == "StrictModeError" {
					if _, ok := err.(*StrictModeError); !ok {
						t.Errorf("Expected StrictModeError, got %T: %v", err, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
					return
				}

				if result == "" && tt.template != "" {
					t.Error("Expected non-empty result")
				}

				// In normal mode, undefined variables should be replaced with "<no value>"
				if !tt.strictMode && strings.Contains(tt.template, ".undefined") && !strings.Contains(result, "<no value>") {
					t.Errorf("Expected '<no value>' in result for undefined variable in normal mode, got: %s", result)
				}
			}
		})
	}
}

func TestExtractVariableFromError(t *testing.T) {
	tests := []struct {
		name     string
		errorMsg string
		expected string
	}{
		{
			name:     "map has no entry error",
			errorMsg: `template: test:1:2: executing "test" at <.undefined>: map has no entry for key "undefined"`,
			expected: "undefined",
		},
		{
			name:     "can't evaluate field error",
			errorMsg: `template: test:1:2: executing "test" at <.app.missing>: can't evaluate field missing in type map[string]interface {}`,
			expected: "missing",
		},
		{
			name:     "unknown error format",
			errorMsg: "some other error",
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractVariableFromError(tt.errorMsg)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestStrictModeWithSprigFunctions(t *testing.T) {
	tests := []struct {
		name        string
		template    string
		data        map[string]any
		strictMode  bool
		expectError bool
	}{
		{
			name:        "sprig function with defined variable - strict mode",
			template:    "{{.name | upper}}",
			data:        map[string]any{"name": "hello"},
			strictMode:  true,
			expectError: false,
		},
		{
			name:        "sprig function with undefined variable - strict mode",
			template:    "{{.undefined | upper}}",
			data:        map[string]any{"name": "hello"},
			strictMode:  true,
			expectError: true,
		},
		{
			name:        "default function with undefined variable - strict mode",
			template:    "{{default \"fallback\" .undefined}}",
			data:        map[string]any{"name": "hello"},
			strictMode:  true,
			expectError: true,
		},
		{
			name:        "default function with undefined variable - normal mode",
			template:    "{{default \"fallback\" .undefined}}",
			data:        map[string]any{"name": "hello"},
			strictMode:  false,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl := NewStrictTemplate("test", tt.strictMode)
			parsedTemplate, err := tmpl.ParseTemplate(tt.template)
			if err != nil {
				t.Fatalf("Failed to parse template: %v", err)
			}

			result, err := parsedTemplate.ExecuteTemplate(tt.data)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result == "" {
					t.Error("Expected non-empty result")
				}
			}
		})
	}
}
