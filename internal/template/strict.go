package template

import (
	"fmt"
	"strings"
	"text/template"
)

// StrictModeError represents an error that occurs when a template variable is undefined in strict mode.
type StrictModeError struct {
	Variable string
	Template string
}

func (e *StrictModeError) Error() string {
	return fmt.Sprintf("undefined variable '%s' in template '%s' (strict mode enabled)", e.Variable, e.Template)
}

// StrictTemplate wraps a template to provide strict mode validation.
type StrictTemplate struct {
	*template.Template
	StrictMode bool
}

// NewStrictTemplate creates a new template wrapper with strict mode support.
func NewStrictTemplate(name string, strictMode bool) *StrictTemplate {
	tmpl := template.New(name).Funcs(GetTemplateFuncs())

	// In strict mode, we need to add a custom function that catches undefined variables
	if strictMode {
		// Override the default "missingkey" option to error instead of printing "<no value>"
		tmpl = tmpl.Option("missingkey=error")
	}

	return &StrictTemplate{
		Template:   tmpl,
		StrictMode: strictMode,
	}
}

// ParseTemplate parses template content with strict mode considerations.
func (st *StrictTemplate) ParseTemplate(content string) (*StrictTemplate, error) {
	tmpl, err := st.Template.Parse(content)
	if err != nil {
		return nil, err
	}

	return &StrictTemplate{
		Template:   tmpl,
		StrictMode: st.StrictMode,
	}, nil
}

// ExecuteTemplate executes the template with strict mode validation.
func (st *StrictTemplate) ExecuteTemplate(data any) (string, error) {
	var result strings.Builder

	if st.StrictMode {
		// In strict mode, template execution will fail with "missingkey=error" option
		// if any undefined variables are encountered
		err := st.Template.Execute(&result, data)
		if err != nil {
			// Check if it's a missing key error and wrap it appropriately
			if strings.Contains(err.Error(), "map has no entry for key") ||
				strings.Contains(err.Error(), "can't evaluate field") {
				return "", &StrictModeError{
					Variable: extractVariableFromError(err.Error()),
					Template: st.Template.Name(),
				}
			}
			return "", err
		}
	} else {
		// Normal mode - missing keys will be replaced with "<no value>"
		err := st.Template.Execute(&result, data)
		if err != nil {
			return "", err
		}
	}

	return result.String(), nil
}

// extractVariableFromError attempts to extract the variable name from a template error message.
func extractVariableFromError(errMsg string) string {
	// Try to extract variable name from common error patterns
	if strings.Contains(errMsg, "map has no entry for key") {
		// Example: "template: test:1:2: executing \"test\" at <.undefined>: map has no entry for key \"undefined\""
		parts := strings.Split(errMsg, "key \"")
		if len(parts) > 1 {
			endQuote := strings.Index(parts[1], "\"")
			if endQuote > 0 {
				return parts[1][:endQuote]
			}
		}
	}

	if strings.Contains(errMsg, "can't evaluate field") {
		// Example: "template: test:1:2: executing \"test\" at <.app.undefined>: can't evaluate field undefined in type map[string]interface {}"
		parts := strings.Split(errMsg, "field ")
		if len(parts) > 1 {
			spaceParts := strings.Split(parts[1], " ")
			if len(spaceParts) > 0 {
				return spaceParts[0]
			}
		}
	}

	// If we can't extract the specific variable, return a generic message
	return "unknown"
}
