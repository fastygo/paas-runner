# 101 — CLI reference

Invocation: `paas <command> [options]`.

Global help:

```bash
paas
paas --help
paas help
```

## `paas run <extension>`

Runs the named extension (file name without `.yml` or with).

| Flag | Description |
|------|-------------|
| `--server <name>` | Use the named entry from `~/.config/paas/servers.yml`. If omitted, uses project `server:` from `.paas/config.yml` when remote steps exist. |
| `--input key=value` | Repeatable. Sets `INPUT_*` after normalization (see [DSL](../101-dsl/README.md)). |
| `--dry-run` | Resolves env and prints commands; does not execute steps. SSH may still be skipped or partially resolved depending on extension. |

**Local Bash:** If any step has `local: true` and execution is not dry-run, Bash must be discoverable (`findBash`); otherwise the command fails early with a clear error.

**Remote:** If any step is non-local, a server must be resolvable and SSH must succeed unless dry-run avoids connection (implementation may still load server env for masking when configured).

## `paas validate <extension>`

Static validation only: YAML parse, DSL validation, `when` syntax, variable reference shape. No execution.

## `paas list`

Lists extensions discoverable via project dir, user dir, and embedded names, with source hint (`project`, `user`, or `embedded`).

## `paas init`

Creates `.paas/` if needed and writes `.paas/config.yml` **only if that file does not exist**.

| Flag | Description |
|------|-------------|
| `-extract <name>` | Copies embedded `<name>.yml` into `.paas/extensions/`. |

## `paas servers`

Prints configured server names and connection summary **without** printing secrets.

## Environment variables

Inputs can be satisfied by process environment: keys `INPUT_<NAME>` match declared inputs after the same normalization as CLI.

## Related

- [Configuration](../101-configuration/README.md)  
- [Extensions](../101-extensions/README.md)  
