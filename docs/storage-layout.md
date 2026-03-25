# Storage Layout

Phase 1 persists runtime snapshots to local files under the configured data directory.

Directory layout:

- `data/jobs/<job-id>.json`: current visible job state, including stage status and summary metrics
- `data/results/<job-id>.json`: final or latest result summary for a job
- `data/events/<job-id>.json`: ordered lifecycle and execution events
- `data/endpoints/<endpoint-id>.json`: saved endpoint configurations for reuse

Storage rules:

- no database dependency
- JSON files are written atomically through temp-file rename
- in-memory state is rebuilt from snapshots on startup
- jobs found in `running` state during restart recovery are marked `interrupted`
- jobs found in `cancelling` state during restart recovery are marked `cancelled`
