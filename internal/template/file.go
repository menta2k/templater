package template

// File represents a template file to be processed.
type File struct {
	SourcePath   string // Full path to the source template file
	RelativePath string // Relative path from the template directory
	OutputPath   string // Full path for the output file
}
