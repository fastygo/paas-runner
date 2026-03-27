# 101 — Configuration

`paas` uses two YAML configuration layers: **project** (per repository) and **user** (per account on the machine).

## Project config: `.paas/config.yml`

**Path:** `.paas/config.yml` (see `config.ProjectConfigPath()`).

**Purpose:** Defaults for the current project, optional default server name, and where to look for project-scoped extensions.

| Field | Type | Meaning |
|-------|------|---------|
| `server` | string | Name of a server entry in user config (`servers.yml`). Used when `paas run` needs a remote connection and `--server` is omitted. |
| `defaults` | map | Key/value strings merged into the execution environment **before** process env and inputs (see [DSL — variable merge](../101-dsl/README.md)). |
| `extensions_dir` | string | Directory for project extensions. Default: `.paas/extensions`. |

If the file is missing, built-in defaults apply (including `PAAS_VERSION` and default `extensions_dir`).

**`paas init`** creates `.paas/config.yml` only if it does **not** already exist (non-destructive).

## User config: `~/.config/paas/servers.yml`

**Path:** `~/.config/paas/servers.yml` on Unix; on Windows, under the user profile (`.config/paas/servers.yml`).

**Purpose:** Named SSH targets and optional dashboard-related fields exposed as environment variables.

Top-level shape:

```yaml
servers:
  production:
    host: deploy.example.com
    port: 22
    user: root
    key: ~/.ssh/id_ed25519
    dashboard_user: ""
    dashboard_pass: ""
    host_key_check: strict
```

### `ServerConfig` fields

| Field | Default | Meaning |
|-------|---------|---------|
| `host` | required for use | SSH hostname or IP |
| `port` | `22` | SSH port |
| `user` | `root` | SSH username |
| `key` | empty | Path to private key; if empty, agent then default keys may be used |
| `dashboard_user` / `dashboard_pass` | optional | Copied to `DASHBOARD_USER` / `DASHBOARD_PASS` in env |
| `host_key_check` | `strict` | Host key policy: `strict`, `tofu`, or `insecure` |

### Environment mapping

`ServerConfig.ToEnv()` exposes (when set): `SERVER_HOST`, `SERVER_PORT`, `SERVER_USER`, `SERVER_KEY`, `DASHBOARD_USER`, `DASHBOARD_PASS`, `HOST_KEY_CHECK`.

## User extensions directory

**Path:** `~/.config/paas/extensions/`

YAML files here are the **second** tier in extension resolution (after the project directory). See [Extensions](../101-extensions/README.md).

## Precedence summary

For **extension files**, lookup order is:

1. Project: `<extensions_dir>/<name>.yml` (default `.paas/extensions/<name>.yml`)
2. User: `~/.config/paas/extensions/<name>.yml`
3. Embedded built-ins in the binary

For **execution environment** (high level), merging is implemented in `dsl.BuildBaseEnv` — see [DSL reference](../101-dsl/README.md).

## Related

- [CLI — `servers`](../101-cli/README.md)  
- [SSH and remote](../101-ssh-and-remote/README.md)  
