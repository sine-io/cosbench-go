# Worktree JSON Generated-At Design

## Goal

Add a consistent top-level `generated_at` timestamp to the machine-readable worktree helpers so artifacts can be inspected without relying on filesystem metadata.

## Scope

- Add `generated_at` to `worktree-audit-json`
- Add `generated_at` to `worktree-prune-plan-json`
- Add `generated_at` to `worktree-cleanup-report-json`
- Keep text output unchanged

## Design

Each JSON helper will emit a top-level RFC 3339 UTC timestamp with second precision, for example:

```json
{
  "generated_at": "2026-03-26T12:34:56Z",
  "...": "..."
}
```

The field should be added at the top level only. Existing `summary` and `rows` structures stay intact.

## Testing

- Extend the existing Go make-target tests to require `generated_at`
- Re-run `go test ./cmd/cosbench-go`
- Re-run `go test ./...`
- Re-run `go build ./...`
