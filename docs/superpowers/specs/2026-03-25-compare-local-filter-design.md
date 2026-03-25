# Compare Local Filter Design

## Goal

Allow local and manual compare-local runs to target one curated fixture instead of always running the full set.

## Recommended Approach

- add `COMPARE_LOCAL_FILTER` to the Makefile target
- treat an empty filter as “run all fixtures”
- when set, only run the matching fixture from `testdata/workloads/compare-local-fixtures.txt`
- update the manual workflow to expose a `fixture` input and pass it through

This keeps the existing default behavior intact while making targeted reruns cheaper.

## Success Criteria

1. `make compare-local COMPARE_LOCAL_FILTER=<name>` only writes results for the selected fixture
2. the manual workflow can run one fixture via `workflow_dispatch`
3. tests fail if the filter still writes extra entries
4. `go test ./...` and `go build ./...` remain green
