# CLI Reference

Complete command reference for Updatectl.

## updatectl

Root command.

```bash
updatectl [command]
```

### Available Commands

- `init` - Initialize configuration and daemon
- `watch` - Run update daemon manually
- `version` - Show version information

## init

Initialize Updatectl configuration and set up the daemon.

```bash
updatectl init
```

Creates config file and systemd service (Linux) or Task Scheduler job (Windows).

## watch

Run the update daemon. Checks all projects for updates at configured intervals.

```bash
updatectl watch
```

Use for manual testing or when daemon is not running.

## version

Display version information.

```bash
updatectl version
```

## Global Flags

- `--help` - Show help
- `--version` - Show version