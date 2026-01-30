#!/usr/bin/env python3
"""
Validate Go code blocks in tech specs documentation.

This script checks that:
1. Each Go code block has at most one type or interface definition
2. Each Go code block has at most one function or method definition
3. Type definitions and function definitions are mutually exclusive in a code block
4. Each Go code block is under a different heading
5. Function headings should NOT include the function name in backticks
   (e.g., NewPackage, not `NewPackage`)
6. Method headings should NOT include type and method name in backticks
   (e.g., FileEntry.GetProcessingState, not `FileEntry.GetProcessingState`)
7. Type/Interface/Struct headings should NOT include the type name in backticks
   (e.g., Package, not `Package`)
8. All type, interface, struct, function, and method definitions have
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

import sys
from pathlib import Path
from collections import defaultdict
from typing import List, Tuple, Dict, Optional

scripts_dir = Path(__file__).parent
lib_dir = scripts_dir / "lib"

# Import shared utilities
for module_path in (str(scripts_dir), str(lib_dir)):
    if module_path not in sys.path:
        sys.path.insert(0, module_path)

from lib._validation_utils import (  # noqa: E402
    OutputBuilder, parse_no_color_flag,
    find_markdown_files, parse_paths, get_validation_exit_code,
    find_heading_for_code_block, format_issue_message,
    has_backticks, get_backticks_error_message,
    ValidationIssue,
    DOCS_DIR, TECH_SPECS_DIR
)
from lib._go_code_utils import (  # noqa: E402
    find_go_code_blocks as find_go_code_blocks_base,
    is_example_code,
    is_example_signature_name,
    check_kind_word_after,
    find_definition_line_index,
    is_example_definition,
    count_go_definitions,
    find_first_definition
)


def find_go_code_blocks(content: str, file_path: str) -> List[Tuple[int, int, str, Optional[str]]]:
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
    if def_line_index == 0:
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
        code, code_lines, start_line, lines, heading,
        is_type, file_path, issues
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
            if not is_valid:
                for error_msg in error_messages:
                    extra_fields = {}
                    if is_type:
                        extra_fields['type_name'] = def_name
                        extra_fields['kind'] = kind
                    else:
                        extra_fields['func_name'] = def_name
                        extra_fields['receiver_type'] = receiver_type
                        extra_fields['is_method'] = (def_kind == 'method')
                    issues.append(ValidationIssue(
                        'heading_format',
                        file_path,
                        start_line,
                        end_line,
                        error_msg,
                        severity='error',
                        suggestion=suggestion,
                        heading=heading,
                        **extra_fields
                    ))


def check_missing_comment(
    code: str,
    code_lines: List[str],
    start_line: int,
    lines: List[str],
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
        issues.append(ValidationIssue(
            'missing_comment',
            file_path,
            def_line_num,
            def_line_num,
            f'{def_kind.capitalize()} definition `{def_display}` '
            'does not have a preceding comment',
            severity='error',
            suggestion=(
                f'Add a comment before the {def_kind} definition '
                f'`{def_display}` to document its purpose'
            ),
            heading=heading,
            def_name=def_display,
            def_kind=def_kind
        ))


def suggest_heading(heading: str, search_term: str, kind_word: str) -> str:
    """
    Suggest a corrected heading for a definition.

    Normalizes the heading to the format: [number] search_term kind_word remaining_text
    by removing backticks and preserving the numbering prefix (if present).

    Args:
        heading: Current heading text (may include numbering like "2.5 Heading Text")
        search_term: The term to search for (e.g., "Package", "FileEntry.GetState", "NewPackage")
        kind_word: The kind word to include (e.g., "Method", "Function", "Struct", "Interface")

    Returns:
        Suggested heading in format: {number} {search_term} {kind_word} {remaining_text}
        (number is preserved if present in original heading)
    """
    import re

    # Pattern to match numbered headings: "2.5" or "2.5.3" etc.
    numbered_pattern = re.compile(r'^([0-9]+(?:\.[0-9]+)*)\.?\s+(.+)$')

    # Extract numbering prefix if present
    numbering_prefix = None
    heading_without_number = heading
    match = numbered_pattern.match(heading.strip())
    if match:
        numbering_prefix = match.group(1)
        heading_without_number = match.group(2)

    # Step 1: Remove all backticks from heading text (not from numbering)
    heading_without_number = heading_without_number.replace('`', '')

    # Step 2: Remove search_term (case-insensitive)
    if '.' in search_term:
        # Method format - exact match
        pattern = re.compile(re.escape(search_term), re.IGNORECASE)
    else:
        # Simple name - use word boundaries
        pattern = re.compile(r'\b' + re.escape(search_term) + r'\b', re.IGNORECASE)
    heading_without_number = pattern.sub('', heading_without_number, count=1)

    # Step 3: Remove kind_word (case-insensitive, whole word)
    kind_pattern = re.compile(r'\b' + re.escape(kind_word) + r'\b', re.IGNORECASE)
    heading_without_number = kind_pattern.sub('', heading_without_number, count=1)

    # Step 4: Clean up extra whitespace
    remaining = ' '.join(heading_without_number.split())

    # Step 5: Build the suggested heading
    suggested = f"{search_term} {kind_word} {remaining}".strip()

    # Step 6: Add numbering prefix if it was present
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
    - Functions: Heading should include function name NOT in backticks (e.g., NewPackage)
      and should include "Function" immediately after the name
    - Methods: Heading should include type and method name in format
      Type.MethodName NOT in backticks (e.g., FileEntry.GetProcessingState)
      and should include "Method" immediately after Type.MethodName
    - Types/Interfaces/Structs: Heading should include type name NOT in backticks (e.g., Package)
      and should include the kind word ("Interface", "Struct", or "Type") immediately after the name

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
    # This preserves the original kind (e.g., 'alias' -> 'Alias', not 'Type')
    kind_word = kind.capitalize()

    # Check for backticks in heading (common check for all cases)
    if has_backticks(heading):
        errors.append(get_backticks_error_message())

    # Determine search term and display name based on kind
    if kind == 'method' and receiver_type:
        search_term = f'{receiver_type}.{name}'
        display_name = search_term
    else:
        search_term = name
        display_name = name

    # Check if search term is present in heading (common for all cases except methods)
    if kind != 'method' and name not in heading:
        errors.append(
            f'Heading should include {kind.capitalize()} name: {name}'
        )

    match kind:
        case 'method':
            if receiver_type:
                if search_term not in heading:
                    errors.append(
                        f'Method heading should include {display_name} (without backticks)'
                    )
                else:
                    check_kind_word_after(
                        heading, search_term, kind_word, display_name,
                        'Method heading', errors
                    )

        case _:  # Default: types and functions
            if name in heading:
                check_kind_word_after(
                    heading, search_term, kind_word, display_name,
                    f'{kind.capitalize()} heading', errors
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

    return (len(errors) == 0, errors, suggestion)


def audit_file(file_path: Path) -> Dict:
    """Audit a single markdown file for Go code block issues."""
    issues: List[ValidationIssue] = []
    code_blocks = []

    try:
        content = file_path.read_text(encoding='utf-8')
        lines = content.split('\n')
        blocks = find_go_code_blocks(content, str(file_path))

        # Track headings used by code blocks
        heading_usage = defaultdict(list)

        for start_line, end_line, code, heading in blocks:
            counts = count_go_definitions(code)
            type_count = counts['type']
            func_count = counts['func'] + counts['method']
            func_type_count = counts['func_type']

            # Cache code_lines once per code block to avoid repeated splits
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

            # Check if this is example code (skip validation for example code blocks)
            is_example = is_example_code(
                code, start_line,
                lines=lines,
                heading_text=heading,
                check_prose_before_block=True
            )

            # Skip all validation checks if this is example code
            if is_example:
                # Track heading usage even for example code
                if heading:
                    heading_usage[heading].append((start_line, end_line))
                continue

            # Check: at most one type/interface definition
            if type_count > 1:
                issues.append(ValidationIssue(
                    'multiple_types',
                    file_path,
                    start_line,
                    end_line,
                    f'Code block has {type_count} type/interface definitions '
                    '(max 1 allowed)',
                    severity='error',
                    suggestion=(
                        'Split into separate code blocks, '
                        'one per type/interface definition'
                    ),
                    heading=heading,
                    type_count=type_count
                ))

            # Check: at most one function/method definition
            if func_count > 1:
                issues.append(ValidationIssue(
                    'multiple_funcs',
                    file_path,
                    start_line,
                    end_line,
                    f'Code block has {func_count} func definitions '
                    '(max 1 allowed)',
                    severity='error',
                    suggestion=(
                        'Split into separate code blocks, '
                        'one per function/method definition'
                    ),
                    heading=heading,
                    func_count=func_count
                ))

            # Check: type and func definitions are mutually exclusive
            if type_count > 0 and func_count > 0:
                issues.append(ValidationIssue(
                    'type_func_exclusive',
                    file_path,
                    start_line,
                    end_line,
                    f'Code block has both {type_count} type definition(s) and '
                    f'{func_count} func definition(s) (must be exclusive)',
                    severity='error',
                    suggestion=(
                        'Separate type definitions and function definitions '
                        'into different code blocks'
                    ),
                    heading=heading,
                    type_count=type_count,
                    func_count=func_count
                ))

            # Check: function types (warnings for review)
            if func_type_count > 0:
                issues.append(ValidationIssue(
                    'function_type_warning',
                    file_path,
                    start_line,
                    end_line,
                    f'Code block has {func_type_count} function type definition(s) '
                    '(review recommended)',
                    severity='warning',
                    suggestion=(
                        'Review if function type definitions are intentional '
                        'and properly documented'
                    ),
                    heading=heading,
                    func_type_count=func_type_count
                ))

            # Check: type definitions have preceding comments
            if type_count > 0:
                check_missing_comment(
                    code, code_lines, start_line, lines, heading,
                    is_type=True, file_path=file_path, issues=issues
                )

            # Check: function/method definitions have preceding comments
            if func_count > 0:
                check_missing_comment(
                    code, code_lines, start_line, lines, heading,
                    is_type=False, file_path=file_path, issues=issues
                )

            # Check: heading format and comments for single definitions
            # Unified validation for both types and functions
            if (type_count == 1 and func_count == 0) or (func_count == 1 and type_count == 0):
                is_type = type_count == 1
                validate_single_definition(
                    code, code_lines, start_line, end_line, lines,
                    heading, is_type, file_path, issues
                )

            # Track heading usage
            if heading:
                heading_usage[heading].append((start_line, end_line))

        # Check: each code block should be under a different heading
        for heading, blocks_under_heading in heading_usage.items():
            if len(blocks_under_heading) > 1:
                # Calculate line range: first block start to last block end
                first_block_start = min(block[0] for block in blocks_under_heading)
                last_block_end = max(block[1] for block in blocks_under_heading)

                issues.append(ValidationIssue(
                    'multiple_blocks_per_heading',
                    file_path,
                    first_block_start,
                    last_block_end,
                    (
                        f'Heading "{heading}" has {len(blocks_under_heading)} '
                        'Go code blocks (each should be under a different heading)'
                    ),
                    severity='error',
                    suggestion='Move each code block to a separate heading',
                    heading=heading,
                    blocks=blocks_under_heading
                ))

        # Check for code blocks without headings
        for block in code_blocks:
            if block['heading'] is None:
                issues.append(ValidationIssue(
                    'no_heading',
                    file_path,
                    block['start_line'],
                    block['end_line'],
                    'Code block is not under any heading',
                    severity='error',
                    suggestion='Add a heading above the code block'
                ))

    except Exception as e:
        issues.append(ValidationIssue(
            'error',
            file_path,
            1,
            1,
            f'Error reading file: {e}',
            severity='error'
        ))

    return {
        'file': str(file_path),
        'code_blocks': code_blocks,
        'issues': issues
    }


def generate_report(results: List[Dict], output_path: Path) -> None:
    """Generate markdown report from audit results."""
    report_lines = []

    report_lines.append('# Go Code Blocks Validation Report')
    report_lines.append('')
    report_lines.append('This report validates all Go code blocks in the tech specs documentation.')
    report_lines.append('')
    report_lines.append('## Requirements')
    report_lines.append('')
    report_lines.append(
        '1. Each Go code block should have at most one type or interface '
        'definition'
    )
    report_lines.append(
        '2. Each Go code block should have at most one function or method '
        'definition'
    )
    report_lines.append(
        '3. Type definitions and function definitions are mutually exclusive '
        'in a code block'
    )
    report_lines.append('4. Each Go code block should be under a different heading')
    report_lines.append(
        '5. Function headings should NOT include the function name in backticks '
        '(e.g., NewPackage, not `NewPackage`)'
    )
    report_lines.append(
        '6. Method headings should NOT include type and method name in backticks '
        '(e.g., FileEntry.GetProcessingState, not `FileEntry.GetProcessingState`)'
    )
    report_lines.append(
        '7. Type/Interface/Struct headings should NOT include the type name in backticks '
        '(e.g., Package, not `Package`)'
    )
    report_lines.append(
        '8. All type, interface, struct, function, and method definitions should have '
        'comments preceding them'
    )
    report_lines.append('')

    # Summary
    total_files = len(results)
    total_blocks = sum(len(r['code_blocks']) for r in results)
    total_issues = sum(len(r['issues']) for r in results)

    report_lines.append('## Summary')
    report_lines.append('')
    report_lines.append(f'- Files audited: {total_files}')
    report_lines.append(f'- Total Go code blocks found: {total_blocks}')
    report_lines.append(f'- Total issues found: {total_issues}')
    report_lines.append('')

    if total_issues == 0:
        report_lines.append('✅ All Go code blocks comply with the requirements!')
        report_lines.append('')
    else:
        report_lines.append('## Issues Found')
        report_lines.append('')

        # Group issues by type
        issues_by_type = defaultdict(list)
        for result in results:
            for issue in result['issues']:
                # Convert ValidationIssue to dict if needed
                if isinstance(issue, ValidationIssue):
                    issue = issue.to_dict()
                issues_by_type[issue['type']].append((result['file'], issue))

        for issue_type, issues in sorted(issues_by_type.items()):
            report_lines.append(f'### {issue_type.replace("_", " ").title()} Issues')
            report_lines.append('')

            for file_path, issue in issues:
                report_lines.append(f'**File:** `{file_path}`')
                if 'start_line' in issue:
                    report_lines.append(f'**Lines:** {issue["start_line"]}-{issue["end_line"]}')
                if 'heading' in issue:
                    report_lines.append(f'**Heading:** {issue["heading"]}')
                if 'type_count' in issue:
                    report_lines.append(f'**Type definitions found:** {issue["type_count"]}')
                if 'func_count' in issue:
                    report_lines.append(f'**Func definitions found:** {issue["func_count"]}')
                if 'func_type_count' in issue:
                    report_lines.append(
                        f'**Function type definitions found:** '
                        f'{issue["func_type_count"]}'
                    )
                if 'blocks' in issue:
                    block_info = ', '.join(f'lines {s}-{e}' for s, e in issue['blocks'])
                    report_lines.append(f'**Code blocks:** {block_info}')
                if 'def_name' in issue:
                    report_lines.append(f'**Definition:** {issue["def_name"]}')
                if 'def_kind' in issue:
                    report_lines.append(f'**Definition kind:** {issue["def_kind"]}')
                report_lines.append(f'**Issue:** {issue["message"]}')
                report_lines.append('')

        # Detailed file-by-file breakdown
        report_lines.append('## Detailed File Breakdown')
        report_lines.append('')

        for result in sorted(results, key=lambda x: x['file']):
            if result['code_blocks'] or result['issues']:
                file_name = Path(result["file"]).stem
                report_lines.append(f'### {file_name}')
                report_lines.append('')
                report_lines.append(f'**File path:** `{result["file"]}`')
                report_lines.append(f'**Code blocks:** {len(result["code_blocks"])}')
                report_lines.append(f'**Issues:** {len(result["issues"])}')
                report_lines.append('')

                if result['code_blocks']:
                    report_lines.append(f'#### {file_name} Code Blocks')
                    report_lines.append('')
                    for i, block in enumerate(result['code_blocks'], 1):
                        report_lines.append(
                            f'Code block {i}: Lines '
                            f'{block["start_line"]}-{block["end_line"]}'
                        )
                        report_lines.append('')
                        report_lines.append(f'- Heading: {block["heading"] or "(none)"}')
                        report_lines.append(f'- Type definitions: {block["type_count"]}')
                        report_lines.append(f'- Func definitions: {block["func_count"]}')
                        if block.get("func_type_count", 0) > 0:
                            report_lines.append(
                                f'- Function type definitions: '
                                f'{block["func_type_count"]}'
                            )
                        report_lines.append('')

                if result['issues']:
                    report_lines.append(f'#### {file_name} Issues')
                    report_lines.append('')
                    for issue in result['issues']:
                        report_lines.append(f'- {issue["message"]}')
                        if 'start_line' in issue:
                            report_lines.append(
                                f'  - Lines: {issue["start_line"]}-'
                                f'{issue["end_line"]}'
                            )
                    report_lines.append('')

    # Write report
    output_path.parent.mkdir(parents=True, exist_ok=True)
    output_path.write_text('\n'.join(report_lines), encoding='utf-8')


def print_summary(results, output=None, verbose=False, no_color=False):
    """
    Print summary of audit results.

    Args:
        results: List of audit results
        output: Optional OutputBuilder instance (creates new one if None)
        verbose: Verbose mode flag
        no_color: Disable colors flag
    """
    if output is None:
        output = OutputBuilder(no_color=no_color, verbose=verbose)
        # If creating new output, add header
        output.add_header("Go Code Blocks Validation",
                          "Validates Go code blocks in tech specs")

    total_files = len(results)
    total_blocks = sum(len(r['code_blocks']) for r in results)
    total_issues = sum(len(r['issues']) for r in results)

    # Summary section
    output.add_summary_header()
    summary_items = [
        ("Files audited:", total_files),
        ("Total code blocks:", total_blocks),
        ("Total issues found:", total_issues),
    ]
    output.add_summary_section(summary_items)

    # Group issues by type (includes all types: duplicate_heading, multiple_funcs, etc.)
    issues_by_type = defaultdict(int)
    for result in results:
        for issue in result['issues']:
            # Convert ValidationIssue to get type
            if isinstance(issue, ValidationIssue):
                issue_type = issue.issue_type
            else:
                issue_type = issue.get('type', 'unknown')
            issues_by_type[issue_type] += 1

    # Print breakdown of all issue types found (non-verbose summary)
    if issues_by_type:
        output.add_blank_line("summary")
        output.add_line('Breakdown by issue type:', section="summary")
        breakdown_items = [
            (issue_type.replace('_', ' ').title() + ':', count)
            for issue_type, count in sorted(issues_by_type.items())
        ]
        output.add_summary_section(breakdown_items)

    # Group issues by type for display - uses format_issue_message
    # Errors are always shown, warnings are verbose-only
    if issues_by_type:
        # Only add errors header if there are errors (not just warnings)
        has_errors = any(
            (isinstance(issue, ValidationIssue) and
             issue.issue_type != "function_type_warning") or
            (not isinstance(issue, ValidationIssue) and
             issue.get('type') != "function_type_warning")
            for result in results
            for issue in result['issues']
        )
        if has_errors:
            output.add_errors_header()
            output.add_blank_line("error")

        issues_by_type_list = defaultdict(list)
        for result in results:
            for issue in result['issues']:
                # Convert ValidationIssue to get type
                if isinstance(issue, ValidationIssue):
                    issue_type = issue.issue_type
                    issue_dict = issue.to_dict()
                else:
                    issue_type = issue.get('type', 'unknown')
                    issue_dict = issue
                issues_by_type_list[issue_type].append((result['file'], issue_dict))

        for issue_type, issues in sorted(issues_by_type_list.items()):
            # Determine severity based on issue type
            severity = "warning" if issue_type == "function_type_warning" else "error"

            for file_path, issue in issues:
                # Build message details
                message_parts = []
                if 'heading' in issue:
                    message_parts.append(f'Heading: {issue["heading"]}')
                if 'type_count' in issue:
                    message_parts.append(f'Type definitions: {issue["type_count"]}')
                if 'func_count' in issue:
                    message_parts.append(f'Func definitions: {issue["func_count"]}')
                if 'func_type_count' in issue:
                    message_parts.append(f'Function type definitions: {issue["func_type_count"]}')
                if 'func_name' in issue:
                    message_parts.append(f'Function/Method: {issue["func_name"]}')
                if 'receiver_type' in issue and issue['receiver_type']:
                    message_parts.append(f'Receiver: {issue["receiver_type"]}')
                if 'type_name' in issue:
                    message_parts.append(f'Type: {issue["type_name"]}')
                if 'kind' in issue:
                    message_parts.append(f'Kind: {issue["kind"]}')
                if 'def_name' in issue:
                    message_parts.append(f'Definition: {issue["def_name"]}')
                if 'def_kind' in issue:
                    message_parts.append(f'Definition kind: {issue["def_kind"]}')
                if 'blocks' in issue:
                    block_info = ', '.join(f'lines {s}-{e}' for s, e in issue['blocks'])
                    message_parts.append(f'Code blocks: {block_info}')

                # Use the main message, append details if any
                message = issue.get('message', '')
                if message_parts:
                    message = f"{message} ({', '.join(message_parts)})"

                # Format issue type name
                issue_type_name = issue_type.replace('_', ' ').title()

                # Use format_issue_message
                line_num = issue.get('start_line')
                suggestion = issue.get('suggestion')
                formatted_msg = format_issue_message(
                    severity=severity,
                    issue_type=issue_type_name,
                    file_path=file_path,
                    line_num=line_num,
                    message=message,
                    suggestion=suggestion,
                    no_color=no_color
                )

                # Add to appropriate output section
                # Errors are always shown, warnings are verbose-only
                if severity == "warning":
                    output.add_warning_line(formatted_msg, verbose_only=True)
                else:
                    output.add_error_line(formatted_msg, verbose_only=False)

    # Add final message (mutually exclusive)
    if total_issues == 0:
        output.add_success_message("All Go code blocks comply with the requirements!")
    else:
        output.add_failure_message("Validation failed. Please fix the errors above.")

    return output


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

    # Exit with error code if issues found (unless --no-fail is set)
    total_issues = sum(len(r['issues']) for r in results)
    has_errors = total_issues > 0
    return get_validation_exit_code(has_errors, no_fail)


if __name__ == '__main__':
    sys.exit(main())
