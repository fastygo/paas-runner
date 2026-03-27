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

## Known MVP limitations

- `timeout` is not implemented in the DSL yet.
- Native Windows Pageant support is not implemented yet; agent support currently depends on `SSH_AUTH_SOCK`.
- Only literals inside `when:` are parsed with Go-style escaping via `strconv.Unquote`.
- YAML-loaded values such as input defaults and `step.env` values are not unquoted again after YAML parsing.

