package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/menta2k/templater/internal/cli"
	"github.com/menta2k/templater/internal/config"
	"github.com/menta2k/templater/internal/processor"
)

func main() {
	var (
		templateFile = flag.String("template", "", "Path to the template file or directory (required)")
		valuesFile   = flag.String("values", "", "Path to the YAML values file (optional)")
		outputFile   = flag.String("output", "output", "Path to the output file or directory")
		setVals      = cli.SetValues{}
		help         = flag.Bool("help", false, "Show help message")
		strict       = flag.Bool("strict", false, "Enable strict mode - exit on undefined values")
	)

	flag.Var(&setVals, "set", "Set values on the command line (can be used multiple times or comma-separated)")
	flag.Parse()

	if *help {
		fmt.Println("Go Template Processor - A Helm-like template processing tool")
		fmt.Println("\nUsage:")
		flag.PrintDefaults()
		fmt.Println("\nExamples:")
		fmt.Println("  # Process single template file")
		fmt.Println("  go run main.go -template=config.tmpl -values=values.yaml -output=config.txt")
		fmt.Println("  ")
		fmt.Println("  # Process directory of templates")
		fmt.Println("  go run main.go -template=./templates -values=values.yaml -output=./output")
		fmt.Println("  ")
		fmt.Println("  # Use --set values")
		fmt.Println("  go run main.go -template=./templates --set app.name=myapp,app.version=2.0")
		fmt.Println("  go run main.go -template=config.tmpl --set app.name=myapp --set debug=true")
		fmt.Println("  go run main.go -template=./templates -values=values.yaml --set database.password=secret")
		fmt.Println("  # Directory processing with templated paths")
		fmt.Println("  go run main.go -template=./templates -values=values.yaml -output=./output \\")
		fmt.Println("    --set app.name=myapp,environment=production")
		fmt.Println("  # Template: srv/{{.app.name}}/config.tpl -> Output: srv/myapp/config")
		fmt.Println("  ")
		fmt.Println("  # Strict mode - exit on undefined values")
		fmt.Println("  go run main.go -template=config.tmpl -values=values.yaml --strict")
		fmt.Println("  ")
		fmt.Println("\nTemplate discovery:")
		fmt.Println("  - Single file: processes the specified .tpl file")
		fmt.Println("  - Directory: recursively finds all *.tpl files and processes them")
		fmt.Println("  - Output maintains directory structure from source")
		fmt.Println("  - Directory paths can contain template variables (e.g., srv/{{.app.name}}/config.tpl)")
		fmt.Println("  - Templated paths are processed with the same variables as file contents")
		fmt.Println("\nTemplate variables can come from (in order of precedence):")
		fmt.Println("  1. --set values (highest precedence)")
		fmt.Println("  2. Environment variables (converted to camelCase)")
		fmt.Println("  3. YAML values file (lowest precedence)")
		fmt.Println("\nEnvironment variable conversion examples:")
		fmt.Println("  DATABASE_HOST → databaseHost")
		fmt.Println("  APP_VERSION → appVersion")
		fmt.Println("  MAX_CONNECTIONS → maxConnections")
		fmt.Println("\nSet value formats:")
		fmt.Println("  --set key=value")
		fmt.Println("  --set key1=value1,key2=value2")
		fmt.Println("  --set nested.key=value")
		fmt.Println("  --set debug=true (converts to boolean)")
		fmt.Println("  --set port=8080 (converts to integer)")
		return
	}

	if *templateFile == "" {
		fmt.Println("Error: template file or directory is required")
		fmt.Println("Use -help for usage information")
		os.Exit(1)
	}

	// Check if template file/directory exists
	if _, err := os.Stat(*templateFile); os.IsNotExist(err) {
		fmt.Printf("Error: template path '%s' does not exist\n", *templateFile)
		os.Exit(1)
	}

	// Check if values file exists (if specified)
	if *valuesFile != "" {
		if _, err := os.Stat(*valuesFile); os.IsNotExist(err) {
			fmt.Printf("Error: values file '%s' does not exist\n", *valuesFile)
			os.Exit(1)
		}
	}

	// Determine if template is a directory
	fileInfo, err := os.Stat(*templateFile)
	if err != nil {
		fmt.Printf("Error: cannot stat template path '%s': %v\n", *templateFile, err)
		os.Exit(1)
	}

	cfg := config.NewConfig(*templateFile, *valuesFile, *outputFile, []string(setVals), fileInfo.IsDir(), *strict)

	processor := processor.NewTemplateProcessor(cfg)

	err = processor.Process()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
