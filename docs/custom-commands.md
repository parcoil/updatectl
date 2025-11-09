# Custom Build Commands

Advanced configuration for all project types.

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

PM2 projects automatically restart the process after running the build command.

### Build and Restart

```yaml
buildCommand: npm install && npm run build
```

### Custom Deployment

```yaml
buildCommand: |
  npm ci
  npm run lint
  npm run test
  npm run build
```

## Static Site Commands

For static sites, build commands typically generate the site.

### Hugo

```yaml
buildCommand: hugo --minify
```

### Next.js Static Export

```yaml
buildCommand: npm run build && npm run export
```

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
