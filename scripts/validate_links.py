#!/usr/bin/env python3
"""
Validate all markdown links and anchors in the documentation.

This script validates internal links and anchors in markdown files,
properly excluding code blocks and inline code to avoid false positives.

Usage:
    python3 validate_links.py [options]
Options:
    --verbose, -v           Show detailed progress information
    --output, -o FILE       Write detailed output to FILE
    --check-coverage        Check that all requirements reference tech specs
    --path, -p PATHS        Check only the specified file(s) or
                            directory(ies) (recursive). Can be a single
                            path or comma-separated list of paths
    --nocolor, --no-color   Disable colored output
    --help, -h              Show this help message
Examples:
    # Basic validation
    python3 tmp/validate_links.py
    # Save output to file
    python3 tmp/validate_links.py --output tmp/validation_report.txt
    # Check requirements coverage
    python3 tmp/validate_links.py --check-coverage
    # Verbose with coverage check
    python3 tmp/validate_links.py --verbose --check-coverage

    # Check specific file
    python3 tmp/validate_links.py --path docs/tech_specs/api_file_management.md

    # Check specific directory
    python3 tmp/validate_links.py --path docs/requirements

    # Check multiple paths
    python3 tmp/validate_links.py --path docs/requirements,docs/tech_specs
"""

import os
import sys
from pathlib import Path
from collections import defaultdict
from dataclasses import dataclass
from typing import Dict, List, Optional, Tuple

from lib._validation_utils import (
    OutputBuilder, parse_no_color_flag, format_issue_message,
    parse_paths, ValidationIssue, find_markdown_files,
    is_safe_path, validate_file_name, validate_anchor,
    FileContentCache
)
from lib._link_extraction import extract_links
from lib._validate_links_helpers import extract_headings, suggest_anchor


def resolve_path(base_file, relative_path, repo_root: Path = None):
    """Resolve a relative path from a base file with security validation."""
    base_dir = os.path.dirname(base_file)
    resolved = os.path.normpath(os.path.join(base_dir, relative_path))

    # Convert to Path for safety check
    resolved_path = Path(resolved)
    if repo_root and not is_safe_path(resolved_path, repo_root):
        return None

    return resolved


def categorize_files(md_files):
    """Categorize markdown files by type."""
    requirements = []
    tech_specs = []
    root_docs = []
    other = []

    for f in md_files:
        path_str = str(f)
        if 'requirements/' in path_str:
            requirements.append(f)
        elif 'tech_specs/' in path_str:
            tech_specs.append(f)
        elif path_str.startswith('docs/') and '/' not in path_str[5:]:
            # Root-level docs (e.g., docs/README.md)
            root_docs.append(f)
        else:
            other.append(f)

    return requirements, tech_specs, root_docs, other


def validate_requirements_coverage(
    requirements_files, _tech_specs_files, _anchor_cache, file_cache=None
):
    """
    Check that requirements files properly reference tech specs.

    Args:
        requirements_files: List of requirement file paths
        tech_specs_files: List of tech spec file paths
        anchor_cache: Dictionary mapping file paths to heading dictionaries
        file_cache: Optional FileContentCache instance to use for reading files

    Returns list of issues found.
    """
    issues = []

    # Use cache if provided, otherwise create temporary one
    if file_cache is None:
        file_cache = FileContentCache()

    for req_file in requirements_files:
        links = extract_links(str(req_file), file_cache)
        has_tech_spec_link = False

        for _link_text, link_target, _line_num in links:
            if 'tech_specs/' in link_target:
                has_tech_spec_link = True
                break

        if not has_tech_spec_link:
            issues.append(ValidationIssue.create(
                'no_tech_spec_refs',
                req_file,
                1,
                1,
                message='No tech spec references found',
                severity='warning'
            ))

    return issues


@dataclass(frozen=True)
class LinkValidationContext:
    """Shared context for link validation helpers."""
    md_file: Path
    anchor_cache: dict
    file_cache: FileContentCache
    repo_root: Path
    verbose: bool


