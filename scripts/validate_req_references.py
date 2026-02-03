#!/usr/bin/env python3
"""
Validate requirement references in feature files.

This script checks that all @REQ-* tags in feature files reference valid
requirements defined in the requirements documentation.

Usage:
    python3 scripts/validate_req_references.py [options]

Options:
    --verbose, -v       Show detailed progress information
    --path, -p PATHS    Check only the specified file(s) or
                        directory(ies) (recursive). Can be a single
                        path or comma-separated list of paths
    --help, -h          Show this help message

Examples:
    # Basic validation
    python3 scripts/validate_req_references.py

    # Verbose output
    python3 scripts/validate_req_references.py --verbose

    # Check specific file
    python3 scripts/validate_req_references.py --path \\
        features/file_management/write.feature

    # Check specific directory
    python3 scripts/validate_req_references.py --path \\
        features/file_management

    # Check multiple paths
    python3 scripts/validate_req_references.py --path \\
        features/file_management,features/compression
"""

import sys
import re
import argparse
from pathlib import Path
from collections import defaultdict
from typing import Dict, List, Optional, Tuple

from lib._validation_utils import (
    OutputBuilder, parse_no_color_flag,
    is_in_dot_directory, get_workspace_root, parse_paths,
    ValidationIssue, DOCS_DIR, REQUIREMENTS_DIR, FEATURES_DIR
)


# Compiled regex patterns for performance (module level)
_RE_REQ_TAG_PATTERN = re.compile(r'@(REQ-[A-Z_]+-\d+)')
_RE_REQ_DEFINITION_PATTERN = re.compile(r'-\s*(?:~~)?\s*(REQ-[A-Z_]+-\d+):')
_RE_TYPE_TAG_PATTERN = re.compile(r'\[type:\s*(obsolete|documentation-only)\]')
_RE_REQ_FORMAT_PATTERN = re.compile(r'^REQ-[A-Z_]+-[0-9]+$')


# Mapping from canonical REQ category prefixes to requirement file names.
#
# NOTE: Keep this list strict.
# Similar/legacy prefixes should not be accepted as valid categories.
CATEGORY_TO_FILE = {
    'API_BASIC': 'basic_ops.md',
    'COMPR': 'compression.md',
    'CORE': 'core.md',
    'CRYPTO': 'security_encryption.md',
    'DEDUP': 'dedup.md',
    'FILEFMT': 'file_format.md',
    'FILEMGMT': 'file_mgmt.md',
    'FILETYPES': 'file_types.md',
    'GEN': 'generics.md',
    'META': 'metadata.md',
    'METASYS': 'metadata_system.md',
    'PIPELINE': 'transformation_pipeline.md',
    'SEC': 'security.md',
    'SIG': 'signatures.md',
    'STREAM': 'streaming.md',
    'TEST': 'testing.md',
    'VALID': 'validation.md',
    'WRITE': 'writing.md',
}

# Legacy or lookalike category prefixes that must not be used.
# If these appear in feature tags or requirements docs, we emit a clear error
# suggesting the canonical prefix to use instead.
DEPRECATED_CATEGORY_PREFIXES = {
    'FILEFORMAT': 'FILEFMT',
    'FILETYPE': 'FILETYPES',
    'GENERIC': 'GEN',
    'SECURITY': 'SEC',
    'SEC_ENC': 'CRYPTO',
    # TRACE is not a requirements domain file; traceability.md is a matrix.
    'TRACE': None,
}


