# Root-level Makefile for NovusPack multi-language repository
# This Makefile delegates to language-specific makefiles in api/<language>/<version>/
# Each language implementation has its own Makefile with version-specific targets.

.PHONY: test test-go test-go-v1 bdd bdd-go bdd-go-v1 bdd-ci bdd-lint test-unified

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

bdd-lint:
	$(MAKE) -C api/go/v1 bdd-lint

# Unified test runner (future implementation)
# This will run tests across all language implementations
test-unified:
	@echo "Unified test runner not yet implemented"
	# Future: $(MAKE) -C test-runner run
