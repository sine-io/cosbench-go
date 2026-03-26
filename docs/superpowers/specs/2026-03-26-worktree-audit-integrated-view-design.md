# Worktree Audit Integrated View Design

## Goal

Expose integrated worktrees as a first-class view instead of making users infer them from prune-plan commands or the full audit table.

## Scope

- Add `integrated-only` filtering to the audit helper
- Add make targets for text and JSON integrated views
- Add an `## Integrated` section and JSON payload to the cleanup report
- Keep existing merged, stale, and prune-plan behavior unchanged

## Design

`scripts/worktree_audit.py` will accept `--integrated-only`, parallel to `--merged-only` and `--stale-only`.

The Makefile will add:

- `make --no-print-directory worktree-audit-integrated`
- `make --no-print-directory worktree-audit-integrated-json`

`scripts/worktree_cleanup_report.py` will consume the new integrated-only view in both formats:

- Markdown report gains an `## Integrated` text section
- JSON report gains an `integrated` object alongside `merged`, `stale`, and `prune_plan`

## Testing

- Add make-target tests for the integrated text and JSON views
- Extend cleanup-report tests to require the integrated section and JSON key
- Re-run `go test ./cmd/cosbench-go`
- Re-run `go test ./...`
- Re-run `go build ./...`
