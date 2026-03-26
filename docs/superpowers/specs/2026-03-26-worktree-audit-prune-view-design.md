# Worktree Audit Prune View Design

## Goal

Expose prune candidates as a first-class audit view so users can review which worktrees are considered safe-to-prune before reading command-oriented prune-plan output.

## Scope

- Add a prune-only audit filter
- Add text and JSON Make targets for that filter
- Surface the structured prune-candidate view inside cleanup-report markdown and JSON
- Keep the existing prune-plan command output unchanged

## Design

`scripts/worktree_audit.py` will accept `--prune-only`. A row qualifies when:

- `state` is `merged` or `integrated`
- `branch` is not `main` or `master`
- the row is not the current worktree

The cleanup report will gain:

- a `## Prune Candidates` text section
- a top-level `prune_candidates` JSON object

This keeps the audit-style review output separate from the command list in `prune_plan`.

## Testing

- Add make-target tests for `worktree-audit-prune` and `worktree-audit-prune-json`
- Extend cleanup-report tests to require the prune-candidates section and JSON key
- Re-run `go test ./cmd/cosbench-go`
- Re-run `go test ./...`
- Re-run `go build ./...`
