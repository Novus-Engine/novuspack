#!/usr/bin/env python3
"""
Audit feature file coverage for requirements.

This script checks which requirements have feature files that reference them
using @REQ-* tags, and validates that @spec() references in feature files
match the specs referenced by their requirements.

The script:
1. Extracts all requirements from requirements files (excluding obsolete and
   documentation-only requirements)
2. Extracts spec references from each requirement (from markdown links)
3. Checks that each requirement has at least one feature file with @REQ-* tag
4. Validates that @spec() tags in feature files match the specs referenced
   by the requirements (by filename, allowing for anchor differences)

Usage:
    python3 scripts/audit_feature_coverage.py [options]

Options:
    --verbose, -v       Show detailed progress information
    --path, -p PATHS    Check only the specified file(s) or
                        directory(ies) (recursive). Can be a single
                        path or comma-separated list of paths
    --help, -h          Show this help message

Examples:
    # Basic audit
    python3 scripts/audit_feature_coverage.py

    # Verbose output
    python3 scripts/audit_feature_coverage.py --verbose

    # Check specific file
    python3 scripts/audit_feature_coverage.py --path \\
        features/file_management/write.feature

    # Check specific directory
    python3 scripts/audit_feature_coverage.py --path \\
        features/file_management

    # Check multiple paths
    python3 scripts/audit_feature_coverage.py --path \\
        features/file_management,features/compression
"""

import sys
import argparse
import re
from dataclasses import dataclass
from pathlib import Path
from typing import List, Set, Tuple, Optional, Dict

# Import shared utilities.
#
# Note: when executed as `python3 scripts/audit_feature_coverage.py`, Python adds the
# script's directory (`scripts/`) to sys.path automatically, so `lib.*` imports work.
from lib._validation_utils import (
    OutputBuilder, parse_no_color_flag,
    get_workspace_root, parse_paths,
    ValidationIssue, find_feature_files, find_markdown_files,
    DOCS_DIR, REQUIREMENTS_DIR, FEATURES_DIR
)


# Compiled regex patterns for performance (module level)
_RE_SPEC_LINK_PATTERN = re.compile(r'\[([^\]]+)\]\(\.\./tech_specs/([^\)]+)\)')
_RE_REQ_PATTERN = re.compile(r'-\s*(?:~~)?\s*(REQ-[A-Z_]+-\d+):')
_RE_TYPE_PATTERN = re.compile(r'\[type:\s*(obsolete|documentation-only)\]')
_RE_SPEC_TAG_PATTERN = re.compile(r'@spec\(([^\)]+)\)')
_RE_REQ_TAG_PATTERN = re.compile(r'@(REQ-[A-Z_]+-\d+)')

_CROSS_CUTTING_SPECS = {'api_core.md', 'file_validation.md'}


@dataclass(frozen=True)
class _FeatureIndex:
    """
    Precomputed index of feature file tags to avoid O(requirements * features) scans.

    - req_occurrence_counts: total occurrences of @REQ-* tags, summed across all feature files.
    - req_to_feature_files: maps req_id -> list of feature files that reference it at least once.
    - feature_spec_refs: maps feature file -> set of @spec() tag values in that file.
    """

    req_occurrence_counts: Dict[str, int]
    req_to_feature_files: Dict[str, List[Path]]
    feature_spec_refs: Dict[Path, Set[str]]


def extract_spec_references_from_line(line: str) -> Set[str]:
    """
    Extract spec references from a requirements line.

    Args:
        line: The line text to parse

    Returns:
        Set of spec references in format "filename#anchor" or "filename"
    """
    spec_refs = set()
    # Pattern to match markdown links: [text](../tech_specs/filename#anchor)
    # or [text](../tech_specs/filename)
    # Link text may contain anchor: [filename#anchor](../tech_specs/filename)

    for match in _RE_SPEC_LINK_PATTERN.finditer(line):
        # e.g., "api_file_management.md#614-updatefile-behavior"
        link_text = match.group(1)
        # e.g., "api_file_mgmt_file_entry.md" or "api_file_mgmt_file_entry.md#anchor"
        spec_path = match.group(2)

        # If spec_path already has an anchor, use it as-is
        if '#' in spec_path:
            spec_refs.add(spec_path)
        else:
            # Extract anchor from link text if available
            if '#' in link_text:
                anchor_part = link_text.split('#', 1)[1]  # Get everything after #
                full_ref = f"{spec_path}#{anchor_part}"
                spec_refs.add(full_ref)
            else:
                # No anchor in either link text or path
                spec_refs.add(spec_path)

    return spec_refs


