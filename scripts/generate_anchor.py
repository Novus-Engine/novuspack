#!/usr/bin/env python3
"""
Generate markdown anchors from markdown headings.

This script generates GitHub-style markdown anchors from markdown headings.
It supports generating anchors for:

- A specific heading line in a file (via --line)
- All headings in a file (via --file)

Options:
    --file, -f FILE    Markdown file to scan and print anchors for all headings
    --line, -l LINE    File + line reference in the format: path.md:224
    --help, -h         Show this help message

Examples:
    # Generate anchor for a specific heading line in a file
    python3 scripts/generate_anchor.py --line docs/tech_specs/api_core.md:42
    # Output: #some-heading-anchor

    # Print anchors for all headings in a file
    python3 scripts/generate_anchor.py --file docs/tech_specs/api_core.md
    # Output (one per heading):
    # docs/tech_specs/api_core.md:1: H1 Title => #title
"""

import argparse
import sys
from pathlib import Path
from typing import Tuple

from lib._validation_utils import extract_headings_from_file, generate_anchor_from_heading


def _parse_line_ref(line_ref: str) -> Tuple[Path, int]:
    """
    Parse a file:line reference (e.g., "docs/x.md:224").
    """
    if not line_ref or ":" not in line_ref:
        raise ValueError("LINE must be in the format: path.md:224")

    path_str, line_str = line_ref.rsplit(":", 1)
    path_str = path_str.strip()
    line_str = line_str.strip()

    if not path_str:
        raise ValueError("LINE must include a file path before ':'")

    try:
        line_num = int(line_str)
    except ValueError as e:
        raise ValueError("LINE must end with an integer line number") from e

    if line_num <= 0:
        raise ValueError("LINE number must be >= 1")

    return Path(path_str), line_num


def _generate_anchor_for_line(file_path: Path, line_num: int) -> str:
    if not file_path.exists():
        raise FileNotFoundError(f"File not found: {file_path}")

    headings = extract_headings_from_file(file_path, skip_code_blocks=True)
    for heading_text, _level, heading_line in headings:
        if heading_line == line_num:
            return generate_anchor_from_heading(heading_text, include_hash=True)

    raise ValueError(
        f"No markdown heading found at {file_path}:{line_num} "
        "(note: headings inside code blocks are ignored)"
    )


def _print_anchors_for_file(file_path: Path) -> None:
    if not file_path.exists():
        raise FileNotFoundError(f"File not found: {file_path}")

    headings = extract_headings_from_file(file_path, skip_code_blocks=True)
    for heading_text, level, line_num in headings:
        anchor = generate_anchor_from_heading(heading_text, include_hash=True)
        print(f"{file_path}:{line_num}: H{level} {heading_text} => {anchor}")


def main():
    """Main entry point for the script."""
    parser = argparse.ArgumentParser(
        description='Generate markdown anchors from markdown headings.',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
    Examples:
        # Generate anchor for a specific heading line in a file
        %(prog)s --line docs/tech_specs/api_core.md:42

        # Print anchors for all headings in a file
        %(prog)s --file docs/tech_specs/api_core.md
    """
    )
    group = parser.add_mutually_exclusive_group(required=True)
    group.add_argument(
        '--file', '-f',
        type=str,
        help='Markdown file to scan and print anchors for all headings'
    )
    group.add_argument(
        '--line', '-l',
        type=str,
        help='File + line reference in the format: path.md:224'
    )

    args = parser.parse_args()

    try:
        if args.line:
            file_path, line_num = _parse_line_ref(args.line)
            anchor = _generate_anchor_for_line(file_path, line_num)
            print(anchor)
        else:
            _print_anchors_for_file(Path(args.file))
    except (OSError, ValueError) as e:
        print(f"Error: {e}", file=sys.stderr)
        return 1

    return 0


if __name__ == '__main__':
    sys.exit(main())
