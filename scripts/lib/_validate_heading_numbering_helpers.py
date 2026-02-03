"""
Structure validation and check helpers for heading numbering.
"""

import re
from pathlib import Path

from lib._validation_utils import ValidationIssue
from lib._validate_heading_numbering_models import MAX_ORGANIZATIONAL_PROSE_LINES


def _is_first_h2_numbered(first_h2):
    """Return True if the first H2 heading has valid numbering."""
    return (
        first_h2.numbers and
        len(first_h2.numbers) > 0 and
        first_h2.original_number != "MISSING"
    )


def _record_unnumbered_issues(
    filepath, unnumbered_headings, issues, first_error_line
):
    """Append heading_missing_numbering issues for each unnumbered heading."""
    for line_num, level, _heading_text, full_line, heading_info in unnumbered_headings:
        msg = (
            f"H{level} heading is missing numbering. "
            "This document uses numbered headings, so all headings must be numbered."
        )
        error = ValidationIssue.create(
            "heading_missing_numbering",
            Path(filepath),
            line_num,
            line_num,
            message=msg,
            severity='error',
            heading=full_line,
            heading_info=heading_info
        )
        issues.append(error)
        heading_info.issue = error
        if first_error_line[filepath] is None:
            first_error_line[filepath] = line_num


def _check_first_h2_value(filepath, first_h2, issues):
    """If first H2 number is not 0 or 1, append issue and return True. Else return False."""
    if first_h2.numbers[0] in [0, 1]:
        return False
    error = ValidationIssue.create(
        "heading_first_h2_numbering",
        Path(filepath),
        first_h2.line_num,
        first_h2.line_num,
        message=(
            f"First H2 heading must be numbered '0' or '1', "
            f"got '{first_h2.numbers[0]}'. "
            "Please run a markdown linter to fix basic heading order, "
            "then re-run this script."
        ),
        severity='error',
        heading=first_h2.full_line,
        heading_info=first_h2
    )
    issues.append(error)
    first_h2.issue = error
    return True


def _check_heading_parents(filepath, headings, issues, log_fn=None):
    """Append heading_no_parent issues for headings with no parent; log H2 parent warning."""
    for heading in headings:
        if heading.level == 2 and heading.parent is not None:
            if log_fn:
                log_fn(f"  Warning: H2 heading at line {heading.line_num} has a parent")
        elif heading.level > 2 and heading.parent is None:
            error = ValidationIssue.create(
                "heading_no_parent",
                Path(filepath),
                heading.line_num,
                heading.line_num,
                message=(
                    f"H{heading.level} heading has no parent. "
                    "Please run a markdown linter to fix basic heading order, "
                    "then re-run this script."
                ),
                severity='error',
                heading=heading.full_line,
                heading_info=heading
            )
            issues.append(error)
            heading.issue = error


def validate_heading_structure(
    filepath, headings, unnumbered_headings, *, issues, first_error_line, log_fn=None
):
    """
    Validate heading structure after parsing. Mutates headings (sets .issue)
    and first_error_line. Appends to issues.

    Returns:
        List of headings (may be modified with issues).
    """
    if not headings:
        return []
    h2_headings = [h for h in headings if h.level == 2]
    if not h2_headings:
        return headings
    first_h2 = min(h2_headings, key=lambda h: h.line_num)
    if _is_first_h2_numbered(first_h2) and unnumbered_headings:
        _record_unnumbered_issues(filepath, unnumbered_headings, issues, first_error_line)
    if not _is_first_h2_numbered(first_h2):
        return headings
    if _check_first_h2_value(filepath, first_h2, issues):
        return headings
    _check_heading_parents(filepath, headings, issues, log_fn)
    return headings


def is_go_code_related_heading(heading_text):
    """
    Return True if the heading appears to reference a Go code element.
    """
    if not heading_text:
        return False
    camel_case_pattern = r'\b[a-z][a-zA-Z]*[A-Z][a-zA-Z]*\b'
    if re.search(camel_case_pattern, heading_text):
        return True
    method_pattern = r'\b[a-z][a-zA-Z]*\.[A-Z][a-zA-Z]*\b'
    if re.search(method_pattern, heading_text):
        return True
    go_kind_words = ['Struct', 'Function', 'Method', 'Interface', 'Type']
    for kind_word in go_kind_words:
        pattern = rf'\b[a-z][a-zA-Z]*\s+{kind_word}\b'
        if re.search(pattern, heading_text):
            return True
        method_kind_pattern = rf'\b[a-z][a-zA-Z]*\.[A-Z][a-zA-Z]*\s+{kind_word}\b'
        if re.search(method_kind_pattern, heading_text):
            return True
    return False


def get_max_organizational_prose_lines():
    """Return constant for organizational heading check (for callers that need it)."""
    return MAX_ORGANIZATIONAL_PROSE_LINES
