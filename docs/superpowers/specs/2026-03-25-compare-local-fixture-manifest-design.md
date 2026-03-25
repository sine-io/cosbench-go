# Compare Local Fixture Manifest Design

## Goal

Reduce drift around `compare-local` by moving its curated fixture set into one source of truth.

This should make the Makefile target, tests, and docs easier to keep aligned without changing which fixtures run.

## Scope

### In Scope

- add a repository file that lists the curated `compare-local` fixtures
- update `make compare-local` to read from that manifest
- add test coverage that the manifest entries remain valid
- document the manifest location

### Out of Scope

- changing the fixture set itself
- changing compare-local result semantics
- adding live endpoint logic

## Approaches Considered

### 1. Keep the list duplicated

Leave the fixture list inline in the Makefile and continue updating tests and docs manually.

This is the least work today but keeps drift risk high.

### 2. Add a simple manifest file

Store the fixture list in a small text file that the Makefile and tests can read.

This keeps the solution lightweight, easy to inspect, and good enough for the current repository size.

### 3. Build a dedicated orchestration tool

Move compare-local fixture handling into a separate Go or shell tool.

This would work but adds unnecessary complexity.

## Recommended Approach

Use approach 2.

Add a plain-text manifest under `testdata/workloads/` that records the output name and workload path for each curated compare-local fixture. Update `make compare-local` to iterate over that file. Add a focused test that loads the manifest and verifies that the referenced workload files exist and still parse cleanly.

## Success Criteria

1. `make compare-local` reads its curated fixture set from one manifest file
2. automated tests fail if the manifest references a missing or unparsable workload
3. docs point contributors to the manifest when adjusting compare-local coverage
4. `go test ./...` and `go build ./...` remain green