def build_anchor_suggestion(
    link_text,
    anchor,
    file_path,
    ctx,
    prefix
):
    """Build a suggested anchor string and details for verbose output."""
    suggestion_result = suggest_anchor(
        link_text, anchor, file_path, ctx.anchor_cache, verbose=ctx.verbose
    )
    if not suggestion_result:
        return None, None

    if ctx.verbose and len(suggestion_result) == 3:
        suggested_anchor, _confidence, details = suggestion_result
        return f"{prefix}{suggested_anchor}", details

    suggested_anchor, _confidence = suggestion_result
    return f"{prefix}{suggested_anchor}", None


def validate_internal_anchor_link(
    link_text,
    link_target,
    line_num,
    ctx
):
    """Validate an anchor-only link within the same file."""
    anchor = link_target[1:]
    if not validate_anchor(anchor):
        return [ValidationIssue.create(
            'unsafe_anchor',
            ctx.md_file,
            line_num,
            line_num,
            message=f"Unsafe internal anchor: #{anchor}",
            severity='error',
            suggestion=None,
            target=link_target,
            link_text=link_text
        )]

    file_headings = ctx.anchor_cache.get(str(ctx.md_file), {})
    if anchor in file_headings:
        return []

    suggestion, details = build_anchor_suggestion(
        link_text, anchor, ctx.md_file, ctx, "#"
    )
    return [ValidationIssue.create(
        'internal_anchor',
        ctx.md_file,
        line_num,
        line_num,
        message=f"Broken internal anchor #{anchor}",
        severity='error',
        suggestion=suggestion,
        target=link_target,
        link_text=link_text,
        suggestion_details=details
    )]


def split_link_target(link_target):
    """Split a link target into path and optional anchor."""
    if '#' in link_target:
        return link_target.split('#', 1)
    return link_target, None


def validate_file_link(
    link_text,
    link_target,
    line_num,
    ctx
):
    """Validate a file link, including optional anchors."""
    link_path, anchor = split_link_target(link_target)
    if anchor and not validate_anchor(anchor):
        return [ValidationIssue.create(
            'unsafe_anchor',
            ctx.md_file,
            line_num,
            line_num,
            message=f"Unsafe anchor: {anchor}",
            severity='error',
            suggestion=None,
            target=link_target,
            link_text=link_text
        )]

    if not link_path.endswith('/'):
        filename = os.path.basename(link_path)
        if filename and not validate_file_name(filename):
            return [ValidationIssue.create(
                'unsafe_path',
                ctx.md_file,
                line_num,
                line_num,
                message=f"Unsafe or invalid filename: {filename}",
                severity='error',
                suggestion=None,
                target=link_target,
                link_text=link_text
            )]

    resolved_path = resolve_path(str(ctx.md_file), link_path, ctx.repo_root)
    if resolved_path is None:
        return [ValidationIssue.create(
            'unsafe_path',
            ctx.md_file,
            line_num,
            line_num,
            message=f"Unsafe or invalid path: {link_path}",
            severity='error',
            suggestion=None,
            target=link_target,
            link_text=link_text
        )]

    if not os.path.exists(resolved_path):
        return [ValidationIssue.create(
            'missing_file',
            ctx.md_file,
            line_num,
            line_num,
            message=f"File not found: {link_path}",
            severity='error',
            suggestion=None,
            target=link_target,
            link_text=link_text
        )]

    if not anchor:
        return []

    if str(resolved_path) not in ctx.anchor_cache:
        headings_dict = extract_headings(str(resolved_path), ctx.file_cache)
        ctx.anchor_cache[str(resolved_path)] = headings_dict

    file_headings = ctx.anchor_cache.get(str(resolved_path), {})
    if anchor in file_headings:
        return []

    suggestion, details = build_anchor_suggestion(
        link_text, anchor, resolved_path, ctx, f"{link_path}#"
    )
    return [ValidationIssue.create(
        'broken_anchor',
        ctx.md_file,
        line_num,
        line_num,
        message=f"Broken anchor in {link_path}#{anchor}",
        severity='error',
        suggestion=suggestion,
        target=link_target,
        link_text=link_text,
        suggestion_details=details
    )]


