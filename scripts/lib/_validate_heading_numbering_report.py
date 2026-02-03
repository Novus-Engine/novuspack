"""
Report/output helpers for heading numbering validation.

Emit error blocks and numbering-display section used by validate_heading_numbering.
"""

import re
from collections import defaultdict

from lib._validation_utils import ValidationIssue, format_issue_message


def emit_org_errors_block(
    output_builder, errors_by_type, errors_by_line, *,
    rel_path_fn, build_corrected_full_line_fn, no_color
):
    """Emit organizational heading errors section."""
    org_errors = errors_by_type.get("organizational_heading", [])
    output_builder.add_line(
        f"Organizational Heading Errors ({len(org_errors)}):",
        section="error"
    )
    output_builder.add_blank_line("error")
    org_errors_by_line = {
        k: v for k, v in errors_by_line.items()
        if any(
            isinstance(e, ValidationIssue) and
            e.issue_type == "organizational_heading"
            for e in v
        )
    }
    sorted_org_lines = sorted(
        org_errors_by_line.keys(), key=lambda k: (k[0], k[1])
    )
    for file, line_num in sorted_org_lines:
        line_errors = [
            e for e in org_errors_by_line[(file, line_num)]
            if isinstance(e, ValidationIssue) and
            e.issue_type == "organizational_heading"
        ]
        if not line_errors:
            continue
        rel_file = rel_path_fn(file)
        messages = [
            e.message if isinstance(e, ValidationIssue) else e.get('message', '')
            for e in line_errors
        ]
        combined_message = "; ".join(messages)
        suggestion = None
        for error in line_errors:
            heading_info = error.extra_fields.get('heading_info')
            if heading_info:
                suggestion = build_corrected_full_line_fn(heading_info)
                break
        error_msg = format_issue_message(
            "error",
            "Organizational heading",
            rel_file,
            line_num=line_num,
            message=combined_message,
            suggestion=suggestion,
            no_color=no_color
        )
        output_builder.add_error_line(error_msg)
    output_builder.add_blank_line("error")


def emit_formatting_errors_block(
    output_builder, errors_by_type, errors_by_line, org_errors, *,
    rel_path_fn, build_corrected_full_line_fn, no_color
):
    """Emit heading formatting errors section."""
    formatting_errors = errors_by_type.get("heading_formatting", [])
    if org_errors:
        output_builder.add_separator(section="error")
        output_builder.add_blank_line("error")
    output_builder.add_line(
        f"Heading Formatting Errors ({len(formatting_errors)}):",
        section="error"
    )
    output_builder.add_blank_line("error")
    formatting_errors_by_line = {
        k: v for k, v in errors_by_line.items()
        if any(
            isinstance(e, ValidationIssue) and
            e.issue_type == "heading_formatting"
            for e in v
        )
    }
    sorted_formatting_lines = sorted(
        formatting_errors_by_line.keys(), key=lambda k: (k[0], k[1])
    )
    for file, line_num in sorted_formatting_lines:
        line_errors = [
            e for e in formatting_errors_by_line[(file, line_num)]
            if isinstance(e, ValidationIssue) and
            e.issue_type == "heading_formatting"
        ]
        if not line_errors:
            continue
        rel_file = rel_path_fn(file)
        messages = [
            e.message if isinstance(e, ValidationIssue) else e.get('message', '')
            for e in line_errors
        ]
        combined_message = "; ".join(messages)
        suggestion = None
        for error in line_errors:
            heading_info = error.extra_fields.get('heading_info')
            if heading_info:
                suggestion = build_corrected_full_line_fn(heading_info)
                break
        error_msg = format_issue_message(
            "error",
            "Heading formatting",
            rel_file,
            line_num=line_num,
            message=combined_message,
            suggestion=suggestion,
            no_color=no_color
        )
        output_builder.add_error_line(error_msg)
    output_builder.add_blank_line("error")


