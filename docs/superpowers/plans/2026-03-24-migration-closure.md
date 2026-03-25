# Migration Closure Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Close the repository's migration on a self-consistent, local-only v1 boundary by aligning migration docs and implementing the remaining local execution/reporting gaps.

**Architecture:** Keep the current single-process control-plane plus executor design. Treat remote controller/worker split as explicitly deferred, add only the missing local behaviors (`mfilewrite`, work-level result visibility, preflight validation, real delay behavior), and expose them through the existing snapshot/export/UI seams.

**Tech Stack:** Go 1.26, standard library HTTP/templates, package-local unit tests, repository docs under `docs/`

---

### Task 1: Converge Migration Documents

**Files:**
- Modify: `README.md`
- Modify: `BOARD.md`
- Modify: `TODO.md`
- Modify: `docs/migration-spec-v1.md`
- Modify: `docs/migration-gap-analysis.md`
- Modify: `docs/xml-compat-matrix.md`

- [ ] **Step 1: Re-read the approved spec and mark document mismatches**

Read:
- `docs/superpowers/specs/2026-03-24-migration-closure-design.md`
- `docs/migration-spec-v1.md`
- `docs/migration-gap-analysis.md`
- `docs/xml-compat-matrix.md`
- `README.md`
- `BOARD.md`
- `TODO.md`

Expected mismatches to fix:
- remote controller/driver split still presented as v1 scope in some docs
- export gaps still listed even though JSON/CSV export exists
- current completion target not clearly described as local-only

- [ ] **Step 2: Rewrite migration scope to a single local-only v1 story**

Update wording so all migration-facing docs agree on:
- active XML subset compatibility
- S3/SIO storage support
- local single-process control plane and executor
- snapshot persistence, web inspection, JSON/CSV export
- remote split explicitly deferred

- [ ] **Step 3: Update checklist/board status to match the new boundary**

Make `BOARD.md` and `TODO.md` reflect:
- remote split is deferred, not blocking this phase
- local-only gaps remain open until code tasks below land

- [ ] **Step 4: Verify document consistency**

Run:
```bash
rg -n "controller-driver HTTP|sample upload|driver registration|export formats.*missing|remote worker" README.md BOARD.md TODO.md docs
```

Expected:
- only deferred/future sections mention remote split
- no document still claims JSON/CSV export is missing

### Task 2: Add Work-Level Result Persistence and Exposure

**Files:**
- Modify: `internal/domain/job.go`
- Modify: `internal/controlplane/manager.go`
- Modify: `internal/controlplane/manager_test.go`
- Modify: `internal/web/handler.go`
- Modify: `internal/web/handler_test.go`
- Modify: `web/templates/job_detail.html`

- [ ] **Step 1: Write failing tests for work-level result visibility**

Add tests that assert:
- `internal/controlplane/manager_test.go`: final `JobResult` includes per-stage work summaries
- `internal/web/handler_test.go`: JSON export contains work totals, CSV export includes `scope=work`, and job detail renders a work summary table

- [ ] **Step 2: Run targeted tests to verify they fail for the right reason**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/controlplane ./internal/web
```

Expected:
- failures because work-level result structures are absent from result assembly/export/rendering

- [ ] **Step 3: Extend the result model and control-plane aggregation**

Implement:
- a work-level DTO in `internal/domain/job.go`
- `JobResult` stage entries that can carry per-work summaries
- control-plane accumulation of each `WorkResult` into persisted `JobResult`

- [ ] **Step 4: Expose work-level results in exports and UI**

Implement:
- JSON export includes the new work-level payload automatically through `JobResult`
- CSV export writes rows with `scope=work`
- `web/templates/job_detail.html` renders a work summary table grouped by stage

- [ ] **Step 5: Re-run targeted tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/controlplane ./internal/web
```

Expected:
- target packages pass with work-level visibility covered

### Task 3: Implement Real `mfilewrite` and `delay` Semantics

**Files:**
- Modify: `internal/domain/execution/opconfig.go`
- Modify: `internal/domain/execution/opconfig_test.go`
- Modify: `internal/domain/execution/engine.go`
- Modify: `internal/domain/execution/engine_test.go`

- [ ] **Step 1: Write failing tests for file-backed multipart upload and real delay**

Add tests that assert:
- `mfilewrite` reads bytes from a configured local file and uploads them through the storage adapter
- `delay` consumes measurable wall-clock time instead of returning immediately
- invalid file configuration returns a concrete error

- [ ] **Step 2: Run targeted execution tests to verify the failures**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/domain/execution
```

Expected:
- failures because `mfilewrite` and `delay` are currently placeholder no-ops

- [ ] **Step 3: Implement minimal real behavior**

Implement:
- file-source resolution in `ParsedOpConfig` using existing `files` / `fileselection` fields
- `mfilewrite` as multipart upload using the selected file's content and size
- `delay` as a real `time.Sleep`/context-aware wait using config-derived duration when present, otherwise one interval tick

Keep scope minimal:
- no new storage port methods
- no remote-driver abstractions
- no speculative file-discovery subsystem beyond what tests require

- [ ] **Step 4: Re-run targeted execution tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/domain/execution
```

Expected:
- execution package passes with new behavior covered

### Task 4: Add Start-Time Preflight Validation

**Files:**
- Modify: `internal/controlplane/manager.go`
- Modify: `internal/controlplane/manager_test.go`

- [ ] **Step 1: Write failing tests for preflight rejection**

Add tests that assert `StartJob` returns an error before spawning execution when:
- effective storage cannot be resolved
- adapter creation/init fails immediately
- operation config is invalid in a way detectable before execution
- `mfilewrite` references an unreadable file

- [ ] **Step 2: Run targeted control-plane tests to verify they fail**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/controlplane
```

Expected:
- failures because `StartJob` currently transitions to running without preflight checks

- [ ] **Step 3: Implement preflight validation in the control plane**

Implement a synchronous validation pass in `StartJob` or a dedicated helper that:
- resolves effective storage for every work
- instantiates and initializes adapters once for validation
- parses operation config before background execution
- checks file readability for file-backed ops

On failure:
- return an error from `StartJob`
- keep job state out of `running`

- [ ] **Step 4: Re-run targeted control-plane tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/controlplane
```

Expected:
- preflight tests pass and existing lifecycle tests stay green

### Task 5: Final Verification and Status Refresh

**Files:**
- Modify: `BOARD.md`
- Modify: `TODO.md`
- Modify: `docs/migration-gap-analysis.md` (if implementation changed any stated open gaps)

- [ ] **Step 1: Reconcile board/checklist with implemented state**

Update task tracking after code lands so the repository clearly shows what closed in this phase and what remains deferred.

- [ ] **Step 2: Run the full test suite**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./...
```

Expected:
- all packages pass

- [ ] **Step 3: Run the full build**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go build ./...
```

Expected:
- repository builds cleanly

- [ ] **Step 4: Review diff for scope discipline**

Run:
```bash
git status --short
git diff -- docs README.md BOARD.md TODO.md internal web
```

Expected:
- only migration-closure files changed
- no accidental refactors outside the planned scope

- [ ] **Step 5: Prepare branch-completion handoff**

Because the user asked to proceed inline, finish by summarizing:
- what changed
- what remains deferred
- exact verification commands and outcomes