def validate_link_target(
    link_text,
    link_target,
    line_num,
    ctx
):
    """Validate a single link target and return any issues."""
    if link_target.startswith('#'):
        return validate_internal_anchor_link(
            link_text, link_target, line_num, ctx
        )
    return validate_file_link(
        link_text,
        link_target,
        line_num,
        ctx
    )


def _parse_cli_args(
    argv: List[str],
) -> Tuple[bool, bool, bool, bool, Optional[str], Optional[str]]:
    verbose = '--verbose' in argv or '-v' in argv
    check_coverage = '--check-coverage' in argv
    no_color = parse_no_color_flag(argv)
    no_fail = '--no-fail' in argv
    output_file = None
    target_paths_str = None
    for i, arg in enumerate(argv):
        if arg in ('--output', '-o') and i + 1 < len(argv):
            output_file = argv[i + 1]
        elif arg in ('--path', '-p') and i + 1 < len(argv):
            target_paths_str = argv[i + 1]
    return verbose, check_coverage, no_color, no_fail, output_file, target_paths_str


def _build_output(
    verbose: bool,
    no_color: bool,
    output_file: Optional[str],
) -> OutputBuilder:
    return OutputBuilder(
        "Link and Anchor Validation",
        "Validates all markdown links and anchors",
        no_color=no_color,
        verbose=verbose,
        output_file=output_file
    )


def _default_exclude_dirs() -> set[str]:
    return {
        '.git', 'node_modules', 'vendor', '.venv', 'venv',
        '__pycache__', '.pytest_cache', 'dist', 'build',
        '.idea', '.vscode', 'tmp', '.cache'
    }


def _warn_non_markdown_targets(
    target_paths: List[str],
    output: OutputBuilder,
) -> None:
    for target_path in target_paths:
        target = Path(target_path)
        if target.exists() and target.is_file() and target.suffix != '.md':
            output.add_warning_line(
                f"Target file is not a markdown file: {target_path}"
            )


def _build_anchor_cache(
    md_files: List[Path],
    file_cache: FileContentCache,
    output: OutputBuilder,
    verbose: bool,
) -> Dict[str, dict]:
    output.add_verbose_line("Building anchor cache...")
    anchor_cache: Dict[str, dict] = {}
    for md_file in md_files:
        headings_dict = extract_headings(str(md_file), file_cache)
        anchor_cache[str(md_file)] = headings_dict
    if verbose:
        total_anchors = sum(len(headings) for headings in anchor_cache.values())
        output.add_verbose_line(
            f"  Cached {total_anchors} anchors from {len(anchor_cache)} files"
        )
    output.add_blank_line("working_verbose")
    return anchor_cache


def _validate_links(
    md_files: List[Path],
    anchor_cache: Dict[str, dict],
    file_cache: FileContentCache,
    root_dir: Path,
    verbose: bool,
) -> Tuple[List[ValidationIssue], int, int]:
    broken_links: List[ValidationIssue] = []
    total_links = 0
    files_with_links = 0
    for md_file in md_files:
        links = extract_links(str(md_file), file_cache)
        if links:
            files_with_links += 1
        ctx = LinkValidationContext(
            md_file=md_file,
            anchor_cache=anchor_cache,
            file_cache=file_cache,
            repo_root=root_dir,
            verbose=verbose
        )
        for link_text, link_target, line_num in links:
            total_links += 1
            broken_links.extend(validate_link_target(
                link_text,
                link_target,
                line_num,
                ctx
            ))
    return broken_links, total_links, files_with_links


def _check_coverage_if_requested(
    check_coverage: bool,
    *,
    requirements_files: List[Path],
    tech_spec_files: List[Path],
    anchor_cache: Dict[str, dict],
    file_cache: FileContentCache,
    output: OutputBuilder,
) -> List:
    if not check_coverage:
        return []
    output.add_verbose_line("Checking requirements coverage...")
    coverage_issues = validate_requirements_coverage(
        requirements_files,
        tech_spec_files,
        anchor_cache,
        file_cache
    )
    output.add_blank_line("working_verbose")
    return coverage_issues


