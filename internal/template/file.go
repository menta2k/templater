package template

// TemplateFile represents a template file to be processed
type TemplateFile struct {
	SourcePath   string // Full path to the source template file
	RelativePath string // Relative path from the template directory
	OutputPath   string // Full path for the output file
}