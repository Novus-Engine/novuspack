# Makefile for NovusPack Go API v1
# This Makefile provides build and test targets for the Go v1 implementation.
# The BDD tests reference feature files from ../../../features/ (repository root).

# Ensure MAKE is set to the actual make binary for recursive calls
# Check if MAKE matches the system make (using restricted PATH to avoid AppImage interference)
# If not, force it to the system make
SYSTEM_MAKE := $(shell PATH="/usr/bin:/usr/local/bin:/bin:$$PATH" command -v make 2>/dev/null || echo /usr/bin/make)
ifneq ($(MAKE),$(SYSTEM_MAKE))
  # MAKE doesn't match system make, force it
  override MAKE := $(SYSTEM_MAKE)
endif

.PHONY: test bdd bdd-ci bdd-domain ci tidy go-fmt coverage coverage-html coverage-report validate-go-spec-references apply-go-spec-references

BDD_TAGS ?= '~@skip && ~@wip'
BDD_DOMAIN ?= ''
BDD_OUTPUT_DIR ?= ../../tmp
BDD_OUTPUT_FILE ?= $(BDD_OUTPUT_DIR)/bdd_test_output_$(shell date +%Y%m%d_%H%M%S).txt

# Tidy dependencies - must use bdd build tag to preserve BDD dependencies
tidy:
	@echo "Running go mod tidy with bdd build tag..."
	@GOFLAGS=-tags=bdd go mod tidy
	@echo "Dependencies tidied successfully."

# Format Go source files
go-fmt:
	@echo "Running go fmt..."
	@go fmt ./...
	@echo "Formatting complete."

# Run all unit tests
# Set GOCACHE to a writable location if the default cache is not writable
# Test actual write capability by attempting to create a test file
test:
	@if ( [ -w "$(HOME)/.cache/go-build" ] 2>/dev/null && touch "$(HOME)/.cache/go-build/.write-test" 2>/dev/null && rm -f "$(HOME)/.cache/go-build/.write-test" ) || [ -z "$(HOME)" ]; then \
		go test -v ./...; \
	else \
		CACHE_DIR="$$(cd ../.. && pwd)/tmp/go-cache"; \
		mkdir -p "$$CACHE_DIR"; \
		GOCACHE="$$CACHE_DIR" go test -v ./...; \
	fi

# Run tests with coverage and generate coverage.out file
coverage:
	@echo "Running tests with coverage..."
	@go clean -testcache
	@go test -coverprofile=coverage.out ./...
	@echo ""
	@echo "Coverage profile saved to: coverage.out"
	@echo "Run 'make coverage-report' to view coverage in terminal"
	@echo "Run 'make coverage-html' to generate HTML report"

# Generate HTML coverage report and open in browser
coverage-html: coverage
	@echo "Generating HTML coverage report..."
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report saved to: coverage.html"
	@if command -v xdg-open >/dev/null 2>&1; then \
		xdg-open coverage.html; \
	elif command -v open >/dev/null 2>&1; then \
		open coverage.html; \
	else \
		echo "Open coverage.html in your browser to view the report"; \
	fi

# Display coverage report in terminal
coverage-report: coverage
	@echo "=== Coverage Report ==="
	@go tool cover -func=coverage.out
	@echo ""
	@echo "=== Modified Files Coverage ==="
	@go tool cover -func=coverage.out | grep -E "package_writer\.go|package_path_metadata.*\.go|package_reader\.go" || echo "No modified files found"
	@echo ""
	@echo "=== Total Coverage ==="
	@go tool cover -func=coverage.out | grep "total:" || echo "Coverage data not available"

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

