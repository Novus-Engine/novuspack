from __future__ import annotations

import re
from pathlib import Path
from typing import Dict, List, Optional, Tuple

from lib._go_code_utils import normalize_generic_name
from lib._validation_utils import (
    FileContentCache,
    ValidationIssue,
    extract_h2_plus_headings_with_sections,
    generate_anchor_from_heading,
    is_safe_path,
    validate_anchor,
    validate_file_name,
)


def _split_section_number_from_heading(heading_text: str) -> Tuple[Optional[str], str]:
    """
    Split a heading like \"2.1 AddFile Package Method\" into (\"2.1\", \"AddFile Package Method\").
    If no section prefix exists, returns (None, heading_text).
    """
    m = re.match(r"^(\\d+(?:\\.\\d+)*)\\s+(.+)$", heading_text)
    if not m:
        return (None, heading_text)
    return (m.group(1).strip(), m.group(2).strip())


def _find_heading_before_line(
    headings: List[Tuple[int, str, int, str, Optional[str]]],
    line_num: int,
) -> Optional[Tuple[int, str, int, str, Optional[str]]]:
    """
    Given extract_h2_plus_headings_with_sections output, find the last heading
    with line_num <= target line.
    """
    best: Optional[Tuple[int, str, int, str, Optional[str]]] = None
    for h in headings:
        _level, _text, h_line, _anchor, _section_anchor = h
        if h_line <= line_num:
            best = h
        else:
            break
    return best


def find_section_for_definition(
    definition_name: str,
    target_file: str,
    all_definitions: Dict[str, List[Tuple[str, int]]],
    tech_specs_dir: Path,
    file_cache: Optional[FileContentCache] = None,
) -> Optional[Tuple[Optional[str], str]]:
    """
    Find the section heading for a definition in a target file.

    Returns:
        (section_number, heading_text_without_section_prefix) or None.
    """
    if definition_name not in all_definitions:
        return None

    if not validate_file_name(target_file):
        return None

    target_path = tech_specs_dir / target_file
    repo_root = tech_specs_dir.parent.parent
    if not is_safe_path(target_path, repo_root):
        return None
    if not target_path.exists():
        return None

    definition_line: Optional[int] = None
    for filename, line_num in all_definitions[definition_name]:
        if filename == target_file:
            definition_line = line_num
            break
    if definition_line is None:
        return None

    # Use shared heading extractor with section anchors.
    headings = extract_h2_plus_headings_with_sections(
        target_path, skip_code_blocks=True, file_cache=file_cache
    )
    if not headings:
        return None

    heading = _find_heading_before_line(headings, definition_line)
    if not heading:
        return None

    _level, heading_text, _line, _anchor, _section_anchor = heading
    section_num, heading_without_section = _split_section_number_from_heading(heading_text)
    return (section_num, heading_without_section)


def check_missing_section_anchors(
    index_content: str,
    entry_to_target_md: Dict[str, str],
    _entry_to_link_anchor: Dict[str, Optional[str]],
    all_definitions: Dict[str, List[Tuple[str, int]]],
    tech_specs_dir: Path,
    *,
    index_filename: str = "api_go_defs_index.md",
    file_cache: Optional[FileContentCache] = None,
) -> List[ValidationIssue]:
    """
    Check for index entries that point to files without section anchors.

    Returns list of ValidationIssue objects.
    """
    issues: List[ValidationIssue] = []
    lines = index_content.split("\n")
    entry_pattern = r"^\\s*-\\s+\\*\\*`([^`]+)`\\*\\*"
    link_pattern = r"\\[([^\\]]+)\\]\\(([^)]+)\\)"

    for line_num, line in enumerate(lines, 1):
        entry_match = re.match(entry_pattern, line)
        if not entry_match:
            continue

        entry_name = normalize_generic_name(entry_match.group(1))
        target_file = entry_to_target_md.get(entry_name)
        if not target_file:
            continue

        link_match = re.search(link_pattern, line)
        if not link_match:
            continue

        current_link_full = link_match.group(2)
        current_anchor = None
        if "#" in current_link_full:
            current_file_part, current_anchor = current_link_full.split("#", 1)
            if not validate_anchor(current_anchor):
                continue
        else:
            current_file_part = current_link_full

        if not validate_file_name(current_file_part):
            continue

        section_info = find_section_for_definition(
            entry_name,
            target_file,
            all_definitions,
            tech_specs_dir,
            file_cache=file_cache,
        )
        if not section_info:
            continue

        section_num, heading_text = section_info
        if section_num:
            section_anchor = section_num.replace(".", "")
            heading_anchor = generate_anchor_from_heading(heading_text, include_hash=False)
            correct_anchor = f"{section_anchor}-{heading_anchor}"
        else:
            correct_anchor = generate_anchor_from_heading(heading_text, include_hash=False)

        if not validate_anchor(correct_anchor):
            continue
        if not validate_file_name(target_file):
            continue

        if not current_anchor or current_anchor != correct_anchor:
            if current_anchor and not validate_anchor(current_anchor):
                continue
            link_text = link_match.group(1)
            suggested_link = f"[{link_text}]({target_file}#{correct_anchor})"
            current_link = link_match.group(0)
            index_file_path = tech_specs_dir / index_filename
            issues.append(
                ValidationIssue.create(
                    "incorrect_section_anchor",
                    index_file_path,
                    line_num,
                    line_num,
                    message=f"Index entry '{entry_name}' has incorrect or missing section anchor",
                    severity="error",
                    suggestion=suggested_link,
                    entry_name=entry_name,
                    target_file=target_file,
                    current_link=current_link,
                    correct_anchor=correct_anchor,
                )
            )

    return issues


