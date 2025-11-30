# AI Coding Instructions for NovusPack Development

## Project Overview

**Project:** NovusPack - File Asset Package Management System
**Context:** Part of the larger novus-engine-go project
**Focus:** Development of novuspack library
**Methodology:** BDD/TDD with strict adherence to specifications

## Core Principles

### 1. Specification Adherence

- **Primary Rule:** Follow specifications exactly as written
- **No Modifications:** AI cannot modify or correct specification documents
- **Halt on Issues:** Report specification problems to Product Owner, do not proceed
- **Source of Truth:** Specification documents in `docs/tech_specs/` are authoritative

### 2. Code Preservation

- **Golden Rule:** Preserve existing working code
- **Minimal Changes:** Make only necessary modifications
- **No Refactoring:** Avoid refactoring unless directly required
- **Regression Prevention:** All existing tests must continue to pass

### 3. BDD/TDD Workflow

- **Red-Green-Refactor:** Follow strict TDD methodology
- **Test-First:** Write failing tests before implementation
- **Incremental:** Small, focused changes with immediate feedback
- **Quality Gates:** All tests must pass before proceeding

## Development Workflow

### Phase 1: Analysis & Planning

1. **Environment Setup**

    - Verify Go 1.23.0+ installation
    - Check PATH configuration
    - Install required tools (gofmt, golangci-lint)

2. **Repository Management**

    - Only work on the specified branch
      If not branch is specified, make a new branch to commit to
    - Fetch latest changes and review specifications

3. **Code Analysis**

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

- **Format all Go files:** `find ./ -type f -iname '*.go' -exec gofmt -w {} \;`
- **Static analysis:** `go vet ./...`
- **Linting:** `golangci-lint run`
- **Fix all issues before committing**

#### Regression Testing

- Run ALL existing tests
- Verify no functionality was broken
- Fix any regressions immediately
- **Critical:** All tests must pass

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
2. **Run tests:** Ensure all tests pass
3. **Commit:** Use conventional commit format
4. **Push:** Immediately push to development branch

## Agent Roles

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

- **Go Version:** 1.23.0+ (no exceptions)
- **PATH:** Go must be accessible from any directory
- **Tools:** gofmt, golangci-lint installed and working

### Code Quality

- **Formatting:** All Go files formatted with gofmt before commits
- **Static Analysis:** Pass go vet and golangci-lint
- **Testing:** All tests must pass (unit and functional)
- **Documentation:** Proper documentation in dev_docs/

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
