# Root-level Makefile for NovusPack multi-language repository
# This Makefile delegates to language-specific makefiles in api/<language>/<version>/
# Each language implementation has its own Makefile with version-specific targets.

# Ensure MAKE is set to the actual make binary for recursive calls
# Check if MAKE matches the system make (using restricted PATH to avoid AppImage interference)
# If not, force it to the system make
SYSTEM_MAKE := $(shell PATH="/usr/bin:/usr/local/bin:/bin:$$PATH" command -v make 2>/dev/null || echo /usr/bin/make)
ifneq ($(MAKE),$(SYSTEM_MAKE))
  # MAKE doesn't match system make, force it
  override MAKE := $(SYSTEM_MAKE)
endif

.PHONY: test test-go test-go-v1 test-nvpkg bdd bdd-go bdd-go-v1 bdd-ci test-unified ci ci-go ci-go-v1 ci-nvpkg venv lint-markdown lint-python lint-nvpkg validate-links validate-heading-numbering apply-heading-corrections generate-anchor validate-req-references validate-go-defs-index validate-go-code-blocks validate-go-spec-signature-consistency validate-go-signatures validate-go-signatures-go validate-go-signatures-go-v1 validate-go-spec-references apply-go-spec-references audit-feature-coverage audit-requirements-coverage audit-coverage docs-check lint lint-go lint-go-v1 coverage coverage-go coverage-go-v1 coverage-nvpkg coverage-html coverage-html-go coverage-html-go-v1 coverage-html-nvpkg coverage-report coverage-report-go coverage-report-go-v1 coverage-report-nvpkg build build-nvpkg build-dev-nvpkg clean clean-go clean-nvpkg

# Test targets - delegate to language implementations
test: test-go test-nvpkg
test-go: test-go-v1
test-go-v1:
	$(MAKE) -C api/go test

test-nvpkg:
	$(MAKE) -C cli/nvpkg test

# BDD test targets - delegate to language implementations
bdd: bdd-go
bdd-go: bdd-go-v1
bdd-go-v1:
	$(MAKE) -C api/go bdd

bdd-ci:
	$(MAKE) -C api/go bdd-ci

# CI check targets - perform same checks as GitHub Actions workflows
# NOTE: The 'ci' target must be kept in sync with GitHub Actions workflows.
#       When adding or modifying CI checks, update both this Makefile and the
#       corresponding workflow files to ensure local 'make ci' matches CI behavior.
#       Current workflows: docs-check.yml, python-lint.yml, go-ci.yml, nvpkg-ci.yml.
# NOTE: Order matters - validate-heading-numbering runs before validate-links because
#       fixing heading numbering changes anchor IDs, which would break link validation.
ci: docs-check lint-python ci-go ci-nvpkg
ci-go: ci-go-v1
ci-go-v1:
	/usr/bin/make -C api/go ci

ci-nvpkg:
	$(MAKE) -C cli/nvpkg ci

# Python venv for lint tooling - creates .venv and installs scripts/requirements-lint.txt
# Run once (or after adding/updating scripts/requirements-lint.txt) so make lint-python uses the venv.
# Usage: make venv
venv:
	@command -v python3 >/dev/null 2>&1 || { \
		echo "Error: python3 not found. Install Python 3 to create the venv."; \
		exit 1; \
	}
	@python3 -m venv .venv
	@.venv/bin/pip install -q --upgrade pip
	@.venv/bin/pip install -q -r scripts/requirements-lint.txt
	@echo "Created .venv with lint tooling. Use 'make lint-python' (it will use .venv when present)."

