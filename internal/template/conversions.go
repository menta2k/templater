package template

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/BurntSushi/toml"
	yaml3 "gopkg.in/yaml.v3"
	"sigs.k8s.io/yaml"
)

// YAML conversion functions.

// toYAML takes an interface, marshals it to yaml, and returns a string. It will.
// always return a string, even on marshal error (empty string).
func toYAML(v any) string {
	// Convert map[interface{}]interface{} to map[string]any if needed
	converted := convertMapKeys(v)

	data, err := yaml.Marshal(converted)
	if err != nil {
		// Swallow errors inside of a template.
		return ""
	}
	return strings.TrimSuffix(string(data), "\n")
}

// mustToYAML takes an interface, marshals it to yaml, and returns a string.
// It will panic if there is an error.
func mustToYAML(v any) string {
	converted := convertMapKeys(v)
	data, err := yaml.Marshal(converted)
	if err != nil {
		panic(err)
	}
	return strings.TrimSuffix(string(data), "\n")
}

// toYAMLPretty takes an interface, marshals it to pretty yaml, and returns a string.
func toYAMLPretty(v any) string {
	var data bytes.Buffer
	encoder := yaml3.NewEncoder(&data)
	encoder.SetIndent(2)
	err := encoder.Encode(v)
	if err != nil {
		// Swallow errors inside of a template.
		return ""
	}
	return strings.TrimSuffix(data.String(), "\n")
}

// fromYAML converts a YAML document into a map[string]any.
func fromYAML(str string) map[string]any {
	m := map[string]any{}

	if err := yaml.Unmarshal([]byte(str), &m); err != nil {
		m["Error"] = err.Error()
	}
	return m
}

// fromYAMLArray converts a YAML array into a []any.
func fromYAMLArray(str string) []any {
	a := []any{}

	if err := yaml.Unmarshal([]byte(str), &a); err != nil {
		a = []any{err.Error()}
	}
	return a
}

// JSON conversion functions.

// toJSON takes an interface, marshals it to json, and returns a string.
func toJSON(v any) string {
	converted := convertMapKeys(v)
	data, err := json.Marshal(converted)
	if err != nil {
		// Swallow errors inside of a template.
		return ""
	}
	return string(data)
}

// mustToJSON takes an interface, marshals it to json, and returns a string.
// It will panic if there is an error.
func mustToJSON(v any) string {
	converted := convertMapKeys(v)
	data, err := json.Marshal(converted)
	if err != nil {
		panic(err)
	}
	return string(data)
}

// fromJSON converts a JSON document into a map[string]any.
func fromJSON(str string) map[string]any {
	m := make(map[string]any)

	if err := json.Unmarshal([]byte(str), &m); err != nil {
		m["Error"] = err.Error()
	}
	return m
}

// fromJSONArray converts a JSON array into a []any.
func fromJSONArray(str string) []any {
	a := []any{}

	if err := json.Unmarshal([]byte(str), &a); err != nil {
		a = []any{err.Error()}
	}
	return a
}

// TOML conversion functions.

// toTOML takes an interface, marshals it to toml, and returns a string.
func toTOML(v any) string {
	converted := convertMapKeys(v)
	b := bytes.NewBuffer(nil)
	e := toml.NewEncoder(b)
	err := e.Encode(converted)
	if err != nil {
		return err.Error()
	}
	return b.String()
}

// fromTOML converts a TOML document into a map[string]any.
func fromTOML(str string) map[string]any {
	m := make(map[string]any)

	if err := toml.Unmarshal([]byte(str), &m); err != nil {
		m["Error"] = err.Error()
	}
	return m
}

// convertMapKeys recursively converts map[interface{}]interface{} to map[string]any.
// This is needed because YAML unmarshaling creates interface{} keys which cause.
// JSON marshaling to fail.
func convertMapKeys(v any) any {
	switch x := v.(type) {
	case map[interface{}]interface{}:
		m := make(map[string]any)
		for k, v := range x {
			if keyStr, ok := k.(string); ok {
				m[keyStr] = convertMapKeys(v)
			}
		}
		return m
	case map[string]interface{}:
		m := make(map[string]any)
		for k, v := range x {
			m[k] = convertMapKeys(v)
		}
		return m
	case []interface{}:
		var converted []any
		for _, item := range x {
			converted = append(converted, convertMapKeys(item))
		}
		return converted
	default:
		return v
	}
}
