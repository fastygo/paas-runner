# 101 — Overview

## What is `paas`?

`paas` is a single static binary written in Go. It reads **YAML extensions** (a small domain-specific language) and executes each **step** either:

- **locally**, under Bash on the operator machine, or  
- **remotely**, over SSH on a Linux host, with commands wrapped as `bash -lc '…'`.

There is no multi-tenant scheduler, no cluster orchestrator, and no embedded business logic in Go for deployments. The DSL is **shell-oriented**: every executable step is a Bash snippet in `run:`.

## Design goals

- **Determinism** — Conditions use a fixed `when:` grammar with no shell fallback. The same extension and inputs should behave the same in validate and dry-run as far as parsing and substitution allow.
- **Bash everywhere** — Local execution uses Bash (including discovery of Git Bash or WSL on Windows). Remote execution always uses Bash on the server.
- **Non-interactive `run`** — Required inputs must be supplied via flags or environment; there are no interactive prompts in the default `paas run`.
- **Predictable variables** — `${VAR}` substitution fails if `VAR` is missing from the execution environment; empty string is allowed and substitutes to nothing.
- **Safe output** — Secret-like values are redacted in printed output where configured.

## Operating model

- One operator, one machine (or one SSH target per run).
- Typical use: deploy or operate a single application stack as **root** or a user with Docker access.
- Extensions are **one file per workflow** (e.g. `deploy.yml`, `preflight.yml`).

## Repository layout (high level)

| Area | Role |
|------|------|
| `cmd/paas/` | Program entrypoint |
| `internal/cli/` | Commands and flags |
| `internal/config/` | Project and user YAML config |
| `internal/dsl/` | Parse extensions, `when`, `${…}` substitution |
| `internal/runner/` | Local/remote execution and step loop |
| `internal/sshclient/` | SSH dial, keys, agent, host keys |
| `internal/output/` | Printer and secret masker |
| `internal/extensions/` | **Embedded** built-in YAML (compiled into the binary) |
| `extensions/` (repo root) | Optional **mirror** of built-ins for browsing; runtime uses embed + project/user paths |

## Related documentation

- Next: [Getting started](../101-getting-started/README.md)  
- DSL details: [DSL reference](../101-dsl/README.md)  
- Limits: [Limitations](../101-limitations/README.md)  
