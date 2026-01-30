from __future__ import annotations

import re
from pathlib import Path
from typing import Dict, List, Optional

from lib._index_utils import IndexEntry, ParsedIndex
from lib._validation_utils import OutputBuilder, ValidationIssue


def check_entry_descriptions(
    parsed_index: ParsedIndex,
    index_file: Path,
    output: Optional[OutputBuilder] = None,
) -> int:
    """
    Phase 5: Validate that all index entries have descriptive text (>= 20 characters).

    Returns:
        Number of errors found.
    """
    error_count = 0

    if not output:
        return 0

    # Track descriptions for duplicate detection.
    description_to_entries: Dict[str, List[str]] = {}

    index_entries: Dict[str, IndexEntry] = {}
    expected_entries: Dict[str, IndexEntry] = {}
    for section_path in parsed_index.section_order:
        section = parsed_index.sections.get(section_path)
        if not section:
            continue
        index_entries.update(section.current_entries)
        expected_entries.update(section.expected_entries)

    for entry_name, entry in index_entries.items():
        if entry_name not in expected_entries:
            continue

        if not entry.has_description:
            suggested_entry = expected_entries.get(entry_name)
            comments = suggested_entry.def_comments if suggested_entry else None
            suggestion_parts: List[str] = []
            if comments:
                suggestion_parts.append("Review the definition comments for potential summary:")
                suggestion_parts.append(comments)
            else:
                suggestion_parts.append(
                    "Add descriptive text (minimum 20 characters) below this entry."
                )

            suggestion = "\n".join(suggestion_parts)
            suggestion = re.sub(r"\s+", " ", suggestion).strip()

            if entry.description_text is None:
                error_type = "Missing description"
                error_msg_text = (
                    f"Entry `{entry_name}` has no descriptive text "
                    f"(minimum 20 characters required)"
                )
            else:
                error_type = "Description too short"
                error_msg_text = (
                    f"Entry `{entry_name}` has descriptive text that is too short "
                    f"({len(entry.description_text)} characters, minimum 20 required)"
                )

            output.add_error_line(
                ValidationIssue(
                    error_type,
                    index_file,
                    entry.line_number,
                    entry.line_number,
                    error_msg_text,
                    severity="error",
                    suggestion=suggestion,
                ).format_message(no_color=output.no_color)
            )
            error_count += 1

        if entry.has_description and entry.description_text:
            desc_text = entry.description_text.strip()
            description_to_entries.setdefault(desc_text, []).append(entry_name)

    for desc_text, entry_names in description_to_entries.items():
        if len(entry_names) <= 1:
            continue

        entries_list = ", ".join(f"`{name}`" for name in entry_names)
        first_entry = index_entries[entry_names[0]]
        output.add_error_line(
            ValidationIssue(
                "Duplicate description",
                index_file,
                first_entry.line_number,
                first_entry.line_number,
                f"Multiple entries share the same description: {entries_list}",
                severity="error",
                suggestion="Each entry should have a unique description",
            ).format_message(no_color=output.no_color)
        )
        error_count += 1

    return error_count
