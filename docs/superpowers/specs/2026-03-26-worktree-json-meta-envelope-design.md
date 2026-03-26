# Worktree JSON Meta Envelope Design

## Goal

Reduce metadata drift across the machine-readable worktree helpers by giving them a consistent top-level `meta` object.

## Scope

- Add `meta` to `worktree-audit-json`
- Add `meta` to `worktree-prune-plan-json`
- Add `meta` to `worktree-cleanup-report-json`
- Keep existing top-level metadata fields for backward compatibility

## Design

Each helper will now include:

```json
{
  "generated_at": "...",
  "meta": {
    "generated_at": "...",
    "base_ref": "origin/main",
    "current_worktree": "/abs/path"
  }
}
```

The `meta` container is additive. Current top-level fields such as `generated_at`, `summary`, `rows`, and `views` remain intact so existing consumers do not break.

## Testing

- Extend existing Go make-target tests to require `meta`
- Re-run `go test ./cmd/cosbench-go`
- Re-run `go test ./...`
- Re-run `go build ./...`