def extract_requirements_from_file(
    req_file: Path, verbose: bool = False
) -> List[Tuple[str, int, Optional[str], Set[str]]]:
    """
    Extract all requirement IDs and their spec references from a requirements file.

    Args:
        req_file: Path to the requirements file
        verbose: Whether to show detailed progress

    Returns:
        List of tuples (req_id, line_number, req_type, spec_refs) where:
        - req_id is the requirement ID (e.g., 'REQ-API_BASIC-028')
        - line_number is the line where it appears
        - req_type is 'obsolete', 'documentation-only', or None
        - spec_refs is a set of spec references (format: "filename#anchor" or "filename")
    """
    # Match requirements with or without strikethrough
    # Format: "- REQ-XXX-123:" or "- ~~REQ-XXX-123:"
    # Pattern to extract type tag: [type: obsolete] or [type: documentation-only]
    req_definitions = []

    try:
        with open(req_file, 'r', encoding='utf-8') as f:
            for line_num, line in enumerate(f, start=1):
                match = _RE_REQ_PATTERN.search(line)
                if match:
                    req_id = match.group(1)
                    # Check for type tag
                    type_match = _RE_TYPE_PATTERN.search(line)
                    req_type = None
                    if type_match:
                        req_type = type_match.group(1)
                    # Extract spec references from the line
                    spec_refs = extract_spec_references_from_line(line)
                    req_definitions.append((req_id, line_num, req_type, spec_refs))

            if verbose:
                print(
                    f"  Found {len(req_definitions)} requirement(s) "
                    f"in {req_file.name}"
                )

    except (IOError, OSError) as e:
        if verbose:
            print(f"  Warning: Could not read {req_file}: {e}")
    except UnicodeDecodeError as e:
        if verbose:
            print(f"  Warning: Could not decode {req_file} (encoding issue): {e}")

    return req_definitions


def _compute_missing_spec_refs_for_feature(
    req_spec_refs: Set[str],
    feature_spec_refs: Set[str],
) -> List[str]:
    """
    Compute which requirement spec refs are missing from a feature file's @spec() tags.

    Matching is by filename only (anchors may differ due to heading renumbering).
    Cross-cutting specs are ignored.
    """
    missing: List[str] = []
    if not req_spec_refs:
        return missing

    feature_filenames = {spec.split('#')[0] for spec in feature_spec_refs}
    for req_spec in req_spec_refs:
        req_filename = req_spec.split('#')[0]
        if req_filename in _CROSS_CUTTING_SPECS:
            continue
        if req_filename not in feature_filenames:
            missing.append(req_spec)

    return sorted(missing)


def build_feature_index(
    features_dir: Path,
    *,
    verbose: bool = False,
    target_paths: Optional[List[str]] = None,
) -> _FeatureIndex:
    """
    Build an index of @REQ-* and @spec() tags across feature files.

    This replaces repeated per-requirement scans and prevents apparent "hangs"
    caused by scanning all feature files once per requirement.
    """
    feature_files = find_feature_files(
        target_paths=target_paths,
        default_dir=features_dir,
        verbose=verbose
    )
    req_occurrence_counts: Dict[str, int] = {}
    req_to_feature_files: Dict[str, List[Path]] = {}
    feature_spec_refs: Dict[Path, Set[str]] = {}

    for feature_file in feature_files:
        try:
            with open(feature_file, 'r', encoding='utf-8') as f:
                content = f.read()
        except (IOError, OSError) as e:
            if verbose:
                print(f"  Warning: Could not read {feature_file}: {e}", file=sys.stderr)
            continue
        except UnicodeDecodeError as e:
            if verbose:
                print(
                    f"  Warning: Could not decode {feature_file} (encoding issue): {e}",
                    file=sys.stderr
                )
            continue

        specs = set(_RE_SPEC_TAG_PATTERN.findall(content))
        feature_spec_refs[feature_file] = specs

        req_matches = _RE_REQ_TAG_PATTERN.findall(content)
        if not req_matches:
            continue

        per_file_counts: Dict[str, int] = {}
        for req_id in req_matches:
            per_file_counts[req_id] = per_file_counts.get(req_id, 0) + 1

        for req_id, occurrences in per_file_counts.items():
            req_occurrence_counts[req_id] = req_occurrence_counts.get(req_id, 0) + occurrences
            req_to_feature_files.setdefault(req_id, []).append(feature_file)

    return _FeatureIndex(
        req_occurrence_counts=req_occurrence_counts,
        req_to_feature_files=req_to_feature_files,
        feature_spec_refs=feature_spec_refs,
    )


