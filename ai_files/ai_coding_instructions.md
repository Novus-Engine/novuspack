# AI Coding Instructions for NovusPack Development

## Project Overview

**Project:** NovusPack - File Asset Package Management System
**Context:** Part of the larger novus-engine-go project
**Focus:** Development of novuspack library
**Methodology:** BDD/TDD with strict adherence to specifications

## 1. Core Principles

The following principles guide all development work on this project.

### 1.1 Specification Adherence

- **Primary Rule:** Follow specifications exactly as written
- **No Modifications:** AI cannot modify or correct specification documents
- **Halt on Issues:** Report specification problems to Product Owner, do not proceed
- **Source of Truth:** Specification documents in `docs/tech_specs/` are authoritative

### 1.2 Code Preservation

- **Golden Rule:** Preserve existing working code
- **Minimal Changes:** Make only necessary modifications
- **No Refactoring:** Avoid refactoring unless directly required
- **Regression Prevention:** All existing tests must continue to pass

### 1.3 BDD/TDD Workflow

- **Red-Green-Refactor:** Follow strict TDD methodology
- **Test-First:** Write failing tests before implementation
- **Incremental:** Small, focused changes with immediate feedback
- **Quality Gates:** All tests must pass before proceeding

### 1.4 No Lint Disabling

- **Do not disable linter checks.**
- Fix the underlying issue instead of suppressing it.
- Do not bypass linting in CI, Make targets, or tooling configuration.
- Do not add ignore directives (e.g., `# noqa`, `# pylint: disable=...`, `//nolint`) to silence failures unless explicitly instructed.

### 1.5 Shell Quoting with Backticks

- **Critical Rule:** Avoid passing markdown headings (especially ones containing backticks like `` `code` ``) as command-line string arguments.
- **Why:** Backticks are command substitution in shells, and quoting/escaping is error-prone for both humans and AI agents.
- **Fallback (when text must be on the command line):** Use single quotes (not double quotes) wherever possible and escape backticks in all other cases.
- **Correct Examples (preferred, file-based):**
  - `` make generate-anchor LINE='docs/tech_specs/api_core.md:224' ``
  - `` make generate-anchor FILE='docs/tech_specs/api_core.md' ``
- **Applies To:** Any tooling that needs to work with markdown headings containing backticks.

## 2. Development Workflow

This section outlines the development workflow from initial analysis through implementation and quality assurance.

### 2.1 Phase 1: Analysis & Planning

This phase involves setting up the development environment and analyzing the codebase and specifications.

#### 2.1.1 Environment Setup

- Verify Go 1.25.0+ installation
- Check PATH configuration
- Install required tools:
  - **Go tools:** gofmt, golangci-lint
  - **Markdown linting:** markdownlint-cli2 (`npm install -g markdownlint-cli2`)
  - **Python 3:** Required for validation and audit scripts
- Verify Make is available for build automation

#### 2.1.2 Repository Management

- Only work on the specified branch
  If not branch is specified, make a new branch to commit to
- Fetch latest changes and review specifications

#### 2.1.3 Code Analysis

- Map existing code to specifications
- Identify gaps and required changes
- Plan minimal necessary modifications

### Phase 2: BDD/TDD Implementation

#### Step 1: Functional Test (Red Phase)

- Write high-level functional test that defines the feature
- Test should initially fail (expected)
- Document BDD findings in `dev_docs/` directory
- Use naming convention: `YYYY-MM-DD_bdd_findings_feature_description.md`

#### Step 2: Unit Test Implementation (Inner TDD Loop)

- Write failing unit tests for specific components
- **Prefer updating existing tests** over creating new ones
- Focus on one component at a time
- Each test should fail initially

#### Step 3: Implementation (Green Phase)

- Write minimal code to make unit tests pass
- Preserve existing functionality
- Make only necessary changes
- Run tests after each change

#### Step 4: Refactoring (Refactor Phase)

- Clean up code while maintaining functionality
- Ensure all tests still pass
- No breaking changes to existing features

#### Step 5: Integration Testing

- Run functional test from Step 1
- Should now pass if implementation is correct
- If not, return to inner TDD loop

### Phase 3: Quality Assurance

#### Code Quality Checks

Use the provided Makefile targets for all quality checks:

- **Run linters:** `make lint` (Go + markdown + Python) or `make lint-go` (Go only)
  - Checks code formatting (gofmt)
  - Runs go vet
  - Runs golangci-lint (including BDD tests)
- **Fix all issues before committing**

#### Regression Testing

Use the provided Makefile targets for testing:

- **Run all unit tests:** `make test` or `make test-go`
- **Run all BDD tests:** `make bdd` or `make bdd-go`
- **Run BDD tests for specific domain (Go implementation):** `make -C api/go bdd-domain BDD_DOMAIN='@domain:xxx'`
  - Available domains: basic_ops, core, file_format, file_mgmt, file_types, compression, signatures, streaming, dedup, metadata, metadata_system, security_validation, security_encryption, generics, validation, testing, writing
- **Run BDD tests in CI mode:** `make bdd-ci`
- **Critical:** All tests must pass

#### Coverage Analysis

Use coverage targets to ensure adequate test coverage:

- **Generate coverage report:** `make coverage`
  - Creates coverage.out file
- **View HTML coverage report:** `make coverage-html`
  - Opens interactive HTML report in browser
- **View terminal coverage report:** `make coverage-report`
  - Shows coverage statistics in terminal

#### Documentation Quality

For complete documentation on all validation scripts, their options, and usage examples, see [scripts/README.md](../scripts/README.md).

Quick reference:

- **Run all documentation checks:** `make docs-check [PATHS="file1.md,dir1,file2.md"] [VERBOSE=1] [OUTPUT="file.txt"] [CHECK_COVERAGE=1]`
  - Runs all documentation validation checks in the correct order
  - Note: Some checks (signatures index, requirement references, coverage audits) are skipped when PATHS is specified, as they require checking all files
  - See [scripts/README.md](../scripts/README.md) for all available options
- **Validate markdown links:** `make validate-links [PATHS="file1.md,dir1"] [VERBOSE=1] [OUTPUT="file.txt"] [CHECK_COVERAGE=1]`
- **Validate Go code blocks:** `make validate-go-code-blocks [PATHS="file1.md,dir1"] [VERBOSE=1] [OUTPUT="file.txt"]`
- **Validate heading numbering:** `make validate-heading-numbering [PATHS="file1.md,dir1"] [VERBOSE=1] [OUTPUT="file.txt"]`
- **Apply heading corrections:** `make apply-heading-corrections [INPUT="file.txt"] [DRY_RUN=1] [VERBOSE=1]`
- **Validate Go definitions index:** `make validate-go-defs-index [VERBOSE=1]`
  - Note: This check is skipped when PATHS is specified, as it requires checking all tech specs
- **Validate requirement references:** `make validate-req-references [VERBOSE=1]`
  - Note: This check is skipped when PATHS is specified, as it requires checking all feature files
- **Lint markdown files:** `make lint-markdown [PATHS="file1.md,dir1"]`
- **Lint Python scripts:** `make lint-python [PATHS="path1,path2"]`

#### Coverage Audits

- **Audit feature coverage:** `make audit-feature-coverage [VERBOSE=1]`
  - Note: This check is skipped when PATHS is specified, as it requires checking all requirements and feature files
- **Audit requirements coverage:** `make audit-requirements-coverage [VERBOSE=1]`
  - Note: This check is skipped when PATHS is specified, as it requires checking all tech specs and requirements
- **Run all coverage audits:** `make audit-coverage`

See [scripts/README.md](../scripts/README.md) for detailed information about all options and usage examples.

#### Security Validation

- Verify quantum-safe cryptography implementations
- Check key management and storage
- Validate existing security measures remain intact

### Phase 4: Documentation & Version Control

#### Documentation Requirements

- **Location:** All development docs in `dev_docs/` directory only
- **Naming:** `YYYY-MM-DD_document_type_description.md`
- **Content:** Document decisions, analysis, and progress
- **Version Control:** Include in all commits

#### Commit Process

1. **Format code:** Execute gofmt on all Go files
2. **Run CI checks locally:** `make ci`
   - Runs documentation checks (Go code blocks, Go spec signature consistency, heading numbering, Go definitions index, requirement references, coverage audits, link validation, markdown linting)
   - Runs coverage audits
   - Runs Python linting
   - Verifies dependencies
   - Runs all unit tests
   - Builds code
   - Checks formatting
   - Runs static analysis (go vet)
   - Runs linters
3. **Fix any issues:** Address all failures before committing
4. **Commit:** Use conventional commit format
5. **Push:** Immediately push to development branch

#### CI Integration

The `make ci` target runs the exact same checks as GitHub Actions workflows:

- **Markdown Linting:** Checks markdown style and formatting
- **Go Code Blocks Validation:** Validates Go code blocks follow conventions
- **Heading Numbering Validation:** Validates markdown heading numbering consistency
- **Go Definitions Index Validation:** Validates all Go definitions in tech specs are present in the index
- **Link Validation:** Validates all internal documentation links
- **Requirement Reference Validation:** Validates REQ references in feature files
- **Coverage Audits:** Ensures specs have feature and requirement coverage
- **Python Linting:** Lints and audits Python scripts used by tooling
- **Go CI:** Runs dependency verification, tests, build, formatting, vet, and lint

