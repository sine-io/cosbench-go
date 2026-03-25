# Compare Local Metrics Design

## Goal

Make compare-local artifacts and workflow summaries immediately useful by including the key per-fixture metrics in `index.json`.

## Recommended Approach

Stop hand-assembling `index.json` in the Makefile. Instead:

1. run the existing compare-local fixture loop and write each summary JSON file
2. run a short Python step that reads:
   - `testdata/workloads/compare-local-fixtures.txt`
   - `.artifacts/compare-local/*.json`
3. emit an enriched `index.json` with:
   - `name`
   - `workload`
   - `summary`
   - `stages`
   - `works`
   - `samples`
   - `errors`

Then update the manual workflow summary table to show those metrics directly.

## Success Criteria

1. `index.json` includes the key per-fixture metrics
2. the workflow summary table displays those metrics
3. tests fail if the enriched fields disappear
4. `go test ./...` and `go build ./...` remain green
