# Templater

A powerful, Helm-like template processing tool for Go that supports advanced template functions, multiple data formats, and strict validation.

## Features

- 🚀 **100+ Template Functions** - Full Sprig library integration
- 📄 **Multi-format Support** - JSON, YAML, TOML conversions
- 🔍 **Strict Mode** - Exit on undefined variables for production safety
- 📁 **Directory Processing** - Recursive template processing with dynamic paths
- 🔧 **Multiple Value Sources** - YAML files, environment variables, command-line values
- 🎯 **Production Ready** - Comprehensive error handling and validation

## Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/menta2k/templater.git
cd templater

# Build the binary
go build ./cmd/templater

# Or run directly
go run ./cmd/templater -help
```

### Basic Usage

```bash
# Process a single template file
./templater -template config.tpl -values values.yaml -output config.txt

# Process a directory of templates
./templater -template ./templates -values values.yaml -output ./output

# Use command-line values
./templater -template config.tpl --set app.name=myapp --set debug=true

# Enable strict mode (exit on undefined variables)
./templater -template config.tpl -values values.yaml --strict
```

## Template Syntax

### Basic Variables

```yaml
# template.tpl
App Name: {{.app.name}}
Version: {{.app.version}}
Debug: {{default false .debug}}
```

### Sprig Functions (100+ available)

```yaml
# String functions
Quoted: {{.name | quote}}
Upper: {{.name | upper}}
Title: {{.name | title}}
Slugified: {{.name | replace " " "-" | lower}}

# Date/time functions
Current: {{now | date "2006-01-02 15:04:05"}}
Custom: {{.timestamp | date "Monday, January 2, 2006"}}

# Math functions
Total: {{add .cpu .memory}}
Percentage: {{mul .ratio 100}}

# Logic functions
{{- if eq .environment "production"}}
Production Environment
{{- else}}
Development Environment
{{- end}}
```

### Format Conversions

```yaml
# JSON conversion
Config: {{.config | toJson}}
Pretty: {{.config | mustToJson}}

# YAML conversion
Database:
{{.database | toYaml | indent 2}}

# TOML conversion
[app]
{{.app | toToml}}

# Parse formats
{{$config := fromJson .jsonString}}
Host: {{$config.host}}
```

## Value Sources

Values are merged from multiple sources in order of precedence:

### 1. Command-line --set values (highest precedence)

```bash
./templater -template app.tpl --set app.name=myapp --set debug=true
```

### 2. Environment variables (converted to camelCase)

```bash
export DATABASE_HOST=localhost
export MAX_CONNECTIONS=100
# Becomes: databaseHost and maxConnections in templates
```

### 3. YAML values file (lowest precedence)

```yaml
# values.yaml
app:
  name: myapp
  version: 1.0.0
database:
  host: localhost
  port: 5432
```

## Directory Processing

Process entire directory trees with templated paths:

```
templates/
├── config.tpl
├── {{.environment}}/
│   └── {{.app.name}}/
│       └── deployment.tpl
└── services/
    └── {{.app.name}}/
        └── config.tpl
```

```bash
./templater -template ./templates -values values.yaml -output ./output \
  --set app.name=myapp --set environment=production
```

Output:
```
output/
├── config
├── production/
│   └── myapp/
│       └── deployment
└── services/
    └── myapp/
        └── config
```

## Strict Mode

Enable strict validation to catch undefined variables:

```bash
# Normal mode - undefined variables become "<no value>"
./templater -template config.tpl -values values.yaml

# Strict mode - exit on undefined variables
./templater -template config.tpl -values values.yaml --strict
```

**Example Error:**
```
Error: strict mode error in config.tpl: undefined variable 'database' in template 'config.tpl' (strict mode enabled)
```

**Benefits:**
- Catch configuration errors early
- Prevent silent failures in production
- Ensure all required variables are provided
- Perfect for CI/CD pipelines

## Advanced Examples

### Complete Application Template

```yaml
# app-template.tpl
# {{.app.name}} Configuration
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{.app.name}}-config
  namespace: {{.namespace | default "default"}}
data:
  app.properties: |
    app.name={{.app.name}}
    app.version={{.app.version}}
    debug={{default false .debug}}
    
    # Database configuration
    {{- if .database}}
    database.host={{.database.host}}
    database.port={{.database.port | default 5432}}
    database.ssl={{.database.ssl | default true}}
    {{- end}}
    
    # Feature flags
    {{- range $key, $value := .features}}
    feature.{{$key}}={{$value}}
    {{- end}}

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.app.name}}
  namespace: {{.namespace | default "default"}}
spec:
  replicas: {{.replicas | default 1}}
  selector:
    matchLabels:
      app: {{.app.name}}
  template:
    metadata:
      labels:
        app: {{.app.name}}
        version: {{.app.version}}
    spec:
      containers:
      - name: {{.app.name}}
        image: {{.image.repository}}:{{.image.tag | default .app.version}}
        ports:
        - containerPort: {{.service.port | default 8080}}
        env:
        - name: APP_NAME
          value: {{.app.name | quote}}
        - name: APP_VERSION
          value: {{.app.version | quote}}
        {{- if .database}}
        - name: DATABASE_URL
          value: {{printf "postgresql://%s:%d/%s" .database.host (.database.port | default 5432) .database.name | quote}}
        {{- end}}
        resources:
          requests:
            memory: {{.resources.memory | default "256Mi"}}
            cpu: {{.resources.cpu | default "250m"}}
          limits:
            memory: {{.resources.memory | default "512Mi"}}
            cpu: {{.resources.cpu | default "500m"}}
