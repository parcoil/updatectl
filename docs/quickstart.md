# Quick Start Guide

This guide will get you up and running with Updatectl in minutes.

## Prerequisites

- Go installed
- Git installed
- Docker or PM2 depending on your projects

## Installation

Download the latest release or build from source:

```bash
go build -o updatectl
sudo mv updatectl /usr/local/bin/
```

## First Setup

1. Initialize the configuration:

```bash
sudo updatectl init
```

2. Edit the config file at `/etc/updatectl/updatectl.yaml`

3. Add your first project:

```yaml
intervalMinutes: 15
projects:
  - name: myproject
    path: /path/to/project
    repo: https://github.com/user/project.git
    type: docker
    buildCommand: docker compose up -d --build
```

4. The daemon will start automatically and check for updates every 15 minutes.