GO ?= /snap/bin/go
PYTHON ?= python3
PYTHON_ENV ?= PYTHONDONTWRITEBYTECODE=1
export PYTHON
COMPARE_LOCAL_OUTPUT_DIR ?= .artifacts/compare-local
COMPARE_LOCAL_MANIFEST ?= testdata/workloads/compare-local-fixtures.txt
COMPARE_LOCAL_FILTER ?=
WORKTREE_AUDIT_BASE_REF ?= origin/main
WORKTREE_CLEANUP_REPORT_PATH ?= .artifacts/worktree-cleanup-report.md

.PHONY: build compare-local compare-local-list compare-local-list-json fmt smoke-local smoke-ready smoke-ready-json smoke-ready-validate smoke-ready-validate-json smoke-s3 smoke-remote-local test tidy validate vet worktree-audit worktree-audit-json worktree-audit-merged worktree-audit-merged-json worktree-audit-integrated worktree-audit-integrated-json worktree-audit-prune worktree-audit-prune-json worktree-prune-plan worktree-prune-plan-json worktree-audit-stale worktree-cleanup-report worktree-cleanup-report-json

build:
	$(GO) build ./...

worktree-audit:
	@$(PYTHON_ENV) $(PYTHON) ./scripts/worktree_audit.py "$(WORKTREE_AUDIT_BASE_REF)"

worktree-audit-merged:
	@$(PYTHON_ENV) $(PYTHON) ./scripts/worktree_audit.py --merged-only "$(WORKTREE_AUDIT_BASE_REF)"

worktree-audit-json:
	@$(PYTHON_ENV) $(PYTHON) ./scripts/worktree_audit.py --json "$(WORKTREE_AUDIT_BASE_REF)"

worktree-audit-merged-json:
	@$(PYTHON_ENV) $(PYTHON) ./scripts/worktree_audit.py --json --merged-only "$(WORKTREE_AUDIT_BASE_REF)"

worktree-audit-integrated:
	@$(PYTHON_ENV) $(PYTHON) ./scripts/worktree_audit.py --integrated-only "$(WORKTREE_AUDIT_BASE_REF)"

worktree-audit-integrated-json:
	@$(PYTHON_ENV) $(PYTHON) ./scripts/worktree_audit.py --json --integrated-only "$(WORKTREE_AUDIT_BASE_REF)"

worktree-audit-prune:
	@$(PYTHON_ENV) $(PYTHON) ./scripts/worktree_audit.py --prune-only "$(WORKTREE_AUDIT_BASE_REF)"

worktree-audit-prune-json:
	@$(PYTHON_ENV) $(PYTHON) ./scripts/worktree_audit.py --json --prune-only "$(WORKTREE_AUDIT_BASE_REF)"

worktree-audit-stale:
	@$(PYTHON_ENV) $(PYTHON) ./scripts/worktree_audit.py --stale-only "$(WORKTREE_AUDIT_BASE_REF)"

worktree-prune-plan:
	@$(PYTHON_ENV) $(PYTHON) ./scripts/worktree_prune_plan.py "$(WORKTREE_AUDIT_BASE_REF)"

worktree-prune-plan-json:
	@$(PYTHON_ENV) $(PYTHON) ./scripts/worktree_prune_plan.py --json "$(WORKTREE_AUDIT_BASE_REF)"

worktree-cleanup-report:
	@mkdir -p "$$(dirname "$(WORKTREE_CLEANUP_REPORT_PATH)")"
	@$(PYTHON_ENV) $(PYTHON) ./scripts/worktree_cleanup_report.py "$(WORKTREE_AUDIT_BASE_REF)" "$(WORKTREE_CLEANUP_REPORT_PATH)"

worktree-cleanup-report-json:
	@$(PYTHON_ENV) $(PYTHON) ./scripts/worktree_cleanup_report.py --json "$(WORKTREE_AUDIT_BASE_REF)"

compare-local-list:
	@$(PYTHON_ENV) $(PYTHON) ./scripts/list_compare_local_fixtures.py "$(COMPARE_LOCAL_MANIFEST)" --names "$(COMPARE_LOCAL_FILTER)"

compare-local-list-json:
	@$(PYTHON_ENV) $(PYTHON) ./scripts/list_compare_local_fixtures.py "$(COMPARE_LOCAL_MANIFEST)" "$(COMPARE_LOCAL_FILTER)"

compare-local:
	@dir_base="$$(basename -- "$(COMPARE_LOCAL_OUTPUT_DIR)")"; \
	if [ "$$dir_base" != "compare-local" ]; then \
		echo "COMPARE_LOCAL_OUTPUT_DIR must end with compare-local: $(COMPARE_LOCAL_OUTPUT_DIR)" >&2; \
		exit 1; \
	fi
	@if [ -n "$(COMPARE_LOCAL_FILTER)" ]; then \
		$(PYTHON_ENV) $(PYTHON) ./scripts/validate_compare_local_filter.py "$(COMPARE_LOCAL_MANIFEST)" "$(COMPARE_LOCAL_FILTER)"; \
	fi
	@mkdir -p $(COMPARE_LOCAL_OUTPUT_DIR)
	@find "$(COMPARE_LOCAL_OUTPUT_DIR)" -mindepth 1 -maxdepth 1 -exec rm -rf -- {} +
	@echo "== compare-local results =="
	@echo "$(COMPARE_LOCAL_OUTPUT_DIR)"
	@fixtures_file="$$(mktemp)"; \
	trap 'rm -f "$$fixtures_file"' EXIT; \
	$(PYTHON_ENV) $(PYTHON) ./scripts/list_compare_local_fixtures.py "$(COMPARE_LOCAL_MANIFEST)" --pairs "$(COMPARE_LOCAL_FILTER)" > "$$fixtures_file"; \
	while read -r name fixture; do \
		if [ -z "$$name" ] || [ "$${name#\#}" != "$$name" ]; then \
			continue; \
		fi; \
		echo "== $$name =="; \
		$(GO) run ./cmd/cosbench-go "$$fixture" -backend mock -json -quiet -summary-file "$(COMPARE_LOCAL_OUTPUT_DIR)/$$name.json"; \
	done < "$$fixtures_file"
	@$(PYTHON_ENV) $(PYTHON) ./scripts/build_compare_local_index.py "$(COMPARE_LOCAL_MANIFEST)" "$(COMPARE_LOCAL_OUTPUT_DIR)" "$(COMPARE_LOCAL_FILTER)"

smoke-s3:
	$(GO) test ./internal/driver/s3 -run Smoke -v

smoke-local:
	@$(PYTHON_ENV) $(PYTHON) ./scripts/smoke_local.py

smoke-remote-local:
	@$(PYTHON_ENV) $(PYTHON) ./scripts/smoke_remote_local.py

smoke-ready:
	@$(PYTHON_ENV) $(PYTHON) ./scripts/smoke_ready.py

smoke-ready-json:
	@$(PYTHON_ENV) $(PYTHON) ./scripts/smoke_ready.py --json

smoke-ready-validate:
	@$(PYTHON_ENV) $(PYTHON) ./scripts/validate_smoke_ready_schema.py

smoke-ready-validate-json:
	@$(PYTHON_ENV) $(PYTHON) ./scripts/validate_smoke_ready_schema.py --json

test:
	$(GO) test ./...

fmt:
	$(GO) fmt ./...

tidy:
	$(GO) mod tidy

vet:
	$(GO) vet ./...

validate: vet test build