```

### Values File

```yaml
# values.yaml
app:
  name: "my-awesome-app"
  version: "2.1.0"

namespace: "production"
replicas: 3

image:
  repository: "myregistry/myapp"
  tag: "v2.1.0"

service:
  port: 8080

database:
  host: "postgres.production.svc.cluster.local"
  port: 5432
  name: "myapp_prod"
  ssl: true

resources:
  memory: "1Gi"
  cpu: "500m"

features:
  analytics: true
  monitoring: true
  debug: false
```

### Multi-Environment Setup

```bash
# Development
./templater -template k8s-template.tpl \
  --set app.name=myapp \
  --set namespace=development \
  --set replicas=1 \
  --set database.host=localhost

# Staging  
./templater -template k8s-template.tpl -values values-staging.yaml \
  --set namespace=staging \
  --set replicas=2

# Production (with strict validation)
./templater -template k8s-template.tpl -values values-production.yaml \
  --set namespace=production \
  --strict
```

## Function Reference

### Built-in Functions

- `default` - Provide default values
- `upper` - Convert to uppercase
- `lower` - Convert to lowercase

### Sprig Functions (100+)

**String Functions:**
- `quote`, `squote` - Add quotes
- `trim`, `trimPrefix`, `trimSuffix` - String trimming
- `replace`, `regexReplace` - String replacement  
- `split`, `join` - String/array operations
- `indent`, `nindent` - Text indentation

**Math Functions:**
- `add`, `sub`, `mul`, `div` - Basic arithmetic
- `mod`, `max`, `min` - Additional math operations
- `floor`, `ceil`, `round` - Rounding functions

**Date Functions:**
- `now` - Current timestamp
- `date` - Format dates
- `dateInZone` - Timezone-aware formatting

**Encoding Functions:**
- `b64enc`, `b64dec` - Base64 encoding/decoding
- `urlquery` - URL encoding

### Format Conversion Functions

**JSON:**
- `toJson` - Convert to JSON (safe)
- `mustToJson` - Convert to JSON (panic on error)
- `fromJson` - Parse JSON to object
- `fromJsonArray` - Parse JSON to array

**YAML:**
- `toYaml` - Convert to YAML (safe)
- `mustToYaml` - Convert to YAML (panic on error)
- `toYamlPretty` - Convert to pretty YAML
- `fromYaml` - Parse YAML to object
- `fromYamlArray` - Parse YAML to array

**TOML:**
- `toToml` - Convert to TOML
- `fromToml` - Parse TOML to object

## Command Line Options

```
Usage: ./templater [options]

Options:
  -template string
        Path to the template file or directory (required)
  -values string
        Path to the YAML values file (optional)
  -output string
        Path to the output file or directory (default "output")
  -set value
        Set values on the command line (can be used multiple times or comma-separated)
  -strict
        Enable strict mode - exit on undefined values
  -help
        Show help message
```

## Use Cases

### Configuration Management
- Kubernetes manifests
- Docker Compose files
- Application configuration files
- Environment-specific settings

### Infrastructure as Code
- Terraform variable files
- Ansible playbooks
- CloudFormation templates
- Helm chart alternatives

### CI/CD Pipelines
- Build configuration generation
- Deployment manifests
- Environment promotion
- Multi-stage deployments

### Development Workflows
- Local development setup
- Testing configurations
- Documentation generation
- Code templating

## Security

- **Sandboxed Execution** - No file system or network access from templates
- **Safe Functions** - Dangerous functions (`env`, `expandenv`) are disabled
- **Input Validation** - Comprehensive error handling and validation
- **No Code Execution** - Templates cannot execute arbitrary code

## Performance

- **Optimized for Speed** - Efficient template processing
- **Memory Efficient** - Handles large templates and data sets
- **Concurrent Processing** - Fast directory processing
- **Minimal Dependencies** - Lightweight binary

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/new-feature`)
3. Write tests for your changes
4. Ensure all tests pass (`go test ./...`)
5. Commit your changes (`git commit -am 'Add new feature'`)
6. Push to the branch (`git push origin feature/new-feature`)
7. Create a Pull Request

## Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test suites
go test ./internal/template -v
go test ./internal/processor -v
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Sprig](https://masterminds.github.io/sprig/) - Template function library
- [Go Templates](https://pkg.go.dev/text/template) - Core template engine
- [Helm](https://helm.sh/) - Inspiration for template functionality

---

## Examples Repository

Check out the `example/` directory for more templates and use cases:

- `enhanced-template.tpl` - Showcase of all functions
- `strict-demo.tpl` - Strict mode demonstration  
- `enhanced-values.yaml` - Comprehensive values file
- Various Kubernetes manifests and configuration examples

Happy templating! 🚀