def find_all_requirements(
    requirements_dir: Path, verbose: bool = False
) -> Dict[str, Tuple[Path, int, Set[str]]]:
    """
    Find all requirements from all requirements files, excluding obsolete
    and documentation-only requirements.

    Args:
        requirements_dir: Path to the requirements directory
        verbose: Whether to show detailed progress

    Returns:
        Dict mapping req_id -> (req_file, line_number, spec_refs) tuple where:
        - req_file is the requirements file path
        - line_number is the line where the requirement appears
        - spec_refs is a set of spec references (format: "filename#anchor" or "filename")
    """
    all_requirements = {}

    # Exclude README.md and traceability.md
    exclude_files = ['README.md', 'traceability.md']

    # Use shared utility function to find markdown files
    req_files = find_markdown_files(
        default_dir=requirements_dir,
        verbose=verbose
    )

    for req_file in req_files:
        if req_file.name in exclude_files:
            continue

        if verbose:
            print(f"  Scanning {req_file.name}...")

        req_definitions = extract_requirements_from_file(req_file, verbose)

        for req_id, line_num, req_type, spec_refs in req_definitions:
            # Skip obsolete and documentation-only requirements
            if req_type in ('obsolete', 'documentation-only'):
                if verbose:
                    print(f"    Skipping {req_id} (type: {req_type})")
                continue

            if req_id in all_requirements:
                if verbose:
                    print(
                        f"    Warning: Duplicate requirement {req_id} "
                        f"(already found in {all_requirements[req_id][0].name})"
                    )
            else:
                all_requirements[req_id] = (req_file, line_num, spec_refs)

    return all_requirements


def extract_spec_tags_from_feature(feature_file: Path) -> Set[str]:
    """
    Extract all @spec() tags from a feature file.

    Args:
        feature_file: Path to the feature file

    Returns:
        Set of spec references (format: "filename#anchor" or "filename")
    """
    spec_refs = set()
    # Pattern to match @spec(filename#anchor) or @spec(filename)

    try:
        with open(feature_file, 'r', encoding='utf-8') as f:
            content = f.read()
            for match in _RE_SPEC_TAG_PATTERN.finditer(content):
                spec_ref = match.group(1)
                spec_refs.add(spec_ref)
    except (IOError, OSError, UnicodeDecodeError):
        # File read/encoding errors - return empty set, caller will handle
        return spec_refs

    return spec_refs


def check_feature_references_for_requirement(
    req_id, req_spec_refs, features_dir, verbose=False, target_paths=None
):
    """
    Check how many features reference a requirement using @REQ-* tags,
    and validate that @spec() tags match the requirement's spec references.

    Args:
        req_id: The requirement ID (e.g., 'REQ-API_BASIC-028')
        req_spec_refs: Set of spec references from the requirement
        features_dir: Path to the features directory
        verbose: Whether to show detailed progress
        target_paths: Optional list of specific files or directories to check

    Returns:
        Tuple (count, spec_mismatches) where:
        - count is the number of feature files referencing the requirement
        - spec_mismatches is a list of tuples:
          (feature_file, feature_specs, all_req_specs, missing_specs)
    """
    # Backward-compatible wrapper; kept for external callers.
    # Prefer `build_feature_index()` + `_check_requirement_against_index()` for performance.
    index = build_feature_index(features_dir, verbose=verbose, target_paths=target_paths)
    return _check_requirement_against_index(req_id, req_spec_refs, index, verbose=verbose)


