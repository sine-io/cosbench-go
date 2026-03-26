# Worktree Cleanup Report JSON Views Design

## Goal

Expose a single consistent section container in `worktree-cleanup-report-json` so downstream consumers can parse the report through one stable entrypoint instead of hard-coding multiple top-level keys.

## Scope

- Add a top-level `views` object to `worktree-cleanup-report-json`
- Keep existing top-level section keys for backward compatibility
- Do not change Markdown output

## Design

`worktree-cleanup-report-json` will now include:

```json
{
  "generated_at": "...",
  "summary": { ... },
  "views": {
    "merged": { ... },
    "integrated": { ... },
    "stale": { ... },
    "prune_candidates": { ... },
    "prune_plan": { ... }
  },
  "merged": { ... },
  "integrated": { ... },
  "stale": { ... },
  "prune_candidates": { ... },
  "prune_plan": { ... }
}
```

The duplicated top-level keys remain temporarily so existing consumers do not break. New consumers should prefer `views`.

## Testing

- Extend the cleanup-report JSON test to require `views`
- Re-run `go test ./cmd/cosbench-go`
- Re-run `go test ./...`
- Re-run `go build ./...`