def _scan_block_for_definition(
    lines_content: List[str],
    anchor_section_line: int,
    next_heading_line: int,
    code_block_start_line: int,
    entry_name: str,
) -> tuple:
    """Scan section for Go block and definition name. Returns (found, in_block, end_line)."""
    definition_found_in_block = False
    in_target_block = False
    code_block_end: Optional[int] = None
    end_row = min(next_heading_line - 1, len(lines_content))
    for i in range(anchor_section_line - 1, end_row):
        line_text = lines_content[i]
        current_line_num = i + 1
        if current_line_num == code_block_start_line:
            in_target_block = line_text.strip() == "```go"
            continue
        if in_target_block:
            if line_text.strip() == "```":
                code_block_end = current_line_num
                break
            if "." in entry_name:
                _receiver, method = entry_name.split(".", 1)
                if (entry_name in line_text) or re.search(
                    rf"\\b{re.escape(method)}\\s*\\(", line_text
                ):
                    definition_found_in_block = True
            else:
                if re.search(rf"\\b{re.escape(entry_name)}\\b", line_text):
                    definition_found_in_block = True
    return (definition_found_in_block, in_target_block, code_block_end)


def _find_anchor_section_line(headings, anchor: str) -> Optional[int]:
    """Return heading line number for matching anchor, or None."""
    for (
        _heading_level,
        _heading_text,
        heading_line,
        heading_anchor,
        section_anchor,
    ) in headings:
        if section_anchor and section_anchor == anchor:
            return heading_line
        if heading_anchor == anchor:
            return heading_line
    return None


def _validate_entry_and_get_definition_line(
    entry_name: str,
    target_file: Optional[str],
    anchor: Optional[str],
    all_definitions: Dict[str, List[Tuple[str, int]]],
    index_file_path: Path,
    *,
    line_num: int,
) -> Tuple[Optional[ValidationIssue], Optional[int]]:
    """Return (issue, None) if invalid, else (None, definition_line_in_target)."""
    if not target_file or not anchor:
        return (None, None)
    if not validate_file_name(target_file) or not validate_anchor(anchor):
        return (None, None)
    if entry_name not in all_definitions:
        return (
            ValidationIssue.create(
                "definition_not_found",
                index_file_path,
                line_num,
                line_num,
                message=f"Definition '{entry_name}' not found in any tech spec file",
                severity="error",
                entry_name=entry_name,
                target_file=target_file,
                anchor=anchor,
            ),
            None,
        )
    definition_line_in_target = None
    for filename, def_line in all_definitions[entry_name]:
        if filename == target_file:
            definition_line_in_target = def_line
            break
    if definition_line_in_target is None:
        return (
            ValidationIssue.create(
                "definition_not_in_target",
                index_file_path,
                line_num,
                line_num,
                message=f"Definition '{entry_name}' not found in target file '{target_file}'",
                severity="error",
                entry_name=entry_name,
                target_file=target_file,
                anchor=anchor,
            ),
            None,
        )
    return (None, definition_line_in_target)