# Markdown linting - performs same checks as GitHub Actions workflow
# NOTE: This target must be kept in sync with .github/workflows/docs-check.yml.
#       When adding or modifying markdown linting, update both this Makefile and
#       the workflow file to ensure local 'make lint-markdown' matches CI behavior.
#       Requires: npm install -g markdownlint-cli2
#       NOTE: The workflow runs: markdownlint-cli2 "**/*.md"
#             Excludes are configured in .markdownlint-cli2.jsonc (e.g. tmp/).
# Usage: make lint-markdown [PATHS="file1.md,dir1,file2.md"] [MDL_CONFIG="path/to/.markdownlint-cli2.jsonc"]
#        - PATHS: Comma-separated list of files/directories to check (default: all .md files)
#        - MDL_CONFIG: Optional config file path to pass via --config
lint-markdown:
	@command -v markdownlint-cli2 >/dev/null 2>&1 || { \
		echo "Error: markdownlint-cli2 not found. Install with: npm install -g markdownlint-cli2"; \
		exit 1; \
	}
	@if [ -n "$(PATHS)" ]; then \
		LINT_PATHS=$$(echo "$(PATHS)" | tr ',' ' '); \
		if [ -n "$(MDL_CONFIG)" ]; then \
			NODE_OPTIONS="--no-warnings=MODULE_TYPELESS_PACKAGE_JSON" markdownlint-cli2 --config "$(MDL_CONFIG)" $$LINT_PATHS; \
		else \
			NODE_OPTIONS="--no-warnings=MODULE_TYPELESS_PACKAGE_JSON" markdownlint-cli2 $$LINT_PATHS; \
		fi; \
	else \
		if [ -n "$(MDL_CONFIG)" ]; then \
			NODE_OPTIONS="--no-warnings=MODULE_TYPELESS_PACKAGE_JSON" markdownlint-cli2 --config "$(MDL_CONFIG)" "**/*.md"; \
		else \
			NODE_OPTIONS="--no-warnings=MODULE_TYPELESS_PACKAGE_JSON" markdownlint-cli2 "**/*.md"; \
		fi; \
	fi

