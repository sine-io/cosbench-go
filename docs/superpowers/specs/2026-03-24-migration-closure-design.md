# Migration Closure Design

## Goal

Close the current migration on a self-consistent, local-only v1 boundary:

- XML compatibility for the active COSBench/SineIO subset used by this repository
- S3 and SIO storage support
- single-process control plane + executor flow
- persistent job state, browser-visible history, and JSON/CSV exports

This design explicitly does **not** include the remote controller/driver split. That work remains a later seam, not part of migration closure for this phase.

## Scope Boundary

### In Scope

- align migration-facing documents to the actual intended boundary
- keep remote worker/controller split deferred
- implement missing local-only execution behavior required for a coherent v1 story
- expose work-level results in the same result model used by UI and exports
- strengthen preflight validation before execution starts
- add regression coverage for the newly closed gaps

### Out of Scope

- remote worker registration, heartbeats, mission dispatch, sample upload
- non-S3/SIO drivers
- legacy Freemarker UI parity
- full historical XML breadth beyond the active subset
- real-environment smoke tests against external S3/SIO systems

## Current Reality

The repository already contains more than the README and some migration docs admit:

- XML parsing and normalization are implemented
- S3/SIO config parsing and client wiring are implemented
- local execution, snapshots, web control plane, and result export are implemented
- `go test ./...` is green

The main remaining mismatch is not that the code is non-functional. It is that the migration boundary is inconsistent across docs, and a few local-only behavior gaps are still open inside that boundary.

## Problems To Close

### 1. Migration docs disagree about the target

`docs/migration-spec-v1.md` still describes a controller/driver HTTP/JSON split as part of v1. Other docs already treat remote workers as deferred. This creates false incompleteness and makes it impossible to say whether migration is done.

### 2. Some “supported” behavior is still partial

The execution layer advertises support for operations such as `mfilewrite`, `localwrite`, and `delay`, but at least part of that support is placeholder behavior. That is acceptable for scaffolding, but not for a closed migration boundary.

### 3. Result reporting is missing a work-level view

Job and stage summaries exist, and operation-level rollups exist. Work-level summaries are still absent even though they are one of the most useful local diagnostics for operator-facing inspection.

### 4. Validation happens too late

Some failures are only discovered after a job starts running. A local-only v1 should reject obviously invalid or non-runnable jobs before execution begins where practical.

## Proposed Approach

### A. Document Convergence

Update migration-facing docs to define the current completion target as:

- local single-process controller + executor
- XML active subset compatibility
- S3/SIO support
- web control plane, snapshots, and exports

Update the gap analysis and checklists so they describe only real remaining gaps inside that boundary. Remote split work should move fully into deferred/future sections.

### B. Local Execution Closure

Implement the missing behavior needed for the local boundary to be credible:

- `mfilewrite` performs a real multipart upload from a local file selected from configured file inputs
- `delay` performs real waiting behavior instead of a no-op
- `localwrite` remains local-only and should have explicit semantics documented rather than silently pretending to be a full remote-driver feature

The implementation should stay minimal and follow current code shape. This is not a refactor pass.

### C. Work-Level Result Model

Extend the result model so each stage can expose per-work summaries. The control plane already runs works one-by-one within a stage and receives a summary per work, so the missing piece is persistence and exposure:

- add work-level summary structures to the domain result model
- persist them through snapshots
- include them in JSON export
- add them to CSV export with a distinct scope
- surface them in the job detail page

### D. Preflight Validation

Add a control-plane preflight step before starting execution to catch issues such as:

- missing effective storage
- unsupported adapter configuration
- invalid operation config that can be detected without running the job

This should reduce “job started then immediately failed” cases.

## Architecture Notes

No new subsystem is required. The changes stay within existing seams:

- `docs/` for migration boundary updates
- `internal/domain/execution` for missing op semantics
- `internal/controlplane` for preflight and result assembly
- `internal/domain` for result DTO expansion
- `internal/web` for export/detail rendering
- targeted tests in current package-local test files

## Testing Strategy

Use TDD for each missing behavior:

- add failing tests for `mfilewrite`
- add failing tests for work-level result/export rendering
- add failing tests for preflight failures
- run targeted package tests red/green before the full suite

Final verification is:

- `go test ./...`
- `go build ./...`

## Success Criteria

Migration closure for this phase is reached when all of the following are true:

1. migration docs consistently define a local-only v1 boundary
2. the local boundary no longer claims placeholder operation support where real behavior is required
3. work-level summaries are persisted and visible through exports/UI
4. obvious non-runnable jobs fail during preflight rather than mid-run
5. full test suite and full build both pass

## Review Constraint

The usual spec-review subagent loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review and verification will be used instead.
