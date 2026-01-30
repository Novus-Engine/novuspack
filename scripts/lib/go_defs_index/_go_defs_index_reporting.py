from __future__ import annotations

import re
from pathlib import Path
from lib._index_utils import ParsedIndex
from lib._validation_utils import OutputBuilder, ValidationIssue


def map_implementation_to_interface(receiver_type: str) -> str:
    """
    Map implementation types to their interface types.

    Returns the interface name if the receiver is an implementation,
    otherwise returns the receiver as-is.
    """
    if not receiver_type:
        return receiver_type

    receiver_lower = receiver_type.lower()

    # Package interface implementations
    if receiver_lower in ["readonlypackage", "filepackage"]:
        return "Package"

    return receiver_type


def normalize_section_for_sort(section_name: str) -> str:
    """
    Normalize a section name for numeric sorting.

    Extracts numeric prefixes and pads them for proper numeric ordering.
    Example: "1. Core" => "0001. Core", "10. Package" => "0010. Package".
    Handles nested sections like "10. Package > 10.1 Comment" =>
    "0010. Package > 0010.0001. Comment".
    """
    if " > " in section_name:
        parts = section_name.split(" > ")
        normalized_parts = [normalize_section_for_sort(part) for part in parts]
        return " > ".join(normalized_parts)

    match = re.match(r"^(\d+)\.(?:(\d+)\s+)?(.*)$", section_name)
    if match:
        first_num = int(match.group(1))
        second_num = match.group(2)
        rest = match.group(3)
        if second_num:
            second_num_int = int(second_num)
            return f"{first_num:04d}.{second_num_int:04d} {rest}"
        return f"{first_num:04d}. {rest}"

    return section_name


def get_section_display_name(section_name: str) -> str:
    """
    Get the display name for a section.

    For nested sections like "3. Generics > 3.1 Core Generic Types",
    returns just "3.1 Core Generic Types".
    For top-level sections, returns as-is.
    """
    if " > " in section_name:
        return section_name.split(" > ")[-1]
    return section_name


def generate_report(
    parsed_index: "ParsedIndex",
    output: OutputBuilder,
    index_file_name: str,
) -> None:
    """
    Generate output showing all changes needed.
    """
    added_entries = parsed_index.get_added_entries()
    moved_entries = parsed_index.get_moved_entries()
    orphaned_entries = parsed_index.get_orphans()
    removed_entries = parsed_index.get_removed_entries()
    link_updates = parsed_index.get_link_update_entries()
    unresolved_entries = parsed_index.get_unresolved_entries()
    has_issues = (
        added_entries
        or moved_entries
        or orphaned_entries
        or removed_entries
        or link_updates
        or unresolved_entries
    )
    if not has_issues:
        return

    output.add_errors_header()

    if added_entries:
        output.add_blank_line("error")
        output.add_line(
            f"Found {len(added_entries)} high-confidence sorted definition(s) not in index:",
            section="error",
        )
        output.add_blank_line("error")
        for entry in sorted(added_entries, key=lambda item: item.name):
            section_str = entry.current_section or "(unresolved)"
            confidence_str = (
                f"{int(entry.confidence_score * 100)}%"
                if entry.confidence_score is not None
                else "N/A"
            )
            canonical_str = entry.link_target() if entry.link_file else "(no canonical link)"
            output.add_error_line(f"  {entry.name}")
            output.add_error_line(f"    - Kind: {entry.kind}")
            if entry.source_file:
                output.add_error_line(f"    - File: {entry.source_file}:{entry.source_line}")
            output.add_error_line(
                f"    - Suggested section: {section_str} (confidence: {confidence_str})"
            )
            output.add_error_line(f"    - Canonical location: {canonical_str}")
        output.add_blank_line("error")

    if orphaned_entries:
        output.add_line(
            f"Found {len(orphaned_entries)} orphaned entry/entries in index:",
            section="error",
        )
        output.add_blank_line("error")
        for entry in sorted(orphaned_entries, key=lambda e: e.name):
            output.add_error_line(
                ValidationIssue(
                    "Orphaned entry",
                    Path(index_file_name),
                    entry.line_number,
                    entry.line_number,
                    f"`{entry.name}` not found in any tech spec file",
                    severity="error",
                ).format_message(no_color=output.no_color)
            )
        output.add_blank_line("error")

    if moved_entries:
        output.add_line(
            f"Found {len(moved_entries)} definition(s) in wrong section:",
            section="error",
        )
        output.add_blank_line("error")
        for entry in sorted(moved_entries, key=lambda e: e.name):
            current_section = parsed_index.find_section_by_current_entry(entry.name)
            current_path = current_section.path_label() if current_section else "(unknown)"
            output.add_error_line(
                ValidationIssue(
                    "Wrong section",
                    Path(index_file_name),
                    entry.line_number,
                    entry.line_number,
                    f"`{entry.name}` in '{current_path}'",
                    severity="error",
                    suggestion=f"Move to '{entry.current_section}'",
                ).format_message(no_color=output.no_color)
            )
        output.add_blank_line("error")

    if link_updates:
        output.add_line(
            f"Found {len(link_updates)} entry/entries with incorrect links:",
            section="error",
        )
        output.add_blank_line("error")
        for entry in sorted(link_updates, key=lambda item: item.name):
            current_link = entry.link_target()
            suggested_link = entry.expected_link_file or ""
            if entry.expected_link_anchor:
                suggested_link = f"{suggested_link}#{entry.expected_link_anchor}"
            output.add_error_line(
                ValidationIssue(
                    "Incorrect link",
                    Path(index_file_name),
                    entry.line_number,
                    entry.line_number,
                    f"`{entry.name}`: {current_link}",
                    severity="error",
                    suggestion=f"Update to: {suggested_link}",
                ).format_message(no_color=output.no_color)
            )
        output.add_blank_line("error")

    if unresolved_entries:
        output.add_line(
            (
                f"Found {len(unresolved_entries)} definition(s) with low confidence "
                f"(< 75%) not in index:"
            ),
            section="error",
        )
        output.add_blank_line("error")
        for entry in sorted(unresolved_entries, key=lambda item: item.name):
            confidence_str = (
                f"{int(entry.confidence_score * 100)}%"
                if entry.confidence_score is not None
                else "0%"
            )
            reasoning_str = (
                ", ".join(entry.confidence_reasoning)
                if entry.confidence_reasoning
                else "no matches"
            )
            suggested = entry.suggested_section or "(unresolved)"
            output.add_error_line(f"  {entry.name}")
            if entry.source_file:
                output.add_error_line(f"    - File: {entry.source_file}:{entry.source_line}")
            output.add_error_line(
                f"    - Suggested section: {suggested} (confidence: {confidence_str})"
            )
            output.add_error_line(f"    - Reasoning: {reasoning_str}")
            output.add_error_line(
                "    - Manual review required - confidence too low "
                "for automatic placement"
            )
        output.add_blank_line("error")

    output.add_blank_line("error")
    output.add_line("Expected index (full tree):", section="error")
    output.add_blank_line("error")
    for line in parsed_index.render_full_tree():
        output.add_error_line(line)
    output.add_blank_line("error")

# End of module.
