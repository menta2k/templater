package processor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/menta2k/templater/internal/config"
	templatepkg "github.com/menta2k/templater/internal/template"
)

func TestNewTemplateProcessor(t *testing.T) {
	cfg := config.NewConfig("test.tpl", "values.yaml", "output.txt", []string{}, false, false)
	processor := NewTemplateProcessor(cfg)

	if processor == nil {
		t.Error("NewTemplateProcessor should return a non-nil processor")
	}
	if processor.config != cfg {
		t.Error("Processor should store the provided config")
	}
	if processor.valuesLoader == nil {
		t.Error("Processor should have a values loader")
	}
}

func TestProcessTemplatePath(t *testing.T) {
	cfg := config.NewConfig("", "", "", []string{}, false, false)
	processor := NewTemplateProcessor(cfg)

	tests := []struct {
		name         string
		pathTemplate string
		values       map[string]interface{}
		expected     string
		wantError    bool
	}{
		{
			name:         "simple path",
			pathTemplate: "{{.app.name}}/config.txt",
			values:       map[string]interface{}{"app": map[string]interface{}{"name": "myapp"}},
			expected:     "myapp/config.txt",
			wantError:    false,
		},
		{
			name:         "multiple variables",
			pathTemplate: "{{.environment}}/{{.app.name}}/deployment.yaml",
			values: map[string]interface{}{
				"environment": "production",
				"app":         map[string]interface{}{"name": "myapp"},
			},
			expected:  "production/myapp/deployment.yaml",
			wantError: false,
		},
		{
			name:         "with functions",
			pathTemplate: "{{upper .app.name}}/{{lower .environment}}.txt",
			values: map[string]interface{}{
				"app":         map[string]interface{}{"name": "myapp"},
				"environment": "PRODUCTION",
			},
			expected:  "MYAPP/production.txt",
			wantError: false,
		},
		{
			name:         "invalid template",
			pathTemplate: "{{.invalid.syntax",
			values:       map[string]interface{}{},
			expected:     "",
			wantError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := processor.processTemplatePath(tt.pathTemplate, tt.values)
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
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestEnsureOutputDir(t *testing.T) {
	cfg := config.NewConfig("", "", "", []string{}, false, false)
	processor := NewTemplateProcessor(cfg)

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test-output-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	outputPath := filepath.Join(tempDir, "nested", "deep", "path", "file.txt")

	err = processor.ensureOutputDir(outputPath)
	if err != nil {
		t.Errorf("ensureOutputDir failed: %v", err)
	}

	// Check that the directory structure was created
	expectedDir := filepath.Dir(outputPath)
	if _, err := os.Stat(expectedDir); os.IsNotExist(err) {
		t.Errorf("Expected directory %s was not created", expectedDir)
	}
}

func TestProcessTemplateFile(t *testing.T) {
	cfg := config.NewConfig("", "", "", []string{}, false, false)
	processor := NewTemplateProcessor(cfg)

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test-template-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test template file
	templateContent := `App: {{.app.name}}
Version: {{.app.version}}
Debug: {{default false .debug}}
Upper: {{upper .app.name}}`

	templatePath := filepath.Join(tempDir, "test.tpl")
	err = os.WriteFile(templatePath, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write template file: %v", err)
	}

	outputPath := filepath.Join(tempDir, "output.txt")
	templateFile := templatepkg.TemplateFile{
		SourcePath:   templatePath,
		RelativePath: "test.tpl",
		OutputPath:   outputPath,
	}

	values := map[string]interface{}{
		"app": map[string]interface{}{
			"name":    "test-app",
			"version": "1.0.0",
		},
		"debug": true,
	}

	err = processor.processTemplateFile(templateFile, values)
	if err != nil {
		t.Errorf("processTemplateFile failed: %v", err)
	}

	// Check that output file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}

	// Check output content
	outputContent, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	expectedSubstrings := []string{
		"App: test-app",
		"Version: 1.0.0",
		"Debug: true",
		"Upper: TEST-APP",
	}

	for _, expected := range expectedSubstrings {
		if !strings.Contains(string(outputContent), expected) {
			t.Errorf("Expected output to contain %s, but got: %s", expected, string(outputContent))
		}
	}
}

func TestFindTemplateFiles(t *testing.T) {
	cfg := config.NewConfig("", "", "", []string{}, false, false)
	processor := NewTemplateProcessor(cfg)

	// Create a temporary directory structure for testing
	tempDir, err := os.MkdirTemp("", "test-templates-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test template files
	templateFiles := []string{
		"config.tpl",
		"srv/myapp/config.tpl",
		"docs/readme.txt", // not a template
		"nested/deep/template.tpl",
	}

	for _, file := range templateFiles {
		fullPath := filepath.Join(tempDir, file)
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)
		if err != nil {
			t.Fatalf("Failed to create directory for %s: %v", file, err)
		}
		err = os.WriteFile(fullPath, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", file, err)
		}
	}

	outputDir := filepath.Join(tempDir, "output")
	values := map[string]interface{}{"app": map[string]interface{}{"name": "myapp"}}

	foundFiles, err := processor.findTemplateFiles(tempDir, outputDir, values)
	if err != nil {
		t.Errorf("findTemplateFiles failed: %v", err)
	}

	// Should find 3 .tpl files (excluding docs/readme.txt)
	expectedCount := 3
	if len(foundFiles) != expectedCount {
		t.Errorf("Expected %d template files, got %d", expectedCount, len(foundFiles))
	}

	// Check that all found files have .tpl extension in source path
	for _, file := range foundFiles {
		if !strings.HasSuffix(file.SourcePath, ".tpl") {
			t.Errorf("Found file %s should have .tpl extension", file.SourcePath)
		}
		if strings.HasSuffix(file.OutputPath, ".tpl") {
			t.Errorf("Output path %s should not have .tpl extension", file.OutputPath)
		}
	}
}

func TestProcessSingleFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test-single-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test template file
	templateContent := `Hello {{.name}}!`
	templatePath := filepath.Join(tempDir, "test.tpl")
	err = os.WriteFile(templatePath, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write template file: %v", err)
	}

	// Create a test values file
	valuesContent := `name: World`
	valuesPath := filepath.Join(tempDir, "values.yaml")
	err = os.WriteFile(valuesPath, []byte(valuesContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write values file: %v", err)
	}

	outputPath := filepath.Join(tempDir, "output.txt")

	cfg := config.NewConfig(templatePath, valuesPath, outputPath, []string{"name=World"}, false, false)
	processor := NewTemplateProcessor(cfg)

	err = processor.Process()
	if err != nil {
		t.Errorf("Process failed: %v", err)
	}

	// Check that output file was created with correct content
	outputContent, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	expected := "Hello World!"
	actual := strings.TrimSpace(string(outputContent))
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestProcessDirectory(t *testing.T) {
	// Create a temporary directory structure for testing
	tempDir, err := os.MkdirTemp("", "test-directory-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create template directory structure
	templateDir := filepath.Join(tempDir, "templates")
	err = os.MkdirAll(filepath.Join(templateDir, "{{.app.name}}"), 0755)
	if err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}

	// Create test template files
	templates := map[string]string{
		"config.tpl":                   "Config for {{.app.name}}",
		"{{.app.name}}/deployment.tpl": "Deploy {{.app.name}} version {{.app.version}}",
	}

	for file, content := range templates {
		fullPath := filepath.Join(templateDir, file)
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)
		if err != nil {
			t.Fatalf("Failed to create directory for %s: %v", file, err)
		}
		err = os.WriteFile(fullPath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create template file %s: %v", file, err)
		}
	}

	// Create values file
	valuesContent := `app:
  name: myapp
  version: 1.0.0`
	valuesPath := filepath.Join(tempDir, "values.yaml")
	err = os.WriteFile(valuesPath, []byte(valuesContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write values file: %v", err)
	}

	outputDir := filepath.Join(tempDir, "output")

	cfg := config.NewConfig(templateDir, valuesPath, outputDir, []string{}, true, false)
	processor := NewTemplateProcessor(cfg)

	err = processor.Process()
	if err != nil {
		t.Errorf("Process failed: %v", err)
	}

	// Check that output files were created
	expectedFiles := []string{
		filepath.Join(outputDir, "config"),
		filepath.Join(outputDir, "myapp", "deployment"),
	}

	for _, expectedFile := range expectedFiles {
		if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
			t.Errorf("Expected output file %s was not created", expectedFile)
		}
	}

	// Check content of one of the files
	configContent, err := os.ReadFile(filepath.Join(outputDir, "config"))
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	expectedContent := "Config for myapp"
	if strings.TrimSpace(string(configContent)) != expectedContent {
		t.Errorf("Expected %s, got %s", expectedContent, strings.TrimSpace(string(configContent)))
	}
}

func TestStrictModeProcessing(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test-strict-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name        string
		template    string
		values      string
		strictMode  bool
		expectError bool
		errorString string
	}{
		{
			name:        "defined variable - strict mode",
			template:    "Hello {{.name}}!",
			values:      "name: world",
			strictMode:  true,
			expectError: false,
		},
		{
			name:        "defined variable - normal mode",
			template:    "Hello {{.name}}!",
			values:      "name: world",
			strictMode:  false,
			expectError: false,
		},
		{
			name:        "undefined variable - strict mode",
			template:    "Hello {{.undefined}}!",
			values:      "name: world",
			strictMode:  true,
			expectError: true,
			errorString: "strict mode error",
		},
		{
			name:        "undefined variable - normal mode",
			template:    "Hello {{.undefined}}!",
			values:      "name: world",
			strictMode:  false,
			expectError: false,
		},
		{
			name:        "nested undefined variable - strict mode",
			template:    "App: {{.app.undefined}}",
			values:      "app:\n  name: myapp",
			strictMode:  true,
			expectError: true,
			errorString: "strict mode error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create template file
			templatePath := filepath.Join(tempDir, tt.name+".tpl")
			err := os.WriteFile(templatePath, []byte(tt.template), 0644)
			if err != nil {
				t.Fatalf("Failed to create template file: %v", err)
			}

			// Create values file
			valuesPath := filepath.Join(tempDir, tt.name+".yaml")
			err = os.WriteFile(valuesPath, []byte(tt.values), 0644)
			if err != nil {
				t.Fatalf("Failed to create values file: %v", err)
			}

			outputPath := filepath.Join(tempDir, tt.name+".out")

			cfg := config.NewConfig(templatePath, valuesPath, outputPath, []string{}, false, tt.strictMode)
			processor := NewTemplateProcessor(cfg)

			err = processor.Process()

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
					return
				}
				if tt.errorString != "" && !strings.Contains(err.Error(), tt.errorString) {
					t.Errorf("Expected error to contain '%s', got: %v", tt.errorString, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestStrictModePathProcessing(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "test-strict-path-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create template directory structure
	templateDir := filepath.Join(tempDir, "templates")
	err = os.MkdirAll(filepath.Join(templateDir, "{{.environment}}"), 0755)
	if err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}

	// Create template file with templated path
	templatePath := filepath.Join(templateDir, "{{.environment}}", "config.tpl")
	templateContent := "Config for {{.app.name}}"
	err = os.WriteFile(templatePath, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create template file: %v", err)
	}

	tests := []struct {
		name        string
		values      string
		strictMode  bool
		expectError bool
		errorString string
	}{
		{
			name: "defined path variable - strict mode",
			values: `app:
  name: myapp
environment: production`,
			strictMode:  true,
			expectError: false,
		},
		{
			name: "undefined path variable - strict mode",
			values: `app:
  name: myapp`,
			strictMode:  true,
			expectError: true,
			errorString: "strict mode error in path template",
		},
		{
			name: "undefined path variable - normal mode",
			values: `app:
  name: myapp`,
			strictMode:  false,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create values file
			valuesPath := filepath.Join(tempDir, tt.name+".yaml")
			err := os.WriteFile(valuesPath, []byte(tt.values), 0644)
			if err != nil {
				t.Fatalf("Failed to create values file: %v", err)
			}

			outputDir := filepath.Join(tempDir, tt.name+"-output")

			cfg := config.NewConfig(templateDir, valuesPath, outputDir, []string{}, true, tt.strictMode)
			processor := NewTemplateProcessor(cfg)

			err = processor.Process()

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
					return
				}
				if tt.errorString != "" && !strings.Contains(err.Error(), tt.errorString) {
					t.Errorf("Expected error to contain '%s', got: %v", tt.errorString, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}