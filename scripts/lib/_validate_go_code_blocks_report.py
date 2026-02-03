"""Report generation helpers for Go code blocks validation."""

from collections import defaultdict
from pathlib import Path
from typing import Dict, List

from lib._validation_utils import (
    OutputBuilder,
    ValidationIssue,
    format_issue_message,
)


def _append_issues_found_section(
    report_lines: List[str], results: List[Dict]
) -> None:
    """Append '## Issues Found' and grouped issues to report_lines."""
    report_lines.append('## Issues Found')
    report_lines.append('')

    issues_by_type = defaultdict(list)
    for result in results:
        for issue in result['issues']:
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


def _append_detailed_breakdown_section(
    report_lines: List[str], results: List[Dict]
) -> None:
    """Append '## Detailed File Breakdown' to report_lines."""
    report_lines.append('## Detailed File Breakdown')
    report_lines.append('')

    for result in sorted(results, key=lambda x: x['file']):
        if not (result['code_blocks'] or result['issues']):
            continue
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
                issue_dict = issue.to_dict() if isinstance(issue, ValidationIssue) else issue
                report_lines.append(f'- {issue_dict["message"]}')
                if 'start_line' in issue_dict:
                    report_lines.append(
                        f'  - Lines: {issue_dict["start_line"]}-'
                        f'{issue_dict["end_line"]}'
                    )
            report_lines.append('')


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
        '5. Headings for Go definitions should include the definition name and kind word; '
        'definition names are preferred in backticks (e.g. `` `Package.Write` Method ``). '
        'Case inside backticks is ignored for validation.'
    )
    report_lines.append(
        '6. All type, interface, struct, function, and method definitions should have '
        'comments preceding them'
    )
    report_lines.append('')

    total_files = len(results)
    total_blocks = sum(len(r['code_blocks']) for r in results)
    total_issues = sum(len(r['issues']) for r in results)

    report_lines.append('## Summary')
    report_lines.append('')
    report_lines.append(f'- Files audited: {total_files}')
    report_lines.append(f'- Total Go code blocks found: {total_blocks}')
    report_lines.append(f'- Total issues found: {total_issues}')
    report_lines.append('')

    if not total_issues:
        report_lines.append('✅ All Go code blocks comply with the requirements!')
        report_lines.append('')
    else:
        _append_issues_found_section(report_lines, results)
        _append_detailed_breakdown_section(report_lines, results)

    output_path.parent.mkdir(parents=True, exist_ok=True)
    output_path.write_text('\n'.join(report_lines), encoding='utf-8')


def _message_parts_for_issue(issue) -> List[str]:
    """Build list of message part strings from an issue (ValidationIssue or dict)."""
    if isinstance(issue, ValidationIssue):
        issue = issue.to_dict()
    parts = []
    key_labels = [
        ('heading', 'Heading'),
        ('type_count', 'Type definitions'),
        ('func_count', 'Func definitions'),
        ('func_type_count', 'Function type definitions'),
        ('func_name', 'Function/Method'),
        ('receiver_type', 'Receiver'),
        ('type_name', 'Type'),
        ('kind', 'Kind'),
        ('def_name', 'Definition'),
        ('def_kind', 'Definition kind'),
    ]
    for key, label in key_labels:
        if key in issue and issue[key]:
            parts.append(f'{label}: {issue[key]}')
    if 'blocks' in issue and issue['blocks']:
        block_info = ', '.join(f'lines {s}-{e}' for s, e in issue['blocks'])
        parts.append(f'Code blocks: {block_info}')
    return parts


def _issues_by_type(results) -> Dict[str, int]:
    """Return dict of issue_type -> count from results."""
    counts = defaultdict(int)
    for result in results:
        for issue in result['issues']:
            it = (
                issue.issue_type if isinstance(issue, ValidationIssue)
                else issue.get('type', 'unknown')
            )
            counts[it] += 1
    return dict(counts)


WARNING_ISSUE_TYPES = ("function_type_warning", "heading_prefer_backticks")


def _has_non_warning_errors(results) -> bool:
    """Return True if any issue is not a warning (e.g. heading_prefer_backticks)."""
    return any(
        (isinstance(i, ValidationIssue) and i.issue_type not in WARNING_ISSUE_TYPES) or
        (not isinstance(i, ValidationIssue) and i.get('type') not in WARNING_ISSUE_TYPES)
        for r in results for i in r['issues']
    )


def _issues_by_type_list(results) -> Dict[str, List[tuple]]:
    """Return dict of issue_type -> [(file_path, issue_dict), ...]."""
    out = defaultdict(list)
    for result in results:
        for issue in result['issues']:
            if isinstance(issue, ValidationIssue):
                it, issue_dict = issue.issue_type, issue.to_dict()
            else:
                it, issue_dict = issue.get('type', 'unknown'), issue
            out[it].append((result['file'], issue_dict))
    return dict(out)


def _emit_issues_section(output, results, no_color) -> None:
    """Emit errors header and all issues grouped by type."""
    issues_by_type_list = _issues_by_type_list(results)
    for issue_type, issues in sorted(issues_by_type_list.items()):
        severity = "warning" if issue_type in WARNING_ISSUE_TYPES else "error"
        # Show backtick recommendations in default output (not verbose-only)
        show_by_default = severity == "warning" and issue_type == "heading_prefer_backticks"
        for file_path, issue in issues:
            message_parts = _message_parts_for_issue(issue)
            message = issue.get('message', '')
            if message_parts:
                message = f"{message} ({', '.join(message_parts)})"
            formatted_msg = format_issue_message(
                severity=severity,
                issue_type=issue_type.replace('_', ' ').title(),
                file_path=file_path,
                line_num=issue.get('start_line'),
                message=message,
                suggestion=issue.get('suggestion'),
                no_color=no_color,
            )
            if severity == "warning":
                output.add_warning_line(formatted_msg, verbose_only=not show_by_default)
            else:
                output.add_error_line(formatted_msg, verbose_only=False)


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
        output = OutputBuilder(
            "Go Code Blocks Validation",
            "Validates Go code blocks in tech specs",
            no_color=no_color,
            verbose=verbose,
        )
        output.add_header("Go Code Blocks Validation", "Validates Go code blocks in tech specs")

    total_files = len(results)
    total_blocks = sum(len(r['code_blocks']) for r in results)
    total_issues = sum(len(r['issues']) for r in results)
    output.add_summary_header()
    output.add_summary_section([
        ("Files audited:", total_files),
        ("Total code blocks:", total_blocks),
        ("Total issues found:", total_issues),
    ])

    issues_by_type = _issues_by_type(results)
    if issues_by_type:
        output.add_blank_line("summary")
        output.add_line('Breakdown by issue type:', section="summary")
        output.add_summary_section([
            (t.replace('_', ' ').title() + ':', c)
            for t, c in sorted(issues_by_type.items())
        ])
    if issues_by_type:
        if _has_non_warning_errors(results):
            output.add_errors_header()
            output.add_blank_line("error")
        _emit_issues_section(output, results, no_color)

    if not total_issues:
        output.add_success_message("All Go code blocks comply with the requirements!")
    elif _has_non_warning_errors(results):
        output.add_failure_message("Validation failed. Please fix the errors above.")
    else:
        output.add_warnings_only_message(
            message="All Go code blocks comply (see recommendations above).",
        )
    return output