def extract_req_tags_from_feature(feature_file, verbose=False):
    """
    Extract all @REQ-* tags from a feature file.

    Args:
        feature_file: Path to the feature file
        verbose: Whether to show detailed progress

    Returns:
        Set of requirement IDs (e.g., 'REQ-API_BASIC-028')
    """
    req_tags = set()

    try:
        with open(feature_file, 'r', encoding='utf-8') as f:
            content = f.read()
            matches = _RE_REQ_TAG_PATTERN.findall(content)
            req_tags.update(matches)

            if verbose and matches:
                print(f"  Found {len(matches)} REQ tag(s) in {feature_file.name}")

    except (IOError, OSError) as e:
        # File read errors - log warning if verbose
        if verbose:
            print(f"  Warning: Could not read {feature_file}: {e}", file=sys.stderr)
    except UnicodeDecodeError as e:
        # Encoding errors - log warning if verbose
        if verbose:
            print(
                f"  Warning: Could not decode {feature_file} (encoding issue): {e}",
                file=sys.stderr
            )
    except (ValueError, KeyError, TypeError, AttributeError, RuntimeError) as e:
        # Unexpected errors - log warning if verbose
        if verbose:
            print(f"  Warning: Unexpected error reading {feature_file}: {e}", file=sys.stderr)

    return req_tags


def extract_req_definitions_from_requirements(req_file, verbose=False):
    """
    Extract all REQ-* definitions from a requirements file.

    Args:
        req_file: Path to the requirements file
        verbose: Whether to show detailed progress

    Returns:
        List of tuples (req_id, line_number, req_type) where:
        - req_id is the requirement ID (e.g., 'REQ-API_BASIC-028')
        - line_number is the line where it appears
        - req_type is 'obsolete', 'documentation-only', or None
    """
    # Match requirements with or without strikethrough
    # Format: "- REQ-XXX-123:" or "- ~~REQ-XXX-123:"
    # Pattern to extract type tag: [type: obsolete] or [type: documentation-only]
    req_definitions = []

    try:
        with open(req_file, 'r', encoding='utf-8') as f:
            for line_num, line in enumerate(f, start=1):
                match = _RE_REQ_DEFINITION_PATTERN.search(line)
                if match:
                    req_id = match.group(1)
                    # Check for type tag
                    type_match = _RE_TYPE_TAG_PATTERN.search(line)
                    req_type = None
                    if type_match:
                        req_type = type_match.group(1)
                    req_definitions.append((req_id, line_num, req_type))

            if verbose:
                print(
                    f"  Found {len(req_definitions)} requirement(s) "
                    f"in {req_file.name}"
                )

    except (ValueError, KeyError, TypeError, AttributeError, RuntimeError) as e:
        if verbose:
            print(f"  Warning: Could not read {req_file}: {e}")

    return req_definitions


def get_category_from_req_id(req_id):
    """
    Extract the category from a requirement ID.

    Args:
        req_id: Requirement ID (e.g., 'REQ-API_BASIC-028')

    Returns:
        Category string (e.g., 'API_BASIC') or None if invalid
    """
    match = re.match(r'REQ-([A-Z_]+)-\d+', req_id)
    if match:
        return match.group(1)
    return None


def get_expected_req_file(category):
    """
    Get the expected requirements file for a category.

    Args:
        category: Category string (e.g., 'API_BASIC')

    Returns:
        Expected requirements filename or None if unknown
    """
    return CATEGORY_TO_FILE.get(category)


def validate_req_format(req_id):
    """
    Validate that a requirement ID matches the expected format.

    Args:
        req_id: Requirement ID to validate

    Returns:
        True if format is valid, False otherwise
    """
    return bool(_RE_REQ_FORMAT_PATTERN.match(req_id))