def _check_requirement_against_index(
    req_id: str,
    req_spec_refs: Set[str],
    index: _FeatureIndex,
    *,
    verbose: bool = False,
) -> Tuple[int, List[Tuple[Path, Set[str], Set[str], List[str]]]]:
    """
    Check a requirement's feature coverage and @spec() completeness using a prebuilt index.
    """
    if verbose:
        print(f"  Searching for: @{req_id}")

    count = index.req_occurrence_counts.get(req_id, 0)
    spec_mismatches: List[Tuple[Path, Set[str], Set[str], List[str]]] = []

    for feature_file in index.req_to_feature_files.get(req_id, []):
        feature_specs = index.feature_spec_refs.get(feature_file, set())
        missing_specs = _compute_missing_spec_refs_for_feature(req_spec_refs, feature_specs)
        if missing_specs:
            spec_mismatches.append((
                feature_file,
                feature_specs,
                req_spec_refs,
                missing_specs,
            ))
        if verbose and count > 0:
            # Preserve the old behavior of reporting "Found N in <file>" lines,
            # but without recomputing per-file occurrences (the total count is what matters).
            print(f"    Found in {feature_file}")

    return count, spec_mismatches


def _parse_args():
    parser = argparse.ArgumentParser(
        description='Audit feature file coverage for requirements',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog=__doc__
    )
    parser.add_argument(
        '-v', '--verbose',
        action='store_true',
        help='Show detailed progress information'
    )
    parser.add_argument(
        '-p', '--path',
        type=str,
        help='Check only the specified file(s) or directory(ies) (comma-separated list)'
    )
    parser.add_argument(
        '--output', '-o',
        type=str,
        metavar='FILE',
        help='Write detailed output to FILE'
    )
    parser.add_argument(
        '--nocolor', '--no-color',
        action='store_true',
        help='Disable colored output'
    )
    parser.add_argument(
        '--no-fail',
        action='store_true',
        help='Exit with code 0 even if errors are found'
    )
    return parser.parse_args()


def _resolve_paths(args):
    workspace_root = get_workspace_root()
    requirements_dir = workspace_root / DOCS_DIR / REQUIREMENTS_DIR
    features_dir = workspace_root / FEATURES_DIR
    target_paths = parse_paths(args.path)
    no_color = args.nocolor or parse_no_color_flag(sys.argv)
    return requirements_dir, features_dir, target_paths, no_color


def _validate_dirs(requirements_dir: Path, features_dir: Path, args) -> int:
    if not requirements_dir.exists():
        print(f"Error: Requirements directory not found: {requirements_dir}")
        return 1
    if not args.path and not features_dir.exists():
        print(f"Error: Features directory not found: {features_dir}")
        return 1
    return 0


def _collect_coverage_issues(
    all_requirements: Dict[str, Tuple[Path, int, Set[str]]],
    feature_index: _FeatureIndex,
    *,
    verbose: bool,
    output: OutputBuilder,
) -> List[ValidationIssue]:
    issues: List[ValidationIssue] = []

    for req_id, (req_file, line_num, req_spec_refs) in sorted(all_requirements.items()):
        if verbose:
            output.add_verbose_line(f"Checking {req_id}...")

        count, spec_mismatches = _check_requirement_against_index(
            req_id,
            req_spec_refs,
            feature_index,
            verbose=verbose
        )

        if not count:
            issues.append(ValidationIssue.create(
                "missing_feature_coverage",
                req_file,
                line_num,
                line_num,
                message=f"Requirement {req_id} has no feature files referencing it",
                severity='error',
                req_id=req_id
            ))
        elif verbose:
            output.add_verbose_line(f"✓  {req_id}: {count} feature(s)")

        for (feature_file, feature_specs, expected_specs, missing_specs) in spec_mismatches:
            message = (
                f"Feature file {feature_file.name} is missing @spec() tags for "
                f"requirement {req_id}. "
                f"Required: {', '.join(sorted(expected_specs))}. "
                f"Missing: {', '.join(sorted(missing_specs))}"
            )
            issues.append(ValidationIssue.create(
                "spec_mismatch",
                req_file,
                line_num,
                line_num,
                message=message,
                severity='error',
                req_id=req_id,
                feature_file=str(feature_file),
                feature_specs=list(feature_specs),
                expected_specs=list(expected_specs),
                missing_specs=list(missing_specs)
            ))

    return issues


def _split_issues(issues: List[ValidationIssue]):
    missing_issues = []
    spec_mismatch_issues = []
    for issue in issues:
        if issue.matches(issue_type='missing_feature_coverage'):
            missing_issues.append(issue)
        if issue.matches(issue_type='spec_mismatch'):
            spec_mismatch_issues.append(issue)
    return missing_issues, spec_mismatch_issues


