# Configuration

Updatectl uses a YAML configuration file to define update intervals and projects.

## Location

- Linux: `/etc/updatectl/updatectl.yaml`
- Windows: `%ProgramData%\updatectl\updatectl.yaml`

## Schema

```yaml
intervalMinutes: 10  # Check interval in minutes
projects:
  - name: string      # Project identifier
    path: string      # Local filesystem path
    repo: string      # Git repository URL
    type: string      # Project type (docker/pm2/static)
    buildCommand: string  # Build command (for docker type)
```

## Examples

### Docker Project

```yaml
projects:
  - name: webapp
    path: /srv/webapp
    repo: https://github.com/company/webapp.git
    type: docker
    buildCommand: docker compose up -d --build
```

### PM2 Project

```yaml
projects:
  - name: api
    path: /srv/api
    repo: https://github.com/company/api.git
    type: pm2
```

### Static Site

```yaml
projects:
  - name: docs
    path: /var/www/docs
    repo: https://github.com/company/docs.git
    type: static
```