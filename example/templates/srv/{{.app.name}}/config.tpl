[application]
name = "{{.app.name}}"
version = "{{.app.version}}"
environment = "{{default "development" .environment}}"

[server]
port = {{.app.port}}
debug = {{.config.debug}}

[database]
host = "{{.database.host}}"
port = {{.database.port}}
name = "{{.database.name}}"