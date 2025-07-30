package values

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

// Loader handles loading values from various sources.
type Loader struct{}

// NewLoader creates a new values loader.
func NewLoader() *Loader {
	return &Loader{}
}

// LoadYAMLValues loads values from a YAML file.
func (l *Loader) LoadYAMLValues(valuesFile string) (map[string]any, error) {
	values := make(map[string]any)

	if valuesFile == "" {
		return values, nil
	}

	data, err := os.ReadFile(valuesFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read values file: %w", err)
	}

	err = yaml.Unmarshal(data, &values)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return values, nil
}

// LoadEnvValues loads values from environment variables and converts keys to camelCase.
func (l *Loader) LoadEnvValues() map[string]any {
	envValues := make(map[string]any)

	for _, env := range os.Environ() {
		// Split environment variable into key-value pair
		for i := 0; i < len(env); i++ {
			if env[i] == '=' {
				key := env[:i]
				value := env[i+1:]

				// Convert environment variable key to camelCase
				camelKey := l.toCamelCase(key)
				envValues[camelKey] = value
				break
			}
		}
	}

	return envValues
}

// ParseSetValues parses command-line set values.
func (l *Loader) ParseSetValues(setValues []string) (map[string]any, error) {
	parsedValues := make(map[string]any)

	for _, setValue := range setValues {
		// Split by comma to handle comma-separated key=value pairs
		pairs := strings.Split(setValue, ",")

		for _, pair := range pairs {
			pair = strings.TrimSpace(pair)
			if pair == "" {
				continue
			}

			// Split by = to get key and value
			parts := strings.SplitN(pair, "=", 2)
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid set value format: %s (expected key=value)", pair)
			}

			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// Convert value to appropriate type
			convertedValue := l.convertValue(value)

			// Handle nested keys (e.g., app.name=value)
			err := l.setNestedValue(parsedValues, key, convertedValue)
			if err != nil {
				return nil, fmt.Errorf("error setting nested value for key %s: %w", key, err)
			}
		}
	}

	return parsedValues, nil
}

// MergeValues merges YAML values with environment variables and --set values.
// Precedence: --set values > environment variables > YAML values.
func (l *Loader) MergeValues(yamlValues, envValues, setValues, configValues map[string]any) map[string]any {
	merged := make(map[string]any)

	// Start with YAML values (lowest precedence)
	l.deepMerge(merged, yamlValues)

	// Override with environment variables
	l.deepMerge(merged, envValues)

	// Override with --set values (highest precedence)
	l.deepMerge(merged, setValues)

	// Add any values from config
	l.deepMerge(merged, configValues)

	return merged
}

// convertValue attempts to convert string values to appropriate types.
func (l *Loader) convertValue(value string) any {
	// Try to convert to boolean
	if strings.ToLower(value) == "true" {
		return true
	}
	if strings.ToLower(value) == "false" {
		return false
	}

	// Try to convert to integer
	if intVal, err := strconv.Atoi(value); err == nil {
		return intVal
	}

	// Try to convert to float
	if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
		return floatVal
	}

	// Return as string if no conversion possible
	return value
}

// setNestedValue sets a value in a nested map structure using dot notation.
func (l *Loader) setNestedValue(values map[string]any, key string, value any) error {
	keys := strings.Split(key, ".")

	// Navigate to the parent of the final key
	current := values
	for i := 0; i < len(keys)-1; i++ {
		k := keys[i]

		if _, exists := current[k]; !exists {
			current[k] = make(map[string]any)
		}

		// Type assertion to ensure we have a map
		if nextMap, ok := current[k].(map[string]any); ok {
			current = nextMap
		} else {
			return fmt.Errorf("key %s is not a map, cannot set nested value", strings.Join(keys[:i+1], "."))
		}
	}

	// Set the final value
	finalKey := keys[len(keys)-1]
	current[finalKey] = value

	return nil
}

// toCamelCase converts UPPER_CASE_WITH_UNDERSCORES to camelCase.
func (l *Loader) toCamelCase(s string) string {
	if s == "" {
		return s
	}

	// Split by underscores
	parts := strings.Split(strings.ToLower(s), "_")

	// First part stays lowercase, capitalize first letter of subsequent parts
	result := parts[0]
	for i := 1; i < len(parts); i++ {
		if len(parts[i]) > 0 {
			result += strings.ToUpper(string(parts[i][0])) + parts[i][1:]
		}
	}

	return result
}

// deepMerge recursively merges source map into destination map.
func (l *Loader) deepMerge(dst, src map[string]any) {
	for k, v := range src {
		if srcMap, srcIsMap := v.(map[string]any); srcIsMap {
			if dstMap, dstIsMap := dst[k].(map[string]any); dstIsMap {
				// Both are maps, merge recursively
				l.deepMerge(dstMap, srcMap)
			} else {
				// Destination is not a map, replace it
				newMap := make(map[string]any)
				l.deepMerge(newMap, srcMap)
				dst[k] = newMap
			}
		} else {
			// Not a map, just set the value
			dst[k] = v
		}
	}
}
