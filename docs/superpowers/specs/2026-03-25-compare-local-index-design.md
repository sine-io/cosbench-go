# Compare Local Index Design

## Goal

Make compare-local artifacts easier to consume by adding one stable index file next to the per-fixture JSON summaries.

## Scope

### In Scope

- generate `.artifacts/compare-local/index.json`
- list the curated fixture name, workload path, and summary file for each compare-local run
- add regression coverage that the index is created
- document the index as the top-level artifact entrypoint

### Out of Scope

- changing fixture semantics
- changing summary payload contents
- live endpoint automation

## Recommended Approach

Keep the current manifest-driven compare-local flow and add one more generated artifact: `index.json`.

The index should be small and stable:

- `name`
- `workload`
- `summary`

This gives local users and the manual GitHub workflow a predictable first file to inspect without replacing the per-fixture summaries.

## Success Criteria

1. `make compare-local` writes `.artifacts/compare-local/index.json`
2. the index lists every fixture from `testdata/workloads/compare-local-fixtures.txt`
3. tests fail if the index is missing or incomplete
4. `go test ./...` and `go build ./...` remain green
