# updatectl

A CLI tool for automating project updates. It periodically pulls the latest changes from Git repositories and rebuilds/restarts projects based on their type (PM2 or Docker).

> **Warning:** This project is a work in progress and not ready for production use.

## Installation

### Using Installers

Installer scripts are provided in the repository for easy installation.

- **Linux**: Run `./install.sh` (requires sudo for installation)
- **Windows**: Run `install.bat` as administrator

### Manual Installation

1. Clone or download the repository.
2. Build the executable: `go build -o updatectl main.go`
3. Move `updatectl` to a directory in your PATH (e.g., `/usr/local/bin/` on Linux or `C:\Program Files\updatectl\` on Windows).

or build it and upload to your server

## Usage

### Initialize Configuration

Run `updatectl init` to create the default configuration file and set up the daemon.

- On Linux: Creates a systemd service.
- On Windows: Creates a Task Scheduler job.

### Start Watching

The daemon runs automatically after init. To run manually: `updatectl watch`

## Configuration

Configuration is stored in:

- Linux: `/etc/updatectl/updatectl.yaml`
- Windows: `%ProgramData%\updatectl\updatectl.yaml`

Example config:

```yaml
intervalMinutes: 10
projects:
  - name: myapp
    path: /srv/myapp
    repo: https://github.com/user/myapp.git
    type: docker
    buildCommand: docker compose up -d --build
  - name: webserver
    path: /srv/webserver
    repo: https://github.com/user/webserver.git
    type: pm2
    buildCommand: "" # Not used for PM2
```

- `intervalMinutes`: How often to check for updates (in minutes).
- `projects`: List of projects to monitor.
  - `name`: Project name.
  - `path`: Local path to the project.
  - `repo`: Git repository URL.
  - `type`: "pm2" or "docker".
  - `buildCommand`: Command to run after pulling (for Docker).

## Supported Project Types

- **PM2**: Restarts the PM2 process with the project name.
- **Docker**: Runs the specified build command (e.g., Docker Compose rebuild or direct Docker commands).

### Docker Without Compose

For projects using Docker without Compose, set the `buildCommand` to build and run the container directly. Example:

```yaml
projects:
  - name: myapp
    path: /srv/myapp
    repo: https://github.com/user/myapp.git
    type: docker
    buildCommand: docker build -t myapp . && docker stop myapp || true && docker rm myapp || true && docker run -d --name myapp -p 8080:8080 myapp
```

## Requirements

- Go
- Git
- PM2 (for PM2 projects)
- Docker (for Docker projects)
