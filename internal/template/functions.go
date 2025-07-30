package template

import (
	"maps"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

// GetTemplateFuncs returns the enhanced set of template functions including sprig functions.
// and format conversion functions (JSON, YAML, TOML).
func GetTemplateFuncs() template.FuncMap {
	f := sprig.TxtFuncMap()
	// Remove potentially dangerous functions
	delete(f, "env")
	delete(f, "expandenv")

	const notImplementedStr = "not implemented"

	// Add format conversion and additional utility functions
	extra := template.FuncMap{
		// YAML functions
		"toYaml":        toYAML,
		"mustToYaml":    mustToYAML,
		"toYamlPretty":  toYAMLPretty,
		"fromYaml":      fromYAML,
		"fromYamlArray": fromYAMLArray,

		// JSON functions
		"toJson":        toJSON,
		"mustToJson":    mustToJSON,
		"fromJson":      fromJSON,
		"fromJsonArray": fromJSONArray,

		// TOML functions
		"toToml":   toTOML,
		"fromToml": fromTOML,

		// Placeholder functions for advanced features
		"include":  func(string, any) string { return notImplementedStr },
		"tpl":      func(string, any) any { return notImplementedStr },
		"required": func(string, any) (any, error) { return notImplementedStr, nil },
		"lookup": func(string, string, string, string) (map[string]any, error) {
			return map[string]any{}, nil
		},
	}

	// Merge sprig functions with our custom functions
	maps.Copy(f, extra)

	return f
}