def emit_numbering_errors_block(
    output_builder, errors_by_line, *,
    org_errors, formatting_errors, numbering_errors,
    rel_path_fn, build_corrected_full_line_fn, no_color
):
    """Emit heading numbering errors section."""
    if org_errors or formatting_errors:
        output_builder.add_separator(section="error")
        output_builder.add_blank_line("error")
    output_builder.add_line(
        f"Heading Numbering Errors ({len(numbering_errors)}):",
        section="error"
    )
    output_builder.add_blank_line("error")
    numbering_errors_by_line = {
        k: v for k, v in errors_by_line.items()
        if any(
            isinstance(e, ValidationIssue) and
            e.issue_type not in (
                "organizational_heading", "heading_formatting"
            )
            for e in v
        )
    }
    sorted_numbering_lines = sorted(
        numbering_errors_by_line.keys(), key=lambda k: (k[0], k[1])
    )
    for file, line_num in sorted_numbering_lines:
        line_errors = numbering_errors_by_line[(file, line_num)]
        if not line_errors:
            continue
        rel_file = rel_path_fn(file)
        messages = [
            e.message if isinstance(e, ValidationIssue) else e.get('message', '')
            for e in line_errors
        ]
        combined_message = "; ".join(messages)
        suggestion = None
        for error in line_errors:
            heading_info = error.extra_fields.get('heading_info')
            if heading_info:
                suggestion = build_corrected_full_line_fn(heading_info)
                break
        error_msg = format_issue_message(
            "error",
            "Heading numbering",
            rel_file,
            line_num=line_num,
            message=combined_message,
            suggestion=suggestion,
            no_color=no_color
        )
        output_builder.add_error_line(error_msg)
    output_builder.add_blank_line("error")


def filter_headings_with_numbering_errors(errored_headings):
    """Return list of headings that have numbering errors."""
    result = []
    for heading in errored_headings:
        if heading.original_number == "MISSING" and heading.corrected_number:
            result.append(heading)
        elif heading.original_number and heading.corrected_number:
            current = heading.original_number.rstrip('.')
            correct = heading.corrected_number.rstrip('.')
            if current != correct:
                result.append(heading)
    return result


def emit_numbering_display_block(
    output_builder, _filepath, first_error_line, rel_file,
    headings_with_numbering_errors, *,
    _build_corrected_full_line_fn=None
):
    """
    Emit the 'Sorted headings from first error' block (format for apply_heading_corrections).
    """
    output_builder.add_separator(section="error")
    output_builder.add_line(
        f"Sorted headings from first error (line {first_error_line}) "
        f"in {rel_file}:",
        section="error"
    )
    output_builder.add_separator(section="error")
    output_builder.add_blank_line("error")
    output_builder.add_line(
        "The following headings should be in this order "
        "(sorted by numeric values):",
        section="error"
    )
    output_builder.add_blank_line("error")
    output_builder.add_line(
        "Format: Line X: [CURRENT] -> [CORRECT] Title",
        section="error"
    )
    output_builder.add_blank_line("error")

    sorted_headings = sorted(
        headings_with_numbering_errors,
        key=lambda h: (h.line_num, h.sort_key())
    )
    max_line_num = (
        max(h.line_num for h in sorted_headings)
        if sorted_headings else 0
    )
    line_num_width = len(str(max_line_num))

    h2_headings_in_output = [h for h in sorted_headings if h.level == 2]
    display_period = False
    if h2_headings_in_output:
        first_h2 = min(h2_headings_in_output, key=lambda h: h.line_num)
        display_period = first_h2.has_period

    for heading in sorted_headings:
        current_number_str = heading.original_number
        if heading.corrected_number is None:
            correct_number_str = current_number_str
        else:
            correct_number_str = heading.corrected_number

        if current_number_str == "MISSING":
            current_display = "MISSING"
        elif heading.level == 2 and display_period:
            current_display = f"{current_number_str}."
        else:
            current_display = current_number_str

        current_for_comparison = current_number_str.rstrip('.')
        correct_for_comparison = correct_number_str.rstrip('.')
        needs_change = current_for_comparison != correct_for_comparison

        is_duplicate_error = (
            heading.issue and
            isinstance(heading.issue, ValidationIssue) and
            heading.issue.matches(issue_type="heading_duplicate")
        )

        if heading.level == 2 and display_period:
            correct_display = f"{correct_number_str}."
        else:
            correct_display = correct_number_str

        heading_text_display = heading.heading_text
        if heading.corrected_capitalization:
            heading_text_display = heading.corrected_capitalization

        if is_duplicate_error and not needs_change:
            output_builder.add_error_line(
                f"Line {heading.line_num:{line_num_width}d}: "
                f"{'#' * heading.level} [{current_display}] (DUPLICATE) "
                f"{heading_text_display}"
            )
        else:
            output_builder.add_error_line(
                f"Line {heading.line_num:{line_num_width}d}: "
                f"{'#' * heading.level} [{current_display}] -> "
                f"[{correct_display}] {heading_text_display}"
            )

    output_builder.add_blank_line("error")