**Important:** Always run `make ci` before pushing to ensure your changes will pass CI.

## Available Tooling

The repository provides comprehensive Make targets and Python scripts for all development tasks. **Always use Makefile targets instead of running commands directly.**

### Primary Make Targets

#### Complete CI Check

- **`make ci`** - Run all CI checks locally
  - Note: Keep `Makefile` CI targets in sync with GitHub Actions workflows (`docs-check.yml`, `python-lint.yml`, `go-ci.yml`).
  - Documentation checks (`make docs-check`)
    - Go code blocks validation
    - Go spec signature consistency validation
    - Heading numbering validation
    - Go definitions index validation
    - Requirement reference validation
    - Coverage audits (features and requirements)
    - Link validation
    - Markdown linting
  - Python linting (`make lint-python` or `make lint-python PATHS="path1,path2"`)
  - Go CI (`make ci-go`)
  - **Use this before every commit to ensure CI will pass**

#### Testing Targets

- **`make test`** - Run all unit tests
- **`make bdd`** - Run all BDD tests (output saved to tmp/ directory)
- **`make bdd-ci`** - Run BDD tests in CI mode with tag filtering (~@skip && ~@wip)
- **`make -C api/go bdd-domain BDD_DOMAIN='@domain:xxx'`** - Run BDD tests for specific domain (Go implementation)
  - Available domains: basic_ops, core, file_format, file_mgmt, file_types, compression, signatures, streaming, dedup, metadata, metadata_system, security_validation, security_encryption, generics, validation, testing, writing
  - Example: `make -C api/go bdd-domain BDD_DOMAIN='@domain:basic_ops'`

#### Coverage Targets

- **`make coverage`** - Generate test coverage report (creates coverage.out)
- **`make coverage-html`** - Generate and open HTML coverage report in browser
- **`make coverage-report`** - Display coverage statistics in terminal

#### Linting Targets

- **`make lint`** - Run all linters (Go + markdown + Python)
- **`make lint-go`** - Run Go linters only (gofmt, go vet, golangci-lint)
- **`make lint-markdown`** - Lint markdown files for style compliance
- **`make lint-python [PATHS="path1,path2"]`** - Lint Python scripts (flake8, pylint, xenon -b C, radon mi gate; plus non-gating radon/vulture/bandit)
  - PATHS: Comma-separated list of files/directories to check (default: scripts)
  - Fails if any block has cyclomatic complexity > C (via xenon -b C)
  - Fails if any module has maintainability index rank C (MI 0-9)

#### Documentation Validation Targets

For complete documentation on all validation scripts, their options, and usage examples, see [scripts/README.md](../scripts/README.md).

Quick reference:

- **`make docs-check [PATHS="file1.md,dir1,file2.md"] [VERBOSE=1] [OUTPUT="file.txt"] [CHECK_COVERAGE=1]`** - Run all documentation validation checks
- **`make markdown-lint [PATHS="file1.md,dir1,file2.md"]`** - Lint markdown files for style compliance
  - Supports PATHS for checking specific files or directories
- **`make validate-links [PATHS="file1.md,dir1"] [VERBOSE=1] [OUTPUT="file.txt"] [CHECK_COVERAGE=1]`** - Validate all internal markdown links and anchors
- **`make validate-go-code-blocks [PATHS="file1.md,dir1"] [VERBOSE=1] [OUTPUT="file.txt"]`** - Validate Go code blocks in tech specs
- **`make validate-heading-numbering [PATHS="file1.md,dir1"] [VERBOSE=1] [OUTPUT="file.txt"]`** - Validate markdown heading numbering
- **`make apply-heading-corrections [INPUT="file.txt"] [DRY_RUN=1] [VERBOSE=1]`** - Apply heading numbering corrections
- **`make validate-go-defs-index [VERBOSE=1]`** - Validate Go API definitions index
  - Note: Skipped when PATHS is specified (requires all tech specs)
- **`make validate-req-references [VERBOSE=1]`** - Validate requirement references
  - Note: Skipped when PATHS is specified (requires all feature files)
