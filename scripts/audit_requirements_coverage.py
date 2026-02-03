#!/usr/bin/env python3
"""
Audit requirements coverage for tech specs.

This script checks which tech specs have requirements that reference
them using inline markdown links. It also checks that each H2+ heading
in tech specs has at least one requirement referencing it.

Usage:
    python3 scripts/audit_requirements_coverage.py [options]

Options:
    --verbose, -v       Show detailed progress information
    --path, -p PATHS    Check only the specified file(s) or
                        directory(ies) (recursive). Can be a single
                        path or comma-separated list of paths
    --help, -h          Show this help message

Examples:
    # Basic audit
    python3 scripts/audit_requirements_coverage.py

    # Verbose output
    python3 scripts/audit_requirements_coverage.py --verbose

    # Check specific file
    python3 scripts/audit_requirements_coverage.py --path docs/requirements/core.md

    # Check specific directory
    python3 scripts/audit_requirements_coverage.py --path docs/requirements

    # Check multiple paths
    python3 scripts/audit_requirements_coverage.py --path \\
        docs/requirements/core.md,docs/requirements/security.md
"""

import sys
import argparse
import re
from dataclasses import dataclass
from pathlib import Path
from typing import List, Optional, Set, Tuple

from lib._validation_utils import (
    OutputBuilder, parse_no_color_flag, format_issue_message,
    ValidationIssue,
    get_workspace_root, parse_paths,
    build_heading_hierarchy, is_organizational_heading,
    extract_h2_plus_headings_with_sections,
    FileContentCache, DOCS_DIR, TECH_SPECS_DIR, REQUIREMENTS_DIR
)
import lib._go_code_utils as go_code_utils

from lib.audit_requirements._audit_requirements_config import (
    MAX_ORGANIZATIONAL_PROSE_LINES,
    HEADING_EXCLUSION_PATTERNS,
    FUNCTIONAL_BEHAVIOR_KEYWORDS,
    _SCRIPT_ERROR_EXCEPTIONS,
)
from lib.audit_requirements._audit_requirements_scan import (
    find_tech_specs,
    check_index_file_references,
    get_requirement_files,
    check_heading_referenced,
    get_requirement_spec_anchor_links,
    count_requirement_references,
    check_anchor_in_text_missing_href_anchor,
    check_requirement_tech_spec_link_thresholds,
)
from lib.audit_requirements._audit_requirements_classify import (
    extract_section_content,
    classify_heading,
)


@dataclass(frozen=True)
class _SpecHeadingContext:
    spec_content: str
    spec_relative_path: Path
    spec_basename: str
    headings_for_hierarchy: list
    spec_lines: list
    hierarchy: list


@dataclass(frozen=True)
class _HeadingInfo:
    heading_level: int
    heading_text: str
    line_num: int
    anchor: str
    section_anchor: Optional[str]


def extract_h2_plus_headings(spec_file, file_cache=None):
    """Extract all H2+ headings from a tech spec file."""
    file_path = spec_file if isinstance(spec_file, Path) else Path(spec_file)
    return extract_h2_plus_headings_with_sections(file_path, file_cache=file_cache)


def _run_index_ref_check(
    requirements_dir: Path,
    file_cache: FileContentCache,
    verbose: bool,
    target_paths: Optional[List[str]],
    *,
    output: OutputBuilder,
    no_color: bool,
) -> Optional[List]:
    """Check for index file references. Return list of errors or None if ok."""
    index_ref_errors = check_index_file_references(
        requirements_dir, file_cache, verbose, target_paths
    )
    if not index_ref_errors:
        return None
    output.add_errors_header()
    output.add_line(
        "Requirements should NOT reference index files (*_index.md).",
        section="error"
    )
    output.add_line(
        "Index files are navigation aids, not source specifications.",
        section="error"
    )
    output.add_blank_line("error")
    for req_file, line_num, reference in index_ref_errors:
        error_msg = format_issue_message(
            "error", "Index file ref", req_file,
            line_num=line_num, message=reference, no_color=no_color
        )
        output.add_error_line(error_msg)
    output.add_blank_line("error")
    output.add_line(
        "Note: After fixing these references, verify that each updated "
        "reference points to the correct content.",
        section="error"
    )
    return index_ref_errors


