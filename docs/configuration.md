# Configuration

Updatectl uses a YAML configuration file to define update intervals and projects.

## Location

- Linux: `/etc/updatectl/updatectl.yaml`
- Windows: `%ProgramData%\updatectl\updatectl.yaml`

## Schema

```yaml
interval: 600  # Check interval in seconds (recommended)
intervalMinutes: 10  # Deprecated: Use interval instead
projects:
  - name: string      # Project identifier
    path: string      # Local filesystem path (required for git-based types)
    repo: string      # Git repository URL (required for git-based types)
    type: string      # Project type (docker/pm2/static/image)
    buildCommand: string  # Optional build command (runs after git pull for git-based types)
    image: string     # Docker image to pull (required for image type, e.g., "ghcr.io/user/app:main")
    port: string      # Port mapping (optional for image type, e.g., "80:80" or "3000:80")
    env:              # Environment variables (optional for image type)
      KEY: value
    containerName: string  # Optional custom container name (defaults to project name for image type)
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
    buildCommand: npm install && npm run build  # Optional: run before restart
```

### Static Site

```yaml
projects:
  - name: docs
    path: /var/www/docs
    repo: https://github.com/company/docs.git
    type: static
    buildCommand: npm run build  # Optional: run after git pull
```

### Image-based Project

For projects deployed as Docker images from registries like Docker Hub or GitHub Container Registry.

```yaml
projects:
  - name: vite-app
    type: image
    image: ghcr.io/user/vite-app:main
    port: "80:80"
    env:
      NODE_ENV: production
      API_URL: https://api.example.com
    containerName: my-vite-app  # Optional: defaults to project name
```