# Python linting - performs same checks as GitHub Actions workflow
# NOTE: This target must be kept in sync with .github/workflows/python-lint.yml.
#       When adding or modifying Python linting, update both this Makefile and
#       the workflow file to ensure local 'make lint-python' matches CI behavior.
#       Requires: pip install flake8 pylint radon xenon vulture bandit
#       NOTE: This target includes gate-style linting (flake8, pylint, xenon -b C,
#             radon mi fail on rank C) and code smell tooling (radon, vulture, bandit).
#       Xenon (radon-based) fails if any block has cyclomatic complexity > C.
#       Radon mi fails if any module has maintainability index rank C (MI 0-9).
# Usage: make lint-python [PATHS="path1,path2,path3"]
#        - PATHS: Comma-separated list of files/directories to check (default: scripts)
lint-python:
	@command -v python3 >/dev/null 2>&1 || { \
		echo "Error: python3 not found. Install Python 3 to run Python linting."; \
		exit 1; \
	}
	@command -v flake8 >/dev/null 2>&1 || [ -x .venv/bin/flake8 ] || { \
		echo "Error: flake8 not found. Install with: pip install flake8 or run 'make venv'"; \
		exit 1; \
	}
	@command -v pylint >/dev/null 2>&1 || [ -x .venv/bin/pylint ] || { \
		echo "Error: pylint not found. Install with: pip install pylint or run 'make venv'"; \
		exit 1; \
	}
	@command -v radon >/dev/null 2>&1 || [ -x .venv/bin/radon ] || { \
		echo "Error: radon not found. Install with: pip install radon or run 'make venv'"; \
		exit 1; \
	}
	@command -v xenon >/dev/null 2>&1 || [ -x .venv/bin/xenon ] || { \
		echo "Error: xenon not found. Install with: pip install xenon or run 'make venv'"; \
		exit 1; \
	}
	@command -v vulture >/dev/null 2>&1 || [ -x .venv/bin/vulture ] || { \
		echo "Error: vulture not found. Install with: pip install vulture or run 'make venv'"; \
		exit 1; \
	}
	@command -v bandit >/dev/null 2>&1 || [ -x .venv/bin/bandit ] || { \
		echo "Error: bandit not found. Install with: pip install bandit or run 'make venv'"; \
		exit 1; \
	}
	@if [ -n "$(PATHS)" ]; then \
		LINT_PATHS=$$(echo "$(PATHS)" | tr ',' ' '); \
	else \
		LINT_PATHS="scripts"; \
	fi; \
	if [ -d .venv ]; then PATH="$(CURDIR)/.venv/bin:$$PATH"; export PATH; fi; \
	export PYTHONPATH="$(CURDIR)/scripts"; \
	echo "Running flake8 on Python scripts..."; \
	flake8 $$LINT_PATHS --jobs=1; FLAKE8_RESULT=$$?; \
	echo "Running pylint on Python scripts..."; \
	pylint --rcfile=.pylintrc $$LINT_PATHS; PYLINT_RESULT=$$?; \
	echo "Running radon complexity (non-gating)..."; \
	radon cc -s -a $$LINT_PATHS || true; \
	echo "Running xenon cyclomatic complexity check (fail if any block > C)..."; \
	xenon -b C $$LINT_PATHS; XENON_RESULT=$$?; \
	echo "Running radon maintainability metrics (non-gating)..."; \
	radon mi -s $$LINT_PATHS || true; \
	echo "Running radon maintainability check (fail if any module MI rank C)..."; \
	TMP_MI=$$(mktemp); \
	radon mi -j $$LINT_PATHS -O $$TMP_MI; \
	python3 -c "import sys, json; d=json.load(open(sys.argv[1])); bad=[k for k,v in d.items() if v.get('rank')=='C']; [print('MI rank C (low maintainability):', f) for f in bad]; sys.exit(1 if bad else 0)" $$TMP_MI; \
	MI_RESULT=$$?; rm -f $$TMP_MI; \
	echo "Running vulture unused code detection (non-gating)..."; \
	vulture $$LINT_PATHS --min-confidence 80 || true; \
	echo "Running bandit security scan (non-gating)..."; \
	bandit -r $$LINT_PATHS; BANDIT_RESULT=$$?; \
	echo ""; echo "Lint exit codes: flake8=$$FLAKE8_RESULT pylint=$$PYLINT_RESULT xenon=$$XENON_RESULT radon_mi=$$MI_RESULT bandit=$$BANDIT_RESULT"; \
	[ $$FLAKE8_RESULT -ne 0 ] || [ $$PYLINT_RESULT -ne 0 ] || [ $$XENON_RESULT -ne 0 ] || [ $$MI_RESULT -ne 0 ] || [ $$BANDIT_RESULT -ne 0 ] && exit 1; exit 0

# Link validation - validates all internal markdown links and anchors
# NOTE: This target must be kept in sync with .github/workflows/docs-check.yml.
#       When adding or modifying link validation, update both this Makefile and
#       the workflow file to ensure local 'make validate-links' matches CI behavior.
#       Requires: Python 3
# Usage: make validate-links [PATHS="file1.md,dir1,file2.md"] [VERBOSE=1] [OUTPUT="file.txt"] [CHECK_COVERAGE=1] [NO_FAIL=1] [NO_COLOR=1]
#        - PATHS: Comma-separated list of files/directories to check
#        - VERBOSE: Set to 1 for verbose output
#        - OUTPUT: Path to output file for detailed report
#        - CHECK_COVERAGE: Set to 1 to check that all requirements reference tech specs
#        - NO_FAIL: Set to 1 to exit with code 0 even if errors are found
#        - NO_COLOR: Set to 1 to disable colored output
validate-links:
	@command -v python3 >/dev/null 2>&1 || { \
		echo "Error: python3 not found. Install Python 3 to run link validation."; \
		exit 1; \
	}
	@ARGS=""; \
	if [ -n "$(PATHS)" ]; then ARGS="$$ARGS --path \"$(PATHS)\""; fi; \
	if [ -n "$(VERBOSE)" ]; then ARGS="$$ARGS --verbose"; fi; \
	if [ -n "$(OUTPUT)" ]; then ARGS="$$ARGS --output \"$(OUTPUT)\""; fi; \
	if [ -n "$(CHECK_COVERAGE)" ]; then ARGS="$$ARGS --check-coverage"; fi; \
	if [ -n "$(NO_FAIL)" ]; then ARGS="$$ARGS --no-fail"; fi; \
	if [ -n "$(NO_COLOR)" ]; then ARGS="$$ARGS --no-color"; fi; \
	eval python3 scripts/validate_links.py $$ARGS

