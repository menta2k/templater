package processor

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/menta2k/templater/internal/config"
	templatepkg "github.com/menta2k/templater/internal/template"
	"github.com/menta2k/templater/internal/values"
)

// TemplateProcessor handles template processing operations.
type TemplateProcessor struct {
	config       *config.Config
	valuesLoader *values.Loader
}

// NewTemplateProcessor creates a new template processor.
func NewTemplateProcessor(cfg *config.Config) *TemplateProcessor {
	return &TemplateProcessor{
		config:       cfg,
		valuesLoader: values.NewLoader(),
	}
}

// Process processes the template(s) with merged values.
func (tp *TemplateProcessor) Process() error {
	// Load values from YAML file
	yamlValues, err := tp.valuesLoader.LoadYAMLValues(tp.config.ValuesFile)
	if err != nil {
		return fmt.Errorf("error loading YAML values: %w", err)
	}

	// Load values from environment variables
	envValues := tp.valuesLoader.LoadEnvValues()

	// Parse --set values
	setValues, err := tp.valuesLoader.ParseSetValues(tp.config.SetValues)
	if err != nil {
		return fmt.Errorf("error parsing set values: %w", err)
	}

	// Merge all values (--set values have highest precedence)
	allValues := tp.valuesLoader.MergeValues(yamlValues, envValues, setValues, tp.config.Values)

	// Check if template is a directory or file
	fileInfo, err := os.Stat(tp.config.TemplateFile)
	if err != nil {
		return fmt.Errorf("failed to stat template path: %w", err)
	}

	if fileInfo.IsDir() {
		// Process directory of templates
		return tp.processDirectory(allValues)
	} else {
		// Process single template file
		return tp.processSingleFile(allValues)
	}
}

// processTemplatePath processes a path that may contain template variables.
func (tp *TemplateProcessor) processTemplatePath(pathTemplate string, allValues map[string]any) (string, error) {
	// Create strict template wrapper for path processing
	strictTemplate := templatepkg.NewStrictTemplate("path", tp.config.StrictMode)

	// Parse path template
	parsedTemplate, err := strictTemplate.ParseTemplate(pathTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse path template '%s': %w", pathTemplate, err)
	}

	// Execute path template with strict mode support
	result, err := parsedTemplate.ExecuteTemplate(allValues)
	if err != nil {
		// Check if it's a strict mode error
		var strictErr *templatepkg.StrictModeError
		if errors.As(err, &strictErr) {
			return "", fmt.Errorf("strict mode error in path template '%s': %s", pathTemplate, strictErr.Error())
		}
		return "", fmt.Errorf("failed to execute path template '%s': %w", pathTemplate, err)
	}

	return result, nil
}

// findTemplateFiles recursively finds all *.tpl files in a directory.
func (tp *TemplateProcessor) findTemplateFiles(templateDir, outputDir string, allValues map[string]any) ([]templatepkg.File, error) {
	var templateFiles []templatepkg.File

	err := filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if file has .tpl extension
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".tpl") {
			// Calculate relative path from template directory
			relativePath, err := filepath.Rel(templateDir, path)
			if err != nil {
				return fmt.Errorf("failed to calculate relative path: %w", err)
			}

			// Process the relative path as a template to handle templated directory names
			processedRelativePath, err := tp.processTemplatePath(relativePath, allValues)
			if err != nil {
				return fmt.Errorf("failed to process path template '%s': %w", relativePath, err)
			}

			// Create output path by replacing .tpl extension and joining with output directory
			outputName := strings.TrimSuffix(processedRelativePath, ".tpl")
			outputPath := filepath.Join(outputDir, outputName)

			templateFiles = append(templateFiles, templatepkg.File{
				SourcePath:   path,
				RelativePath: relativePath,
				OutputPath:   outputPath,
			})
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error walking template directory: %w", err)
	}

	return templateFiles, nil
}

// ensureOutputDir creates the output directory structure for a file.
func (tp *TemplateProcessor) ensureOutputDir(outputPath string) error {
	outputDir := filepath.Dir(outputPath)
	return os.MkdirAll(outputDir, 0o755)
}

// processTemplateFile processes a single template file.
func (tp *TemplateProcessor) processTemplateFile(templateFile templatepkg.File, allValues map[string]any) error {
	// Load and parse template
	templateContent, err := os.ReadFile(templateFile.SourcePath)
	if err != nil {
		return fmt.Errorf("failed to read template file %s: %w", templateFile.SourcePath, err)
	}

	// Create strict template wrapper
	strictTemplate := templatepkg.NewStrictTemplate(filepath.Base(templateFile.SourcePath), tp.config.StrictMode)

	// Parse template content
	parsedTemplate, err := strictTemplate.ParseTemplate(string(templateContent))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", templateFile.SourcePath, err)
	}

	// Ensure output directory exists
	err = tp.ensureOutputDir(templateFile.OutputPath)
	if err != nil {
		return fmt.Errorf("failed to create output directory for %s: %w", templateFile.OutputPath, err)
	}

	// Execute template with strict mode support
	result, err := parsedTemplate.ExecuteTemplate(allValues)
	if err != nil {
		// Check if it's a strict mode error
		var strictErr *templatepkg.StrictModeError
		if errors.As(err, &strictErr) {
			return fmt.Errorf("strict mode error in %s: %s", templateFile.SourcePath, strictErr.Error())
		}
		return fmt.Errorf("failed to execute template %s: %w", templateFile.SourcePath, err)
	}

	// Create output file
	outputFile, err := os.Create(templateFile.OutputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %w", templateFile.OutputPath, err)
	}
	defer outputFile.Close()

	// Write result to file
	_, err = outputFile.WriteString(result)
	if err != nil {
		return fmt.Errorf("failed to write output file %s: %w", templateFile.OutputPath, err)
	}

	fmt.Printf("Processed: %s -> %s\n", templateFile.RelativePath, templateFile.OutputPath)
	return nil
}

// processDirectory processes all *.tpl files in a directory recursively.
func (tp *TemplateProcessor) processDirectory(allValues map[string]any) error {
	templateDir := tp.config.TemplateFile
	outputDir := tp.config.OutputFile

	// Find all template files (now with templated path processing)
	templateFiles, err := tp.findTemplateFiles(templateDir, outputDir, allValues)
	if err != nil {
		return err
	}

	if len(templateFiles) == 0 {
		fmt.Printf("No *.tpl files found in directory: %s\n", templateDir)
		return nil
	}

	fmt.Printf("Found %d template file(s) in directory: %s\n", len(templateFiles), templateDir)

	// Process each template file
	for _, templateFile := range templateFiles {
		err := tp.processTemplateFile(templateFile, allValues)
		if err != nil {
			return err
		}
	}

	fmt.Printf("\nSuccessfully processed %d template file(s). Output directory: %s\n", len(templateFiles), outputDir)
	return nil
}

// processSingleFile processes a single template file.
func (tp *TemplateProcessor) processSingleFile(allValues map[string]any) error {
	templateFile := templatepkg.File{
		SourcePath:   tp.config.TemplateFile,
		RelativePath: filepath.Base(tp.config.TemplateFile),
		OutputPath:   tp.config.OutputFile,
	}

	err := tp.processTemplateFile(templateFile, allValues)
	if err != nil {
		return err
	}

	fmt.Printf("Template processed successfully. Output written to: %s\n", tp.config.OutputFile)
	return nil
}
