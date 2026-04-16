# phase-eval (POC spike)

`examples/phase-eval` is an empirical spike to measure lifecycle behavior across
DevPod and Codespaces.

## Important scope note

This is **not** a canonical decomk usage example.
It is a **POC spike** used to validate assumptions before changing decomk
stage-0 lifecycle design.

## What it evaluates

The spike captures evidence for:

- whether `updateContentCommand` runs,
- whether `postCreateCommand` runs,
- whether `GITHUB_USER` is populated,

across:

- `devpod build`,
- `devpod up`,
- Codespaces prebuild-list/API visibility,
- Codespaces prebuild workflow trigger/wait/log capture,
- Codespaces create/start behavior.

## Files

- `run.sh` — orchestrates scenarios and writes artifacts.
- `hook_probe.sh` — called by lifecycle hooks; emits marker lines and env snapshots.
- `lib.sh` — shared logging + command capture helpers.
- `.devcontainer/phase-eval/devcontainer.json` — evaluation-specific lifecycle config.

## Artifacts

By default, output is written under:

- `/tmp/decomk-phase-eval.<runid>/`

Key files:

- `summary.json`
- `raw/*.stdout.log`, `raw/*.stderr.log`, `raw/*.rc`
- `codespaces-prebuild.events.log` (when Codespaces prebuild logs contain hook markers)
- `scenario-notes.tsv`
- `diagnostics.complete`

## Usage

From repository root:

```bash
examples/phase-eval/run.sh --platform both
```

DevPod only:

```bash
examples/phase-eval/run.sh --platform devpod
```

Codespaces only:

```bash
examples/phase-eval/run.sh --platform codespaces
```

Keep remote workspace/codespace when failing:

```bash
examples/phase-eval/run.sh --keep-on-fail
```

## Prerequisites

- DevPod + Docker for DevPod scenarios.
- `gh` auth for Codespaces scenarios.
- Repository branch pushed for Codespaces create checks.
- Codespaces prebuilds configured for the repo (the harness triggers the prebuild workflow and waits for completion).
- `.devcontainer/phase-eval/devcontainer.json` committed on the evaluated branch (Codespaces resolves this path from remote repo contents).

## Interpretation

This spike intentionally favors explicit evidence over assumptions. If a platform
API is unavailable, the run records that condition in artifacts and exits
non-zero for the requested platform evaluation.