# Heading numbering validation - validates markdown heading numbering consistency
# NOTE: This target must be kept in sync with .github/workflows/docs-check.yml.
#       When adding or modifying heading numbering validation, update both this Makefile and
#       the workflow file to ensure local 'make validate-heading-numbering' matches CI behavior.
#       Requires: Python 3
# Usage: make validate-heading-numbering [PATHS="file1.md,dir1,file2.md"] [VERBOSE=1] [OUTPUT="file.txt"] [NO_FAIL=1] [NO_COLOR=1]
#        - PATHS: Comma-separated list of files/directories to check
#        - VERBOSE: Set to 1 for verbose output
#        - OUTPUT: Path to output file for detailed report
#        - NO_FAIL: Set to 1 to exit with code 0 even if errors are found
#        - NO_COLOR: Set to 1 to disable colored output
validate-heading-numbering:
	@command -v python3 >/dev/null 2>&1 || { \
		echo "Error: python3 not found. Install Python 3 to run heading numbering validation."; \
		exit 1; \
	}
	@ARGS=""; \
	if [ -n "$(PATHS)" ]; then ARGS="$$ARGS --path \"$(PATHS)\""; fi; \
	if [ -n "$(VERBOSE)" ]; then ARGS="$$ARGS --verbose"; fi; \
	if [ -n "$(OUTPUT)" ]; then ARGS="$$ARGS --output \"$(OUTPUT)\""; fi; \
	if [ -n "$(NO_FAIL)" ]; then ARGS="$$ARGS --no-fail"; fi; \
	if [ -n "$(NO_COLOR)" ]; then ARGS="$$ARGS --no-color"; fi; \
	eval python3 scripts/validate_heading_numbering.py $$ARGS

# Apply heading corrections - automatically applies corrections from validate-heading-numbering output
# NOTE: This target is a convenience wrapper for apply_heading_corrections.py.
#       Requires: Python 3
# Usage: make apply-heading-corrections [INPUT="tmp/report.txt"] [DRY_RUN=1] [VERBOSE=1]
#        - INPUT: Path to validation output file (default: reads from stdin)
#        - DRY_RUN: Set to 1 for dry-run mode (shows changes without applying)
#        - VERBOSE: Set to 1 for verbose output
apply-heading-corrections:
	@command -v python3 >/dev/null 2>&1 || { \
		echo "Error: python3 not found. Install Python 3 to apply heading corrections."; \
		exit 1; \
	}
	@ARGS=""; \
	if [ -n "$(INPUT)" ]; then ARGS="$$ARGS --input \"$(INPUT)\""; fi; \
	if [ -n "$(DRY_RUN)" ]; then ARGS="$$ARGS --dry-run"; fi; \
	if [ -n "$(VERBOSE)" ]; then ARGS="$$ARGS --verbose"; fi; \
	eval python3 scripts/apply_heading_corrections.py $$ARGS

