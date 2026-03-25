# CI Automation Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a minimal repository-local CI workflow that runs the existing validation path automatically on pushes and pull requests.

**Architecture:** Keep CI deliberately thin: one GitHub Actions workflow that installs Go from `go.mod` and runs `make validate` with `GO=go`. Do not add live smoke coverage or alternate validation logic to default CI.

**Tech Stack:** GitHub Actions YAML, existing `Makefile`, Markdown docs

---

### Task 1: Add the Workflow File

**Files:**
- Create: `.github/workflows/ci.yml`

- [ ] **Step 1: Add a minimal workflow skeleton**

The workflow should:
- trigger on `push` and `pull_request`
- checkout the repository
- set up Go from `go.mod`
- run `GO=go make validate`

- [ ] **Step 2: Validate the YAML visually**

Run:
```bash
sed -n '1,220p' .github/workflows/ci.yml
```

Expected:
- workflow is minimal and only wraps the existing validation path

### Task 2: Update Repository Docs

**Files:**
- Modify: `README.md`
- Modify: `AGENTS.md`
- Modify: `BOARD.md`
- Modify: `TODO.md`

- [ ] **Step 1: Document the CI path**

Add short notes explaining:
- CI runs `make validate`
- smoke tests stay opt-in and are not part of default CI

- [ ] **Step 2: Update board/checklist state**

Reflect that external automation now exists for the validation path.

### Task 3: Final Verification

**Files:**
- Review only: `.github/workflows/ci.yml`
- Review only: `README.md`
- Review only: `AGENTS.md`
- Review only: `BOARD.md`
- Review only: `TODO.md`

- [ ] **Step 1: Re-read the workflow and docs**

Run:
```bash
sed -n '1,220p' .github/workflows/ci.yml
```

Expected:
- workflow clearly runs only `make validate`

- [ ] **Step 2: Run local validation**

Run:
```bash
GO=$(which go || echo /snap/bin/go) make validate
```

Expected:
- local validation still passes

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
git diff -- .github/workflows/ci.yml README.md AGENTS.md BOARD.md TODO.md
```

Expected:
- slice remains automation/doc focused
