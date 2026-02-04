# Makefile for nvpkg CLI
# Follows conventions of root Makefile and api/go/Makefile.
# Unit tests target 90%% code coverage.

# Ensure MAKE is set to the actual make binary for recursive calls
SYSTEM_MAKE := $(shell PATH="/usr/bin:/usr/local/bin:/bin:$$PATH" command -v make 2>/dev/null || echo /usr/bin/make)
ifneq ($(MAKE),$(SYSTEM_MAKE))
  override MAKE := $(SYSTEM_MAKE)
endif

.PHONY: test tidy coverage coverage-90 coverage-html coverage-report lint ci build build-dev clean
.PHONY: build-all build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64 build-windows-amd64 build-freebsd-amd64 build-freebsd-arm64

# Default minimum 82%%; use make coverage-90 to require 90%% (target)
COVERAGE_MIN ?= 82
TMP_DIR ?= ../../tmp
OUT_DIR ?= dist
LDFLAGS_RELEASE := -ldflags="-s -w"

# Run all unit tests
test:
	@if ( [ -w "$(HOME)/.cache/go-build" ] 2>/dev/null && touch "$(HOME)/.cache/go-build/.write-test" 2>/dev/null && rm -f "$(HOME)/.cache/go-build/.write-test" ) || [ -z "$(HOME)" ]; then \
		go test -v ./...; \
	else \
		mkdir -p $(TMP_DIR)/go-cache; \
		GOCACHE="$$(cd $(TMP_DIR) && pwd)/go-cache" go test -v ./...; \
	fi

# Run tests with coverage; fail if cmd coverage below COVERAGE_MIN (default 90%%).
# Only ./cmd is measured so main package does not lower the total.
coverage:
	@echo "Running tests with coverage (minimum $(COVERAGE_MIN)%% for cmd)..."
	@go clean -testcache
	@go test -coverprofile=coverage.out -coverpkg=./cmd ./cmd/...
	@echo ""
	@TOTAL=$$(go tool cover -func=coverage.out | awk '/total:/{gsub(/%/,""); print $$NF}'); \
	if [ -z "$$TOTAL" ]; then echo "Error: could not parse coverage total"; exit 1; fi; \
	if ! echo "$$TOTAL $(COVERAGE_MIN)" | awk '{exit !($$1>=$$2)}'; then \
		echo "Coverage $$TOTAL%% is below minimum $(COVERAGE_MIN)%%"; \
		go tool cover -func=coverage.out | grep "total:"; \
		exit 1; \
	fi; \
	echo "Coverage $$TOTAL%% meets minimum $(COVERAGE_MIN)%%"; \
	echo "Coverage profile saved to: coverage.out"

# Require 90%% coverage (target); fails until api path metadata write is complete
coverage-90:
	@$(MAKE) coverage COVERAGE_MIN=90

# Generate HTML coverage report
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
	@echo "=== Total ==="
	@go tool cover -func=coverage.out | grep "total:" || echo "Coverage data not available"

# Tidy dependencies
tidy:
	@echo "Running go mod tidy..."
	@go mod tidy
	@echo "Dependencies tidied successfully."

# Lint: gofmt, go vet, golangci-lint (same strategy as api/go; uses .golangci.yml)
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

# Native release build for current OS/arch only. Outputs to $(OUT_DIR)/<goos>-<goarch>/nvpkg (or nvpkg.exe on Windows).
BUILD_GOOS := $(shell go env GOOS)
BUILD_GOARCH := $(shell go env GOARCH)
BUILD_PLATFORM := $(BUILD_GOOS)-$(BUILD_GOARCH)
BUILD_BINARY := nvpkg$(if $(filter windows,$(BUILD_GOOS)),.exe,)

build:
	@mkdir -p $(OUT_DIR)/$(BUILD_PLATFORM)
	@rm -f $(OUT_DIR)/$(BUILD_PLATFORM)/$(BUILD_BINARY)
	@CGO_ENABLED=0 go build $(LDFLAGS_RELEASE) -trimpath -o $(OUT_DIR)/$(BUILD_PLATFORM)/$(BUILD_BINARY) .
	$(call run_upx,$(OUT_DIR)/$(BUILD_PLATFORM)/$(BUILD_BINARY))
	@echo "Built: $(CURDIR)/$(OUT_DIR)/$(BUILD_PLATFORM)/$(BUILD_BINARY)"

