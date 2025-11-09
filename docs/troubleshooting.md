# Troubleshooting

Common issues and their solutions.

## Daemon Not Starting

**Symptoms:** `systemctl status updatectl` shows failed

**Solutions:**

- Check systemd service file: `cat /etc/systemd/system/updatectl.service`
- Verify user has permissions for project paths
- Check logs: `journalctl -u updatectl`
- Try running `updatectl init` again

## Git Pull Failures

**Symptoms:** "Git pull failed" in logs

**Solutions:**

- Ensure SSH keys are set up for private repos
- Check repository permissions
- Verify the path exists and is a Git repository

## Build Command Failures

**Symptoms:** Docker/PM2 commands fail

**Solutions:**

- Test commands manually in the project directory
- Check for missing dependencies (Docker, PM2)
- Verify environment variables are available

## Permission Issues

**Symptoms:** Access denied errors

**Solutions:**

- Run updatectl as appropriate user (not root if possible)
- Ensure project directories are writable by the service user
- Check file ownership: `ls -la /path/to/project`
