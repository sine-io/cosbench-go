# Compare Local Artifact Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make the manual compare-local workflow upload its output as an artifact.

**Architecture:** Keep the current workflow and `make compare-local` command intact; just tee the output to a file and upload that file as an artifact.

**Tech Stack:** GitHub Actions YAML, Markdown docs

---

### Task 1: Update the Workflow

**Files:**
- Modify: `.github/workflows/compare-local.yml`

- [ ] **Step 1: Capture output to a file**

Change the run step so it writes compare-local output to a file while preserving console output.

- [ ] **Step 2: Upload the file as an artifact**

Use `actions/upload-artifact`.

### Task 2: Update Docs

**Files:**
- Modify: `README.md`
- Modify: `AGENTS.md`
- Modify: `BOARD.md`
- Modify: `TODO.md`

- [ ] **Step 1: Mention the artifact**

Document that the manual compare-local workflow leaves downloadable output.

### Task 3: Final Verification

**Files:**
- Review only: `.github/workflows/compare-local.yml`
- Review only: `README.md`
- Review only: `AGENTS.md`
- Review only: `BOARD.md`
- Review only: `TODO.md`

- [ ] **Step 1: Re-read the workflow**

Run:
```bash
sed -n '1,220p' .github/workflows/compare-local.yml
```

Expected:
- the workflow both runs compare-local and uploads an artifact

- [ ] **Step 2: Run local verification**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./...
GO=$(which go || echo /snap/bin/go) go build ./...
```

Expected:
- local verification still passes
