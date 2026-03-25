# S3/SIO Smoke Workflow Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add an opt-in live-endpoint smoke-test workflow for the S3/SIO adapter that skips by default and is runnable through a dedicated `make` target.

**Architecture:** Keep the smoke workflow inside `internal/driver/s3` so it reuses the real adapter and config parsing rather than building a new CLI harness. Activation is environment-variable based; missing configuration yields `t.Skip`, while configured runs execute a minimal happy-path object lifecycle and optional SIO multipart coverage.

**Tech Stack:** Go 1.26, package-level Go tests, AWS SDK v2 S3 client, Makefile, repository docs

---

### Task 1: Add Default-Skip Smoke Tests

**Files:**
- Create: `internal/driver/s3/smoke_test.go`
- Test: `internal/driver/s3/smoke_test.go`

- [ ] **Step 1: Write the failing smoke test skeleton**

Add tests that:
- load smoke configuration from environment
- skip when required variables are absent
- define one common object lifecycle smoke path
- define one SIO-only multipart smoke path

- [ ] **Step 2: Run smoke tests to verify the red state**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/driver/s3 -run Smoke -v
```

Expected:
- tests fail because the smoke harness/config helpers do not exist yet

- [ ] **Step 3: Implement minimal smoke harness**

Implement:
- env-backed config loader using the names from the spec
- adapter construction via existing `Adapter.Init`
- unique bucket/object naming
- best-effort cleanup through `t.Cleanup`
- common happy-path assertions
- conditional multipart coverage when backend is `sio`

- [ ] **Step 4: Re-run the smoke tests without env**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/driver/s3 -run Smoke -v
```

Expected:
- the smoke tests skip cleanly with a clear reason

### Task 2: Add a Dedicated Make Target

**Files:**
- Modify: `Makefile`

- [ ] **Step 1: Add a failing expectation check for the make target**

Run:
```bash
make -n smoke-s3
```

Expected:
- target missing before implementation

- [ ] **Step 2: Add the minimal make target**

Implement:
- `smoke-s3` target running `$(GO) test ./internal/driver/s3 -run Smoke -v`

- [ ] **Step 3: Verify the target command expansion**

Run:
```bash
GO=$(which go || echo /snap/bin/go) make -n smoke-s3
```

Expected:
- prints the exact smoke test command

### Task 3: Document the Smoke Workflow

**Files:**
- Modify: `README.md`
- Modify: `AGENTS.md`
- Modify: `BOARD.md`
- Modify: `TODO.md`

- [ ] **Step 1: Document activation and usage**

Update docs with:
- required and optional environment variables
- `make smoke-s3` usage
- default-skip behavior
- note that live failures indicate endpoint/credential/config issues

- [ ] **Step 2: Update board/checklist state**

Mark the real S3/SIO smoke workflow item as done once the code and docs land.

### Task 4: Final Verification

**Files:**
- Review only: `internal/driver/s3/smoke_test.go`
- Review only: `Makefile`
- Review only: `README.md`
- Review only: `AGENTS.md`
- Review only: `BOARD.md`
- Review only: `TODO.md`

- [ ] **Step 1: Verify smoke tests skip cleanly**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./internal/driver/s3 -run Smoke -v
```

Expected:
- `SKIP` when env vars are not configured

- [ ] **Step 2: Verify the make target**

Run:
```bash
GO=$(which go || echo /snap/bin/go) make smoke-s3
```

Expected:
- same skip behavior when env vars are absent

- [ ] **Step 3: Run the full test suite**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./...
```

Expected:
- all packages pass without requiring live credentials

- [ ] **Step 4: Run the full build**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go build ./...
```

Expected:
- repository builds cleanly

- [ ] **Step 5: Review scope**

Run:
```bash
git diff -- internal/driver/s3 Makefile README.md AGENTS.md BOARD.md TODO.md
```

Expected:
- only smoke-workflow files changed