def validate_req_format_errors(requirements_data, feature_tags, workspace_root, _features_dir):
    """
    Validate format consistency of all requirement IDs.

    Args:
        requirements_data: Dict mapping req_file -> list of (req_id, line_num) tuples
        feature_tags: Dict mapping req_id -> list of feature files
        workspace_root: Path to workspace root (for relative paths)
        features_dir: Path to features directory (for relative paths)

    Returns:
        List of format error dicts with keys: type, file, line, req_id, reason
    """
    format_errors = []

    # Helper to get relative path
    def get_relative_path(file_path, base_dir):
        try:
            return str(file_path.relative_to(base_dir))
        except ValueError:
            return str(file_path)

    # Check requirements from requirements files
    for req_file, req_list in requirements_data.items():
        for req_id, line_num, _ in req_list:
            if not _RE_REQ_FORMAT_PATTERN.match(req_id):
                format_errors.append(ValidationIssue.create(
                    'format',
                    Path(get_relative_path(req_file, workspace_root)),
                    line_num,
                    line_num,
                    message=f"{req_id}: Invalid format: does not match REQ-[A-Z_]+-[0-9]+",
                    severity='error',
                    req_id=req_id,
                    reason='Invalid format: does not match REQ-[A-Z_]+-[0-9]+'
                ))

    # Check feature file tags
    for req_id, feature_files in feature_tags.items():
        if not _RE_REQ_FORMAT_PATTERN.match(req_id):
            for feature_file in feature_files:
                format_errors.append(ValidationIssue.create(
                    'format',
                    Path(get_relative_path(feature_file, workspace_root)),
                    1,
                    1,
                    message=f"{req_id}: Invalid format: does not match REQ-[A-Z_]+-[0-9]+",
                    severity='error',
                    req_id=req_id,
                    reason='Invalid format: does not match REQ-[A-Z_]+-[0-9]+'
                ))

    return format_errors


def check_duplicate_requirements(requirements_data, workspace_root):
    """
    Check for duplicate requirement IDs within each requirements file.

    Args:
        requirements_data: Dict mapping req_file -> list of (req_id, line_num) tuples
        workspace_root: Path to workspace root (for relative paths)

    Returns:
        List of duplicate error dicts with keys: type, file, req_id, lines
    """
    duplicate_errors = []

    # Helper to get relative path
    def get_relative_path(file_path, base_dir):
        try:
            return str(file_path.relative_to(base_dir))
        except ValueError:
            return str(file_path)

    for req_file, req_list in requirements_data.items():
        # Track occurrences of each req_id
        req_occurrences = defaultdict(list)
        for req_id, line_num, _ in req_list:
            req_occurrences[req_id].append(line_num)

        # Find duplicates
        for req_id, line_nums in req_occurrences.items():
            if len(line_nums) > 1:
                duplicate_errors.append({
                    'type': 'duplicate',
                    'file': get_relative_path(req_file, workspace_root),
                    'req_id': req_id,
                    'lines': sorted(line_nums)
                })

    return duplicate_errors


def check_sequential_numbering(requirements_data, workspace_root):
    """
    Check for sequential numbering gaps within each category prefix per file.

    Args:
        requirements_data: Dict mapping req_file -> list of (req_id, line_num) tuples
        workspace_root: Path to workspace root (for relative paths)

    Returns:
        List of sequential numbering warning dicts with keys: type, file, category, missing_numbers
    """
    sequential_warnings = []

    # Helper to get relative path
    def get_relative_path(file_path, base_dir):
        try:
            return str(file_path.relative_to(base_dir))
        except ValueError:
            return str(file_path)

    for req_file, req_list in requirements_data.items():
        # Group by category prefix
        category_reqs = defaultdict(list)
        for req_id, line_num, _req_type in req_list:
            category = get_category_from_req_id(req_id)
            if category:
                # Extract numeric suffix
                match = re.search(r'-(\d+)$', req_id)
                if match:
                    num = int(match.group(1))
                    category_reqs[category].append((num, req_id, line_num))

        # Check each category for sequential numbering
        for category, req_nums in category_reqs.items():
            if not req_nums:
                continue

            # Extract just the numbers and sort
            numbers = sorted([num for num, _, _ in req_nums])
            if not numbers:
                continue

            # Find gaps in sequence
            min_num = numbers[0]
            max_num = numbers[-1]
            expected_numbers = set(range(min_num, max_num + 1))
            actual_numbers = set(numbers)
            missing_numbers = sorted(expected_numbers - actual_numbers)

            if missing_numbers:
                missing_str = ', '.join(str(num) for num in missing_numbers)
                sequential_warnings.append(ValidationIssue.create(
                    'sequential',
                    Path(get_relative_path(req_file, workspace_root)),
                    1,
                    1,
                    message=f"REQ-{category}-* (Missing numbers: {missing_str})",
                    severity='warning',
                    category=category,
                    missing_numbers=missing_numbers
                ))

    return sequential_warnings