def _partition_errors_warnings(validator_issues):
    """Return (errors, warnings) from validator_issues."""
    errors = [i for i in validator_issues if i.matches(severity='error')]
    warnings = [i for i in validator_issues if i.matches(severity='warning')]
    return (errors, warnings)


def _emit_errors_section(
    output_builder,
    errors,
    headings_from_first_error,
    first_error_line,
    *,
    rel_path_fn,
    build_corrected_full_line_fn,
    no_color,
):
    """Emit errors header and all error blocks (org, formatting, numbering, display)."""
    output_builder.add_errors_header()
    output_builder.add_line(f"Found {len(errors)} error(s):", section="error")
    output_builder.add_blank_line("error")
    errors_by_type = defaultdict(list)
    for error in errors:
        errors_by_type[error.issue_type].append(error)
    errors_by_line = defaultdict(list)
    for error in errors:
        errors_by_line[(error.file, error.start_line)].append(error)
    org_errors = errors_by_type.get("organizational_heading", [])
    if org_errors:
        emit_org_errors_block(
            output_builder, errors_by_type, errors_by_line,
            rel_path_fn=rel_path_fn,
            build_corrected_full_line_fn=build_corrected_full_line_fn,
            no_color=no_color,
        )
    formatting_errors = errors_by_type.get("heading_formatting", [])
    if formatting_errors:
        emit_formatting_errors_block(
            output_builder, errors_by_type, errors_by_line, org_errors,
            rel_path_fn=rel_path_fn,
            build_corrected_full_line_fn=build_corrected_full_line_fn,
            no_color=no_color,
        )
    numbering_errors = [
        e for e in errors
        if isinstance(e, ValidationIssue) and
        e.issue_type not in ("organizational_heading", "heading_formatting")
    ]
    if numbering_errors:
        emit_numbering_errors_block(
            output_builder, errors_by_line,
            org_errors=org_errors,
            formatting_errors=formatting_errors,
            numbering_errors=numbering_errors,
            rel_path_fn=rel_path_fn,
            build_corrected_full_line_fn=build_corrected_full_line_fn,
            no_color=no_color,
        )
    if numbering_errors:
        for filepath in sorted(headings_from_first_error.keys()):
            errored_headings = headings_from_first_error[filepath]
            if not errored_headings:
                continue
            headings_with_numbering_errors = filter_headings_with_numbering_errors(
                errored_headings
            )
            if not headings_with_numbering_errors:
                continue
            first_error_line_val = first_error_line[filepath]
            rel_file = rel_path_fn(filepath)
            emit_numbering_display_block(
                output_builder, filepath, first_error_line_val, rel_file,
                headings_with_numbering_errors,
                _build_corrected_full_line_fn=build_corrected_full_line_fn,
            )


