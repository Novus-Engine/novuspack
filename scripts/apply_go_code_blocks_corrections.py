#!/usr/bin/env python3
"""
Apply corrections from validate_go_code_blocks.py output.

This script reads the output from validate_go_code_blocks.py and automatically
applies fixable corrections:
- Missing comments: Adds a comment before function/method/type definitions
- Heading format: Fixes heading text to remove backticks and ensure proper format

Usage:
    # Read from stdin
    python3 scripts/validate_go_code_blocks.py --nocolor | \
        python3 scripts/apply_go_code_blocks_corrections.py

    # Read from file
    python3 scripts/apply_go_code_blocks_corrections.py --input tmp/go_code_blocks_report.txt

    # Dry run (show what would be changed without modifying files)
    python3 scripts/apply_go_code_blocks_corrections.py \
        --input tmp/go_code_blocks_report.txt --dry-run

Options:
    --input, -i FILE       Read corrections from FILE (default: stdin)
    --dry-run, -d          Show what would be changed without modifying files
    --verbose, -v          Show detailed progress information
    --repo-root DIR        Repository root directory (default: parent of script directory)
    --help, -h             Show this help message

The input file is typically produced by:
    make validate-go-code-blocks OUTPUT=tmp/go_code_blocks_report.txt NO_COLOR=1
"""

import argparse
import re
import sys
from collections import defaultdict
from pathlib import Path
from typing import List, Optional

scripts_dir = Path(__file__).parent
lib_dir = scripts_dir / "lib"

# Import shared utilities
if str(scripts_dir) not in sys.path:
    sys.path.insert(0, str(scripts_dir))

from lib._go_code_utils import (  # noqa: E402  # pylint: disable=wrong-import-position
    find_go_code_blocks,
    find_first_definition,
)


# Pattern to match error/warning lines with ANSI color codes stripped
# Format: ERROR/WARNING: Issue Type: file:line: message -> suggestion
ISSUE_PATTERN = re.compile(
    r'^(?:ERROR|WARNING):\s+([^:]+):\s+([^:]+):(\d+):\s+(.+?)(?:\s+->\s+(.+))?$'
)

# Pattern to match heading lines in markdown
HEADING_PATTERN = re.compile(r'^(#{1,6})\s+(.+)$')


def strip_ansi_codes(text: str) -> str:
    """Remove ANSI color codes from text."""
    ansi_escape = re.compile(r'\x1B(?:[@-Z\\-_]|\[[0-?]*[ -/]*[@-~])')
    return ansi_escape.sub('', text)


class Correction:
    """Represents a single correction to apply."""

    def __init__(  # pylint: disable=too-many-positional-arguments
        self,
        filepath: str,
        line_num: int,
        issue_type: str,
        message: str,
        suggestion: Optional[str] = None,
    ):
        self.filepath = filepath
        self.line_num = line_num
        self.issue_type = issue_type
        self.message = message
        self.suggestion = suggestion

    def __repr__(self):
        return (f"Correction(filepath={self.filepath!r}, line={self.line_num}, "
                f"type={self.issue_type!r})")


class CorrectionParser:
    """Parses output from validate_go_code_blocks.py."""

    def __init__(self, repo_root: Path):
        self.repo_root = repo_root
        self.corrections: List[Correction] = []

    def parse(self, input_lines: List[str]) -> List[Correction]:
        """Parse input lines and extract fixable corrections."""
        for line in input_lines:
            line = line.rstrip()
            if not line:
                continue

            # Strip ANSI color codes
            line = strip_ansi_codes(line)

            # Skip non-error/warning lines
            if not line.startswith(('ERROR:', 'WARNING:')):
                continue

            # Parse the issue line
            match = ISSUE_PATTERN.match(line)
            if not match:
                continue

            issue_type_name, file_path_str, line_str, message, suggestion = match.groups()

            # Normalize issue type (remove spaces, convert to lowercase with underscores)
            issue_type = issue_type_name.lower().replace(' ', '_')

            # Only process fixable issue types
            if issue_type not in ('missing_comment', 'heading_format'):
                continue

            try:
                line_num = int(line_str)
            except ValueError:
                continue

            # Resolve file path relative to repo root
            file_path = self.repo_root / file_path_str
            if not file_path.exists():
                continue

            correction = Correction(
                file_path,
                line_num,
                issue_type,
                message,
                suggestion
            )
            self.corrections.append(correction)

        return self.corrections


