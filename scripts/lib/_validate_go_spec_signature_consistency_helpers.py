"""Helpers for Go signature consistency validation: extraction and counting."""

import re
from pathlib import Path
from typing import Dict, List, Optional, Tuple

from lib._go_code_utils import (
    InterfaceParser,
    Signature,
    is_public_name,
    parse_go_def_signature,
    is_example_code,
)
from lib._validation_utils import OutputBuilder, ValidationIssue, parse_no_color_flag

_RE_INTERFACE_PATTERN = re.compile(r'^\s*type\s+\w+(?:\s*\[[^\]]+\])?\s+interface\s*\{')
_RE_STRUCT_PATTERN = re.compile(r'^\s*type\s+(\w+)(?:\s*(\[[^\]]+\]))?\s+struct\s*\{')


def count_interface_methods(content: str, start_line: int, end_line: int) -> int:
    """Count methods in an interface definition within a specific line range."""
    lines = content.split('\n')
    method_count = 0
    interface_parser = InterfaceParser()

    for i in range(start_line - 1, min(end_line, len(lines))):
        line = lines[i]
        interface_name = interface_parser.check_interface_start(line)
        if interface_name:
            continue
        if interface_parser.is_in_interface():
            still_in_interface = interface_parser.update_brace_depth(line)
            if still_in_interface and interface_parser.brace_depth > 0:
                sig = parse_go_def_signature(line, location="")
                if sig and sig.kind in ('func', 'method'):
                    method_count += 1
            if not still_in_interface:
                break
    return method_count


def count_struct_fields(content: str, start_line: int, end_line: int) -> int:
    """Count fields in a struct definition."""
    lines = content.split('\n')
    field_count = 0
    in_struct = False
    brace_depth = 0

    for i in range(start_line - 1, min(end_line, len(lines))):
        line = lines[i]
        stripped = line.strip()
        if re.match(r'^\s*type\s+\w+\s+struct\s*\{', line):
            in_struct = True
            brace_depth = stripped.count('{') - stripped.count('}')
            continue
        if in_struct:
            brace_depth += stripped.count('{') - stripped.count('}')
            if brace_depth > 0 and stripped and not stripped.startswith('//'):
                if re.match(r'^\s*\w+\s+\w+', stripped):
                    if not re.match(r'^\s*func\s+', stripped):
                        field_count += 1
            if brace_depth <= 0:
                break
    return field_count


def _count_struct_fields_in_block(
    block_lines: List[str], start_index: int, initial_brace_depth: int
) -> int:
    """Count struct fields from block_lines starting at start_index."""
    field_count = 0
    temp_brace_depth = initial_brace_depth
    for j in range(start_index + 1, len(block_lines)):
        temp_line = block_lines[j]
        temp_stripped = temp_line.strip()
        if not temp_stripped or temp_stripped.startswith('//'):
            continue
        temp_brace_depth += temp_stripped.count('{') - temp_stripped.count('}')
        if temp_brace_depth > 0 and temp_stripped:
            if re.match(r'^\s*\w+\s+\w+', temp_stripped):
                if not re.match(r'^\s*func\s+', temp_stripped):
                    field_count += 1
        if temp_brace_depth <= 0:
            break
    return field_count


