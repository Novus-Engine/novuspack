from __future__ import annotations

import re
from pathlib import Path
from typing import Optional, Tuple

from lib._validation_utils import (
    FileContentCache,
    OutputBuilder,
    ValidationIssue,
    extract_headings,
    generate_anchor_from_heading,
    get_workspace_root,
    is_safe_path,
    validate_file_name,
)


def resolve_canonical_reference(
    definition,
    content: str,
    tech_specs_dir: Path,
    output: Optional[OutputBuilder] = None,
    file_cache: Optional[FileContentCache] = None,
) -> Tuple[str, str, str]:
    """
    Resolve canonical reference if present.

    Only follows one link step and validates paths stay within repository.

    Returns:
        (canonical_file, canonical_heading, canonical_anchor)
    """
    repo_root = get_workspace_root()

    if not definition.section_content:
        # No section content, use current heading
        anchor = generate_anchor_from_heading(definition.heading, include_hash=True)
        return (definition.file, definition.heading, anchor)

    section_content_lower = definition.section_content.lower()

    # Exception: if section contains "this is the canonical", treat as canonical
    if "this is the canonical" in section_content_lower:
        anchor = generate_anchor_from_heading(definition.heading, include_hash=True)
        return (definition.file, definition.heading, anchor)

    # Search for "canonical" keyword (case-insensitive)
    canonical_match = re.search(r"\bcanonical\b", definition.section_content, re.IGNORECASE)
    if not canonical_match:
        # No canonical reference, use current heading
        anchor = generate_anchor_from_heading(definition.heading, include_hash=True)
        return (definition.file, definition.heading, anchor)

    # Find link on same line or following lines
    # Calculate line offset of canonical match
    lines_before_match = definition.section_content[: canonical_match.start()].count("\n")
    canonical_line_start = definition.heading_line - 1 + lines_before_match  # 0-indexed

    # Search from canonical line to end of section for markdown link
    search_start = canonical_line_start
    content_lines = len(content.split("\n"))
    section_lines = definition.section_content.count("\n")
    search_end = min(content_lines, definition.heading_line - 1 + section_lines)
    search_lines = content.split("\n")[search_start:search_end + 1]
    search_content = "\n".join(search_lines)

    link_pattern = r"\[([^\]]+)\]\(([^)]+)\)"
    link_match = re.search(link_pattern, search_content)

    if link_match:
        link_target = link_match.group(2)  # e.g., "api_core.md#20-package-interface"

        # Parse link target
        if "#" in link_target:
            file_part, anchor_part = link_target.split("#", 1)
            canonical_file = file_part if file_part else definition.file
            canonical_anchor = "#" + anchor_part
        else:
            canonical_file = link_target if link_target else definition.file
            canonical_anchor = ""

        # Validate file name
        if not validate_file_name(canonical_file):
            if output:
                output.add_error_line(
                    ValidationIssue(
                        "Invalid canonical link",
                        Path(definition.file),
                        definition.heading_line,
                        definition.heading_line,
                        f"Canonical link points to invalid file: {canonical_file}",
                        severity="error",
                    ).format_message(no_color=output.no_color)
                )
            anchor = generate_anchor_from_heading(definition.heading, include_hash=True)
            return (definition.file, definition.heading, anchor)

        # Resolve file path
        canonical_file_path = tech_specs_dir / canonical_file
        if not canonical_file_path.exists():
            if output:
                output.add_error_line(
                    ValidationIssue(
                        "Canonical file not found",
                        Path(definition.file),
                        definition.heading_line,
                        definition.heading_line,
                        f"Canonical link points to non-existent file: {canonical_file}",
                        severity="error",
                    ).format_message(no_color=output.no_color)
                )
            anchor = generate_anchor_from_heading(definition.heading, include_hash=True)
            return (definition.file, definition.heading, anchor)

        # Validate path is within repository
        if not is_safe_path(canonical_file_path, repo_root):
            if output:
                output.add_error_line(
                    ValidationIssue(
                        "Unsafe canonical path",
                        Path(definition.file),
                        definition.heading_line,
                        definition.heading_line,
                        f"Canonical link points outside repository: {canonical_file}",
                        severity="error",
                    ).format_message(no_color=output.no_color)
                )
            anchor = generate_anchor_from_heading(definition.heading, include_hash=True)
            return (definition.file, definition.heading, anchor)

        content_cache = file_cache or FileContentCache()

        # Read target file and find heading for anchor
        try:
            canonical_content = content_cache.get_content(canonical_file_path)
            canonical_headings = extract_headings(
                canonical_content, skip_code_blocks=True
            )

            # Find heading that matches the anchor
            anchor_without_hash = canonical_anchor.lstrip("#")
            canonical_heading_text = None

            for heading_text, level, line_num in canonical_headings:
                heading_anchor = generate_anchor_from_heading(
                    heading_text, include_hash=True
                )
                if heading_anchor.lstrip("#") == anchor_without_hash:
                    canonical_heading_text = heading_text
                    break

            if canonical_heading_text:
                return (canonical_file, canonical_heading_text, canonical_anchor)
            else:
                # Anchor not found, use current heading
                if output:
                    output.add_error_line(
                        ValidationIssue(
                            "Canonical anchor not found",
                            Path(definition.file),
                            definition.heading_line,
                            definition.heading_line,
                            f"Canonical anchor '{canonical_anchor}' not found in {canonical_file}",
                            severity="error",
                        ).format_message(no_color=output.no_color)
                    )
                anchor = generate_anchor_from_heading(
                    definition.heading, include_hash=True
                )
                return (definition.file, definition.heading, anchor)

        except Exception as e:
            if output:
                output.add_error_line(
                    ValidationIssue(
                        "Error reading canonical file",
                        Path(definition.file),
                        definition.heading_line,
                        definition.heading_line,
                        f"Could not read canonical file {canonical_file}: {e}",
                        severity="error",
                    ).format_message(no_color=output.no_color)
                )
            anchor = generate_anchor_from_heading(
                definition.heading, include_hash=True
            )
            return (definition.file, definition.heading, anchor)

    # No link found after "canonical", use current heading
    anchor = generate_anchor_from_heading(definition.heading, include_hash=True)
    return (definition.file, definition.heading, anchor)
