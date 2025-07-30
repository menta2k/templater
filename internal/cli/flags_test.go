package cli

import (
	"strings"
	"testing"
)

func TestSetValues_String(t *testing.T) {
	tests := []struct {
		name     string
		values   SetValues
		expected string
	}{
		{
			name:     "empty values",
			values:   SetValues{},
			expected: "",
		},
		{
			name:     "single value",
			values:   SetValues{"key=value"},
			expected: "key=value",
		},
		{
			name:     "multiple values",
			values:   SetValues{"key1=value1", "key2=value2", "key3=value3"},
			expected: "key1=value1,key2=value2,key3=value3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.values.String()
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestSetValues_Set(t *testing.T) {
	tests := []struct {
		name        string
		initial     SetValues
		setValue    string
		expected    SetValues
		expectError bool
	}{
		{
			name:        "add to empty",
			initial:     SetValues{},
			setValue:    "key=value",
			expected:    SetValues{"key=value"},
			expectError: false,
		},
		{
			name:        "add to existing",
			initial:     SetValues{"existing=value"},
			setValue:    "new=value",
			expected:    SetValues{"existing=value", "new=value"},
			expectError: false,
		},
		{
			name:        "add complex value",
			initial:     SetValues{},
			setValue:    "app.name=test-app",
			expected:    SetValues{"app.name=test-app"},
			expectError: false,
		},
		{
			name:        "add comma-separated values",
			initial:     SetValues{},
			setValue:    "key1=value1,key2=value2",
			expected:    SetValues{"key1=value1,key2=value2"},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values := make(SetValues, len(tt.initial))
			copy(values, tt.initial)

			err := values.Set(tt.setValue)
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(values) != len(tt.expected) {
				t.Errorf("Expected length %d, got %d", len(tt.expected), len(values))
				return
			}

			for i, expected := range tt.expected {
				if values[i] != expected {
					t.Errorf("Expected values[%d] = %s, got %s", i, expected, values[i])
				}
			}
		})
	}
}

func TestSetValues_MultipleSetCalls(t *testing.T) {
	var values SetValues

	// Simulate multiple flag calls: -set key1=value1 -set key2=value2
	err := values.Set("key1=value1")
	if err != nil {
		t.Fatalf("First Set call failed: %v", err)
	}

	err = values.Set("key2=value2")
	if err != nil {
		t.Fatalf("Second Set call failed: %v", err)
	}

	err = values.Set("nested.key=nested-value")
	if err != nil {
		t.Fatalf("Third Set call failed: %v", err)
	}

	expected := SetValues{"key1=value1", "key2=value2", "nested.key=nested-value"}
	if len(values) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(values))
	}

	for i, expectedValue := range expected {
		if values[i] != expectedValue {
			t.Errorf("Expected values[%d] = %s, got %s", i, expectedValue, values[i])
		}
	}

	// Test String() method with multiple values
	stringResult := values.String()
	expectedString := "key1=value1,key2=value2,nested.key=nested-value"
	if stringResult != expectedString {
		t.Errorf("Expected string representation %s, got %s", expectedString, stringResult)
	}
}

func TestSetValues_EmptyStringSet(t *testing.T) {
	var values SetValues

	err := values.Set("")
	if err != nil {
		t.Errorf("Setting empty string should not error: %v", err)
	}

	expected := SetValues{""}
	if len(values) != len(expected) || values[0] != expected[0] {
		t.Errorf("Expected %v, got %v", expected, values)
	}
}

func TestSetValues_SpecialCharacters(t *testing.T) {
	var values SetValues

	specialValues := []string{
		"key=value with spaces",
		"url=https://example.com:8080/path?query=value",
		"json={\"key\":\"value\"}",
		"path=/usr/local/bin",
		"regex=^[a-zA-Z0-9]+$",
	}

	for _, setValue := range specialValues {
		err := values.Set(setValue)
		if err != nil {
			t.Errorf("Failed to set value %s: %v", setValue, err)
		}
	}

	if len(values) != len(specialValues) {
		t.Errorf("Expected %d values, got %d", len(specialValues), len(values))
	}

	// Check that all values were preserved correctly
	for i, expected := range specialValues {
		if values[i] != expected {
			t.Errorf("Expected values[%d] = %s, got %s", i, expected, values[i])
		}
	}

	// Test that String() properly joins with commas
	result := values.String()
	expectedParts := len(specialValues) - 1 // Number of commas should be len-1
	actualCommas := strings.Count(result, ",")
	if actualCommas != expectedParts {
		t.Errorf("Expected %d commas in string representation, got %d", expectedParts, actualCommas)
	}
}