def _emit_summary(
    output: OutputBuilder,
    *,
    total: int,
    covered: int,
    missing: int,
    spec_mismatches: int,
):
    summary_items = [
        ("Total requirements:", total),
        ("Covered:", covered),
        ("Missing:", missing),
    ]
    if spec_mismatches:
        summary_items.append(("Spec mismatches:", spec_mismatches))
    output.add_summary_header()
    output.add_summary_section(summary_items)


def _emit_errors(
    output: OutputBuilder,
    *,
    missing_issues: List[ValidationIssue],
    spec_mismatch_issues: List[ValidationIssue],
    no_color: bool,
):
    if missing_issues:
        output.add_errors_header()
        for issue in missing_issues:
            output.add_error_line(issue.format_message(no_color=no_color))
        output.add_blank_line("error")

    if not spec_mismatch_issues:
        return

    if not missing_issues:
        output.add_errors_header()
    output.add_line(
        "Feature files missing @spec() tags for requirement spec references:",
        section="error"
    )
    output.add_blank_line("error")
    for issue in spec_mismatch_issues:
        req_id = issue.extra_fields.get('req_id', '')
        feature_file = issue.extra_fields.get('feature_file', '')
        feature_specs = issue.extra_fields.get('feature_specs', [])
        expected_specs = issue.extra_fields.get('expected_specs', [])
        missing_specs = issue.extra_fields.get('missing_specs', [])

        output.add_line(
            f"Requirement {req_id} (in {issue.file}:{issue.start_line})",
            section="error"
        )
        output.add_line(
            f"  Feature file: {feature_file}",
            section="error"
        )
        if feature_specs:
            output.add_line(
                f"  Has @spec(): {', '.join(sorted(feature_specs))}",
                section="error"
            )
        else:
            output.add_line(
                "  Has @spec(): (none)",
                section="error"
            )
        output.add_line(
            f"  Required (from requirement): {', '.join(sorted(expected_specs))}",
            section="error"
        )
        output.add_line(
            f"  Missing @spec() tags: {', '.join(sorted(missing_specs))}",
            section="error"
        )
        output.add_blank_line("error")


def main():
    """Main function to audit feature coverage."""
    args = _parse_args()
    requirements_dir, features_dir, target_paths, no_color = _resolve_paths(args)
    exit_code = _validate_dirs(requirements_dir, features_dir, args)
    if exit_code:
        return exit_code

    # Create output builder (header streams immediately if verbose)
    output = OutputBuilder(
        "Feature File Coverage Audit",
        "Audits feature file coverage for requirements",
        no_color=no_color,
        verbose=args.verbose,
        output_file=args.output
    )

    # Get all requirements (excluding obsolete and documentation-only)
    if args.verbose:
        output.add_verbose_line("Extracting requirements...")
        output.add_blank_line("working_verbose")

    all_requirements = find_all_requirements(requirements_dir, args.verbose)

    if not all_requirements:
        output.add_line("No requirements found.")
        output.print()
        return 0

    if args.verbose:
        output.add_verbose_line(f"Found {len(all_requirements)} requirement(s)")
        output.add_verbose_line("Checking Coverage...")
        output.add_blank_line("working_verbose")

    if args.verbose:
        output.add_verbose_line("Indexing feature files...")
        output.add_blank_line("working_verbose")

    feature_index = build_feature_index(
        features_dir,
        verbose=args.verbose,
        target_paths=target_paths,
    )
    issues = _collect_coverage_issues(
        all_requirements,
        feature_index,
        verbose=args.verbose,
        output=output,
    )
    missing_issues, spec_mismatch_issues = _split_issues(issues)
    errors = [issue for issue in issues if issue.matches(severity='error')]

    total = len(all_requirements)
    missing_count = len(missing_issues)
    covered = total - missing_count
    _emit_summary(
        output,
        total=total,
        covered=covered,
        missing=missing_count,
        spec_mismatches=len(spec_mismatch_issues),
    )
    _emit_errors(
        output,
        missing_issues=missing_issues,
        spec_mismatch_issues=spec_mismatch_issues,
        no_color=no_color,
    )

    if errors:
        output.add_line(
            "Note: After adding feature references, verify that each "
            "reference points to the correct requirement and that @spec() "
            "tags match the specs referenced by the requirements.",
            section="error"
        )
        output.add_failure_message("Validation failed. Please fix the errors above.")
    else:
        output.add_success_message("All requirements have feature coverage!")

    output.print()
    return output.get_exit_code(args.no_fail)


if __name__ == '__main__':
    sys.exit(main())
