# PAAS Go Runner

`paas` is a single-binary Go CLI that executes a deterministic YAML DSL locally or over SSH.

**Documentation:** introductory guides live under [docs/](docs/README.md) (101-level sections in separate folders).

## Built-in extension sources

Runtime extension lookup order is:

1. `.paas/extensions/<name>.yml`
2. `~/.config/paas/extensions/<name>.yml`
3. embedded built-ins from `internal/extensions/*.yml`

The repo-root `extensions/*.yml` files are mirror copies for local visibility. The embedded files in `internal/extensions/` are the canonical source used by the binary.

## Built-in workflow goal

The built-in `deploy` extension mirrors the reference deployment flow from `.project/github-deploy.yml` in a local CLI form:

- validate required inputs
- prepare image tags locally
- log in to the registry
- build and push Docker images locally
- render compose with `APP_IMAGE`
- build the dashboard payload safely with `jq`
- update the dashboard app
- trigger dashboard deploy
- optionally run a smoke test

Project repositories can override that built-in `deploy` by placing `.paas/extensions/deploy.yml` next to the target application. That is how the `@twelve-factor` project switches from local image builds to a server-build-over-SSH flow.

## Real-world deployment notes

The first successful Windows + Git Bash + remote Linux deployment exposed a few important operational details:

- If a project override uses `git archive | ssh ...` for upload, that step depends on the local system `ssh`, not only on the Go SSH client embedded in `paas`.
- On Windows, remote steps can break if the full process environment is forwarded as-is. Variables such as `ProgramFiles(x86)` are not valid Bash export identifiers, and a Windows `PATH` can hide remote Linux binaries. Until filtering is implemented in the runner, use a small allowlisted shell environment when running remote deploys.
- The safest pattern for secrets is exported `INPUT_`* variables rather than `--input` flags in command history.
- Shell scripts copied into Linux images should still be normalized during build, even if the repository uses `.gitattributes` to force LF. A defensive `sed -i 's/\r$//'` in the image build saved a real deployment.

Example Windows wrapper command used successfully with a project-level remote deploy override:

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

## Known MVP limitations

- `timeout` is not implemented in the DSL yet.
- Native Windows Pageant support is not implemented yet; agent support currently depends on `SSH_AUTH_SOCK`.
- Only literals inside `when:` are parsed with Go-style escaping via `strconv.Unquote`.
- YAML-loaded values such as input defaults and `step.env` values are not unquoted again after YAML parsing.

