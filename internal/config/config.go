package config

// Config holds the configuration for template processing
type Config struct {
	TemplateFile string
	ValuesFile   string
	OutputFile   string
	SetValues    []string
	Values       map[string]any
	IsDirectory  bool
	StrictMode   bool
}

// NewConfig creates a new configuration instance
func NewConfig(templateFile, valuesFile, outputFile string, setValues []string, isDirectory, strictMode bool) *Config {
	return &Config{
		TemplateFile: templateFile,
		ValuesFile:   valuesFile,
		OutputFile:   outputFile,
		SetValues:    setValues,
		Values:       make(map[string]any),
		IsDirectory:  isDirectory,
		StrictMode:   strictMode,
	}
}