def _collect_feature_files_for_stub_check(features_dir, target_paths, verbose):
    """Return list of feature file paths to check for stubs."""
    if target_paths:
        out = []
        for target_path in target_paths:
            target = Path(target_path)
            if not target.exists():
                if verbose:
                    print(f"Warning: Target path does not exist: {target_path}")
                continue
            if target.is_file():
                if target.suffix == '.feature':
                    out.append(target)
                elif verbose:
                    print(f"Warning: Target file is not a .feature file: {target_path}")
            else:
                out.extend(f for f in target.rglob('*.feature') if not is_in_dot_directory(f))
        return out
    return [f for f in features_dir.rglob('*.feature') if not is_in_dot_directory(f)]


def _count_content_lines(path: Path) -> int:
    """Return number of non-empty, non-comment lines in a feature file."""
    count = 0
    with open(path, 'r', encoding='utf-8') as f:
        for line in f:
            stripped = line.strip()
            if stripped and not stripped.startswith('#'):
                count += 1
    return count


def check_feature_stubs(features_dir, verbose=False, target_paths=None):
    """
    Check for feature files that are stubs (8 or fewer content lines).

    A content line is a line that is not empty and not a comment.
    Comments in Gherkin files start with '#'.

    Args:
        features_dir: Path to the features directory
        verbose: Whether to show detailed progress
        target_paths: Optional list of specific files or directories to check

    Returns:
        List of feature files that are stubs (each entry is a dict with
        'file' and 'line_count' keys)
    """
    feature_files = _collect_feature_files_for_stub_check(
        features_dir, target_paths, verbose
    )
    stub_files = []
    for feature_file in feature_files:
        try:
            content_line_count = _count_content_lines(feature_file)
            if content_line_count <= 8:
                stub_files.append(ValidationIssue.create(
                    'stub_file', feature_file, 1, 1,
                    message=f"Feature file is a stub ({content_line_count} content lines)",
                    severity='warning', line_count=content_line_count,
                ))
                if verbose:
                    print(f"  Found stub: {feature_file.name} ({content_line_count} content lines)")
        except (IOError, OSError) as e:
            if verbose:
                print(f"  Warning: Could not read {feature_file}: {e}", file=sys.stderr)
        except UnicodeDecodeError as e:
            if verbose:
                print(
                    f"  Warning: Could not decode {feature_file} (encoding issue): {e}",
                    file=sys.stderr,
                )
        except (ValueError, KeyError, TypeError, AttributeError, RuntimeError) as e:
            if verbose:
                print(f"  Warning: Unexpected error reading {feature_file}: {e}", file=sys.stderr)
    if verbose:
        print()
    return stub_files


def _load_requirement_definitions(
    requirements_dir: Path,
    output: OutputBuilder,
    verbose: bool,
) -> Tuple[Dict[str, str], Dict[Path, List[Tuple[str, int, int]]]]:
    all_requirements: Dict[str, str] = {}
    req_files: Dict[str, str] = {}
    requirements_data: Dict[Path, List[Tuple[str, int, int]]] = {}

    output.add_verbose_line("Loading requirement definitions...")
    if verbose:
        output.add_blank_line("working_verbose")

    for req_file in sorted(requirements_dir.glob('*.md')):
        if is_in_dot_directory(req_file):
            continue
        if req_file.name == 'README.md':
            continue

        req_definitions = extract_req_definitions_from_requirements(req_file, verbose)
        requirements_data[req_file] = req_definitions
        for req_id, _, _ in req_definitions:
            all_requirements[req_id] = req_file.name
            category = get_category_from_req_id(req_id)
            if category:
                req_files[category] = req_file.name

    if verbose:
        output.add_blank_line("working_verbose")
    output.add_verbose_line(
        f"Loaded {len(all_requirements)} requirement definitions "
        f"from {len(req_files)} files"
    )
    output.add_blank_line("working_verbose")
    return all_requirements, requirements_data


