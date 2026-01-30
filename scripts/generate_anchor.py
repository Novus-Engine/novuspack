#!/usr/bin/env python3
"""
Generate markdown anchor from heading text.

This script generates GitHub-style markdown anchors from heading text,
which can be used in markdown links to reference specific sections.

Usage:
    python3 scripts/generate_anchor.py "Heading Text"
    python3 scripts/generate_anchor.py --text "Heading Text"
    echo "Heading Text" | python3 scripts/generate_anchor.py

Options:
    --text, -t TEXT    Heading text to convert to anchor
    --help, -h         Show this help message

Examples:
    # From command line argument
    python3 scripts/generate_anchor.py "1.2.3 AddFile Package Method"
    # Output: #123-addfile-package-method

    # From stdin
    echo "1.2.3 AddFile Package Method" | python3 scripts/generate_anchor.py
    # Output: #123-addfile-package-method

    # Using --text switch
    python3 scripts/generate_anchor.py --text "File Management"
    # Output: #file-management

    # Via Makefile
    make generate-anchor TEXT="1.2.3 AddFile Package Method"
    # Output: #123-addfile-package-method

    # Headings with backticks (use single quotes to preserve backticks)
    python3 scripts/generate_anchor.py '1.2.3 AddFile with `code` example'
    # Output: #123-addfile-with-code-example

    # Headings with backticks via --text (single quotes recommended)
    python3 scripts/generate_anchor.py --text 'File Management with `Package` type'
    # Output: #file-management-with-package-type

    # Headings with backticks via Makefile (use single quotes)
    make generate-anchor TEXT='1.2.3 AddFile with `code` example'
    # Output: #123-addfile-with-code-example

Note: When headings contain backticks (e.g., `code`), use single quotes
      around the heading text to preserve the backticks. The script will
      automatically remove backticks and their contents when generating
      the anchor, as per GitHub markdown anchor generation rules.
"""

import sys
import argparse
from pathlib import Path

scripts_dir = Path(__file__).parent
lib_dir = scripts_dir / "lib"

# Import shared utilities
if str(scripts_dir) not in sys.path:
    sys.path.insert(0, str(scripts_dir))

from lib._validation_utils import (  # noqa: E402
    generate_anchor_from_heading,
)


def main():
    """Main entry point for the script."""
    parser = argparse.ArgumentParser(
        description='Generate markdown anchor from heading text.',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
    Examples:
        %(prog)s "1.2.3 AddFile Package Method"
        echo "File Management" | %(prog)s
        %(prog)s --text "File Management"

        # Headings with backticks (use single quotes to preserve backticks)
        %(prog)s '1.2.3 AddFile with `code` example'
        %(prog)s --text 'File Management with `Package` type'

    Note: When headings contain backticks (e.g., `code`), use single quotes
        around the heading text to preserve the backticks. The script will
        automatically remove backticks and their contents when generating
        the anchor, as per GitHub markdown anchor generation rules.
    """
    )
    parser.add_argument(
        '--text', '-t',
        type=str,
        help='Heading text to convert to anchor (use single quotes if heading contains backticks)'
    )
    parser.add_argument(
        'heading',
        nargs='?',
        type=str,
        help=(
            'Heading text to convert to anchor (alternative to --text; '
            'use single quotes if heading contains backticks)'
        )
    )

    args = parser.parse_args()

    # Get heading text from argument, --text switch, or stdin
    heading_text = None
    if args.text:
        heading_text = args.text
    elif args.heading:
        heading_text = args.heading
    else:
        # Read from stdin
        try:
            heading_text = sys.stdin.read().strip()
        except (EOFError, KeyboardInterrupt):
            parser.print_help()
            sys.exit(1)

    if not heading_text:
        parser.print_help()
        sys.exit(1)

    # Generate and output anchor (with '#' prefix for CLI output)
    anchor = generate_anchor_from_heading(heading_text, include_hash=True)
    print(anchor)

    return 0


if __name__ == '__main__':
    sys.exit(main())
