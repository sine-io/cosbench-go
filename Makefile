GO ?= /snap/bin/go
COMPARE_LOCAL_OUTPUT_DIR ?= .artifacts/compare-local
COMPARE_LOCAL_MANIFEST ?= testdata/workloads/compare-local-fixtures.txt

.PHONY: build compare-local fmt smoke-s3 test tidy validate vet

build:
	$(GO) build ./...

compare-local:
	@dir_base="$$(basename -- "$(COMPARE_LOCAL_OUTPUT_DIR)")"; \
	if [ "$$dir_base" != "compare-local" ]; then \
		echo "COMPARE_LOCAL_OUTPUT_DIR must end with compare-local: $(COMPARE_LOCAL_OUTPUT_DIR)" >&2; \
		exit 1; \
	fi
	@mkdir -p $(COMPARE_LOCAL_OUTPUT_DIR)
	@find "$(COMPARE_LOCAL_OUTPUT_DIR)" -mindepth 1 -maxdepth 1 -exec rm -rf -- {} +
	@echo "== compare-local results =="
	@echo "$(COMPARE_LOCAL_OUTPUT_DIR)"
	@printf '{\n  "fixtures": [\n' > "$(COMPARE_LOCAL_OUTPUT_DIR)/index.json"; \
	first=1; \
	while read -r name fixture; do \
		if [ -z "$$name" ] || [ "$${name#\#}" != "$$name" ]; then \
			continue; \
		fi; \
		echo "== $$name =="; \
		$(GO) run ./cmd/cosbench-go "$$fixture" -backend mock -json -quiet -summary-file "$(COMPARE_LOCAL_OUTPUT_DIR)/$$name.json"; \
		if [ $$first -eq 0 ]; then \
			printf ',\n' >> "$(COMPARE_LOCAL_OUTPUT_DIR)/index.json"; \
		fi; \
		printf '    {"name":"%s","workload":"%s","summary":"%s.json"}' "$$name" "$$fixture" "$$name" >> "$(COMPARE_LOCAL_OUTPUT_DIR)/index.json"; \
		first=0; \
	done < $(COMPARE_LOCAL_MANIFEST); \
	printf '\n  ]\n}\n' >> "$(COMPARE_LOCAL_OUTPUT_DIR)/index.json"

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
