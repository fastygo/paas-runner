# PAAS documentation

Introductory (**101**) guides live in separate folders below. Each folder is self-contained; read them in any order, or follow the suggested path.

## Suggested reading order

1. [Overview](101-overview/README.md) — what `paas` is and design goals  
2. [Getting started](101-getting-started/README.md) — build, install, first commands  
3. [Configuration](101-configuration/README.md) — project and user config, precedence  
4. [DSL reference](101-dsl/README.md) — extensions, steps, variables, `when`, capture  
5. [CLI reference](101-cli/README.md) — `run`, `validate`, `list`, `init`, `servers`  
6. [Extensions](101-extensions/README.md) — lookup order, built-ins, authoring  
7. [SSH and remote execution](101-ssh-and-remote/README.md) — keys, agent, host keys  
8. [Execution and output](101-execution-and-output/README.md) — runners, streaming, dry-run  
9. [Security](101-security/README.md) — masking, secrets, operational hygiene  
10. [Limitations](101-limitations/README.md) — MVP gaps and future work  

The docs now also reflect a real Windows + Git Bash deployment that:

- uploaded source with `git archive | ssh`
- built the image on the remote Linux server
- pushed it to a private registry
- updated the dashboard API
- exposed the site successfully over HTTPS

## Quick links

| Topic | Folder |
|-------|--------|
| Product overview | [101-overview](101-overview/README.md) |
| Install and first run | [101-getting-started](101-getting-started/README.md) |
| `.paas/config.yml` and `servers.yml` | [101-configuration](101-configuration/README.md) |
| YAML DSL grammar | [101-dsl](101-dsl/README.md) |
| Command-line flags | [101-cli](101-cli/README.md) |
| Built-in and custom extensions | [101-extensions](101-extensions/README.md) |
| SSH client behavior | [101-ssh-and-remote](101-ssh-and-remote/README.md) |
| How steps run and print | [101-execution-and-output](101-execution-and-output/README.md) |
| Redaction and credentials | [101-security](101-security/README.md) |
| Known gaps | [101-limitations](101-limitations/README.md) |
