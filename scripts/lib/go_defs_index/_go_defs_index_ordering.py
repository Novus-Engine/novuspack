from __future__ import annotations

from pathlib import Path
from typing import Optional
from lib._index_utils import ParsedIndex
from lib._validation_utils import OutputBuilder, ValidationIssue


def determine_ordering(
    parsed_index: ParsedIndex,
    output: Optional[OutputBuilder] = None,
) -> None:
    """
    Phase 6: Determine correct ordering within sections.
    """
    if output:
        _emit_ordering_warnings(parsed_index, output)
    parsed_index.sort_expected_entries()


def _emit_ordering_warnings(
    parsed_index: ParsedIndex,
    output: OutputBuilder,
) -> None:
    max_warnings_per_section = 5
    index_path = Path("api_go_defs_index.md")
    for section_path in parsed_index.section_order:
        section = parsed_index.sections.get(section_path)
        if not section or len(section.entries) < 2:
            continue
        expected = sorted(section.entries, key=lambda entry: entry.sort_key())
        if [entry.name for entry in section.entries] == [entry.name for entry in expected]:
            continue
        mismatches = 0
        for idx, entry in enumerate(section.entries):
            if idx >= len(expected):
                break
            expected_entry = expected[idx]
            if entry.name == expected_entry.name:
                continue
            if entry.entry_status in (None, "present"):
                entry.entry_status = "reordered"
            expected_index_entry = section.expected_entries.get(entry.name)
            if expected_index_entry and expected_index_entry.entry_status in (None, "present"):
                expected_index_entry.entry_status = "reordered"
            suggestion = (
                "Reorder entries to maintain alphabetical ordering by name."
            )
            message = (
                f"`{entry.raw_name}` appears before `{expected_entry.raw_name}`"
            )
            output.add_warning_line(
                ValidationIssue.create(
                    "Incorrect entry order",
                    index_path,
                    entry.line_number,
                    entry.line_number,
                    message=message,
                    severity="warning",
                    suggestion=suggestion,
                ).format_message(no_color=output.no_color)
            )
            mismatches += 1
            if mismatches >= max_warnings_per_section:
                output.add_warning_line(
                    f"WARNING: Additional ordering issues in '{section_path}' omitted."
                )
                break
