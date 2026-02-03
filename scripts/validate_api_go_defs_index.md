# Validate API Go Definitions Index

## 1. Overview

This document describes the current business logic implemented by [validate_api_go_defs_index.py](validate_api_go_defs_index.py).
If this document and the implementation disagree, the implementation is authoritative.

The validator scans Go code blocks (` ```go `) in tech specs and ensures that all discovered Go API definitions are present in the Go definitions index.
This is intentionally scoped to Go API definitions only.
Other language code blocks are ignored.
Constants and variables are intentionally excluded.

## 2. Inputs and Outputs

Inputs:

- Index file (default): [docs/tech_specs/api_go_defs_index.md](../docs/tech_specs/api_go_defs_index.md).
- Tech specs directory (fixed): [docs/tech_specs](../docs/tech_specs).

Outputs:

- Structured output (errors, warnings, summary) to stdout, and optionally to `--output`.
- Exit code is determined by the output builder and is forced to 0 when `--no-fail` is set.
- `--apply` may write the index file, but it does not change the process exit code.

Command line options:

- `--verbose` / `-v`: Include verbose details (including placement details and the full expected tree at the end).
- `--index-file`: Override the index file path (relative to the repo root).
- `--output` / `-o FILE`: Write detailed output to `FILE`.
- `--no-color` / `--nocolor`: Disable colored output.
- `--no-fail`: Exit with code 0 even if errors are found.
- `--apply`: Apply high-confidence updates and reordering to the index file (interactive confirmation required).

## 3. ParsedIndex Model

ParsedIndex is produced by [scripts/lib/_index_utils.py](lib/_index_utils.py) via
[scripts/lib/go_defs_index/_go_defs_index_indexfile.py](lib/go_defs_index/_go_defs_index_indexfile.py).
It contains:

- sections: Map of section path to IndexSection.
- section_order: Ordered list of numbered section paths in index order.
- overview: ProseSection tree for the overview content.
- unsorted_paths: Path list for unsorted types, methods, and functions.
- title: Document title from the first H1.

During validation, each section may contain:

- current_entries: Entries that currently exist in the index file.
- expected_entries: Entries the validator expects based on discovered definitions and placement.

## 4. Validation Phases

Each phase builds on the data prepared by the previous phase.

### 4.1 Discovery

The validator scans all `*.md` files directly under [docs/tech_specs](../docs/tech_specs), excluding the index file itself.
It extracts Go definitions from ` ```go ` code blocks only.
Example code is excluded using heuristics (for example, example headings and single-line example checks).

Discovered definitions are normalized into `DetectedDefinition` objects.
Discovery supports:

- Types (including structs and interfaces).
- Methods (receiver methods).
- Functions, including extraction of referenced types and referenced methods from function signatures to help placement.

Notes and limitations:

- Some placement behavior is intentionally hard-coded (for example, receiver-specific method categorization rules and receiver normalization rules).
- This means the validator can drift from expectations as the API surface changes.
  - Renames or additions to receiver types can change placement confidence and increase unresolved entries.
  - New method families or reorganized sections may require updates to the categorization logic to avoid wrong-section suggestions.
- Sorting is deterministic once placement is complete, but placement quality depends on these rules and the current index structure.

Discovery also resolves canonical references (file, heading, anchor) for each definition.
If a definition name appears in multiple different tech spec files, discovery emits duplicate-definition errors.

### 4.2 Index Parsing

The validator reads the index file and parses it once into a `ParsedIndex`.
Each numbered section becomes an `IndexSection` with `current_entries` populated from the index file.
Entry description text is captured from the index file and attached to `IndexEntry` objects.

### 4.3 Placement

Placement is kind-first and uses confidence scoring.
Definitions are processed in this order: types, then methods, then functions.
The confidence threshold for high-confidence placement is 0.75.

Types:

- Types are placed only in type sections.
- If no high-confidence section match is found, the type is placed into the unsorted types section with status `unresolved`.

Methods:

- Methods require a receiver type.
- Methods are constrained to the receiver type section's method subsections (including nested method subsections).
- Receiver types may be normalized using implementation-to-interface mappings (for example, `ReadonlyPackage` and `FilePackage` are treated as `Package`).
- Category rules are applied for some receiver types (for example `Package` and `FileEntry`) to prefer structure-first placement when a category subsection exists.
- If no high-confidence match is found, the method is placed into the unsorted methods section with status `unresolved` and a suggested section (if any).

Functions:

- Functions may be placed into function subsections under a related type section when the function signature references types or methods.
- If no relationship is found, scoring selects among all function sections.
- If no high-confidence match is found, the function is placed into the unsorted functions section with status `unresolved` and a suggested section (if any).

### 4.4 Comparison

Comparison reconciles `expected_entries` against `current_entries` and sets statuses.
Comparison also moves expected entries out of unsorted sections if the same entry already exists in the index (so that section status can be evaluated in the correct section).

Statuses applied during comparison:

- For current index entries:
  - `orphaned`: The entry is not expected anywhere (it does not match any discovered definition).
  - `removed`: The entry exists, but the expected section for that entry is different from the current section.
  - `present`: The entry matches the expected section.
- For expected entries:
  - `added`: The entry is expected but does not exist in the current index.
  - `moved`: The entry exists in the index but in a different section.
  - `present`: The entry exists in the correct section.
  - `unresolved`: The entry could not be confidently placed (typically in an unsorted section).

Link update detection is performed by comparing expected and current link targets.
When a mismatch is detected, the current entry is marked as needing a link update and the expected link target is recorded.

### 4.5 Description Validation

Description validation enforces:

- Each indexed entry that is expected must have description text with a minimum length of 20 characters.
- Description text should be unique.
  Multiple entries sharing the exact same description are treated as errors.

When an entry is missing a description, the validator suggests using the definition's doc comments when available.

### 4.6 Ordering

Ordering performs two actions:

- Emits warnings when entries in the current index file are not in the expected alphabetical order.
  These are warnings (not errors) and are capped per section.
- Sorts expected entries within each section via `ParsedIndex.sort_expected_entries`.

### 4.7 Reporting

Reporting emits:

- Missing high-confidence definitions (expected entries with status `added`).
- Orphaned entries (current entries with status `orphaned`).
- Wrong-section entries (expected entries with status `moved`).
- Incorrect link targets (current entries flagged for link updates).
- Low-confidence unresolved definitions (expected entries with status `unresolved` and confidence < 75%).

When `--verbose` is set, the script also prints the full expected index tree from `ParsedIndex.render_full_tree`.
The final summary includes a breakdown of issue counts, and also includes counts of zero-confidence placements from discovery and placement.

## 5. Apply Switch

The `--apply` switch overwrites the index file with a full render of the expected structure.
The regenerated file includes the overview and table of contents.

Interactive confirmation requirements:

- The user must type `yes` to confirm.
- If stdin is a TTY, confirmation is read from stdin.
- If stdin is not a TTY, confirmation is read from `/dev/tty` when available.
- If neither is possible, `--apply` exits with code 1 and prints an error explaining that an interactive terminal is required.

Apply triggers:

- `--apply` only writes when there are pending high-confidence changes, including description fix candidates.
- If there are no pending changes, it prints `No high-confidence updates to apply.` and does not write.

Description application behavior:

- Before writing, the validator syncs expected descriptions and may populate missing expected entry descriptions derived from definition doc comments.
- This allows `--apply` to fix missing descriptions even when there are no structural or link changes, as long as suitable doc comments are available.

`--apply` is executed after output is printed and does not change the process exit code.

## 6. Related Files

- [scripts/validate_api_go_defs_index.py](validate_api_go_defs_index.py)
- [scripts/lib/_index_utils.py](lib/_index_utils.py)
- [scripts/lib/go_defs_index/_go_defs_index_discovery.py](lib/go_defs_index/_go_defs_index_discovery.py)
- [scripts/lib/go_defs_index/_go_defs_index_indexfile.py](lib/go_defs_index/_go_defs_index_indexfile.py)
- [scripts/lib/go_defs_index/_go_defs_index_matching.py](lib/go_defs_index/_go_defs_index_matching.py)
- [scripts/lib/go_defs_index/_go_defs_index_comparison.py](lib/go_defs_index/_go_defs_index_comparison.py)
- [scripts/lib/go_defs_index/_go_defs_index_descriptions.py](lib/go_defs_index/_go_defs_index_descriptions.py)
- [scripts/lib/go_defs_index/_go_defs_index_ordering.py](lib/go_defs_index/_go_defs_index_ordering.py)
- [scripts/lib/go_defs_index/_go_defs_index_reporting.py](lib/go_defs_index/_go_defs_index_reporting.py)
- [scripts/lib/go_defs_index/_go_defs_index_apply_parsing.py](lib/go_defs_index/_go_defs_index_apply_parsing.py) (used when applying index updates)

## 7. Validation Commands

- make validate-go-defs-index
- make docs-check
