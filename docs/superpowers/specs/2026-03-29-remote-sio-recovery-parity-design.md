# Remote SIO Recovery Parity Design

## Goal

Extend the existing remote recovery smoke path so the same recovery mechanics can be validated for `sio` as well as `s3`.

This slice is intended to close the remaining recovery-evidence gap between the two supported storage targets.

## Problem

The repository now has recovery evidence for:

- local MinIO-backed `s3 + recovery`
- a dedicated manual `Remote Smoke Recovery` workflow for `s3 + recovery`

But `sio + recovery` still lacks the same class of evidence. That leaves recovery verification asymmetric even though the project scope remains centered on **S3 and SIO**.

## Scope

### In Scope

- add one `sio` recovery workload fixture
- extend helper selection so `backend=sio`, `scenario=recovery` is valid
- add helper tests for the new combination
- document the local and manual invocation for `sio + recovery`
- validate the new path locally against MinIO

### Out Of Scope

- changing recovery orchestration logic
- adding a separate dedicated `Remote Smoke Recovery SIO` workflow
- changing the existing `Remote Smoke Recovery` fixed-`s3` workflow
- changing the default `CI` workflow

## Recommended Approach

Keep the existing recovery implementation backend-agnostic and only add the missing fixture plus selection wiring for `sio + recovery`.

### Why this is the right boundary

- the current recovery helper already proves lease expiry and reassignment mechanics
- recovery orchestration should not need backend-specific branching if the `sio` path is truly parity-complete
- the smallest meaningful proof is to reuse the same scenario with a `sio` fixture and local MinIO

## Fixture Strategy

Retain:

- `testdata/workloads/remote-smoke-s3-recovery-two-driver.xml`

Add:

- `testdata/workloads/remote-smoke-sio-recovery-two-driver.xml`

The new fixture should mirror the S3 recovery fixture closely:

- `storage type="sio"`
- one stage
- one work
- `workers="2"`
- `delay` operation with an explicit long duration

## Helper Behavior

The helper should continue to use:

- `SMOKE_REMOTE_LOCAL_BACKEND=s3|sio`
- `SMOKE_REMOTE_LOCAL_SCENARIO=single|multistage|recovery`

The only new supported combination is:

- `backend=sio`
- `scenario=recovery`

The summary contract should remain unchanged. No new fields are required for this parity slice.

## Validation Strategy

The minimum useful verification is:

1. helper tests proving `fixture_for_selection("sio", "recovery")`
2. local `s3 + recovery` rerun to prove no regression
3. local `sio + recovery` smoke against MinIO

After the local path is green, the existing parameterized `Remote Smoke Local` workflow can already be used for remote proof with:

- `gh workflow run "Remote Smoke Local" --repo sine-io/cosbench-go -f backend=sio -f scenario=recovery`

That means no workflow YAML changes are required in this slice.

## Success Criteria

This slice is complete when:

1. a `remote-smoke-sio-recovery-two-driver.xml` fixture exists
2. the helper accepts `backend=sio, scenario=recovery`
3. `SMOKE_REMOTE_LOCAL_BACKEND=sio SMOKE_REMOTE_LOCAL_SCENARIO=recovery make --no-print-directory smoke-remote-local` succeeds locally
4. the current `s3 + recovery` path still succeeds

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