def _run_file_level_coverage(
    tech_specs: List[Path],
    requirements_dir: Path,
    file_cache: FileContentCache,
    target_paths: Optional[List[str]],
    *,
    output: OutputBuilder,
    args,
    no_color: bool,
) -> List[str]:
    """Check file-level coverage; output errors; return missing_specs."""
    missing_specs = []

    if args.verbose:
        output.add_verbose_line("Checking File-Level Coverage...")
        output.add_blank_line("working_verbose")

    for spec in tech_specs:
        spec_basename = spec.name
        if args.verbose:
            output.add_verbose_line(f"Checking {spec_basename}...")

        count = count_requirement_references(
            spec_basename, requirements_dir, file_cache, args.verbose, target_paths
        )

        if not count:
            error_msg = format_issue_message(
                "error", "Spec without Req", spec_basename,
                message="NO REQUIREMENTS", no_color=no_color
            )
            output.add_error_line(error_msg)
            missing_specs.append(spec_basename)
        elif args.verbose:
            output.add_verbose_line(f"✓  {spec_basename}: {count} requirement(s)")

    if missing_specs:
        output.add_errors_header()
        for spec_basename in missing_specs:
            error_msg = format_issue_message(
                "error", "Spec without Req", spec_basename,
                message="NO REQUIREMENTS", no_color=no_color
            )
            output.add_error_line(error_msg)

    return missing_specs


def _run_heading_coverage(
    tech_specs: List[Path],
    requirement_files: List[Path],
    file_cache: FileContentCache,
    workspace_root: Path,
    *,
    output: OutputBuilder,
    args,
    no_color: bool,
) -> Tuple[List[ValidationIssue], int, Set[Tuple[str, str]]]:
    """Check heading-level coverage; return (issues, total_headings, excluded_headings_set)."""
    issues: List[ValidationIssue] = []
    total_headings = 0
    excluded_headings_set: Set[Tuple[str, str]] = set()

    if args.verbose:
        output.add_verbose_line("Checking Heading-Level Coverage (H2+)")
        output.add_blank_line("working_verbose")

    for spec in tech_specs:
        spec_basename = spec.name
        try:
            spec_relative_path = spec.relative_to(workspace_root)
        except ValueError:
            spec_relative_path = spec

        try:
            spec_content = file_cache.get_content(spec)
        except (IOError, OSError) as e:
            if args.verbose:
                output.add_verbose_line(f"  Warning: Could not read {spec_basename}: {e}")
            issues.append(ValidationIssue.create(
                "file_read_error", spec_relative_path, 0, 0,
                message=f"Could not read file: {e}",
                severity='error'
            ))
            continue
        except UnicodeDecodeError as e:
            if args.verbose:
                output.add_verbose_line(
                    f"  Warning: Could not decode {spec_basename} (encoding issue): {e}"
                )
            issues.append(ValidationIssue.create(
                "file_encoding_error", spec_relative_path, 0, 0,
                message=f"Could not decode file (encoding issue): {e}",
                severity='error'
            ))
            continue
        except _SCRIPT_ERROR_EXCEPTIONS as e:
            if args.verbose:
                output.add_verbose_line(f"  Warning: Unexpected error reading {spec_basename}: {e}")
            issues.append(ValidationIssue.create(
                "unexpected_error", spec_relative_path, 0, 0,
                message=f"Unexpected error reading file: {e}",
                severity='error'
            ))
            continue

        part_issues, part_total, part_excluded = _process_headings_for_spec(
            spec, spec_content, spec_relative_path, spec_basename,
            requirement_files=requirement_files, file_cache=file_cache,
            output=output, args=args, no_color=no_color
        )
        issues.extend(part_issues)
        total_headings += part_total
        excluded_headings_set.update(part_excluded)

    return (issues, total_headings, excluded_headings_set)