def _issue_if_target_missing(
    target_path: Path,
    repo_root: Path,
    index_file_path: Path,
    line_num: int,
    entry_name: str,
    *,
    target_file: str,
    anchor: str,
):
    """Return a single ValidationIssue if target path is unsafe or missing, else None."""
    if not is_safe_path(target_path, repo_root):
        return None
    if target_path.exists():
        return None
    return ValidationIssue.create(
        "target_file_not_found",
        index_file_path,
        line_num,
        line_num,
        message=f"Target file '{target_file}' does not exist",
        severity="error",
        entry_name=entry_name,
        target_file=target_file,
        anchor=anchor,
    )


def _issue_if_anchor_or_position_bad(
    lines_content: List[str],
    headings: list,
    anchor: str,
    definition_line_in_target: int,
    index_file_path: Path,
    *,
    line_num: int,
    entry_name: str,
    target_file: str,
) -> Tuple[Optional[ValidationIssue], Optional[int], Optional[int]]:
    """Return (issue, None, None) if bad, else (None, anchor_section_line, next_heading_line)."""
    anchor_section_line = _find_anchor_section_line(headings, anchor)
    if anchor_section_line is None:
        return (
            ValidationIssue.create(
                "anchor_no_match",
                index_file_path,
                line_num,
                line_num,
                message=f"Anchor '{anchor}' does not match any heading in target file",
                severity="error",
                entry_name=entry_name,
                target_file=target_file,
                anchor=anchor,
            ),
            None,
            None,
        )
    next_heading_line = len(lines_content) + 1
    for _level, _text, h_line, _a, _sa in headings:
        if h_line > anchor_section_line:
            next_heading_line = h_line
            break
    if definition_line_in_target < anchor_section_line:
        return (
            ValidationIssue.create(
                "definition_before_anchor",
                index_file_path,
                line_num,
                line_num,
                message=(
                    f"Definition at line {definition_line_in_target} is before anchor "
                    f"section at line {anchor_section_line}"
                ),
                severity="error",
                entry_name=entry_name,
                target_file=target_file,
                anchor=anchor,
                definition_line=definition_line_in_target,
                anchor_line=anchor_section_line,
            ),
            None,
            None,
        )
    if definition_line_in_target >= next_heading_line:
        return (
            ValidationIssue.create(
                "definition_after_anchor",
                index_file_path,
                line_num,
                line_num,
                message=(
                    f"Definition at line {definition_line_in_target} is after anchor "
                    f"section (next heading at line {next_heading_line})"
                ),
                severity="error",
                entry_name=entry_name,
                target_file=target_file,
                anchor=anchor,
                definition_line=definition_line_in_target,
                next_heading_line=next_heading_line,
            ),
            None,
            None,
        )
    code_block_start_line = definition_line_in_target
    if (
        code_block_start_line < anchor_section_line
        or code_block_start_line >= next_heading_line
    ):
        return (
            ValidationIssue.create(
                "code_block_outside_section",
                index_file_path,
                line_num,
                line_num,
                message=(
                    f"Definition code block at line {code_block_start_line} is not "
                    f"within anchor section (lines {anchor_section_line}-"
                    f"{next_heading_line - 1})"
                ),
                severity="error",
                entry_name=entry_name,
                target_file=target_file,
                anchor=anchor,
                code_block_line=code_block_start_line,
                section_start=anchor_section_line,
                section_end=next_heading_line - 1,
            ),
            None,
            None,
        )
    return (None, anchor_section_line, next_heading_line)