def _format_broken_link_issue(
    broken: ValidationIssue,
    file_path: str,
    no_color: bool,
) -> str:
    issue_msg = broken.message
    if " in " in issue_msg:
        link_info = issue_msg.split(" in ", 1)[1]
    else:
        link_info = issue_msg
    return format_issue_message(
        "error",
        "Broken link",
        file_path,
        line_num=broken.start_line,
        message=link_info,
        suggestion=broken.suggestion,
        no_color=no_color
    )


def _emit_suggestion_details(
    output: OutputBuilder,
    broken: ValidationIssue,
    verbose: bool,
) -> None:
    if not (verbose and broken.extra_fields.get('suggestion_details')):
        return
    details = broken.extra_fields['suggestion_details']
    scores = details.get('scores', {})
    output.add_verbose_line(
        f"    Suggestion scores: word_match={scores.get('word_match', 0):.1f}, "
        f"anchor_similarity={scores.get('anchor_similarity', 0):.1f}, "
        f"context={scores.get('context', 0):.1f}, "
        f"normalization={scores.get('normalization', 0):.1f}, "
        f"total={details.get('total_score', 0):.1f}"
    )


def _emit_broken_links_group(
    output: OutputBuilder,
    section_title: str,
    *,
    file_paths: List[str],
    by_file: Dict[str, List[ValidationIssue]],
    no_color: bool,
    verbose: bool,
    add_blank: bool = False,
) -> None:
    if not file_paths:
        return
    if add_blank:
        output.add_blank_line("error")
    output.add_line(section_title, section="error")
    for file_path in file_paths:
        output.add_error_line(f"{file_path}:")
        for broken in by_file[file_path]:
            error_output = _format_broken_link_issue(broken, file_path, no_color)
            output.add_error_line(error_output)
            _emit_suggestion_details(output, broken, verbose)


def _emit_coverage_warnings(
    output: OutputBuilder,
    coverage_issues: List,
    no_color: bool,
) -> None:
    if not coverage_issues:
        return
    output.add_warnings_header()
    output.add_line(
        "The following requirements files don't reference any tech specs:",
        section="warning"
    )
    output.add_blank_line("warning")
    for issue in coverage_issues:
        if isinstance(issue, ValidationIssue):
            warning_msg = issue.format_message(no_color=no_color)
        else:
            warning_msg = format_issue_message(
                "warning",
                "No tech spec refs",
                issue.get('file', ''),
                message=issue.get('issue', 'No tech spec references found'),
                no_color=no_color
            )
        output.add_warning_line(warning_msg)
    output.add_blank_line("warning")
    output.add_line(
        "Note: After adding tech spec references, verify that each "
        "reference points to the correct content.",
        section="warning"
    )


def _emit_broken_links(
    output: OutputBuilder,
    broken_links: List[ValidationIssue],
    no_color: bool,
    verbose: bool,
) -> None:
    if not broken_links:
        return
    output.add_errors_header()
    by_file: Dict[str, List[ValidationIssue]] = defaultdict(list)
    for broken in broken_links:
        by_file[broken.file].append(broken)

    req_files = sorted([f for f in by_file.keys() if 'requirements/' in f])
    spec_files = sorted([f for f in by_file.keys() if 'tech_specs/' in f])
    docs_root_files = sorted([
        f for f in by_file.keys()
        if f.startswith('docs/') and '/' not in f[5:]
    ])
    other_files_broken = sorted([
        f for f in by_file.keys()
        if f not in req_files
        and f not in spec_files
        and f not in docs_root_files
    ])

    _emit_broken_links_group(
        output,
        "## Requirements Files",
        file_paths=req_files,
        by_file=by_file,
        no_color=no_color,
        verbose=verbose
    )
    _emit_broken_links_group(
        output,
        "## Tech Spec Files",
        file_paths=spec_files,
        by_file=by_file,
        no_color=no_color,
        verbose=verbose,
        add_blank=bool(req_files)
    )
    _emit_broken_links_group(
        output,
        "## Root Documentation Files",
        file_paths=docs_root_files,
        by_file=by_file,
        no_color=no_color,
        verbose=verbose,
        add_blank=bool(req_files or spec_files)
    )
    _emit_broken_links_group(
        output,
        "## Other Files",
        file_paths=other_files_broken,
        by_file=by_file,
        no_color=no_color,
        verbose=verbose,
        add_blank=bool(req_files or spec_files or docs_root_files)
    )

    has_tech_spec_links = any(
        'tech_specs/' in (
            broken.extra_fields.get('target', '')
            if isinstance(broken, ValidationIssue)
            else broken.get('target', '')
        )
        for broken in broken_links
    )
    if has_tech_spec_links:
        output.add_blank_line("error")
        output.add_line(
            "Note: After fixing broken links to tech specs, verify that each "
            "updated reference points to the correct content.",
            section="error"
        )


