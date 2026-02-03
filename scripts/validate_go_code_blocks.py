#!/usr/bin/env python3
"""
Validate Go code blocks in tech specs documentation.

This script checks that:
1. Each Go code block has at most one type or interface definition
2. Each Go code block has at most one function or method definition
3. Type definitions and function definitions are mutually exclusive in a code block
4. Each Go code block is under a different heading
5. Function/Method/Type headings should include the definition name and kind word;
   definition names are preferred in backticks (e.g. `` `Package.Write` Method ``).
   Case inside backticks is ignored for validation.
6. All type, interface, struct, function, and method definitions have
   comments preceding them

Usage:
    python3 scripts/validate_go_code_blocks.py [options]

Options:
    --verbose, -v           Show detailed progress information
    --output, -o FILE       Write detailed report to FILE
    --path, -p PATHS        Check only the specified file(s) or
                            directory(ies) (recursive). Can be a single
                            path or comma-separated list of paths
    --nocolor, --no-color   Disable colored output
    --help, -h              Show this help message

Examples:
    # Basic validation
    python3 scripts/validate_go_code_blocks.py

    # Save report to file
    python3 scripts/validate_go_code_blocks.py --output dev_docs/go_code_blocks_audit.md

    # Verbose output
    python3 scripts/validate_go_code_blocks.py --verbose

    # Check specific file
    python3 scripts/validate_go_code_blocks.py --path docs/tech_specs/api_file_management.md

    # Check specific directory
    python3 scripts/validate_go_code_blocks.py --path docs/tech_specs

    # Check multiple paths
    python3 scripts/validate_go_code_blocks.py --path \\
        docs/tech_specs/api_file_management.md,docs/tech_specs/api_core.md
"""

import re
import sys
from pathlib import Path
from collections import defaultdict
from typing import List, Tuple, Dict, Optional

from lib._validation_utils import (
    OutputBuilder, parse_no_color_flag,
    find_markdown_files, parse_paths, get_validation_exit_code,
    find_heading_for_code_block,
    remove_backticks_keep_content,
    ValidationIssue,
    DOCS_DIR, TECH_SPECS_DIR
)
from lib._validate_go_code_blocks_report import (
    generate_report,
    print_summary,
    _has_non_warning_errors,
)
from lib._go_code_utils import (
    find_go_code_blocks as find_go_code_blocks_base,
    is_example_code,
    is_example_signature_name,
    check_kind_word_after,
    find_definition_line_index,
    is_example_definition,
    count_go_definitions,
    find_first_definition
)


def find_go_code_blocks(content: str, _file_path: str) -> List[Tuple[int, int, str, Optional[str]]]:
    """
    Find all Go code blocks in markdown content with heading context.

    Returns list of tuples: (start_line, end_line, code_content, heading)
    """
    go_blocks = []

    # Use base function to find blocks, then add heading context
    base_blocks = find_go_code_blocks_base(content)

    for start_line, end_line, code_content in base_blocks:
        # Use shared utility to find heading before this code block
        current_heading = find_heading_for_code_block(content, start_line)
        go_blocks.append((start_line, end_line, code_content, current_heading))

    return go_blocks


def has_preceding_comment(code_lines: List[str], def_line_index: int) -> bool:
    """
    Check if a definition at the given line index has a preceding comment.

    Looks for comments immediately before the definition line.
    Comments can be:
    - Single-line comments (//) on the line immediately before (skipping blank lines)
    - Multi-line comments (/* */) ending on or before the definition line

    Args:
        code_lines: List of code lines
        def_line_index: Index of the definition line (0-indexed)

    Returns:
        True if there's a comment preceding the definition, False otherwise
    """
    if not def_line_index:
        # First line of code block, no preceding comment possible
        return False

    # Find the last non-blank line before the definition
    last_non_blank_idx = None
    for i in range(def_line_index - 1, -1, -1):
        if code_lines[i].strip():
            last_non_blank_idx = i
            break

    if last_non_blank_idx is None:
        # Only blank lines before definition
        return False

    last_line = code_lines[last_non_blank_idx]
    last_stripped = last_line.strip()

    # Check for single-line comment
    if last_stripped.startswith('//'):
        return True

    # Check for multi-line comment on this line (/* ... */)
    if '/*' in last_line and '*/' in last_line:
        return True

    # Check if we're inside a multi-line comment by looking backwards
    # Track if we're in a comment block
    in_comment = False
    for i in range(last_non_blank_idx, -1, -1):
        line = code_lines[i]
        stripped = line.strip()

        # If we find a comment end, we're in a comment
        if '*/' in line:
            in_comment = True
            # Check if comment also starts on this line
            if '/*' in line:
                return True
            continue

        # If we find a comment start
        if '/*' in line:
            if in_comment:
                # We found the start of the comment we're in
                return True
            # Comment starts here, check if it ends on same line
            if '*/' in line:
                return True
            # Comment starts but doesn't end on same line - not a preceding comment
            break

        # If we're in a comment, continue looking backwards
        if in_comment:
            continue

        # If we hit a non-empty line that's not a comment and not in a comment,
        # there's no preceding comment
        if stripped and not stripped.startswith('//'):
            break

    return False


