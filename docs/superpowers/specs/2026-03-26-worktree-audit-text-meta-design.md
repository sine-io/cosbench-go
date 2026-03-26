# Worktree Audit Text Metadata Design

## Goal

Make the human-readable audit helpers self-describing by adding the same header metadata already used by the text prune-plan helper.

## Scope

- Add `Generated at` to text audit output
- Add `Base ref` to text audit output
- Add `Current worktree` to text audit output
- Apply consistently to the base audit and filtered audit views
- Keep JSON behavior unchanged

## Design

The text output remains table-first, but gets three comment-style header lines before the table:

```text
# Generated at: 2026-03-26T12:34:56Z
# Base ref: origin/main
# Current worktree: /abs/path
PATH    BRANCH    CURRENT    STATE    DETAILS
...
```

This keeps the output copy/paste friendly and makes artifacts understandable without separate context.

## Testing

- Extend the text audit tests to require the three header lines
- Re-run `go test ./cmd/cosbench-go`
- Re-run `go test ./...`
- Re-run `go build ./...`