def _emit_warnings_section(
    output_builder, warnings, *, rel_path_fn, build_corrected_full_line_fn, no_color
):
    """Emit warnings header and warning lines."""
    output_builder.add_warnings_header()
    output_builder.add_line(f"Found {len(warnings)} warning(s):", section="warning")
    output_builder.add_blank_line("warning")
    warnings_by_line = defaultdict(list)
    for warning in warnings:
        if isinstance(warning, ValidationIssue):
            key = (warning.file, warning.start_line)
        else:
            key = (warning.get('file', ''), warning.get('line_num', 0))
        warnings_by_line[key].append(warning)

    def _category_for_issue_type(issue_type):
        if issue_type == "heading_capitalization":
            return "Heading capitalization"
        if issue_type == "organizational_heading":
            return "Organizational heading"
        return "Heading numbering"

    for file, line_num in sorted(warnings_by_line.keys(), key=lambda k: (k[0], k[1])):
        line_warnings = warnings_by_line[(file, line_num)]
        rel_file = rel_path_fn(file)
        messages = []
        first_issue_type = None
        for warning in line_warnings:
            if isinstance(warning, ValidationIssue):
                if first_issue_type is None:
                    first_issue_type = warning.issue_type
                msg = warning.message
            else:
                msg = warning.get('message', '')
            if "expected" in msg.lower():
                msg = re.sub(r",\s*expected\s+['\"][^'\"]+['\"]", "", msg)
            messages.append(msg)
        combined_message = "; ".join(messages)
        category = _category_for_issue_type(first_issue_type or "")
        suggestion = None
        for warning in line_warnings:
            heading_info = (
                warning.extra_fields.get('heading_info')
                if isinstance(warning, ValidationIssue)
                else getattr(warning, 'heading_info', None)
            )
            if heading_info:
                is_cap = (
                    isinstance(warning, ValidationIssue)
                    and warning.issue_type == "heading_capitalization"
                )
                if is_cap:
                    suggestion = (
                        f"{heading_info.full_line} => "
                        f"{build_corrected_full_line_fn(heading_info)}"
                    )
                else:
                    suggestion = build_corrected_full_line_fn(heading_info)
                break
        warning_msg = format_issue_message(
            "warning", category, rel_file,
            line_num=line_num, message=combined_message,
            suggestion=suggestion, no_color=no_color,
        )
        output_builder.add_warning_line(warning_msg)


def _emit_final_message(output_builder, errors, warnings, all_headings):
    """Emit success/final_message/failure based on errors and warnings."""
    if not errors and not warnings:
        files_checked = len(all_headings)
        total_headings = sum(len(h) for h in all_headings.values())
        output_builder.add_summary_header()
        output_builder.add_summary_section([
            ("Files checked:", files_checked),
            ("Headings checked:", total_headings),
        ])
        output_builder.add_success_message("All heading numbering is valid!")
    elif not errors:
        output_builder.add_warnings_only_message(
            message="No heading numbering errors found (only warnings).",
        )
    else:
        output_builder.add_failure_message(
            "Validation failed. Please fix the errors above."
        )


def print_summary(
    validator_issues,
    all_headings,
    headings_from_first_error,
    first_error_line,
    *,
    rel_path_fn,
    build_corrected_full_line_fn,
    no_color,
    output_builder
):
    """
    Print validation summary. Uses validator's rel_path and build_corrected_full_line
    via the passed functions; issues/headings/state passed as data.
    """
    errors, warnings = _partition_errors_warnings(validator_issues)
    if errors:
        _emit_errors_section(
            output_builder, errors, headings_from_first_error, first_error_line,
            rel_path_fn=rel_path_fn,
            build_corrected_full_line_fn=build_corrected_full_line_fn,
            no_color=no_color,
        )
    if warnings:
        _emit_warnings_section(
            output_builder, warnings,
            rel_path_fn=rel_path_fn,
            build_corrected_full_line_fn=build_corrected_full_line_fn,
            no_color=no_color,
        )
    _emit_final_message(output_builder, errors, warnings, all_headings)