def validate_single_definition(
    code: str,
    code_lines: List[str],
    start_line: int,
    end_line: int,
    lines: List[str],
    *,
    heading: Optional[str],
    is_type: bool,
    file_path: Path,
    issues: List[ValidationIssue]
) -> None:
    """
    Validate a code block that contains exactly one definition.

    Performs all validations that apply to single-definition blocks:
    - Missing comment check
    - Heading format validation
    - Example code detection (skips validation if example)
    """
    # Skip if example code
    if is_example_definition(code, start_line, lines, heading, is_type):
        return

    # Check for missing comment
    check_missing_comment(
        code, code_lines, start_line, lines,
        heading=heading, is_type=is_type, file_path=file_path, issues=issues
    )

    # Validate heading format
    if heading:
        sig = find_first_definition(code, is_type=is_type)
        if sig:
            if is_type:
                def_name = sig.normalized_type_name()
                kind = sig.kind
                def_kind = kind
                receiver_type = None
            else:
                func_name = sig.name
                receiver_type = sig.receiver
                is_method = sig.is_method()
                def_name = func_name
                def_kind = 'method' if is_method else 'function'

            is_valid, error_messages, suggestion = validate_heading_format(
                heading, name=def_name, kind=def_kind, receiver_type=receiver_type
            )
            search_term = (
                f'{receiver_type}.{def_name}' if def_kind == 'method' else def_name
            )
            kind_word = def_kind.capitalize()
            if not is_valid:
                for error_msg in error_messages:
                    extra_fields = {}
                    if is_type:
                        extra_fields['type_name'] = def_name
                        extra_fields['kind'] = kind
                    else:
                        extra_fields['func_name'] = def_name
                        extra_fields['receiver_type'] = receiver_type
                        extra_fields['is_method'] = def_kind == 'method'
                    issues.append(ValidationIssue.create(
                        'heading_format',
                        file_path,
                        start_line,
                        end_line,
                        message=error_msg,
                        severity='error',
                        suggestion=suggestion,
                        heading=heading,
                        **extra_fields
                    ))
            elif not _heading_has_name_in_backticks(heading, search_term):
                suggested_heading = suggest_heading(heading, search_term, kind_word)
                extra_fields = {}
                if is_type:
                    extra_fields['type_name'] = def_name
                    extra_fields['kind'] = kind
                else:
                    extra_fields['func_name'] = def_name
                    extra_fields['receiver_type'] = receiver_type
                    extra_fields['is_method'] = def_kind == 'method'
                issues.append(ValidationIssue.create(
                    'heading_prefer_backticks',
                    file_path,
                    start_line,
                    end_line,
                    message='Prefer backticks for definition name in heading',
                    severity='warning',
                    suggestion=f'Suggested: {suggested_heading}',
                    heading=heading,
                    **extra_fields
                ))


def check_missing_comment(
    code: str,
    code_lines: List[str],
    start_line: int,
    lines: List[str],
    *,
    heading: Optional[str],
    is_type: bool,
    file_path: Path,
    issues: List[Dict]
) -> None:
    """Unified function to check for missing comments on definitions."""
    def_line_idx = find_definition_line_index(code, is_type=is_type)
    if def_line_idx is None:
        return

    sig = find_first_definition(code, is_type=is_type)
    if not sig:
        return

    if is_type:
        def_name = sig.normalized_type_name()
        def_kind = sig.kind
        def_display = def_name
    else:
        func_name = sig.name
        receiver_type = sig.receiver
        is_method = sig.is_method()
        def_kind = 'method' if is_method else 'function'
        def_display = f'{receiver_type}.{func_name}' if is_method else func_name

    # Skip example code
    if is_example_signature_name(sig.name):
        return

    # Use is_example_definition helper to check if this definition is example code
    if is_example_definition(code, start_line, lines, heading, is_type):
        return

    # Check for comment
    if not has_preceding_comment(code_lines, def_line_idx):
        def_line_num = start_line + def_line_idx
        issues.append(ValidationIssue.create(
            'missing_comment',
            file_path,
            def_line_num,
            def_line_num,
            message=(
                f'{def_kind.capitalize()} definition `{def_display}` '
                'does not have a preceding comment'
            ),
            severity='error',
            suggestion=(
                f'Add a comment before the {def_kind} definition '
                f'`{def_display}` to document its purpose'
            ),
            heading=heading,
            def_name=def_display,
            def_kind=def_kind
        ))


