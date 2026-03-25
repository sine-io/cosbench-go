GO ?= /snap/bin/go
COMPARE_LOCAL_OUTPUT_DIR ?= .artifacts/compare-local
COMPARE_LOCAL_MANIFEST ?= testdata/workloads/compare-local-fixtures.txt

.PHONY: build compare-local fmt smoke-s3 test tidy validate vet

build:
	$(GO) build ./...

compare-local:
	@rm -rf "$(COMPARE_LOCAL_OUTPUT_DIR)"
	@mkdir -p $(COMPARE_LOCAL_OUTPUT_DIR)
	@echo "== compare-local results =="
	@echo "$(COMPARE_LOCAL_OUTPUT_DIR)"
	@while read -r name fixture; do \
		if [ -z "$$name" ] || [ "$${name#\#}" != "$$name" ]; then \
			continue; \
		fi; \
		echo "== $$name =="; \
		$(GO) run ./cmd/cosbench-go "$$fixture" -backend mock -json -quiet -summary-file "$(COMPARE_LOCAL_OUTPUT_DIR)/$$name.json"; \
	done < $(COMPARE_LOCAL_MANIFEST)

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
