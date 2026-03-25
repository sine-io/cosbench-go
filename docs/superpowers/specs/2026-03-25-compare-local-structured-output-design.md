# Compare Local Structured Output Design

## Goal

Make local comparison runs easier to reuse by producing stable per-fixture JSON result files instead of relying only on console text.

This should improve local review and GitHub Actions artifacts without changing which fixtures `compare-local` runs.

## Scope

### In Scope

- add a CLI option to write the summary JSON to a file
- update `make compare-local` to create a results directory with one JSON file per fixture
- update the manual `compare-local` workflow to upload the structured result directory
- document the new output layout

### Out of Scope

- changing the fixture set behind `compare-local`
- automatic updates to the legacy comparison matrix
- live endpoint automation

## Approaches Considered

### 1. Makefile-only redirection

Redirect `-json` stdout into files directly from `make compare-local`.

This is simple but keeps the CLI less reusable and makes it harder to combine human-readable console output with file-backed JSON output.

### 2. CLI `-summary-file` plus structured compare-local output

Add a dedicated CLI flag that writes the same summary payload to a file, then have `compare-local` populate a stable output directory.

This keeps the file-writing behavior in one place, is easy to test in Go, and lets local or CI callers reuse the same mechanism.

### 3. Separate helper tool

Build a dedicated compare-local helper command that orchestrates the runs and writes artifacts.

This would work but adds an unnecessary new entrypoint.

## Recommended Approach

Use approach 2.

Add `-summary-file <path>` to `cmd/cosbench-go`. When present, the CLI should write the same summary JSON payload that `-json` emits to the requested path. `make compare-local` should create a stable directory such as `.artifacts/compare-local/` and write one JSON file per fixture there. The manual workflow should upload that directory as the artifact instead of a single text file.

## Success Criteria

1. `go run ./cmd/cosbench-go ... -summary-file <path>` writes a valid summary JSON file
2. `make compare-local` creates stable per-fixture JSON outputs in a documented directory
3. the manual compare-local workflow uploads the structured result directory
4. `go test ./...` and `go build ./...` remain green
