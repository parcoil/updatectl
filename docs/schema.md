# Configuration Schema

Detailed YAML configuration reference.

## Root Level

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `intervalMinutes` | integer | Yes | Minutes between update checks |
| `projects` | array | Yes | List of projects to monitor |

## Project Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Unique project identifier |
| `path` | string | Yes | Local filesystem path |
| `repo` | string | Yes | Git repository URL |
| `type` | string | Yes | Project type: `docker`, `pm2`, `static` |
| `buildCommand` | string | No | Build command (Docker only) |

## Validation Rules

- `intervalMinutes`: Must be positive integer
- `path`: Must exist and be writable
- `repo`: Must be valid Git URL
- `type`: Must be one of supported types
- `buildCommand`: Required for `docker` type

## Example

```yaml
intervalMinutes: 30
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
```