def _issue_if_block_invalid(
    lines_content: List[str],
    anchor_section_line: int,
    next_heading_line: int,
    code_block_start_line: int,
    entry_name: str,
    *,
    index_file_path: Path,
    line_num: int,
    target_file: str,
    anchor: str,
):
    """Return one ValidationIssue if block invalid/not go/unterminated/missing def, else None."""
    definition_found_in_block, in_target_block, code_block_end = _scan_block_for_definition(
        lines_content,
        anchor_section_line,
        next_heading_line,
        code_block_start_line,
        entry_name,
    )
    if not in_target_block:
        return ValidationIssue.create(
            "code_block_not_go",
            index_file_path,
            line_num,
            line_num,
            message=(
                f"Definition code block at line {code_block_start_line} is not a "
                "Go code block (expected ```go)"
            ),
            severity="error",
            entry_name=entry_name,
            target_file=target_file,
            anchor=anchor,
            code_block_line=code_block_start_line,
        )
    if code_block_end is None:
        return ValidationIssue.create(
            "code_block_unterminated",
            index_file_path,
            line_num,
            line_num,
            message=(
                f"Go code block starting at line {code_block_start_line} is not closed"
            ),
            severity="error",
            entry_name=entry_name,
            target_file=target_file,
            anchor=anchor,
            code_block_line=code_block_start_line,
        )
    if not definition_found_in_block:
        return ValidationIssue.create(
            "definition_not_in_block",
            index_file_path,
            line_num,
            line_num,
            message=(
                f"Definition '{entry_name}' not found in code block at line "
                f"{code_block_start_line} within anchor section"
            ),
            severity="error",
            entry_name=entry_name,
            target_file=target_file,
            anchor=anchor,
            code_block_line=code_block_start_line,
        )
    return None


def check_anchor_points_to_definition(
    index_content: str,
    entry_to_target_md: Dict[str, str],
    entry_to_link_anchor: Dict[str, Optional[str]],
    all_definitions: Dict[str, List[Tuple[str, int]]],
    tech_specs_dir: Path,
    *,
    index_filename: str = "api_go_defs_index.md",
    file_cache: Optional[FileContentCache] = None,
) -> List[ValidationIssue]:
    """
    Check that index entry anchors point to sections containing the definition in Go code blocks.
    """
    issues: List[ValidationIssue] = []
    content_cache = file_cache or FileContentCache()
    index_file_path = tech_specs_dir / index_filename
    lines = index_content.split("\n")
    entry_pattern = r"^\\s*-\\s+\\*\\*`([^`]+)`\\*\\*"
    repo_root = tech_specs_dir.parent.parent

    for line_num, line in enumerate(lines, 1):
        entry_match = re.match(entry_pattern, line)
        if not entry_match:
            continue
        entry_name = normalize_generic_name(entry_match.group(1))
        target_file = entry_to_target_md.get(entry_name)
        anchor = entry_to_link_anchor.get(entry_name)
        entry_issue, definition_line_in_target = _validate_entry_and_get_definition_line(
            entry_name, target_file, anchor, all_definitions,
            index_file_path, line_num=line_num,
        )
        if entry_issue is not None:
            issues.append(entry_issue)
            continue
        if definition_line_in_target is None:
            continue
        target_path = tech_specs_dir / target_file
        target_issue = _issue_if_target_missing(
            target_path, repo_root, index_file_path, line_num,
            entry_name, target_file=target_file, anchor=anchor,
        )
        if target_issue is not None:
            issues.append(target_issue)
            continue
        try:
            lines_content = content_cache.get_lines(target_path)
            headings = extract_h2_plus_headings_with_sections(
                target_path, skip_code_blocks=True, file_cache=file_cache
            )
            pos_issue, anchor_section_line, next_heading_line = (
                _issue_if_anchor_or_position_bad(
                    lines_content,
                    headings,
                    anchor,
                    definition_line_in_target,
                    index_file_path,
                    line_num=line_num,
                    entry_name=entry_name,
                    target_file=target_file,
                )
            )
            if pos_issue is not None:
                issues.append(pos_issue)
                continue
            block_issue = _issue_if_block_invalid(
                lines_content,
                anchor_section_line,
                next_heading_line,
                definition_line_in_target,
                entry_name,
                index_file_path=index_file_path,
                line_num=line_num,
                target_file=target_file,
                anchor=anchor,
            )
            if block_issue is not None:
                issues.append(block_issue)
        except (ValueError, KeyError, TypeError, AttributeError, RuntimeError) as e:
            issues.append(
                ValidationIssue.create(
                    "anchor_check_error",
                    index_file_path,
                    line_num,
                    line_num,
                    message=f"Error checking anchor: {e}",
                    severity="error",
                    entry_name=entry_name,
                    target_file=target_file,
                    anchor=anchor,
                    error=str(e),
                )
            )

    return issues