def _collect_feature_files(
    features_dir: Path,
    target_paths: Optional[List[str]],
    output: OutputBuilder,
) -> List[Path]:
    feature_files: List[Path] = []
    if target_paths:
        for target_path in target_paths:
            target = Path(target_path)
            if not target.exists():
                print(f"Warning: Target path does not exist: {target_path}")
                continue

            if target.is_file():
                if target.suffix == '.feature':
                    feature_files.append(target)
                else:
                    output.add_warning_line(
                        f"Target file is not a .feature file: {target_path}"
                    )
            else:
                feature_files.extend([
                    f for f in target.rglob('*.feature')
                    if not is_in_dot_directory(f)
                ])
    else:
        feature_files = [
            f for f in features_dir.rglob('*.feature')
            if not is_in_dot_directory(f)
        ]
    return feature_files


def _collect_req_tags(
    feature_files: List[Path],
    verbose: bool,
) -> Dict[str, List[Path]]:
    all_req_tags: Dict[str, List[Path]] = defaultdict(list)
    for feature_file in feature_files:
        req_tags = extract_req_tags_from_feature(feature_file, verbose)
        for req_id in req_tags:
            all_req_tags[req_id].append(feature_file)
    return all_req_tags


def _validate_req_tags(
    all_req_tags: Dict[str, List[Path]],
    all_requirements: Dict[str, str],
) -> Tuple[List[ValidationIssue], List[ValidationIssue], List[ValidationIssue]]:
    invalid_refs: List[ValidationIssue] = []
    missing_refs: List[ValidationIssue] = []
    errors: List[ValidationIssue] = []

    for req_id in sorted(all_req_tags.keys()):
        category = get_category_from_req_id(req_id)

        if not category:
            first_file = list(all_req_tags[req_id])[0] if all_req_tags[req_id] else None
            if first_file:
                errors.append(ValidationIssue.create(
                    'invalid_req_format',
                    first_file,
                    1,
                    1,
                    message=f"{req_id}: Invalid REQ ID format",
                    severity='error',
                    req_id=req_id,
                    reason='Invalid REQ ID format',
                    files=list(all_req_tags[req_id])
                ))
            continue

        # Block legacy/lookalike prefixes early with a clear fix hint.
        if category in DEPRECATED_CATEGORY_PREFIXES:
            replacement = DEPRECATED_CATEGORY_PREFIXES.get(category)
            first_file = list(all_req_tags[req_id])[0] if all_req_tags[req_id] else None
            if first_file:
                if replacement:
                    errors.append(ValidationIssue.create(
                        'deprecated_category',
                        first_file,
                        1,
                        1,
                        message=(
                            f"{req_id}: Deprecated category prefix: {category} "
                            f"(use {replacement})"
                        ),
                        severity='error',
                        req_id=req_id,
                        reason=f'Deprecated category prefix: {category} (use {replacement})',
                        files=list(all_req_tags[req_id]),
                        category=category,
                        suggested_category=replacement,
                    ))
                else:
                    errors.append(ValidationIssue.create(
                        'invalid_category',
                        first_file,
                        1,
                        1,
                        message=f"{req_id}: Invalid category prefix: {category}",
                        severity='error',
                        req_id=req_id,
                        reason=f'Invalid category prefix: {category}',
                        files=list(all_req_tags[req_id]),
                        category=category,
                    ))
            continue

        expected_file = get_expected_req_file(category)

        if not expected_file:
            first_file = list(all_req_tags[req_id])[0] if all_req_tags[req_id] else None
            if first_file:
                errors.append(ValidationIssue.create(
                    'unknown_category',
                    first_file,
                    1,
                    1,
                    message=f"{req_id}: Unknown category: {category}",
                    severity='error',
                    req_id=req_id,
                    reason=f'Unknown category: {category}',
                    files=list(all_req_tags[req_id])
                ))
            continue

        if req_id not in all_requirements:
            first_file = list(all_req_tags[req_id])[0] if all_req_tags[req_id] else None
            if first_file:
                missing_refs.append(ValidationIssue.create(
                    'missing_ref',
                    first_file,
                    1,
                    1,
                    message=f"{req_id} not found in {expected_file}",
                    severity='error',
                    req_id=req_id,
                    expected_file=expected_file,
                    files=list(all_req_tags[req_id])
                ))
        elif all_requirements[req_id] != expected_file:
            first_file = list(all_req_tags[req_id])[0] if all_req_tags[req_id] else None
            if first_file:
                invalid_refs.append(ValidationIssue.create(
                    'invalid_ref',
                    first_file,
                    1,
                    1,
                    message=(
                        f"{req_id} found in {all_requirements[req_id]}, "
                        f"expected {expected_file}"
                    ),
                    severity='error',
                    req_id=req_id,
                    expected_file=expected_file,
                    actual_file=all_requirements[req_id],
                    files=list(all_req_tags[req_id])
                ))
    return invalid_refs, missing_refs, errors


