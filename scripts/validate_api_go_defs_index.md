# Validate API Go Definitions Index

## 1. Overview

This document describes the business logic for [validate_api_go_defs_index.py](validate_api_go_defs_index.py).
It is the source of truth for how the Go definitions index is parsed, compared, ordered, reported, and optionally updated.

## 2. Inputs and Outputs

Inputs:

- Index file: [docs/tech_specs/api_go_defs_index.md](../docs/tech_specs/api_go_defs_index.md).
- Tech specs directory: [docs/tech_specs](../docs/tech_specs).

Outputs:

- Structured report to stdout or to the optional output file.
- Exit code 0 on success, 1 on validation failure unless --no-fail is set.

## 3. ParsedIndex Model

ParsedIndex is produced by [scripts/lib/_index_utils.py](lib/_index_utils.py).
It contains:

- sections: Map of section path to IndexSection.
- section_order: Ordered list of numbered section paths in index order.
- overview: ProseSection tree for the overview content.
- unsorted_paths: Path list for unsorted types, methods, and functions.
- title: Document title from the first H1.

## 4. Validation Phases

Each phase builds on the data prepared by the previous phase.

### 4.1 Discovery

The validator scans all markdown in [docs/tech_specs](../docs/tech_specs) except the index file.
It extracts Go definitions from ` ```go ` code blocks and creates DetectedDefinition objects.
Heading resolution is performed during discovery, so each definition includes canonical file and anchor data.
`api_file_mgmt_errors.md` participates in discovery like other tech specs.

### 4.2 Index Parsing

The validator parses the index once using [scripts/lib/_index_utils.py](lib/_index_utils.py).
Each numbered section becomes an IndexSection with current_entries populated from the index file.
Entry description lines are captured and attached to the current IndexEntry objects.

### 4.3 Placement

Placement is structure-first and kind-first.
Definitions are placed in this order: types, then methods, then functions.

Types:

- Types are placed only in type sections.
- Unresolvable types go into unsorted types.

Methods:

- Methods are constrained only to the receiver type section's child method subsections.
- Structure-first categorization picks the exact subsection when possible.
- Signature-related Package methods fall back to Package Other Methods when no signature subsection exists.
- Unresolvable methods go into unsorted methods.

Functions:

- Functions are placed using referenced types or methods first.
- If no relationship is found, scoring chooses among helper sections.
- Scoring logic is split into focused modules under `scripts/lib/go_defs_index/`.
- Unresolvable functions go into unsorted functions.

### 4.4 Comparison

Comparison reconciles unsorted expected entries with existing current entries.
It sets entry_status for current and expected entries as added, moved, removed, orphaned, present, or unresolved.
Link updates are detected by comparing expected and current link targets.

### 4.5 Description Validation

Description checks require at least 20 characters of description text per entry.
If a description is missing, def comments from expected entries are used as suggestions.

### 4.6 Ordering

Ordering sorts expected entries within each section.
Sorting uses ParsedIndex.sort_expected_entries with capitals-first, case-insensitive ordering.

### 4.7 Reporting

Reporting uses entry_status to list missing, moved, orphaned, and link issues.
Low-confidence entries are reported from unsorted expected entries.
The full expected tree is rendered from ParsedIndex.render_full_tree.

## 5. Apply Switch

The --apply switch overwrites the index file with a full render of the expected structure.
It requires a TTY and an explicit "yes" confirmation before writing.
The regeneration includes the overview and table of contents.
Unsorted sections are never applied and unresolved items are excluded by design.

## 6. Related Files

- [scripts/validate_api_go_defs_index.py](validate_api_go_defs_index.py)
- [scripts/lib/_index_utils.py](lib/_index_utils.py)
- [scripts/lib/go_defs_index/_go_defs_index_matching.py](lib/go_defs_index/_go_defs_index_matching.py)
- [scripts/lib/go_defs_index/_go_defs_index_comparison.py](lib/go_defs_index/_go_defs_index_comparison.py)
- [scripts/lib/go_defs_index/_go_defs_index_ordering.py](lib/go_defs_index/_go_defs_index_ordering.py)
- [scripts/lib/go_defs_index/_go_defs_index_reporting.py](lib/go_defs_index/_go_defs_index_reporting.py)
- [scripts/lib/go_defs_index/_go_defs_index_config.py](lib/go_defs_index/_go_defs_index_config.py)
- [scripts/lib/go_defs_index/_go_defs_index_scoring.py](lib/go_defs_index/_go_defs_index_scoring.py)
- [scripts/lib/go_defs_index/_go_defs_index_scoring_domain.py](lib/go_defs_index/_go_defs_index_scoring_domain.py)
- [scripts/lib/go_defs_index/_go_defs_index_scoring_text.py](lib/go_defs_index/_go_defs_index_scoring_text.py)
- [scripts/lib/go_defs_index/_go_defs_index_scoring_rules_core.py](lib/go_defs_index/_go_defs_index_scoring_rules_core.py)
- [scripts/lib/go_defs_index/_go_defs_index_scoring_rules_core_base.py](lib/go_defs_index/_go_defs_index_scoring_rules_core_base.py)
- [scripts/lib/go_defs_index/_go_defs_index_scoring_rules_core_domain.py](lib/go_defs_index/_go_defs_index_scoring_rules_core_domain.py)
- [scripts/lib/go_defs_index/_go_defs_index_scoring_rules_methods.py](lib/go_defs_index/_go_defs_index_scoring_rules_methods.py)
- [scripts/lib/go_defs_index/_go_defs_index_scoring_rules_sections.py](lib/go_defs_index/_go_defs_index_scoring_rules_sections.py)
- [scripts/lib/go_defs_index/_go_defs_index_scoring_rules_penalties.py](lib/go_defs_index/_go_defs_index_scoring_rules_penalties.py)

## 7. Validation Commands

- make validate-go-defs-index
- make docs-check
