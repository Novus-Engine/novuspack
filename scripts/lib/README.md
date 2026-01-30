# Scripts Lib Utilities

- [1. Purpose](#1-purpose)
- [2. Shared Utilities](#2-shared-utilities)
  - [2.1 Go Code Utils Module](#21-go-code-utils-module)
  - [2.2 Validation Utils Module](#22-validation-utils-module)
- [3. Go Definitions Index Subsystem](#3-go-definitions-index-subsystem)
- [4. Requirements Coverage Audit Subsystem](#4-requirements-coverage-audit-subsystem)
- [5. Tests](#5-tests)
- [6. Usage Notes](#6-usage-notes)

## 1. Purpose

This directory contains shared utility modules used by scripts in [`scripts/`](../).

## 2. Shared Utilities

These modules provide common parsing and validation helpers for the entrypoint scripts.

### 2.1 Go Code Utils Module

Shared utilities for parsing and processing Go code blocks in markdown files.

Location: [`scripts/lib/_go_code_utils.py`](../lib/_go_code_utils.py).

Core functionality:

- `find_go_code_blocks()` - Find Go code blocks in markdown content.
- `parse_go_def_signature()` - Parse Go definition signatures (functions, methods, types).
- `normalize_go_signature()` - Normalize Go signatures for comparison.
- `normalize_go_signature_with_params()` - Normalize signatures while preserving parameter names.
- `is_example_code()` - Detect example code blocks.
- `is_example_signature_name()` - Check if signature name indicates example code.
- `remove_go_comments()` - Remove Go comments from text.
- `determine_type_kind()` - Determine the kind of a Go type definition.
- `extract_receiver_type()` - Extract receiver type from receiver string.
- `normalize_generic_name()` - Normalize generic type names.
- `InterfaceParser` - Class for parsing Go interface definitions.
- `Signature` - Dataclass for representing Go signatures.

Used by:

- `validate_go_signature_sync.py`
- `validate_go_spec_signature_consistency.py`
- `validate_api_go_defs_index.py`
- `validate_go_code_blocks.py`

### 2.2 Validation Utils Module

Shared utilities for validation scripts, including output formatting, error reporting, and common helpers.

Location: [`scripts/lib/_validation_utils.py`](../lib/_validation_utils.py).

Core functionality:

- `OutputBuilder` - Class for building formatted validation output.
- `ValidationIssue` - Unified class for representing validation errors and warnings.
- `get_workspace_root()` - Find repository root directory.
- `find_markdown_files()` - Find markdown files in specified paths.
- `find_feature_files()` - Find feature files (.feature) in specified paths.
- `parse_paths()` - Parse comma-separated path arguments.
- `generate_anchor_from_heading()` - Generate GitHub-style markdown anchor from heading text.
- `remove_backticks_keep_content()` - Remove backticks from text while preserving content.
- `has_backticks()` - Check if text contains backticks.
- `get_backticks_error_message()` - Get standard error message for backticks in headings.
- `build_heading_hierarchy()` - Build hierarchical structure of headings in markdown files.
- `is_organizational_heading()` - Detect organizational headings with no direct content.
- `import_module_with_fallback()` - Import a module with fallback to file-based import.
- `extract_headings()` - Extract all headings from markdown content.
- `extract_headings_from_file()` - Extract all headings from a markdown file.
- `extract_headings_with_anchors()` - Extract headings and generate anchors with optional level filtering.
- `extract_h2_plus_headings_with_sections()` - Extract H2+ headings with anchors and section numbers.
- `extract_headings_with_section_numbers()` - Extract headings with section numbers, returns anchors set and sections dict.
- `is_safe_path()` - Check if a path is safe (within repo and no traversal).
- `validate_file_name()` - Validate that filename is safe (no path traversal, no separators).
- `validate_spec_file_name()` - Validate that spec file name is safe (includes .md extension check).
- `validate_anchor()` - Validate that anchor is safe (no path traversal, no separators).
- `FileContentCache` - Cache for file contents to avoid repeated reads.

Used by:

- [`scripts/apply_heading_corrections.py`](../apply_heading_corrections.py)
- [`scripts/audit_feature_coverage.py`](../audit_feature_coverage.py)
- [`scripts/audit_requirements_coverage.py`](../audit_requirements_coverage.py)
- [`scripts/generate_anchor.py`](../generate_anchor.py)
- [`scripts/validate_api_go_defs_index.py`](../validate_api_go_defs_index.py)
- [`scripts/validate_go_code_blocks.py`](../validate_go_code_blocks.py)
- [`scripts/validate_go_signature_sync.py`](../validate_go_signature_sync.py)
- [`scripts/validate_go_spec_references.py`](../validate_go_spec_references.py)
- [`scripts/validate_go_spec_signature_consistency.py`](../validate_go_spec_signature_consistency.py)
- [`scripts/validate_heading_numbering.py`](../validate_heading_numbering.py)
- [`scripts/validate_links.py`](../validate_links.py)
- [`scripts/validate_req_references.py`](../validate_req_references.py)

## 3. Go Definitions Index Subsystem

Modules in [`scripts/lib/go_defs_index/`](../lib/go_defs_index/) implement the phased pipeline for `validate_api_go_defs_index.py`.

### 3.1 Pipeline Modules

The pipeline is split into focused modules for discovery, parsing, placement, comparison, and reporting.

- **\_go_defs_index_discovery.py** - Locates spec files and extracts definitions from Go code blocks.
- **\_go_defs_index_indexfile.py** - Parses the index file into sections and entries.
- **\_go_defs_index_matching.py** - Places definitions into index sections and tracks unresolved entries.
- **\_go_defs_index_comparison.py** - Compares expected vs actual index entries and reports differences.
- **\_go_defs_index_descriptions.py** - Validates and suggests entry descriptions from comments.
- **\_go_defs_index_ordering.py** - Verifies ordering rules within sections.
- **\_go_defs_index_reporting.py** - Formats the final summary and detailed output.

### 3.2 Scoring Modules

Placement scoring is split into domain, text, and rule helpers.

- **\_go_defs_index_scoring.py** - Orchestrates scoring and aggregates reasoning.
- **\_go_defs_index_scoring_domain.py** - Domain detection and subsection keyword extraction helpers.
- **\_go_defs_index_scoring_text.py** - Comment/prose/heading keyword extraction and matching.
- **\_go_defs_index_scoring_rules_core_base.py** - Base scoring rules (kind checks, exact matches).
- **\_go_defs_index_scoring_rules_core_domain.py** - Domain scoring and constructor rules.
- **\_go_defs_index_scoring_rules_methods.py** - Method pattern and preference scoring.
- **\_go_defs_index_scoring_rules_sections.py** - Section and keyword matching rules.
- **\_go_defs_index_scoring_rules_penalties.py** - Penalties and special-case rules.

## 4. Requirements Coverage Audit Subsystem

Modules in [`scripts/lib/audit_requirements/`](../lib/audit_requirements/) implement the requirements-coverage audit used by `audit_requirements_coverage.py`.

- **\_audit_requirements_config.py** – Constants (exclusion patterns, keywords, prose thresholds, weights, compiled regexes).
- **\_audit_requirements_scan.py** – Spec/requirement discovery, index-file ref check, heading-ref check, requirement ref counts.
- **\_audit_requirements_classify.py** – Section content and heading classification (prose metrics, architectural score, classify_heading).

## 5. Tests

`_go_code_utils_test.py` provides a small harness for exercising parsing helpers in `_go_code_utils.py`.

Location: [`scripts/lib/_go_code_utils_test.py`](../lib/_go_code_utils_test.py).

## 6. Usage Notes

These modules are internal implementation details.

Do not call them directly from the command line.

Use the entrypoint scripts in [`scripts/`](../) or the corresponding Make targets.
