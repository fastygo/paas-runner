# 101 — Getting started

## Prerequisites

- **Go** 1.22+ (to build from source).
- **Bash** on the machine where you run local steps: on Windows, install Git for Windows (includes Bash) or use WSL. `paas` does not execute steps with `cmd.exe` or PowerShell.
- For **remote** steps: OpenSSH-compatible server on the target, and local tools to manage keys (`~/.ssh` or ssh-agent).

## Build

From the repository root:

```bash
go build -o paas ./cmd/paas
```

Cross-compile example:

```bash
GOOS=linux GOARCH=amd64 go build -o paas ./cmd/paas
```

## Install (optional)

Place the binary on your `PATH` as `paas`, or invoke it with a full path.

## First-time project setup

1. Create a project config scaffold (fails if `.paas/config.yml` already exists — non-destructive):

   ```bash
   ./paas init
   ```

2. Optionally extract a built-in extension into the project tree:

   ```bash
   ./paas init -extract deploy
   ```

3. Configure servers for SSH in `~/.config/paas/servers.yml` (see [Configuration](../101-configuration/README.md)).

4. List available extensions:

   ```bash
   ./paas list
   ```

5. Validate an extension without running commands:

   ```bash
   ./paas validate deploy
   ```

6. Dry-run a workflow (prints resolved commands, no execution):

   ```bash
   ./paas run deploy --dry-run --input app_id=YOUR_APP_ID ...
   ```

   For this MVP, prefer either:

   ```bash
   ./paas run --input app_id=YOUR_APP_ID deploy
   ```

   or exported `INPUT_*` variables. Flags placed after the extension name can be unreliable because the CLI uses Go's standard `flag` parsing.

## Typical local workflow (built-in `deploy`)

The built-in `deploy` extension expects Docker, `docker compose`, `curl`, `jq`, and a git repository at the current directory. Pass required inputs (registry, dashboard URL, credentials, etc.) via `--input` or environment. See [Extensions — built-in deploy](../101-extensions/README.md).

## Typical remote-build workflow (project override)

A project can override the built-in `deploy` with `.paas/extensions/deploy.yml` and change the model completely. In the `@twelve-factor` project, the successful flow was:

1. run `git archive` locally
2. upload the repository snapshot over SSH
3. build and push the image on the remote Linux server
4. update the dashboard API
5. trigger deploy

On Windows + Git Bash, the proven wrapper command was:

```bash
env -i \
  HOME="$HOME" \
  USERPROFILE="$USERPROFILE" \
  HOMEDRIVE="$HOMEDRIVE" \
  HOMEPATH="$HOMEPATH" \
  PATH="$PATH" \
  TERM="${TERM:-xterm-256color}" \
  LANG="${LANG:-en_EN.UTF-8}" \
  SSH_AUTH_SOCK="$SSH_AUTH_SOCK" \
  SSH_AGENT_PID="$SSH_AGENT_PID" \
  INPUT_REGISTRY_USERNAME="$INPUT_REGISTRY_USERNAME" \
  INPUT_REGISTRY_PASSWORD="$INPUT_REGISTRY_PASSWORD" \
  ./paas.exe run deploy
```

## Where to go next

- [Configuration](../101-configuration/README.md) — `.paas/config.yml` and `servers.yml`  
- [CLI reference](../101-cli/README.md) — all commands and flags  
- [DSL reference](../101-dsl/README.md) — writing your own extension  
