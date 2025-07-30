package template

import (
	"strings"
	"testing"
	"text/template"
)

func TestGetTemplateFuncs(t *testing.T) {
	funcs := GetTemplateFuncs()

	// Test that sprig functions are present (sample check)
	sprigFuncs := []string{"quote", "squote", "cat", "indent", "nindent", "replace"}
	for _, funcName := range sprigFuncs {
		if _, exists := funcs[funcName]; !exists {
			t.Errorf("Expected sprig function %s to be present in template funcs", funcName)
		}
	}

	// Test that our custom format conversion functions are present
	customFuncs := []string{
		"toYaml", "mustToYaml", "toYamlPretty", "fromYaml", "fromYamlArray",
		"toJson", "mustToJson", "fromJson", "fromJsonArray",
		"toToml", "fromToml",
	}
	for _, funcName := range customFuncs {
		if _, exists := funcs[funcName]; !exists {
			t.Errorf("Expected custom function %s to be present in template funcs", funcName)
		}
	}

	// Test that placeholder functions are present
	placeholderFuncs := []string{"include", "tpl", "required", "lookup"}
	for _, funcName := range placeholderFuncs {
		if _, exists := funcs[funcName]; !exists {
			t.Errorf("Expected placeholder function %s to be present in template funcs", funcName)
		}
	}

	// Test that dangerous functions are removed
	dangerousFuncs := []string{"env", "expandenv"}
	for _, funcName := range dangerousFuncs {
		if _, exists := funcs[funcName]; exists {
			t.Errorf("Dangerous function %s should be removed from template funcs", funcName)
		}
	}
}


func TestEnhancedTemplateFunctionsInTemplate(t *testing.T) {
	templateText := `
Default: {{default "fallback" .value}}
Quote: {{quote .text}}
Indent: {{indent 2 .yaml}}
JSON: {{toJson .data}}
YAML: {{toYaml .data}}
TOML: {{toToml .data}}
`

	tmpl := template.New("test").Funcs(GetTemplateFuncs())
	tmpl, err := tmpl.Parse(templateText)
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	data := map[string]any{
		"value": "",
		"text":  "Hello World",
		"yaml":  "line1\nline2",
		"data": map[string]any{
			"name":    "test",
			"version": 1,
		},
	}

	var result strings.Builder
	err = tmpl.Execute(&result, data)
	if err != nil {
		t.Fatalf("Failed to execute template: %v", err)
	}

	output := result.String()

	// Check that functions worked correctly
	if !strings.Contains(output, "Default: fallback") {
		t.Error("Default function did not work correctly")
	}
	if !strings.Contains(output, `Quote: "Hello World"`) {
		t.Error("Quote function did not work correctly")
	}
	if !strings.Contains(output, "JSON: {") {
		t.Error("toJson function did not work correctly")
	}
	if !strings.Contains(output, "YAML: name: test") {
		t.Error("toYaml function did not work correctly")
	}
	if !strings.Contains(output, "TOML: name = ") {
		t.Error("toToml function did not work correctly")
	}
}

func TestSprigFunctionsInTemplate(t *testing.T) {
	templateText := `
Upper: {{upper .text}}
Title: {{title .text}}
Trim: {{trim .padded}}
Replace: {{replace " " "-" .text}}
Repeat: {{repeat 3 "x"}}
`

	tmpl := template.New("test").Funcs(GetTemplateFuncs())
	tmpl, err := tmpl.Parse(templateText)
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	data := map[string]any{
		"text":   "hello world",
		"padded": "  trimmed  ",
	}

	var result strings.Builder
	err = tmpl.Execute(&result, data)
	if err != nil {
		t.Fatalf("Failed to execute template: %v", err)
	}

	output := result.String()

	// Check that sprig functions worked correctly
	if !strings.Contains(output, "Upper: HELLO WORLD") {
		t.Error("Sprig upper function did not work correctly")
	}
	if !strings.Contains(output, "Title: Hello World") {
		t.Error("Sprig title function did not work correctly")
	}
	if !strings.Contains(output, "Trim: trimmed") {
		t.Error("Sprig trim function did not work correctly")
	}
	if !strings.Contains(output, "Replace: hello-world") {
		t.Error("Sprig replace function did not work correctly")
	}
	if !strings.Contains(output, "Repeat: xxx") {
		t.Error("Sprig repeat function did not work correctly")
	}
}