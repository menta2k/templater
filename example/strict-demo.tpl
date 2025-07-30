# Strict Mode Demo

## Defined Variables (will work in both modes)
App Name: {{.app.name}}
Version: {{.app.version}}
Environment: {{.environment}}

## Undefined Variable (will fail in strict mode)
Database Password: {{.database.password}}

## Nested Undefined (will fail in strict mode)
Cache Config: {{.cache.redis.host}}:{{.cache.redis.port}}

---
Generated at: {{now | date "2006-01-02 15:04:05"}}