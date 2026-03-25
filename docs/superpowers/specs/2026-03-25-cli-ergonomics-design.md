# CLI Ergonomics Design

## Goal

Improve the local `cosbench-go` CLI workflow without changing its core role as a simple local workload runner.

This slice should make the current CLI easier to discover and script while remaining backward-compatible.

## Problem

Current CLI friction points:

- workload input only works through `-workload`
- there is no short alias like `-f`
- there is no positional path fallback
- `-json` output is not clean machine-readable JSON because progress lines are printed before the JSON summary
- usage/help is minimal

These are usability issues, not missing product features.

## Scope

### In Scope

- `-f` alias for workload path
- positional workload path fallback when flags are omitted
- pure JSON output when `-json` is requested
- clearer usage/help text
- tests covering the new invocation shapes
- small documentation updates

### Out of Scope

- subcommand architecture
- new benchmark/reporting features
- config files for the CLI itself
- interactive prompts

## Recommended Approach

Keep the existing single-command CLI and make it friendlier.

The CLI should accept workload input in this priority order:

1. `-workload`
2. `-f`
3. first positional argument

If none are provided, print usage and exit as before.

## JSON Output Rule

When `-json` is set:

- write only the JSON summary to stdout
- send progress/logging lines to stderr or suppress them entirely

The important property is that stdout becomes machine-readable JSON with no leading text.

This makes shell scripting and CI integration predictable.

## Backward Compatibility

Existing usage must continue to work:

- `go run ./cmd/cosbench-go -workload <file> -backend mock -json`

The new forms should be additive:

- `go run ./cmd/cosbench-go -f <file> -backend mock`
- `go run ./cmd/cosbench-go <file> -backend mock`

## Tests To Add

Add CLI tests proving:

- `-workload` still works
- `-f` works
- positional workload path works
- `-json` output is pure JSON
- missing workload path still returns an error/usage path

## Documentation Update

Update repository-facing docs so examples show the friendlier forms and mention that `-json` is now safe for machine parsing.

Good targets:

- `README.md`
- `AGENTS.md`
- `BOARD.md`
- `TODO.md`

## Success Criteria

This slice is complete when:

1. workload path can be supplied via `-workload`, `-f`, or positional argument
2. `-json` emits machine-readable stdout with no progress noise
3. the old `-workload` form still works
4. CLI tests cover the new invocation shapes
5. `go test ./...` and `go build ./...` remain green

## Review Constraint

The usual delegated spec-review loop is not being used here because delegation is restricted unless explicitly requested by the user. Manual review will be used instead.