class CorrectionApplier:
    """Applies corrections to markdown files."""

    def __init__(self, repo_root: Path, dry_run: bool = False, verbose: bool = False):
        self.repo_root = repo_root
        self.dry_run = dry_run
        self.verbose = verbose
        self.applied_count = 0
        self.skipped_count = 0
        self.failed_count = 0

    def log(self, message: str) -> None:
        """Print message if verbose mode is enabled."""
        if self.verbose:
            print(message, file=sys.stderr)

    def apply_corrections(self, corrections: List[Correction]) -> None:
        """Apply corrections grouped by file."""
        # Group corrections by file
        corrections_by_file = defaultdict(list)
        for correction in corrections:
            corrections_by_file[correction.filepath].append(correction)

        # Apply corrections for each file
        for filepath, file_corrections in corrections_by_file.items():
            self.apply_file_corrections(filepath, file_corrections)

    def apply_file_corrections(self, filepath: str, corrections: List[Correction]) -> None:
        """Apply corrections to a single file."""
        filepath = Path(filepath)
        if not filepath.exists():
            print(f"Error: file not found: {filepath}", file=sys.stderr)
            self.failed_count += len(corrections)
            return

        self.log(f"Processing {filepath} ({len(corrections)} corrections)...")

        # Read file
        try:
            content = filepath.read_text(encoding='utf-8')
            lines = content.splitlines(keepends=True)
        except (OSError, ValueError) as e:
            print(f"Error reading {filepath}: {e}", file=sys.stderr)
            self.failed_count += len(corrections)
            return

        # Sort corrections by line number (descending) to avoid line number shifts
        corrections.sort(key=lambda c: c.line_num, reverse=True)

        file_modified = False

        # Apply corrections
        for correction in corrections:
            if correction.issue_type == 'missing_comment':
                modified = self.apply_missing_comment(
                    filepath, lines, correction
                )
                if modified:
                    file_modified = True
            elif correction.issue_type == 'heading_format':
                modified = self.apply_heading_format(
                    filepath, lines, correction
                )
                if modified:
                    file_modified = True

        # Write file if modified
        if file_modified and not self.dry_run:
            try:
                filepath.write_text(''.join(lines), encoding='utf-8')
                self.log(f"  Updated {filepath}")
            except (OSError, ValueError) as e:
                print(f"Error writing {filepath}: {e}", file=sys.stderr)
                self.failed_count += len(corrections)

    def apply_missing_comment(
        self, filepath: Path, lines: List[str], correction: Correction
    ) -> bool:
        """Apply missing comment correction."""
        if correction.line_num < 1 or correction.line_num > len(lines):
            print(
                f"Warning: line out of range in {filepath}: {correction.line_num}",
                file=sys.stderr,
            )
            self.failed_count += 1
            return False

        # correction.line_num is the line number in the markdown file where the definition is
        file_line_idx = correction.line_num - 1  # Convert to 0-based

        # Verify this line is inside a code block
        content = ''.join(lines)
        code_blocks = find_go_code_blocks(content)

        # Find which code block contains the target line
        # Note: start_line is the ```go line, end_line is the ``` line
        target_block = None
        for start_line, end_line, code_content in code_blocks:
            # The code content is between start_line+1 and end_line-1
            if start_line + 1 <= correction.line_num <= end_line - 1:
                target_block = (start_line, end_line, code_content)
                break

        if not target_block:
            print(
                f"Warning: could not find code block for line {correction.line_num} "
                f"in {filepath}",
                file=sys.stderr,
            )
            self.failed_count += 1
            return False

        start_line, end_line, code_content = target_block

        # Check if there's already a comment on the previous line
        if file_line_idx > 0:
            prev_line = lines[file_line_idx - 1].rstrip()
            if prev_line.startswith('//'):
                self.skipped_count += 1
                if self.verbose:
                    print(f"Already has comment: {filepath}:{correction.line_num}")
                return False

        # Get the definition to generate appropriate comment
        # Try function/method first, then type
        sig = find_first_definition(code_content, is_type=False)
        is_type_def = False
        if not sig:
            sig = find_first_definition(code_content, is_type=True)
            is_type_def = True
            if not sig:
                print(
                    f"Warning: could not parse definition at line {correction.line_num} "
                    f"in {filepath}",
                    file=sys.stderr,
                )
                self.failed_count += 1
                return False

        # Generate comment text based on definition type
        if is_type_def and hasattr(sig, 'normalized_type_name'):
            def_name = sig.normalized_type_name()
            kind = getattr(sig, 'kind', 'type')
            comment_text = f"// {def_name} represents a {kind}."
        else:
            def_name = sig.name
            if hasattr(sig, 'receiver') and sig.receiver:
                comment_text = f"// {sig.receiver}.{def_name} ..."
            else:
                comment_text = f"// {def_name} ..."

        # Get indentation from the definition line
        def_line = lines[file_line_idx]
        match = re.match(r'^(\s*)', def_line)
        indent = match.group(1) if match else ''

        comment_line = f"{indent}{comment_text}\n"

        if self.dry_run:
            print(f"Would add comment at line {correction.line_num} in {filepath}:")
            print(f"  {comment_line.rstrip()}")
        else:
            lines.insert(file_line_idx, comment_line)
            self.applied_count += 1
            if self.verbose:
                print(f"Added comment: {filepath}:{correction.line_num}")

        return True

    def apply_heading_format(
        self, filepath: Path, lines: List[str], correction: Correction
    ) -> bool:
        """Apply heading format correction."""
        if correction.line_num < 1 or correction.line_num > len(lines):
            print(
                f"Warning: line out of range in {filepath}: {correction.line_num}",
                file=sys.stderr,
            )
            self.failed_count += 1
            return False

        # Extract suggested heading from suggestion
        # Format: "Suggested: {heading_text}"
        if not correction.suggestion:
            print(
                f"Warning: no suggestion provided for heading format correction "
                f"at {filepath}:{correction.line_num}",
                file=sys.stderr,
            )
            self.failed_count += 1
            return False

        # Parse suggestion
        suggestion_match = re.match(r'^Suggested:\s*(.+)$', correction.suggestion)
        if not suggestion_match:
            print(
                f"Warning: could not parse suggestion: {correction.suggestion}",
                file=sys.stderr,
            )
            self.failed_count += 1
            return False

        suggested_heading_text = suggestion_match.group(1).strip()

        # Find the heading line
        # The correction line_num points to the code block, but we need to find
        # the heading that precedes it
        heading_line_idx = None
        for i in range(correction.line_num - 1, -1, -1):
            if i >= len(lines):
                continue
            line = lines[i].rstrip()
            match = HEADING_PATTERN.match(line)
            if match:
                heading_line_idx = i
                break

        if heading_line_idx is None:
            print(
                f"Warning: could not find heading before line {correction.line_num} "
                f"in {filepath}",
                file=sys.stderr,
            )
            self.failed_count += 1
            return False

        # Get current heading
        current_line = lines[heading_line_idx].rstrip()
        match = HEADING_PATTERN.match(current_line)
        if not match:
            print(
                f"Warning: line {heading_line_idx + 1} is not a heading in {filepath}",
                file=sys.stderr,
            )
            self.failed_count += 1
            return False

        heading_prefix, current_heading_text = match.groups()

        # Check if heading already matches suggestion
        if current_heading_text.strip() == suggested_heading_text:
            self.skipped_count += 1
            if self.verbose:
                print(f"Already correct: {filepath}:{heading_line_idx + 1}")
            return False

        # Create new heading line
        new_line = f"{heading_prefix} {suggested_heading_text}\n"

        if self.dry_run:
            print(f"Would change heading at line {heading_line_idx + 1} in {filepath}:")
            print(f"  Old: {current_line}")
            print(f"  New: {new_line.rstrip()}")
        else:
            lines[heading_line_idx] = new_line
            self.applied_count += 1
            if self.verbose:
                print(f"Updated heading: {filepath}:{heading_line_idx + 1}")

        return True


