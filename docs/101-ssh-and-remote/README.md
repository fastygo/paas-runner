# 101 — SSH and remote execution

Remote steps are executed over **SSH** using `golang.org/x/crypto/ssh`.

## Connection

- **Address:** `host:port` from `ServerConfig` (default port 22).
- **User:** defaults to `root` if omitted.
- **Auth order (conceptual):**
  1. If `key` is set in server config: load that private key (passphrase via terminal if encrypted).
  2. Else try **SSH agent** signers.
  3. Else default keys under `~/.ssh` (e.g. `id_ed25519`, `id_rsa`).

## SSH agent

Agent support uses `SSH_AUTH_SOCK` and a **Unix domain socket** (`LoadAgentSigners`). The implementation returns signers and a **closable connection** so the socket stays open until dialing completes.

**Windows note:** Native Pageant or named-pipe agents are **not** implemented in this MVP. Agent support assumes `SSH_AUTH_SOCK` points to a compatible socket (often available in WSL or Git Bash environments that expose it).

## Host key verification

`host_key_check` on the server entry:

| Value | Behavior |
|-------|----------|
| `strict` (default) | Uses `knownhosts` with `~/.ssh/known_hosts`. |
| `tofu` | Trust on first use style (append unknown keys). |
| `insecure` | Disables verification (warning printed). |

## Remote command wrapping

Remote commands are not executed as raw shell one-liners at the top level; the runner builds a shell that exports environment variables safely and runs:

```text
bash -lc '<command>'
```

Workdir is applied with `cd` before the command when set.

## Related

- [Configuration](../101-configuration/README.md)  
- [Limitations](../101-limitations/README.md)  
