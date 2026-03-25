# Worktree Prune Plan JSON Summary Design

## Goal

Align `make worktree-prune-plan-json` with the structured shape already used by `make worktree-audit-json` so automation can consume prune-plan metadata without recomputing counts from the row array.

## Scope

This change is limited to the machine-readable prune-plan helper and its direct consumers.

- Keep the text output of `scripts/worktree_prune_plan.py` unchanged
- Change `--json` output from a bare array to an object with `summary` and `rows`
- Preserve the existing per-row fields, including the recently added branch context
- Update repo-local consumers and tests to use the structured payload
- Update repository docs to describe the new JSON shape

## Design

`scripts/worktree_prune_plan.py --json` will return:

```json
{
  "summary": {
    "base_ref": "origin/main",
    "current_worktree": "/abs/path",
    "total": 3,
    "merged": 1,
    "integrated": 2
  },
  "rows": [
    {
      "path": "/abs/path/.worktrees/example",
      "branch": "example",
      "state": "integrated",
      "details": "patch-equivalent to origin/main",
      "ahead": 1,
      "behind": 0,
      "commands": [
        "git worktree remove '/abs/path/.worktrees/example'",
        "git branch -D example"
      ]
    }
  ]
}
```

The `summary` block is intentionally small: it should describe the prune-plan result itself, not duplicate the full audit summary.

## Rationale

- The current array-only output is harder to reuse in GitHub workflows and report generators
- `worktree-audit-json` already established a `summary + rows` convention in this repository
- The prune-plan helper is explicitly repo-local; updating the schema together with tests and docs is a safe cleanup

## Testing

- Update the existing make target test to expect `summary` and `rows`
- Add a direct script test covering the structured summary values and row context
- Re-run `go test ./cmd/cosbench-go`, `go test ./...`, and `go build ./...`
