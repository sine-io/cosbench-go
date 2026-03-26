# Worktree Cleanup Report Metadata Design

## Goal

Make the human-readable cleanup report self-describing by showing when it was generated and which current worktree was excluded from prune-candidate calculations.

## Scope

- Add `Generated at` to the Markdown cleanup report summary
- Add `Current worktree` to the Markdown cleanup report summary
- Keep the existing report sections and command outputs unchanged
- Keep the JSON payload machine-readable and add top-level `generated_at` for consistency

## Design

`scripts/worktree_cleanup_report.py` already has everything it needs:

- the current worktree path is available via `prune_plan["summary"]["current_worktree"]`
- the report can generate its own UTC timestamp once per run

The Summary section will now include:

- `Generated at`
- `Base ref`
- `Current worktree`
- existing count lines

## Testing

- Extend the Markdown cleanup report test to require the two new summary lines in stdout and the written report file
- Re-run `go test ./cmd/cosbench-go`
- Re-run `go test ./...`
- Re-run `go build ./...`