# Generate markdown anchors from markdown headings
# NOTE: This is a utility wrapper for scripts/generate_anchor.py.
#       Useful for creating links to specific sections in markdown files.
#       Requires: Python 3
# Usage: make generate-anchor FILE="path/to/file.md"
#        make generate-anchor LINE="path/to/file.md:224"
#        - FILE: Print anchors for all headings in the file
#        - LINE: Print anchor for the heading at a specific line in the file
generate-anchor:
	@command -v python3 >/dev/null 2>&1 || { \
		echo "Error: python3 not found. Install Python 3 to generate anchor."; \
		exit 1; \
	}
	@if [ -z "$(FILE)" ] && [ -z "$(LINE)" ]; then \
		echo "Error: FILE or LINE is required."; \
		echo ""; \
		echo "Usage: make generate-anchor FILE=\"path/to/file.md\""; \
		echo "       make generate-anchor LINE=\"path/to/file.md:224\""; \
		exit 1; \
	fi
	@if [ -n "$(FILE)" ] && [ -n "$(LINE)" ]; then \
		echo "Error: FILE and LINE are mutually exclusive. Provide only one."; \
		echo ""; \
		echo "Usage: make generate-anchor FILE=\"path/to/file.md\""; \
		echo "       make generate-anchor LINE=\"path/to/file.md:224\""; \
		exit 1; \
	fi
	@if [ -n "$(LINE)" ]; then \
		python3 scripts/generate_anchor.py --line "$(LINE)"; \
	else \
		python3 scripts/generate_anchor.py --file "$(FILE)"; \
	fi

# Requirement reference validation - validates REQ references in feature files
# NOTE: This target must be kept in sync with .github/workflows/docs-check.yml.
#       When adding or modifying requirement reference validation, update both this Makefile and
#       the workflow file to ensure local 'make validate-req-references' matches CI behavior.
#       Requires: Python 3
#       NOTE: This script requires checking all feature files to validate references properly,
#             so it is skipped when PATHS is specified (use without PATHS to run this check).
# Usage: make validate-req-references [VERBOSE=1] [NO_FAIL=1] [NO_COLOR=1]
#        - VERBOSE: Set to 1 for verbose output
#        - NO_FAIL: Set to 1 to exit with code 0 even if errors are found
#        - NO_COLOR: Set to 1 to disable colored output
validate-req-references:
	@if [ -n "$(PATHS)" ]; then \
		echo "Skipping requirement reference validation (requires checking all feature files, not individual paths)"; \
	else \
		command -v python3 >/dev/null 2>&1 || { \
			echo "Error: python3 not found. Install Python 3 to run requirement reference validation."; \
			exit 1; \
		}; \
		ARGS=""; \
		if [ -n "$(VERBOSE)" ]; then ARGS="$$ARGS --verbose"; fi; \
		if [ -n "$(NO_FAIL)" ]; then ARGS="$$ARGS --no-fail"; fi; \
		if [ -n "$(NO_COLOR)" ]; then ARGS="$$ARGS --no-color"; fi; \
		eval python3 scripts/validate_req_references.py $$ARGS; \
	fi

# Go definitions index validation - validates that all Go definitions in tech specs are in the index
# NOTE: This target must be kept in sync with .github/workflows/docs-check.yml.
#       This should run after validate-heading-numbering but before validate-links.
#       Requires: Python 3
#       NOTE: This script requires checking all tech specs to validate the index, so it is
#             skipped when PATHS is specified (use without PATHS to run this check).
# Usage: make validate-go-defs-index [VERBOSE=1] [OUTPUT="file.txt"] [NO_FAIL=1] [NO_COLOR=1] [APPLY=1]
#        - VERBOSE: Set to 1 for verbose output
#        - OUTPUT: Path to output file for detailed report
#        - NO_FAIL: Set to 1 to exit with code 0 even if errors are found
#        - NO_COLOR: Set to 1 to disable colored output
#        - APPLY: Set to 1 to apply high-confidence index updates (interactive)
validate-go-defs-index:
	@if [ -n "$(PATHS)" ]; then \
		echo "Skipping Go definitions index validation (requires checking all tech specs, not individual paths)"; \
	else \
		command -v python3 >/dev/null 2>&1 || { \
			echo "Error: python3 not found. Install Python 3 to run Go definitions index validation."; \
			exit 1; \
		}; \
		ARGS=""; \
		if [ -n "$(VERBOSE)" ]; then ARGS="$$ARGS --verbose"; fi; \
		if [ -n "$(OUTPUT)" ]; then ARGS="$$ARGS --output \"$(OUTPUT)\""; fi; \
		if [ -n "$(NO_FAIL)" ]; then ARGS="$$ARGS --no-fail"; fi; \
		if [ -n "$(NO_COLOR)" ]; then ARGS="$$ARGS --no-color"; fi; \
		if [ -n "$(APPLY)" ]; then ARGS="$$ARGS --apply"; fi; \
		eval python3 scripts/validate_api_go_defs_index.py $$ARGS; \
	fi

