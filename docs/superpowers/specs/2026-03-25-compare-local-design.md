# Compare Local Design

## Goal

Add a repeatable local comparison entrypoint that runs a small, safe set of representative `mock`-backed workloads and emits their summaries in one place.

This slice is not about live endpoint validation. It is about making the existing comparison runbook easier to execute and repeat.

## Why This Slice

The repository already has:

- a legacy comparison matrix
- a runbook that explains how to run `cosbench-go` locally against representative fixtures
- multiple safe `mock`-backed representative fixtures

What it still lacks is one command that runs the local comparison set consistently, without asking contributors to copy multiple `go run` commands by hand.

## Scope

### In Scope

- a thin `make compare-local` target
- documentation updates pointing contributors at that target
- optional matrix updates if the command is used to refresh local evidence

### Out of Scope

- live endpoint comparison
- automatic editing of the comparison matrix from command output
- release automation
- benchmark orchestration across the legacy Java system

## Recommended Approach

Add a single `make compare-local` target that invokes the existing CLI against a curated list of safe fixtures using:

- `-backend mock`
- `-json`

The target should stay intentionally simple and human-oriented. It is a convenience wrapper, not a reporting subsystem.

## Fixture Set

Use fixtures that are:

- already within the current local-only v1 scope
- safe to run without live credentials
- representative of meaningful structure

Recommended initial set:

- `testdata/workloads/s3-active-subset.xml`
- `testdata/workloads/mock-stage-aware.xml`
- `testdata/workloads/mock-reusedata-subset.xml`
- `testdata/workloads/xml-splitrw-subset.xml`

## Output Shape

The target can print lightweight section headers around each fixture run.

The important property is that each underlying CLI invocation still emits the same JSON summary it would emit on its own.

This makes the command useful for:

- manual comparison refresh
- shell-level capture if needed

## Documentation Update

Update:

- `README.md`
- `AGENTS.md`
- `docs/legacy-comparison-matrix.md`
- `BOARD.md`
- `TODO.md`

to mention:

- the new `make compare-local` command
- that it is mock-backed and local-only
- that live comparison remains a separate workflow

## Success Criteria

This slice is complete when:

1. `make compare-local` exists
2. it runs the curated fixture set through the current CLI
3. documentation points contributors at the new command
4. `go test ./...` and `go build ./...` remain green

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
