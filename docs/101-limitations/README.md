# 101 — Limitations and known gaps

This section documents **behavioral gaps** relative to a full-featured orchestration tool. They are intentional or pending for MVP.

## DSL

| Topic | Status |
|-------|--------|
| `timeout` per step | **Not implemented** — no `timeout` field on `Step`; steps run until the process exits or context cancels. |
| `when:` shell fallback | **By design** — only the documented grammar is accepted; no shell evaluation. |

## String parsing

| Topic | Status |
|-------|--------|
| `when:` literals | Parsed with **`strconv.Unquote`** (Go string literal rules). |
| YAML `default:` / `step.env` | Interpreted by **YAML only**; no second `Unquote` pass in Go. |

## SSH / agent

| Topic | Status |
|-------|--------|
| `SSH_AUTH_SOCK` Unix agent | **Supported** (connection kept alive for signing). |
| Windows Pageant / named pipes | **Not implemented** — use environments that expose `SSH_AUTH_SOCK` or use explicit keys. |

## CLI

| Topic | Status |
|-------|--------|
| `paas init` overwrite | **Refuses** to overwrite existing `.paas/config.yml` (non-destructive). |

## Related

- [Overview](../101-overview/README.md)  
- [SSH and remote](../101-ssh-and-remote/README.md)  