def _process_one_heading(
    ctx: _SpecHeadingContext,
    heading: _HeadingInfo,
    *,
    requirement_files: List[Path],
    file_cache: FileContentCache,
    output: OutputBuilder,
    args,
    no_color: bool,
) -> Tuple[List[ValidationIssue], Set[Tuple[str, str]]]:
    """Process one heading; return (issues_to_append, excluded_pairs)."""
    issues: List[ValidationIssue] = []
    excluded: Set[Tuple[str, str]] = set()

    try:
        section_content = extract_section_content(
            ctx.spec_content, heading.line_num, ctx.headings_for_hierarchy, lines=ctx.spec_lines
        )
    except _SCRIPT_ERROR_EXCEPTIONS as e:
        error_msg = format_issue_message(
            "error", "Analysis error", ctx.spec_relative_path,
            line_num=heading.line_num,
            message=(
                f"Failed to extract section content for heading "
                f"'{heading.heading_text}': {str(e)}"
            ),
            no_color=no_color
        )
        output.add_error_line(error_msg)
        issues.append(ValidationIssue.create(
            "analysis_error", Path(ctx.spec_relative_path), heading.line_num, heading.line_num,
            message=(
                f"Failed to extract section content for heading "
                f"'{heading.heading_text}': {str(e)}"
            ),
            severity='error',
            heading_level=heading.heading_level,
            heading_text=heading.heading_text,
            anchor=heading.anchor
        ))
        return (issues, excluded)

    org_result = {'is_organizational': False, 'is_empty': False, 'sentence_count': 0}
    try:
        org_result = is_organizational_heading(
            ctx.spec_content,
            heading.line_num,
            heading.heading_level,
            ctx.headings_for_hierarchy,
            ctx.hierarchy,
            max_prose_lines=MAX_ORGANIZATIONAL_PROSE_LINES
        )
    except (ValueError, IndexError, KeyError) as e:
        error_msg = format_issue_message(
            "error", "Analysis error", ctx.spec_relative_path,
            line_num=heading.line_num,
            message=(
                f"Failed to check if organizational for heading "
                f"'{heading.heading_text}': {str(e)}"
            ),
            no_color=no_color
        )
        output.add_error_line(error_msg)
        issues.append(ValidationIssue.create(
            "analysis_error", Path(ctx.spec_relative_path), heading.line_num, heading.line_num,
            message=(
                f"Failed to check if organizational for heading "
                f"'{heading.heading_text}': {str(e)}"
            ),
            severity='error',
            heading_level=heading.heading_level,
            heading_text=heading.heading_text,
            anchor=heading.anchor
        ))
        return (issues, excluded)
    except _SCRIPT_ERROR_EXCEPTIONS as e:
        error_msg = format_issue_message(
            "error", "Analysis error", ctx.spec_relative_path,
            line_num=heading.line_num,
            message=(
                "Unexpected error checking organizational heading "
                f"'{heading.heading_text}': {str(e)}"
            ),
            no_color=no_color
        )
        output.add_error_line(error_msg)
        issues.append(ValidationIssue.create(
            "analysis_error", Path(ctx.spec_relative_path), heading.line_num, heading.line_num,
            message=(
                f"Unexpected error checking organizational heading "
                f"'{heading.heading_text}': {str(e)}"
            ),
            severity='error',
            heading_level=heading.heading_level,
            heading_text=heading.heading_text,
            anchor=heading.anchor
        ))
        return (issues, excluded)

    if org_result['is_organizational']:
        excluded.add((ctx.spec_basename, heading.anchor))
        if heading.section_anchor:
            excluded.add((ctx.spec_basename, heading.section_anchor))
        if not check_heading_referenced(
            ctx.spec_basename,
            heading.anchor,
            heading.section_anchor,
            requirement_files,
            file_cache=file_cache, verbose=args.verbose
        ):
            content_note = (
                " (no direct content)" if org_result['is_empty'] else " (minor content)"
            )
            issues.append(ValidationIssue.create(
                "organizational_heading", Path(ctx.spec_relative_path),
                heading.line_num, heading.line_num,
                message=f"{heading.heading_text} (#{heading.anchor}){content_note}",
                severity='warning',
                heading_level=heading.heading_level,
                heading_text=heading.heading_text,
                anchor=heading.anchor,
                is_empty=org_result['is_empty']
            ))
        return (issues, excluded)

    try:
        classification = classify_heading(
            heading.heading_text, heading.heading_level, section_content,
            HEADING_EXCLUSION_PATTERNS,
            _heading_line=heading.line_num,
            functional_keywords=FUNCTIONAL_BEHAVIOR_KEYWORDS,
            _max_org_lines=MAX_ORGANIZATIONAL_PROSE_LINES,
            go_code_utils_module=go_code_utils
        )
    except _SCRIPT_ERROR_EXCEPTIONS as e:
        error_msg = format_issue_message(
            "error", "Analysis error", ctx.spec_relative_path,
            line_num=heading.line_num,
            message=f"Failed to classify heading '{heading.heading_text}': {str(e)}",
            no_color=no_color
        )
        output.add_error_line(error_msg)
        issues.append(ValidationIssue.create(
            "analysis_error", Path(ctx.spec_relative_path), heading.line_num, heading.line_num,
            message=f"Failed to classify heading '{heading.heading_text}': {str(e)}",
            severity='error',
            heading_level=heading.heading_level,
            heading_text=heading.heading_text,
            anchor=heading.anchor
        ))
        return (issues, excluded)

    issue_or_none, excluded_pairs, skip_ref_check = _handle_classification_result(
        classification,
        ctx.spec_basename,
        heading.anchor,
        heading.section_anchor,
        ctx.spec_relative_path,
        line_num=heading.line_num,
        heading_text=heading.heading_text,
        heading_level=heading.heading_level
    )
    for pair in excluded_pairs:
        excluded.add(pair)
    if issue_or_none is not None:
        issues.append(issue_or_none)
        return (issues, excluded)
    if skip_ref_check:
        return (issues, excluded)

    if check_heading_referenced(
        ctx.spec_basename,
        heading.anchor,
        heading.section_anchor,
        requirement_files,
        file_cache=file_cache, verbose=args.verbose
    ):
        return (issues, excluded)
    severity = classification['severity_if_missing']
    issues.append(ValidationIssue.create(
        "missing_requirement", Path(ctx.spec_relative_path),
        heading.line_num, heading.line_num,
        message=f"{heading.heading_text} (#{heading.anchor}) - {classification['reason']}",
        severity=severity,
        heading_level=heading.heading_level,
        heading_text=heading.heading_text,
        anchor=heading.anchor,
        reason=classification['reason']
    ))
    return (issues, excluded)


