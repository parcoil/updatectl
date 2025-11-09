# Project Types

Updatectl supports different types of projects with varying update strategies.

## Docker

For containerized applications using Docker or Docker Compose.

**Process:**

1. Pull latest Git changes
2. Execute the `buildCommand`

**Example:**

```yaml
type: docker
buildCommand: docker compose up -d --build
```

**Use cases:** Web apps, APIs, databases in containers

## PM2

For Node.js applications managed by PM2 process manager.

**Process:**

1. Pull latest Git changes
2. Restart the PM2 process

**Example:**

```yaml
type: pm2
```

**Requirements:** PM2 must be installed and the app started with `pm2 start`

example: `pm2 start index.js --name my-app <br/>
name must match name in updatectl config

## Static

For static websites or projects that only need Git pulls.

**Process:**

1. Pull latest Git changes
2. No further action

**Example:**

```yaml
type: static
```

**Use cases:** Static site generators, documentation sites
