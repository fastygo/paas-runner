# 101 — Security

## Threat model (brief)

`paas` is an operator tool: it runs arbitrary Bash from YAML and can connect to remote hosts. Trust extensions and configs like you would trust shell scripts and SSH keys.

## Secret masking

`output.SecretMasker` redacts substrings in printed output:

- Values registered explicitly via `AddSecret` (e.g. password-type inputs).
- **Heuristic** env keys: names containing `pass`, `password`, `token`, or `secret` (case-insensitive) via `AddFromEnv`.

**Note:** SSH key **paths** (`SERVER_KEY`) are not treated as secrets by the current heuristic so they remain visible for debugging.

Short secrets are still masked when added (empty string is skipped).

## Passphrases

Encrypted private keys may prompt for a passphrase via `golang.org/x/term` when stdin is a TTY.

## Non-interactive `run`

`paas run` does not prompt for missing inputs; required inputs must be supplied via `--input` or environment.

## Reporting

For vulnerability reporting, see the repository `SECURITY.md` if maintained for this project; otherwise follow your organization’s process.

## Related

- [CLI](../101-cli/README.md)  
- [Configuration](../101-configuration/README.md)  
