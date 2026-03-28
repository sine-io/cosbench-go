# Controller Web/API Closure Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add the deferred controller-facing web and API capabilities to the unified Go service without recreating the legacy Java UI structure.

**Architecture:** Introduce controller read models and timeline persistence first, then add API-first controller routes, then render controller pages against those APIs and read models, and finally expose Prometheus and artifact downloads through the same controller-owned query layer.

**Tech Stack:** Go 1.26, existing `internal/controlplane`, `internal/reporting`, `internal/snapshot`, `internal/web`, Go templates under `web/templates`

---

### Task 1: Add Controller Read Models And Timeline Persistence

**Files:**
- Create: `internal/domain/timeline.go`
- Create: `internal/controlplane/readmodel.go`
- Create: `internal/reporting/timeline.go`
- Modify: `internal/controlplane/manager.go`
- Modify: `internal/snapshot/store.go`
- Modify: `internal/controlplane/manager_test.go`
- Create: `internal/reporting/timeline_test.go`

- [ ] **Step 1: Write failing tests for matrix and timeline read concerns**

Add tests proving:
- controller can expose matrix-ready job summaries
- execution samples can be aggregated into persisted timeline buckets
- stage/job timeline queries survive snapshot reload

- [ ] **Step 2: Run focused controller/reporting tests to confirm the red state**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/controlplane ./internal/reporting
```

Expected:
- failures because current state only persists summary results and events

- [ ] **Step 3: Implement timeline bucket persistence and controller read models**

Add the minimal persisted/queryable structures needed for:
- matrix summaries
- stage drilldown
- timeline queries
- artifact references introduced later

- [ ] **Step 4: Re-run focused tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/controlplane ./internal/reporting
```

Expected:
- controller and reporting packages pass with timeline/read-model coverage

- [ ] **Step 5: Commit the state-model slice**

Run:
```bash
git add internal/domain/timeline.go internal/controlplane/readmodel.go internal/reporting/timeline.go internal/controlplane/manager.go internal/snapshot/store.go internal/controlplane/manager_test.go internal/reporting/timeline_test.go
git commit -m "feat: add controller timeline read models"
```

### Task 2: Add Controller API Routes

**Files:**
- Create: `internal/web/controller_api.go`
- Create: `internal/web/controller_api_test.go`
- Modify: `internal/web/handler.go`
- Modify: `internal/app/app.go`

- [ ] **Step 1: Write failing HTTP tests for controller API endpoints**

Cover:
- jobs list/detail
- config and advanced-config payloads
- stage detail
- timeline JSON
- timeline CSV

- [ ] **Step 2: Run web tests to confirm the red state**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/web
```

Expected:
- failures because the new controller API routes do not exist yet

- [ ] **Step 3: Implement controller API handlers and route wiring**

Expose:
- `/api/controller/jobs`
- `/api/controller/jobs/:id`
- `/api/controller/jobs/:id/config`
- `/api/controller/jobs/:id/config/advanced`
- `/api/controller/jobs/:id/stages/:stage`
- `/api/controller/jobs/:id/timeline`
- `/api/controller/jobs/:id/timeline.csv`

- [ ] **Step 4: Re-run web tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/web
```

Expected:
- controller API coverage passes

- [ ] **Step 5: Commit the API slice**

Run:
```bash
git add internal/web/controller_api.go internal/web/controller_api_test.go internal/web/handler.go internal/app/app.go
git commit -m "feat: add controller api routes"
```

### Task 3: Add Controller Pages And Templates

**Files:**
- Create: `internal/web/controller_pages.go`
- Create: `web/templates/controller_matrix.html`
- Create: `web/templates/controller_job_config.html`
- Create: `web/templates/controller_advanced_config.html`
- Create: `web/templates/controller_stage.html`
- Create: `web/templates/controller_timeline.html`
- Modify: `web/templates/layout.html`
- Modify: `internal/web/handler_test.go`

- [ ] **Step 1: Write failing page-render tests for controller pages**

Add coverage for:
- `/controller/matrix`
- `/controller/jobs/:id/config`
- `/controller/jobs/:id/config/advanced`
- `/controller/jobs/:id/stages/:stage`
- `/controller/jobs/:id/timeline`

- [ ] **Step 2: Run web tests to confirm the red state**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/web
```

Expected:
- failures because the controller page routes/templates are not present

- [ ] **Step 3: Implement controller page handlers and templates**

Render the new pages from the same controller read models used by the APIs.

- [ ] **Step 4: Re-run web tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/web
```

Expected:
- page-render coverage passes

- [ ] **Step 5: Commit the page slice**

Run:
```bash
git add internal/web/controller_pages.go web/templates/controller_matrix.html web/templates/controller_job_config.html web/templates/controller_advanced_config.html web/templates/controller_stage.html web/templates/controller_timeline.html web/templates/layout.html internal/web/handler_test.go
git commit -m "feat: add controller pages"
```

### Task 4: Add Prometheus, Artifact Downloads, And Final Verification

**Files:**
- Create: `internal/web/prometheus.go`
- Create: `internal/web/controller_artifacts.go`
- Modify: `internal/web/controller_api_test.go`
- Modify: `internal/web/handler_test.go`
- Modify: `README.md`

- [ ] **Step 1: Write failing tests for Prometheus and controller artifact routes**

Cover:
- `/api/controller/metrics/prometheus`
- `/api/controller/jobs/:id/artifacts/log`
- `/api/controller/jobs/:id/artifacts/config`

- [ ] **Step 2: Run web tests to confirm the red state**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/web
```

Expected:
- failures because Prometheus and artifact handlers do not exist yet

- [ ] **Step 3: Implement exporters and artifact handlers**

Use controller-owned timeline/read state rather than ad hoc route logic.

- [ ] **Step 4: Run the full test suite**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./...
```

Expected:
- all packages pass

- [ ] **Step 5: Run the full build**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go build ./...
```

Expected:
- repository builds cleanly

- [ ] **Step 6: Commit the closure slice**

Run:
```bash
git add internal/web/prometheus.go internal/web/controller_artifacts.go internal/web/controller_api_test.go internal/web/handler_test.go README.md
git commit -m "feat: close controller web api surface"
```
