# Legacy Live Run Checklist Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a dedicated live-run checklist so legacy-comparison work is process-ready once an endpoint and credentials are available.

**Architecture:** Keep this slice documentation-only. Create one focused checklist document, keep the comparison matrix as the findings ledger, and update repository docs/status so the remaining blocker is clearly environment availability rather than missing process.

**Tech Stack:** Markdown documentation under `docs/`

---

### Task 1: Create the Checklist Document

**Files:**
- Create: `docs/legacy-live-run-checklist.md`

- [ ] **Step 1: Draft the checklist**

Include:
- preconditions
- smoke precheck
- recommended run order
- recording rules
- known watchpoints

- [ ] **Step 2: Re-read the checklist**

Run:
```bash
sed -n '1,260p' docs/legacy-live-run-checklist.md
```

Expected:
- the document is concrete and actionable

### Task 2: Link the Checklist from Existing Docs

**Files:**
- Modify: `docs/legacy-comparison-matrix.md`
- Modify: `README.md`
- Modify: `BOARD.md`
- Modify: `TODO.md`

- [ ] **Step 1: Link from the matrix and README**

Make it clear that:
- the matrix records findings
- the new checklist explains how to run future live checks

- [ ] **Step 2: Update board/todo wording**

Reflect that:
- live comparison is process-ready
- the remaining blocker is live environment availability

### Task 3: Final Verification

**Files:**
- Review only: `docs/legacy-live-run-checklist.md`
- Review only: `docs/legacy-comparison-matrix.md`
- Review only: `README.md`
- Review only: `BOARD.md`
- Review only: `TODO.md`

- [ ] **Step 1: Re-read the new checklist**

Run:
```bash
sed -n '1,260p' docs/legacy-live-run-checklist.md
```

Expected:
- the runbook is coherent and executable once environment is available

- [ ] **Step 2: Run the full test suite**

Run:
```bash
GO=$(which go || echo /snap/bin/go) go test ./...
```

Expected:
- all packages pass unchanged

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
git diff -- docs/legacy-live-run-checklist.md docs/legacy-comparison-matrix.md README.md BOARD.md TODO.md
```

Expected:
- the slice remains documentation/status only
