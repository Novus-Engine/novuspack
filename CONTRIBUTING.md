# Contributing to NovusPack

Thank you for your interest in contributing to NovusPack.

## 0. Overview

This repository is specification-driven.
The technical specifications in [`docs/tech_specs/`](docs/tech_specs/) are the source of truth.
The shared Gherkin feature files in [`features/`](features/) describe expected behavior for all implementations.

If you are implementing behavior, start from the relevant technical specification and the corresponding feature files.

## 1. Repository Structure

This repository is organized to support multiple language implementations while sharing a single specification and test surface.

- **Specifications**: [`docs/tech_specs/`](docs/tech_specs/)
- **Requirements**: [`docs/requirements/`](docs/requirements/)
- **Shared BDD feature files**: [`features/`](features/)
- **Language implementations**: `api/<language>/<version>/`
  - Go v1 implementation: [`api/go/`](api/go/)

## 2. Development Prerequisites

- **Go**: Go 1.25.0 or later.
- **Python**: Python 3.x (used by documentation tooling).
- **Node**: `markdownlint-cli2` installed globally for markdown linting.
- **Make**: Required to run the standard workflows.

## 3. AI-Assisted Development

This repository is designed to facilitate AI-assisted development with comprehensive validation tooling that provides clear, actionable feedback to help AI agents be self-correcting.

This does however put particular emphasis on making sure the repository technical specifications, feature files, and requirements references are correct and up-to-date; a spec change can have cascading effects that requires everything to be in-sync for all validations to pass.

### 3.1 Self-Correcting Feedback System

The repository provides a comprehensive validation and feedback system through Make targets and Python validation scripts that:

- **Clearly identify issues**: All validation scripts provide specific error messages with file paths and line numbers.
- **Call out specification violations**: Scripts explicitly validate against technical specifications and report identified deviations (namely definition signature differences).
- **Provide actionable feedback**: Error messages include examples of correct formats and suggestions for fixes.
- **Enable iterative improvement**: AI agents can run validation, fix issues, and re-validate in a feedback loop using the make targets (see [Makefile](./Makefile)).

### 3.2 Make Targets for Validation

The repository provides unified Make targets that wrap Python validation scripts, making it easy for AI agents to:

- **Run comprehensive checks**: `make ci` runs all validation checks locally, matching CI behavior.
- **Get specific feedback**: Individual targets like `make validate-links` or `make validate-heading-numbering` provide focused validation.
- **Target specific files**: Most targets support `PATHS` parameter for validating only modified files.

Key validation targets:

- `make ci` - Run all CI checks (recommended before every commit).
- `make docs-check` - Validate all documentation (markdown linting, heading numbering, links, code blocks, etc.).
- `make validate-go-signatures` - Validate Go implementation matches technical specifications (from `api/go/` directory).
- `make lint` - Run all linters (Go and markdown).
- `make test` - Run all unit tests.
- `make bdd` - Run BDD tests.

See [Makefile](./Makefile) for more details.

### 3.3 Python Validation Scripts

The `scripts/` directory contains Python validation and audit scripts that provide detailed feedback:

- **Specification compliance**: Scripts validate code and documentation against technical specifications.
- **Clear error reporting**: All scripts provide detailed error messages with file paths, line numbers, and specific issues.
- **Actionable suggestions**: Error messages include examples of correct formats and how to fix issues.
- **Coverage validation**: Audit scripts ensure specifications are covered by tests and requirements.

For complete documentation on all validation scripts, see [`scripts/README.md`](scripts/README.md).

### 3.4 How AI Agents Can Use This System

1. **Make changes**: AI agents can make code or documentation changes.
2. **Run validation**: Execute `make ci` or specific validation targets to check for issues.
3. **Review feedback**: Validation scripts provide clear error messages identifying what needs to be fixed.
4. **Fix issues**: Use the specific feedback to make corrections.
5. **Re-validate**: Run validation again to confirm fixes.
6. **Iterate**: Repeat until all validations pass.

This feedback loop enables AI agents to self-correct by providing clear, actionable information about what needs to be fixed, rather than requiring human intervention for every issue.

## 4. Working on Code

This section covers guidelines for working with code in the repository.

### 4.1. Where to Start

When implementing new functionality or making changes:

1. **Read the relevant technical specification** in [`docs/tech_specs/`](docs/tech_specs/) - this is the source of truth for implementation details.
2. **Read the relevant requirements** in [`docs/requirements/`](docs/requirements/) to understand the functional requirements.
3. **Find the corresponding feature coverage** in [`features/`](features/) using `@spec(...)` tags that reference the tech spec.
4. **Prefer small, focused changes** that keep existing behavior working.

Important workflow for changes:

- **If tech specs need updating**: Update the technical specifications first (they are the source of truth), then ensure requirements reference the updated specs.
  - **Note**: The validation scripts expect type, function, and method signatures to be in the tech specs for comparison and will flag deviations from the tech specs in the implementations.
- **If requirements need updating**: Update the requirements documentation, ensuring they reference the correct tech specs using relative links.
- **Update feature files**: Ensure feature files use correct `@spec(...)` tags to reference tech specs and `@REQ-*` tags to reference requirements.
- **Validate references**: Run `make validate-req-references` and `make validate-links` to ensure all references are correct.
- **Then implement**: Only after specifications and references are correct should you implement the code changes.