def main():
    """Main validation function."""
    # Show help if requested
    if '--help' in sys.argv or '-h' in sys.argv:
        print(__doc__)
        return 0

    verbose, check_coverage, no_color, no_fail, output_file, target_paths_str = (
        _parse_cli_args(sys.argv)
    )

    # Parse comma-separated paths
    target_paths = parse_paths(target_paths_str)

    # Create output builder (header streams immediately if verbose)
    output = _build_output(verbose, no_color, output_file)

    # Find all markdown files in the repository
    # Start from current directory (repository root when called from Makefile)
    root_dir = Path(".")

    # Directories to exclude from scanning (only when no target path is specified)
    exclude_dirs = _default_exclude_dirs()

    # Find all markdown files using shared utility
    md_files = find_markdown_files(
        target_paths=target_paths,
        root_dir=root_dir,
        exclude_dirs=exclude_dirs,
        verbose=verbose
    )

    # Handle warnings for non-markdown files when target_paths is specified
    if target_paths:
        _warn_non_markdown_targets(target_paths, output)

    if not md_files:
        output.add_error_line("No markdown files found")
        output.print()
        return 1

    # Create file content cache to avoid repeated reads
    file_cache = FileContentCache()

    # Categorize files
    requirements_files, tech_spec_files, root_docs, other_files = (
        categorize_files(md_files)
    )

    output.add_verbose_line(f"Scanning {len(md_files)} markdown files:")
    output.add_verbose_line(f"  Requirements: {len(requirements_files)}")
    output.add_verbose_line(f"  Tech Specs:   {len(tech_spec_files)}")
    output.add_verbose_line(f"  Root Docs:    {len(root_docs)}")
    output.add_verbose_line(f"  Other:        {len(other_files)}")
    output.add_blank_line("working_verbose")

    # Build anchor cache (now includes heading text and metadata)
    anchor_cache = _build_anchor_cache(md_files, file_cache, output, verbose)

    # Validate links
    broken_links, total_links, files_with_links = _validate_links(
        md_files,
        anchor_cache,
        file_cache,
        root_dir,
        verbose
    )

    # Check requirements coverage if requested
    coverage_issues = _check_coverage_if_requested(
        check_coverage,
        requirements_files=requirements_files,
        tech_spec_files=tech_spec_files,
        anchor_cache=anchor_cache,
        file_cache=file_cache,
        output=output
    )

    summary_items = [
        ("Files scanned:", len(md_files)),
        ("  Requirements:", len(requirements_files)),
        ("  Tech Specs:", len(tech_spec_files)),
        ("  Root Docs:", len(root_docs)),
        ("  Other:", len(other_files)),
        ("Files with links:", files_with_links),
        ("Total links checked:", total_links),
        ("Broken links found:", len(broken_links)),
    ]
    output.add_summary_section(summary_items)

    # Report coverage issues first
    _emit_coverage_warnings(output, coverage_issues, no_color)

    # Report broken links
    _emit_broken_links(output, broken_links, no_color, verbose)

    # Final status
    if broken_links:
        output.add_failure_message("Validation failed. Please fix the errors above.")
    elif coverage_issues:
        msg = "All links are valid. Review the warnings above."
        if check_coverage:
            msg = (
                "All links are valid. All requirements reference tech specs. "
                "Review the warnings above."
            )
        output.add_warnings_only_message(
            message=msg,
            verbose_hint="Run with --verbose to see the full warning details.",
        )
    else:
        output.add_success_message("All links are valid!")
        if check_coverage:
            output.add_success_message("All requirements reference tech specs!")
    output.print()
    return output.get_exit_code(no_fail)


if __name__ == "__main__":
    sys.exit(main())
