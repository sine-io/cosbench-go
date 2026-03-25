GO ?= /snap/bin/go
COMPARE_LOCAL_OUTPUT_DIR ?= .artifacts/compare-local

.PHONY: build compare-local fmt smoke-s3 test tidy validate vet

build:
	$(GO) build ./...

compare-local:
	@mkdir -p $(COMPARE_LOCAL_OUTPUT_DIR)
	@echo "== compare-local results =="
	@echo "$(COMPARE_LOCAL_OUTPUT_DIR)"
	@echo "== s3-active-subset =="
	@$(GO) run ./cmd/cosbench-go testdata/workloads/s3-active-subset.xml -backend mock -json -quiet -summary-file $(COMPARE_LOCAL_OUTPUT_DIR)/s3-active-subset.json
	@echo "== mock-stage-aware =="
	@$(GO) run ./cmd/cosbench-go testdata/workloads/mock-stage-aware.xml -backend mock -json -quiet -summary-file $(COMPARE_LOCAL_OUTPUT_DIR)/mock-stage-aware.json
	@echo "== mock-reusedata-subset =="
	@$(GO) run ./cmd/cosbench-go testdata/workloads/mock-reusedata-subset.xml -backend mock -json -quiet -summary-file $(COMPARE_LOCAL_OUTPUT_DIR)/mock-reusedata-subset.json
	@echo "== xml-splitrw-subset =="
	@$(GO) run ./cmd/cosbench-go testdata/workloads/xml-splitrw-subset.xml -backend mock -json -quiet -summary-file $(COMPARE_LOCAL_OUTPUT_DIR)/xml-splitrw-subset.json

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
