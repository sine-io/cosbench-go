GO ?= /snap/bin/go
COMPARE_LOCAL_OUTPUT_DIR ?= .artifacts/compare-local
COMPARE_LOCAL_MANIFEST ?= testdata/workloads/compare-local-fixtures.txt
COMPARE_LOCAL_FILTER ?=

.PHONY: build compare-local compare-local-list compare-local-list-json fmt smoke-s3 test tidy validate vet worktree-audit worktree-audit-json worktree-audit-merged worktree-audit-merged-json worktree-prune-plan worktree-prune-plan-json worktree-audit-stale

build:
	$(GO) build ./...

worktree-audit:
	@python3 ./scripts/worktree_audit.py origin/main

worktree-audit-merged:
	@python3 ./scripts/worktree_audit.py --merged-only origin/main

worktree-audit-json:
	@python3 ./scripts/worktree_audit.py --json origin/main

worktree-audit-merged-json:
	@python3 ./scripts/worktree_audit.py --json --merged-only origin/main

worktree-audit-stale:
	@python3 ./scripts/worktree_audit.py --stale-only origin/main

worktree-prune-plan:
	@python3 ./scripts/worktree_prune_plan.py

worktree-prune-plan-json:
	@python3 ./scripts/worktree_prune_plan.py --json

compare-local-list:
	@python3 ./scripts/list_compare_local_fixtures.py "$(COMPARE_LOCAL_MANIFEST)" --names "$(COMPARE_LOCAL_FILTER)"

compare-local-list-json:
	@python3 ./scripts/list_compare_local_fixtures.py "$(COMPARE_LOCAL_MANIFEST)" "$(COMPARE_LOCAL_FILTER)"

compare-local:
	@dir_base="$$(basename -- "$(COMPARE_LOCAL_OUTPUT_DIR)")"; \
	if [ "$$dir_base" != "compare-local" ]; then \
		echo "COMPARE_LOCAL_OUTPUT_DIR must end with compare-local: $(COMPARE_LOCAL_OUTPUT_DIR)" >&2; \
		exit 1; \
	fi
	@if [ -n "$(COMPARE_LOCAL_FILTER)" ]; then \
		python3 ./scripts/validate_compare_local_filter.py "$(COMPARE_LOCAL_MANIFEST)" "$(COMPARE_LOCAL_FILTER)"; \
	fi
	@mkdir -p $(COMPARE_LOCAL_OUTPUT_DIR)
	@find "$(COMPARE_LOCAL_OUTPUT_DIR)" -mindepth 1 -maxdepth 1 -exec rm -rf -- {} +
	@echo "== compare-local results =="
	@echo "$(COMPARE_LOCAL_OUTPUT_DIR)"
	@while read -r name fixture; do \
		if [ -z "$$name" ] || [ "$${name#\#}" != "$$name" ]; then \
			continue; \
		fi; \
		if [ -n "$(COMPARE_LOCAL_FILTER)" ] && [ "$(COMPARE_LOCAL_FILTER)" != "all" ]; then \
			case ",$(COMPARE_LOCAL_FILTER)," in \
				*,"$$name",*) ;; \
				*) continue ;; \
			esac; \
		fi; \
		echo "== $$name =="; \
		$(GO) run ./cmd/cosbench-go "$$fixture" -backend mock -json -quiet -summary-file "$(COMPARE_LOCAL_OUTPUT_DIR)/$$name.json"; \
	done < $(COMPARE_LOCAL_MANIFEST)
	@python3 ./scripts/build_compare_local_index.py "$(COMPARE_LOCAL_MANIFEST)" "$(COMPARE_LOCAL_OUTPUT_DIR)" "$(COMPARE_LOCAL_FILTER)"

smoke-s3:
	$(GO) test ./internal/driver/s3 -run Smoke -v

test:
	$(GO) test ./...

fmt:
	$(GO) fmt ./...

tidy:
	$(GO) mod tidy

vet:
	$(GO) vet ./...

validate: vet test build