def _display_path(feature_file: Path, base_dir: Path) -> str:
    try:
        return str(feature_file.relative_to(base_dir))
    except ValueError:
        return str(feature_file)


def _emit_issue_list(
    output: OutputBuilder,
    issues: List[ValidationIssue],
    no_color: bool,
    *,
    header: str,
    features_dir: Optional[Path] = None,
    show_files: bool = False,
) -> None:
    if not issues:
        return
    if header == "warning":
        output.add_warnings_header()
    else:
        output.add_errors_header()
    for issue in issues:
        issue_msg = issue.format_message(no_color=no_color)
        if header == "warning":
            output.add_warning_line(issue_msg)
        else:
            output.add_error_line(issue_msg)
        if show_files and features_dir:
            files = issue.extra_fields.get('files', [])
            for feature_file in files[1:]:
                output.add_error_line(
                    f"    Also in: {_display_path(feature_file, features_dir)}"
                )


def validate_req_references(
    features_dir, requirements_dir, output, *, verbose=False, target_paths=None, no_color=False
):
    """
    Validate that all REQ tags in feature files exist in requirements files.

    Args:
        features_dir: Path to the features directory
        requirements_dir: Path to the requirements directory
        verbose: Whether to show detailed progress
        target_paths: Optional list of specific files or directories to check
        no_color: Whether to disable colored output

    Returns:
        Tuple of (total_refs, invalid_refs, missing_refs, errors, format_errors,
                 duplicate_errors, sequential_warnings)
    """
    # Determine workspace root (requirements_dir is docs/requirements, so parent.parent is root)
    workspace_root = requirements_dir.parent.parent

    all_requirements, requirements_data = _load_requirement_definitions(
        requirements_dir, output, verbose
    )

    # Scan all feature files for REQ tags
    output.add_verbose_line("Scanning feature files...")
    if verbose:
        output.add_blank_line("working_verbose")

    feature_files = _collect_feature_files(features_dir, target_paths, output)
    all_req_tags = _collect_req_tags(feature_files, verbose)

    if verbose:
        output.add_blank_line("working_verbose")
    output.add_verbose_line(
        f"Found {len(all_req_tags)} unique REQ tags "
        f"across {len(feature_files)} feature files"
    )
    output.add_blank_line("working_verbose")

    # Run format validation first
    format_errors = validate_req_format_errors(
        requirements_data, all_req_tags, workspace_root, features_dir
    )

    # Check for duplicates
    duplicate_errors = check_duplicate_requirements(requirements_data, workspace_root)

    # Check sequential numbering (warnings, not errors)
    sequential_warnings = check_sequential_numbering(requirements_data, workspace_root)

    invalid_refs, missing_refs, errors = _validate_req_tags(
        all_req_tags, all_requirements
    )

    # Report warnings first (sequential numbering gaps)
    _emit_issue_list(output, sequential_warnings, no_color, header="warning")
    _emit_issue_list(output, format_errors, no_color, header="error")
    _emit_issue_list(output, duplicate_errors, no_color, header="error")
    _emit_issue_list(
        output,
        errors,
        no_color,
        header="error",
        features_dir=features_dir,
        show_files=True
    )
    _emit_issue_list(
        output,
        invalid_refs,
        no_color,
        header="error",
        features_dir=features_dir,
        show_files=True
    )
    _emit_issue_list(
        output,
        missing_refs,
        no_color,
        header="error",
        features_dir=features_dir,
        show_files=True
    )

    # Summary
    total_refs = len(all_req_tags)
    valid_refs = total_refs - len(invalid_refs) - len(missing_refs) - len(errors)

    summary_items = [
        ("Total unique REQ references:", total_refs),
        ("Valid references:", valid_refs),
        ("Warnings (sequential gaps):", len(sequential_warnings)),
        ("Errors (format/category):", len(errors)),
        ("Format errors:", len(format_errors)),
        ("Duplicate requirements:", len(duplicate_errors)),
        ("Invalid references (wrong file):", len(invalid_refs)),
        ("Missing references:", len(missing_refs)),
    ]
    output.add_summary_header()
    output.add_summary_section(summary_items)

    return (
        total_refs, invalid_refs, missing_refs, errors, format_errors,
        duplicate_errors, sequential_warnings
    )