def _process_headings_for_spec(
    spec: Path,
    spec_content: str,
    spec_relative_path: Path,
    spec_basename: str,
    *,
    requirement_files: List[Path],
    file_cache: FileContentCache,
    output: OutputBuilder,
    args,
    no_color: bool,
) -> Tuple[List[ValidationIssue], int, Set[Tuple[str, str]]]:
    """Process all headings in one spec; return (issues, total_headings, excluded_set)."""
    issues: List[ValidationIssue] = []
    total_headings = 0
    excluded_headings_set: Set[Tuple[str, str]] = set()

    headings = extract_h2_plus_headings(spec, file_cache)
    if not headings:
        return (issues, 0, excluded_headings_set)

    headings_for_hierarchy = [
        (line_num, heading_level, heading_text)
        for heading_level, heading_text, line_num, _anchor, _section_anchor in headings
    ]
    hierarchy = build_heading_hierarchy(headings_for_hierarchy)
    spec_lines = spec_content.split('\n')
    ctx = _SpecHeadingContext(
        spec_content=spec_content,
        spec_relative_path=spec_relative_path,
        spec_basename=spec_basename,
        headings_for_hierarchy=headings_for_hierarchy,
        spec_lines=spec_lines,
        hierarchy=hierarchy,
    )

    if args.verbose:
        output.add_verbose_line(f"Checking headings in {spec_basename}...")

    for heading_level, heading_text, line_num, anchor, section_anchor in headings:
        if heading_text.lower() == "table of contents":
            continue
        heading_text_lower = heading_text.lower()
        heading_without_section = re.sub(r'^\d+(?:\.\d+)*\s+', '', heading_text_lower).strip()
        if heading_without_section == "overview":
            continue
        section_match = re.match(r'^(\d+(?:\.\d+)*)', heading_text)
        if section_match:
            section_num = section_match.group(1)
            if section_num.startswith("0.") or section_num == "0":
                continue

        total_headings += 1
        part_issues, part_excluded = _process_one_heading(
            ctx,
            _HeadingInfo(
                heading_level=heading_level,
                heading_text=heading_text,
                line_num=line_num,
                anchor=anchor,
                section_anchor=section_anchor,
            ),
            requirement_files=requirement_files, file_cache=file_cache,
            output=output, args=args, no_color=no_color
        )
        issues.extend(part_issues)
        excluded_headings_set.update(part_excluded)

    return (issues, total_headings, excluded_headings_set)


