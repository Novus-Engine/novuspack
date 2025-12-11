# Makefile for NovusPack Go API v1
# This Makefile provides build and test targets for the Go v1 implementation.
# The BDD tests reference feature files from ../../../features/ (repository root).

.PHONY: test bdd bdd-ci bdd-domain ci tidy

BDD_TAGS ?= '~@skip && ~@wip'
BDD_DOMAIN ?= ''
BDD_OUTPUT_DIR ?= ../../../tmp
BDD_OUTPUT_FILE ?= $(BDD_OUTPUT_DIR)/bdd_test_output_$(shell date +%Y%m%d_%H%M%S).txt

# Tidy dependencies - must use bdd build tag to preserve BDD dependencies
tidy:
	@echo "Running go mod tidy with bdd build tag..."
	@GOFLAGS=-tags=bdd go mod tidy
	@echo "Dependencies tidied successfully."

# Run all unit tests
test:
	go test -v ./...

# Run BDD tests using shared feature files from repository root
# Feature files are located at ../../../features/ relative to this directory
# Output is saved to a timestamped file in tmp/ for review
bdd:
	@echo "Running BDD tests... Output will be saved to $(BDD_OUTPUT_FILE)"
	@mkdir -p $(BDD_OUTPUT_DIR)
	@GOFLAGS=-tags=bdd go test -v ./_bdd > $(BDD_OUTPUT_FILE) 2>&1; \
	exit_code=$$?; \
	cat $(BDD_OUTPUT_FILE); \
	echo ""; \
	echo "BDD test output saved to: $(BDD_OUTPUT_FILE)"; \
	echo "Exit code: $$exit_code"; \
	exit $$exit_code

# Run BDD tests with tag filtering (for CI)
# Output is saved to a timestamped file in tmp/ for review
bdd-ci:
	@echo "Running BDD tests (CI mode)... Output will be saved to $(BDD_OUTPUT_FILE)"
	@mkdir -p $(BDD_OUTPUT_DIR)
	@GOFLAGS=-tags=bdd go test -v ./_bdd -args --godog.tags=$(BDD_TAGS) > $(BDD_OUTPUT_FILE) 2>&1; \
	exit_code=$$?; \
	cat $(BDD_OUTPUT_FILE); \
	echo ""; \
	echo "BDD test output saved to: $(BDD_OUTPUT_FILE)"; \
	echo "Exit code: $$exit_code"; \
	exit $$exit_code

# Run BDD tests for a specific domain
# Usage: make bdd-domain BDD_DOMAIN='@domain:file_format'
# Available domains: @domain:basic_ops, @domain:core, @domain:file_format, @domain:file_mgmt,
#                    @domain:file_types, @domain:compression, @domain:signatures, @domain:streaming,
#                    @domain:dedup, @domain:metadata, @domain:metadata_system, @domain:security_validation,
#                    @domain:security_encryption, @domain:generics, @domain:validation, @domain:testing,
#                    @domain:writing
bdd-domain:
	@if [ -z "$(BDD_DOMAIN)" ]; then \
		echo "Error: BDD_DOMAIN must be set. Example: make bdd-domain BDD_DOMAIN='@domain:file_format'"; \
		exit 1; \
	fi
	@echo "Running BDD tests for domain $(BDD_DOMAIN)... Output will be saved to $(BDD_OUTPUT_FILE)"
	@mkdir -p $(BDD_OUTPUT_DIR)
	@tags="$(BDD_DOMAIN) && $(BDD_TAGS)"; \
	GOFLAGS=-tags=bdd go test -v ./_bdd -args --godog.tags=$$tags > $(BDD_OUTPUT_FILE) 2>&1; \
	exit_code=$$?; \
	cat $(BDD_OUTPUT_FILE); \
	echo ""; \
	echo "BDD test output saved to: $(BDD_OUTPUT_FILE)"; \
	echo "Exit code: $$exit_code"; \
	exit $$exit_code

lint:
	@echo "Checking code formatting..."
	@if [ "$$(gofmt -s -l . | wc -l)" -gt 0 ]; then \
		echo "Code is not formatted. Run 'go fmt ./...'"; \
		gofmt -s -d .; \
		exit 1; \
	fi
	@echo ""
	@echo "Running go vet..."
	@go vet ./...
	@echo ""
	@echo "Running golangci-lint..."
	@golangci-lint run ./...
	@golangci-lint run --build-tags=bdd ./...

# CI checks - performs same checks as GitHub Actions workflow
# NOTE: This target must be kept in sync with .github/workflows/go-ci.yml.
#       When modifying CI checks, update both this Makefile and the workflow file.
# Runs: dependency verification, unit tests, build, formatting check, go vet, and golangci-lint
ci:
	@echo "Running CI checks (same as GitHub Actions workflow)..."
	@echo ""
	@echo "1. Verifying dependencies..."
	@go mod verify
	@echo ""
	@echo "2. Running unit tests..."
	@$(MAKE) test
	@echo ""
	@echo "3. Building code..."
	@go build ./...
	@echo ""
	@echo "4. Checking formatting..."
	@if [ "$$(gofmt -s -l . | wc -l)" -gt 0 ]; then \
		echo "Code is not formatted. Run 'go fmt ./...'"; \
		gofmt -s -d .; \
		exit 1; \
	fi
	@echo ""
	@echo "5. Running go vet..."
	@go vet ./...
	@echo ""
	@echo "6. Running golangci-lint..."
	@golangci-lint run ./...
	@golangci-lint run --build-tags=bdd ./...
