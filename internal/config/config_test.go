package config

import (
	"reflect"
	"testing"
)

func TestNewConfig(t *testing.T) {
	templateFile := "test.tpl"
	valuesFile := "values.yaml"
	outputFile := "output.txt"
	setValues := []string{"key=value", "app.name=test"}
	isDirectory := true
	strictMode := true

	config := NewConfig(templateFile, valuesFile, outputFile, setValues, isDirectory, strictMode)

	if config.TemplateFile != templateFile {
		t.Errorf("Expected TemplateFile %s, got %s", templateFile, config.TemplateFile)
	}
	if config.ValuesFile != valuesFile {
		t.Errorf("Expected ValuesFile %s, got %s", valuesFile, config.ValuesFile)
	}
	if config.OutputFile != outputFile {
		t.Errorf("Expected OutputFile %s, got %s", outputFile, config.OutputFile)
	}
	if !reflect.DeepEqual(config.SetValues, setValues) {
		t.Errorf("Expected SetValues %v, got %v", setValues, config.SetValues)
	}
	if config.IsDirectory != isDirectory {
		t.Errorf("Expected IsDirectory %t, got %t", isDirectory, config.IsDirectory)
	}
	if config.StrictMode != strictMode {
		t.Errorf("Expected StrictMode %t, got %t", strictMode, config.StrictMode)
	}
	if config.Values == nil {
		t.Error("Expected Values to be initialized")
	}
}

func TestConfig_Fields(t *testing.T) {
	config := &Config{
		TemplateFile: "template.tpl",
		ValuesFile:   "values.yaml",
		OutputFile:   "output.txt",
		SetValues:    []string{"key=value"},
		Values:       map[string]any{"test": "value"},
		IsDirectory:  false,
		StrictMode:   true,
	}

	if config.TemplateFile != "template.tpl" {
		t.Error("TemplateFile field not set correctly")
	}
	if config.ValuesFile != "values.yaml" {
		t.Error("ValuesFile field not set correctly")
	}
	if config.OutputFile != "output.txt" {
		t.Error("OutputFile field not set correctly")
	}
	if len(config.SetValues) != 1 || config.SetValues[0] != "key=value" {
		t.Error("SetValues field not set correctly")
	}
	if config.Values["test"] != "value" {
		t.Error("Values field not set correctly")
	}
	if config.IsDirectory != false {
		t.Error("IsDirectory field not set correctly")
	}
	if config.StrictMode != true {
		t.Error("StrictMode field not set correctly")
	}
}