def main():
    """Main function to validate requirement references."""
    parser = argparse.ArgumentParser(
        description='Validate requirement references in feature files',
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

    args = parser.parse_args()

    # Determine the workspace root
    workspace_root = get_workspace_root()

    features_dir = workspace_root / FEATURES_DIR
    requirements_dir = workspace_root / DOCS_DIR / REQUIREMENTS_DIR

    # Verify requirements directory exists (before opening output file)
    if not requirements_dir.exists():
        print(f"Error: Requirements directory not found: {requirements_dir}")
        return 1

    # If no specific path provided, verify default features directory exists
    if not args.path and not features_dir.exists():
        print(f"Error: Features directory not found: {features_dir}")
        return 1

    # Parse comma-separated paths
    target_paths = parse_paths(args.path)

    no_color = args.nocolor or parse_no_color_flag(sys.argv)

    # Create output builder (header streams immediately if verbose)
    output = OutputBuilder(
        "Requirement Reference Validation",
        "Validates REQ references in feature files",
        no_color=no_color,
        verbose=args.verbose,
        output_file=args.output
    )

    # Check for feature stubs first
    if args.verbose:
        output.add_verbose_line("Checking for feature stubs...")
    stub_files = check_feature_stubs(features_dir, args.verbose, target_paths)

    if stub_files:
        output.add_errors_header()
        output.add_line(
            f"Found {len(stub_files)} feature file(s) with 8 or fewer content lines:",
            section="error"
        )
        output.add_blank_line("error")

        # Helper function to display file paths relative to features_dir or absolute
        def display_path(feature_file, base_dir):
            try:
                return str(feature_file.relative_to(base_dir))
            except ValueError:
                return str(feature_file)

        for stub in stub_files:
            # stub is a ValidationIssue
            output.add_error_line(
                f"  {display_path(Path(stub.file), features_dir)} "
                f"({stub.extra_fields.get('line_count', 0)} content lines)"
            )

    # Validate requirement references
    (
        _total, invalid, missing, errors, format_errors, duplicate_errors,
        _sequential_warnings
    ) = validate_req_references(
        features_dir, requirements_dir, output,
        verbose=args.verbose, target_paths=target_paths, no_color=no_color
    )

    # Return error code only if errors found (warnings don't cause failure)
    has_errors = (
        stub_files or invalid or missing or errors or format_errors
        or duplicate_errors
    )

    if not has_errors:
        if output.has_warnings():
            output.add_warnings_only_message(
                verbose_hint="Run with --verbose to see the full warning details.",
            )
        else:
            output.add_success_message("All requirement references are valid!")
    else:
        output.add_failure_message("Validation failed. Please fix the errors above.")

    output.print()
    return output.get_exit_code(args.no_fail)


if __name__ == '__main__':
    sys.exit(main())
