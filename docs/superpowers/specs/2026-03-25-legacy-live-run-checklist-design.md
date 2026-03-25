# Legacy Live Run Checklist Design

## Goal

Add a dedicated, actionable checklist for future live comparisons against `cosbench-sineio`.

The checklist should make the remaining live-comparison work process-ready even when the environment is not yet available.

## Why This Slice

The repository already has:

- a legacy comparison matrix
- local comparison evidence
- code-level storage driver comparison notes
- an opt-in smoke workflow

What is still missing is a single, focused document that answers:

> "When live credentials and an endpoint become available, exactly what do I check, in what order, and what do I write back into the matrix?"

## Scope

### In Scope

- one dedicated live-run checklist document
- links from the matrix and README
- status wording updates in board/todo

### Out of Scope

- running live comparisons in this slice
- adding new automation for live credentials
- changing application code

## Recommended Approach

Keep responsibilities split:

- `docs/legacy-comparison-matrix.md` remains the system of record for findings
- the new checklist becomes the execution guide for future live runs

This keeps the matrix readable and makes the live process explicit.

## Checklist Content

The checklist should include:

### 1. Preconditions

- required `COSBENCH_SMOKE_*` environment variables
- endpoint reachability assumption
- a reminder not to claim live results without actual environment access

### 2. Smoke Precheck

- run `make smoke-s3`
- record whether connectivity passes before workload-level comparison starts

### 3. Recommended Execution Order

1. `testdata/legacy/sio-config-sample.xml`
2. `testdata/legacy/s3-config-sample.xml`
3. storage-level `part_size` / `restore_days`
4. cleanup/list-sensitive scenarios

### 4. Recording Rules

For each run, update the matrix with:

- execution outcome category
- result surface availability
- notable semantic differences
- `match` / `acceptable delta` / `mismatch`

### 5. Known Watchpoints

Reference the current driver comparison notes:

- SIO path-style default
- delete tolerance
- list output shape
- storage-level `part_size`
- storage-level `restore_days`
- slash-containing SIO bucket names

## Documentation Update

Update:

- `docs/legacy-comparison-matrix.md`
- `README.md`
- `BOARD.md`
- `TODO.md`

so the repository clearly communicates:

- a live comparison process exists
- the remaining blocker is environment availability

## Success Criteria

This slice is complete when:

1. `docs/legacy-live-run-checklist.md` exists
2. the matrix links to it
3. the README points contributors at it
4. board/todo reflect that live comparison is process-ready but environment-blocked
5. `go test ./...` and `go build ./...` remain green

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
