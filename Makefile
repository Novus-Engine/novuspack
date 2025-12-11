# Root-level Makefile for NovusPack multi-language repository
# This Makefile delegates to language-specific makefiles in api/<language>/<version>/
# Each language implementation has its own Makefile with version-specific targets.

.PHONY: test test-go test-go-v1 bdd bdd-go bdd-go-v1 bdd-ci test-unified ci ci-go ci-go-v1 markdown-lint lint lint-go lint-go-v1

# Test targets - delegate to language implementations
test: test-go
test-go: test-go-v1
test-go-v1:
	$(MAKE) -C api/go/v1 test

# BDD test targets - delegate to language implementations
bdd: bdd-go
bdd-go: bdd-go-v1
bdd-go-v1:
	$(MAKE) -C api/go/v1 bdd

bdd-ci:
	$(MAKE) -C api/go/v1 bdd-ci

# CI check targets - perform same checks as GitHub Actions workflows
# NOTE: The 'ci' target in api/go/v1/Makefile must be kept in sync with
#       .github/workflows/go-ci.yml. When adding or modifying CI checks,
#       update both files to ensure local 'make ci' matches CI behavior.
ci: markdown-lint ci-go
ci-go: ci-go-v1
ci-go-v1:
	$(MAKE) -C api/go/v1 ci

# Markdown linting - performs same checks as GitHub Actions workflow
# NOTE: This target must be kept in sync with .github/workflows/markdown-lint.yml.
#       When adding or modifying markdown linting, update both this Makefile and
#       the workflow file to ensure local 'make markdown-lint' matches CI behavior.
#       Requires: npm install -g markdownlint-cli2
#       NOTE: This target must be kept in sync with .github/workflows/markdown-lint.yml.
#             The workflow runs: markdownlint-cli2 "**/*.md"
#             This relies on .markdownlint-cli2.jsonc to ignore tmp/** and other local dirs.
markdown-lint:
	@command -v markdownlint-cli2 >/dev/null 2>&1 || { \
		echo "Error: markdownlint-cli2 not found. Install with: npm install -g markdownlint-cli2"; \
		exit 1; \
	}
	@NODE_OPTIONS="--no-warnings=MODULE_TYPELESS_PACKAGE_JSON" markdownlint-cli2 "**/*.md"

# Unified test runner (future implementation)
# This will run tests across all language implementations
test-unified:
	@echo "Unified test runner not yet implemented"
	# Future: $(MAKE) -C test-runner run

lint: lint-go markdown-lint
lint-go: lint-go-v1
lint-go-v1:
	$(MAKE) -C api/go/v1 lint
