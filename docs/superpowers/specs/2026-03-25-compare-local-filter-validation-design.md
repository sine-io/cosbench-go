# Compare Local Filter Validation Design

## Goal

Fail fast when a compare-local filter does not match any curated fixture.

## Recommended Approach

- validate `COMPARE_LOCAL_FILTER` against `testdata/workloads/compare-local-fixtures.txt` before running the compare-local loop
- if the filter is unknown, exit non-zero with a short error that includes the known fixture names
- add a regression test that proves unknown filters fail

## Success Criteria

1. unknown compare-local filters fail with a clear message
2. valid filtered runs still work
3. tests fail if the filter quietly produces an empty result set again
4. `go test ./...` and `go build ./...` remain green