def _handle_classification_result(
    classification: dict,
    spec_basename: str,
    anchor: str,
    section_anchor: Optional[str],
    spec_relative_path: Path,
    *,
    line_num: int = 0,
    heading_text: str = '',
    heading_level: int = 0,
) -> Tuple[Optional[ValidationIssue], List[Tuple[str, str]], bool]:
    """
    Map classification to issue (or None), excluded pairs to add, and whether to skip ref check.
    Returns (issue, [(spec_basename, anchor), ...], skip_ref_check).
    """
    excluded_pairs = [(spec_basename, anchor)]
    if section_anchor:
        excluded_pairs.append((spec_basename, section_anchor))

    kind = classification['classification']
    reason = classification['reason']
    spec_path = Path(spec_relative_path)

    if kind == 'architectural':
        return (
            ValidationIssue.create(
                "architectural_heading", spec_path, line_num, line_num,
                message=f"{heading_text} (architectural: {reason})",
                severity='warning',
                heading=heading_text, heading_level=heading_level, anchor=anchor, reason=reason
            ),
            excluded_pairs,
            True,
        )
    if kind == 'signature_only':
        return (
            ValidationIssue.create(
                "signature_only_heading", spec_path, line_num, line_num,
                message=f"{heading_text} (#{anchor}) - {reason}",
                severity='warning',
                heading_level=heading_level, heading_text=heading_text, anchor=anchor, reason=reason
            ),
            excluded_pairs,
            True,
        )
    if kind == 'example_only':
        return (
            ValidationIssue.create(
                "example_only_heading", spec_path, line_num, line_num,
                message=f"{heading_text} (#{anchor}) - {reason}",
                severity='warning',
                heading_level=heading_level, heading_text=heading_text, anchor=anchor, reason=reason
            ),
            excluded_pairs,
            True,
        )
    if kind == 'non_prose':
        return (
            ValidationIssue.create(
                "non_prose_heading", spec_path, line_num, line_num,
                message=f"{heading_text} (#{anchor}) - {reason}",
                severity='warning',
                heading_level=heading_level, heading_text=heading_text, anchor=anchor, reason=reason
            ),
            excluded_pairs,
            True,
        )
    if kind == 'excluded':
        return (None, excluded_pairs, True)
    if not classification['needs_requirement']:
        return (None, [], True)
    # needs_requirement: no issue yet, not skipped; caller will check ref and maybe add issue
    return (None, [], False)


def _run_requirement_refs_to_excluded(
    requirement_files: List[Path],
    excluded_headings_set: Set[Tuple[str, str]],
    file_cache: FileContentCache,
    workspace_root: Path,
    *,
    output: OutputBuilder,
    args,
) -> List[ValidationIssue]:
    """Warn when a requirement references an excluded heading."""
    issues = []
    req_dir_path = workspace_root / DOCS_DIR / REQUIREMENTS_DIR
    try:
        for req_file in requirement_files:
            req_path = req_file if isinstance(req_file, Path) else req_dir_path / req_file
            if not req_path.is_absolute():
                req_path = workspace_root / req_path
            try:
                req_relative = req_path.relative_to(workspace_root)
            except ValueError:
                req_relative = req_path
            for spec_basename, anchor, line_num in get_requirement_spec_anchor_links(
                req_path, file_cache
            ):
                if (spec_basename, anchor) in excluded_headings_set:
                    issues.append(ValidationIssue.create(
                        "requirement_references_excluded_heading", Path(req_relative), line_num,
                        line_num,
                        message=(
                            f"references excluded heading {spec_basename}#{anchor} - "
                            "requirement may be out of date or overly implementation-specific"
                        ),
                        severity='warning', spec_basename=spec_basename, anchor=anchor
                    ))
    except _SCRIPT_ERROR_EXCEPTIONS as e:
        if args.verbose:
            output.add_verbose_line(
                f"  Warning: Error scanning requirements for excluded refs: {e}"
            )
    return issues


# (issue_type, severity) -> (counter_key, 'error'|'warning', count_empty_org)
_ISSUE_OUTPUT_DISPATCH = {
    ('missing_requirement', 'error'): ('missing_headings_count', 'error', False),
    ('missing_requirement', 'warning'): ('warning_headings_count', 'warning', False),
    ('architectural_heading', 'warning'): ('architectural_headings_count', 'warning', False),
    ('organizational_heading', 'warning'): ('organizational_headings_count', 'warning', True),
    ('analysis_error', 'error'): ('analysis_errors_count', None, False),
    ('signature_only_heading', 'warning'): ('signature_only_headings_count', 'warning', False),
    ('example_only_heading', 'warning'): ('example_only_headings_count', 'warning', False),
    ('non_prose_heading', 'warning'): ('non_prose_headings_count', 'warning', False),
    ('requirement_references_excluded_heading', 'warning'): (
        'requirement_refs_excluded_count', 'warning', False
    ),
    ('anchor_in_text_missing_href_anchor', 'error'): (
        'anchor_in_text_missing_href_anchor_count', 'error', False
    ),
    ('too_many_tech_spec_links', 'warning'): (
        'too_many_tech_spec_links_warning_count', 'warning', False
    ),
    ('too_many_tech_spec_links', 'error'): (
        'too_many_tech_spec_links_error_count', 'error', False
    ),
}


