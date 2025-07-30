Application Name: {{.app.name}}
Version: {{.app.version}}
Port: {{.app.port}}

# Database Configuration
Database Host: {{.database.host}}
Database Port: {{.database.port}}
Database Name: {{.database.name}}
Database User: {{.database.user}}
{{if .database.password}}Database Password: {{.database.password}}{{end}}

# Deployment Configuration
Replicas: {{.replicas}}

# Application Settings
Debug Mode: {{.config.debug}}
Log Level: {{.config.logLevel}}
Max Connections: {{.config.maxConnections}}

# Environment variables (converted to camelCase)
{{if .home}}Home Directory: {{.home}}{{end}}
{{if .user}}Current User: {{.user}}{{end}}
{{if .databasePassword}}Database Password from ENV: {{.databasePassword}}{{end}}
{{if .appVersion}}App Version from ENV: {{.appVersion}}{{end}}
{{if .maxConnections}}Max Connections from ENV: {{.maxConnections}}{{end}}

# Using template functions
Default Environment: {{default "development" .environment}}

# Command-line set values (highest precedence)
{{if .override}}Override Value: {{.override}}{{end}}
{{if .features}}{{if .features.experimental}}Experimental Features: {{.features.experimental}}{{end}}{{end}}