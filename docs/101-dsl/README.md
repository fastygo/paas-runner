# 101 — DSL reference

Extensions are YAML documents describing **inputs** and **steps**. The root object is an **extension**; each step runs a Bash command string.

## Extension document

| Field | Required | Description |
|-------|----------|-------------|
| `id` | yes | Stable identifier (used in errors and UX). |
| `name` | yes | Short display name. |
| `description` | no | Longer description. |
| `inputs` | no | List of declared inputs. |
| `steps` | yes | Ordered list of steps. |

## Inputs

```yaml
inputs:
  - name: app_id
    label: Application ID
    type: text
    required: true
  - name: tag
    type: text
    default: ""
```

| Field | Description |
|-------|-------------|
| `name` | Logical name; normalized to an env key `INPUT_<NAME>` (see below). |
| `label` | Optional UI hint (CLI does not prompt interactively). |
| `type` | One of: `text`, `select`, `confirm`, `password`. |
| `default` | Optional; use YAML `default:` — pointer semantics in code distinguish “unset” from “empty string”. |
| `required` | If true, value must be non-empty after defaults and env merge. |
| `options` | Required for `type: select`. |

### Input name normalization

Input names are normalized with `dsl.NormalizeInputEnvKey`: hyphens become underscores, letters uppercased, invalid characters stripped. Example: `registry-host` → `INPUT_REGISTRY_HOST`.

CLI flags `--input registry-host=value` map to the same key.

## Steps

| Field | Description |
|-------|-------------|
| `id` | Optional; must be unique within the extension. |
| `run` | **Required.** Bash script or one-liner executed under Bash. |
| `description` | Optional; shown in output. |
| `local` | If `true`, run on the local runner; if `false` or omitted, run on the remote runner when SSH is configured. |
| `when` | Optional; deterministic condition (see below). Empty/absent means always run. |
| `capture` | Optional; name fragment for `STEP_<CAPTURE>` (uppercased). |
| `ignore_error` | If `true`, non-zero exit is treated as a **warning** and the pipeline continues. |
| `workdir` | Working directory for the step (local/remote runners). |
| `env` | Map of extra env vars for **this step only** (values support `${VAR}` substitution). |

There is **no** `timeout` field in the `Step` struct yet; see [Limitations](../101-limitations/README.md).

## Variable substitution: `${VAR}`

- Pattern: `${NAME}` where `NAME` matches `[A-Z_][A-Z0-9_]*`.
- If `NAME` is **missing** from the merged environment map, substitution **errors**.
- If `NAME` is present with an **empty** value, substitution yields an empty string (no error).

This applies to `run:` and values in `step.env`.

## `when:` grammar (deterministic)

No shell execution for conditions. Supported forms only:

| Form | Meaning |
|------|---------|
| *(empty or omitted)* | Always true |
| `VAR` | True if `isTruthy(env[VAR])` |
| `not VAR` | True if not `isTruthy(env[VAR])` |
| `VAR == "literal"` | String equality after literal parsing |
| `VAR != "literal"` | String inequality |

**Truthy** (`isTruthy`):

```text
value != "" && value != "0" && value != "false"
```

**Literals** in `==` / `!=` must be double-quoted. The literal token is parsed with **`strconv.Unquote`** (Go string literal semantics: `\n`, `\t`, `\"`, etc.).

**Important:** Only these `when:` literals use `strconv.Unquote`. Values loaded from YAML for `default:` or `step.env` follow **YAML** rules only; they are not unquoted again in Go.

## Capture: `STEP_*`

If `capture: foo` is set:

- After the step, `STEP_FOO` is set in the environment for later steps.
- Capture uses **stdout** lines only; the value is the **last non-empty** trimmed line.
- If there is no non-empty line, `STEP_FOO` is still set to **`""`** (so `${STEP_FOO}` never errors as “undefined” for that capture).

## Execution environment merge (`BuildBaseEnv`)

Order (later overrides earlier where applicable):

1. Project defaults from code (`DefaultConfig().Defaults`) merged with `.paas/config.yml` `defaults`
2. Server-derived env from selected `ServerConfig.ToEnv()`
3. Process environment (`os.Environ()`)
4. For each declared input: default from DSL, then process env for `INPUT_*`, then CLI `--input`
5. `stepCaptures` map (`STEP_*` from prior steps — applied in the executor across steps)

The executor merges `step.env` onto the current step’s environment before substituting `run:`.

## Validation

`dsl.ValidateExtension` checks: non-empty `id`, at least one step, valid `when`, valid capture names, valid `${VAR}` references in `run`/`env`, valid inputs and types.

## Related

- [Execution and output](../101-execution-and-output/README.md) — ignore_error, printer  
- [CLI reference](../101-cli/README.md) — `--input`  