def extract_signatures_from_block(
    block_lines: List[str],
    start_line: int,
    relative_path: Path,
    code_content: str,
    lines: List[str],
) -> List[Signature]:
    """Extract struct/func/method/type signatures from one Go code block."""
    result: List[Signature] = []
    for i, line in enumerate(block_lines):
        line_num = start_line + i
        stripped = line.strip()
        if not stripped or stripped.startswith('//'):
            continue
        if _RE_INTERFACE_PATTERN.match(stripped):
            continue
        is_example = is_example_code(
            code_content, start_line, lines=lines, check_single_line=i
        )
        struct_match = _RE_STRUCT_PATTERN.match(stripped)
        if struct_match:
            if is_example:
                continue
            name = struct_match.group(1)
            generic_params = struct_match.group(2)
            is_public = is_public_name(name) if name else False
            brace_depth = stripped.count('{') - stripped.count('}')
            has_full_body = brace_depth > 0
            field_count = (
                _count_struct_fields_in_block(block_lines, i, brace_depth)
                if has_full_body
                else 0
            )
            result.append(Signature(
                name=name,
                kind='type',
                location=f"{relative_path}:{line_num}",
                is_public=is_public,
                has_body=has_full_body,
                field_count=field_count,
                generic_params=generic_params
            ))
            continue
        sig = parse_go_def_signature(line, location=f"{relative_path}:{line_num}")
        if not sig:
            continue
        if sig.kind in ('func', 'method'):
            result.append(Signature(
                name=sig.name,
                kind=sig.kind,
                receiver=sig.receiver,
                params=sig.params,
                returns=sig.returns,
                location=f"{relative_path}:{line_num}",
                is_public=sig.is_public,
                has_body=True
            ))
        else:
            if is_example:
                continue
            if sig.kind != 'interface':
                result.append(Signature(
                    name=sig.name,
                    kind=sig.kind,
                    location=f"{relative_path}:{line_num}",
                    is_public=sig.is_public,
                    has_body=False,
                    generic_params=sig.generic_params
                ))
    return result


def parse_cli_args(
    argv: List[str],
) -> Tuple[bool, bool, bool, Optional[str], Optional[str]]:
    """Parse CLI flags for the signature consistency validator."""
    verbose = '--verbose' in argv or '-v' in argv
    no_color = parse_no_color_flag(argv)
    no_fail = '--no-fail' in argv
    output_file = None
    target_paths_str = None
    for i, arg in enumerate(argv):
        if arg in ('--output', '-o') and i + 1 < len(argv):
            output_file = argv[i + 1]
        elif arg in ('--path', '-p') and i + 1 < len(argv):
            target_paths_str = argv[i + 1]
    return verbose, no_color, no_fail, output_file, target_paths_str


def build_output(
    verbose: bool,
    no_color: bool,
    output_file: Optional[str],
) -> OutputBuilder:
    """Create an OutputBuilder for the signature consistency validator."""
    return OutputBuilder(
        "Go Signature Consistency",
        "Validates signature consistency within tech specs",
        no_color=no_color,
        verbose=verbose,
        output_file=output_file
    )


def split_issues(
    issues: List[ValidationIssue],
) -> Tuple[List[ValidationIssue], List[ValidationIssue]]:
    """Split issues into error and warning lists."""
    errors: List[ValidationIssue] = []
    warnings: List[ValidationIssue] = []
    for issue in issues:
        if issue.matches(severity='error'):
            errors.append(issue)
        if issue.matches(severity='warning'):
            warnings.append(issue)
    return errors, warnings


def emit_issues(
    output: OutputBuilder,
    errors: List[ValidationIssue],
    warnings: List[ValidationIssue],
    no_color: bool,
) -> None:
    """Emit issue lines to the output builder."""
    for error in errors:
        output.add_error_line(error.format_message(no_color=no_color))
    for warning in warnings:
        output.add_warning_line(warning.format_message(no_color=no_color))


def emit_summary(
    output: OutputBuilder,
    all_signatures: List[Signature],
    signatures_by_key: Dict[str, List[Signature]],
    errors: List[ValidationIssue],
    warnings: List[ValidationIssue],
) -> None:
    """Emit summary section for signature validation."""
    summary_items = [
        ("Signatures checked:", len(all_signatures)),
        ("Unique definitions:", len(signatures_by_key)),
    ]
    if errors:
        summary_items.append(("Errors found:", len(errors)))
    if warnings:
        summary_items.append(("Warnings found:", len(warnings)))
    output.add_summary_header()
    output.add_summary_section(summary_items)


def emit_final_message(
    output: OutputBuilder,
    errors: List[ValidationIssue],
    warnings: List[ValidationIssue],
) -> None:
    """Emit the final success, warnings-only, or failure message."""
    if errors:
        output.add_failure_message("Validation failed. Please fix the errors above.")
        return
    if warnings:
        output.add_warnings_only_message(
            message="All signatures are consistent! (Some warnings were found - see above)",
            verbose_hint="Run with --verbose to see the full warning details.",
        )
    else:
        output.add_success_message('All signatures are consistent!')
