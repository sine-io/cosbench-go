# Driver Shared Token Auth Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Protect remote driver write endpoints with one shared bearer token while keeping combined-mode loopback execution working.

**Architecture:** Keep auth at the transport boundary. Add shared-token config to app startup, a reusable driver-write auth helper in the web layer, and automatic bearer-header injection in the driver HTTP client. Do not push auth concerns into controller business logic.

**Tech Stack:** Go 1.26, existing `internal/app`, `internal/web`, `internal/driver/agent`, `cmd/server`

---

### Task 1: Add Failing Driver Write Auth Tests

**Files:**
- Modify: `internal/web/driver_api_test.go`
- Modify: `internal/driver/agent/agent_test.go`
- Modify: `internal/app/remote_integration_test.go`

- [ ] **Step 1: Write failing tests for missing, wrong, and correct token behavior**

Cover:
- protected driver write endpoints reject when controller token is missing
- missing `Authorization` header returns `401`
- wrong token returns `403`
- correct token preserves existing success behavior
- combined-mode loopback still works when token is configured

- [ ] **Step 2: Run focused tests to confirm the red state**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/web ./internal/driver/agent ./internal/app
```

Expected:
- failures because driver write endpoints are currently unauthenticated

### Task 2: Add App-Level Shared Token Wiring

**Files:**
- Modify: `internal/app/app.go`
- Modify: `internal/app/mode.go`
- Modify: `cmd/server/main.go`

- [ ] **Step 1: Extend app config with driver shared token**

Add a config field and environment-variable path for the shared token.

- [ ] **Step 2: Inject the token into combined-mode loopback setup**

Ensure loopback agent requests use the same auth path that remote drivers will use.

- [ ] **Step 3: Re-run focused app tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/app
```

Expected:
- app tests still pass, with token-aware combined mode covered

### Task 3: Add Web-Layer Driver Write Auth Guard

**Files:**
- Create: `internal/web/driver_auth.go`
- Modify: `internal/web/driver_api.go`

- [ ] **Step 1: Add a reusable auth helper for protected driver write endpoints**

The helper should:
- read configured shared token
- reject with `503` if token is unset
- reject with `401` for missing or malformed headers
- reject with `403` for wrong tokens

- [ ] **Step 2: Apply the helper to all driver write endpoints**

Protect:
- register
- heartbeat
- claim
- events
- samples
- complete

- [ ] **Step 3: Re-run focused web tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/web
```

Expected:
- driver write auth tests pass

- [ ] **Step 4: Commit the controller auth slice**

Run:
```bash
git add internal/web/driver_auth.go internal/web/driver_api.go internal/web/driver_api_test.go internal/app/app.go internal/app/mode.go cmd/server/main.go internal/app/remote_integration_test.go
git commit -m "feat: require shared token for driver writes"
```

### Task 4: Add Agent Bearer Token Support

**Files:**
- Modify: `internal/driver/agent/http_client.go`
- Modify: `internal/driver/agent/agent.go`
- Modify: `internal/driver/agent/agent_test.go`

- [ ] **Step 1: Extend the driver HTTP client with bearer token support**

Write requests should automatically attach:

```http
Authorization: Bearer <token>
```

- [ ] **Step 2: Wire token usage through the agent**

No call site should assemble auth headers manually.

- [ ] **Step 3: Re-run focused agent tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/driver/agent
```

Expected:
- agent tests pass with protected write flows

- [ ] **Step 4: Commit the agent auth slice**

Run:
```bash
git add internal/driver/agent/http_client.go internal/driver/agent/agent.go internal/driver/agent/agent_test.go
git commit -m "feat: send bearer token on driver writes"
```

### Task 5: Final Verification And Docs Refresh

**Files:**
- Modify: `README.md`
- Modify: `docs/remote-worker-seams.md`

- [ ] **Step 1: Document the shared-token requirement**

Describe:
- controller-side token config
- driver-side token config
- combined-mode behavior

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
git diff -- internal/app internal/web internal/driver/agent cmd/server README.md docs/remote-worker-seams.md
```

Expected:
- the slice stays focused on shared-token auth for driver write endpoints

- [ ] **Step 5: Commit the docs slice**

Run:
```bash
git add README.md docs/remote-worker-seams.md
git commit -m "docs: record driver shared token auth"
```