def _is_safe_repo_relative_path(repo_root: Path, candidate: Path) -> bool:
    """Check if candidate path is safely within repo_root."""
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


def main() -> int:
    """Main entry point."""
    parser = argparse.ArgumentParser(
        description=(
            "Apply suggested corrections from validate_go_code_blocks.py output"
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

    # Read input lines
    if args.input and args.input != "-":
        inp = Path(args.input)
        if not inp.exists():
            print(f"Error: input file not found: {inp}", file=sys.stderr)
            return 1
        try:
            input_lines = inp.read_text(encoding="utf-8").splitlines()
        except (OSError, ValueError) as e:
            print(f"Error: failed to read input file {inp}: {e}", file=sys.stderr)
            return 1
    else:
        input_lines = sys.stdin.read().splitlines()

    # Parse corrections
    parser_obj = CorrectionParser(args.repo_root)
    corrections = parser_obj.parse(input_lines)

    if not corrections:
        print("No fixable corrections found in input.", file=sys.stderr)
        return 0

    if args.verbose:
        print(f"Found {len(corrections)} fixable corrections to apply.", file=sys.stderr)

    # Apply corrections
    applier = CorrectionApplier(
        args.repo_root, dry_run=args.dry_run, verbose=args.verbose
    )
    applier.apply_corrections(corrections)

    # Print summary
    if args.dry_run:
        print(
            f"Dry run complete: {applier.applied_count} correction(s) would be applied, "
            f"{applier.skipped_count} skipped, {applier.failed_count} failed."
        )
    else:
        print(
            f"Applied {applier.applied_count} correction(s), "
            f"{applier.skipped_count} skipped, {applier.failed_count} failed."
        )

    return 1 if applier.failed_count > 0 else 0


if __name__ == '__main__':
    sys.exit(main())
