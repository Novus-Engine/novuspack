"""
Spec and requirement file discovery; index-file ref check; heading ref checks.
"""

import re
import sys
from pathlib import Path
from typing import List, Tuple, Optional, Dict

from lib._validation_utils import (
    is_in_dot_directory,
    find_markdown_files,
    FileContentCache,
    ValidationIssue,
)
from lib._link_extraction import extract_links

from lib.audit_requirements._audit_requirements_config import _SCRIPT_ERROR_EXCEPTIONS


# Per-requirement tech spec link density thresholds.
#
# Used by `check_requirement_tech_spec_link_thresholds()`:
# - Emit WARNING if a single requirement contains at least WARN_THRESHOLD links
# - Emit ERROR   if a single requirement contains at least ERROR_THRESHOLD links
#
# Keeping these at module scope makes it easy to tune the audit without changing logic.
REQ_TECH_SPEC_LINK_WARN_THRESHOLD = 2
REQ_TECH_SPEC_LINK_ERROR_THRESHOLD = 4


def is_index_file(filename: str) -> bool:
    """Return True if filename is an index file (_index.md)."""
    return '_index.md' in filename or filename.endswith('_index.md')


def find_tech_specs(tech_specs_dir: Path) -> List[Path]:
    """Find all tech spec markdown files, excluding index files."""
    tech_specs = []
    exclude_patterns = ['func_signatures_index', '_overview', '_main']

    for md_file in sorted(tech_specs_dir.glob('*.md')):
        if is_in_dot_directory(md_file):
            continue
        if is_index_file(md_file.name):
            continue
        if not any(pattern in md_file.name for pattern in exclude_patterns):
            tech_specs.append(md_file)

    return tech_specs


def _get_requirement_files_for_index_check(
    requirements_dir: Path,
    target_paths: Optional[List[str]] = None
) -> List[Path]:
    """Return list of requirement files to scan for index refs."""
    if target_paths:
        req_files = []
        for target_path in target_paths:
            target = Path(target_path)
            if not target.exists():
                continue
            if target.is_file():
                if target.suffix == '.md':
                    req_files.append(target)
            else:
                req_files.extend([
                    f for f in target.glob('*.md')
                    if not is_in_dot_directory(f)
                ])
        return req_files
    return [
        f for f in requirements_dir.glob('*.md')
        if not is_in_dot_directory(f)
    ]


def _collect_index_refs_in_file(
    req_file: Path,
    file_cache: FileContentCache,
    index_pattern: re.Pattern,
) -> List[Tuple[Path, int, str]]:
    """Scan one requirement file for index refs; return (file, line_num, reference)."""
    errors = []
    lines = file_cache.get_lines(req_file)
    for line_num, line in enumerate(lines, 1):
        if '_index.md' not in line:
            continue
        for match in index_pattern.finditer(line):
            errors.append((req_file, line_num, match.group(0)))
    return errors


def check_index_file_references(
    requirements_dir: Path,
    file_cache: Optional[FileContentCache] = None,
    verbose: bool = False,
    target_paths: Optional[List[str]] = None,
) -> List[Tuple[Path, int, str]]:
    """
    Check for references to index files in requirements.

    Returns list of (file_path, line_num, reference) for index file references.
    """
    errors = []
    index_pattern = re.compile(r'\(\.\./tech_specs/[^)]*_index\.md[^)]*\)')

    if file_cache is None:
        file_cache = FileContentCache()

    req_files = _get_requirement_files_for_index_check(requirements_dir, target_paths)

    for req_file in req_files:
        try:
            found = _collect_index_refs_in_file(req_file, file_cache, index_pattern)
            errors.extend(found)
        except (IOError, OSError) as e:
            if verbose:
                print(f"  Warning: Could not read {req_file}: {e}", file=sys.stderr)
        except UnicodeDecodeError as e:
            if verbose:
                print(
                    f"  Warning: Could not decode {req_file} (encoding issue): {e}",
                    file=sys.stderr
                )
        except _SCRIPT_ERROR_EXCEPTIONS as e:
            if verbose:
                print(f"  Warning: Unexpected error reading {req_file}: {e}", file=sys.stderr)

    return errors


def get_requirement_files(
    requirements_dir: Path,
    target_paths: Optional[List[str]] = None,
) -> List[Path]:
    """Return list of requirement files to scan."""
    return find_markdown_files(
        target_paths=target_paths,
        default_dir=requirements_dir,
        verbose=False
    )


