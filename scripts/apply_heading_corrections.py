#!/usr/bin/env python3
"""
Apply heading numbering corrections from validate_heading_numbering.py output.

This script reads the output from validate_heading_numbering.py and automatically
applies the corrections to the markdown files.

Usage:
    # Read from stdin
    python3 scripts/validate_heading_numbering.py --path file.md | \
        python3 scripts/apply_heading_corrections.py

    # Read from file
    python3 scripts/apply_heading_corrections.py --input tmp/heading_report.txt

    # Dry run (show what would be changed without modifying files)
    python3 scripts/apply_heading_corrections.py --input tmp/heading_report.txt --dry-run

Options:
    --input, -i FILE       Read corrections from FILE (default: stdin)
    --dry-run, -d          Show what would be changed without modifying files
    --verbose, -v          Show detailed progress information
    --help, -h             Show this help message
"""

import argparse
import re
import sys
from collections import defaultdict
from pathlib import Path
from typing import List, NamedTuple

from lib._validation_utils import format_issue_message


class _CorrectionData(NamedTuple):
    """Data for a single heading correction."""

    filepath: str
    line_num: int
    current_number: str
    corrected_number: str
    heading_text: str
    level: int


class Correction:
    """Represents a single heading correction."""

    def __init__(self, data: _CorrectionData):
        self.filepath = data.filepath
        self.line_num = data.line_num
        self.current_number = data.current_number
        self.corrected_number = data.corrected_number
        self.heading_text = data.heading_text
        self.level = data.level

    def __repr__(self):
        return (f"Correction(filepath={self.filepath!r}, line={self.line_num}, "
                f"current={self.current_number!r}, corrected={self.corrected_number!r})")


class CorrectionParser:
    """Parses output from validate_heading_numbering.py."""

    # Pattern to match: "Line 123: ## [1.2] -> [2.3] Title"
    CORRECTION_PATTERN = re.compile(
        r'^Line\s+(\d+):\s+(#{2,})\s+\[([^\]]+)\]\s+->\s+\[([^\]]+)\]\s+(.+)$'
    )

    # Pattern to match file path in section header:
    # "Sorted headings from first error (line X) in filepath:"
    FILE_PATTERN = re.compile(
        r'Sorted headings from first error \(line \d+\) in (.+):'
    )

    def __init__(self):
        self.corrections = []  # List of Correction objects
        self.current_file = None

    def parse(self, input_lines: List[str]) -> List['Correction']:
        """Parse input lines and extract corrections."""
        for line in input_lines:
            line = line.rstrip()

            # Check for file path in section header
            file_match = self.FILE_PATTERN.search(line)
            if file_match:
                self.current_file = file_match.group(1).strip()
                continue

            # Check for correction line
            match = self.CORRECTION_PATTERN.match(line)
            if match:
                line_num = int(match.group(1))
                heading_prefix = match.group(2)
                current_number = match.group(3)
                corrected_number = match.group(4)
                heading_text = match.group(5)

                # Determine heading level from prefix
                level = len(heading_prefix)

                if self.current_file:
                    correction = Correction(_CorrectionData(
                        self.current_file,
                        line_num,
                        current_number,
                        corrected_number,
                        heading_text,
                        level,
                    ))
                    self.corrections.append(correction)
                else:
                    warning_msg = format_issue_message(
                        "warning",
                        "Correction without file context",
                        "unknown",
                        message=line.strip(),
                        no_color=False,
                    )
                    print(warning_msg, file=sys.stderr)

        return self.corrections


