# NovusPack Functional Requirements

This directory contains functional requirements derived from `../tech_specs/`.

**Note**: All links in requirements files use relative paths `../tech_specs/` to reference the technical specifications.

Each requirement is traceable to spec anchors and mapped to Gherkin features.

See `traceability.md` for end to end mapping.

## Structure

Requirements are organized by domain, with each domain having its own file:

- [`basic_ops.md`](basic_ops.md) - Basic operations API requirements
- [`compression.md`](compression.md) - Compression requirements
- [`core.md`](core.md) - Core API requirements
- [`dedup.md`](dedup.md) - Deduplication requirements
- [`file_format.md`](file_format.md) - File format requirements
- [`file_mgmt.md`](file_mgmt.md) - File management requirements
- [`file_types.md`](file_types.md) - File type system requirements
- [`generics.md`](generics.md) - Generic API requirements
- [`metadata.md`](metadata.md) - Metadata requirements
- [`metadata_system.md`](metadata_system.md) - Metadata system requirements
- [`security.md`](security.md) - Security requirements
- [`security_encryption.md`](security_encryption.md) - Security encryption requirements
- [`signatures.md`](signatures.md) - Digital signature requirements
- [`streaming.md`](streaming.md) - Streaming API requirements
- [`testing.md`](testing.md) - Testing requirements
- [`validation.md`](validation.md) - Validation requirements
- [`writing.md`](writing.md) - Writing API requirements

## Conventions

### Basic Conventions

- Requirement IDs use `REQ-<DOMAIN>-NNN`.
- Requirements list `@spec(file#anchor)` sources.
- Functional requirements link to at least one feature scenario.
- Documentation-only requirements are marked with strikethrough and explicit "DO NOT CREATE FEATURE FILE" notation.

### Requirement Type Classification

All requirements must be classified by type:

- **Functional** (default) - Observable, testable behavior.
  - Must have feature file with testable scenarios.
  - Type can be implied (no explicit tag needed for standard functional requirements).

- **Documentation-only** - Examples, guidance, best practices.
  - Marked with strikethrough: `~~REQ-XXX-NNN: Description~~`
  - Include explicit notation: `(documentation-only: [reason] - DO NOT CREATE FEATURE FILE)`
  - Examples, guidance, and best practices belong in tech specs only.
  - See Documentation-Only Requirements section below.

- **Non-functional** - Performance, security characteristics, quality attributes.
  - Must have measurable criteria.
  - May have feature files if testable aspects exist.

- **Architectural** - System design, structure, organization.
  - Typically documented in tech specs.
  - May have feature files if design decisions are testable.

- **Constraint** - Limitations, rules, restrictions.
  - Must be testable (constraint violations return errors).
  - Should have feature files testing constraint enforcement.

#### Format for Type Classification

```markdown
- REQ-XXX-NNN: Description. [spec#anchor]
- REQ-XXX-NNN: Description [type: non-functional]. [spec#anchor]
- ~~REQ-XXX-NNN: Description~~ [type: documentation-only] (documentation-only: examples - DO NOT CREATE FEATURE FILE). [spec#anchor]
```

### Organization and Grouping

Requirements within each domain file should be organized by feature or API method:

- Group related requirements under section headers (e.g., `## Package Creation`, `## Package Opening`).
- Use logical groupings that reflect the API structure or feature boundaries.
- For domains with few requirements, a flat list is acceptable.
- Cross-cutting concerns (error handling, validation) may be grouped separately.

#### Example Structure

```markdown
# Basic Operations API Requirements

## Package Creation

- REQ-API_BASIC-001: NewPackage creates an empty, valid container in memory. [api_basic_operations.md#41-package-constructor](../tech_specs/api_basic_operations.md#41-package-constructor)
- REQ-API_BASIC-006: Create validates path and directory, configures package in memory. [api_basic_operations.md#42-create-method](../tech_specs/api_basic_operations.md#42-create-method)
- ~~REQ-API_BASIC-032: Create example usage demonstrates creation usage~~ [type: documentation-only] (documentation-only: examples - DO NOT CREATE FEATURE FILE). [api_basic_operations.md#424-create-example-usage](../tech_specs/api_basic_operations.md#424-create-example-usage)

## Package Opening

- REQ-API_BASIC-002: OpenPackage validates format and returns structured errors. [api_basic_operations.md#51-open-method](../tech_specs/api_basic_operations.md#51-open-method)
- REQ-API_BASIC-009: OpenWithValidation opens package and performs full validation. [api_basic_operations.md#52-open-with-validation](../tech_specs/api_basic_operations.md#52-open-with-validation)

## Error Handling

- REQ-API_BASIC-016: Error handling returns structured errors for all failure cases. [api_basic_operations.md#8-error-handling](../tech_specs/api_basic_operations.md#8-error-handling)
```

## Documentation-Only Requirements

Some requirements are documentation-only (examples, guidance, best practices) and should NOT have feature files:

- Marked with strikethrough: `~~REQ-XXX-NNN: Description~~`
- Include explicit notation: `(documentation-only: [reason] - DO NOT CREATE FEATURE FILE)`
- Examples, guidance, and best practices belong in tech specs only.
- See `dev_docs/2025-11-07_documentation_only_requirements_policy.md--` for full policy.

### Example

```markdown
- ~~REQ-API_BASIC-032: Create example usage demonstrates creation usage~~ [type: documentation-only] (documentation-only: examples - DO NOT CREATE FEATURE FILE). [api_basic_operations.md#424-create-example-usage](../tech_specs/api_basic_operations.md#424-create-example-usage)
```

## Related Documents

- `dev_docs/2025-11-07_documentation_only_requirements_policy.md--` - Documentation-only requirements policy
- `dev_docs/2025-11-07_requirements_categorization_proposal.md` - Full categorization proposal and rationale
