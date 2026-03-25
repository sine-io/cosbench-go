# Compare Local Output Safety Design

## Goal

Make compare-local output cleanup safer.

The current pruning change clears the configured output directory, but `COMPARE_LOCAL_OUTPUT_DIR` is user-overridable. That means cleanup needs a tighter safety boundary than a raw directory delete.

## Recommended Approach

- require `COMPARE_LOCAL_OUTPUT_DIR` to end with `compare-local`
- clear only that dedicated directory's contents instead of deleting the directory path itself
- add a regression test that an unsafe output directory is rejected

This keeps the target flexible enough for temp directories like `/tmp/.../compare-local` while making accidental broad cleanup much harder.

## Success Criteria

1. unsafe `COMPARE_LOCAL_OUTPUT_DIR` values are rejected
2. safe compare-local directories are still refreshed correctly
3. tests cover both the stale-output path and the unsafe-directory rejection path
4. `go test ./...` and `go build ./...` remain green