class CorrectionApplier:
    """Applies corrections to markdown files."""

    # Pattern to match numbered headings: "## 1.2 Title" or "## 1 Title"
    NUMBERED_HEADING_PATTERN = re.compile(
        r'^(#{2,})\s+([0-9]+(?:\.[0-9]+)*)\.?\s+(.+)$'
    )

    def __init__(self, dry_run=False, verbose=False):
        self.dry_run = dry_run
        self.verbose = verbose
        self.applied_count = 0
        self.failed_count = 0

    def log(self, message: str) -> None:
        """Print message if verbose mode is enabled."""
        if self.verbose:
            print(message)

    def apply_corrections(self, corrections: List['Correction']) -> None:
        """Apply corrections grouped by file."""
        # Group corrections by file
        corrections_by_file = defaultdict(list)
        for correction in corrections:
            corrections_by_file[correction.filepath].append(correction)

        # Apply corrections for each file
        for filepath, file_corrections in corrections_by_file.items():
            self.apply_file_corrections(filepath, file_corrections)

    def apply_file_corrections(self, filepath, corrections):
        """Apply corrections to a single file."""
        filepath = Path(filepath)
        if not filepath.exists():
            error_msg = format_issue_message(
                "error",
                "File not found",
                str(filepath),
                no_color=False,
            )
            print(error_msg, file=sys.stderr)
            self.failed_count += len(corrections)
            return

        self.log(f"Processing {filepath} ({len(corrections)} corrections)...")

        # Read file
        try:
            with open(filepath, 'r', encoding='utf-8') as f:
                lines = f.readlines()
        except (OSError, ValueError) as e:
            print(f"Error reading {filepath}: {e}", file=sys.stderr)
            self.failed_count += len(corrections)
            return

        # Sort corrections by line number (descending) to avoid line number shifts
        corrections.sort(key=lambda c: c.line_num, reverse=True)

        # Track if file was modified
        file_modified = False

        # Apply corrections
        for correction in corrections:
            if correction.line_num < 1 or correction.line_num > len(lines):
                error_msg = format_issue_message(
                    "error",
                    "Line out of range",
                    str(filepath),
                    line_num=correction.line_num,
                    message=(
                        f"Line number {correction.line_num} is out of range "
                        f"(file has {len(lines)} lines)"
                    ),
                    no_color=False,
                )
                print(error_msg, file=sys.stderr)
                self.failed_count += 1
                continue

            line_idx = correction.line_num - 1  # Convert to 0-based index
            original_line = lines[line_idx]
            modified_line = self.apply_correction_to_line(
                original_line, correction
            )

            if modified_line != original_line:
                if self.dry_run:
                    print(f"Would change line {correction.line_num} in {filepath}:")
                    print(f"  Old: {original_line.rstrip()}")
                    print(f"  New: {modified_line.rstrip()}")
                else:
                    lines[line_idx] = modified_line
                    file_modified = True
                self.applied_count += 1
            else:
                warning_msg = format_issue_message(
                    "warning",
                    "Could not apply correction",
                    str(filepath),
                    line_num=correction.line_num,
                    message=f"Line: {original_line.rstrip()}",
                    no_color=False,
                )
                print(warning_msg, file=sys.stderr)
                self.failed_count += 1

        # Write file if modified
        if file_modified and not self.dry_run:
            try:
                with open(filepath, 'w', encoding='utf-8') as f:
                    f.writelines(lines)
                self.log(f"  Updated {filepath}")
            except (OSError, ValueError) as e:
                print(f"Error writing {filepath}: {e}", file=sys.stderr)
                self.failed_count += len(corrections)

    def apply_correction_to_line(self, line, correction):
        """Apply a single correction to a line."""
        # Match the heading pattern
        match = self.NUMBERED_HEADING_PATTERN.match(line.rstrip())
        if not match:
            return line  # Not a numbered heading, return unchanged

        heading_prefix = match.group(1)
        current_number_str = match.group(2)
        heading_text = match.group(3)

        # Check if this line matches the correction
        # Strip periods for comparison
        current_stripped = current_number_str.rstrip('.')
        correction_current_stripped = correction.current_number.rstrip('.')

        if current_stripped != correction_current_stripped:
            # Line doesn't match expected current number
            return line

        # Determine if we need to add a period
        # Check if the corrected number in the output had a period
        corrected_has_period = correction.corrected_number.endswith('.')
        corrected_number_stripped = correction.corrected_number.rstrip('.')

        # Build the new heading
        if corrected_has_period:
            new_line = f"{heading_prefix} {corrected_number_stripped}. {heading_text}\n"
        else:
            new_line = f"{heading_prefix} {corrected_number_stripped} {heading_text}\n"

        return new_line


def main():
    """Main entry point."""
    parser = argparse.ArgumentParser(
        description='Apply heading numbering corrections from validate_heading_numbering.py',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog=__doc__
    )
    parser.add_argument(
        '--input', '-i',
        type=str,
        help='Read corrections from FILE (default: stdin)'
    )
    parser.add_argument(
        '--dry-run', '-d',
        action='store_true',
        help='Show what would be changed without modifying files'
    )
    parser.add_argument(
        '--verbose', '-v',
        action='store_true',
        help='Show detailed progress information'
    )

    args = parser.parse_args()

    # Read input
    if args.input:
        try:
            with open(args.input, 'r', encoding='utf-8') as f:
                input_lines = f.readlines()
        except (OSError, ValueError) as e:
            print(f"Error reading input file {args.input}: {e}", file=sys.stderr)
            sys.exit(1)
    else:
        input_lines = sys.stdin.readlines()

    # Parse corrections
    parser_obj = CorrectionParser()
    corrections = parser_obj.parse(input_lines)

    if not corrections:
        print("No corrections found in input.", file=sys.stderr)
        sys.exit(0)

    if args.verbose:
        print(f"Found {len(corrections)} corrections to apply.")

    # Apply corrections
    applier = CorrectionApplier(dry_run=args.dry_run, verbose=args.verbose)
    applier.apply_corrections(corrections)

    # Print summary
    if args.dry_run:
        print(f"\nDry run complete: {applier.applied_count} corrections would be applied, "
              f"{applier.failed_count} failed.")
    else:
        print(f"\nApplied {applier.applied_count} corrections, "
              f"{applier.failed_count} failed.")

    if applier.failed_count > 0:
        sys.exit(1)
    sys.exit(0)


if __name__ == '__main__':
    main()
