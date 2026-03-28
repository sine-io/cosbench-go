# Driver Pages In Unified Service Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add driver-facing pages to the unified Go service so operators can inspect driver health, missions, workers, and logs without maintaining a separate driver-web package.

**Architecture:** Reuse the remote split APIs and state model, add driver read helpers where needed, render driver pages under `/driver/...`, and update shared layout/navigation so controller and driver surfaces coexist cleanly in one service.

**Tech Stack:** Go 1.26, existing `internal/web`, Go templates under `web/templates`, remote state and APIs from the remote split project

---

### Task 1: Add Driver Read Models And API Completion

**Files:**
- Create: `internal/controlplane/driver_readmodel.go`
- Create: `internal/controlplane/driver_readmodel_test.go`
- Modify: `internal/web/driver_api.go`
- Modify: `internal/web/driver_api_test.go`

- [ ] **Step 1: Write failing tests for driver overview, mission list, and worker read models**

Cover:
- driver self-summary
- active/recent mission summaries
- worker-pool state
- driver log/health excerpts

- [ ] **Step 2: Run focused controller and web tests to confirm the red state**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/controlplane ./internal/web
```

Expected:
- failures because driver-focused read helpers and API payloads are incomplete

- [ ] **Step 3: Implement driver read models and any missing API fields**

Keep the API surface authoritative and page-friendly before adding HTML rendering.

- [ ] **Step 4: Re-run focused tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/controlplane ./internal/web
```

Expected:
- controller and web packages pass with driver read coverage

- [ ] **Step 5: Commit the read-model slice**

Run:
```bash
git add internal/controlplane/driver_readmodel.go internal/controlplane/driver_readmodel_test.go internal/web/driver_api.go internal/web/driver_api_test.go
git commit -m "feat: add driver read models"
```

### Task 2: Add Driver Page Routes And Templates

**Files:**
- Create: `internal/web/driver_pages.go`
- Create: `web/templates/driver_dashboard.html`
- Create: `web/templates/driver_missions.html`
- Create: `web/templates/driver_mission_detail.html`
- Create: `web/templates/driver_workers.html`
- Modify: `internal/web/handler.go`
- Modify: `internal/web/handler_test.go`

- [ ] **Step 1: Write failing page-render tests for the new driver routes**

Cover:
- `/driver`
- `/driver/missions`
- `/driver/missions/:id`
- `/driver/workers`

- [ ] **Step 2: Run web tests to confirm the red state**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/web
```

Expected:
- failures because driver page routes and templates do not exist

- [ ] **Step 3: Implement driver page handlers and templates**

Render all pages from the unified driver APIs and read models rather than direct ad hoc state assembly.

- [ ] **Step 4: Re-run web tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/web
```

Expected:
- driver page coverage passes

- [ ] **Step 5: Commit the page slice**

Run:
```bash
git add internal/web/driver_pages.go web/templates/driver_dashboard.html web/templates/driver_missions.html web/templates/driver_mission_detail.html web/templates/driver_workers.html internal/web/handler.go internal/web/handler_test.go
git commit -m "feat: add driver pages"
```

### Task 3: Add Driver Logs View And Shared Navigation

**Files:**
- Create: `web/templates/driver_logs.html`
- Modify: `internal/web/driver_pages.go`
- Modify: `web/templates/layout.html`
- Modify: `web/templates/dashboard.html`
- Modify: `web/templates/endpoints.html`
- Modify: `internal/web/handler_test.go`

- [ ] **Step 1: Write failing tests for driver logs and role-aware navigation**

Cover:
- `/driver/logs`
- shared layout links for controller and driver areas
- graceful empty states when no driver logs or missions exist

- [ ] **Step 2: Run web tests to confirm the red state**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/web
```

Expected:
- failures because driver logs page and shared navigation are not implemented

- [ ] **Step 3: Implement driver logs page and layout updates**

Keep controller and driver surfaces visually distinct while still using one shared layout system.

- [ ] **Step 4: Re-run web tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/web
```

Expected:
- web package passes with route and navigation coverage

- [ ] **Step 5: Commit the navigation slice**

Run:
```bash
git add web/templates/driver_logs.html internal/web/driver_pages.go web/templates/layout.html web/templates/dashboard.html web/templates/endpoints.html internal/web/handler_test.go
git commit -m "feat: add driver logs page"
```

### Task 4: Final Verification And Documentation Refresh

**Files:**
- Modify: `README.md`
- Modify: `docs/migration-gap-analysis.md`

- [ ] **Step 1: Update docs to describe unified driver pages**

Document:
- new driver routes
- unified service navigation
- dependency on the remote split API/state model

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

- [ ] **Step 4: Review final scope**

Run:
```bash
git diff -- internal/controlplane internal/web web/templates README.md docs/migration-gap-analysis.md
```

Expected:
- the slice stays focused on driver-facing pages and shared navigation

- [ ] **Step 5: Commit the documentation updates**

Run:
```bash
git add README.md docs/migration-gap-analysis.md
git commit -m "docs: record unified driver pages"
```
