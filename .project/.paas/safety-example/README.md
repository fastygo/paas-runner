# PAAS Safety Example

This folder contains safe example files for reusing `paas.exe` in a new project.

These examples follow one rule:

- keep long-lived machine access in the user profile
- keep project defaults in the project
- keep secrets out of committed files

## Recommended file split

### Global machine-level file

Use a user-scoped file for SSH connection settings:

- `C:/Users/<you>/.config/paas/servers.yml`

That file should contain:

- server host
- SSH port
- SSH user
- private key path
- host key verification mode

It should not contain:

- registry passwords
- dashboard passwords
- one-off project secrets

### Project-level file

Use a project-scoped file for stable non-secret defaults:

- `.paas/config.yml`

That file should contain:

- app name
- app id
- registry host
- image repository
- dashboard URL
- dashboard user
- optional healthcheck URL

It should not contain:

- registry password
- dashboard password
- API tokens

### Runtime-only values

Export secrets only in the shell session before deploy:

```bash
export INPUT_REGISTRY_USERNAME="..."
export INPUT_REGISTRY_PASSWORD="..."
export INPUT_DASHBOARD_PASS="..."
```

## Recommended deploy command

For the current Windows + Git Bash workflow:

```bash
env -i \
  HOME="$HOME" \
  USERPROFILE="$USERPROFILE" \
  HOMEDRIVE="$HOMEDRIVE" \
  HOMEPATH="$HOMEPATH" \
  PATH="$PATH" \
  TERM="${TERM:-xterm-256color}" \
  LANG="${LANG:-en_US.UTF-8}" \
  SSH_AUTH_SOCK="$SSH_AUTH_SOCK" \
  SSH_AGENT_PID="$SSH_AGENT_PID" \
  INPUT_REGISTRY_USERNAME="$INPUT_REGISTRY_USERNAME" \
  INPUT_REGISTRY_PASSWORD="$INPUT_REGISTRY_PASSWORD" \
  INPUT_DASHBOARD_PASS="$INPUT_DASHBOARD_PASS" \
  ./paas.exe run deploy
```

## Why this layout is safer

- the SSH private key stays outside the project
- the server inventory stays outside the project
- the project can be copied without dragging machine secrets with it
- passwords do not end up in git history or shell command history
- the `env -i` wrapper prevents Windows-only variables from breaking remote Bash
