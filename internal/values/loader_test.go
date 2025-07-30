package values

import (
	"os"
	"reflect"
	"testing"
)

func TestNewLoader(t *testing.T) {
	loader := NewLoader()
	if loader == nil {
		t.Error("NewLoader should return a non-nil loader")
	}
}

func TestLoadYAMLValues(t *testing.T) {
	loader := NewLoader()

	// Test with empty file path
	values, err := loader.LoadYAMLValues("")
	if err != nil {
		t.Errorf("Expected no error for empty file path, got %v", err)
	}
	if len(values) != 0 {
		t.Error("Expected empty values for empty file path")
	}

	// Create a temporary YAML file
	content := `app:
  name: test-app
  version: 1.0.0
database:
  host: localhost
  port: 5432
`
	tmpFile, err := os.CreateTemp("", "test-values-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	values, err = loader.LoadYAMLValues(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to load YAML values: %v", err)
	}

	expected := map[string]interface{}{
		"app": map[interface{}]interface{}{
			"name":    "test-app",
			"version": "1.0.0",
		},
		"database": map[interface{}]interface{}{
			"host": "localhost",
			"port": 5432,
		},
	}

	if !reflect.DeepEqual(values["app"], expected["app"]) {
		t.Errorf("Expected app values %v, got %v", expected["app"], values["app"])
	}
}

func TestLoadEnvValues(t *testing.T) {
	loader := NewLoader()

	// Set some test environment variables
	os.Setenv("TEST_APP_NAME", "test-app")
	os.Setenv("DATABASE_HOST", "localhost")
	os.Setenv("MAX_CONNECTIONS", "100")
	defer func() {
		os.Unsetenv("TEST_APP_NAME")
		os.Unsetenv("DATABASE_HOST")
		os.Unsetenv("MAX_CONNECTIONS")
	}()

	values := loader.LoadEnvValues()

	// Check that environment variables were converted to camelCase
	if values["testAppName"] != "test-app" {
		t.Errorf("Expected testAppName to be 'test-app', got %v", values["testAppName"])
	}
	if values["databaseHost"] != "localhost" {
		t.Errorf("Expected databaseHost to be 'localhost', got %v", values["databaseHost"])
	}
	if values["maxConnections"] != "100" {
		t.Errorf("Expected maxConnections to be '100', got %v", values["maxConnections"])
	}
}

