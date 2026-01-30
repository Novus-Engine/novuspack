#!/usr/bin/env python3
"""
Apply specification reference updates from validate_go_spec_references.py output.

This script reads the parseable output lines:
  path:line: old_ref -> new_ref
and for each line with a suggestion, replaces the Go file
  // Specification: ...
comment content with the suggested reference.

Usage:
    # Read from stdin
    python3 scripts/validate_go_spec_references.py | \
        python3 scripts/apply_go_spec_references.py

    # Read from file
    python3 scripts/apply_go_spec_references.py --input tmp/spec_refs_report.txt

    # Dry run (show what would be changed without modifying files)
    python3 scripts/apply_go_spec_references.py \
        --input tmp/spec_refs_report.txt --dry-run

Options:
    --input,   -i FILE       Read corrections from FILE (default: stdin)
    --dry-run, -d          Show what would be changed without modifying files
    --verbose, -v          Show detailed progress information
    --repo-root DIR        Repository root directory (default: parent of script directory)
    --help, -h             Show this help message

The input file is typically produced by:
    make -C api/go validate-go-spec-references OUTPUT=/tmp/spec_ref_errors.txt
"""

import argparse
import re
import sys
from pathlib import Path


CORRECTION_PATTERN = re.compile(r'^([^:]+):(\d+):\s+(.+?)\s+->\s+(.+)$')
SPEC_PATTERN = re.compile(r'^(\s*)//\s*Specification:\s*(.+)$')


def _is_safe_repo_relative_path(repo_root: Path, candidate: Path) -> bool:
    try:
        resolved_candidate = candidate.resolve()
        resolved_repo_root = repo_root.resolve()
    except (OSError, ValueError):
        return False
    try:
        resolved_candidate.relative_to(resolved_repo_root)
        return True
    except ValueError:
        return False


def _normalize_ref(ref: str) -> str:
    return re.sub(r'\s+', ' ', ref.strip())


def _read_input_lines(args: argparse.Namespace) -> tuple[list[str], int | None]:
    """Read input lines from file or stdin. Returns (lines, None) or ([], exit_code)."""
    if args.input and args.input != "-":
        inp = Path(args.input)
        if not inp.exists():
            print(f"Error: input file not found: {inp}", file=sys.stderr)
            return [], 1
        try:
            return inp.read_text(encoding="utf-8").splitlines(), None
        except (OSError, ValueError) as e:
            print(f"Error: failed to read input file {inp}: {e}", file=sys.stderr)
            return [], 1
    return sys.stdin.read().splitlines(), None


def _parse_corrections(
    input_lines: list[str],
    repo_root: Path,
    verbose: bool,
) -> list[tuple[Path, int, str, str]]:
    """Parse path:line: old -> new lines. Returns (file_path, line_num, old_ref, new_ref)."""
    corrections: list[tuple[Path, int, str, str]] = []
    for line in input_lines:
        line = line.rstrip()
        if not line or "(no suggestion)" in line or "(no valid suggestion" in line:
            continue
        match = CORRECTION_PATTERN.match(line)
        if not match:
            if verbose:
                print(f"Skipping unparseable line: {line}", file=sys.stderr)
            continue
        rel_path_str, line_str, old_ref, new_ref = match.groups()
        try:
            line_num = int(line_str)
        except ValueError:
            if verbose:
                print(f"Skipping invalid line number: {line}", file=sys.stderr)
            continue
        file_path = repo_root / rel_path_str
        if not _is_safe_repo_relative_path(repo_root, file_path):
            if verbose:
                print(f"Skipping unsafe path: {rel_path_str}", file=sys.stderr)
            continue
        corrections.append((file_path, line_num, old_ref.strip(), new_ref.strip()))
    return corrections


