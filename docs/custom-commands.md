# Custom Build Commands

Advanced configuration for Docker and PM2 projects.

## Docker Commands

### Docker Compose

```yaml
buildCommand: docker compose up -d --build
```

### Direct Docker

```yaml
buildCommand: |
  docker build -t myapp . &&
  docker stop myapp || true &&
  docker rm myapp || true &&
  docker run -d --name myapp -p 3000:3000 myapp
```

### Multi-stage Builds

```yaml
buildCommand: |
  docker build -t myapp:latest -t myapp:$(git rev-parse --short HEAD) . &&
  docker-compose up -d
```

## PM2 Commands

PM2 projects automatically restart the process. For custom behavior:

Currently, PM2 projects use the default restart. For advanced PM2 usage, consider using Docker type with PM2 commands.

## Environment Variables

Pass environment variables to build commands:

```yaml
buildCommand: ENV_VAR=value docker compose up -d --build
```

## Pre/Post Commands

For complex workflows, use shell scripts:

```yaml
buildCommand: ./scripts/deploy.sh
```

Where `deploy.sh` contains:

```bash
#!/bin/bash
npm install
npm run build
docker compose up -d --build
```