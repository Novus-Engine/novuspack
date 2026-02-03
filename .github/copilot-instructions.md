---
alwaysApply: true
---

# AI Instructions

## General Rules

- Always check existing files before making changes.
- When creating new files, use `touch` to create the file first, then edit it.
- Check the actual date using the `date` command before writing the date.
- **Review all README files for context** - see [README Files for Context](#readme-files-for-context) below.
- See [../ai_files/](../ai_files/) for AI assisted coding instructions.
- See [../docs/tech_specs/](../docs/tech_specs/) for technical specifications.
- For markdown files, always abide by the [markdown standards](#markdown-standards) below.
- **Use Make targets for all development tasks** - see [Available Tooling](#available-tooling) below.
  - Do NOT use direct script calls; use the `make` targets instead.
- Do not create or call scripts directly unless instructed to do so.
- **Do not disable linter checks.**
  - Do not bypass linting in CI, Make targets, or tooling configuration.
  - Do not add ignore directives (e.g., `# noqa`, `# pylint: disable=...`, `//nolint`) to silence failures unless explicitly instructed.
  - Fix the underlying issue instead of suppressing it.
- **Shell quoting with backticks:** Avoid passing markdown headings directly as command-line string arguments.
  Prefer file-based inputs to eliminate quoting issues (especially backticks).
  If you must pass text with backticks on the commandline, use single quotes (not double-quotes) wherever possible
  and properly escape the backticks in all other cases.
- Avoid using commands which require approval.
  - Commands that do not require approval:
    - awk
    - basename
    - cat
    - cd
    - comm
    - date
    - echo
    - git diff
    - git log
    - git show
    - git status
    - go build
    - go list
    - go test
    - go tool cover
    - grep
    - head
    - ls
    - make
    - pwd
    - realpath
    - sha256sum
    - sort
    - tail
    - test
    - uniq
    - wc
  - Unapproved commands:
    - sed (**use grep instead where ever possible**)
    - find
    - xargs
    - tee

## README Files for Context

Before starting work on any task, review the relevant README files for context and project conventions:

- **[README.md](../README.md)** - Project overview, setup, and main documentation
- **[docs/requirements/README.md](../docs/requirements/README.md)** - Requirements documentation standards, numbering conventions, and best practices
- **[scripts/README.md](../scripts/README.md)** - Documentation validation scripts, options, and usage examples
- **[api/go/README.md](../api/go/README.md)** - Go API implementation documentation and conventions
- **[api/go/_bdd/README.md](../api/go/_bdd/README.md)** - Go BDD test structure and conventions
- **[dev_docs/README.md](../dev_docs/README.md)** - Development documentation and working notes

These README files contain essential context about project structure, conventions, and workflows that should be understood before making changes.

## Markdown Standards

- Avoid pseudo-headings or heading-like lines.
  - Use proper Markdown headings instead of `^\*\*.*:\*\*$`, `^\*\*.*\*\*:$`, `^[0-9]+\. \*\*.*\*\*$` or similar.
- Use "=>" instead of "→".
- Avoid using non-ASCII characters, with the following exceptions:
  - The following may be used in work tracking documents (dev_docs): ✅, ❌, 📊, ⚠️
- Always include a blank line after any heading lines.
- Always include a blank line before and after a list.
- Put one sentence on a line, except in tables.
- Code blocks that are part of a list should be indented by four spaces.
- All headings within a document must be unique
- All headings must have proper numbering.
- Ensure there is no content duplication; use references instead.
- When referencing a file or path, make it a link.
- Whenever you create or modify a markdown file, run `make docs-check PATHS=<path/to/markdown.md>`
  - Fix errors and re-run the make command until no errors remain.

### Markdown Standards Exceptions

- The top level [README](../README.md) file should have useful emoji, badges, and other "nice" visual elements.

## Tech Specs Docs

- Use prose to describe the technical implementation.
- Tech specs docs should have minimal code aside from:
  - Function signatures
  - Constants definitions
  - Generic type definitions (type parameters, constraints)
  - Basic usage examples
- There must be only one source of truth for each specification.
  Any related specifications should refer to the singular source of truth and link to it instead of re-stating it.
- All error handling should be using the latest structured errors approach.
  References to legacy sentinel errors should be removed and replaced with structured errors.

## Tech Specs vs. Implementation

In any case where there are deviations between the technical specifications (docs/tech_specs) and the actual implementation,
the implementation MUST be brought into compliance with the tech specs.
In cases where the tech_specs are unclear would cause issues during implementation, STOP and ASK what you should do.

## Available Tooling

The repository provides comprehensive Make targets and Python scripts for all development tasks.
**Always use Makefile targets instead of running commands directly.**

### Primary Make Targets

#### Complete CI Check

- **`make ci`** - Run all CI checks locally (use before every commit)
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

#### Testing

- **`make test`** - Run all unit tests
- **`make bdd`** - Run all BDD tests (output saved to tmp/ directory)
- **`make bdd-ci`** - Run BDD tests in CI mode with tag filtering
- **`make -C api/go bdd-domain BDD_DOMAIN='@domain:xxx'`** - Run domain-specific BDD tests (Go implementation)
  - Available domains: basic_ops, core, file_format, file_mgmt, file_types, compression, signatures, streaming, dedup, metadata, metadata_system, security_validation, security_encryption, generics, validation, testing, writing

#### Coverage Analysis

- **`make coverage`** - Generate test coverage report
- **`make coverage-html`** - Generate and open HTML coverage report
- **`make coverage-report`** - Display coverage statistics in terminal

#### Linting and Validation

- **`make lint`** - Run all linters (Go + markdown + Python)
- **`make lint-go`** - Run Go linters only (gofmt, go vet, golangci-lint)
- **`make lint-markdown`** - Lint markdown files
- **`make lint-python [PATHS="path1,path2"]`** - Lint Python scripts (flake8, pylint, xenon -b C, radon mi gate; plus non-gating radon/vulture/bandit)
  - PATHS: Comma-separated list of files/directories to check (default: scripts)
  - Fails if any block has cyclomatic complexity > C (via xenon -b C)
  - Fails if any module has maintainability index rank C (MI 0-9)

#### Documentation Validation

For complete documentation on all validation scripts, their options, and usage examples, see [scripts/README.md](../scripts/README.md).

Quick reference:

- **`make docs-check [PATHS="file1.md,dir1,file2.md"] [VERBOSE=1] [OUTPUT="file.txt"] [CHECK_COVERAGE=1]`** - Run all documentation validation checks
  - Runs all checks in the correct order
  - Note: Some checks (signatures index, requirement references, coverage audits) are skipped when PATHS is specified, as they require checking all files
  - See [scripts/README.md](../scripts/README.md) for all available options
- **`make markdown-lint [PATHS="file1.md,dir1,file2.md"]`** - Lint markdown files for style compliance
  - Supports PATHS for checking specific files or directories
- **`make validate-links [PATHS="file1.md,dir1"] [VERBOSE=1] [OUTPUT="file.txt"] [CHECK_COVERAGE=1]`** - Validate all internal markdown links and anchors
- **`make validate-heading-numbering [PATHS="file1.md,dir1"] [VERBOSE=1] [OUTPUT="file.txt"]`** - Validate markdown heading numbering consistency
- **`make apply-heading-corrections [INPUT="file.txt"] [DRY_RUN=1] [VERBOSE=1]`** - Apply heading numbering corrections
- **`make validate-go-code-blocks [PATHS="file1.md,dir1"] [VERBOSE=1] [OUTPUT="file.txt"]`** - Validate Go code blocks in tech specs follow conventions
- **`make validate-go-defs-index [VERBOSE=1]`** - Validate that all Go definitions in tech specs are in the Go API definitions index
  - Note: Skipped when PATHS is specified (requires all tech specs to validate the index)
- **`make validate-req-references [VERBOSE=1]`** - Validate REQ references in feature files
  - Note: Skipped when PATHS is specified (requires all feature files to validate references)
- **`make generate-anchor FILE='path/to/file.md'`** - Print anchors for all headings in a file
- **`make generate-anchor LINE='path/to/file.md:224'`** - Print anchor for the heading at a specific line in a file

#### Coverage Audits

- **`make audit-feature-coverage [VERBOSE=1]`** - Check requirements have BDD feature coverage
  - Note: Skipped when PATHS is specified (requires all tech specs and feature files)
- **`make audit-requirements-coverage [VERBOSE=1]`** - Check tech specs referenced by requirements
  - Note: Skipped when PATHS is specified (requires all tech specs and requirements)
- **`make audit-coverage`** - Run both coverage audits

See [scripts/README.md](../scripts/README.md) for detailed information about all options and usage examples.

### Approved Make Commands

All make commands are pre-approved for use:

- `make ci` - Complete CI check
- `make test` - Unit tests
- `make bdd` - BDD tests
- `make bdd-ci` - BDD tests (CI mode)
- `make -C api/go bdd-domain BDD_DOMAIN='@domain:xxx'` - Domain-specific BDD tests (Go implementation)
- `make lint` - All linters
- `make lint-go` - Go linters
- `make markdown-lint` - Markdown linter
- `make lint-python [PATHS="path1,path2"]` - Python linter (default: scripts)
- `make coverage` - Coverage report
- `make coverage-html` - HTML coverage report
- `make coverage-report` - Terminal coverage report
- `make validate-links` - Link validation
- `make validate-go-code-blocks` - Go code blocks validation
- `make validate-heading-numbering` - Heading numbering validation
- `make validate-go-defs-index` - Go definitions index validation
- `make audit-feature-coverage` - Feature coverage audit
- `make audit-requirements-coverage` - Requirements coverage audit
- `make audit-coverage` - All coverage audits
- `make generate-anchor FILE='path/to/file.md'` - Print anchors for all headings in a file
- `make generate-anchor LINE='path/to/file.md:224'` - Print anchor for a heading at a specific line
