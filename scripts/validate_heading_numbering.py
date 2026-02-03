#!/usr/bin/env python3
"""
Validate markdown heading numbering consistency and capitalization.

This script checks all markdown files for heading numbering consistency and capitalization.
For headings (H2 and beyond) that start with numbers, it validates:
- H2 headings have 1 number (e.g., "## 1 Title")
- H3 headings have 2 numbers (e.g., "### 1.1 Subtitle")
- H4 headings have 3 numbers (e.g., "#### 1.1.1 Subtitle")
- Numbers are sequential within each parent section
- Child heading numbers match their parent section number
- Headings are not overly-deeply nested (flags H6 and beyond)
- Headings follow Title Case capitalization (warnings with suggestions)
- Organizational headings with no content are flagged as warnings
- H1 headings warn on numbering and error on duplicates

Usage:
    python3 validate_heading_numbering.py [options]

Options:
    --verbose, -v           Show detailed progress information
    --output, -o FILE       Write detailed output to FILE
    --path, -p PATHS        Check only the specified file(s) or
                            directory(ies) (recursive). Can be a single
                            path or comma-separated list of paths
    --nocolor, --no-color   Disable colored output
    --help, -h              Show this help message

Examples:
    # Basic validation
    python3 scripts/validate_heading_numbering.py

    # Save output to file
    python3 scripts/validate_heading_numbering.py --output tmp/numbering_report.txt

    # Verbose output
    python3 scripts/validate_heading_numbering.py --verbose

    # Check specific file
    python3 scripts/validate_heading_numbering.py --path \\
        docs/tech_specs/api_file_management.md

    # Check specific directory
    python3 scripts/validate_heading_numbering.py --path \\
        docs/requirements

    # Check multiple paths
    python3 scripts/validate_heading_numbering.py --path \\
        docs/requirements,docs/tech_specs
"""

import argparse
import sys
from pathlib import Path
from collections import defaultdict
from typing import List

from lib._validation_utils import (
    OutputBuilder, parse_no_color_flag,
    get_workspace_root, parse_paths,
    ValidationIssue, find_markdown_files
)
from lib._validate_heading_numbering_title_case import to_title_case as _to_title_case
from lib._validate_heading_numbering_report import print_summary as _print_summary_report
from lib._validate_heading_numbering_models import (
    HeadingInfo,
    RE_HEADING_PATTERN,
    RE_NUMBERED_HEADING_PATTERN,
)
from lib._validate_heading_numbering_helpers import (
    validate_heading_structure as _validate_heading_structure,
    is_go_code_related_heading as _is_go_code_related_heading,
)
from lib.heading_numbering import (
    check_duplicate_headings as _check_duplicate_headings,
    check_excessive_numbering as _check_excessive_numbering,
    check_h2_period_consistency as _check_h2_period_consistency,
    check_heading_capitalization as _check_heading_capitalization,
    check_organizational_headings as _check_organizational_headings,
    check_single_word_headings as _check_single_word_headings,
)


