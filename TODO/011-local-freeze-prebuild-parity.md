# TODO 011 - local freeze pipeline for prebuild parity

## Decision Intent Log

ID: DI-011-20260417-135200
Date: 2026-04-17 13:52:00
Status: active
Decision: Adopt Option 2 (controlled local/CI freeze pipeline) as the active image-management path for decomk blocks, and explicitly track parity verification against Codespaces prebuild behavior.
Intent: Ensure frozen images include both Dockerfile and decomk `updateContent` effects while using a path we can operate today.
Constraints: Treat direct prebuild export/promotion as unavailable for now; keep lifecycle split (`updateContent` common/prebuild, `postCreate` runtime/user) intact; keep parity-proof model decision deferred to a dedicated bakeoff step.
Affects: `doc/image-management.md`, `TODO/011-local-freeze-prebuild-parity.md`, `TODO/010-codespaces-block-prebuild-profiles.md`, `TODO/007-devpod-gcp-selfhost-migration.md`, `TODO/TODO.md`.

ID: DI-011-20260418-115406
Date: 2026-04-18 11:54:06
Status: active
Decision: Implement TODO 011.3 as a first-class `decomk checkpoint` subcommand (not an external script), using existing decomk CLI structure (`main.go` switch + `flag.FlagSet`) and running prebuild/common lifecycle via `devcontainer up --prebuild` before checkpointing with Docker.
Intent: Make image checkpointing an operator-facing decomk workflow with one command path that is easy to automate and consistent with the rest of the CLI.
Constraints: Keep TODO 011.1/011.2 deferred; require caller-provided `--tag`; emit JSON to stdout; default to removing the container after checkpoint; no Cobra migration in this phase.
Affects: `cmd/decomk/main.go`, `cmd/decomk/checkpoint.go`, `cmd/decomk/checkpoint_test.go`, `README.md`, `TODO/011-local-freeze-prebuild-parity.md`.

## Goal

Implement a deterministic local/CI freeze process that produces
decomk-managed block images with measurable parity against Codespaces
prebuild behavior.

## Background

See `doc/image-management.md` for rationale, option analysis, and
design constraints.

## Scope

In scope:

- Define parity artifacts and comparison contract.
- Implement a freeze runner that executes Dockerfile + decomk
  `updateContent` phase.
- Add parity checks and acceptance gates.
- Document operator workflow for freeze/verify/promote.

Out of scope:

- Assuming availability of direct Codespaces prebuild export/promotion.
- Folding runtime/user customization (`postCreate`) into frozen shared
  block images.
- GCP rollout execution details (tracked in TODO 007).

## Dependencies and links

- Lifecycle evidence baseline: `TODO/009-phase-eval-lifecycle-spike.md`
- Block profile selection model: `TODO/010-codespaces-block-prebuild-profiles.md`
- Self-host migration track: `TODO/007-devpod-gcp-selfhost-migration.md`

## Subtasks

- [ ] 011.1 Define canonical parity artifact schema (image metadata, manifests, lifecycle markers, and provenance fields).
- [ ] 011.2 Run a parity proof-model bakeoff and lock the acceptance model for this TODO.
- [ ] 011.3 Implement `decomk checkpoint` (v1) to run prebuild/common lifecycle with `devcontainer up --prebuild`, then commit a checkpoint image.
- [ ] 011.3.1 Add `checkpoint` subcommand wiring in `cmd/decomk/main.go` and usage text updates.
- [ ] 011.3.2 Add `cmdCheckpoint(args, stdout, stderr)` in a new `cmd/decomk/checkpoint.go` using `flag.FlagSet`.
- [ ] 011.3.3 Require `--tag`; support `--workspace` (default `.`), optional devcontainer config override, and `--keep-container`.
- [ ] 011.3.4 Execute `devcontainer up --prebuild` and parse result metadata needed to identify the created container.
- [ ] 011.3.5 Commit container to caller tag via Docker and inspect the resulting image ID.
- [ ] 011.3.6 Emit a machine-readable JSON result on stdout (success/error, tag, container/image identifiers, cleanup status).
- [ ] 011.3.7 Default to container cleanup after checkpoint, with explicit retain behavior when `--keep-container` is set.
- [ ] 011.3.8 Add focused unit tests with command stubs for success, required-flag errors, prebuild failures, commit failures, and cleanup failures.
- [ ] 011.4 Enforce phase separation so frozen images exclude runtime/user-phase (`postCreate`) side effects.
- [ ] 011.5 Add deterministic pinning checks (base digest, package/tool versions, git refs, and other mutable inputs).
- [ ] 011.6 Implement parity comparator tooling and machine-readable failure reports.
- [ ] 011.7 Add repeatable operator/CI entrypoints and documentation for freeze + parity verification.
- [ ] 011.8 Execute acceptance matrix runs and record evidence paths/results in this TODO.

## 011.3 v1 command contract

- Command: `decomk checkpoint --tag <image:tag> [flags]`
- Runner model: prebuild/common lifecycle only (`devcontainer up --prebuild`) for this phase.
- Output: JSON on stdout (not files by default).
- Cleanup: remove checkpoint container by default; keep only with `--keep-container`.
- Deferred by design: parity schema/proof model (`011.1`, `011.2`) and comparator/gating (`011.6`).

## Acceptance criteria

- Repeated freezes with identical inputs produce stable outputs under the
  locked parity model.
- Parity checks against Codespaces baseline pass under the locked model.
- Drift is classified explicitly; unapproved drift fails the gate.
- Freeze workflow is documented and runnable without hidden/manual steps.
