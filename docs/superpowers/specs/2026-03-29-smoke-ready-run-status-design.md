# Smoke Ready Run Status Design

## Goal

Extend `smoke-ready` and `smoke-ready-json` so they report the latest known GitHub workflow run status for the repository’s smoke workflows, not just whether those workflows exist.

## Problem

The repository now has a broad smoke surface:

- `Smoke Local`
- `Remote Smoke Local`
- `Remote Smoke Matrix`
- `Remote Smoke Recovery`
- `Remote Smoke Recovery Matrix`

The current `smoke-ready` helper already knows those workflows exist, but it cannot answer:

- when the latest run happened
- whether the latest run succeeded
- whether remote happy-path evidence is currently green
- whether remote recovery evidence is currently green

That leaves operators with presence information but not the latest signal.

## Scope

### In Scope

- query the latest run for each smoke workflow via `gh run list`
- include latest run metadata in JSON output
- include a concise latest-run view in text output
- add summary booleans for:
  - remote happy-path latest success
  - remote recovery latest success
- add focused tests using mocked workflow-run data

### Out Of Scope

- triggering workflows from `smoke-ready`
- aggregating artifacts
- changing existing smoke workflows

## Recommended Approach

Keep `smoke-ready` as a fast read-only helper and add one lightweight per-workflow latest-run lookup.

### Why latest-run status is the right next step

- presence alone is no longer enough now that the workflow surface is mature
- latest-run status gives immediately useful operational signal
- this remains much lighter than a full audit or artifact crawl

## Data Model

For each smoke workflow, report the latest available run record with fields such as:

- `status`
- `conclusion`
- `created_at`
- `url`

Add these under a new structure such as:

- `workflows.latest`

Summary should include:

- `remote_happy_latest_success`
- `remote_recovery_latest_success`

The existing availability-oriented fields should remain.

## Classification

Happy-path smoke workflows:

- `Remote Smoke Local`
- `Remote Smoke Matrix`

Recovery smoke workflows:

- `Remote Smoke Recovery`
- `Remote Smoke Recovery Matrix`

`Smoke Local` remains the local workflow readiness anchor, but its latest status should still be reported.

## Success Criteria

This slice is complete when:

1. `smoke-ready-json` includes latest run metadata per smoke workflow
2. `smoke-ready` text output shows a readable latest-run section
3. summary exposes latest-success booleans for remote happy and remote recovery
4. tests cover both JSON and text output with mocked workflow-run data

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
