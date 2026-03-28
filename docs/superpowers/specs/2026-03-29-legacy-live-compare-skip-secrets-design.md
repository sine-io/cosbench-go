# Legacy Live Compare Skip-Secrets Design

## Goal

Make the manual `Legacy Live Compare` workflow behave like `make smoke-s3` when live endpoint secrets are absent: emit a clear skipped result instead of rendering empty storage credentials and failing at runtime.

## Current Problem

- `.github/workflows/legacy-live-compare.yml` always renders the selected legacy XML fixture with `COSBENCH_SMOKE_*` values from repository secrets.
- When those secrets are unset, the renderer writes empty `accesskey`, `secretkey`, and `endpoint` values into the generated XML.
- The later `go run ./cmd/cosbench-go ...` step fails with storage preflight errors such as `sio endpoint is required`.
- The failure is operationally misleading because the real issue is missing workflow configuration, not workload incompatibility.

## Desired Behavior

If any required live credential is missing:

- the workflow should not execute the rendered workload
- the workflow should write a short skipped summary to both the job summary and the uploaded artifact
- the workflow should complete successfully rather than failing

If the required credentials are present:

- the workflow should keep the existing render-and-run behavior
- the artifact layout should remain stable under `.artifacts/legacy-live-compare/`

## Scope

In scope:

- `Legacy Live Compare` workflow control flow
- workflow contract tests
- runbook/readme notes that describe the new skip semantics

Out of scope:

- changes to `cmd/cosbench-go`
- changes to XML parsing or runtime storage validation
- new backend support
- automatic live comparison against the legacy Java runtime

## Design

Add a shell preflight step near the start of `.github/workflows/legacy-live-compare.yml` that checks:

- `COSBENCH_SMOKE_ENDPOINT`
- `COSBENCH_SMOKE_ACCESS_KEY`
- `COSBENCH_SMOKE_SECRET_KEY`

If any are missing, the step should:

- create `.artifacts/legacy-live-compare/summary.json` with a `skipped` status and reason
- create `.artifacts/legacy-live-compare/run.log` with the same reason
- expose a workflow output such as `should_run=false`

The render and execute steps should run only when `should_run == 'true'`.

The always-on summary and artifact upload steps should remain, but the summary block should include the skip status when present.

## Acceptance Criteria

- Triggering `Legacy Live Compare` without configured repository secrets results in a successful workflow run with explicit `skipped` evidence.
- Triggering it with configured secrets preserves the current render-and-run path.
- The artifact still uploads on both skip and run paths.
- README and legacy live run guidance mention that the workflow may skip cleanly when secrets are unavailable.