def check_heading_referenced(
    spec_basename: str,
    heading_anchor: str,
    section_anchor: Optional[str],
    requirement_files: List[Path],
    *,
    file_cache: Optional[FileContentCache] = None,
    verbose: bool = False,
) -> bool:
    """Return True if the heading anchor is referenced in any requirement file."""
    escaped_basename = re.escape(spec_basename)
    escaped_anchor = re.escape(heading_anchor)

    patterns = []
    patterns.append(
        rf'\[[^\]]*{escaped_basename}#{escaped_anchor}[^\]]*\]'
        rf'\(\.\./tech_specs/{escaped_basename}(?:#{escaped_anchor})?\)'
    )
    if section_anchor:
        escaped_section_anchor = re.escape(section_anchor)
        patterns.append(
            rf'\[[^\]]*{escaped_basename}#{escaped_section_anchor}[^\]]*\]'
            rf'\(\.\./tech_specs/{escaped_basename}'
            rf'(?:#{escaped_section_anchor})?\)'
        )
    patterns.append(
        rf'\[[^\]]+\]\(\.\./tech_specs/{escaped_basename}#{escaped_anchor}\)'
    )
    if section_anchor:
        escaped_section_anchor = re.escape(section_anchor)
        patterns.append(
            rf'\[[^\]]+\]\(\.\./tech_specs/{escaped_basename}'
            rf'#{escaped_section_anchor}\)'
        )
    patterns.append(
        rf'\[[^\]]*{escaped_basename}#\d+-{escaped_anchor}[^\]]*\]'
        rf'\(\.\./tech_specs/{escaped_basename}'
        rf'(?:#\d+-{escaped_anchor})?\)'
    )

    if file_cache is None:
        file_cache = FileContentCache()

    for req_file in requirement_files:
        try:
            content = file_cache.get_content(req_file)
            for pattern in patterns:
                if re.search(pattern, content):
                    if verbose:
                        anchor_used = (
                            section_anchor
                            if section_anchor and section_anchor in content
                            else heading_anchor
                        )
                        print(f"    Found reference to #{anchor_used} in {req_file.name}")
                    return True
        except (IOError, OSError) as e:
            if verbose:
                print(f"  Warning: Could not read {req_file}: {e}", file=sys.stderr)
        except UnicodeDecodeError as e:
            if verbose:
                print(
                    f"  Warning: Could not decode {req_file} (encoding issue): {e}",
                    file=sys.stderr
                )
        except _SCRIPT_ERROR_EXCEPTIONS as e:
            if verbose:
                print(f"  Warning: Unexpected error reading {req_file}: {e}", file=sys.stderr)

    return False


_RE_TECH_SPEC_ANCHOR = re.compile(
    r'\.\./tech_specs/([^/]+\.md)#([^)\s]+)'
)

_RE_TECH_SPEC_HREF = re.compile(r'^\.\./tech_specs/([^/#\s]+\.md)(?:#([^)\s]+))?$')

_RE_REQUIREMENT_START = re.compile(r'^\s*-\s*(?:~~)?(REQ-[A-Z0-9_]+-\d{3})\b')


def get_requirement_spec_anchor_links(
    req_file: Path,
    file_cache: Optional[FileContentCache] = None,
) -> List[Tuple[str, str, int]]:
    """Extract (spec_basename, anchor, line_num) for links to tech_specs with #anchor."""
    result = []
    if file_cache is None:
        file_cache = FileContentCache()
    for _link_text, link_target, line_num in extract_links(req_file, file_cache):
        if '../tech_specs/' not in link_target or '#' not in link_target:
            continue
        match = _RE_TECH_SPEC_ANCHOR.search(link_target)
        if match:
            result.append((match.group(1), match.group(2), line_num))
    return result


def _build_requirement_line_map(
    req_lines: List[str],
) -> Tuple[List[Optional[str]], Dict[str, int]]:
    """
    Build a mapping from line number -> requirement id for requirement list items.

    Requirements are list items starting with "- REQ-...".
    Continuation lines belong to the most recent requirement until the next "- REQ-..." line.

    Returns:
        (line_to_req, req_start_lines) where:
        - line_to_req is a 1-indexed list (index 0 unused) of requirement ids or None
        - req_start_lines maps requirement id -> starting line number
    """
    line_to_req: List[Optional[str]] = [None] * (len(req_lines) + 1)
    req_start_lines: Dict[str, int] = {}
    current_req: Optional[str] = None

    for line_num, line in enumerate[str](req_lines, 1):
        match = _RE_REQUIREMENT_START.match(line)
        if match:
            current_req = match.group(1)
            req_start_lines.setdefault(current_req, line_num)
        line_to_req[line_num] = current_req

    return (line_to_req, req_start_lines)