# Go signature validation - validates Go signatures in implementation against tech specs
# NOTE: This target must be kept in sync with .github/workflows/docs-check.yml.
#       When adding or modifying signature validation, update both this Makefile and
#       the workflow file to ensure local 'make validate-go-signatures' matches CI behavior.
#       Requires: Python 3
# Usage: make validate-go-signatures [VERBOSE=1] [OUTPUT="file.txt"] [NO_FAIL=1] [NO_COLOR=1] [SPECS_DIR="dir"] [IMPL_DIR="dir"]
#        - VERBOSE: Set to 1 for verbose output
#        - OUTPUT: Path to output file for validation report
#        - NO_FAIL: Set to 1 to exit with code 0 even if errors are found
#        - NO_COLOR: Set to 1 to disable colored output
#        - SPECS_DIR: Directory containing tech specs (default: docs/tech_specs)
#        - IMPL_DIR: Directory containing Go implementation (default: api/go)
validate-go-signatures:
	@command -v python3 >/dev/null 2>&1 || { \
		echo "Error: python3 not found. Install Python 3 to run signature validation."; \
		exit 1; \
	}
	@echo "Validating Go signatures against tech specs..."
	@ARGS=""; \
	if [ -n "$(VERBOSE)" ]; then ARGS="$$ARGS --verbose"; fi; \
	if [ -n "$(OUTPUT)" ]; then ARGS="$$ARGS --output \"$(OUTPUT)\""; fi; \
	if [ -n "$(NO_FAIL)" ]; then ARGS="$$ARGS --no-fail"; fi; \
	if [ -n "$(NO_COLOR)" ]; then ARGS="$$ARGS --no-color"; fi; \
	if [ -n "$(SPECS_DIR)" ]; then ARGS="$$ARGS --specs-dir \"$(SPECS_DIR)\""; fi; \
	if [ -n "$(IMPL_DIR)" ]; then ARGS="$$ARGS --impl-dir \"$(IMPL_DIR)\""; fi; \
	eval python3 ../../scripts/validate_go_signature_sync.py --impl-dir api/go --specs-dir docs/tech_specs $$ARGS

# Specification reference validation - validates Go file specification references
# Validates that all "Specification:" comments in Go files follow the format:
#   file_name.md: section_number heading_text
# Requires: Python 3
# Usage: make validate-go-spec-references [VERBOSE=1] [OUTPUT="file.txt"] [NO_FAIL=1] [NO_COLOR=1] [CHECK_INDEX=1] [REPO_ROOT="dir"]
#        - VERBOSE: Set to 1 for verbose output
#        - OUTPUT: Set to file path to save output to file (also prints to stdout)
#        - NO_FAIL: Set to 1 to exit with code 0 even if errors are found
#        - NO_COLOR: Set to 1 to disable colored output
#        - CHECK_INDEX: Set to 1 to also validate against api_go_defs_index.md (not yet implemented)
#        - REPO_ROOT: Repository root directory (default: parent of script directory)
validate-go-spec-references:
	@command -v python3 >/dev/null 2>&1 || { \
		echo "Error: python3 not found. Install Python 3 to run spec reference validation."; \
		exit 1; \
	}
	@echo "Validating specification references in Go files..."
	@ARGS=""; \
	if [ -n "$(VERBOSE)" ]; then ARGS="$$ARGS --verbose"; fi; \
	if [ -n "$(OUTPUT)" ]; then ARGS="$$ARGS --output $(OUTPUT)"; fi; \
	if [ -n "$(NO_FAIL)" ]; then ARGS="$$ARGS --no-fail"; fi; \
	if [ -n "$(NO_COLOR)" ]; then ARGS="$$ARGS --no-color"; fi; \
	if [ -n "$(CHECK_INDEX)" ]; then ARGS="$$ARGS --check-index"; fi; \
	if [ -n "$(REPO_ROOT)" ]; then ARGS="$$ARGS --repo-root \"$(REPO_ROOT)\""; fi; \
	python3 ../../scripts/validate_go_spec_references.py --repo-root ../.. $$ARGS