def _heading_contains_name(heading: str, search_term: str) -> bool:
    """
    Return True if heading contains the definition name (case-insensitive).
    Strips backticks for comparison so `` `Package.Write` `` matches Package.Write.
    """
    normalized = remove_backticks_keep_content(heading)
    return search_term.lower() in normalized.lower()


def _heading_has_name_in_backticks(heading: str, search_term: str) -> bool:
    """
    Return True if heading contains the definition name inside backticks.
    Used to recommend backticks when the heading is valid but plain.
    """
    pattern = re.compile(r'`' + re.escape(search_term) + r'`', re.IGNORECASE)
    return pattern.search(heading) is not None


def suggest_heading(heading: str, search_term: str, kind_word: str) -> str:
    """
    Suggest a corrected heading for a definition.
    Prefers definition name in backticks: `` `Name` Kind remaining ``.

    Preserves the numbering prefix if present.

    Args:
        heading: Current heading text (may include numbering like "2.5 Heading Text")
        search_term: The term to use (e.g., "Package", "FileEntry.GetState", "NewPackage")
        kind_word: The kind word to include (e.g., "Method", "Function", "Struct", "Interface")

    Returns:
        Suggested heading in format: {number} `{search_term}` {kind_word} {remaining}
        (number is preserved if present in original heading)
    """
    # Pattern to match numbered headings: "2.5" or "2.5.3" etc.
    numbered_pattern = re.compile(r'^([0-9]+(?:\.[0-9]+)*)\.?\s+(.+)$')

    # Extract numbering prefix if present
    numbering_prefix = None
    heading_without_number = heading
    match = numbered_pattern.match(heading.strip())
    if match:
        numbering_prefix = match.group(1)
        heading_without_number = match.group(2)

    # Normalize (strip backticks) for extracting remaining text
    normalized = remove_backticks_keep_content(heading_without_number)

    # Remove search_term (case-insensitive)
    if '.' in search_term:
        pattern = re.compile(re.escape(search_term), re.IGNORECASE)
    else:
        pattern = re.compile(r'\b' + re.escape(search_term) + r'\b', re.IGNORECASE)
    normalized = pattern.sub('', normalized, count=1)

    # Remove kind_word (case-insensitive, whole word)
    kind_pattern = re.compile(r'\b' + re.escape(kind_word) + r'\b', re.IGNORECASE)
    normalized = kind_pattern.sub('', normalized, count=1)
    # Also remove common long form (e.g. "Structure" when kind is "Struct")
    if kind_word == 'Struct':
        normalized = re.sub(r'\bStructure\b', '', normalized, count=1, flags=re.IGNORECASE)

    remaining = ' '.join(normalized.split())

    # Prefer definition name in backticks
    suggested = f"`{search_term}` {kind_word} {remaining}".strip()

    if numbering_prefix:
        suggested = f"{numbering_prefix} {suggested}"

    return suggested


