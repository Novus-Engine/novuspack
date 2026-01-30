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

import re
import os
import sys
from pathlib import Path
from collections import defaultdict
from dataclasses import dataclass
from typing import List

scripts_dir = Path(__file__).parent
lib_dir = scripts_dir / "lib"

# Import shared utilities
for module_path in (str(scripts_dir), str(lib_dir)):
    if module_path not in sys.path:
        sys.path.insert(0, module_path)

# Import shared utilities
from lib._validation_utils import (  # noqa: E402
    OutputBuilder, parse_no_color_flag, format_issue_message,
    parse_paths, ValidationIssue, find_markdown_files,
    is_safe_path, validate_file_name, validate_anchor,
    extract_headings_with_anchors, FileContentCache
)
from lib._link_extraction import extract_links  # noqa: E402


# Compiled regex patterns for performance (module level)
_RE_MARKDOWN_FORMAT = re.compile(r'[*_`]')
_RE_SPECIAL_CHARS = re.compile(r'[^a-zA-Z0-9 .-]')
_RE_SPLIT_WORDS = re.compile(r'[\s.\-]+')
_RE_CAMEL_CASE = re.compile(r'([a-z])([A-Z])')
_RE_NUMBERING_PREFIX = re.compile(r'^([0-9]+(?:\.[0-9]+)*)\.?\s+(.+)$')


# extract_headings() now uses shared utility extract_headings_with_anchors()
def extract_headings(file_path, file_cache=None):
    """
    Extract all headings from a markdown file and generate anchors.

    Args:
        file_path: Path to the file
        file_cache: Optional FileContentCache instance to use for reading files

    Returns:
        dict: Mapping of anchor -> (heading_text, heading_level, line_number)
    """
    return extract_headings_with_anchors(Path(file_path), file_cache=file_cache)


def normalize_text_for_matching(text: str) -> str:
    """
    Normalize text for matching by removing markdown formatting,
    converting to lowercase, and removing common words.

    Args:
        text: Text to normalize

    Returns:
        Normalized text string
    """
    # Remove markdown formatting
    text = _RE_MARKDOWN_FORMAT.sub('', text)
    # Convert to lowercase
    text = text.lower()
    # Remove special characters except spaces, hyphens, and dots
    text = _RE_SPECIAL_CHARS.sub('', text)
    return text.strip()


def extract_words(text: str) -> List[str]:
    """
    Extract words from text, handling various separators.

    Handles:
    - Spaces: "Add File" -> ["add", "file"]
    - Dots: "Package.AddFile" -> ["package", "add", "file"]
    - Hyphens: "add-file" -> ["add", "file"]
    - CamelCase: "AddFile" -> ["add", "file"]
    - Mixed: "Package.AddFile" -> ["package", "add", "file"]

    Args:
        text: Text to extract words from

    Returns:
        List of normalized words (all lowercase)
    """
    # Remove markdown formatting first (preserve case for camelCase detection)
    text = _RE_MARKDOWN_FORMAT.sub('', text)

    # Split on spaces, dots, and hyphens (preserve case)
    words = _RE_SPLIT_WORDS.split(text)

    # Further split camelCase words and normalize to lowercase
    all_words = []
    for word in words:
        if not word:
            continue
        # Split camelCase: "AddFile" -> ["Add", "File"], then lowercase
        camel_split = _RE_CAMEL_CASE.sub(r'\1 \2', word)
        camel_words = camel_split.split()
        # Convert to lowercase and add
        all_words.extend([w.lower() for w in camel_words if w])

    # Filter out empty strings and common stop words
    stop_words = {'the', 'a', 'an', 'and', 'or', 'of', 'in', 'on', 'at', 'to', 'for', 'with', 'by'}
    return [w for w in all_words if w and w not in stop_words]


def strip_numbering_prefix(text: str) -> str:
    """
    Strip numbering prefix from heading text (e.g., "1.2.3 Add File" -> "Add File").

    Args:
        text: Heading text that may contain numbering

    Returns:
        Text with numbering prefix removed
    """
    # Match pattern like "1.2.3 " or "1 " at the start
    match = _RE_NUMBERING_PREFIX.match(text)
    if match:
        return match.group(2).strip()
    return text