def check_anchor_in_text_missing_href_anchor(
    requirement_files: List[Path],
    file_cache: Optional[FileContentCache] = None,
    *,
    workspace_root: Optional[Path] = None,
) -> List[ValidationIssue]:
    """
    Error when link text includes `file.md#anchor` but href is `../tech_specs/file.md`.

    This is misleading: the visible anchor implies clicking will navigate to the anchor.
    Without the fragment in the href, it will not.
    """
    issues: List[ValidationIssue] = []
    if file_cache is None:
        file_cache = FileContentCache()

    for req_file in requirement_files:
        try:
            req_lines = file_cache.get_lines(req_file)
        except (IOError, OSError, UnicodeDecodeError):
            # Other parts of the audit will report file read issues; don't duplicate here.
            continue

        line_to_req, _req_start_lines = _build_requirement_line_map(req_lines)
        for link_text, link_target, line_num in extract_links(req_file, file_cache):
            req_id = line_to_req[line_num] if line_num < len(line_to_req) else None

            match = _RE_TECH_SPEC_HREF.match(link_target.strip())
            if not match:
                continue
            spec_basename = match.group(1)
            href_anchor = match.group(2)

            # Only care about "href has no anchor".
            if href_anchor:
                continue

            # Detect anchor in link text that matches the same spec basename.
            text_match = re.search(
                rf'\b{re.escape(spec_basename)}#([A-Za-z0-9_\-]+)\b',
                link_text
            )
            if not text_match:
                continue
            anchor = text_match.group(1)

            try:
                rel_path = req_file.relative_to(workspace_root) if workspace_root else req_file
            except ValueError:
                rel_path = req_file

            prefix = f"{req_id}: " if req_id else ""
            extra_fields = {
                'spec_basename': spec_basename,
                'anchor': anchor,
            }
            if req_id:
                extra_fields['requirement_id'] = req_id

            issues.append(ValidationIssue.create(
                "anchor_in_text_missing_href_anchor",
                Path(rel_path),
                line_num,
                line_num,
                message=(
                    f"{prefix}link text includes '{spec_basename}#{anchor}' "
                    f"but href has no '#{anchor}'"
                ),
                severity="error",
                suggestion=f"Change href to '../tech_specs/{spec_basename}#{anchor}'",
                **extra_fields,
            ))

    return issues


def check_requirement_tech_spec_link_thresholds(
    requirement_files: List[Path],
    file_cache: Optional[FileContentCache] = None,
    *,
    warn_threshold: int = REQ_TECH_SPEC_LINK_WARN_THRESHOLD,
    error_threshold: int = REQ_TECH_SPEC_LINK_ERROR_THRESHOLD,
    workspace_root: Optional[Path] = None,
) -> List[ValidationIssue]:
    """
    Warn/error when a single requirement contains too many tech spec links.

    - warning: count >= warn_threshold
    - error: count >= error_threshold
    """
    issues: List[ValidationIssue] = []
    if file_cache is None:
        file_cache = FileContentCache()

    for req_file in requirement_files:
        try:
            req_lines = file_cache.get_lines(req_file)
        except (IOError, OSError, UnicodeDecodeError):
            continue

        line_to_req, req_start_lines = _build_requirement_line_map(req_lines)
        counts: Dict[str, int] = {}

        for _link_text, link_target, line_num in extract_links(req_file, file_cache):
            req_id = line_to_req[line_num] if line_num < len(line_to_req) else None
            if req_id is None:
                continue
            if not link_target.strip().startswith('../tech_specs/'):
                continue
            counts[req_id] = counts.get(req_id, 0) + 1

        for req_id, count in counts.items():
            if count < warn_threshold:
                continue
            severity = "error" if count >= error_threshold else "warning"
            start_line = req_start_lines.get(req_id, 1)
            try:
                rel_path = req_file.relative_to(workspace_root) if workspace_root else req_file
            except ValueError:
                rel_path = req_file

            issues.append(ValidationIssue.create(
                "too_many_tech_spec_links",
                Path(rel_path),
                start_line,
                start_line,
                message=(
                    f"{req_id}: has {count} tech spec link(s) "
                    f"(warning >= {warn_threshold}, error >= {error_threshold})"
                ),
                severity=severity,
                suggestion=(
                    "Split into multiple requirements or reduce anchors "
                    "to keep link density low"
                ),
                requirement_id=req_id,
                count=count,
            ))

    return issues


def count_requirement_references(
    spec_basename: str,
    requirements_dir: Path,
    file_cache: Optional[FileContentCache] = None,
    verbose: bool = False,
    target_paths: Optional[List[str]] = None,
) -> int:
    """Count how many requirements reference a spec using inline markdown links."""
    count = 0
    search_pattern = f"(../tech_specs/{spec_basename}"

    if verbose:
        print(f"  Searching for: {search_pattern}")

    req_files = get_requirement_files(requirements_dir, target_paths)

    if file_cache is None:
        file_cache = FileContentCache()

    for req_file in req_files:
        try:
            content = file_cache.get_content(req_file)
            occurrences = content.count(search_pattern)
            count += occurrences
            if verbose and occurrences > 0:
                print(f"    Found {occurrences} in {req_file}")
        except (IOError, OSError) as e:
            if verbose:
                print(f"  Warning: Could not read {req_file}: {e}", file=sys.stderr)
        except UnicodeDecodeError as e:
            if verbose:
                print(
                    f"  Warning: Could not decode {req_file} (encoding issue): {e}",
                    file=sys.stderr
                )
        except _SCRIPT_ERROR_EXCEPTIONS as e:
            if verbose:
                print(f"  Warning: Unexpected error reading {req_file}: {e}", file=sys.stderr)

    return count