# Go code blocks validation - validates Go code blocks in tech specs follow conventions
# NOTE: This target must be kept in sync with .github/workflows/docs-check.yml.
#       Requires: Python 3
# Usage: make validate-go-code-blocks [PATHS="file1.md,dir1,file2.md"] [VERBOSE=1] [OUTPUT="file.txt"] [NO_FAIL=1] [NO_COLOR=1]
#        - PATHS: Comma-separated list of files/directories to check
#        - VERBOSE: Set to 1 for verbose output
#        - OUTPUT: Path to output file for detailed report
#        - NO_FAIL: Set to 1 to exit with code 0 even if errors are found
#        - NO_COLOR: Set to 1 to disable colored output
validate-go-code-blocks:
	@command -v python3 >/dev/null 2>&1 || { \
		echo "Error: python3 not found. Install Python 3 to run Go code blocks validation."; \
		exit 1; \
	}
	@ARGS=""; \
	if [ -n "$(PATHS)" ]; then ARGS="$$ARGS --path \"$(PATHS)\""; fi; \
	if [ -n "$(VERBOSE)" ]; then ARGS="$$ARGS --verbose"; fi; \
	if [ -n "$(OUTPUT)" ]; then ARGS="$$ARGS --output \"$(OUTPUT)\""; fi; \
	if [ -n "$(NO_FAIL)" ]; then ARGS="$$ARGS --no-fail"; fi; \
	if [ -n "$(NO_COLOR)" ]; then ARGS="$$ARGS --no-color"; fi; \
	eval python3 scripts/validate_go_code_blocks.py $$ARGS

# Go signature consistency validation - validates signature consistency within tech specs
# NOTE: This target must be kept in sync with .github/workflows/docs-check.yml.
#       Requires: Python 3
# Usage: make validate-go-spec-signature-consistency [PATHS="file1.md,dir1,file2.md"] [VERBOSE=1] [OUTPUT="file.txt"] [NO_FAIL=1] [NO_COLOR=1]
#        - PATHS: Comma-separated list of files/directories to check
#        - VERBOSE: Set to 1 for verbose output
#        - OUTPUT: Path to output file for detailed report
#        - NO_FAIL: Set to 1 to exit with code 0 even if errors are found
#        - NO_COLOR: Set to 1 to disable colored output
validate-go-spec-signature-consistency:
	@command -v python3 >/dev/null 2>&1 || { \
		echo "Error: python3 not found. Install Python 3 to run signature consistency validation."; \
		exit 1; \
	}
	@ARGS=""; \
	if [ -n "$(PATHS)" ]; then ARGS="$$ARGS --path \"$(PATHS)\""; fi; \
	if [ -n "$(VERBOSE)" ]; then ARGS="$$ARGS --verbose"; fi; \
	if [ -n "$(OUTPUT)" ]; then ARGS="$$ARGS --output \"$(OUTPUT)\""; fi; \
	if [ -n "$(NO_FAIL)" ]; then ARGS="$$ARGS --no-fail"; fi; \
	if [ -n "$(NO_COLOR)" ]; then ARGS="$$ARGS --no-color"; fi; \
	eval python3 scripts/validate_go_spec_signature_consistency.py $$ARGS

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
validate-go-signatures: validate-go-signatures-go
validate-go-signatures-go: validate-go-signatures-go-v1
validate-go-signatures-go-v1:
	$(MAKE) -C api/go validate-go-signatures VERBOSE="$(VERBOSE)" OUTPUT="$(OUTPUT)" NO_FAIL="$(NO_FAIL)" NO_COLOR="$(NO_COLOR)" SPECS_DIR="$(SPECS_DIR)" IMPL_DIR="$(IMPL_DIR)"