### 4.2 Tests

Run unit tests:

```bash
make test
```

Run BDD tests:

```bash
make bdd
```

Run BDD tests for a single domain:

```bash
make bdd-domain BDD_DOMAIN='@domain:file_mgmt'
```

## 5. Documentation Workflows

The documentation set has a canonical entry point.
See [`docs/tech_specs/_main.md`](docs/tech_specs/_main.md).

### 5.1. Markdown Standards

Documentation is validated in CI.
Follow the markdown standards described in [`.github/copilot-instructions.md`](.github/copilot-instructions.md).

### 5.2. Documentation Validation

Prefer using `make` targets rather than calling scripts directly.
For full details on the validators, see [`scripts/README.md`](scripts/README.md).

Run the full local CI suite (recommended before every commit):

```bash
make ci
```

Run only documentation checks:

```bash
make docs-check
```

Target specific files or directories using `PATHS`:

```bash
make docs-check PATHS="README.md,docs/tech_specs"
```

## 6. Contribution Guidelines

This section outlines guidelines for contributing to the project.

### 6.1. Keep Specs As the Source of Truth

Avoid duplicating specification content in multiple places.
If you need to summarize, link to the canonical source-of-truth document instead.

### 6.2. Error Handling

Prefer the repository's structured error approach.
If you change error behavior, update specs and tests accordingly.

### 6.3. External Library Usage

NovusPack follows a policy of minimizing external library dependencies across all language implementations.
When external libraries are necessary, make sure the external library is using a compatible open-source license (see below).

#### 6.3.1. Minimizing Dependencies

- **Prefer standard library**: Use language standard libraries whenever possible.
- **Justify external libraries**: External libraries should only be used when they provide essential functionality that cannot be reasonably implemented using standard libraries.
- **Evaluate alternatives**: Before adding a dependency, consider if the functionality can be achieved with existing dependencies or standard library features.
- **Assess library health**: Ensure the library is well maintained, actively developed, and not deprecated before adding it as a dependency.

#### 6.3.2. License Compatibility Requirements

All external libraries and their entire dependency chain must use compatible open-source licenses.
This requirement applies to all language implementations (Go, Rust, Zig, and future languages).

**Compatible licenses** (most common):

- **MIT License**: Fully compatible, very permissive.
- **Apache License 2.0**: Compatible, includes patent grants.
- **BSD 2-Clause and 3-Clause**: Compatible, very permissive.
- **LGPL (v2.1/v3)**: Compatible for library use, allows linking from proprietary code.
- **MPL-2.0 (Mozilla Public License 2.0)**: Compatible, weak copyleft (file-level).
  MPL-2.0 files must remain under MPL-2.0, but can be combined with other licensed code in a larger work.

**Incompatible licenses** (must be avoided):

- **GPL v2**: Incompatible with Apache 2.0 in both directions; cannot be used.
  Apache 2.0 code cannot be included in GPL v2 projects, and GPL v2 code cannot be incorporated into Apache 2.0 projects.
- **GPL v3**: Incompatible with Apache 2.0 for incorporation into NovusPack.
  Apache 2.0 code cannot incorporate GPL v3 code (though Apache 2.0 code can be used in GPL v3 projects).
- **AGPL**: Incompatible; cannot be used.
- **Proprietary licenses**: Cannot be used.

#### 6.3.3. Dependency Chain Verification

When adding an external library, you must verify:

1. **Direct dependency license**: The library itself must have a compatible license.
2. **Transitive dependencies**: All dependencies in the entire dependency chain must also have compatible licenses.
3. **No external dependencies preferred**: Libraries with no external dependencies (like `samber/lo` in the Go implementation) are preferred.
4. **Library is well maintained/heathy**: Avoid using deprecated libraries or libraries where maintenance is poor or unclear.

**Example**: The Go implementation uses `samber/lo` which:

- Has an MIT license (compatible).
- Has no external dependencies (ideal).
- Is well maintained as of January 2026.

#### 6.3.4. Adding a New Dependency

Before adding a new external library:

1. Verify the library's license is in the compatible list above.
2. Check the library's dependency chain using language-specific tools:
   - **Go**: Use `go list -m all` to see all dependencies, then verify each license.
   - **Rust**: Use `cargo tree` and verify licenses in `Cargo.toml` files.
   - **Zig**: Check package manifests and verify licenses.
3. Assess library maintenance and health:
   - Check recent commit activity and release frequency.
   - Verify the library is not deprecated or in maintenance-only mode.
   - Review issue tracker for signs of active maintenance.
   - Ensure the library has a clear maintenance status and roadmap.
4. Document the justification in your pull request, including maintenance assessment.
5. Update the relevant `README.md` or dependency documentation if needed.

#### 6.3.5. License Verification Tools

Use language-specific tools to verify licenses:

- **Go**: `go-licenses` or manual verification via `go list -m all`.
- **Rust**: `cargo-deny` or `cargo-license`.
- **Zig**: Manual verification of package manifests.

**Note**: When evaluating dependencies, also assess their maintenance status using repository activity metrics, recent releases, and issue tracker responsiveness.

### 6.4 Commits

Run `make ci` before committing.
Use clear, conventional commit messages that explain intent.

## 7. Pull Requests

Keep pull requests scoped.
Include a short summary and test plan in the PR description.
If you changed docs, ensure `make docs-check` passes locally.
