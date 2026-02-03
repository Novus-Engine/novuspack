"""
Heading numbering check functions.

Extracted from validate_heading_numbering.py for C0302 module splitting.
Callers pass mutable issues list and optional first_error_line dict; these
functions append issues and update first_error_line where relevant.
"""

from collections import defaultdict
from pathlib import Path
from typing import Callable, List, Optional

from lib._validation_utils import (
    ValidationIssue,
    build_heading_hierarchy,
    is_organizational_heading,
)
from lib._validate_heading_numbering_helpers import is_go_code_related_heading
from lib._validate_heading_numbering_models import (
    MAX_HEADING_NUMBER_SEGMENT,
    MAX_ORGANIZATIONAL_PROSE_LINES,
)
from lib._validate_heading_numbering_title_case import to_title_case


def check_excessive_numbering(
    issues: List,
    filepath: str,
    headings: List,
) -> None:
    """
    Check for H3+ headings where the depth-specific number exceeds 20.
    Performed after corrected heading numbering is calculated.
    """
    if not headings:
        return
    h3_plus_headings = [h for h in headings if h.level >= 3]
    if not h3_plus_headings:
        return
    for heading in h3_plus_headings:
        if not heading.corrected_number:
            continue
        try:
            number_segments = [int(n) for n in heading.corrected_number.split('.')]
        except (ValueError, AttributeError):
            continue
        segment_index = heading.level - 2
        if segment_index >= len(number_segments):
            continue
        segment = number_segments[segment_index]
        if segment > MAX_HEADING_NUMBER_SEGMENT:
            msg = (
                f"H{heading.level} heading has numbering "
                f"'{heading.corrected_number}' where number {segment} "
                f"(at depth {heading.level - 1}) exceeds "
                f"{MAX_HEADING_NUMBER_SEGMENT}. "
                "Consider restructuring the document to reduce nesting depth."
            )
            warning = ValidationIssue.create(
                "heading_excessive_numbering",
                Path(filepath),
                heading.line_num,
                heading.line_num,
                message=msg,
                severity='warning',
                heading=heading.full_line,
                heading_info=heading
            )
            issues.append(warning)


def check_single_word_headings(
    issues: List,
    filepath: str,
    headings: List,
) -> None:
    """
    Check for H4+ headings where the title (after numbering) is a single word.
    """
    if not headings:
        return
    h4_plus_headings = [h for h in headings if h.level >= 4]
    if not h4_plus_headings:
        return
    for heading in h4_plus_headings:
        if not heading.heading_text:
            continue
        title = heading.heading_text.strip()
        if title and ' ' not in title:
            msg = (f"H{heading.level} heading has a single-word title '{title}'. "
                   "Consider using a more descriptive multi-word heading.")
            warning = ValidationIssue.create(
                "heading_single_word",
                Path(filepath),
                heading.line_num,
                heading.line_num,
                message=msg,
                severity='warning',
                heading=heading.full_line,
                heading_info=heading
            )
            issues.append(warning)


def check_duplicate_headings(
    issues: List,
    first_error_line: dict,
    filepath: str,
    headings: List,
) -> None:
    """
    Check for duplicate headings (excluding numbering) across all levels.
    All occurrences after the first are flagged as errors.
    """
    if not headings:
        return
    heading_groups = defaultdict(list)
    for heading in headings:
        if not heading.heading_text:
            continue
        normalized_title = heading.heading_text.strip().lower()
        if normalized_title:
            heading_groups[normalized_title].append(heading)
    for normalized_title, heading_list in heading_groups.items():
        if len(heading_list) > 1:
            heading_list.sort(key=lambda h: h.line_num)
            for duplicate_heading in heading_list[1:]:
                other_locations = [
                    f"line {h.line_num}" for h in heading_list
                    if h.line_num != duplicate_heading.line_num
                ]
                other_locations_str = ", ".join(other_locations)
                msg = (f"Duplicate heading title '{duplicate_heading.heading_text}' "
                       f"(also appears at {other_locations_str}). "
                       "Each heading should have a unique title.")
                error = ValidationIssue.create(
                    "heading_duplicate",
                    Path(filepath),
                    duplicate_heading.line_num,
                    duplicate_heading.line_num,
                    message=msg,
                    severity='error',
                    heading=duplicate_heading.full_line,
                    heading_info=duplicate_heading
                )
                issues.append(error)
                if duplicate_heading.issue is None:
                    duplicate_heading.issue = error
                if first_error_line.get(filepath) is None:
                    first_error_line[filepath] = duplicate_heading.line_num