# Go spec reference validation - validates Specification: comments in Go files
# NOTE: This target must be kept in sync with api/go/Makefile and go-ci workflow.
# Usage: make validate-go-spec-references [VERBOSE=1] [OUTPUT="file.txt"] [NO_FAIL=1] [NO_COLOR=1] [CHECK_INDEX=1] [REPO_ROOT="dir"]
validate-go-spec-references:
	$(MAKE) -C api/go validate-go-spec-references VERBOSE="$(VERBOSE)" OUTPUT="$(OUTPUT)" NO_FAIL="$(NO_FAIL)" NO_COLOR="$(NO_COLOR)" CHECK_INDEX="$(CHECK_INDEX)" REPO_ROOT="$(REPO_ROOT)"

# Apply spec reference updates from validate output
# Usage: make apply-go-spec-references INPUT="file.txt"
apply-go-spec-references:
	$(MAKE) -C api/go apply-go-spec-references INPUT="$(INPUT)"

# Feature coverage audit - checks requirements are covered by feature files
# NOTE: This target must be kept in sync with .github/workflows/docs-check.yml.
#       When adding or modifying coverage audits, update both this Makefile and
#       the workflow file to ensure local 'make audit-feature-coverage' matches CI behavior.
#       Requires: Python 3
#       NOTE: This script requires checking all requirements and feature files to validate coverage,
#             so it is skipped when PATHS is specified (use without PATHS to run this check).
# Usage: make audit-feature-coverage [VERBOSE=1] [NO_FAIL=1] [NO_COLOR=1]
#        - VERBOSE: Set to 1 for verbose output
#        - NO_FAIL: Set to 1 to exit with code 0 even if errors are found
#        - NO_COLOR: Set to 1 to disable colored output
audit-feature-coverage:
	@if [ -n "$(PATHS)" ]; then \
		echo "Skipping feature coverage audit (requires checking all requirements and feature files, not individual paths)"; \
	else \
		command -v python3 >/dev/null 2>&1 || { \
			echo "Error: python3 not found. Install Python 3 to run coverage audit."; \
			exit 1; \
		}; \
		ARGS=""; \
		if [ -n "$(VERBOSE)" ]; then ARGS="$$ARGS --verbose"; fi; \
		if [ -n "$(NO_FAIL)" ]; then ARGS="$$ARGS --no-fail"; fi; \
		if [ -n "$(NO_COLOR)" ]; then ARGS="$$ARGS --no-color"; fi; \
		eval python3 scripts/audit_feature_coverage.py $$ARGS; \
	fi

# Requirements coverage audit - checks tech specs are referenced by requirements
# NOTE: This target must be kept in sync with .github/workflows/docs-check.yml.
#       When adding or modifying coverage audits, update both this Makefile and
#       the workflow file to ensure local 'make audit-requirements-coverage' matches CI behavior.
#       Requires: Python 3
#       NOTE: This script requires checking all tech specs and requirements to validate coverage,
#             so it is skipped when PATHS is specified (use without PATHS to run this check).
# Usage: make audit-requirements-coverage [VERBOSE=1] [NO_FAIL=1] [NO_COLOR=1]
#        - VERBOSE: Set to 1 for verbose output
#        - NO_FAIL: Set to 1 to exit with code 0 even if errors are found
#        - NO_COLOR: Set to 1 to disable colored output
audit-requirements-coverage:
	@if [ -n "$(PATHS)" ]; then \
		echo "Skipping requirements coverage audit (requires checking all tech specs and requirements, not individual paths)"; \
	else \
		command -v python3 >/dev/null 2>&1 || { \
			echo "Error: python3 not found. Install Python 3 to run coverage audit."; \
			exit 1; \
		}; \
		ARGS=""; \
		if [ -n "$(VERBOSE)" ]; then ARGS="$$ARGS --verbose"; fi; \
		if [ -n "$(NO_FAIL)" ]; then ARGS="$$ARGS --no-fail"; fi; \
		if [ -n "$(NO_COLOR)" ]; then ARGS="$$ARGS --no-color"; fi; \
		eval PYTHONPATH=scripts python3 scripts/audit_requirements_coverage.py $$ARGS; \
	fi

