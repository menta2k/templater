apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.app.name}}
  namespace: {{default "default" .environment}}
spec:
  replicas: {{.replicas}}
  selector:
    matchLabels:
      app: {{.app.name}}
  template:
    metadata:
      labels:
        app: {{.app.name}}
        environment: {{default "development" .environment}}
    spec:
      containers:
      - name: {{.app.name}}
        image: {{.app.name}}:{{.app.version}}
        ports:
        - containerPort: {{.app.port}}
        env:
        - name: ENVIRONMENT
          value: "{{default "development" .environment}}"
        - name: DATABASE_HOST
          value: "{{.database.host}}"