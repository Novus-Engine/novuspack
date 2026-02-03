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


def _current_heading_result(definition):
    """Return (file, heading, anchor) for the definition's current heading."""
    anchor = generate_anchor_from_heading(definition.heading, include_hash=True)
    return (definition.file, definition.heading, anchor)


def _parse_canonical_link_target(link_target: str, definition) -> Tuple[str, str]:
    """Parse link target into (canonical_file, canonical_anchor)."""
    if "#" in link_target:
        file_part, anchor_part = link_target.split("#", 1)
        canonical_file = file_part if file_part else definition.file
        canonical_anchor = "#" + anchor_part
    else:
        canonical_file = link_target if link_target else definition.file
        canonical_anchor = ""
    return (canonical_file, canonical_anchor)


def _emit_canonical_error(
    output: Optional[OutputBuilder],
    definition,
    message: str,
    issue_message: str,
) -> None:
    """Emit a canonical validation error line if output is set."""
    if output:
        output.add_error_line(
            ValidationIssue.create(
                message,
                Path(definition.file),
                definition.heading_line,
                definition.heading_line,
                message=issue_message,
                severity="error",
            ).format_message(no_color=output.no_color)
        )


def _find_canonical_heading_in_file(
    content_cache: FileContentCache,
    canonical_file_path: Path,
    canonical_anchor: str,
) -> Optional[str]:
    """Return canonical heading text if anchor found in file, else None."""
    try:
        canonical_content = content_cache.get_content(canonical_file_path)
        canonical_headings = extract_headings(
            canonical_content, skip_code_blocks=True
        )
        anchor_without_hash = canonical_anchor.lstrip("#")
        for heading_text, _level, _line_num in canonical_headings:
            heading_anchor = generate_anchor_from_heading(
                heading_text, include_hash=True
            )
            if heading_anchor.lstrip("#") == anchor_without_hash:
                return heading_text
    except (OSError, UnicodeDecodeError, ValueError, RuntimeError):
        pass
    return None


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
        return _current_heading_result(definition)
    section_content_lower = definition.section_content.lower()
    if "this is the canonical" in section_content_lower:
        return _current_heading_result(definition)
    canonical_match = re.search(
        r"\bcanonical\b", definition.section_content, re.IGNORECASE
    )
    if not canonical_match:
        return _current_heading_result(definition)

    lines_before_match = definition.section_content[: canonical_match.start()].count("\n")
    canonical_line_start = definition.heading_line - 1 + lines_before_match
    content_lines = len(content.split("\n"))
    section_lines = definition.section_content.count("\n")
    search_end = min(content_lines, definition.heading_line - 1 + section_lines)
    search_lines = content.split("\n")[canonical_line_start:search_end + 1]
    search_content = "\n".join(search_lines)
    link_pattern = r"\[([^\]]+)\]\(([^)]+)\)"
    link_match = re.search(link_pattern, search_content)
    if not link_match:
        return _current_heading_result(definition)

    link_target = link_match.group(2)
    canonical_file, canonical_anchor = _parse_canonical_link_target(
        link_target, definition
    )
    if not validate_file_name(canonical_file):
        _emit_canonical_error(
            output,
            definition,
            "Invalid canonical link",
            f"Canonical link points to invalid file: {canonical_file}",
        )
        return _current_heading_result(definition)
    canonical_file_path = tech_specs_dir / canonical_file
    if not canonical_file_path.exists():
        _emit_canonical_error(
            output,
            definition,
            "Canonical file not found",
            f"Canonical link points to non-existent file: {canonical_file}",
        )
        return _current_heading_result(definition)
    if not is_safe_path(canonical_file_path, repo_root):
        _emit_canonical_error(
            output,
            definition,
            "Unsafe canonical path",
            f"Canonical link points outside repository: {canonical_file}",
        )
        return _current_heading_result(definition)

    content_cache = file_cache or FileContentCache()
    canonical_heading_text = _find_canonical_heading_in_file(
        content_cache, canonical_file_path, canonical_anchor
    )
    if canonical_heading_text:
        return (canonical_file, canonical_heading_text, canonical_anchor)
    _emit_canonical_error(
        output,
        definition,
        "Canonical anchor not found",
        f"Canonical anchor '{canonical_anchor}' not found in {canonical_file}",
    )
    return _current_heading_result(definition)
