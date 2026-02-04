# AI Coding Instructions for NovusPack Development

## Project Overview

**Project:** NovusPack - File Asset Package Management System
**Context:** Part of the larger novus-engine-go project
**Focus:** Development of novuspack library
**Methodology:** BDD/TDD with strict adherence to specifications

## Documentation Standards (Markdown)

- Markdown authoring conventions are defined in:
  - [`docs/docs_standards/markdown_conventions.md`](../docs/docs_standards/markdown_conventions.md) (single source of truth)
- Repository-wide AI workflow rules, tooling, and enforcement guidance are defined in:
  - [`.github/copilot-instructions.md`](../.github/copilot-instructions.md)

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

### 1.4 Repository-Wide Conventions

- Do not disable or bypass linters; fix issues instead.
- Prefer Make targets over direct script execution.
- Be careful with shell quoting text containing backticks; use single-quotes instead of double-quotes or avoid passing the text string via shell commands.
- For workflow conventions and guidance, follow [`.github/copilot-instructions.md`](../.github/copilot-instructions.md).
- The canonical command surface is the Makefiles:
  - [`Makefile`](../Makefile)
  - [`api/go/Makefile`](../api/go/Makefile)
  - [`cli/nvpkg/Makefile`](../cli/nvpkg/Makefile)

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

#### Default Gate (Before Committing or Pushing)

- Run `make ci`.
  This is the closest local equivalent to CI and should be treated as authoritative.

#### While Iterating (Pick What Matches Your Change)

- **Go code changes**:
  - `make lint-go`
  - `make test`
  - `make bdd` (or a narrower domain run while iterating)
- **Docs changes**:
  - `make docs-check PATHS=<path/to/file.md>`
  - If links or anchors changed, also run `make validate-links PATHS=<path/to/file.md>`
  - If you need a fast style check, run `make markdown-lint PATHS=<path/to/file.md>`
- **Tooling scripts (Python) changes**:
  - `make lint-python PATHS="scripts"`

#### Focused BDD Runs (Go Implementation)

- Run a domain-specific suite:
  - `make -C api/go bdd-domain BDD_DOMAIN='@domain:xxx'`
  - Example: `make -C api/go bdd-domain BDD_DOMAIN='@domain:core'`

#### Coverage (When Needed)

- `make coverage`
- `make coverage-html`
- `make coverage-report`

#### Where the Full Tooling Reference Lives

- Repository-wide tooling conventions and the canonical command list:
  [`.github/copilot-instructions.md`](../.github/copilot-instructions.md)
- Full documentation for validation scripts and options:
  [`scripts/README.md`](../scripts/README.md)

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
   - Treat this as the authoritative “everything” check (docs validation + lint + tests + build).
   - Details and the canonical list live in [`.github/copilot-instructions.md`](../.github/copilot-instructions.md).
3. **Fix any issues:** Address all failures before committing
4. **Commit:** Use conventional commit format
5. **Push:** Immediately push to development branch

#### CI Integration

The `make ci` target is intended to match GitHub Actions CI behavior.

**Important:** Always run `make ci` before pushing to ensure your changes will pass CI.

## Tooling Reference

This document intentionally does not duplicate the full tooling/command reference.

- Repository tooling and “approved” workflow conventions: [`.github/copilot-instructions.md`](../.github/copilot-instructions.md)
- Validation scripts documentation and options: [`scripts/README.md`](../scripts/README.md)

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
