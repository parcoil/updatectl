# Configuration Schema

Detailed YAML configuration reference.

## Root Level

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `interval` | integer | Yes | Seconds between update checks |
| `intervalMinutes` | integer | No | **Deprecated**: Use `interval` instead |
| `projects` | array | Yes | List of projects to monitor |

## Project Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Unique project identifier |
| `path` | string | For git-based types | Local filesystem path |
| `repo` | string | For git-based types | Git repository URL |
| `type` | string | Yes | Project type: `docker`, `pm2`, `static`, `image` |
| `buildCommand` | string | No | Build command (for git-based types) |
| `image` | string | For image type | Docker image to pull (e.g., `ghcr.io/user/app:main`) |
| `port` | string | No | Port mapping for image type (e.g., `80:80`) |
| `env` | map[string]string | No | Environment variables for image type |
| `containerName` | string | No | Custom container name for image type (defaults to project name) |

## Validation Rules

- `interval`: Must be positive integer (seconds)
- `intervalMinutes`: **Deprecated**: Use `interval` instead
- `path`: Must exist and be writable (required for git-based types)
- `repo`: Must be valid Git URL (required for git-based types)
- `type`: Must be one of supported types: `docker`, `pm2`, `static`, `image`
- `buildCommand`: Optional for git-based types
- `image`: Required for `image` type, must be valid Docker image reference
- `port`: Optional for `image` type, must be valid port mapping format
- `env`: Optional for `image` type, key-value pairs
- `containerName`: Optional for `image` type

## Example

```yaml
interval: 1800  # 30 minutes
projects:
  - name: frontend
    path: /srv/frontend
    repo: https://github.com/company/frontend.git
    type: docker
    buildCommand: docker compose up -d --build
  - name: backend
    path: /srv/backend
    repo: https://github.com/company/backend.git
    type: pm2
  - name: webapp
    type: image
    image: ghcr.io/company/webapp:latest
    port: "3000:80"
    env:
      NODE_ENV: production
      API_URL: https://api.company.com
    containerName: company-webapp
```