def _emit_issues_and_summary(
    issues: List[ValidationIssue],
    total_tech_specs: int,
    missing_count: int,
    total_headings: int,
    *,
    output: OutputBuilder,
    args,
    no_color: bool,
) -> int:
    """Print issues and summary; return exit code."""

    def _init_issue_counters() -> dict:
        return {
            'missing_headings_count': 0,
            'warning_headings_count': 0,
            'architectural_headings_count': 0,
            'organizational_headings_count': 0,
            'empty_org_count': 0,
            'analysis_errors_count': 0,
            'signature_only_headings_count': 0,
            'example_only_headings_count': 0,
            'non_prose_headings_count': 0,
            'requirement_refs_excluded_count': 0,
            'anchor_in_text_missing_href_anchor_count': 0,
            'too_many_tech_spec_links_warning_count': 0,
            'too_many_tech_spec_links_error_count': 0,
        }

    def _accumulate_issues_and_emit_output() -> dict:
        counters = _init_issue_counters()
        for issue in issues:
            key = (issue.issue_type, issue.severity)
            if key not in _ISSUE_OUTPUT_DISPATCH:
                continue
            counter_key, line_kind, count_empty_org = _ISSUE_OUTPUT_DISPATCH[key]
            counters[counter_key] += 1
            if count_empty_org and issue.extra_fields.get('is_empty', False):
                counters['empty_org_count'] += 1
            if line_kind == 'error':
                output.add_error_line(issue.format_message(no_color=no_color))
            elif line_kind == 'warning':
                output.add_warning_line(issue.format_message(no_color=no_color))
        return counters

    def _build_summary(counters: dict) -> List[Tuple[str, int]]:
        minor_org_count = counters['organizational_headings_count'] - counters['empty_org_count']
        covered = total_tech_specs - missing_count
        covered_headings = (
            total_headings
            - counters['missing_headings_count']
            - counters['warning_headings_count']
            - counters['architectural_headings_count']
            - counters['organizational_headings_count']
            - counters['signature_only_headings_count']
            - counters['example_only_headings_count']
            - counters['non_prose_headings_count']
        )
        return [
            ("Total tech specs:", total_tech_specs),
            ("Covered:", covered),
            ("Missing:", missing_count),
            ("Total H2+ headings:", total_headings),
            ("Covered headings:", covered_headings),
            ("Missing headings (errors):", counters['missing_headings_count']),
            ("Architectural headings (warnings):", counters['architectural_headings_count']),
            ("Missing headings (warnings):", counters['warning_headings_count']),
            ("Organizational headings (empty):", counters['empty_org_count']),
            ("Organizational headings (minor content):", minor_org_count),
            ("Signature-only headings (skipped):", counters['signature_only_headings_count']),
            ("Example-only headings (skipped):", counters['example_only_headings_count']),
            ("Non-prose headings (skipped):", counters['non_prose_headings_count']),
            (
                "Requirement refs to excluded (warnings):",
                counters['requirement_refs_excluded_count']
            ),
            (
                "Anchor-in-text missing href anchor (errors):",
                counters['anchor_in_text_missing_href_anchor_count']
            ),
            (
                "Too many spec links per requirement (warnings):",
                counters['too_many_tech_spec_links_warning_count']
            ),
            (
                "Too many spec links per requirement (errors):",
                counters['too_many_tech_spec_links_error_count']
            ),
            ("Analysis errors:", counters['analysis_errors_count']),
        ]

    def _has_warnings_only(counters: dict) -> bool:
        if counters['missing_headings_count']:
            return False
        return (
            counters['warning_headings_count'] > 0
            or counters['architectural_headings_count'] > 0
            or counters['organizational_headings_count'] > 0
            or counters['requirement_refs_excluded_count'] > 0
            or counters['too_many_tech_spec_links_warning_count'] > 0
        )

    def _has_errors(counters: dict) -> bool:
        return (
            missing_count > 0
            or counters['missing_headings_count'] > 0
            or counters['analysis_errors_count'] > 0
            or counters['anchor_in_text_missing_href_anchor_count'] > 0
            or counters['too_many_tech_spec_links_error_count'] > 0
        )

    counters = _accumulate_issues_and_emit_output()

    output.add_summary_header()
    output.add_summary_section(_build_summary(counters))

    if counters['missing_headings_count'] > 0:
        output.add_blank_line("error")
        output.add_line(
            "Note: Functional H2+ headings (those describing testable behavior) "
            "must have at least one requirement referencing them with an anchor link. "
            "Architectural headings are reported as warnings and should also have requirements.",
            section="error"
        )

    if _has_warnings_only(counters) and total_headings > 0:
        output.add_warnings_only_message(
            message="Note: Warnings are shown above, but they do not cause the audit to fail.",
        )
    elif total_headings > 0:
        output.add_success_message("All H2+ headings have requirement coverage!")

    if _has_errors(counters):
        output.add_failure_message("Validation failed. Please fix the errors above.")

    output.print()
    return output.get_exit_code(args.no_fail)