def _apply_to_one_file(
    fp: Path,
    file_corrections: list[tuple[int, str, str]],
    dry_run: bool,
    verbose: bool,
) -> tuple[int, int, int]:
    """Apply corrections to one file. Returns (applied, skipped, failed)."""
    if not fp.exists():
        print(f"Warning: file not found: {fp}", file=sys.stderr)
        return 0, 0, len(file_corrections)
    try:
        lines = fp.read_text(encoding="utf-8").splitlines(keepends=True)
    except (OSError, ValueError) as e:
        print(f"Error: failed to read {fp}: {e}", file=sys.stderr)
        return 0, 0, len(file_corrections)
    file_corrections = sorted(file_corrections, key=lambda t: t[0], reverse=True)
    applied, skipped, failed = 0, 0, 0
    file_modified = False
    for line_num, old_ref, new_ref in file_corrections:
        idx = line_num - 1
        if idx < 0 or idx >= len(lines):
            print(f"Warning: line out of range in {fp}: {line_num}", file=sys.stderr)
            failed += 1
            continue
        orig = lines[idx]
        orig_no_newline = orig.rstrip("\r\n")
        line_ending = orig[len(orig_no_newline):]
        match = SPEC_PATTERN.match(orig_no_newline)
        if not match:
            print(f"Warning: not a Specification comment at {fp}:{line_num}", file=sys.stderr)
            failed += 1
            continue
        indent, current_ref = match.groups()
        current_ref_norm = _normalize_ref(current_ref)
        old_ref_norm = _normalize_ref(old_ref)
        new_ref_norm = _normalize_ref(new_ref)
        if current_ref_norm == new_ref_norm:
            skipped += 1
            if verbose:
                print(f"Already updated: {fp}:{line_num}")
            continue
        if (
            current_ref_norm != old_ref_norm
            and current_ref_norm.lower() != old_ref_norm.lower()
            and old_ref_norm.lower() not in current_ref_norm.lower()
        ):
            print(f"Warning: old ref does not match at {fp}:{line_num}", file=sys.stderr)
            if verbose:
                print(f"  Expected old: {old_ref!r}", file=sys.stderr)
                print(f"  Found:        {current_ref.strip()!r}", file=sys.stderr)
            failed += 1
            continue
        updated_line = f"{indent}// Specification: {new_ref}{line_ending}"
        if dry_run:
            print(f"Would change {fp}:{line_num}")
            print(f"  Old: {orig.rstrip()}")
            print(f"  New: {updated_line.rstrip()}")
        else:
            lines[idx] = updated_line
            file_modified = True
        applied += 1
        if verbose:
            print(f"Updated: {fp}:{line_num}")
    if file_modified and not dry_run:
        try:
            fp.write_text("".join(lines), encoding="utf-8")
        except (OSError, ValueError) as e:
            print(f"Error: failed to write {fp}: {e}", file=sys.stderr)
            failed += applied
    return applied, skipped, failed


def main() -> int:
    parser = argparse.ArgumentParser(
        description=(
            "Apply suggested spec reference updates from "
            "validate_go_spec_references.py output"
        ),
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog=__doc__,
    )
    parser.add_argument(
        "--input", "-i",
        type=str,
        help="Read corrections from FILE (default: stdin)",
    )
    parser.add_argument(
        "--dry-run", "-d",
        action="store_true",
        help="Show what would be changed without modifying files",
    )
    parser.add_argument(
        "--verbose", "-v",
        action="store_true",
        help="Show detailed progress information",
    )
    parser.add_argument(
        "--repo-root",
        type=Path,
        default=Path(__file__).resolve().parent.parent,
        help="Repository root directory (default: parent of script directory)",
    )
    args = parser.parse_args()

    input_lines, exit_code = _read_input_lines(args)
    if exit_code is not None:
        return exit_code

    corrections = _parse_corrections(
        input_lines, args.repo_root, args.verbose
    )
    if not corrections:
        print("No corrections found in input.", file=sys.stderr)
        return 0

    by_file: dict[Path, list[tuple[int, str, str]]] = {}
    for fp, ln, old_ref, new_ref in corrections:
        by_file.setdefault(fp, []).append((ln, old_ref, new_ref))

    applied_count = 0
    skipped_count = 0
    failed_count = 0
    for fp in sorted(by_file.keys()):
        a, s, f = _apply_to_one_file(
            fp, by_file[fp], args.dry_run, args.verbose
        )
        applied_count += a
        skipped_count += s
        failed_count += f

    if args.dry_run:
        print(
            f"Dry run complete: {applied_count} correction(s) would be applied, "
            f"{skipped_count} skipped, {failed_count} failed."
        )
    else:
        print(
            f"Applied {applied_count} specification reference update(s), "
            f"{skipped_count} skipped, {failed_count} failed."
        )
    return 1 if failed_count > 0 else 0


if __name__ == '__main__':
    sys.exit(main())
