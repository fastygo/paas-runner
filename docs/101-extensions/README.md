# 101 — Extensions

## Resolution order

When you run `paas run deploy` (or `validate`, `list`), the extension file is resolved in this order:

1. **Project:** `<extensions_dir>/<name>.yml`  
   Default `extensions_dir` is `.paas/extensions` (from `.paas/config.yml` or built-in default).
2. **User:** `~/.config/paas/extensions/<name>.yml`
3. **Embedded:** Files compiled into the binary from `internal/extensions/*.yml` via `go:embed`.

First match wins. Project overrides user; user overrides embedded.

## Repo-root `extensions/` folder

Files under the repository root `extensions/` are **mirrors** for documentation and diff review. The running binary does **not** read this path unless you copy files into `.paas/extensions` or user extensions. Canonical embedded sources live in `internal/extensions/`.

## Built-in extensions (typical)

| Name | Purpose |
|------|---------|
| `preflight` | Local workstation checks (bash, git, docker, compose, curl, jq, etc.). |
| `deploy` | Full local build/push/dashboard flow (see README in repo root). |
| `logs` | Fetch logs via dashboard HTTP API (inputs: app id, dashboard URL, credentials, line count). |
| `status` | Quick `docker ps` style status. |

Exact YAML lives in `internal/extensions/` and mirrors under `extensions/`.

## Authoring a custom extension

1. Add `myflow.yml` to `.paas/extensions/`.
2. Use `id`, `name`, `steps`, optional `inputs`.
3. Run `paas validate myflow` then `paas run myflow --dry-run`.

Keep commands portable Bash; avoid relying on interactive prompts.

## Related

- [DSL reference](../101-dsl/README.md)  
- [Getting started](../101-getting-started/README.md)  
