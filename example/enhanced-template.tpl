# Enhanced Template Demo
# This template demonstrates the new advanced functions available

## Basic Application Info
App Name: {{.app.name | title}}
Version: {{.app.version}}
Environment: {{.environment | upper}}

## String Functions (from Sprig)
Quoted Name: {{.app.name | quote}}
Slugified: {{.app.name | replace " " "-" | lower}}
Padded: {{.app.name | printf ">>> %s <<<" }}
Repeated: {{repeat 3 "="}}

## Default Values
Debug Mode: {{default false .debug}}
Max Connections: {{default 100 .maxConnections}}
Log Level: {{default "info" .logLevel}}

## Date/Time Functions
Current Time: {{now | date "2006-01-02 15:04:05"}}
{{- if .deployTime}}
Deploy Time: {{.deployTime | date "Monday, January 2, 2006"}}
{{- end}}

## Conditional Logic
{{- if eq .environment "production"}}
🔴 PRODUCTION ENVIRONMENT
{{- else if eq .environment "staging"}}
🟡 STAGING ENVIRONMENT  
{{- else}}
🟢 DEVELOPMENT ENVIRONMENT
{{- end}}

## Lists and Loops
{{- if .services}}
Services:
{{- range .services}}
  - {{. | title}}
{{- end}}
{{- end}}

## Format Conversions

### JSON Output
```json
{{.app | toJson}}
```

### YAML Output
```yaml
{{.database | toYaml}}
```

### TOML Output
```toml
{{.config | toToml}}
```

## Advanced String Processing
{{- if .description}}
Description (trimmed): {{trim .description}}
Description (indented):
{{indent 4 .description}}
{{- end}}

## Math Functions
{{- if and .replicas .cpu}}
Total CPU: {{mul .replicas .cpu}}
{{- end}}

## Network/URL Functions  
{{- if .baseUrl}}
API Endpoint: {{.baseUrl}}/api/v1
Health Check: {{.baseUrl | printf "%s/health"}}
{{- end}}

## Base64 Encoding
{{- if .secret}}
Encoded Secret: {{.secret | b64enc}}
{{- end}}

---
Generated on {{now | date "2006-01-02 15:04:05 MST"}} with enhanced templater