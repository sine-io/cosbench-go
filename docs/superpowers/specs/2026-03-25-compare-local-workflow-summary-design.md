# Compare Local Workflow Summary Design

## Goal

Make the manual compare-local workflow easier to inspect directly in GitHub Actions.

The workflow already uploads the artifact directory; this change adds a concise job summary derived from `index.json`.

## Recommended Approach

After `make compare-local` runs, add a workflow step that reads `.artifacts/compare-local/index.json` and writes a short Markdown section to `$GITHUB_STEP_SUMMARY`.

Keep the summary thin:

- artifact directory path
- one table row per fixture
- links not required

This improves immediate usability without changing the artifact layout.

## Success Criteria

1. the manual compare-local workflow writes a readable job summary
2. the summary reflects the generated `index.json`
3. docs mention the summary alongside the artifact upload
4. `go test ./...` and `go build ./...` remain green