def calculate_word_match_score(link_words: List[str], heading_words: List[str]) -> float:
    """
    Calculate word matching score between link text and heading text.

    Args:
        link_words: List of words from link text
        heading_words: List of words from heading text

    Returns:
        Score from 0-100 based on word matching
    """
    if not link_words or not heading_words:
        return 0.0

    link_set = set(link_words)
    heading_set = set(heading_words)

    # Exact set match
    if link_set == heading_set:
        return 100.0

    # All link words in heading
    if link_set.issubset(heading_set):
        return 90.0

    # Most link words in heading
    matching_words = link_set.intersection(heading_set)
    if matching_words:
        match_ratio = len(matching_words) / len(link_set)
        return 60.0 + (match_ratio * 30.0)  # 60-90 range

    # Partial word matches (substring matching)
    partial_matches = 0
    for link_word in link_words:
        for heading_word in heading_words:
            if link_word in heading_word or heading_word in link_word:
                partial_matches += 1
                break

    if partial_matches > 0:
        partial_ratio = partial_matches / len(link_words)
        return 20.0 + (partial_ratio * 40.0)  # 20-60 range

    return 0.0


def suggest_anchor(link_text, broken_anchor, target_file, heading_cache, verbose=False):
    """
    Suggest correct anchor based on weighted heuristics.

    Args:
        link_text: Text from the markdown link
        broken_anchor: The broken anchor that was not found
        target_file: Path to the target file
        heading_cache: Dictionary mapping file paths to heading dictionaries
                      (anchor -> (heading_text, heading_level, line_number))
        verbose: If True, return detailed scoring information

    Returns:
        Tuple of (suggested_anchor, confidence_score) or None if no good match found.
        If verbose=True, returns (suggested_anchor, confidence_score, score_details)
    """
    # Get headings for target file
    headings_dict = heading_cache.get(str(target_file), {})
    if not headings_dict:
        return None

    # Normalize link text
    normalized_link_text = normalize_text_for_matching(link_text)
    link_words = extract_words(link_text)

    # Extract words from broken anchor (for additional matching)
    broken_anchor_words = extract_words(broken_anchor.replace('-', ' '))

    best_match = None
    best_score = 0.0
    best_details = {}

    for anchor, (heading_text, heading_level, line_num) in headings_dict.items():
        # Skip if anchor is invalid (shouldn't happen, but safety check)
        if not validate_anchor(anchor):
            continue

        # Strip numbering from heading text
        heading_text_no_numbering = strip_numbering_prefix(heading_text)

        # Normalize heading text
        normalized_heading = normalize_text_for_matching(heading_text_no_numbering)
        heading_words = extract_words(heading_text_no_numbering)

        # Calculate scores for different heuristics
        scores = {}

        # Heuristic 1: Word matching (weight: 0.4)
        word_score = calculate_word_match_score(link_words, heading_words)
        scores['word_match'] = word_score
        weighted_word = word_score * 0.4

        # Heuristic 2: Anchor similarity (weight: 0.3)
        anchor_score = 0.0
        if anchor == broken_anchor:
            anchor_score = 100.0
        elif broken_anchor in anchor:
            anchor_score = 70.0
        elif anchor in broken_anchor:
            anchor_score = 50.0
        else:
            # Check word overlap in anchors
            anchor_words = extract_words(anchor.replace('-', ' '))
            anchor_word_score = calculate_word_match_score(broken_anchor_words, anchor_words)
            anchor_score = anchor_word_score * 0.5  # Lower weight for anchor word matching
        scores['anchor_similarity'] = anchor_score
        weighted_anchor = anchor_score * 0.3

        # Heuristic 3: Context matching (weight: 0.2)
        context_score = 0.0
        # Prefer H2/H3 headings (more likely to be main sections)
        if heading_level == 2:
            context_score += 30.0
        elif heading_level == 3:
            context_score += 20.0
        elif heading_level == 4:
            context_score += 10.0
        # Slight preference for earlier headings (but minimal impact)
        if line_num < 100:
            context_score += 5.0
        scores['context'] = context_score
        weighted_context = context_score * 0.2

        # Heuristic 4: Normalization quality (weight: 0.1)
        # Check if normalized versions are similar
        norm_score = 0.0
        if normalized_link_text == normalized_heading:
            norm_score = 100.0
        elif (normalized_link_text in normalized_heading or
              normalized_heading in normalized_link_text):
            norm_score = 60.0
        scores['normalization'] = norm_score
        weighted_norm = norm_score * 0.1

        # Calculate total weighted score
        total_score = weighted_word + weighted_anchor + weighted_context + weighted_norm

        if total_score > best_score:
            best_score = total_score
            best_match = anchor
            best_details = {
                'heading_text': heading_text,
                'heading_level': heading_level,
                'line_num': line_num,
                'scores': scores,
                'total_score': total_score
            }

    # Only return if score is above 70% threshold
    if best_match and best_score >= 70.0:
        if verbose:
            return (best_match, best_score, best_details)
        return (best_match, best_score)

    return None


