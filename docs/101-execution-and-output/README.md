# 101 — Execution and output

## Runners

- **LocalRunner:** `exec` (or WSL) `bash -c` with merged environment and optional workdir.
- **RemoteRunner:** SSH session with streamed stdout/stderr; command wrapped as `bash -lc` (see [SSH](../101-ssh-and-remote/README.md)).

## Streaming

`Runner.Run` invokes a callback per **line** on stdout and stderr. The printer prints lines (masked); capture collects **stdout** lines only when `capture` is set.

## Step outcomes

| Exit code | `ignore_error` | Result |
|-----------|----------------|--------|
| 0 | any | Success |
| non-zero | `false` | Failure; pipeline stops |
| non-zero | `true` | Warning; pipeline continues |

**Dry-run:** Steps are not executed; `STEP_*` for capture may be set to empty for simulation consistency.

## Printer

- Step index, description, command line, streamed output, duration, status (success / warning / failure).
- Summary line counts successes, warnings, and failures.

## Related

- [DSL — capture and when](../101-dsl/README.md)  
- [Security — masking](../101-security/README.md)  