class HeadingValidator:
    """Validates heading numbering in markdown files."""

    HEADING_PATTERN = RE_HEADING_PATTERN
    NUMBERED_HEADING_PATTERN = RE_NUMBERED_HEADING_PATTERN

    def __init__(self, verbose=False, repo_root=None, no_color=False):
        self.verbose = verbose
        self.repo_root = Path(repo_root).resolve() if repo_root else None
        self.no_color = no_color
        self.issues: List[ValidationIssue] = []
        self.h1_counts = defaultdict(int)
        self.h1_first_line = defaultdict(lambda: None)
        # Track all headings for each file (from build_heading_structure)
        self.all_headings = defaultdict(list)  # Maps filepath -> list of HeadingInfo
        # Track headings from first error onward for each file
        self.headings_from_first_error = defaultdict(list)
        self.first_error_line = defaultdict(lambda: None)

    def log(self, message):
        """Print message if verbose mode is enabled."""
        if self.verbose:
            print(message)

    def rel_path(self, filepath):
        """Convert absolute filepath to relative path from repo root."""
        if not self.repo_root:
            return filepath
        try:
            return str(Path(filepath).resolve().relative_to(self.repo_root))
        except ValueError:
            # If path is not under repo_root, return as-is
            return filepath

    def parse_heading_number(self, heading_text):
        """
        Parse heading number from heading text.
        Returns tuple of (number_list, title) or (None, heading_text) if not numbered.

        Examples:
            "1 Overview" => ([1], "Overview")
            "1.2 Details" => ([1, 2], "Details")
            "1.2.3 Subdetails" => ([1, 2, 3], "Subdetails")
            "Overview" => (None, "Overview")
        """
        match = self.NUMBERED_HEADING_PATTERN.match(heading_text)
        if not match:
            return None, heading_text

        number_str = match.group(1)
        title = match.group(2)

        # Parse number string into list of integers
        try:
            numbers = [int(n) for n in number_str.split('.')]
            return numbers, title
        except ValueError:
            return None, heading_text

    def get_heading_level(self, heading_prefix):
        """Get heading level from markdown prefix (##, ###, etc.)."""
        return len(heading_prefix)

    def build_corrected_full_line(self, heading_info):
        """
        Build the corrected full heading line with numbering and capitalization fixes.
        Preserves backticks in heading text (case inside backticks is not changed).

        Returns the full markdown heading line with corrections applied.
        """
        if not heading_info:
            return None

        # Get heading prefix (##, ###, etc.)
        level = heading_info.level
        prefix = '#' * level

        # Determine the corrected number
        corrected_number = heading_info.original_number
        if heading_info.corrected_number:
            corrected_number = heading_info.corrected_number

        # Determine the corrected heading text (capitalization; backticks preserved)
        corrected_text = heading_info.heading_text
        if heading_info.corrected_capitalization:
            corrected_text = heading_info.corrected_capitalization

        # Trim extra whitespace from the text
        corrected_text = corrected_text.strip()

        # If numbering is MISSING, just return the heading without numbering
        if corrected_number == "MISSING":
            return f"{prefix} {corrected_text}".strip()

        # Determine period usage for H2 headings
        # Use the heading's has_period flag, but if corrected_number suggests a change,
        # we might need to adjust. For now, preserve the original period usage unless
        # the heading itself indicates it should change.
        has_period = heading_info.has_period
        if heading_info.level == 2:
            # Check if first H2 has period to maintain consistency
            h2_headings = [h for h in self.all_headings.get(heading_info.file, []) if h.level == 2]
            if h2_headings:
                first_h2 = min(h2_headings, key=lambda h: h.line_num)
                has_period = first_h2.has_period

        # Build the corrected line
        if heading_info.level == 2 and has_period:
            return f"{prefix} {corrected_number}. {corrected_text}"
        return f"{prefix} {corrected_number} {corrected_text}"

    def check_capitalization(self, heading_text):
        """
        Check if heading text follows Title Case and return corrected version.

        Returns:
            tuple: (is_correct, corrected_text)
        """
        if not heading_text:
            return True, heading_text

        corrected = _to_title_case(heading_text)
        is_correct = heading_text == corrected

        return is_correct, corrected

    def _read_file_lines(self, filepath):
        """
        Read file lines with proper error handling.

        Returns:
            List of lines if successful, None if error occurred
        """
        try:
            with open(filepath, 'r', encoding='utf-8') as f:
                return f.readlines()
        except (IOError, OSError) as e:
            # File read errors - create ValidationIssue
            error = ValidationIssue.create(
                "file_read_error",
                Path(filepath),
                0,
                0,
                message=f"Could not read file: {e}",
                severity='error'
            )
            self.issues.append(error)
            self.log(f"  Error reading file: {e}")
            return None
        except UnicodeDecodeError as e:
            # Encoding errors - create ValidationIssue
            error = ValidationIssue.create(
                "file_encoding_error",
                Path(filepath),
                0,
                0,
                message=f"Could not decode file (encoding issue): {e}",
                severity='error'
            )
            self.issues.append(error)
            self.log(f"  Error decoding file: {e}")
            return None
        except (MemoryError, RuntimeError, BufferError) as e:
            # Unexpected errors - create ValidationIssue
            error = ValidationIssue.create(
                "unexpected_error",
                Path(filepath),
                0,
                0,
                message=f"Unexpected error reading file: {e}",
                severity='error'
            )
            self.issues.append(error)
            self.log(f"  Unexpected error reading file: {e}")
            return None

    def _process_heading_line(
        self, filepath, line_num, line, *, stripped_line, heading_stack, in_code_block
    ):
        """
        Process a single line to check if it's a heading and extract heading info.

        Returns:
            Tuple of (heading_info, updated_in_code_block, should_continue)
            heading_info is HeadingInfo if heading found, None otherwise
            should_continue is True if we should skip this line
        """
        # Check for code block boundaries
        if stripped_line.startswith('```'):
            # Toggle code block state
            return None, not in_code_block, True

        # Skip lines inside code blocks
        if in_code_block:
            return None, in_code_block, True

        # Check if line is a heading
        match = self.HEADING_PATTERN.match(stripped_line)
        if not match:
            return None, in_code_block, False

        # Check for leading whitespace (linting error - headings should start at column 0)
        if line != line.lstrip():
            msg = ("Heading has leading whitespace. "
                   "This is likely a linting error that should be fixed.")
            warning = ValidationIssue.create(
                "heading_leading_whitespace",
                Path(filepath),
                line_num,
                line_num,
                message=msg,
                severity='warning',
                heading=line.strip()
            )
            self.issues.append(warning)

        # Extract heading prefix (all # characters at start)
        heading_prefix = match.group(1)
        heading_text = match.group(2)
        level = self.get_heading_level(heading_prefix)

        # Check for overly-deeply nested headings (H6 and beyond)
        if level >= 6:
            msg = (f"H{level} heading is too deeply nested. "
                   "Consider restructuring the document to use "
                   "H2-H5 only.")
            warning = ValidationIssue.create(
                "heading_leading_whitespace",
                Path(filepath),
                line_num,
                line_num,
                message=msg,
                severity='warning',
                heading=line.strip()
            )
            self.issues.append(warning)
            return None, in_code_block, True

        # Handle H1 headings (title)
        if level == 1:
            self._check_h1_heading(filepath, line_num, heading_text, line.strip())
            return None, in_code_block, True

        # Parse heading number
        numbers, _title = self.parse_heading_number(heading_text)

        # If this heading is not numbered, create HeadingInfo for it
        if numbers is None:
            return self._create_unnumbered_heading(
                filepath, line_num, level, heading_text,
                line=line, heading_stack=heading_stack
            ), in_code_block, True

        # Process numbered heading
        return self._create_numbered_heading(
            filepath, line_num, level, heading_text,
            line=line, numbers=numbers, heading_stack=heading_stack
        ), in_code_block, False

    def _create_unnumbered_heading(
        self, filepath, line_num, level, heading_text, *, line, heading_stack
    ):
        """Create HeadingInfo for an unnumbered heading."""
        # Find parent heading from current stack
        parent_heading = None
        if level > 2:
            parent_level = level - 1
            if parent_level in heading_stack:
                parent_heading = heading_stack[parent_level]

        # Create HeadingInfo with empty numbers list and "MISSING" as original_number
        heading_info = HeadingInfo(
            filepath, line_num, level, [],
            heading_text=heading_text,
            full_line=line.strip(), parent=parent_heading, issue=None
        )
        heading_info.original_number = "MISSING"
        heading_info.has_period = False

        return heading_info

    def _check_h1_heading(self, filepath, line_num, heading_text, full_line):
        """Track H1 headings and validate numbering/duplicates."""
        self.h1_counts[filepath] += 1
        if self.h1_first_line[filepath] is None:
            self.h1_first_line[filepath] = line_num
        else:
            msg = "More than one H1 heading found. Only the first H1 heading is valid."
            error = ValidationIssue.create(
                "heading_multiple_h1",
                Path(filepath),
                line_num,
                line_num,
                message=msg,
                severity='error',
                heading=full_line
            )
            self.issues.append(error)
            if self.first_error_line[filepath] is None:
                self.first_error_line[filepath] = line_num

        numbers, _ = self.parse_heading_number(heading_text)
        if numbers is not None:
            msg = "H1 heading should not be numbered."
            warning = ValidationIssue.create(
                "heading_h1_numbering",
                Path(filepath),
                line_num,
                line_num,
                message=msg,
                severity='warning',
                heading=full_line
            )
            self.issues.append(warning)

    def _create_numbered_heading(
        self, filepath, line_num, level, heading_text, *, line, numbers, heading_stack
    ):
        """Create HeadingInfo for a numbered heading."""
        # Extract original number string from heading_text using regex
        original_number_str = None
        match = self.NUMBERED_HEADING_PATTERN.match(heading_text)
        if match:
            original_number_str = match.group(1)
        else:
            # Fallback: reconstruct from numbers (shouldn't happen if parse succeeded)
            original_number_str = '.'.join(map(str, numbers))

        # Extract title from heading_text (remove the number prefix if present)
        number_str = '.'.join(map(str, numbers))
        if heading_text.startswith(number_str):
            title = heading_text[len(number_str):].strip()
            # Remove leading period and space if present
            if title.startswith('.'):
                title = title[1:].strip()
        else:
            title = heading_text

        # Find parent heading from current stack
        parent_heading = None
        if level > 2:
            # H3 and beyond need a parent
            parent_level = level - 1
            if parent_level in heading_stack:
                parent_heading = heading_stack[parent_level]

        # Check if H2 heading has period after number
        has_period = False
        if level == 2:
            number_str = '.'.join(map(str, numbers))
            # Check if there's a period after the number in the original text
            if heading_text.startswith(number_str):
                remaining = heading_text[len(number_str):].strip()
                if remaining.startswith('.'):
                    has_period = True

        heading_info = HeadingInfo(
            filepath, line_num, level, numbers,
            heading_text=title,
            full_line=line.strip(), parent=parent_heading, issue=None
        )
        # Set original_number from the actual string extracted from the file
        heading_info.original_number = original_number_str
        heading_info.has_period = has_period  # Track period for H2 headings

        return heading_info

    def _validate_heading_structure(self, filepath, headings, unnumbered_headings):
        """Validate heading structure after parsing. Returns list of headings."""
        return _validate_heading_structure(
            filepath, headings, unnumbered_headings,
            issues=self.issues,
            first_error_line=self.first_error_line,
            log_fn=self.log
        )

    def build_heading_structure(self, filepath):
        """
        First pass: Build the complete heading structure.
        Returns list of HeadingInfo objects sorted by line number, or None if error.
        """
        lines = self._read_file_lines(filepath)
        if lines is None:
            return None

        headings = []
        heading_stack = {}  # Maps level -> HeadingInfo (current parent at that level)
        unnumbered_headings = []  # Track unnumbered headings to check later

        # Track code block state
        in_code_block = False

        for line_num, line in enumerate(lines, start=1):
            stripped_line = line.strip()

            heading_info, in_code_block, should_continue = self._process_heading_line(
                filepath, line_num, line,
                stripped_line=stripped_line,
                heading_stack=heading_stack,
                in_code_block=in_code_block
            )

            if should_continue:
                continue

            if heading_info is None:
                continue

            # Handle unnumbered headings
            if heading_info.original_number == "MISSING":
                headings.append(heading_info)
                unnumbered_headings.append((
                    line_num, heading_info.level, heading_info.heading_text,
                    heading_info.full_line, heading_info
                ))
                # Don't update heading_stack for unnumbered headings
                # as they shouldn't be parents for numbered headings
                continue

            # Handle numbered headings
            headings.append(heading_info)

            # Update heading stack - set this heading as the current parent at its level
            heading_stack[heading_info.level] = heading_info

            # Clear deeper levels when we move up in hierarchy
            levels_to_clear = [lvl for lvl in heading_stack if lvl > heading_info.level]
            for lvl in levels_to_clear:
                del heading_stack[lvl]

        # Validate heading structure
        return self._validate_heading_structure(filepath, headings, unnumbered_headings)

    def _apply_h2_corrected_numbers(self, h2_headings: list) -> None:
        """Set corrected_number on H2 headings (sorted by line_num)."""
        start_number = 0
        if h2_headings:
            first_h2_numbers = h2_headings[0].numbers
            start_number = 0 if (first_h2_numbers and not first_h2_numbers[0]) else 1
        h2_sequence = start_number - 1
        for heading in h2_headings:
            h2_sequence += 1
            heading.corrected_number = str(h2_sequence)

    def _apply_level_corrected_numbers(self, headings: list, level: int) -> None:
        """Set corrected_number on headings at level (H3–H6) using parent sequences."""
        level_headings = [
            h for h in headings
            if h.level == level and h.numbers and h.original_number != "MISSING"
        ]
        if not level_headings:
            return
        level_headings.sort(key=lambda h: h.line_num)
        parent_sequences = {}
        for heading in level_headings:
            if heading.parent and heading.parent.corrected_number:
                parent_corrected = heading.parent.corrected_number
                parent_sequences[parent_corrected] = parent_sequences.get(parent_corrected, 0) + 1
                seq = parent_sequences[parent_corrected]
                heading.corrected_number = f"{parent_corrected}.{seq}"
            else:
                heading.corrected_number = heading.original_number

    def calculate_corrected_numbers(self, _filepath, headings):
        """
        Calculate corrected numbers for ALL headings level by level.
        This happens BEFORE validation - we determine what the numbers SHOULD be,
        then compare with original numbers to find errors.
        """
        if not headings:
            return
        numbered_headings = [h for h in headings if h.numbers and h.original_number != "MISSING"]
        h2_headings = [h for h in numbered_headings if h.level == 2]
        if not numbered_headings:
            return
        h2_headings.sort(key=lambda h: h.line_num)
        self._apply_h2_corrected_numbers(h2_headings)
        for level in range(3, 7):
            self._apply_level_corrected_numbers(headings, level)

    def _find_prev_heading_same_level_parent(self, headings_by_line, heading):
        """Return previous heading at same level with same parent, or None."""
        if heading.level <= 2 or not heading.parent:
            return None
        for h in headings_by_line:
            if h.line_num >= heading.line_num:
                break
            if h.level == heading.level and h.parent == heading.parent:
                return h
        return None

    def _record_heading_issue(self, filepath, heading, error):
        """Append issue, set heading.issue, and update first_error_line."""
        self.issues.append(error)
        heading.issue = error
        if self.first_error_line[filepath] is None:
            self.first_error_line[filepath] = heading.line_num

    def _record_depth_mismatch_if_needed(
        self, filepath, heading, level, *, expected_depth, actual_depth
    ):
        """If depth mismatch, record error and return True; else return False."""
        if actual_depth == expected_depth:
            return False
        error = ValidationIssue.create(
            "heading_depth_mismatch",
            Path(filepath),
            heading.line_num,
            heading.line_num,
            message=(
                f"H{level} heading has {actual_depth} number(s), "
                f"expected {expected_depth}"
            ),
            severity='error',
            heading=heading.full_line,
            heading_info=heading
        )
        self._record_heading_issue(filepath, heading, error)
        return True

    def _record_no_parent_if_needed(self, filepath, heading, level):
        """If H3+ has no parent, record error and return True; else return False."""
        if level <= 2 or heading.parent is not None:
            return False
        error = ValidationIssue.create(
            "heading_no_parent_in_validation",
            Path(filepath),
            heading.line_num,
            heading.line_num,
            message=(
                f"H{level} heading has no parent. "
                "Please run a markdown linter to fix basic heading order."
            ),
            severity='error',
            heading=heading.full_line,
            heading_info=heading
        )
        self._record_heading_issue(filepath, heading, error)
        return True

    def validate_numbering(self, filepath, headings):
        """
        Second pass: Validate numbering by comparing original_number vs corrected_number.
        Also check depth and parent relationships.
        """
        if not headings:
            return

        # Sort by line number to process sequentially
        headings_by_line = sorted(headings, key=lambda h: h.line_num)

        for heading in headings_by_line:
            # Skip unnumbered headings - they're already handled in build_heading_structure
            if heading.original_number == "MISSING":
                continue

            level = heading.level
            numbers = heading.numbers
            expected_depth = level - 1
            actual_depth = len(numbers)

            if self._record_depth_mismatch_if_needed(
                filepath, heading, level,
                expected_depth=expected_depth, actual_depth=actual_depth
            ):
                continue
            if self._record_no_parent_if_needed(filepath, heading, level):
                continue

            if heading.original_number != heading.corrected_number:
                prev_heading = self._find_prev_heading_same_level_parent(
                    headings_by_line, heading
                )
                if prev_heading:
                    msg = (f"Non-sequential numbering: got '{heading.original_number}', "
                           f"expected '{heading.corrected_number}' "
                           f"(previous was '{prev_heading.original_number}')")
                else:
                    msg = (f"Non-sequential numbering: got '{heading.original_number}', "
                           f"expected '{heading.corrected_number}'")

                error = ValidationIssue.create(
                    "heading_non_sequential",
                    Path(filepath),
                    heading.line_num,
                    heading.line_num,
                    message=msg,
                    severity='error',
                    heading=heading.full_line,
                    heading_info=heading
                )
                self._record_heading_issue(filepath, heading, error)

    def check_excessive_numbering(self, filepath, headings):
        """Delegate to heading_numbering checks module."""
        _check_excessive_numbering(self.issues, filepath, headings)

    def check_single_word_headings(self, filepath, headings):
        """Delegate to heading_numbering checks module."""
        _check_single_word_headings(self.issues, filepath, headings)

    def check_duplicate_headings(self, filepath, headings):
        """Delegate to heading_numbering checks module."""
        _check_duplicate_headings(
            self.issues, self.first_error_line, filepath, headings
        )

    def is_go_code_related_heading(self, heading_text):
        """Return True if heading appears to reference a Go code element."""
        return _is_go_code_related_heading(heading_text)

    def check_heading_capitalization(self, filepath, headings):
        """Delegate to heading_numbering checks module."""
        _check_heading_capitalization(self.issues, filepath, headings)

    def check_organizational_headings(self, filepath, headings):
        """Delegate to heading_numbering checks module."""
        try:
            with open(filepath, 'r', encoding='utf-8') as f:
                content = f.read()
        except (IOError, OSError) as e:
            self.log(f"  Error reading file for organizational check: {e}")
            return
        except UnicodeDecodeError as e:
            self.log(f"  Error decoding file for organizational check (encoding issue): {e}")
            return
        except (MemoryError, RuntimeError, BufferError) as e:
            self.log(f"  Unexpected error reading file for organizational check: {e}")
            return
        _check_organizational_headings(
            self.issues,
            self.first_error_line,
            filepath,
            headings,
            content,
            log_fn=self.log,
        )

    def check_h2_period_consistency(self, filepath, headings):
        """Delegate to heading_numbering checks module."""
        _check_h2_period_consistency(self.issues, filepath, headings)

    def validate_file(self, filepath):
        """Validate heading numbering in a single markdown file."""
        rel_file = self.rel_path(filepath)
        self.log(f"Validating: {rel_file}")

        # First pass: Build heading structure
        headings = self.build_heading_structure(filepath)
        if headings is None:
            return

        # Store all headings for this file
        self.all_headings[filepath] = headings

        # Filter H2 headings once and check if first is numbered
        # Only perform numbering-related checks if the first H2 is numbered
        h2_headings = [h for h in headings if h.level == 2]
        has_numbered_first_h2 = False
        if h2_headings:
            first_h2 = min(h2_headings, key=lambda h: h.line_num)
            has_numbered_first_h2 = (
                first_h2.numbers and
                len(first_h2.numbers) > 0 and
                first_h2.original_number != "MISSING"
            )

        # Only perform numbering-related checks if the first H2 is numbered
        if has_numbered_first_h2:
            # Calculate corrected numbers for ALL headings
            self.calculate_corrected_numbers(filepath, headings)

            # Check for excessive numbering (H3+ with numbers > 20)
            self.check_excessive_numbering(filepath, headings)

            # Second pass: Validate numbering by comparing original vs corrected
            self.validate_numbering(filepath, headings)

            # Check H2 period consistency (only if no errors)
            # Check for errors in a single pass
            has_errors = any(h.issue for h in headings)
            if not has_errors:
                self.check_h2_period_consistency(filepath, headings)

        # Always perform these checks regardless of numbering:
        # Check for duplicate headings (all levels)
        self.check_duplicate_headings(filepath, headings)

        # Check for single-word headings (H4+)
        self.check_single_word_headings(filepath, headings)

        # Check capitalization for all headings
        self.check_heading_capitalization(filepath, headings)

        # Check for organizational headings with no content
        self.check_organizational_headings(filepath, headings)

        # Track headings with errors for output
        errored_headings = [h for h in headings if h.issue]
        if errored_headings:
            if self.first_error_line[filepath] is None:
                self.first_error_line[filepath] = min(h.line_num for h in errored_headings)
            self.headings_from_first_error[filepath].extend(errored_headings)

    def find_markdown_files(self, root_dir, target_paths=None):
        """Find all markdown files in the repository or target paths."""
        # Use shared utility function
        exclude_dirs = {'node_modules', 'vendor', 'tmp'}
        md_files = find_markdown_files(
            target_paths=target_paths,
            root_dir=Path(root_dir) if root_dir else None,
            exclude_dirs=exclude_dirs,
            verbose=self.verbose,
            return_strings=True
        )
        return md_files

    def validate_all(self, root_dir, output, target_paths=None):
        """Validate all markdown files in the repository or target paths."""
        markdown_files = self.find_markdown_files(root_dir, target_paths)

        if self.verbose:
            output.add_verbose_line(f"Scanning {len(markdown_files)} markdown files...")
            output.add_blank_line("working_verbose")

        for filepath in markdown_files:
            self.validate_file(filepath)

        # Check for errors in a single pass
        for issue in self.issues:
            if issue.matches(severity='error'):
                return False
        return True

    def print_summary(self, output_builder):
        """Print validation summary using OutputBuilder."""
        _print_summary_report(
            self.issues,
            self.all_headings,
            self.headings_from_first_error,
            self.first_error_line,
            rel_path_fn=self.rel_path,
            build_corrected_full_line_fn=self.build_corrected_full_line,
            no_color=self.no_color,
            output_builder=output_builder
        )


