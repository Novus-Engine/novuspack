# NovusPack Scripts

- [Overview](#overview)
- [Available Scripts](#available-scripts)
  - [generate\_anchor.py](#generate_anchorpy)
    - [Anchor Generation Purpose](#anchor-generation-purpose)
    - [Anchor Generation Usage](#anchor-generation-usage)
    - [Anchor Generation Features](#anchor-generation-features)
    - [Anchor Generation Integration](#anchor-generation-integration)
    - [Anchor Generation Requirements](#anchor-generation-requirements)
  - [validate\_links.py](#validate_linkspy)
    - [Link Validation Purpose](#link-validation-purpose)
    - [Link Validation Usage](#link-validation-usage)
    - [Link Validation Features](#link-validation-features)
    - [Link Validation Exit Codes](#link-validation-exit-codes)
    - [Link Validation Integration](#link-validation-integration)
    - [Link Validation Requirements](#link-validation-requirements)
  - [audit\_feature\_coverage.py](#audit_feature_coveragepy)
    - [Feature Coverage Audit Purpose](#feature-coverage-audit-purpose)
    - [Feature Coverage Audit Usage](#feature-coverage-audit-usage)
    - [Feature Coverage Audit Features](#feature-coverage-audit-features)
    - [Feature Coverage Audit Exit Codes](#feature-coverage-audit-exit-codes)
    - [Feature Coverage Audit Integration](#feature-coverage-audit-integration)
    - [Feature Coverage Audit Requirements](#feature-coverage-audit-requirements)
  - [validate\_heading\_numbering.py](#validate_heading_numberingpy)
    - [Heading Numbering Validation Purpose](#heading-numbering-validation-purpose)
    - [Heading Numbering Validation Usage](#heading-numbering-validation-usage)
    - [Heading Numbering Validation Features](#heading-numbering-validation-features)
    - [Heading Numbering Validation Exit Codes](#heading-numbering-validation-exit-codes)
    - [Heading Numbering Validation Integration](#heading-numbering-validation-integration)
    - [Heading Numbering Validation Requirements](#heading-numbering-validation-requirements)
    - [Heading Numbering Validation Output Example](#heading-numbering-validation-output-example)
  - [apply\_heading\_corrections.py](#apply_heading_correctionspy)
    - [Heading Corrections Application Purpose](#heading-corrections-application-purpose)
    - [Heading Corrections Application Usage](#heading-corrections-application-usage)
    - [Heading Corrections Application Features](#heading-corrections-application-features)
    - [Heading Corrections Application Exit Codes](#heading-corrections-application-exit-codes)
    - [Heading Corrections Application Requirements](#heading-corrections-application-requirements)
    - [Heading Corrections Application Output Example](#heading-corrections-application-output-example)
    - [Heading Corrections Application Integration](#heading-corrections-application-integration)
    - [Heading Corrections Application Workflow](#heading-corrections-application-workflow)
  - [validate\_req\_references.py](#validate_req_referencespy)
    - [Requirement Reference Validation Purpose](#requirement-reference-validation-purpose)
    - [Requirement Reference Validation Usage](#requirement-reference-validation-usage)
    - [Requirement Reference Validation Features](#requirement-reference-validation-features)
    - [Requirement Reference Validation Exit Codes](#requirement-reference-validation-exit-codes)
    - [Requirement Reference Validation Integration](#requirement-reference-validation-integration)
    - [Requirement Reference Validation Requirements](#requirement-reference-validation-requirements)
    - [Requirement Reference Validation Output Example](#requirement-reference-validation-output-example)
  - [audit\_requirements\_coverage.py](#audit_requirements_coveragepy)
    - [Requirements Coverage Audit Purpose](#requirements-coverage-audit-purpose)
    - [Requirements Coverage Audit Usage](#requirements-coverage-audit-usage)
    - [Requirements Coverage Audit Features](#requirements-coverage-audit-features)
    - [Requirements Coverage Audit Exit Codes](#requirements-coverage-audit-exit-codes)
    - [Requirements Coverage Audit Integration](#requirements-coverage-audit-integration)
    - [Requirements Coverage Audit Classification System](#requirements-coverage-audit-classification-system)
    - [Requirements Coverage Audit Requirements](#requirements-coverage-audit-requirements)
  - [validate\_go\_code\_blocks.py](#validate_go_code_blockspy)
    - [Go Code Blocks Validation Purpose](#go-code-blocks-validation-purpose)
    - [Go Code Blocks Validation Usage](#go-code-blocks-validation-usage)
    - [Go Code Blocks Validation Features](#go-code-blocks-validation-features)
    - [Go Code Blocks Validation Exit Codes](#go-code-blocks-validation-exit-codes)
    - [Go Code Blocks Validation Integration](#go-code-blocks-validation-integration)
    - [Go Code Blocks Validation Requirements](#go-code-blocks-validation-requirements)
    - [Go Code Blocks Validation Output Example](#go-code-blocks-validation-output-example)
  - [validate\_api\_go\_defs\_index.py](#validate_api_go_defs_indexpy)
    - [Go Definitions Index Validation Purpose](#go-definitions-index-validation-purpose)
    - [Go Definitions Index Validation Usage](#go-definitions-index-validation-usage)
    - [Go Definitions Index Validation Features](#go-definitions-index-validation-features)
    - [Go Definitions Index Validation Exit Codes](#go-definitions-index-validation-exit-codes)
    - [Go Definitions Index Validation Requirements](#go-definitions-index-validation-requirements)
  - [validate\_go\_signature\_sync.py](#validate_go_signature_syncpy)
    - [Go Signature Sync Validation Purpose](#go-signature-sync-validation-purpose)
    - [Go Signature Sync Validation Usage](#go-signature-sync-validation-usage)
    - [Go Signature Sync Validation Features](#go-signature-sync-validation-features)
    - [Go Signature Sync Validation Exit Codes](#go-signature-sync-validation-exit-codes)
    - [Go Signature Sync Validation Integration](#go-signature-sync-validation-integration)
    - [Go Signature Sync Validation Requirements](#go-signature-sync-validation-requirements)
    - [Go Signature Sync Validation Output Example](#go-signature-sync-validation-output-example)
  - [validate\_go\_spec\_references.py](#validate_go_spec_referencespy)
    - [Go Specification References Validation Purpose](#go-specification-references-validation-purpose)
    - [Go Specification References Validation Usage](#go-specification-references-validation-usage)
    - [Go Specification References Validation Features](#go-specification-references-validation-features)
    - [Go Specification References Validation Exit Codes](#go-specification-references-validation-exit-codes)
    - [Go Specification References Validation Integration](#go-specification-references-validation-integration)
    - [Go Specification References Validation Requirements](#go-specification-references-validation-requirements)
    - [Go Specification References Validation Output Example](#go-specification-references-validation-output-example)
  - [apply\_go\_spec\_references.py](#apply_go_spec_referencespy)
    - [Go Specification References Application Purpose](#go-specification-references-application-purpose)
    - [Go Specification References Application Usage](#go-specification-references-application-usage)
    - [Go Specification References Application Features](#go-specification-references-application-features)
    - [Go Specification References Application Exit Codes](#go-specification-references-application-exit-codes)
    - [Go Specification References Application Requirements](#go-specification-references-application-requirements)
    - [Go Specification References Application Output Example](#go-specification-references-application-output-example)
    - [Go Specification References Application Integration](#go-specification-references-application-integration)
    - [Go Specification References Application Workflow](#go-specification-references-application-workflow)
  - [validate\_go\_spec\_signature\_consistency.py](#validate_go_spec_signature_consistencypy)
    - [Go Signature Consistency Validation Purpose](#go-signature-consistency-validation-purpose)
    - [Go Signature Consistency Validation Usage](#go-signature-consistency-validation-usage)
    - [Go Signature Consistency Validation Features](#go-signature-consistency-validation-features)
    - [Go Signature Consistency Validation Exit Codes](#go-signature-consistency-validation-exit-codes)
    - [Go Signature Consistency Validation Integration](#go-signature-consistency-validation-integration)
    - [Go Signature Consistency Validation Requirements](#go-signature-consistency-validation-requirements)
    - [Go Signature Consistency Validation Output Example](#go-signature-consistency-validation-output-example)
- [Go Linting Checks](#go-linting-checks)
  - [Context Propagation Check](#context-propagation-check)
    - [Context Propagation Check Purpose](#context-propagation-check-purpose)
    - [Context Propagation Check Usage](#context-propagation-check-usage)
    - [Context Propagation Check Features](#context-propagation-check-features)
    - [Context Propagation Check Configuration](#context-propagation-check-configuration)
    - [Context Propagation Check Exit Codes](#context-propagation-check-exit-codes)
    - [Context Propagation Check Integration](#context-propagation-check-integration)
    - [Context Propagation Check Requirements](#context-propagation-check-requirements)
    - [Context Propagation Check Output Example](#context-propagation-check-output-example)
- [Common Features](#common-features)
  - [Path Targeting](#path-targeting)
- [CI Integration](#ci-integration)
  - [Current CI Integration](#current-ci-integration)
- [Development](#development)
- [Maintenance](#maintenance)
- [Utility Modules](#utility-modules)
- [Related Documentation](#related-documentation)

## Overview

This directory contains utility scripts for the NovusPack project.

## Available Scripts

### generate_anchor.py

Generates GitHub-style markdown anchors from markdown headings.

#### Anchor Generation Purpose

- Generates consistent markdown anchors for documentation links
- Implements GitHub's anchor generation algorithm
- Useful for creating links to specific sections in markdown files
- Supports headings with inline code (backticks)

#### Anchor Generation Usage

```bash
# Generate anchor for the heading at a specific file line
python3 scripts/generate_anchor.py --line docs/tech_specs/api_core.md:42
# Output: #some-heading-anchor

# Print anchors for all headings in a file
python3 scripts/generate_anchor.py --file docs/tech_specs/api_core.md
# Output (one per heading):
# docs/tech_specs/api_core.md:1: H1 Title => #title

# Via Makefile (preferred)
make generate-anchor LINE="docs/tech_specs/api_core.md:42"
make generate-anchor FILE="docs/tech_specs/api_core.md"
```

#### Anchor Generation Features

- Implements GitHub's markdown anchor generation algorithm
- Removes backticks but preserves their content (e.g., `` `code` `` -> `code`)
- Converts to lowercase
- Removes special characters except word characters, spaces, and hyphens
- Collapses sequences of spaces and hyphens into a single hyphen
- Strips leading and trailing hyphens
- Returns anchor with `#` prefix for use in markdown links

#### Anchor Generation Integration

The script is available via Makefile:

```bash
# Print anchors for all headings in a file
make generate-anchor FILE="docs/tech_specs/api_core.md"

# Print anchor for the heading at a specific line in a file
make generate-anchor LINE="docs/tech_specs/api_core.md:42"
```

**Note:** This interface avoids passing heading text through the shell, which eliminates quoting issues (including backticks).
If you must pass text containing backticks on the command line for other tooling, prefer single quotes and escape backticks as needed.

#### Anchor Generation Requirements

- Python 3.x (standard library only, no external dependencies)

### validate_links.py

Validates all internal markdown links and anchors in the documentation.

#### Link Validation Purpose

- Ensures all documentation cross-references are valid
- Prevents broken links in requirements and tech spec files
- Validates links in code comments within code blocks
- Checks that requirements files reference tech specs

#### Link Validation Usage

```bash
# Basic validation
python3 scripts/validate_links.py

# With verbose output
python3 scripts/validate_links.py --verbose

# Save detailed report to file (use tmp/ for reports)
python3 scripts/validate_links.py --output tmp/validation_report.txt

# Check requirements coverage (ensures all requirements reference tech specs)
python3 scripts/validate_links.py --check-coverage

# Check specific file
python3 scripts/validate_links.py --path docs/tech_specs/api_file_management.md

# Check specific directory (recursive)
python3 scripts/validate_links.py --path docs/requirements

# Check multiple paths (comma-separated)
python3 scripts/validate_links.py --path docs/requirements,docs/tech_specs

# Show help
python3 scripts/validate_links.py --help
```

#### Link Validation Features

- Validates 2,600+ internal links across 43+ markdown files
- Categorizes files by type (requirements, tech specs, other)
- Detects broken anchors, missing files, and internal anchor issues
- Excludes false positives from inline code and function signatures
- Validates links in code comments (e.g., `// See [doc](file.md)`)
- Provides line numbers for all broken links
- Organized reporting by file type
- Uses file content caching to reduce I/O overhead
- Regex patterns compiled at module level for performance
- Exit code 1 on failures (suitable for CI)

#### Link Validation Exit Codes

- `0`: All links valid
- `1`: Broken links found or coverage issues detected

#### Link Validation Integration

The script is integrated into the CI pipeline:

```bash
# Run locally via Makefile (all files)
make validate-links

# Check specific files/directories
make validate-links PATHS="docs/requirements/core.md,docs/tech_specs"

# Part of full CI suite
make ci
```

#### Link Validation Requirements

- Python 3.x (standard library only, no external dependencies)

### audit_feature_coverage.py

Audits feature file coverage for requirements by checking which requirements have feature files that reference them using `@REQ-*` tags.

#### Feature Coverage Audit Purpose

- Ensures all requirements are covered by BDD feature files
- Identifies requirements without corresponding feature tests
- Helps maintain test coverage quality
- Validates traceability from requirements to tests

#### Feature Coverage Audit Usage

```bash
# Basic audit
python3 scripts/audit_feature_coverage.py

# With verbose output
python3 scripts/audit_feature_coverage.py --verbose

# Check specific file
python3 scripts/audit_feature_coverage.py --path features/file_management/write.feature

# Check specific directory (recursive)
python3 scripts/audit_feature_coverage.py --path features/file_management

# Check multiple paths (comma-separated)
python3 scripts/audit_feature_coverage.py --path features/file_management,features/compression

# Show help
python3 scripts/audit_feature_coverage.py --help

# Run via Makefile
make audit-feature-coverage

# Part of combined coverage audit
make audit-coverage
```

#### Feature Coverage Audit Features

- Scans all requirements markdown files (excluding README.md and traceability.md)
- Excludes obsolete and documentation-only requirements
- Searches feature files for `@REQ-*` tag references
- Reports requirements without any feature coverage
- Provides summary statistics (total, covered, missing)
- Verbose mode shows detailed search progress
- Lists all uncovered requirements for easy remediation

#### Feature Coverage Audit Exit Codes

- `0`: All requirements have feature coverage
- `1`: One or more requirements have no feature coverage

#### Feature Coverage Audit Integration

The script is integrated into the CI pipeline:

```bash
# Run locally via Makefile (all files)
make audit-feature-coverage

# Part of combined coverage audit
# Note: This check is skipped when PATHS is specified, as it requires checking all requirements and feature files
make audit-coverage

# Part of full CI suite
make ci
```

**Note:** When running `make docs-check PATHS="..."`, this validation is automatically skipped because it requires checking all requirements and feature files to validate coverage properly.

#### Feature Coverage Audit Requirements

- Python 3.x (standard library only, no external dependencies)

### validate_heading_numbering.py

Validates markdown heading numbering consistency across all markdown files in the repository.

#### Heading Numbering Validation Purpose

- Ensures consistent heading numbering structure in all documentation
- Validates that heading levels match number depth (H2=1 number, H3=2 numbers, etc.)
- Checks sequential numbering within each parent section
- Verifies child heading numbers match their parent section number
- Helps maintain professional documentation standards

#### Heading Numbering Validation Usage

```bash
# Basic validation
python3 scripts/validate_heading_numbering.py

# With verbose output
python3 scripts/validate_heading_numbering.py --verbose

# Save detailed report to file
python3 scripts/validate_heading_numbering.py --output report.txt

# Check specific file
python3 scripts/validate_heading_numbering.py --path docs/tech_specs/api_file_management.md

# Check specific directory (recursive)
python3 scripts/validate_heading_numbering.py --path docs/requirements

# Check multiple paths (comma-separated)
python3 scripts/validate_heading_numbering.py --path docs/requirements,docs/tech_specs

# Show help
python3 scripts/validate_heading_numbering.py --help
```

#### Heading Numbering Validation Features

- Scans all markdown files in the repository (excluding hidden directories)
- Validates numbered headings starting at H2 level and beyond
- Checks heading depth matches number depth:
  - H2 sections: 1 number (e.g., "## 1 Title")
  - H3 sections: 2 numbers (e.g., "### 1.1 Subtitle")
  - H4 sections: 3 numbers (e.g., "#### 1.1.1 Details")
  - And so on for deeper levels
- Validates sequential numbering within parent sections
- Ensures child heading numbers match parent prefixes
- Detects duplicate heading titles (excluding numbering) across all levels
- Allows backticks in headings; case inside backticks is not checked for Title Case
- Warns about H3+ headings with numbering exceeding 20 (e.g., "### 3.25")
- Warns about H4+ headings with single-word titles (e.g., "#### 1.2.3 Title")
- Warns about overly-deeply nested headings (H6 and beyond)
- Provides line numbers and detailed error messages
- Organized reporting by error type:
  - Organizational heading errors (headings with no content)
  - Heading numbering errors (numbering issues)
- Sorted headings output shows only headings with numbering errors (excludes duplicate-only errors)
- Regex patterns compiled at module level for performance
- Exit code 1 on failures (suitable for CI)

#### Heading Numbering Validation Exit Codes

- `0`: All heading numbering is valid
- `1`: Heading numbering errors found

#### Heading Numbering Validation Integration

The script is integrated into the CI pipeline:

```bash
# Run locally via Makefile (all files)
make validate-heading-numbering

# Check specific files/directories
make validate-heading-numbering PATHS="docs/requirements,docs/tech_specs"

# Part of full CI suite
make ci
```

#### Heading Numbering Validation Requirements

- Python 3.x (standard library only, no external dependencies)

#### Heading Numbering Validation Output Example

```text
======================================================================
Markdown Heading Numbering Validation
======================================================================

Scanning 58 markdown files...

======================================================================
Summary
======================================================================
Found 149 heading numbering error(s):

/path/to/file.md:
  Line 486: Non-sequential numbering: got '6.1', expected '6.6' (previous was '6.5')
    ### 6.1 Breaking Changes Required
  Line 745: Non-sequential numbering: got '10.1', expected '10.5' (previous was '10.4')
    ### 10.1 Key Changes
  Line 892: Duplicate heading title 'Implementation Details' (also appears at line 234, line 567).
    #### 12.3.7.2 Implementation Details

Found 3 warning(s):

/path/to/file.md:
  Line 234: WARNING: H3 heading has numbering '3.25' where number 25 exceeds 20.
    ### 3.25 Deep Nesting Example
  Line 567: WARNING: H4 heading has a single-word title 'Details'.
    #### 1.2.3 Details

======================================================================
Sorted headings from first error (line 486) in /path/to/file.md:
======================================================================

The following headings should be in this order (sorted by numeric values):

Format: Line X: [CURRENT] -> [CORRECT] Title

Line 486: ## [6.1] -> [6.6] Breaking Changes Required
Line 745: ### [10.1] -> [10.5] Key Changes
```

Note: The example output above shows example headings that are part of the validation report format. These are displayed in code blocks and are not actual headings in the README file structure.

### apply_heading_corrections.py

Automatically applies heading numbering corrections from `validate_heading_numbering.py` output to markdown files.

#### Heading Corrections Application Purpose

- Automatically fixes heading numbering errors identified by validation
- Applies corrections to multiple files in a single run
- Preserves heading formatting (periods, spacing)
- Reduces manual editing effort for large documentation sets

#### Heading Corrections Application Usage

```bash
# Via Makefile - pipe from validation (reads from stdin)
make validate-heading-numbering PATHS="file.md" | make apply-heading-corrections

# Via Makefile - read from file
make apply-heading-corrections INPUT="tmp/heading_report.txt"

# Via Makefile - dry run (preview changes)
make apply-heading-corrections INPUT="tmp/heading_report.txt" DRY_RUN=1

# Direct usage - pipe from validation script (read from stdin)
python3 scripts/validate_heading_numbering.py --path file.md | \
    python3 scripts/apply_heading_corrections.py

# Direct usage - read from saved output file
python3 scripts/apply_heading_corrections.py --input tmp/heading_report.txt

# Direct usage - dry run (preview changes without modifying files)
python3 scripts/apply_heading_corrections.py --input tmp/heading_report.txt --dry-run

# Direct usage - with verbose output
python3 scripts/apply_heading_corrections.py --input tmp/heading_report.txt --verbose

# Show help
python3 scripts/apply_heading_corrections.py --help
```

#### Heading Corrections Application Features

- Parses correction output from `validate_heading_numbering.py`
- Extracts file paths and line numbers from validation output
- Groups corrections by file for efficient processing
- Applies corrections in reverse line order (prevents line number shifts)
- Preserves period formatting (H2 headings with/without periods)
- Handles multiple files in a single run
- Dry-run mode for safe preview of changes
- Detailed error reporting for failed corrections

#### Heading Corrections Application Exit Codes

- `0`: All corrections applied successfully
- `1`: One or more corrections failed to apply

#### Heading Corrections Application Requirements

- Python 3.x (standard library only, no external dependencies)
- Output from `validate_heading_numbering.py` as input

#### Heading Corrections Application Output Example

```text
Found 4 corrections to apply.
Processing tmp/test_multi_file1.md (2 corrections)...
  Updated tmp/test_multi_file1.md
Processing tmp/test_multi_file2.md (2 corrections)...
  Updated tmp/test_multi_file2.md

Applied 4 corrections, 0 failed.
```

#### Heading Corrections Application Integration

The script can be used via Makefile or directly:

```bash
# Via Makefile - pipe from validation (reads from stdin)
make validate-heading-numbering PATHS="file.md" | make apply-heading-corrections

# Via Makefile - read from file
make apply-heading-corrections INPUT="tmp/report.txt"

# Via Makefile - dry run (preview changes)
make apply-heading-corrections INPUT="tmp/report.txt" DRY_RUN=1

# Direct usage - pipe from validation script
python3 scripts/validate_heading_numbering.py --path file.md | \
    python3 scripts/apply_heading_corrections.py

# Direct usage - read from file
python3 scripts/apply_heading_corrections.py --input tmp/report.txt

# Direct usage - dry run first
python3 scripts/apply_heading_corrections.py --input tmp/report.txt --dry-run
```

#### Heading Corrections Application Workflow

The typical workflow is:

1. **Validate headings** and save output:

   ```bash
   make validate-heading-numbering PATHS="docs/" --output tmp/heading_report.txt
   # Or directly:
   python3 scripts/validate_heading_numbering.py --path docs/ \
       --output tmp/heading_report.txt
   ```

2. **Review corrections** (optional dry-run):

   ```bash
   make apply-heading-corrections INPUT="tmp/heading_report.txt" DRY_RUN=1
   # Or directly:
   python3 scripts/apply_heading_corrections.py \
       --input tmp/heading_report.txt --dry-run
   ```

3. **Apply corrections**:

   ```bash
   make apply-heading-corrections INPUT="tmp/heading_report.txt"
   # Or directly:
   python3 scripts/apply_heading_corrections.py \
       --input tmp/heading_report.txt
   ```

4. **Verify fixes**:

   ```bash
   make validate-heading-numbering PATHS="docs/"
   # Or directly:
   python3 scripts/validate_heading_numbering.py --path docs/
   ```

### validate_req_references.py

Validates that all requirement references in feature files are correct and point to existing requirement definitions in the documentation.

#### Requirement Reference Validation Purpose

- Ensures all @REQ-\* tags in feature files reference valid requirements
- Validates that requirement references point to the correct documentation file
- Detects missing requirement definitions
- Helps maintain traceability from feature files to requirements
- Prevents broken requirement references

#### Requirement Reference Validation Usage

```bash
# Basic validation
python3 scripts/validate_req_references.py

# With verbose output
python3 scripts/validate_req_references.py --verbose

# Check specific file
python3 scripts/validate_req_references.py --path features/file_management/write.feature

# Check specific directory (recursive)
python3 scripts/validate_req_references.py --path features/file_management

# Check multiple paths (comma-separated)
python3 scripts/validate_req_references.py --path features/file_management,features/compression

# Show help
python3 scripts/validate_req_references.py --help
```

#### Requirement Reference Validation Features

- Scans all 932 feature files for @REQ-\* tags
- Extracts 1199 requirement definitions from 17 requirement files
- Maps requirement categories to the correct documentation files:
  - REQ-API_BASIC-\* -> basic_ops.md
  - REQ-FILEMGMT-\* -> file_mgmt.md
  - REQ-COMPR-\* -> compression.md
  - REQ-SEC-\* -> security.md
  - And many more categories
- Validates requirement references exist in expected files
- Reports missing requirements
- Reports invalid references (wrong file)
- Reports errors in REQ ID format or unknown categories
- Provides summary statistics
- Verbose mode shows detailed scanning progress

#### Requirement Reference Validation Exit Codes

- `0`: All requirement references are valid
- `1`: Invalid references, missing requirements, or format errors found

#### Requirement Reference Validation Integration

The script is integrated into the CI pipeline:

```bash
# Run locally via Makefile (all files)
make validate-req-references

# Part of full CI suite
# Note: This check is skipped when PATHS is specified, as it requires checking all feature files
make ci
```

**Note:** When running `make docs-check PATHS="..."`, this validation is automatically skipped because it requires checking all feature files to validate references properly.

#### Requirement Reference Validation Requirements

- Python 3.x (standard library only, no external dependencies)

#### Requirement Reference Validation Output Example

```text
=== Requirement Reference Validation ===

Loading requirement definitions...
Loaded 1199 requirement definitions from 17 files

Scanning feature files...
Found 1051 unique REQ tags across 932 feature files

=== Validation Results ===

=== Summary ===
Total unique REQ references: 1051
Valid references: 1026
Invalid references (wrong file): 0
Missing references: 25
Errors (format/category): 0
```

### audit_requirements_coverage.py

Audits requirements coverage for tech specs by checking which tech specs have requirements that reference them using inline markdown links.

#### Requirements Coverage Audit Purpose

- Ensures all tech specs are referenced by requirements documentation
- Identifies tech specs without requirement links
- Helps maintain documentation traceability
- Validates bidirectional linkage between requirements and specs

#### Requirements Coverage Audit Usage

```bash
# Basic audit
python3 scripts/audit_requirements_coverage.py

# With verbose output
python3 scripts/audit_requirements_coverage.py --verbose

# Check specific file
python3 scripts/audit_requirements_coverage.py --path docs/requirements/core.md

# Check specific directory (recursive)
python3 scripts/audit_requirements_coverage.py --path docs/requirements

# Check multiple paths (comma-separated)
python3 scripts/audit_requirements_coverage.py --path docs/requirements/core.md,docs/requirements/security.md

# Show help
python3 scripts/audit_requirements_coverage.py --help

# Run via Makefile
make audit-requirements-coverage

# Part of combined coverage audit
make audit-coverage
```

#### Requirements Coverage Audit Features

- Scans all tech spec markdown files (excluding index files)
- Searches requirement files for relative link references
- Reports specs without any requirement coverage
- **Heading-level coverage checking**: Analyzes H2+ headings within tech specs
- **Intelligent classification**: Categorizes headings as functional, architectural, organizational, or excluded
- **Weighted architectural detection**: Uses scoring system to identify architectural headings with high confidence
- **Content analysis**: Detects function signatures, type definitions, and example code in sections
- **Organizational heading detection**: Identifies purely organizational headings that should be restructured
- Provides detailed summary statistics with breakdown by classification type:
  - Missing headings (errors): Functional headings requiring requirements
  - Architectural headings (warnings): Design/structure headings flagged for review (don't require requirements)
  - Missing headings (warnings): Other headings missing requirements
  - Organizational headings (warnings): Headings with no direct content
- Verbose mode shows detailed search progress
- Lists all uncovered specs and headings for easy remediation

#### Requirements Coverage Audit Exit Codes

- `0`: All tech specs are referenced by requirements and all functional headings have requirement coverage
- `1`: One or more tech specs have no requirement coverage, or functional headings are missing requirements

**Note:** Architectural headings are reported as warnings but do not cause the audit to fail, as they describe design/structure rather than testable behavior.
These are reported for completeness and user review.

#### Requirements Coverage Audit Integration

The script is integrated into the CI pipeline:

```bash
# Run locally via Makefile (all files)
make audit-requirements-coverage

# Part of combined coverage audit
# Note: This check is skipped when PATHS is specified, as it requires checking all tech specs and requirements
make audit-coverage

# Part of full CI suite
make ci
```

**Note:** When running `make docs-check PATHS="..."`, this validation is automatically skipped because it requires checking all tech specs and requirements to validate coverage properly.

#### Requirements Coverage Audit Classification System

The script uses an intelligent classification system to determine which headings need requirements:

##### Functional Headings (Errors if missing requirements)

- Headings with function signatures or type definitions
- Headings with functional keywords (behavior, operation, method, etc.)
- H2 headings by default (unless classified otherwise)
- These describe testable behavior and must have requirements

##### Architectural Headings (Warnings, no requirements needed)

- Detected using weighted scoring system based on keywords
- High confidence indicators: "type definition", "system architecture", "interface definition" (+3 points)
- Medium confidence: "structure", "system", "interface", "architecture" (+2 points)
- Threshold: Score ≥ 2 to classify as architectural
- Negative weights for non-architectural phrases (e.g., "purpose", "usage", "management")
- These describe design/structure and are flagged for review but don't require requirements

##### Organizational Headings (Warnings)

- Headings with no direct content (only subheadings)
- Max 5 lines of prose, no code blocks
- Should be removed or restructured to reduce nesting

##### Excluded Headings

- Headings matching exclusion patterns (e.g., "best practices", "examples", "cross-references")
- Implementation detail headings
- Overview sections (section 0.x)

#### Requirements Coverage Audit Requirements

- Python 3.x (standard library only, no external dependencies)

### validate_go_code_blocks.py

Validates Go code blocks in tech specs documentation to ensure they follow conventions.

#### Go Code Blocks Validation Purpose

- Ensures each Go code block has at most one type or interface definition
- Validates that each Go code block is under a different heading
- Helps maintain consistent documentation structure
- Prevents code blocks from violating documentation standards

#### Go Code Blocks Validation Usage

```bash
# Basic validation
python3 scripts/validate_go_code_blocks.py

# With verbose output
python3 scripts/validate_go_code_blocks.py --verbose

# Save detailed report to file
python3 scripts/validate_go_code_blocks.py --output dev_docs/go_code_blocks_audit.md

# Check specific file
python3 scripts/validate_go_code_blocks.py --path docs/tech_specs/api_file_management.md

# Check specific directory (recursive)
python3 scripts/validate_go_code_blocks.py --path docs/tech_specs

# Check multiple paths (comma-separated)
python3 scripts/validate_go_code_blocks.py --path docs/tech_specs/api_file_management.md,docs/tech_specs/api_core.md

# Show help
python3 scripts/validate_go_code_blocks.py --help
```

#### Go Code Blocks Validation Features

- Scans all markdown files in `docs/tech_specs` (or specified paths)
- Detects Go code blocks (fenced code blocks with `go` language identifier)
- Counts type definitions in each code block:
  - Struct definitions: `type Name struct`
  - Interface definitions: `type Name interface`
  - Generic type definitions: `type Name[T ...] struct/interface`
  - Type aliases: `type Name SomeType` (including built-in types)
- Validates that each code block has at most one type/interface definition
- Checks that each code block is under a different heading
- Validates heading format: definition name and kind word (e.g. `` `Package.Write` Method ``); definition names preferred in backticks; case inside backticks ignored
- Emits `heading_prefer_backticks` warning when a definition name in a heading is not in backticks (suggests corrected heading)
- When only warnings (e.g. heading_prefer_backticks) are found, the script exits 0
- Reports issues with file paths and line numbers
- Provides summary statistics and detailed breakdown

#### Go Code Blocks Validation Exit Codes

- `0`: All Go code blocks comply with requirements, or only warnings (e.g. heading_prefer_backticks) were found
- `1`: One or more code blocks violate requirements (errors present)

#### Go Code Blocks Validation Integration

The script is integrated into the documentation checks:

```bash
# Run locally via Makefile (all files)
make validate-go-code-blocks

# Check specific files/directories
make validate-go-code-blocks PATHS="docs/tech_specs/api_file_management.md,docs/tech_specs"

# Part of documentation checks (runs before heading validation)
make docs-check

# Part of full CI suite
make ci
```

#### Go Code Blocks Validation Requirements

- Python 3.x (standard library only, no external dependencies)

#### Go Code Blocks Validation Output Example

```text
======================================================================
Go Code Blocks Validation Summary
======================================================================
Files audited:        20
Total code blocks:     304
Total issues found:   38

Breakdown by issue type:
  Duplicate Heading........ 6
  Multiple Types........... 32
```

### validate_api_go_defs_index.py

Validates that all Go API definitions in tech specs are listed in the Go definitions index.

Business logic and usage are documented in [scripts/validate_api_go_defs_index.md](validate_api_go_defs_index.md).
Use `make validate-go-defs-index` to run the check.

Implementation details are in [`scripts/lib/go_defs_index/`](../scripts/lib/go_defs_index/), including the placement scoring modules used by the validator.

#### Go Definitions Index Validation Purpose

- Ensures every discovered Go API definition in tech specs appears in the index.
- Detects missing index entries, orphaned entries, wrong-section entries, and incorrect link targets.
- Enforces description rules (minimum length and uniqueness).
- Reports low-confidence placements that require manual review.

#### Go Definitions Index Validation Usage

```bash
# Run the validator (full scan).
make validate-go-defs-index

# Verbose output and write a report file.
make validate-go-defs-index VERBOSE=1 NO_COLOR=1 OUTPUT="tmp/go_defs_index.txt"

# Apply high-confidence index updates (interactive confirmation required).
make validate-go-defs-index APPLY=1
```

#### Go Definitions Index Validation Features

- Scans `docs/tech_specs/*.md` for ` ```go ` code blocks and extracts types, methods, and functions.
- Builds an expected index tree using confidence-scored placement.
- Compares expected entries with current index entries and reports discrepancies.
- Validates index entry descriptions:
  - Missing or too-short descriptions are errors.
  - Duplicate description text across entries is an error.
- Emits ordering warnings when existing entries are out of order.

#### Go Definitions Index Validation Exit Codes

- `0`: No errors found.
- `1`: Errors found.
- With `NO_FAIL=1`, the validator exits with `0` even when errors are found.

#### Go Definitions Index Validation Requirements

- Python 3.x

### validate_go_signature_sync.py

Validates Go type, method, and function signatures in implementation against tech specs.

#### Go Signature Sync Validation Purpose

- Ensures implementation signatures match tech spec definitions
- Detects signature mismatches (different parameters, return types, etc.)
- Identifies missing signatures in implementation (not yet implemented)
- Flags additional signatures in implementation not in specs (likely helpers)
- Validates proper context propagation
- Detects forbidden empty interface parameter types
- Warns about discouraged `any` type usage

#### Go Signature Sync Validation Usage

```bash
# Basic validation
python3 scripts/validate_go_signature_sync.py

# With verbose output
python3 scripts/validate_go_signature_sync.py --verbose

# Save detailed report to file
python3 scripts/validate_go_signature_sync.py --output report.txt

# Specify custom directories
python3 scripts/validate_go_signature_sync.py \
    --specs-dir docs/tech_specs \
    --impl-dir api/go

# Exit with code 0 even if errors are found
python3 scripts/validate_go_signature_sync.py --no-fail

# Show help
python3 scripts/validate_go_signature_sync.py --help
```

#### Go Signature Sync Validation Features

- Extracts all public signatures from Go implementation files
- Extracts all signatures from tech specs markdown files
- Compares normalized signatures for exact matches
- Detects empty interface parameter types (forbidden)
- Warns about `any` type usage (discouraged)
- Filters out high-confidence helper functions (test files, internal packages, etc.)
- Identifies public API methods missing from specs (errors)
- Provides detailed location information for all mismatches
- Supports custom specs and implementation directories

#### Go Signature Sync Validation Exit Codes

- `0`: All signatures are in sync
- `1`: Signature mismatches, missing signatures, or errors found

#### Go Signature Sync Validation Integration

The script is integrated into the CI pipeline:

```bash
# Run locally via Makefile (all files)
make validate-go-signatures

# Part of Go CI checks
make ci-go

# Part of full CI suite
make ci
```

**Note:** This validation is part of the Go CI workflow and runs automatically when Go code changes.

#### Go Signature Sync Validation Requirements

- Python 3.x (standard library only, no external dependencies)

#### Go Signature Sync Validation Output Example

```text
======================================================================
Go Signature Sync Validation
======================================================================

Collecting signatures from Go implementation...
Found 245 public signatures in implementation

Collecting signatures from tech specs...
Found 250 public signatures in specs

Comparing signatures...

======================================================================
Errors
======================================================================

Found 3 signature mismatch(es):

Signature: Package.AddFile
  Implementation: func (p *Package) AddFile(ctx context.Context, path string, source FileSource) error
    Location: api/go/novus_package/package_writer.go:123
  Specification:  func (p *Package) AddFile(ctx context.Context, path string, source FileSource, opts *AddFileOptions) error
    Location: docs/tech_specs/api_file_management.md:456

======================================================================
Warnings
======================================================================

Found 5 signature(s) in specs but not in implementation:

  Package.RemoveFile
    Signature: func (p *Package) RemoveFile(ctx context.Context, path string) error
    Location: docs/tech_specs/api_file_management.md:789

Found 2 signature(s) in implementation but not in specs:

  Package.internalHelper
    Signature: func (p *Package) internalHelper() error
    Location: api/go/novus_package/package_internal.go:45
    (These may be helper functions, but should be checked)
```

### validate_go_spec_references.py

Validates Go file specification references against tech spec files and the API definitions index.

#### Go Specification References Validation Purpose

- Ensures all `// Specification:` comments in Go files reference valid tech spec sections
- Validates that referenced files exist
- Validates that section anchors exist in those files
- Provides suggestions for correct reference format
- Uses API definitions index to suggest correct references based on function/type names

#### Go Specification References Validation Usage

```bash
# Basic validation
python3 scripts/validate_go_spec_references.py

# With verbose output
python3 scripts/validate_go_spec_references.py --verbose

# Save detailed report to file
python3 scripts/validate_go_spec_references.py --output report.txt

# Specify custom repository root
python3 scripts/validate_go_spec_references.py --repo-root /path/to/repo

# Exit with code 0 even if errors are found
python3 scripts/validate_go_spec_references.py --no-fail

# Disable colored output
python3 scripts/validate_go_spec_references.py --no-color

# Show help
python3 scripts/validate_go_spec_references.py --help
```

#### Go Specification References Validation Features

- Scans all Go files for `// Specification:` comments
- Validates reference format: `file_name.md: section_number heading_text`
- Also accepts anchor format: `file_name.md#anchor`
- Validates that spec files exist
- Validates that section numbers exist in spec files
- Validates that heading text matches actual headings
- Uses API definitions index to suggest correct references
- Extracts function/type names from Go code context for better suggestions
- Filters out section 0 and cross-reference sections (not source of truth)
- Provides parseable output format for automated fixes

#### Go Specification References Validation Exit Codes

- `0`: All specification references are valid
- `1`: Invalid references, missing files, or format errors found

#### Go Specification References Validation Integration

The script is integrated into the CI pipeline:

```bash
# Run locally via Makefile (all files)
make validate-go-spec-references

# Part of Go CI checks
make ci-go

# Part of full CI suite
make ci
```

**Note:** This validation is part of the Go CI workflow and runs automatically when Go code changes.

#### Go Specification References Validation Requirements

- Python 3.x (standard library only, no external dependencies)

#### Go Specification References Validation Output Example

```text
======================================================================
Go Specification References Validation
======================================================================

Scanning 156 Go files for specification references...

======================================================================
Errors
======================================================================

api/go/novus_package/package_writer.go:123: Invalid spec ref
  Current: api_file_management.md: 2.1 AddFile
  SUGGESTION: api_file_management.md: 2.1 AddFile Package Method

api/go/novus_package/package_reader.go:456: Invalid spec ref
  Current: api_file_management.md#21-addfile
  SUGGESTION: api_file_management.md: 2.1 AddFile Package Method

api/go/novus_package/package_writer.go:789: Invalid spec ref
  Current: invalid_file.md: 1.1 SomeMethod
  Also, spec file not found: invalid_file.md

Found 3 error(s)
```

### apply_go_spec_references.py

Automatically applies specification reference updates from `validate_go_spec_references.py` output to Go files.

#### Go Specification References Application Purpose

- Automatically fixes specification reference errors identified by validation
- Applies corrections to multiple files in a single run
- Preserves comment formatting and indentation
- Reduces manual editing effort for large codebases

#### Go Specification References Application Usage

```bash
# Via Makefile - read from file
make apply-go-spec-references INPUT="tmp/spec_ref_errors.txt"

# Via Makefile - dry run (preview changes)
make apply-go-spec-references INPUT="tmp/spec_ref_errors.txt" DRY_RUN=1

# Direct usage - read from saved output file
python3 scripts/apply_go_spec_references.py --input tmp/spec_ref_errors.txt

# Direct usage - dry run (preview changes without modifying files)
python3 scripts/apply_go_spec_references.py \
    --input tmp/spec_ref_errors.txt --dry-run

# Direct usage - with verbose output
python3 scripts/apply_go_spec_references.py \
    --input tmp/spec_ref_errors.txt --verbose

# Show help
python3 scripts/apply_go_spec_references.py --help
```

#### Go Specification References Application Features

- Parses correction output from `validate_go_spec_references.py`
- Extracts file paths and line numbers from validation output
- Groups corrections by file for efficient processing
- Applies corrections in reverse line order (prevents line number shifts)
- Preserves comment indentation and formatting
- Handles multiple files in a single run
- Dry-run mode for safe preview of changes
- Detailed error reporting for failed corrections
- Skips corrections that are already applied

#### Go Specification References Application Exit Codes

- `0`: All corrections applied successfully
- `1`: One or more corrections failed to apply

#### Go Specification References Application Requirements

- Python 3.x (standard library only, no external dependencies)
- Output from `validate_go_spec_references.py` as input

#### Go Specification References Application Output Example

```text
Found 4 corrections to apply.
Processing api/go/novus_package/package_writer.go (2 corrections)...
  Updated api/go/novus_package/package_writer.go
Processing api/go/novus_package/package_reader.go (2 corrections)...
  Updated api/go/novus_package/package_reader.go

Applied 4 specification reference update(s), 0 skipped, 0 failed.
```

#### Go Specification References Application Integration

The script can be used via Makefile or directly:

```bash
# Via Makefile - read from file
make apply-go-spec-references INPUT="tmp/spec_ref_errors.txt"

# Via Makefile - dry run (preview changes)
make apply-go-spec-references INPUT="tmp/spec_ref_errors.txt" DRY_RUN=1

# Direct usage - read from file
python3 scripts/apply_go_spec_references.py --input tmp/spec_ref_errors.txt

# Direct usage - dry run first
python3 scripts/apply_go_spec_references.py \
    --input tmp/spec_ref_errors.txt --dry-run
```

#### Go Specification References Application Workflow

The typical workflow is:

1. **Validate references** and save output:

   ```bash
   make validate-go-spec-references OUTPUT=tmp/spec_ref_errors.txt
   # Or directly:
   python3 scripts/validate_go_spec_references.py \
       --output tmp/spec_ref_errors.txt
   ```

2. **Review corrections** (optional dry-run):

   ```bash
   make apply-go-spec-references INPUT="tmp/spec_ref_errors.txt" DRY_RUN=1
   # Or directly:
   python3 scripts/apply_go_spec_references.py \
       --input tmp/spec_ref_errors.txt --dry-run
   ```

3. **Apply corrections**:

   ```bash
   make apply-go-spec-references INPUT="tmp/spec_ref_errors.txt"
   # Or directly:
   python3 scripts/apply_go_spec_references.py \
       --input tmp/spec_ref_errors.txt
   ```

4. **Verify fixes**:

   ```bash
   make validate-go-spec-references
   # Or directly:
   python3 scripts/validate_go_spec_references.py
   ```

### validate_go_spec_signature_consistency.py

Validates Go signature consistency within tech specs documentation.

#### Go Signature Consistency Validation Purpose

- Ensures the same function, method, or type is not defined multiple times with different signatures
- Detects methods defined for interfaces that are NOT in the canonical interface definition
- Detects interface stubs with methods not in canonical definition
- Detects interface stub methods that differ from canonical definitions
- Identifies duplicate identical signatures
- Warns about type/interface stubs (expected when fleshing out methods in separate sections)

#### Go Signature Consistency Validation Usage

```bash
# Basic validation
python3 scripts/validate_go_spec_signature_consistency.py

# With verbose output
python3 scripts/validate_go_spec_signature_consistency.py --verbose

# Save detailed report to file
python3 scripts/validate_go_spec_signature_consistency.py \
    --output tmp/signature_consistency_report.txt

# Check specific file
python3 scripts/validate_go_spec_signature_consistency.py \
    --path docs/tech_specs/api_basic_operations.md

# Check specific directory (recursive)
python3 scripts/validate_go_spec_signature_consistency.py \
    --path docs/tech_specs

# Check multiple paths (comma-separated)
python3 scripts/validate_go_spec_signature_consistency.py \
    --path docs/tech_specs/api_basic_operations.md,docs/tech_specs/api_core.md

# Save detailed report to file
python3 scripts/validate_go_spec_signature_consistency.py \
    --output tmp/signature_consistency_report.txt

# Exit with code 0 even if errors are found
python3 scripts/validate_go_spec_signature_consistency.py --no-fail

# Disable colored output
python3 scripts/validate_go_spec_signature_consistency.py --no-color

# Show help
python3 scripts/validate_go_spec_signature_consistency.py --help
```

#### Go Signature Consistency Validation Features

- Scans all markdown files in `docs/tech_specs` (or specified paths)
- Extracts all Go signatures from code blocks
- Detects duplicate signatures with different definitions
- Identifies canonical signatures using heuristics (heading level, body presence, etc.)
- Validates interface method consistency
- Checks for methods defined outside canonical interface definitions
- Validates stub methods match canonical definitions
- Filters out example code and signatures
- Provides detailed location information for all inconsistencies
- Suggests canonical signatures for duplicates

#### Go Signature Consistency Validation Exit Codes

- `0`: All signatures are consistent
- `1`: Signature inconsistencies found

#### Go Signature Consistency Validation Integration

The script is integrated into the documentation checks:

```bash
# Run locally via Makefile (all files)
make validate-go-spec-signature-consistency

# Check specific files/directories
make validate-go-spec-signature-consistency \
    PATHS="docs/tech_specs/api_file_management.md,docs/tech_specs"

# Part of documentation checks (runs after code blocks validation)
make docs-check

# Part of full CI suite
make ci
```

#### Go Signature Consistency Validation Requirements

- Python 3.x (standard library only, no external dependencies)

#### Go Signature Consistency Validation Output Example

```text
======================================================================
Go Signature Consistency
======================================================================

Found 20 markdown file(s) to audit
Extracting signatures from api_basic_operations.md...
Extracting signatures from api_core.md...
Found 250 total signatures

Checking for signature inconsistencies...

======================================================================
Errors
======================================================================

ERROR: Signature inconsistency for 'Package.AddFile':
  Signature: func (p *Package) AddFile(ctx context.Context, path string, source FileSource) error
    Location: docs/tech_specs/api_file_management.md:456
  Signature: func (p *Package) AddFile(ctx context.Context, path string, source FileSource, opts *AddFileOptions) error
    Location: docs/tech_specs/api_file_management.md:789

ERROR: Method 'RemoveFile' is defined for interface 'PackageWriter' but is not in the canonical interface definition:
  Canonical interface: docs/tech_specs/api_file_management.md:123
  Methods in canonical: ['AddFile', 'GetFile', 'ListFiles']
  Method location: docs/tech_specs/api_file_management.md:999

======================================================================
Warnings
======================================================================

WARNING: Type/interface stub detected for 'FileEntry':
  Canonical: docs/tech_specs/api_file_management.md:123 (has_body=True, field_count=5)
  Stub: docs/tech_specs/api_file_management.md:456 (has_body=False, field_count=0)

======================================================================
Summary
======================================================================
Signatures checked:        250
Unique definitions:        245
Errors found:              2
Warnings found:            1
```

## Go Linting Checks

### Context Propagation Check

The context propagation check ensures that functions calling other functions that require `context.Context` also require `context.Context` as a parameter. This enforces proper context propagation through the call chain, which is essential for handling timeouts, cancellations, and deadlines in Go applications.

#### Context Propagation Check Purpose

- Ensures functions calling context-requiring functions also accept context
- Prevents functions from creating `context.Background()` when they should accept context
- Enforces proper context propagation through the call chain
- Helps maintain Go best practices for context handling

#### Context Propagation Check Usage

The check is automatically run as part of golangci-lint:

```bash
# Run via Makefile (includes contextcheck)
make lint-go

# Run as part of CI checks
make ci-go

# Run directly with golangci-lint
cd api/go
golangci-lint run ./...

# Run only contextcheck linter
golangci-lint run --enable-only=contextcheck ./...
```

#### Context Propagation Check Features

- Automatically detects functions that call context-requiring functions
- Flags functions that create `context.Background()` instead of accepting context
- Validates context parameter presence in function signatures
- Provides file and line number for each violation
- Integrated into golangci-lint for consistent checking

#### Context Propagation Check Configuration

The check is configured in `api/go/.golangci.yml`:

```yaml
version: "2"

linters:
  enable:
    - contextcheck
```

#### Context Propagation Check Exit Codes

- `0`: No context propagation issues found
- `1`: Context propagation violations detected

#### Context Propagation Check Integration

The check is integrated into the CI pipeline:

```bash
# Run locally via Makefile
make lint-go

# Part of full CI suite
make ci
```

The check runs automatically in the `lint-go-v1` job in `.github/workflows/go-ci.yml`.

#### Context Propagation Check Requirements

- golangci-lint v2.7.0 or later
- Go 1.25 or later

#### Context Propagation Check Output Example

```text
generics/concurrency.go:253:38: Non-inherited new context, use function like `context.WithXXX` instead (contextcheck)
novus_package/package_lifecycle.go:109:33: Function `ValidatePath` should pass the context parameter (contextcheck)
novus_package/package_lifecycle.go:190:33: Function `ValidatePath` should pass the context parameter (contextcheck)
novus_package/package_lifecycle.go:481:33: Function `ValidatePath` should pass the context parameter (contextcheck)

4 issues:
* contextcheck: 4
```

## Common Features

### Path Targeting

Most validation and audit scripts support the `--path` (or `-p`) flag for targeted validation, and the Makefile targets support a `PATHS` variable:

- **Single file**: Check one specific file
- **Single directory**: Recursively check all files in a directory
- **Multiple paths**: Comma-separated list of files and/or directories

This feature is useful for:

- Validating only modified files during development
- Optimizing CI pipeline performance by checking specific paths
- Quick validation of a subset of the codebase
- Testing newly added documentation before full validation

**Note:** Some scripts require checking all files to validate properly and will skip when `PATHS` is specified:

- `validate-go-defs-index` - requires all tech specs to validate the index
- `validate-req-references` - requires all feature files to validate references
- `audit-feature-coverage` - requires all requirements and feature files to validate coverage
- `audit-requirements-coverage` - requires all tech specs and requirements to validate coverage
- `validate-go-signature-sync` - requires all implementation and spec files to compare signatures
- `validate-go-spec-references` - requires all Go files to validate references

When running `make docs-check PATHS="..."`, these checks are automatically skipped.

Examples (direct script usage):

```bash
# Single file
python3 scripts/validate_links.py --path docs/requirements/core.md

# Single directory (recursive)
python3 scripts/validate_req_references.py --path features/file_management

# Multiple paths (comma-separated)
python3 scripts/validate_heading_numbering.py --path docs/requirements,docs/tech_specs

# Combine files and directories
python3 scripts/audit_feature_coverage.py --path features/file_management/write.feature,features/compression
```

Examples (Makefile usage):

```bash
# Single file
make validate-links PATHS="docs/requirements/core.md"

# Single directory
make validate-heading-numbering PATHS="docs/tech_specs"

# Multiple paths (comma-separated)
make docs-check PATHS="docs/requirements,docs/tech_specs"

# Combine files and directories
make validate-go-code-blocks PATHS="docs/tech_specs/api_file_management.md,docs/tech_specs/api_core.md"
```

Behavior:

- Gracefully handles non-existent paths (shows warning, continues with valid paths)
- Validates file types (shows warning if wrong extension)
- Maintains backward compatibility (no `--path` or `PATHS` = default full scan)
- Makefile `PATHS` variable is passed to all scripts that support `--path`

## CI Integration

All scripts in this directory should be:

1. **Executable**: `chmod +x scripts/*.py`
2. **Integrated in Makefile**: Add target in root `Makefile`
3. **Added to CI**: Create workflow in `.github/workflows/`
4. **Documented**: Update this README

### Current CI Integration

| Script/Check                                | Makefile Target                               | GitHub Workflow                  | Status    |
| ------------------------------------------- | --------------------------------------------- | -------------------------------- | --------- |
| `validate_links.py`                         | `make validate-links`                         | `validate-links.yml`             | ✅ Active |
| `validate_heading_numbering.py`             | `make validate-heading-numbering`             | `validate-heading-numbering.yml` | ✅ Active |
| `apply_heading_corrections.py`              | `make apply-heading-corrections`              | N/A (utility script)             | ✅ Active |
| `validate_req_references.py`                | `make validate-req-references`                | `validate-req-references.yml`    | ✅ Active |
| `audit_feature_coverage.py`                 | `make audit-feature-coverage`                 | `audit-coverage.yml`             | ✅ Active |
| `audit_requirements_coverage.py`            | `make audit-requirements-coverage`            | `audit-coverage.yml`             | ✅ Active |
| `validate_api_go_defs_index.py`             | `make validate-go-defs-index`                 | `docs-check` (via Makefile)      | ✅ Active |
| `validate_go_code_blocks.py`                | `make validate-go-code-blocks`                | `docs-check` (via Makefile)      | ✅ Active |
| `validate_go_spec_signature_consistency.py` | `make validate-go-spec-signature-consistency` | `docs-check` (via Makefile)      | ✅ Active |
| `validate_go_signature_sync.py`             | `make validate-go-signatures`                 | `go-ci.yml`                      | ✅ Active |
| `validate_go_spec_references.py`            | `make validate-go-spec-references`            | `go-ci.yml`                      | ✅ Active |
| `apply_go_spec_references.py`               | `make apply-go-spec-references`               | N/A (utility script)             | ✅ Active |
| `generate_anchor.py`                        | `make generate-anchor`                        | N/A (utility script)             | ✅ Active |
| Context Propagation Check                   | `make lint-go`                                | `go-ci.yml`                      | ✅ Active |

## Development

When adding new scripts:

1. Place script in `scripts/` directory
2. Make it executable: `chmod +x scripts/your_script.py`
3. Add shebang line: `#!/usr/bin/env python3`
4. Use `ValidationIssue` class from [`scripts/lib/_validation_utils.py`](../scripts/lib/_validation_utils.py) for all error/warning reporting
5. Maintain a single consolidated list of issues (filter by severity/type when displaying)
6. Use shared utilities from [`scripts/lib/_validation_utils.py`](../scripts/lib/_validation_utils.py) instead of duplicating code:
   - Use `import_module_with_fallback()` for module imports
   - Use `find_markdown_files()` or `find_feature_files()` for file finding
   - Use path validation functions (`is_safe_path()`, `validate_file_name()`, etc.)
   - Use heading extraction functions instead of custom implementations
   - Use `generate_anchor_from_heading()` directly (no wrapper functions)
   - Use `FileContentCache` to avoid repeated file reads (pass as parameter to functions that read files)
7. For performance, compile regex patterns at module level instead of inside functions
8. Add type hints to function signatures for better code clarity and maintainability
9. Extract magic numbers and strings to module-level constants for improved readability
10. Keep functions focused and reasonably sized (prefer smaller, focused functions over large monolithic ones)
11. Use consistent error handling with `ValidationIssue` objects and specific exception types
12. Add to root `Makefile` with appropriate target
13. Create GitHub workflow in `.github/workflows/`
14. Update this README
15. Ensure script has `--help` option
16. Use proper exit codes (0 = success, 1 = failure)
17. Use `OutputBuilder` from [`scripts/lib/_validation_utils.py`](../scripts/lib/_validation_utils.py) for consistent output formatting:
    - `add_success_message()` when validation passes with no issues.
    - `add_failure_message()` when there are errors
    - `add_warnings_only_message()` when there are only warnings (exit 0; optional `verbose_hint` for run-with-verbose text).
18. Ensure code passes Python linting: `make lint-python`

## Maintenance

When modifying existing scripts:

1. Update the corresponding Makefile target if needed
2. Update the GitHub workflow if behavior changes
3. Update this README if usage changes
4. Test locally with `make <target>` before committing
5. Ensure backward compatibility or update all references
6. Run `make lint-python` to ensure code quality standards are met
7. Add type hints to new functions or when modifying function signatures
8. Refactor large functions into smaller, focused helper functions when appropriate

## Utility Modules

Shared library details are documented in [`scripts/lib/README.md`](../scripts/lib/README.md).

## Related Documentation

- [Root Makefile](../Makefile) - Build and CI targets
- [GitHub Workflows](../.github/workflows/) - CI pipeline configuration