# Combined coverage audit - runs both feature and requirements coverage audits
# NOTE: This target must be kept in sync with .github/workflows/docs-check.yml.
#       Requires: Python 3
audit-coverage: audit-requirements-coverage audit-feature-coverage

# Documentation checks - runs all markdown and documentation-related validation checks
# NOTE: This target must be kept in sync with .github/workflows/docs-check.yml.
#       Order matters:
#       - validate-go-code-blocks runs first to check code block structure
#       - validate-go-spec-signature-consistency runs after code blocks to check signature consistency
#       - validate-heading-numbering runs before validate-go-defs-index and validate-links
#         because fixing heading numbering changes anchor IDs
#       - validate-go-defs-index runs before validate-links to catch missing definitions early
#       This target groups all documentation checks for convenience.
# Usage: make docs-check [PATHS="file1.md,dir1,file2.md"] [VERBOSE=1] [OUTPUT="file.txt"] [CHECK_COVERAGE=1] [NO_FAIL=1] [NO_COLOR=1]
#        - PATHS: Comma-separated list of files/directories to check (applies to all validation scripts)
#        - VERBOSE: Set to 1 for verbose output (applies to all validation scripts)
#        - OUTPUT: Path to output file for detailed report (applies to validate-links, validate-heading-numbering, validate-go-code-blocks)
#        - CHECK_COVERAGE: Set to 1 to check that all requirements reference tech specs (applies to validate-links only)
#        - NO_FAIL: Set to 1 to exit with code 0 even if errors are found (applies to all validation scripts)
#        - NO_COLOR: Set to 1 to disable colored output (applies to all validation scripts)
docs-check: validate-go-code-blocks validate-go-spec-signature-consistency validate-heading-numbering validate-go-defs-index validate-req-references audit-coverage validate-links lint-markdown

# Unified test runner (future implementation)
# This will run tests across all language implementations
test-unified:
	@echo "Unified test runner not yet implemented"
	# Future: $(MAKE) -C test-runner run

lint: lint-go lint-nvpkg lint-markdown lint-python
lint-go: lint-go-v1
lint-go-v1:
	$(MAKE) -C api/go lint

lint-nvpkg:
	$(MAKE) -C cli/nvpkg lint

# Coverage targets - delegate to language implementations
coverage: coverage-go coverage-nvpkg
coverage-go: coverage-go-v1
coverage-go-v1:
	$(MAKE) -C api/go coverage

coverage-nvpkg:
	$(MAKE) -C cli/nvpkg coverage

# HTML coverage report targets - delegate to language implementations
coverage-html: coverage-html-go coverage-html-nvpkg
coverage-html-go: coverage-html-go-v1
coverage-html-go-v1:
	$(MAKE) -C api/go coverage-html

coverage-html-nvpkg:
	$(MAKE) -C cli/nvpkg coverage-html

# Coverage report targets - delegate to language implementations
coverage-report: coverage-report-go coverage-report-nvpkg
coverage-report-go: coverage-report-go-v1
coverage-report-go-v1:
	$(MAKE) -C api/go coverage-report

coverage-report-nvpkg:
	$(MAKE) -C cli/nvpkg coverage-report

# Build nvpkg CLI (minimal-size release binary)
# Build nvpkg for current OS/arch (delegates to cli/nvpkg).
build: build-nvpkg

build-nvpkg:
	$(MAKE) -C cli/nvpkg build

# Build nvpkg development binary (with debug symbols, outputs nvpkg-dev)
build-dev-nvpkg:
	$(MAKE) -C cli/nvpkg build-dev

# Remove build and coverage artifacts (delegates to all sub-makefiles).
clean: clean-go clean-nvpkg

clean-go:
	$(MAKE) -C api/go clean

clean-nvpkg:
	$(MAKE) -C cli/nvpkg clean
