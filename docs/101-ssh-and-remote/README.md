# 101 — SSH and remote execution

Remote steps are executed over **SSH** using `golang.org/x/crypto/ssh`.

## Connection

- **Address:** `host:port` from `ServerConfig` (default port 22).
- **User:** defaults to `root` if omitted.
- **Auth order (conceptual):**
  1. If `key` is set in server config: load that private key (passphrase via terminal if encrypted).
  2. Else try **SSH agent** signers.
  3. Else default keys under `~/.ssh` (e.g. `id_ed25712`, `id_rsa`).

## SSH agent

Agent support uses `SSH_AUTH_SOCK` and a **Unix domain socket** (`LoadAgentSigners`). The implementation returns signers and a **closable connection** so the socket stays open until dialing completes.

**Windows note:** Native Pageant or named-pipe agents are **not** implemented in this MVP. Agent support assumes `SSH_AUTH_SOCK` points to a compatible socket (often available in WSL or Git Bash environments that expose it).

For project overrides that use the local system `ssh` command inside a `local: true` step (for example `git archive | ssh host ...`), the Go SSH client is not enough by itself. You must also make sure the local shell can authenticate with ordinary `ssh`, typically by running:

```bash
eval "$(ssh-agent -s)"
ssh-add ~/.ssh/id_ed25712
ssh root@server
```

If the final `ssh root@server` succeeds without asking for the server password, the upload step should succeed too.

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

## Windows environment caution

In the current MVP, remote steps inherit process environment values unless the caller trims them first. On Windows this can cause two real problems:

- variables such as `ProgramFiles(x86)` are invalid Bash export names
- a Windows `PATH` can hide `/usr/bin/bash` or other Linux binaries on the remote host

For remote deploys from Windows + Git Bash, the proven workaround is to start `paas` with a small allowlisted environment:

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

## Related

- [Configuration](../101-configuration/README.md)  
- [Limitations](../101-limitations/README.md)  