# Path validation functions now imported from _validation_utils


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
            issues.append(ValidationIssue(
                'no_tech_spec_refs',
                req_file,
                1,
                1,
                'No tech spec references found',
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
        return [ValidationIssue(
            'unsafe_anchor',
            ctx.md_file,
            line_num,
            line_num,
            f"Unsafe internal anchor: #{anchor}",
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
    return [ValidationIssue(
        'internal_anchor',
        ctx.md_file,
        line_num,
        line_num,
        f"Broken internal anchor #{anchor}",
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
        return [ValidationIssue(
            'unsafe_anchor',
            ctx.md_file,
            line_num,
            line_num,
            f"Unsafe anchor: {anchor}",
            severity='error',
            suggestion=None,
            target=link_target,
            link_text=link_text
        )]

    if not link_path.endswith('/'):
        filename = os.path.basename(link_path)
        if filename and not validate_file_name(filename):
            return [ValidationIssue(
                'unsafe_path',
                ctx.md_file,
                line_num,
                line_num,
                f"Unsafe or invalid filename: {filename}",
                severity='error',
                suggestion=None,
                target=link_target,
                link_text=link_text
            )]

    resolved_path = resolve_path(str(ctx.md_file), link_path, ctx.repo_root)
    if resolved_path is None:
        return [ValidationIssue(
            'unsafe_path',
            ctx.md_file,
            line_num,
            line_num,
            f"Unsafe or invalid path: {link_path}",
            severity='error',
            suggestion=None,
            target=link_target,
            link_text=link_text
        )]

    if not os.path.exists(resolved_path):
        return [ValidationIssue(
            'missing_file',
            ctx.md_file,
            line_num,
            line_num,
            f"File not found: {link_path}",
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
    return [ValidationIssue(
        'broken_anchor',
        ctx.md_file,
        line_num,
        line_num,
        f"Broken anchor in {link_path}#{anchor}",
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


def main():
    """Main validation function."""
    # Show help if requested
    if '--help' in sys.argv or '-h' in sys.argv:
        print(__doc__)
        return 0

    # Parse command line arguments
    verbose = '--verbose' in sys.argv or '-v' in sys.argv
    check_coverage = '--check-coverage' in sys.argv
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

    # Create output builder (header streams immediately if verbose)
    output = OutputBuilder(
        "Link and Anchor Validation",
        "Validates all markdown links and anchors",
        no_color=no_color,
        verbose=verbose,
        output_file=output_file
    )

    # Find all markdown files in the repository
    # Start from current directory (repository root when called from Makefile)
    root_dir = Path(".")

    # Directories to exclude from scanning (only when no target path is specified)
    exclude_dirs = {
        '.git', 'node_modules', 'vendor', '.venv', 'venv',
        '__pycache__', '.pytest_cache', 'dist', 'build',
        '.idea', '.vscode', 'tmp', '.cache'
    }

    # Find all markdown files using shared utility
    md_files = find_markdown_files(
        target_paths=target_paths,
        root_dir=root_dir,
        exclude_dirs=exclude_dirs,
        verbose=verbose
    )

    # Handle warnings for non-markdown files when target_paths is specified
    if target_paths:
        for target_path in target_paths:
            target = Path(target_path)
            if target.exists() and target.is_file() and target.suffix != '.md':
                output.add_warning_line(
                    f"Target file is not a markdown file: {target_path}"
                )

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
    output.add_verbose_line("Building anchor cache...")
    anchor_cache = {}
    for md_file in md_files:
        headings_dict = extract_headings(str(md_file), file_cache)
        anchor_cache[str(md_file)] = headings_dict

    if verbose:
        total_anchors = sum(len(headings) for headings in anchor_cache.values())
        output.add_verbose_line(f"  Cached {total_anchors} anchors from {len(anchor_cache)} files")
    output.add_blank_line("working_verbose")

    # Validate links
    broken_links: List[ValidationIssue] = []
    total_links = 0
    files_with_links = 0

    if verbose:
        output.add_verbose_line("Validating links...")
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

    # Check requirements coverage if requested
    coverage_issues = []
    if check_coverage:
        output.add_verbose_line("Checking requirements coverage...")
        coverage_issues = validate_requirements_coverage(
            requirements_files,
            tech_spec_files,
            anchor_cache,
            file_cache
        )
        output.add_blank_line("working_verbose")

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
    if coverage_issues:
        output.add_warnings_header()
        output.add_line(
            "The following requirements files don't reference any tech specs:",
            section="warning"
        )
        output.add_blank_line("warning")

        for issue in coverage_issues:
            # Convert ValidationIssue to format message if needed
            if isinstance(issue, ValidationIssue):
                warning_msg = issue.format_message(no_color=no_color)
            else:
                warning_msg = format_issue_message(
                    "warning",
                    "No tech spec refs",
                    issue.get('file', ''),
                    None,
                    issue.get('issue', 'No tech spec references found'),
                    no_color
                )
            output.add_warning_line(warning_msg)

        output.add_blank_line("warning")
        output.add_line(
            "Note: After adding tech spec references, verify that each "
            "reference points to the correct content.",
            section="warning"
        )

    # Report broken links
    if broken_links:
        output.add_errors_header()

        # Group by file for better readability
        by_file = defaultdict(list)
        for broken in broken_links:
            # broken_links contains ValidationIssue objects
            file_key = broken.file
            by_file[file_key].append(broken)

        # Separate files by category for better organization
        req_files = sorted([f for f in by_file.keys()
                            if 'requirements/' in f])
        spec_files = sorted([f for f in by_file.keys()
                             if 'tech_specs/' in f])
        docs_root_files = sorted([f for f in by_file.keys()
                                  if f.startswith('docs/')
                                  and '/' not in f[5:]])
        other_files_broken = sorted([f for f in by_file.keys()
                                     if f not in req_files
                                     and f not in spec_files
                                     and f not in docs_root_files])

        if req_files:
            output.add_line("## Requirements Files", section="error")
            for file_path in req_files:
                output.add_error_line(f"{file_path}:")
                for broken in by_file[file_path]:
                    # broken is a ValidationIssue object
                    issue_msg = broken.message
                    if " in " in issue_msg:
                        link_info = issue_msg.split(" in ", 1)[1]
                    else:
                        link_info = issue_msg
                    error_output = format_issue_message(
                        "error",
                        "Broken link",
                        file_path,
                        broken.start_line,
                        link_info,
                        broken.suggestion,
                        no_color
                    )
                    output.add_error_line(error_output)

                    # Add verbose output for suggestion details
                    if verbose and broken.extra_fields.get('suggestion_details'):
                        details = broken.extra_fields['suggestion_details']
                        scores = details.get('scores', {})
                        output.add_verbose_line(
                            f"    Suggestion scores: word_match={scores.get('word_match', 0):.1f}, "
                            f"anchor_similarity={scores.get('anchor_similarity', 0):.1f}, "
                            f"context={scores.get('context', 0):.1f}, "
                            f"normalization={scores.get('normalization', 0):.1f}, "
                            f"total={details.get('total_score', 0):.1f}"
                        )

        if spec_files:
            output.add_blank_line("error")
            output.add_line("## Tech Spec Files", section="error")
            for file_path in spec_files:
                output.add_error_line(f"{file_path}:")
                for broken in by_file[file_path]:
                    # broken is a ValidationIssue object
                    issue_msg = broken.message
                    if " in " in issue_msg:
                        link_info = issue_msg.split(" in ", 1)[1]
                    else:
                        link_info = issue_msg
                    error_output = format_issue_message(
                        "error",
                        "Broken link",
                        file_path,
                        broken.start_line,
                        link_info,
                        broken.suggestion,
                        no_color
                    )
                    output.add_error_line(error_output)

                    # Add verbose output for suggestion details
                    if verbose and broken.extra_fields.get('suggestion_details'):
                        details = broken.extra_fields['suggestion_details']
                        scores = details.get('scores', {})
                        output.add_verbose_line(
                            f"    Suggestion scores: word_match={scores.get('word_match', 0):.1f}, "
                            f"anchor_similarity={scores.get('anchor_similarity', 0):.1f}, "
                            f"context={scores.get('context', 0):.1f}, "
                            f"normalization={scores.get('normalization', 0):.1f}, "
                            f"total={details.get('total_score', 0):.1f}"
                        )

        if docs_root_files:
            output.add_blank_line("error")
            output.add_line("## Root Documentation Files", section="error")
            for file_path in docs_root_files:
                output.add_error_line(f"{file_path}:")
                for broken in by_file[file_path]:
                    # broken is a ValidationIssue object
                    issue_msg = broken.message
                    if " in " in issue_msg:
                        link_info = issue_msg.split(" in ", 1)[1]
                    else:
                        link_info = issue_msg
                    error_output = format_issue_message(
                        "error",
                        "Broken link",
                        file_path,
                        broken.start_line,
                        link_info,
                        broken.suggestion,
                        no_color
                    )
                    output.add_error_line(error_output)

                    # Add verbose output for suggestion details
                    if verbose and broken.extra_fields.get('suggestion_details'):
                        details = broken.extra_fields['suggestion_details']
                        scores = details.get('scores', {})
                        output.add_verbose_line(
                            f"    Suggestion scores: word_match={scores.get('word_match', 0):.1f}, "
                            f"anchor_similarity={scores.get('anchor_similarity', 0):.1f}, "
                            f"context={scores.get('context', 0):.1f}, "
                            f"normalization={scores.get('normalization', 0):.1f}, "
                            f"total={details.get('total_score', 0):.1f}"
                        )

        if other_files_broken:
            output.add_blank_line("error")
            output.add_line("## Other Files", section="error")
            for file_path in other_files_broken:
                output.add_error_line(f"{file_path}:")
                for broken in by_file[file_path]:
                    # broken is a ValidationIssue object
                    issue_msg = broken.message
                    if " in " in issue_msg:
                        link_info = issue_msg.split(" in ", 1)[1]
                    else:
                        link_info = issue_msg
                    error_output = format_issue_message(
                        "error",
                        "Broken link",
                        file_path,
                        broken.start_line,
                        link_info,
                        broken.suggestion,
                        no_color
                    )
                    output.add_error_line(error_output)

                    # Add verbose output for suggestion details
                    if verbose and broken.extra_fields.get('suggestion_details'):
                        details = broken.extra_fields['suggestion_details']
                        scores = details.get('scores', {})
                        output.add_verbose_line(
                            f"    Suggestion scores: word_match={scores.get('word_match', 0):.1f}, "
                            f"anchor_similarity={scores.get('anchor_similarity', 0):.1f}, "
                            f"context={scores.get('context', 0):.1f}, "
                            f"normalization={scores.get('normalization', 0):.1f}, "
                            f"total={details.get('total_score', 0):.1f}"
                        )

        # Check if any broken links point to tech specs
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

    # Final status
    if not broken_links and not coverage_issues:
        output.add_success_message("All links are valid!")
        if check_coverage:
            output.add_success_message("All requirements reference tech specs!")
    else:
        output.add_failure_message("Validation failed. Please fix the errors above.")
    output.print()
    return output.get_exit_code(no_fail)


if __name__ == "__main__":
    sys.exit(main())