def main():
    """Main function to audit requirements coverage."""
    parser = argparse.ArgumentParser(
        description='Audit requirements coverage for tech specs',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog=__doc__
    )
    parser.add_argument('-v', '--verbose', action='store_true',
                        help='Show detailed progress information')
    parser.add_argument('-p', '--path', type=str,
                        help='Check only the specified file(s) or directory(ies) (comma-separated)')
    parser.add_argument(
        '--output', '-o', type=str, metavar='FILE',
        help='Write detailed output to FILE'
    )
    parser.add_argument(
        '--nocolor', '--no-color', action='store_true',
        help='Disable colored output'
    )
    parser.add_argument('--no-fail', action='store_true',
                        help='Exit with code 0 even if errors are found')

    args = parser.parse_args()
    workspace_root = get_workspace_root()
    tech_specs_dir = workspace_root / DOCS_DIR / TECH_SPECS_DIR
    requirements_dir = workspace_root / DOCS_DIR / REQUIREMENTS_DIR

    if not tech_specs_dir.exists():
        print(f"Error: Tech specs directory not found: {tech_specs_dir}")
        return 1
    if not args.path and not requirements_dir.exists():
        print(f"Error: Requirements directory not found: {requirements_dir}")
        return 1

    target_paths = parse_paths(args.path)
    no_color = args.nocolor or parse_no_color_flag(sys.argv)

    output = OutputBuilder(
        "Requirements Coverage Audit",
        "Audits requirements coverage for tech specs",
        no_color=no_color,
        verbose=args.verbose,
        output_file=args.output
    )
    file_cache = FileContentCache()

    index_ref_errors = _run_index_ref_check(
        requirements_dir, file_cache, args.verbose, target_paths,
        output=output, no_color=no_color
    )
    if index_ref_errors is not None:
        output.print()
        return 1

    tech_specs = find_tech_specs(tech_specs_dir)
    if not tech_specs:
        output.add_line("No tech spec files found.")
        output.print()
        return 0

    missing_specs = _run_file_level_coverage(
        tech_specs, requirements_dir, file_cache, target_paths,
        output=output, args=args, no_color=no_color
    )
    total_tech_specs = len(tech_specs)
    missing_count = len(missing_specs)

    requirement_files = get_requirement_files(requirements_dir, target_paths)
    issues, total_headings, excluded_headings_set = _run_heading_coverage(
        tech_specs, requirement_files, file_cache, workspace_root,
        output=output, args=args, no_color=no_color
    )
    issues.extend(check_anchor_in_text_missing_href_anchor(
        requirement_files,
        file_cache,
        workspace_root=workspace_root,
    ))
    issues.extend(check_requirement_tech_spec_link_thresholds(
        requirement_files,
        file_cache,
        workspace_root=workspace_root,
    ))
    more_issues = _run_requirement_refs_to_excluded(
        requirement_files, excluded_headings_set, file_cache, workspace_root,
        output=output, args=args
    )
    issues.extend(more_issues)

    return _emit_issues_and_summary(
        issues, total_tech_specs, missing_count, total_headings,
        output=output, args=args, no_color=no_color
    )


if __name__ == '__main__':
    sys.exit(main())
