# Compare Local Output Pruning Design

## Goal

Prevent stale files from surviving across `make compare-local` runs.

Without this, removing or renaming a curated fixture can leave misleading old JSON files in `.artifacts/compare-local/`.

## Scope

### In Scope

- make `compare-local` clear its output directory before regenerating summaries
- add a regression test that proves stale files are removed
- document the refreshed-directory behavior

### Out of Scope

- changing the fixture set
- changing summary JSON content
- changing GitHub Actions trigger behavior

## Recommended Approach

Keep the output directory stable, but recreate it on each `compare-local` run before writing new files.

Add an integration test that:
- pre-seeds a stale file into a temp output directory
- runs `make compare-local COMPARE_LOCAL_OUTPUT_DIR=<tempdir>`
- verifies the stale file is gone and expected JSON files exist

## Success Criteria

1. `make compare-local` removes stale files from its output directory
2. a regression test proves stale files do not survive a fresh run
3. docs explain that compare-local refreshes the directory, not just appends to it
4. `go test ./...` and `go build ./...` remain green