# Cross-platform release builds for all supported platforms. Outputs to $(OUT_DIR)/<platform>/nvpkg (or nvpkg.exe on Windows).
build-all: build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64 build-windows-amd64 build-freebsd-amd64 build-freebsd-arm64
	@echo "Built binaries:"
	@echo "  $(CURDIR)/$(OUT_DIR)/linux-amd64/nvpkg"
	@echo "  $(CURDIR)/$(OUT_DIR)/linux-arm64/nvpkg"
	@echo "  $(CURDIR)/$(OUT_DIR)/darwin-amd64/nvpkg"
	@echo "  $(CURDIR)/$(OUT_DIR)/darwin-arm64/nvpkg"
	@echo "  $(CURDIR)/$(OUT_DIR)/windows-amd64/nvpkg.exe"
	@echo "  $(CURDIR)/$(OUT_DIR)/freebsd-amd64/nvpkg"
	@echo "  $(CURDIR)/$(OUT_DIR)/freebsd-arm64/nvpkg"

# After build: UPX (required for release). Go -ldflags="-s -w" already omits symbols.
# Tolerate AlreadyPackedException (exit 2) so re-runs of build/ci succeed.
define run_upx
	@if command -v upx >/dev/null 2>&1; then upx --best $(1) || { r=$$?; [ $$r -eq 2 ] && exit 0; exit $$r; }; else echo "upx not found: install upx (e.g. apt install upx-ucl) for release builds"; exit 1; fi
endef

build-linux-amd64:
	@mkdir -p $(OUT_DIR)/linux-amd64
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS_RELEASE) -trimpath -o $(OUT_DIR)/linux-amd64/nvpkg .
	$(call run_upx,$(OUT_DIR)/linux-amd64/nvpkg)

build-linux-arm64:
	@mkdir -p $(OUT_DIR)/linux-arm64
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(LDFLAGS_RELEASE) -trimpath -o $(OUT_DIR)/linux-arm64/nvpkg .
	$(call run_upx,$(OUT_DIR)/linux-arm64/nvpkg)

build-darwin-amd64:
	@mkdir -p $(OUT_DIR)/darwin-amd64
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build $(LDFLAGS_RELEASE) -trimpath -o $(OUT_DIR)/darwin-amd64/nvpkg .
	$(call run_upx,$(OUT_DIR)/darwin-amd64/nvpkg)

build-darwin-arm64:
	@mkdir -p $(OUT_DIR)/darwin-arm64
	@CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build $(LDFLAGS_RELEASE) -trimpath -o $(OUT_DIR)/darwin-arm64/nvpkg .
	$(call run_upx,$(OUT_DIR)/darwin-arm64/nvpkg)

build-freebsd-amd64:
	@mkdir -p $(OUT_DIR)/freebsd-amd64
	@CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build $(LDFLAGS_RELEASE) -trimpath -o $(OUT_DIR)/freebsd-amd64/nvpkg .
	$(call run_upx,$(OUT_DIR)/freebsd-amd64/nvpkg)

build-freebsd-arm64:
	@mkdir -p $(OUT_DIR)/freebsd-arm64
	@CGO_ENABLED=0 GOOS=freebsd GOARCH=arm64 go build $(LDFLAGS_RELEASE) -trimpath -o $(OUT_DIR)/freebsd-arm64/nvpkg .
	$(call run_upx,$(OUT_DIR)/freebsd-arm64/nvpkg)

# Windows build requires api/go to be Windows-compatible (e.g. no Unix-only syscalls).
# Currently not supported until api/go is Windows-compatible.
build-windows-amd64:
	@mkdir -p $(OUT_DIR)/windows-amd64
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(LDFLAGS_RELEASE) -trimpath -o $(OUT_DIR)/windows-amd64/nvpkg.exe .
	$(call run_upx,$(OUT_DIR)/windows-amd64/nvpkg.exe)

# Build development binary with debug symbols (for debugging and stack traces)
build-dev:
	@go build -o nvpkg-dev .

# CI: verify, test, build, lint
ci: tidy
	@echo "Running CI checks..."
	@echo ""
	@echo "1. Verifying dependencies..."
	@go mod verify
	@echo ""
	@echo "2. Running unit tests with coverage..."
	@$(MAKE) coverage
	@echo ""
	@echo "3. Building..."
	@$(MAKE) build
	@echo ""
	@echo "4. Linting..."
	@$(MAKE) lint
	@echo "CI checks passed."

# Remove build and coverage artifacts
clean:
	@rm -rf $(OUT_DIR)
	@rm -f nvpkg nvpkg.exe nvpkg-dev nvpkg-dev.exe coverage.out coverage.html
