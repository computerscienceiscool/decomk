# decomk Design (Selftest + isconf-Aligned Selector Semantics)

This document captures the current design direction discussed for decomk selftest behavior after reviewing `doc/isconf-design.md`.

It focuses on:

- how decomk action arguments should be interpreted (isconf-aligned),
- how `examples/decomk-selftest` uses `decomk.conf` and `Makefile`,
- how to avoid selector translation tables while still testing tuple-selector behavior.

---

## 1) Design Goals

1. Keep decomk behavior aligned with isconf’s model:
   - context expansion produces tuples and tokens,
   - action args can map through tuple variables,
   - unknown action args fall back to literal targets.
2. Avoid maintaining selector translation tables in the harness.
3. Make selftest easy to reason about and deterministic.
4. Keep policy in `decomk.conf`, execution graph in `Makefile` + scripts.

---

## 2) Core decomk Action-Arg Model

For `decomk run [ARGS...]`:

1. Resolve context keys (typically `DEFAULT` plus discovered keys, or explicit `-context`).
2. Expand macros to tokens.
3. Partition tokens into:
   - tuples (`NAME=value`),
   - config target tokens.
4. Select final targets:
   - if `ARGS` are present:
     - if arg matches an effective tuple variable name, expand to that tuple’s value (split on whitespace),
     - otherwise treat arg as a literal make target.
   - if no `ARGS`:
     - use config target tokens if present,
     - else fallback to `INSTALL` tuple if present,
     - else make default goal.

This is the same shape documented in `doc/isconf-design.md`, adapted to decomk.

## 2.1) Environment contract (`env.sh` == make env)

- `decomk run` writes `${DECOMK_HOME}/env.sh` and uses the same canonical tuple list to invoke make.
- Canonical order is:
  1. incoming `DECOMK_*` passthrough vars,
  2. resolved config tuples,
  3. decomk-computed vars (last-wins).
- Tuple value sentinel `NAME=$` resolves from incoming env; if missing, decomk falls back to an earlier tuple assignment for `NAME`; if no fallback exists, run/plan fails fast.
- This keeps env exports and runtime make/recipe env behavior consistent even when make runs via sudo.

---

## 3) Selftest Design Decisions

## 3.1 Primary selector mode: literal targets

Primary selftest execution uses literal Makefile target names as decomk args:

- `all`
- `selftest-verify-tool-repo`
- `selftest-verify-conf-repo`
- `selftest-context-repo`
- `selftest-context-shared`
- `selftest-stamp-probe`
- `selftest-stamp-verify`

That means no primary selector mapping table in `decomk.conf`.

## 3.2 Tuple-selector coverage still required

Even with literal-primary mode, selftest explicitly exercises tuple selector expansion:

- tuple-only runs (arg maps to tuple value),
- mixed runs (tuple selector + literal target in same invocation),
- stamp/idempotency tuple paths.

## 3.3 No translation table policy

The harness does not keep a selector→scenario lookup for expected outcomes.
Instead, fixture make/scripts emit machine-readable pass/fail markers and the harness validates those markers from logs.

---

## 4) `decomk.conf` and `Makefile` Responsibilities

## 4.1 `decomk.conf` (policy/config layer)

`DEFAULT` contains runtime fixture tuples and tuple selectors. The workspace repo context (`decomk`) overrides selected tuples while leaving other `DEFAULT` tuples available.

Example pattern:

```conf
DEFAULT: TUPLE_VERIFY_TOOL='selftest-verify-tool-repo'
  TUPLE_VERIFY_CONF='selftest-verify-conf-repo'
  TUPLE_CONTEXT_OVERRIDE='selftest-context-default'
  TUPLE_DEFAULT_SHARED='selftest-context-shared'
  TUPLE_STAMP_PROBE='selftest-stamp-probe'
  TUPLE_STAMP_VERIFY='selftest-stamp-verify'

decomk: TUPLE_CONTEXT_OVERRIDE='selftest-context-repo'
```

## 4.2 `Makefile` (execution graph layer)

`Makefile` defines the executable target graph:

- `all` depends on fixture verification targets,
- verification and stamp targets emit `SELFTEST PASS ...` / `SELFTEST FAIL ...` markers,
- idempotent checks use stamp-backed file targets and explicit stamp verification.

No selector mapping logic is required in `Makefile`; target names are the selectors for literal-primary runs.

---

## 5) Harness Contract

## 5.1 `run.sh`

- Accepts decomk action args as positional passthrough.
- Rejects flag-style args (for example `-context`) so context selection stays automatic from workspace repo identity.
- Uses default tuple args when no args are provided:
  - `TUPLE_VERIFY_TOOL`
  - `TUPLE_VERIFY_CONF`
  - `TUPLE_CONTEXT_OVERRIDE`
  - `TUPLE_DEFAULT_SHARED`
- Renders those args into generated workspace `.devcontainer/devcontainer.json` as `DECOMK_RUN_ARGS`.
- Brings up one DevPod workspace, reads `make.log`, requires expected `SELFTEST PASS ...` markers, and fails on any `SELFTEST FAIL ...`.
- Runs an explicit two-step stamp check sequence:
  - `decomk run TUPLE_STAMP_PROBE`
  - `decomk run TUPLE_STAMP_PROBE TUPLE_STAMP_VERIFY`

## 5.2 `postCreateCommand.sh`

- Ensures a `decomk` binary is available in `PATH` (install-first default; optional clone mode).
- Syncs `DECOMK_CONF_REPO` into `${DECOMK_HOME}/conf`.
- Runs `decomk run ${DECOMK_RUN_ARGS:-all}`.
- Does not perform selftest marker validation; verification stays in fixture make/scripts plus harness log parsing.

This avoids hardcoded selector-expansion tables in the hook.

---

## 6) Test Matrix Requirements

Selftest covers the selector forms and stamp behavior exercised by the current harness:

1. **Literal-only**
   - `decomk run all`
2. **Tuple-only (default harness args)**
   - `decomk run TUPLE_VERIFY_TOOL TUPLE_VERIFY_CONF TUPLE_CONTEXT_OVERRIDE TUPLE_DEFAULT_SHARED`
3. **Stamp/idempotency sequence**
   - `decomk run TUPLE_STAMP_PROBE`
   - `decomk run TUPLE_STAMP_PROBE TUPLE_STAMP_VERIFY`

Success criteria:

- required `SELFTEST PASS ...` markers are present in each run's `make.log`,
- no `SELFTEST FAIL ...` marker appears,
- stamp probe runs once and stamp verify confirms idempotency.

---

## 7) Why this design

This keeps decomk selftest consistent with isconf’s architecture:

- tuples remain first-class for action-variable semantics,
- literal targets remain first-class for direct invocation,
- config stays declarative,
- execution stays in make targets/prereqs,
- no fragile translation table is needed.

---

## 8) Non-goals

- This design does not change decomk core parser grammar.
- This design does not require making action args expand arbitrary top-level macro keys by default.
- This design does not require separate containers per scenario; one container per harness invocation remains acceptable.