def check_heading_capitalization(
    issues: List,
    filepath: str,
    headings: List,
    *,
    is_go_related: Optional[Callable[[str], bool]] = None,
    to_title: Optional[Callable[[str], str]] = None,
) -> None:
    """
    Check if headings follow Title Case.
    Skips headings that reference Go code elements (use actual identifiers).
    """
    if not headings:
        return
    if is_go_related is None:
        is_go_related = is_go_code_related_heading
    if to_title is None:
        to_title = to_title_case
    for heading in headings:
        if not heading.heading_text:
            continue
        if is_go_related(heading.heading_text):
            continue
        corrected = to_title(heading.heading_text)
        if heading.heading_text != corrected:
            heading.corrected_capitalization = corrected
            msg = "Capitalization may not follow title case."
            warning = ValidationIssue.create(
                "heading_capitalization",
                Path(filepath),
                heading.line_num,
                heading.line_num,
                message=msg,
                severity='warning',
                heading=heading.full_line,
                heading_info=heading
            )
            issues.append(warning)


def check_organizational_headings(
    issues: List,
    _first_error_line: dict,
    filepath: str,
    headings: List,
    content: str,
    *,
    log_fn: Optional[Callable[[str], None]] = None,
) -> None:
    """
    Check for organizational headings with no content.
    Warnings for headings that are purely organizational with no content.
    """
    if not headings:
        return
    headings_for_hierarchy = [
        (h.line_num, h.level, h.heading_text)
        for h in headings
    ]
    headings_for_hierarchy.sort(key=lambda x: x[0])
    hierarchy = build_heading_hierarchy(headings_for_hierarchy)
    for heading in headings:
        if heading.issue:
            continue
        try:
            result = is_organizational_heading(
                content,
                heading.line_num,
                heading.level,
                headings_for_hierarchy,
                hierarchy,
                max_prose_lines=MAX_ORGANIZATIONAL_PROSE_LINES
            )
            if result.get('is_organizational') and result.get('is_empty'):
                msg = ("Organizational heading with no content. "
                       "Headings should have substantive content or be removed.")
                warning = ValidationIssue.create(
                    "organizational_heading",
                    Path(filepath),
                    heading.line_num,
                    heading.line_num,
                    message=msg,
                    severity='warning',
                    heading=heading.full_line,
                    heading_info=heading
                )
                issues.append(warning)
        except (ValueError, IndexError, KeyError) as e:
            if log_fn:
                log_fn(f"  Error checking organizational heading at line {heading.line_num}: {e}")
        except (TypeError, AttributeError, RuntimeError) as e:
            if log_fn:
                log_fn(
                    f"  Unexpected error checking organizational heading at line "
                    f"{heading.line_num}: {e}"
                )


def check_h2_period_consistency(
    issues: List,
    filepath: str,
    headings: List,
) -> None:
    """
    Check if H2 headings have consistent period usage.
    If first H2 has period, all should have period; otherwise none should.
    """
    h2_headings = [h for h in headings if h.level == 2]
    if not h2_headings:
        return
    h2_headings.sort(key=lambda h: h.line_num)
    first_h2 = h2_headings[0]
    expected_has_period = first_h2.has_period
    for heading in h2_headings[1:]:
        if heading.has_period != expected_has_period:
            expected_str = "with period" if expected_has_period else "without period"
            actual_str = "with period" if heading.has_period else "without period"
            msg = (f"H2 heading period inconsistency: first H2 is {expected_str}, "
                   f"but this heading is {actual_str}. "
                   f"All H2 headings should match the first one.")
            warning = ValidationIssue.create(
                "heading_period_inconsistency",
                Path(filepath),
                heading.line_num,
                heading.line_num,
                message=msg,
                severity='warning',
                heading=heading.full_line,
                heading_info=heading
            )
            issues.append(warning)