def validate_heading_format(
    heading: str,
    name: str,
    kind: str,
    receiver_type: Optional[str] = None
) -> Tuple[bool, List[str], Optional[str]]:
    """
    Validate that heading format matches the definition in the code block.

    This function handles all definition types: functions, methods, types, interfaces,
    structs, and other type definitions.

    Rules:
    - Definition name may appear in backticks (preferred) or plain; case inside
      backticks is ignored for validation.
    - Functions: Heading should include function name and "Function"
      (e.g. `` `NewPackage` Function ``).
    - Methods: Heading should include Type.MethodName and "Method"
      (e.g. `` `Package.Write` Method ``).
    - Types/Interfaces/Structs: Heading should include type name and kind word
      (e.g. `` `Package` Struct ``).

    Args:
        heading: The heading text (required)
        name: The name of the definition (function name, method name, or type name)
        kind: The kind of definition ('function', 'method', 'interface', 'struct', 'type',
              'alias', 'pointer', 'slice', or 'map')
        receiver_type: Receiver type name (required for methods, None for functions and types)

    Returns:
        Tuple of (is_valid, error_messages, suggestion)
        - is_valid: True if heading format is correct
        - error_messages: List of error messages if invalid, empty list if valid
        - suggestion: Suggested fix if invalid, None if valid
    """
    errors = []

    # Determine expected kind word by capitalizing the kind
    kind_word = kind.capitalize()

    # Determine search term and display name based on kind
    if kind == 'method' and receiver_type:
        search_term = f'{receiver_type}.{name}'
        display_name = f'`{search_term}`'
    else:
        search_term = name
        display_name = f'`{name}`'

    # Check if definition name is present (case-insensitive; backticks stripped for comparison)
    if kind != 'method' and not _heading_contains_name(heading, search_term):
        errors.append(
            f'Heading should include {kind.capitalize()} name: {name} (prefer in backticks)'
        )

    match kind:
        case 'method':
            if receiver_type:
                if not _heading_contains_name(heading, search_term):
                    errors.append(
                        f'Method heading should include {display_name} and {kind_word}'
                    )
                else:
                    check_kind_word_after(
                        heading, search_term, kind_word,
                        display_name=display_name,
                        error_prefix='Method heading',
                        errors=errors
                    )

        case _:  # Default: types and functions
            if _heading_contains_name(heading, search_term):
                check_kind_word_after(
                    heading, search_term, kind_word,
                    display_name=display_name,
                    error_prefix=f'{kind.capitalize()} heading',
                    errors=errors
                )

    # Generate suggestion if there are errors (common logic for all cases)
    suggestion = None
    if errors:
        # Determine search term based on kind
        if kind == 'method' and receiver_type:
            search_term = f'{receiver_type}.{name}'
        else:
            search_term = name

        suggested_heading = suggest_heading(heading, search_term, kind_word)
        suggestion = f'Suggested: {suggested_heading}'

    return (not errors, errors, suggestion)


def _append_block_count_issues(
    issues: List[ValidationIssue],
    file_path: Path,
    start_line: int,
    end_line: int,
    heading: Optional[str],
    *,
    type_count: int,
    func_count: int,
    func_type_count: int,
) -> None:
    """Append validation issues for block definition counts (multiple types/funcs, etc.)."""
    if type_count > 1:
        issues.append(ValidationIssue.create(
            'multiple_types', file_path, start_line, end_line,
            message=f'Code block has {type_count} type/interface definitions (max 1 allowed)',
            severity='error',
            suggestion='Split into separate code blocks, one per type/interface definition',
            heading=heading, type_count=type_count,
        ))
    if func_count > 1:
        issues.append(ValidationIssue.create(
            'multiple_funcs', file_path, start_line, end_line,
            message=f'Code block has {func_count} func definitions (max 1 allowed)',
            severity='error',
            suggestion='Split into separate code blocks, one per function/method definition',
            heading=heading, func_count=func_count,
        ))
    if type_count > 0 and func_count > 0:
        issues.append(ValidationIssue.create(
            'type_func_exclusive', file_path, start_line, end_line,
            message=(
                f'Code block has both {type_count} type definition(s) and '
                f'{func_count} func definition(s) (must be exclusive)'
            ),
            severity='error',
            suggestion=(
                'Separate type definitions and function definitions into different code blocks'
            ),
            heading=heading, type_count=type_count, func_count=func_count,
        ))
    if func_type_count > 0:
        issues.append(ValidationIssue.create(
            'function_type_warning', file_path, start_line, end_line,
            message=(
                f'Code block has {func_type_count} function type definition(s) '
                '(review recommended)'
            ),
            severity='warning',
            suggestion=(
                'Review if function type definitions are intentional and properly documented'
            ),
            heading=heading, func_type_count=func_type_count,
        ))


def _append_heading_usage_issues(
    issues: List[ValidationIssue],
    file_path: Path,
    heading_usage: Dict[str, List[tuple]],
) -> None:
    """Append issues for multiple blocks per heading."""
    for heading, blocks_under_heading in heading_usage.items():
        if len(blocks_under_heading) <= 1:
            continue
        first_block_start = min(b[0] for b in blocks_under_heading)
        last_block_end = max(b[1] for b in blocks_under_heading)
        issues.append(ValidationIssue.create(
            'multiple_blocks_per_heading',
            file_path, first_block_start, last_block_end,
            message=(
                f'Heading "{heading}" has {len(blocks_under_heading)} '
                'Go code blocks (each should be under a different heading)'
            ),
            severity='error',
            suggestion='Move each code block to a separate heading',
            heading=heading, blocks=blocks_under_heading,
        ))