- **`make generate-anchor FILE='path/to/file.md'`** - Print anchors for all headings in a file
- **`make generate-anchor LINE='path/to/file.md:224'`** - Print anchor for the heading at a specific line in a file
  - See [Shell Quoting with Backticks](#15-shell-quoting-with-backticks) for details

#### Coverage Audit Targets

- **`make audit-feature-coverage [VERBOSE=1]`** - Check requirements have BDD feature coverage
  - Note: Skipped when PATHS is specified (requires all requirements and feature files)
- **`make audit-requirements-coverage [VERBOSE=1]`** - Check tech specs are referenced by requirements
  - Note: Skipped when PATHS is specified (requires all tech specs and requirements)
- **`make audit-coverage`** - Run both coverage audits

### Python Scripts

For complete documentation on all Python scripts, their options, and usage examples, see [scripts/README.md](../scripts/README.md).

**Note:** Always prefer Make targets over direct script execution.

### Project Manager

- **Primary Responsibility:** Specification compliance and code preservation
- **Key Tasks:**
  - Verify environment setup
  - Review implementation against specifications
  - Ensure minimal necessary changes
  - Monitor for regressions

### Functionality Tester

- **Primary Responsibility:** High-level feature validation
- **Key Tasks:**
  - Write functional tests that define features
  - Document BDD findings and analysis
  - Verify end-to-end functionality

### Unit Test Engineer

- **Primary Responsibility:** Component-level testing
- **Key Tasks:**
  - Write/update unit tests for components
  - Prefer updating existing tests over creating new ones
  - Ensure comprehensive test coverage

### Lead Developer

- **Primary Responsibility:** Code implementation
- **Key Tasks:**
  - Implement minimal code to pass tests
  - Preserve existing functionality
  - Format code and run quality checks
  - Commit and push changes

### Security Tester

- **Primary Responsibility:** Cryptographic validation
- **Key Tasks:**
  - Verify quantum-safe cryptography
  - Test for vulnerabilities
  - Validate key management

## Critical Requirements

### Environment Setup

- **Go Version:** 1.25.0+ (no exceptions)
- **PATH:** Go must be accessible from any directory
- **Tools:** gofmt, golangci-lint installed and working

### Code Quality

- **Formatting:** All Go files formatted (verified with `make lint`)
- **Static Analysis:** Pass all checks (verified with `make lint` and `make ci`)
- **Testing:** All tests must pass (run `make test` and `make bdd`)
- **Documentation:** Proper documentation in dev_docs/
- **Links:** All markdown links valid (verified with `make validate-links`)
- **Coverage:** Specs covered by features and requirements (verified with `make audit-coverage`)

### Version Control

- **Branch Strategy:** Work on `deepagent_dev`, source from `abacus_dev`
- **Commit Format:** Conventional commits with clear descriptions
- **Push Policy:** Immediate push after each commit
- **Documentation:** All dev docs committed with code

### File Organization

- **Test Packages:** Separate test packages following Go conventions
- **Documentation:** All dev docs in `dev_docs/` directory only
- **Naming:** Follow `YYYY-MM-DD_document_type_description.md` format

## Success Criteria

### Implementation Complete When

- All functional tests pass
- All unit tests pass
- All existing tests still pass (no regressions)
- Code is properly formatted
- Static analysis passes
- Documentation is complete
- Changes are committed and pushed

### Quality Gates

- **No failing tests** (zero tolerance)
- **No unformatted code** (gofmt required)
- **No static analysis issues** (clean linting)
- **No specification deviations** (exact compliance)
- **No breaking changes** (preserve existing functionality)

## Error Handling

### Specification Issues

- **Halt immediately** on specification problems
- **Report to Product Owner** - do not attempt corrections
- **Wait for clarification** before proceeding
- **No assumptions** or workarounds

### Test Failures

- **Fix immediately** - no failing tests acceptable
- **Identify root cause** of each failure
- **Re-run all tests** after fixes
- **Repeat until 100% pass rate**

### Code Quality Issues

- **Fix before committing** - no exceptions
- **Run all quality checks** before each commit
- **Verify formatting** with gofmt
- **Pass static analysis** before proceeding

## Best Practices

### Development Approach

- **Start with failing tests** (Red phase)
- **Implement minimal code** (Green phase)
- **Refactor carefully** (Refactor phase)
- **Preserve existing functionality** always
- **Make smallest possible changes** to achieve goals

### Documentation

- **Document everything** in dev_docs/
- **Use proper naming** conventions
- **Include analysis** and decisions
- **Version control** all documentation

### Testing Strategy

- **Update existing tests** when possible
- **Create new tests** only when necessary
- **Test incrementally** with each change
- **Verify no regressions** after each change

This document provides a streamlined, AI-optimized approach to the NovusPack development process while maintaining the core BDD/TDD methodology and strict adherence to specifications.
