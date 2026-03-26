# Worktree JSON View Envelope Design

## Goal

Reduce shape drift across the machine-readable worktree helpers by giving the single-view helpers a consistent `views` envelope, similar to `worktree-cleanup-report-json`.

## Scope

- Add `views.audit` to `worktree-audit-json`
- Add `views.prune_plan` to `worktree-prune-plan-json`
- Keep existing top-level `summary` and `rows` fields for backward compatibility

## Design

The new additive JSON shapes are:

```json
{
  "generated_at": "...",
  "views": {
    "audit": {
      "summary": { ... },
      "rows": [ ... ]
    }
  },
  "summary": { ... },
  "rows": [ ... ]
}
```

and

```json
{
  "generated_at": "...",
  "views": {
    "prune_plan": {
      "summary": { ... },
      "rows": [ ... ]
    }
  },
  "summary": { ... },
  "rows": [ ... ]
}
```

New consumers should prefer `views`, while old consumers can keep using the top-level aliases.

## Testing

- Extend the existing audit/prune JSON tests to require the new `views` envelope
- Re-run `go test ./cmd/cosbench-go`
- Re-run `go test ./...`
- Re-run `go build ./...`
