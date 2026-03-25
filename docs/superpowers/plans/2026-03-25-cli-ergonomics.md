# CLI Ergonomics Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Improve local `cosbench-go` CLI usability while staying backward-compatible with the current single-command runner.

**Architecture:** Keep the CLI as one command. Add friendlier workload path resolution (`-workload`, `-f`, positional), make `-json` stdout clean for scripting, and update docs to reflect the improved invocation patterns.

**Tech Stack:** Go 1.26, package-local tests in `cmd/cosbench-go`, Markdown docs

---

### Task 1: Add Failing CLI Tests

**Files:**
- Modify: `cmd/cosbench-go/main_test.go`

- [ ] **Step 1: Write failing tests**

Add tests proving:
- `-f` works as a workload path alias
- positional workload path works
- `-json` emits pure JSON without progress text before the payload
- missing workload path still errors cleanly

- [ ] **Step 2: Run CLI tests to verify the red state**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go
```

Expected:
- failures because the current CLI only supports `-workload` and still mixes progress text into JSON mode

### Task 2: Implement Friendlier Argument Handling

**Files:**
- Modify: `cmd/cosbench-go/main.go`

- [ ] **Step 1: Add workload path resolution priority**

Implement workload selection in this order:
- `-workload`
- `-f`
- first positional argument

- [ ] **Step 2: Improve usage/help text**

Make usage show the supported forms clearly.

- [ ] **Step 3: Make JSON mode machine-readable**

Ensure:
- stdout in `-json` mode contains only JSON
- progress lines are redirected elsewhere or suppressed

- [ ] **Step 4: Re-run CLI tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go
```

Expected:
- CLI package passes with the new invocation shapes

### Task 3: Update Repository Docs

**Files:**
- Modify: `README.md`
- Modify: `AGENTS.md`
- Modify: `BOARD.md`
- Modify: `TODO.md`

- [ ] **Step 1: Update examples and command guidance**

Document:
- `-f`
- positional workload path
- pure JSON behavior

- [ ] **Step 2: Update board/checklist state**

Reflect that local CLI ergonomics improved without changing the overall command model.

### Task 4: Final Verification

**Files:**
- Review only: `cmd/cosbench-go/main.go`
- Review only: `cmd/cosbench-go/main_test.go`
- Review only: `README.md`
- Review only: `AGENTS.md`
- Review only: `BOARD.md`
- Review only: `TODO.md`

- [ ] **Step 1: Run CLI package tests**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./cmd/cosbench-go
```

Expected:
- CLI package passes

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

- [ ] **Step 4: Review scope**

Run:
```bash
git diff -- cmd/cosbench-go README.md AGENTS.md BOARD.md TODO.md
```

Expected:
- slice stays limited to ergonomics, not new CLI features
