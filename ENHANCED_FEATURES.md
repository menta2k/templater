# Enhanced Template Functions

The templater now includes powerful template functions from [Sprig](https://masterminds.github.io/sprig/) and additional format conversions for JSON, YAML, and TOML.

## Quick Start

```bash
# Basic usage with enhanced functions
go run ./cmd/templater -template example/enhanced-template.tpl -values example/enhanced-values.yaml -output output.txt

# Strict mode - exit on undefined variables
go run ./cmd/templater -template template.tpl -values values.yaml -output output.txt --strict
```

## Available Function Categories

### 1. Sprig Functions (100+ functions)

All Sprig functions are available except for potentially dangerous ones (`env`, `expandenv`).

#### String Functions
```yaml
# Template examples:
Quote: {{quote .text}}           # "hello world"
Upper: {{upper .text}}           # HELLO WORLD  
Title: {{title .text}}          # Hello World
Trim: {{trim .padded}}          # trimmed
Replace: {{replace " " "-" .text}}  # hello-world
Repeat: {{repeat 3 "x"}}        # xxx
```

#### Date/Time Functions
```yaml
Current: {{now | date "2006-01-02 15:04:05"}}
Custom: {{.timestamp | date "Monday, January 2, 2006"}}
```

#### Math Functions
```yaml
Add: {{add 1 2}}              # 3
Multiply: {{mul .replicas .cpu}}  # 6 (if replicas=3, cpu=2)
```

#### Logic Functions
```yaml
{{- if eq .environment "production"}}
Production Environment
{{- else}}
Non-Production Environment  
{{- end}}
```

### 2. Format Conversion Functions

#### JSON Conversion
```yaml
# Convert to JSON
Config: {{.config | toJson}}
Config (pretty): {{.config | mustToJson}}  # Panics on error

# Parse from JSON
{{$parsed := fromJson .jsonString}}
Name: {{$parsed.name}}

# Parse JSON array
{{range fromJsonArray .jsonArray}}
- {{.}}
{{end}}
```

#### YAML Conversion
```yaml
# Convert to YAML  
Database: |
{{.database | toYaml | indent 2}}

# Pretty YAML with custom indentation
Config: |
{{.config | toYamlPretty}}

# Parse from YAML
{{$parsed := fromYaml .yamlString}}
Host: {{$parsed.host}}

# Parse YAML array
{{range fromYamlArray .yamlArray}}  
- {{.}}
{{end}}
```

#### TOML Conversion
```yaml
# Convert to TOML
[app]
{{.app | toToml}}

# Parse from TOML
{{$parsed := fromToml .tomlString}}
Name: {{$parsed.name}}
```

### 3. Advanced String Processing

#### Indentation and Formatting
```yaml
# Indent text
Config:
{{.config | toYaml | indent 4}}

# No-indent (nindent) - useful for YAML
services:
{{.services | toYaml | nindent 2}}
```

#### Base64 Encoding/Decoding
```yaml
Encoded: {{.secret | b64enc}}
Decoded: {{.encoded | b64dec}}
```

### 4. Default Values and Conditionals

```yaml
# Default values
Debug: {{default false .debug}}
Port: {{default 8080 .port}}
LogLevel: {{default "info" .logLevel}}

# Complex conditionals
{{- if and .database.enabled .database.host}}
Database is configured and enabled
{{- end}}
```

## Example Template

See `example/enhanced-template.tpl` for a comprehensive example showing:

- String manipulation with Sprig functions
- Format conversions (JSON, YAML, TOML)
- Date/time formatting
- Conditional logic
- Loops and iterations
- Math operations
- Base64 encoding
- Advanced text processing

## Template Data Sources

Values can come from multiple sources (in order of precedence):

1. **Command-line `--set` values** (highest precedence)
   ```bash
   --set app.name=myapp --set debug=true
   ```

2. **Environment variables** (converted to camelCase)
   ```bash
   DATABASE_HOST=localhost → databaseHost
   MAX_CONNECTIONS=100 → maxConnections
   ```

3. **YAML values file** (lowest precedence)
   ```yaml
   app:
     name: myapp
     version: 1.0.0
   ```

## Security Features

- Dangerous functions (`env`, `expandenv`) are automatically removed
- Template functions are safely sandboxed
- No access to filesystem or network operations
- No code execution capabilities

## Error Handling

- Format conversion functions gracefully handle errors
- `mustToJson`, `mustToYaml` functions panic on errors (use for strict validation)  
- Invalid input returns error messages in a predictable format
- Template execution continues even with conversion errors

## Performance Notes

- Sprig functions are optimized for template use
- Format conversions handle large data structures efficiently
- Key conversion automatically handles YAML's `map[interface{}]interface{}` types
- Memory usage is optimized for typical template workloads

## Strict Mode

The `--strict` flag enables strict validation mode:

```bash
# Enable strict mode
go run ./cmd/templater -template template.tpl -values values.yaml --strict
```

**Behavior:**
- **Normal mode**: Undefined variables are replaced with `<no value>`
- **Strict mode**: Process exits with error on first undefined variable

**Benefits:**
- Catch configuration errors early
- Ensure all required variables are provided
- Prevent silent failures in production deployments

**Example:**
```yaml
# template.tpl
Hello {{.name}}!
Database: {{.database.host}}

# values.yaml
name: world
# database section missing

# Normal mode result:
Hello world!
Database: <no value>

# Strict mode result:
Error: undefined variable 'database' in template 'template.tpl' (strict mode enabled)
```

## Compatibility

- Fully backward compatible with existing templates
- All original functions (`default`, `upper`, `lower`) continue to work
- Enhanced functions extend rather than replace functionality
- Works with all existing templater features (directory processing, path templating, etc.)
- Strict mode is opt-in and does not affect existing workflows

## Function Reference

For complete Sprig function documentation, see: https://masterminds.github.io/sprig/

Custom functions added:
- `toJson`, `mustToJson`, `fromJson`, `fromJsonArray`
- `toYaml`, `mustToYaml`, `toYamlPretty`, `fromYaml`, `fromYamlArray`  
- `toToml`, `fromToml`
- `include`, `tpl`, `required`, `lookup` (placeholder functions)