# Smoke Ready Upgrade Design

## Goal

Upgrade `smoke-ready` and `smoke-ready-json` so they reflect the repository’s current smoke surface rather than only checking the legacy `Smoke Local` workflow.

## Problem

The repository now has a materially broader smoke surface:

- local live smoke readiness via environment variables
- `Smoke Local`
- `Remote Smoke Local`
- `Remote Smoke Matrix`
- `Remote Smoke Recovery`
- `Remote Smoke Recovery Matrix`

But `scripts/smoke_ready.py` still only tracks:

- local live env presence
- one workflow: `Smoke Local`

That means the readiness helper is now under-reporting the actual verification surface and does not help operators quickly answer whether remote happy-path and recovery workflows are available.

## Scope

### In Scope

- expand workflow discovery in `smoke_ready.py`
- report presence for all current smoke workflows
- add richer summary fields for:
  - local live readiness
  - local workflow readiness
  - remote happy-path readiness
  - remote recovery readiness
- add automated tests for text and JSON modes
- update README wording for the expanded helper

### Out Of Scope

- running workflows from `smoke-ready`
- querying workflow run history
- helper changes for compare-local or other non-smoke workflows

## Recommended Approach

Keep `smoke-ready` as a lightweight availability helper, not an execution or audit tool.

### Why availability-only is the right scope

- the helper should stay fast
- workflow presence is stable enough to check locally
- run-status auditing belongs to separate tools or manual `gh run` inspection

## Output Model

Replace the single-workflow assumption with a named workflow map:

- `Smoke Local`
- `Remote Smoke Local`
- `Remote Smoke Matrix`
- `Remote Smoke Recovery`
- `Remote Smoke Recovery Matrix`

Summary fields should include:

- `local_env_ready`
- `local_workflow_ready`
- `remote_happy_ready`
- `remote_recovery_ready`
- `ready`

Where:

- `ready` should keep the original intent of “can I run at least one smoke path now?”
- the additional fields expose the broader remote surface explicitly

## Backward Compatibility

This slice should preserve:

- `make smoke-ready`
- `make smoke-ready-json`
- the current environment-variable contract

The JSON payload can grow new fields, but should not remove the existing top-level structure.

## Success Criteria

This slice is complete when:

1. `smoke-ready` lists the full smoke workflow surface
2. `smoke-ready-json` exposes the richer readiness summary
3. tests cover the expanded text and JSON output
4. repository tests and build remain green

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
