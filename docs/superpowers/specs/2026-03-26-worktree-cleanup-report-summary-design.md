# Worktree Cleanup Report Summary Design

## Goal

Expose the most useful worktree cleanup counts directly in the Markdown summary section of `make worktree-cleanup-report`.

## Scope

- Keep the existing audit, stale, and prune-plan sections intact
- Expand the Markdown summary with `Integrated`, `Stale`, and `Prune candidates`
- Reuse existing structured JSON helpers instead of recomputing counts ad hoc

## Design

`scripts/worktree_cleanup_report.py` already loads the audit JSON summary and prune-plan text. This change additionally reuses the structured prune-plan JSON payload so the Markdown summary can show:

- `Merged`
- `Integrated`
- `Stale`
- `Prune candidates`
- `Active`
- `Detached`
- `Unknown`

`Prune candidates` should come from the structured prune-plan summary, while the other counts continue to come from the audit summary.

## Testing

- Extend the Markdown cleanup report test to require the three missing summary lines
- Re-run `go test ./cmd/cosbench-go`
- Re-run `go test ./...`
- Re-run `go build ./...`
