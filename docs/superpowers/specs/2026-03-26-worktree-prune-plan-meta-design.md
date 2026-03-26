# Worktree Prune Plan Metadata Design

## Goal

Make the plain-text `worktree-prune-plan` output self-describing so it can be copied into issues or chat threads without losing generation context.

## Scope

- Add `Generated at` to text prune-plan output
- Add `Base ref` to text prune-plan output
- Add `Current worktree` to text prune-plan output
- Keep the command list and JSON output behavior unchanged

## Design

The text output remains comment-oriented:

```text
# Suggested cleanup commands
# Generated at: 2026-03-26T12:34:56Z
# Base ref: origin/main
# Current worktree: /abs/path
...
```

This keeps the output shell-safe for copy/paste while making the artifact understandable on its own.

## Testing

- Extend the prune-plan text test to require the three metadata lines
- Re-run `go test ./cmd/cosbench-go`
- Re-run `go test ./...`
- Re-run `go build ./...`
