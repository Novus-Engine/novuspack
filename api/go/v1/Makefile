# Makefile for NovusPack Go API v1
# This Makefile provides build and test targets for the Go v1 implementation.
# The BDD tests reference feature files from ../../../features/ (repository root).

.PHONY: test bdd bdd-ci bdd-lint

BDD_TAGS ?= '~@skip && ~@wip'

# Run all unit tests
test:
	go test -v ./...

# Run BDD tests using shared feature files from repository root
# Feature files are located at ../../../features/ relative to this directory
bdd:
	GOFLAGS=-tags=bdd go test -v -run "^TestFeatures$$" ./bdd

# Run BDD tests with tag filtering (for CI)
bdd-ci:
	GOFLAGS=-tags=bdd go test -v --godog.tags=$(BDD_TAGS) ./bdd

# Lint BDD feature files
bdd-lint:
	@../../scripts/bdd-lint.sh
