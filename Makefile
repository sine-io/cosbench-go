GO ?= /snap/bin/go
COMPARE_LOCAL_OUTPUT_DIR ?= .artifacts/compare-local
COMPARE_LOCAL_MANIFEST ?= testdata/workloads/compare-local-fixtures.txt
COMPARE_LOCAL_FILTER ?=

.PHONY: build compare-local fmt smoke-s3 test tidy validate vet

build:
	$(GO) build ./...

compare-local:
	@dir_base="$$(basename -- "$(COMPARE_LOCAL_OUTPUT_DIR)")"; \
	if [ "$$dir_base" != "compare-local" ]; then \
		echo "COMPARE_LOCAL_OUTPUT_DIR must end with compare-local: $(COMPARE_LOCAL_OUTPUT_DIR)" >&2; \
		exit 1; \
	fi
	@if [ -n "$(COMPARE_LOCAL_FILTER)" ]; then \
		awk -v want="$(COMPARE_LOCAL_FILTER)" '\
			NF && $$1 !~ /^#/ { \
				names = names "  - " $$1 "\n"; \
				if ($$1 == want) { found = 1 } \
			} \
			END { \
				if (!found) { \
					printf "unknown compare-local fixture: %s\nknown fixtures:\n%s", want, names > "/dev/stderr"; \
					exit 1; \
				} \
			}\
		' "$(COMPARE_LOCAL_MANIFEST)"; \
	fi
	@mkdir -p $(COMPARE_LOCAL_OUTPUT_DIR)
	@find "$(COMPARE_LOCAL_OUTPUT_DIR)" -mindepth 1 -maxdepth 1 -exec rm -rf -- {} +
	@echo "== compare-local results =="
	@echo "$(COMPARE_LOCAL_OUTPUT_DIR)"
	@while read -r name fixture; do \
		if [ -z "$$name" ] || [ "$${name#\#}" != "$$name" ]; then \
			continue; \
		fi; \
		if [ -n "$(COMPARE_LOCAL_FILTER)" ] && [ "$$name" != "$(COMPARE_LOCAL_FILTER)" ]; then \
			continue; \
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