def audit_file(file_path: Path) -> Dict:
    """Audit a single markdown file for Go code block issues."""
    issues: List[ValidationIssue] = []
    code_blocks = []

    try:
        content = file_path.read_text(encoding='utf-8')
        lines = content.split('\n')
        blocks = find_go_code_blocks(content, str(file_path))
        heading_usage = defaultdict(list)

        for start_line, end_line, code, heading in blocks:
            counts = count_go_definitions(code)
            type_count = counts['type']
            func_count = counts['func'] + counts['method']
            func_type_count = counts['func_type']
            code_lines = code.split('\n')

            code_blocks.append({
                'start_line': start_line,
                'end_line': end_line,
                'heading': heading,
                'type_count': type_count,
                'func_count': func_count,
                'func_type_count': func_type_count,
                'code_preview': code[:100] + '...' if len(code) > 100 else code
            })

            if is_example_code(
                code, start_line,
                lines=lines, heading_text=heading, check_prose_before_block=True
            ):
                if heading:
                    heading_usage[heading].append((start_line, end_line))
                continue

            _append_block_count_issues(
                issues, file_path, start_line, end_line, heading,
                type_count=type_count, func_count=func_count, func_type_count=func_type_count,
            )
            if type_count > 0:
                check_missing_comment(
                    code, code_lines, start_line, lines,
                    heading=heading, is_type=True, file_path=file_path, issues=issues
                )
            if func_count > 0:
                check_missing_comment(
                    code, code_lines, start_line, lines,
                    heading=heading, is_type=False, file_path=file_path, issues=issues
                )
            if (type_count == 1 and not func_count) or (func_count == 1 and not type_count):
                validate_single_definition(
                    code, code_lines, start_line, end_line, lines,
                    heading=heading, is_type=(type_count == 1),
                    file_path=file_path, issues=issues
                )
            if heading:
                heading_usage[heading].append((start_line, end_line))

        _append_heading_usage_issues(issues, file_path, dict(heading_usage))
        for block in code_blocks:
            if block['heading'] is None:
                issues.append(ValidationIssue.create(
                    'no_heading', file_path,
                    block['start_line'], block['end_line'],
                    message='Code block is not under any heading',
                    severity='error',
                    suggestion='Add a heading above the code block'
                ))

    except (OSError, ValueError, KeyError, TypeError, AttributeError, RuntimeError) as e:
        issues.append(ValidationIssue.create(
            'error',
            file_path,
            1,
            1,
            message=f'Error reading file: {e}',
            severity='error'
        ))

    return {
        'file': str(file_path),
        'code_blocks': code_blocks,
        'issues': issues
    }


def main():
    """Main entry point."""

    # Show help if requested
    if '--help' in sys.argv or '-h' in sys.argv:
        print(__doc__)
        return 0

    # Parse command line arguments
    verbose = '--verbose' in sys.argv or '-v' in sys.argv
    no_color = parse_no_color_flag(sys.argv)
    no_fail = '--no-fail' in sys.argv
    output_file = None
    target_paths_str = None

    for i, arg in enumerate(sys.argv):
        if arg in ('--output', '-o') and i + 1 < len(sys.argv):
            output_file = sys.argv[i + 1]
        elif arg in ('--path', '-p') and i + 1 < len(sys.argv):
            target_paths_str = sys.argv[i + 1]

    # Parse comma-separated paths
    target_paths = parse_paths(target_paths_str)

    # Find markdown files to audit
    default_dir = Path(DOCS_DIR) / TECH_SPECS_DIR
    md_files = find_markdown_files(
        target_paths=target_paths,
        default_dir=default_dir,
        verbose=verbose
    )

    if not md_files:
        print('Error: No markdown files found', file=sys.stderr)
        return 1

    # Create output builder for all output (header streams immediately if verbose)
    output = OutputBuilder(
        "Go Code Blocks Validation",
        "Validates Go code blocks in tech specs",
        no_color=no_color, verbose=verbose, output_file=output_file
    )

    # 2. Add working verbose output (streams immediately if verbose)
    if verbose:
        output.add_verbose_line(f'Found {len(md_files)} markdown file(s) to audit')
        output.add_blank_line("working_verbose")

    # Audit each file
    results = []
    for md_file in md_files:
        if verbose:
            output.add_verbose_line(f'Auditing {md_file.name}...')
        result = audit_file(md_file)
        results.append(result)

    # Generate report if output file specified
    if output_file:
        output_path = Path(output_file)
        generate_report(results, output_path)

    # 3. Add summary, warnings, errors, and final messages
    print_summary(results, output=output, verbose=verbose, no_color=no_color)

    # Print all output at once
    output.print()

    # Exit with error code only if non-warning issues found (unless --no-fail is set)
    has_errors = _has_non_warning_errors(results)
    return get_validation_exit_code(has_errors, no_fail)


if __name__ == '__main__':
    sys.exit(main())
