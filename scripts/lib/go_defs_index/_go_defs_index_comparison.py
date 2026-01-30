from __future__ import annotations

from typing import Optional

from lib._index_utils import ParsedIndex
from lib._validation_utils import OutputBuilder


def compare_with_index(
    parsed_index: ParsedIndex,
    output: Optional[OutputBuilder] = None,
) -> None:
    """
    Phase 4: Compare detected definitions with current index.
    """
    _ = output

    for path in parsed_index.unsorted_paths:
        section = parsed_index.sections.get(path)
        if not section:
            continue
        for name, entry in list(section.expected_entries.items()):
            current_section = parsed_index.find_section_by_current_entry(name)
            if not current_section:
                continue
            del section.expected_entries[name]
            current_section.expected_entries[name] = entry

    for section_path in parsed_index.section_order:
        section = parsed_index.sections.get(section_path)
        if not section:
            continue
        for name, entry in section.current_entries.items():
            expected_section = parsed_index.find_section_by_expected_entry(name)
            if not expected_section:
                entry.entry_status = "orphaned"
                continue
            if expected_section.path_label() != section_path:
                entry.entry_status = "removed"
                continue
            entry.entry_status = "present"

        for name, entry in section.expected_entries.items():
            current_section = parsed_index.find_section_by_current_entry(name)
            if not current_section:
                entry.entry_status = "added"
                continue
            if current_section.path_label() != section_path:
                entry.entry_status = "moved"
                continue
            entry.entry_status = "present"

        for name, expected_entry in section.expected_entries.items():
            current_entry = section.current_entries.get(name)
            if not current_entry:
                continue
            if expected_entry.link_target() != current_entry.link_target():
                current_entry.needs_link_update = True
                current_entry.expected_link_file = expected_entry.link_file
                current_entry.expected_link_anchor = expected_entry.link_anchor

    for path in parsed_index.unsorted_paths:
        section = parsed_index.sections.get(path)
        if not section:
            continue
        for entry in section.expected_entries.values():
            if entry.entry_status is None:
                entry.entry_status = "unresolved"
