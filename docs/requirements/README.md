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
- [`transformation_pipeline.md`](transformation_pipeline.md) - Multi-stage transformation pipeline requirements
- [`validation.md`](validation.md) - Validation requirements
- [`writing.md`](writing.md) - Writing API requirements

## Conventions

### Basic Conventions

- Requirement IDs use `REQ-<DOMAIN>-NNN`.
- Requirements list `@spec(file#anchor)` sources.
- Functional requirements link to at least one feature scenario.
- See [Documentation-Only Requirements](#documentation-only-requirements) section for documentation-only requirement conventions.

### Domain Prefixes

Each requirements file has a single canonical `REQ-<DOMAIN>-NNN` prefix. The authoritative list of domains and their requirements files is maintained in [Documentation Standards](../docs_standards/requirements_domains.md). Prefix format: `REQ-<DOMAIN>-` where `<DOMAIN>` is the domain tag from that document.

### Disallowed or Legacy Prefixes

Do not introduce additional lookalike prefixes.
See [requirements_domains.md – Disallowed or legacy domain tags](../docs_standards/requirements_domains.md#disallowed-or-legacy-domain-tags) for the canonical list.

`traceability.md` is a matrix document and is not a requirements domain file.

### Multiple Anchor References

A single requirement can reference multiple headings from the same or different spec files.
This is useful when:

- Multiple headings describe aspects of the same requirement.
- A requirement covers related functionality across multiple sections.
- An existing requirement should be extended to cover an additional heading.

#### Format for Multiple Anchors

When a requirement references multiple headings, list them as comma-separated markdown links:

```markdown
- REQ-XXX-NNN: Description. [file#anchor1](../tech_specs/file.md#anchor1), [file#anchor2](../tech_specs/file.md#anchor2)
- REQ-XXX-NNN: Description. [file1#anchor1](../tech_specs/file1.md#anchor1), [file2#anchor2](../tech_specs/file2.md#anchor2)
```

#### When to Add Multiple Anchors

1. **Extending existing requirements**: If a valid requirement already exists for a heading, add the new heading reference to that existing requirement instead of creating a duplicate.

   - Example: If `REQ-API_BASIC-021` already covers "Package structure and loading", and a new heading "3. Package Structure and Loading" describes the same concept, add the new anchor to the existing requirement.

2. **Related functionality**: When multiple headings describe aspects of the same functional requirement.

3. **Avoiding duplication**: Before creating a new requirement, check if an existing requirement already covers the same concept and can be extended with an additional anchor reference.

### Requirement Type Classification

All requirements must be classified by type:

- **Functional** (default) - Observable, testable behavior.
  - Must have feature file with testable scenarios.
  - Type can be implied (no explicit tag needed for standard functional requirements).

- **Documentation-only** - Examples, guidance, best practices.

  - See [Documentation-Only Requirements](#documentation-only-requirements) section for details.

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
- REQ-XXX-NNN: Description [type: documentation-only] (documentation-only: examples - DO NOT CREATE FEATURE FILE). [spec#anchor]
- ~~REQ-XXX-NNN: Description~~ [type: obsolete] (obsolete: replaced by REQ-XXX-YYY - see [reference]). [spec#anchor]
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

- REQ-API_BASIC-001: NewPackage creates an empty, valid container in memory. [api_basic_operations.md#41-package-constructor](../tech_specs/api_basic_operations.md)
- REQ-API_BASIC-006: Create validates path and directory, configures package in memory. [api_basic_operations.md#42-create-method](../tech_specs/api_basic_operations.md)
- REQ-API_BASIC-032: Create example usage demonstrates creation usage [type: documentation-only] (documentation-only: examples - DO NOT CREATE FEATURE FILE). [api_basic_operations.md#424-create-example-usage](../tech_specs/api_basic_operations.md)

## Package Opening

- REQ-API_BASIC-002: OpenPackage validates format and returns structured errors. [api_basic_operations.md#51-openpackage](../tech_specs/api_basic_operations.md)
- REQ-API_BASIC-088: OpenPackageReadOnly opens an existing package and enforces read-only behavior. [api_basic_operations.md#52-openpackagereadonly](../tech_specs/api_basic_operations.md)
- REQ-API_BASIC-092: OpenPackageReadOnly enforces read-only behavior via a wrapper Package without duplicating OpenPackage parsing logic. [api_basic_operations.md#5212-go-v1-reference-implementation-sketch](../tech_specs/api_basic_operations.md)
- REQ-API_BASIC-091: OpenBrokenPackage opens a broken package for repair workflows. [api_basic_operations.md#53-openbrokenpackage](../tech_specs/api_basic_operations.md)

## Error Handling

- REQ-API_BASIC-016: Error handling returns structured errors for all failure cases. [api_basic_operations.md#8-error-handling](../tech_specs/api_basic_operations.md)
```

## Requirement Numbering Best Practices

### Finding Available Numbers

Before adding a new requirement, check for available gap numbers in the sequence:

```bash
make validate-req-references
```

This will show warnings for sequential numbering gaps (missing numbers in the sequence).
Use these gap numbers for new requirements to maintain sequential numbering.

### Numbering Guidelines

1. **Fill gaps first**: When adding new requirements, use the lowest available gap number from the warnings.

   - Example: If warnings show missing numbers `9, 44, 45`, use `9` for the next new requirement.

2. **Avoid duplicates**: Never reuse an existing requirement number, even if the requirement is obsoleted.

   - See [Obsoleted Requirements](#obsoleted-requirements) section for details on obsoleted requirements.
   - Always check for duplicates before committing: `make validate-req-references`

3. **Sequential numbering per category**: Each requirement category (e.g., `REQ-FILEMGMT-*`, `REQ-API_BASIC-*`) should have sequential numbering within each requirements file.

   - Gaps are acceptable but should be minimized.
   - Large gaps (50+ numbers) may indicate significant refactoring occurred.

4. **Validation before commit**: Always run `make validate-req-references` before committing requirement changes:

   - Ensures no duplicate requirements
   - Identifies available gap numbers via warnings
   - Validates format consistency
   - Checks feature file references

### Example Workflow

1. Check for gaps:

   ```bash
   make validate-req-references
   # Look for "Warnings: Sequential Numbering Gaps" section
   ```

2. Use the lowest available gap number for your new requirement.

3. Verify no duplicates:

   ```bash
   make validate-req-references
   # Should show "Duplicate requirements: 0"
   ```

4. Commit your changes.

## Obsoleted Requirements

Obsoleted requirements are requirements that have been replaced by newer requirements or are no longer applicable.
They are marked with strikethrough but remain in the requirements file for historical reference and traceability.

### Marking Obsoleted Requirements

Obsoleted requirements must be clearly marked:

- **Must use strikethrough**: `~~REQ-XXX-NNN: Description~~` (this distinguishes them from documentation-only requirements)
- Include `[type: obsolete]` classification
- Include explicit notation explaining why it's obsolete and what replaces it: `(obsolete: replaced by REQ-XXX-YYY - see [reference])`
- Optionally mark the section header as obsolete: `## ~~Section Name~~ (Obsolete - Use [Alternative])`

**Note**: The validation script uses the combination of strikethrough (`~~`) and `[type: obsolete]` to identify obsoleted requirements.

### Example: Obsoleted Requirements

```markdown
## ~~File Unstage Operations~~ (Obsolete - Use File Removal Operations)

- ~~REQ-FILEMGMT-002: Unstaging a file updates path metadata state or tombstones~~ [type: obsolete] (obsolete: replaced by RemoveFile operations - see REQ-FILEMGMT-136 through REQ-FILEMGMT-141).
- ~~REQ-FILEMGMT-011: UnstageFilePattern unstages files matching patterns~~ [type: obsolete] (obsolete: replaced by RemoveFilePattern operations - see REQ-FILEMGMT-325 through REQ-FILEMGMT-330).
```

### Handling Obsoleted Requirements

1. **Sequential numbering**: Obsoleted requirements remain in the sequential count.

   - They are NOT gaps in the numbering sequence.
   - They should NOT be reused for new requirements.
   - The validation script correctly includes them in gap detection (they are not flagged as gaps).
   - See [Requirement Numbering Best Practices](#requirement-numbering-best-practices) for numbering guidelines.

2. **Replacement references**: Always reference the replacement requirement(s) in the obsolete notation.

   - Use format: `(obsolete: replaced by REQ-XXX-YYY - see [reference])`
   - For multiple replacements, list all: `(obsolete: replaced by REQ-XXX-YYY through REQ-XXX-ZZZ)`

3. **Feature files**: Obsoleted requirements should not have active feature files.

   - Existing feature files referencing obsoleted requirements should be updated to reference the replacement requirements.
   - Remove or update `@REQ-XXX-NNN` tags in feature files when requirements become obsolete.
   - If a feature file only references obsoleted requirements (no valid remaining requirement references), the feature file should be deleted.

4. **Validation**: Run `make validate-req-references` before committing.

   - See [Requirement Numbering Best Practices](#requirement-numbering-best-practices) for validation guidelines.

## Documentation-Only Requirements

Some requirements are documentation-only (examples, guidance, best practices) and should NOT have feature files.

### Marking Documentation-Only Requirements

Documentation-only requirements must be clearly marked:

- Do NOT use strikethrough (unlike obsoleted requirements)
- Format: `- REQ-XXX-NNN: Description [type: documentation-only] ...`
- Include `[type: documentation-only]` classification
- Include explicit notation: `(documentation-only: [reason] - DO NOT CREATE FEATURE FILE)`
- Examples, guidance, and best practices belong in tech specs only.

### Handling Documentation-Only Requirements

1. **Sequential numbering**: Documentation-only requirements remain in the sequential count.

   - They are NOT gaps in the numbering sequence.
   - They should NOT be reused for new requirements.
   - See [Requirement Numbering Best Practices](#requirement-numbering-best-practices) for numbering guidelines.

2. **No feature files**: Documentation-only requirements must NOT have feature files.

   - The notation explicitly states "DO NOT CREATE FEATURE FILE".
   - Examples and guidance belong in tech specs, not in testable feature scenarios.

3. **Validation**: Run `make validate-req-references` before committing.

   - See [Requirement Numbering Best Practices](#requirement-numbering-best-practices) for validation guidelines.

### Example: Documentation-Only Requirements

```markdown
- REQ-API_BASIC-032: Create example usage demonstrates creation usage [type: documentation-only] (documentation-only: examples - DO NOT CREATE FEATURE FILE). [api_basic_operations.md#424-create-example-usage](../tech_specs/api_basic_operations.md)
```
