# Compare Local Workflow Summary Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a readable GitHub Actions job summary for the manual compare-local workflow.

**Architecture:** Keep the artifact upload, then render a small Markdown summary from `.artifacts/compare-local/index.json` into `$GITHUB_STEP_SUMMARY`.

**Tech Stack:** GitHub Actions YAML, inline Python, Markdown docs

---

### Task 1: Update the Manual Workflow

**Files:**
- Modify: `.github/workflows/compare-local.yml`

- [ ] **Step 1: Add a summary step**

Render `index.json` into a Markdown table appended to `$GITHUB_STEP_SUMMARY`.

### Task 2: Update Docs

**Files:**
- Modify: `README.md`
- Modify: `AGENTS.md`
- Modify: `BOARD.md`
- Modify: `TODO.md`

- [ ] **Step 1: Mention the job summary**

Document that the manual workflow now produces both an artifact directory and a GitHub job summary.

### Task 3: Final Verification

- [ ] **Step 1: Re-read the workflow**

Run:
```bash
sed -n '1,240p' .github/workflows/compare-local.yml
```

- [ ] **Step 2: Run repository verification**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./...
GO=$(which go || echo /snap/bin/go) go build ./...
```