func TestToCamelCase(t *testing.T) {
	loader := NewLoader()

	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"TEST", "test"},
		{"TEST_APP", "testApp"},
		{"DATABASE_HOST_NAME", "databaseHostName"},
		{"MAX_CONNECTIONS", "maxConnections"},
		{"A_B_C_D", "aBCD"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := loader.toCamelCase(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestParseSetValues(t *testing.T) {
	loader := NewLoader()

	tests := []struct {
		name      string
		setValues []string
		expected  map[string]interface{}
		wantError bool
	}{
		{
			name:      "simple key-value",
			setValues: []string{"key=value"},
			expected:  map[string]interface{}{"key": "value"},
			wantError: false,
		},
		{
			name:      "multiple values",
			setValues: []string{"key1=value1,key2=value2"},
			expected:  map[string]interface{}{"key1": "value1", "key2": "value2"},
			wantError: false,
		},
		{
			name:      "nested keys",
			setValues: []string{"app.name=test-app"},
			expected: map[string]interface{}{
				"app": map[string]interface{}{"name": "test-app"},
			},
			wantError: false,
		},
		{
			name:      "boolean values",
			setValues: []string{"debug=true,production=false"},
			expected:  map[string]interface{}{"debug": true, "production": false},
			wantError: false,
		},
		{
			name:      "integer values",
			setValues: []string{"port=8080,maxConnections=100"},
			expected:  map[string]interface{}{"port": 8080, "maxConnections": 100},
			wantError: false,
		},
		{
			name:      "float values",
			setValues: []string{"version=1.5,ratio=0.75"},
			expected:  map[string]interface{}{"version": 1.5, "ratio": 0.75},
			wantError: false,
		},
		{
			name:      "invalid format",
			setValues: []string{"invalidformat"},
			expected:  nil,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := loader.ParseSetValues(tt.setValues)
			if tt.wantError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestConvertValue(t *testing.T) {
	loader := NewLoader()

	tests := []struct {
		input    string
		expected interface{}
	}{
		{"true", true},
		{"false", false},
		{"True", true},
		{"FALSE", false},
		{"123", 123},
		{"0", 0},
		{"-42", -42},
		{"3.14", 3.14},
		{"0.5", 0.5},
		{"hello", "hello"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := loader.convertValue(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %v (%T), got %v (%T)", tt.expected, tt.expected, result, result)
			}
		})
	}
}

func TestSetNestedValue(t *testing.T) {
	loader := NewLoader()

	tests := []struct {
		name      string
		key       string
		value     interface{}
		initial   map[string]interface{}
		expected  map[string]interface{}
		wantError bool
	}{
		{
			name:     "simple key",
			key:      "key",
			value:    "value",
			initial:  make(map[string]interface{}),
			expected: map[string]interface{}{"key": "value"},
		},
		{
			name:    "nested key",
			key:     "app.name",
			value:   "test-app",
			initial: make(map[string]interface{}),
			expected: map[string]interface{}{
				"app": map[string]interface{}{"name": "test-app"},
			},
		},
		{
			name:    "deep nested key",
			key:     "database.config.host",
			value:   "localhost",
			initial: make(map[string]interface{}),
			expected: map[string]interface{}{
				"database": map[string]interface{}{
					"config": map[string]interface{}{"host": "localhost"},
				},
			},
		},
		{
			name:  "existing nested structure",
			key:   "app.version",
			value: "2.0",
			initial: map[string]interface{}{
				"app": map[string]interface{}{"name": "test-app"},
			},
			expected: map[string]interface{}{
				"app": map[string]interface{}{"name": "test-app", "version": "2.0"},
			},
		},
		{
			name:      "conflict with non-map",
			key:       "app.version",
			value:     "2.0",
			initial:   map[string]interface{}{"app": "not-a-map"},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := loader.setNestedValue(tt.initial, tt.key, tt.value)
			if tt.wantError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if !reflect.DeepEqual(tt.initial, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, tt.initial)
			}
		})
	}
}

func TestMergeValues(t *testing.T) {
	loader := NewLoader()

	yamlValues := map[string]interface{}{
		"app": map[string]interface{}{
			"name":    "yaml-app",
			"version": "1.0",
		},
		"debug": false,
	}

	envValues := map[string]interface{}{
		"app": map[string]interface{}{
			"name": "env-app",
		},
		"port": 8080,
	}

	setValues := map[string]interface{}{
		"app": map[string]interface{}{
			"version": "2.0",
		},
		"debug": true,
	}

	configValues := map[string]interface{}{
		"maxConnections": 100,
	}

	result := loader.MergeValues(yamlValues, envValues, setValues, configValues)

	expected := map[string]interface{}{
		"app": map[string]interface{}{
			"name":    "env-app", // from env (overrides yaml)
			"version": "2.0",     // from set (overrides yaml)
		},
		"debug":          true, // from set (overrides yaml)
		"port":           8080, // from env
		"maxConnections": 100,  // from config
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestDeepMerge(t *testing.T) {
	loader := NewLoader()

	dst := map[string]interface{}{
		"app": map[string]interface{}{
			"name":    "original",
			"version": "1.0",
		},
		"debug": false,
	}

	src := map[string]interface{}{
		"app": map[string]interface{}{
			"name": "updated",
			"port": 8080,
		},
		"production": true,
	}

	loader.deepMerge(dst, src)

	expected := map[string]interface{}{
		"app": map[string]interface{}{
			"name":    "updated", // overridden
			"version": "1.0",     // preserved
			"port":    8080,      // added
		},
		"debug":      false, // preserved
		"production": true,  // added
	}

	if !reflect.DeepEqual(dst, expected) {
		t.Errorf("Expected %v, got %v", expected, dst)
	}
}