# Apply suggested spec reference updates from validate output.
# Run validate with OUTPUT first, then: make apply-go-spec-references INPUT=/path/to/output
# Usage: make apply-go-spec-references INPUT="file.txt"
apply-go-spec-references:
	@command -v python3 >/dev/null 2>&1 || { \
		echo "Error: python3 not found."; exit 1; \
	}
	@if [ -z "$(INPUT)" ]; then \
		echo "Error: INPUT=path is required. Example: make apply-go-spec-references INPUT=/tmp/spec_ref_errors.txt"; \
		exit 1; \
	fi
	@python3 ../../scripts/apply_go_spec_references.py --input "$(INPUT)" --repo-root ../..

# CI checks - performs same checks as GitHub Actions workflow
# NOTE: This target must be kept in sync with .github/workflows/go-ci.yml.
#       When modifying CI checks, update both this Makefile and the workflow file.
# Runs: dependency verification, unit tests, build, formatting check, go vet, golangci-lint, signature validation, and spec reference validation
ci:
	@echo "Running CI checks (same as GitHub Actions workflow)..."
	@echo ""
	@echo "1. Verifying dependencies..."
	@go mod verify
	@echo ""
	@echo "2. Running unit tests..."
	@/usr/bin/make test
	@echo ""
	@echo "3. Building code..."
	@if ( [ -w "$(HOME)/.cache/go-build" ] 2>/dev/null && touch "$(HOME)/.cache/go-build/.write-test" 2>/dev/null && rm -f "$(HOME)/.cache/go-build/.write-test" ) || [ -z "$(HOME)" ]; then \
		go build ./...; \
	else \
		CACHE_DIR="$$(cd ../.. && pwd)/tmp/go-cache"; \
		mkdir -p "$$CACHE_DIR"; \
		GOCACHE="$$CACHE_DIR" go build ./...; \
	fi
	@echo ""
	@echo "4. Checking formatting..."
	@if [ "$$(gofmt -s -l . | wc -l)" -gt 0 ]; then \
		echo "Code is not formatted. Run 'go fmt ./...'"; \
		gofmt -s -d .; \
		exit 1; \
	fi
	@echo ""
	@echo "5. Running go vet..."
	@if ( [ -w "$(HOME)/.cache/go-build" ] 2>/dev/null && touch "$(HOME)/.cache/go-build/.write-test" 2>/dev/null && rm -f "$(HOME)/.cache/go-build/.write-test" ) || [ -z "$(HOME)" ]; then \
		go vet ./...; \
	else \
		CACHE_DIR="$$(cd ../.. && pwd)/tmp/go-cache"; \
		mkdir -p "$$CACHE_DIR"; \
		GOCACHE="$$CACHE_DIR" go vet ./...; \
	fi
	@echo ""
	@echo "6. Running golangci-lint..."
	@if ( [ -w "$(HOME)/.cache/go-build" ] 2>/dev/null && touch "$(HOME)/.cache/go-build/.write-test" 2>/dev/null && rm -f "$(HOME)/.cache/go-build/.write-test" ) || [ -z "$(HOME)" ]; then \
		golangci-lint run ./...; \
		golangci-lint run --build-tags=bdd ./...; \
	else \
		CACHE_DIR="$$(cd ../.. && pwd)/tmp/go-cache"; \
		mkdir -p "$$CACHE_DIR"; \
		GOCACHE="$$CACHE_DIR" golangci-lint run ./... || { \
			echo "Error: golangci-lint failed. This may be due to Go module cache permissions."; \
			echo "Try running 'go mod tidy' or check Go module cache permissions."; \
			exit 1; \
		}; \
		GOCACHE="$$CACHE_DIR" golangci-lint run --build-tags=bdd ./... || { \
			echo "Error: golangci-lint (BDD) failed. This may be due to Go module cache permissions."; \
			echo "Try running 'go mod tidy' or check Go module cache permissions."; \
			exit 1; \
		}; \
	fi
	@echo ""
	@echo "7. Validating Go signatures against tech specs..."
	@make validate-go-signatures
	@echo ""
	@echo "8. Validating Go specification references..."
	@make validate-go-spec-references
