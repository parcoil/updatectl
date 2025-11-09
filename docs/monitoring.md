# Monitoring

How to monitor Updatectl's activity and performance.

## Logs

### Systemd Logs

```bash
journalctl -u updatectl -f
```

### Manual Runs

When running `updatectl watch` manually, output goes to stdout.

## Metrics

### Update Frequency

Monitor how often projects are updated:

```bash
journalctl -u updatectl | grep "Checking" | wc -l
```

### Success Rate

Check for failures:

```bash
journalctl -u updatectl | grep "failed\|Failed" | tail
```

## Health Checks

### Service Status

```bash
systemctl status updatectl
```

### Configuration Validation

```bash
python3 -c "import yaml; yaml.safe_load(open('/etc/updatectl/updatectl.yaml'))"
```

## Performance

Monitor resource usage:

```bash
ps aux | grep updatectl
```

Adjust `intervalMinutes` based on system load.
