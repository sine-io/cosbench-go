# Compare Local Workflow Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a manual GitHub Actions workflow that runs `make compare-local`.

**Architecture:** Keep the workflow thin and explicit. The existing `compare-local` Make target remains the source of truth; the new workflow only exposes it through `workflow_dispatch`.

**Tech Stack:** GitHub Actions YAML, existing `Makefile`, Markdown docs

---

### Task 1: Add the Manual Workflow

**Files:**
- Create: `.github/workflows/compare-local.yml`

- [ ] **Step 1: Create the workflow**

The workflow should:
- trigger only on `workflow_dispatch`
- checkout the repository
- set up Go from `go.mod`
- run `GO=go make compare-local`

- [ ] **Step 2: Re-read the workflow**

Run:
```bash
sed -n '1,220p' .github/workflows/compare-local.yml
```

Expected:
- the workflow is manual-only and wraps the existing command

### Task 2: Update Docs

**Files:**
- Modify: `README.md`
- Modify: `AGENTS.md`
- Modify: `BOARD.md`
- Modify: `TODO.md`

- [ ] **Step 1: Document the manual workflow**

Clarify:
- default CI still runs `make validate`
- manual workflow runs `make compare-local`
- live smoke remains separate

- [ ] **Step 2: Update board/checklist**

Reflect that broader external automation now exists for local comparison, but not for live smoke.

### Task 3: Final Verification

**Files:**
- Review only: `.github/workflows/compare-local.yml`
- Review only: `README.md`
- Review only: `AGENTS.md`
- Review only: `BOARD.md`
- Review only: `TODO.md`

- [ ] **Step 1: Re-read the new workflow**

Run:
```bash
sed -n '1,220p' .github/workflows/compare-local.yml
```

Expected:
- the workflow only does manual compare-local execution

- [ ] **Step 2: Run local verification**

Run:
```bash
GO=$(which go || echo /snap/bin/go) make compare-local
GO=$(which go || echo /snap/bin/go) go test ./...
GO=$(which go || echo /snap/bin/go) go build ./...
```

Expected:
- all local verification still passes