def show_help():
    """Show help message."""
    print(__doc__)


def main():
    """Main entry point."""
    parser = argparse.ArgumentParser(
        description='Validate markdown heading numbering consistency',
        add_help=False
    )
    parser.add_argument('--verbose', '-v', action='store_true',
                        help='Show detailed progress information')
    parser.add_argument('--output', '-o', metavar='FILE',
                        help='Write detailed output to FILE')
    parser.add_argument('--path', '-p', metavar='PATHS',
                        help='Check only the specified file(s) or '
                             'directory(ies) (comma-separated list)')
    parser.add_argument('--help', '-h', action='store_true',
                        help='Show this help message')
    parser.add_argument('--nocolor', '--no-color', action='store_true',
                        help='Disable colored output')
    parser.add_argument('--no-fail', action='store_true',
                        help='Exit with code 0 even if errors are found')

    args = parser.parse_args()

    if args.help:
        show_help()
        return 0

    # Parse comma-separated paths
    target_paths = parse_paths(args.path)

    # Find repository root
    repo_root = get_workspace_root()

    # Validate
    no_color = args.nocolor or parse_no_color_flag(sys.argv)

    # Create output builder (header streams immediately if verbose)
    output = OutputBuilder(
        "Markdown Heading Numbering Validation",
        "Validates heading numbering consistency",
        no_color=no_color,
        verbose=args.verbose,
        output_file=args.output
    )

    validator = HeadingValidator(verbose=args.verbose, repo_root=str(repo_root), no_color=no_color)
    validator.validate_all(str(repo_root), output, target_paths)
    validator.print_summary(output)
    output.print()

    return output.get_exit_code(args.no_fail)


if __name__ == '__main__':
    sys.exit(main())
