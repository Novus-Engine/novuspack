#!/usr/bin/env python3
"""
Shared utilities for validation scripts.

This module provides common functionality for all validation scripts,
including color support, standardized output formatting, and helper functions.
"""

import os
import sys
import re
import types
import importlib.util
from pathlib import Path
from typing import Optional, List, Set, Tuple, Dict
from dataclasses import dataclass

# Standard directory names used across validation scripts
DOCS_DIR = 'docs'
TECH_SPECS_DIR = 'tech_specs'
REQUIREMENTS_DIR = 'requirements'
FEATURES_DIR = 'features'

# Color codes (ANSI escape sequences)
COLOR_GREEN = "32"
COLOR_RED = "31"
COLOR_YELLOW = "33"
COLOR_RESET = "0"

# Standard separator width
SEPARATOR_WIDTH = 80

# Compiled regex patterns for performance (module level)
_RE_HEADING_PATTERN = re.compile(r'^(#{1,6})\s+(.+)$')
_RE_DECIMAL_PATTERN = re.compile(r'\d+\.\d+')
_RE_SENTENCE_END_PATTERN = re.compile(r'[.!?]+(?=\s+|$)')
_RE_HEADING_NUM_PATTERN = re.compile(r'^\d+(?:\.\d+)*$')


def supports_color(no_color_flag=False):
    """
    Check if colors should be used.

    Args:
        no_color_flag: If True, disable colors regardless of other conditions

    Returns:
        True if colors should be used, False otherwise
    """
    if no_color_flag or 'NO_COLOR' in os.environ:
        return False
    return sys.stdout.isatty()


def colorize(text, color_code, no_color_flag=False):
    """
    Apply color to text if colors are supported.

    Args:
        text: Text to colorize
        color_code: ANSI color code (e.g., "32" for green)
        no_color_flag: If True, disable colors

    Returns:
        Colorized text if colors are supported, otherwise original text
    """
    if supports_color(no_color_flag):
        return f"\033[{color_code}m{text}\033[0m"
    return text


def format_summary_line(label, value, label_width=25, value_width=6):
    """
    Format a summary line with aligned columns.

    Args:
        label: Label text (left-aligned)
        value: Value to display (right-aligned)
        label_width: Width for label column (default: 25)
        value_width: Width for value column (default: 6)

    Returns:
        Formatted string with aligned columns
    """
    return f"{label:<{label_width}} {value:>{value_width}}"


def calculate_label_width(labels, min_width=25, max_width=50):
    """
    Calculate optimal label width for a set of summary labels.

    Args:
        labels: List of label strings
        min_width: Minimum label width (default: 25)
        max_width: Maximum label width (default: 50)

    Returns:
        Optimal label width for formatting
    """
    if not labels:
        return min_width
    max_label_len = max(len(label) for label in labels)
    return min(max(max_label_len + 1, min_width), max_width)


def parse_no_color_flag(args):
    """
    Parse --nocolor or --no-color flag from command line arguments.

    Args:
        args: List of command line arguments (typically sys.argv)

    Returns:
        True if --nocolor or --no-color flag is present, False otherwise
    """
    return '--nocolor' in args or '--no-color' in args


def format_issue_message(
    severity, issue_type, file_path, line_num=None, message=None, suggestion=None, no_color=False
):
    """
    Format an error or warning message with consistent structure.

    Args:
        severity: Either "error" or "warning" (case-insensitive)
        issue_type: Short description of the issue type (e.g., "Heading without Req")
        file_path: Path to the file (will be converted to string)
        line_num: Optional line number
        message: Optional additional message/details
        suggestion: Optional suggestion for fixing the issue (formatted as " -> {suggestion}")
        no_color: If True, disable colors

    Returns:
        Formatted error or warning message string with color applied
    """
    severity_lower = severity.lower()
    if severity_lower not in ('error', 'warning'):
        raise ValueError(f"severity must be 'error' or 'warning', got '{severity}'")

    is_error = severity_lower == 'error'
    prefix = "ERROR" if is_error else "WARNING"
    color_code = COLOR_RED if is_error else COLOR_YELLOW

    file_str = str(file_path)
    if line_num is not None:
        location = f"{file_str}:{line_num}"
    else:
        location = file_str

    # Build the message parts
    parts = [f"{prefix}: {issue_type}: {location}"]

    if message:
        parts.append(message)

    # Build the full message
    if len(parts) > 1:
        issue_msg = ": ".join(parts)
    else:
        issue_msg = parts[0]

    # Add suggestion if present (without extra colon since it already has " -> ")
    if suggestion:
        issue_msg = f"{issue_msg} -> {suggestion}"

    return colorize(issue_msg, color_code, no_color)


class OutputBuilder:
    """
    Builder for consistent script output formatting.

    Handles headers, summaries, success messages, and spacing automatically.
    Tracks line types (error, warning, info) and supports verbose mode filtering.
    Automatically orders output sections in the correct sequence.
    """

    # Line type constants
    LINE_INFO = 'info'
    LINE_ERROR = 'error'
    LINE_WARNING = 'warning'
    LINE_VERBOSE = 'verbose'

    def __init__(self, title, description, no_color=False, output_file=None, verbose=False):
        """
        Initialize the output builder.

        Args:
            title: Script title for header
            description: Brief description for header
            no_color: If True, disable colors
            output_file: Optional file path to write output to
            verbose: If True, include verbose-only lines in output
        """
        # Separate sections for automatic ordering
        self.header_lines = []
        self.working_verbose_lines = []  # Working/progress verbose output
        self.summary_lines = []
        self.warning_lines = []
        self.error_lines = []
        self.final_message_lines = []  # Success messages, etc.

        # Metadata for each section
        self.header_metadata = []
        self.working_verbose_metadata = []
        self.summary_metadata = []
        self.warning_metadata = []
        self.error_metadata = []
        self.final_message_metadata = []

        self.no_color = no_color
        self.output_file = output_file
        self.verbose = verbose
        self._last_was_blank = {}  # Track blank lines per section
        self._has_warnings = False
        self._has_errors = False
        self._header_printed = False  # Track if header has been streamed
        self._streamed_lines = []  # Track streamed lines for file output
        self._streamed_header_count = 0  # Number of header lines streamed
        self._streamed_verbose_count = 0  # Number of working_verbose lines streamed
        self._summary_header_added = False  # Track if summary header has been added
        self._has_success_message = False  # Track if success message has been added
        self._has_failure_message = False  # Track if failure message has been added
        self._errors_header_added = False  # Track if errors header has been added
        self._warnings_header_added = False  # Track if warnings header has been added

        # Add header immediately (will stream if verbose)
        self.add_header(title, description)

    def _add_to_section(self, section, line, line_type=LINE_INFO, verbose_only=False):
        """
        Internal method to add a line to a specific section.

        Strips whitespace-only lines (they should be added via add_blank_line instead).
        """
        # Strip whitespace-only lines - they should be added via add_blank_line
        # Only process non-empty lines (empty lines should use add_blank_line)
        if line and line.strip():
            section_lines = getattr(self, f"{section}_lines")
            section_metadata = getattr(self, f"{section}_metadata")
            section_lines.append(line)
            section_metadata.append((line_type, verbose_only))
            self._last_was_blank[section] = False

    def _add_blank_to_section(self, section):
        """Internal method to add a blank line to a specific section."""
        if not self._last_was_blank.get(section, False):
            section_lines = getattr(self, f"{section}_lines")
            section_metadata = getattr(self, f"{section}_metadata")
            section_lines.append("")
            section_metadata.append((self.LINE_INFO, False))
            self._last_was_blank[section] = True

    def add_header(self, title, description):
        """
        Add script header with separators.

        If verbose=True, prints header immediately. Otherwise buffers it.

        Args:
            title: Script title
            description: Brief description
        """
        separator = "=" * SEPARATOR_WIDTH
        header_text = f"{title} - {description}"
        header_lines = [separator, header_text, separator]

        # Store in header section for final output
        for line in header_lines:
            self._add_to_section("header", line)

        # If verbose, print header immediately
        if self.verbose and not self._header_printed:
            for line in header_lines:
                print(line)
                if self.output_file:
                    self._streamed_lines.append(line)
            self._header_printed = True
            # Count header lines that will be in final output (after filtering)
            filtered_header = self._filter_section(self.header_lines, self.header_metadata)
            self._streamed_header_count = len(filtered_header)

    def add_summary_header(self):
        """Add summary section header."""
        if self._summary_header_added:
            return  # Already added, avoid duplicates
        separator = "=" * SEPARATOR_WIDTH
        self._add_to_section("summary", separator)
        self._add_to_section("summary", "Summary")
        self._add_to_section("summary", separator)
        self._summary_header_added = True

    def add_summary_section(self, items, label_width=None, value_width=6):
        """
        Add summary items with consistent formatting.

        Automatically adds summary header if:
        - There are summary items
        - AND (verbose is True OR there are warnings OR there are errors)
        - AND summary header hasn't been added yet

        Args:
            items: List of (label, value) tuples
            label_width: Optional label width (auto-calculated if None)
            value_width: Value column width (default: 6)
        """
        if not items:
            return

        # Automatically add summary header if conditions are met
        should_show_summary = self.verbose or self._has_warnings or self._has_errors
        if should_show_summary and not self._summary_header_added:
            self.add_summary_header()

        if label_width is None:
            labels = [item[0] for item in items]
            label_width = calculate_label_width(labels)

        for label, value in items:
            line = format_summary_line(label, value, label_width, value_width)
            self._add_to_section("summary", line)

    def add_success_message(self, message):
        """
        Add success message with proper spacing.

        Adds: 1 blank line before, message with ✅ prefix, 1 blank line after.

        Args:
            message: Success message text (✅ will be automatically prepended)
        """
        # Clear failure message if present (mutually exclusive)
        if self._has_failure_message:
            self._clear_final_messages()
        self._has_success_message = True
        self._has_failure_message = False

        self._add_blank_to_section("final_message")
        # Automatically prepend ✅ if not already present
        if not message.startswith("✅ "):
            message = f"✅ {message}"
        colored_msg = colorize(message, COLOR_GREEN, self.no_color)
        self._add_to_section("final_message", colored_msg)
        self._add_blank_to_section("final_message")

    def add_failure_message(self, message):
        """
        Add failure message with proper spacing.

        Adds: 1 blank line before, message with ❌ prefix, 1 blank line after.

        Args:
            message: Failure message text (❌ will be automatically prepended)
        """
        # Clear success message if present (mutually exclusive)
        if self._has_success_message:
            self._clear_final_messages()
        self._has_success_message = False
        self._has_failure_message = True

        self._add_blank_to_section("final_message")
        # Automatically prepend ❌ if not already present
        if not message.startswith("❌ "):
            message = f"❌ {message}"
        colored_msg = colorize(message, COLOR_RED, self.no_color)
        self._add_to_section("final_message", colored_msg)
        self._add_blank_to_section("final_message")

    def add_errors_header(self):
        """
        Add errors section header (standardized, like Summary).

        Note: This only adds the header. The header will only be displayed
        if there are actual error lines (not just the header itself).
        """
        if self._errors_header_added:
            return  # Already added, avoid duplicates
        self._has_errors = True
        separator = "=" * SEPARATOR_WIDTH
        self._add_to_section("error", separator, line_type=self.LINE_ERROR)
        self._add_to_section("error", "Errors", line_type=self.LINE_ERROR)
        self._add_to_section("error", separator, line_type=self.LINE_ERROR)
        self._errors_header_added = True

    def add_warnings_header(self):
        """Add warnings section header (standardized, like Summary)."""
        if self._warnings_header_added:
            return  # Already added, avoid duplicates
        self._has_warnings = True
        separator = "=" * SEPARATOR_WIDTH
        self._add_to_section("warning", separator, line_type=self.LINE_WARNING)
        self._add_to_section("warning", "Warnings", line_type=self.LINE_WARNING)
        self._add_to_section("warning", separator, line_type=self.LINE_WARNING)
        self._warnings_header_added = True

    def add_separator(self, section="summary"):
        """
        Add a separator line to the specified section.

        Args:
            section: Section to add separator to (default: "summary")
        """
        separator = "=" * SEPARATOR_WIDTH
        if section == "error":
            line_type = self.LINE_ERROR
        elif section == "warning":
            line_type = self.LINE_WARNING
        else:
            line_type = self.LINE_INFO
        self._add_to_section(section, separator, line_type=line_type)

    def add_line(self, line, line_type=LINE_INFO, verbose_only=False, section="summary"):
        """
        Add a raw line to output.

        Args:
            line: Line text to add
            line_type: Type of line ('info', 'error', 'warning', 'verbose')
            verbose_only: If True, only include this line when verbose=True
            section: Section to add to ('header', 'working_verbose', 'summary',
                'warning', 'error', 'final_message')
        """
        self._add_to_section(section, line, line_type=line_type, verbose_only=verbose_only)

    def add_error_line(self, line, verbose_only=False):
        """
        Add an error line to output.

        Automatically adds errors header if not already added.

        Args:
            line: Line text to add
            verbose_only: If True, only include this line when verbose=True
        """
        self._has_errors = True
        # Automatically add errors header if not already added
        if not self._errors_header_added:
            self.add_errors_header()
        self._add_to_section("error", line, line_type=self.LINE_ERROR, verbose_only=verbose_only)

    def add_warning_line(self, line, verbose_only=False):
        """
        Add a warning line to output.

        Args:
            line: Line text to add
            verbose_only: If True, only include this line when verbose=True
        """
        self._has_warnings = True
        # Automatically add warnings header if not already added
        if not self._warnings_header_added:
            self.add_warnings_header()
        self._add_to_section(
            "warning", line, line_type=self.LINE_WARNING, verbose_only=verbose_only
        )

    def add_verbose_line(self, line, line_type=LINE_INFO):
        """
        Add a verbose-only line to working verbose output section.

        If verbose=True, prints line immediately (after ensuring header is printed).
        Otherwise buffers it.

        Args:
            line: Line text to add
            line_type: Type of line ('info', 'error', 'warning')
        """
        # Store in working_verbose section for final output
        self._add_to_section("working_verbose", line, line_type=line_type, verbose_only=True)

        # If verbose, print immediately (after header if needed)
        if self.verbose:
            if not self._header_printed:
                # Header hasn't been printed yet, but we're trying to stream
                # This shouldn't happen if scripts call add_header first, but handle gracefully
                pass
            print(line)
            if self.output_file:
                self._streamed_lines.append(line)

    def add_blank_line(self, section="summary"):
        """
        Add a blank line to a specific section.

        If verbose=True and section is "working_verbose", prints immediately.
        Otherwise buffers it.

        Args:
            section: Section to add blank line to
        """
        self._add_blank_to_section(section)

        # If verbose and this is working_verbose, print immediately
        if self.verbose and section == "working_verbose" and self._header_printed:
            print("")
            if self.output_file:
                self._streamed_lines.append("")
            # Track that we've streamed this blank verbose line
            self._streamed_verbose_count += 1

    def _filter_section(self, section_lines, section_metadata):
        """
        Filter a section's lines based on verbose mode.

        Args:
            section_lines: List of lines in the section
            section_metadata: List of (line_type, verbose_only) tuples

        Returns:
            Filtered list of lines
        """
        filtered = []
        for line, (line_type, verbose_only) in zip(section_lines, section_metadata):
            if not verbose_only or self.verbose:
                filtered.append(line)
        return filtered

    def _get_ordered_sections(self):
        """
        Get all sections in the correct order with filtering applied.

        Header and summary are only included if verbose=True OR if there are warnings/errors.

        Returns:
            List of lines in correct order: header, working_verbose, summary,
            warning, error, final_message
        """
        all_lines = []

        # Check if we should show header and summary
        show_header_summary = self.verbose or self._has_warnings or self._has_errors

        # 1. Header (only if verbose or has warnings/errors)
        if show_header_summary:
            all_lines.extend(self._filter_section(self.header_lines, self.header_metadata))

        # 2. Working verbose output
        working_verbose = self._filter_section(
            self.working_verbose_lines, self.working_verbose_metadata
        )
        if working_verbose:
            all_lines.extend(working_verbose)

        # 3. Summary (only if verbose or has warnings/errors)
        if show_header_summary:
            all_lines.extend(self._filter_section(self.summary_lines, self.summary_metadata))

        # 4. Warnings (with header if any warnings exist)
        warnings = self._filter_section(self.warning_lines, self.warning_metadata)
        if warnings:
            # Check if warnings header already exists (from add_warnings_header)
            # Look for separator line followed by "Warnings"
            has_warnings_header = False
            separator = "=" * SEPARATOR_WIDTH
            for i in range(min(3, len(warnings))):
                if warnings[i] == separator:
                    if i + 1 < len(warnings) and warnings[i + 1].strip() == "Warnings":
                        has_warnings_header = True
                        break

            if not has_warnings_header:
                # Add standard warnings header
                all_lines.append(separator)
                all_lines.append("Warnings")
                all_lines.append(separator)
            all_lines.extend(warnings)

        # 5. Errors (with header if any errors exist, but skip if only headers/separators)
        errors = self._filter_section(self.error_lines, self.error_metadata)
        # Filter out empty sections (only headers/separators, no actual error content)
        non_header_errors = [
            line for line in errors
            if line.strip() and line.strip() != "Errors" and
            not (line == "=" * SEPARATOR_WIDTH)
        ]
        if non_header_errors:
            # Check if errors header already exists (from add_errors_header)
            # Look for separator line followed by "Errors"
            has_errors_header = False
            separator = "=" * SEPARATOR_WIDTH
            for i in range(min(3, len(errors))):
                if errors[i] == separator:
                    if i + 1 < len(errors) and errors[i + 1].strip() == "Errors":
                        has_errors_header = True
                        break

            if not has_errors_header:
                # Add standard errors header
                all_lines.append(separator)
                all_lines.append("Errors")
                all_lines.append(separator)
            all_lines.extend(errors)

        # 6. Final messages
        all_lines.extend(
            self._filter_section(self.final_message_lines, self.final_message_metadata)
        )

        return all_lines

    def print(self):
        """
        Print all lines to stdout and optionally to file.

        Outputs sections in correct order: header, working_verbose, summary,
        warning, error, final_message.
        Filters lines based on verbose mode before printing.

        If verbose=True, header and working_verbose have already been streamed,
        so only prints summary, warnings, errors, and final messages.

        After printing, clears all sections.
        """
        all_lines = self._get_ordered_sections()
        if not all_lines:
            return

        # If verbose, we've already streamed header and working_verbose
        # Only print the remaining sections (summary, warnings, errors, final_message)
        if self.verbose and self._header_printed:
            # Skip header and working_verbose (already streamed)
            # Count how many lines to skip: header + working_verbose (after filtering)
            filtered_header = self._filter_section(self.header_lines, self.header_metadata)
            filtered_verbose = self._filter_section(
                self.working_verbose_lines, self.working_verbose_metadata
            )
            skip_count = len(filtered_header) + len(filtered_verbose)

            # Find where summary starts (look for separator or summary items)
            separator = "=" * SEPARATOR_WIDTH
            summary_start = None

            # First try to find "Summary" header
            for i, line in enumerate(all_lines):
                if (line == separator and i + 1 < len(all_lines) and
                        all_lines[i + 1].strip() == "Summary"):
                    summary_start = i
                    break

            # If no summary header found, use skip_count to skip header + verbose
            if summary_start is not None:
                remaining_lines = all_lines[summary_start:]
            elif skip_count > 0:
                # Skip the header and working_verbose lines that were already streamed
                # Ensure we don't skip more than available lines
                if skip_count < len(all_lines):
                    remaining_lines = all_lines[skip_count:]
                else:
                    # If skip_count is >= len(all_lines), everything was already streamed
                    remaining_lines = []
            else:
                # Fallback: try to find first summary-like line (label: value pattern)
                # or first non-header/verbose line
                for i, line in enumerate(all_lines):
                    # Skip separator lines that are part of header (first 3 lines)
                    if line == separator and i < 3:
                        continue
                    # Look for summary items (label: value) or section headers
                    has_colon_with_digit = (
                        ':' in line and
                        any(c.isdigit() for c in line.split(':', 1)[-1].strip())
                    )
                    if has_colon_with_digit or (line == separator and i > 2):
                        summary_start = i
                        break
                if summary_start is not None:
                    remaining_lines = all_lines[summary_start:]
                else:
                    # Last resort: everything was already streamed
                    remaining_lines = []
        else:
            # Not verbose or header not printed yet - print everything
            # But ensure header is only included once
            remaining_lines = all_lines

        if remaining_lines:
            # Collapse consecutive blank lines (max 2 consecutive)
            collapsed_lines = []
            prev_was_blank = False
            for line in remaining_lines:
                is_blank = (line == "" or line.strip() == "")
                if is_blank:
                    # Only add blank line if previous line wasn't blank
                    if not prev_was_blank:
                        collapsed_lines.append("")
                    prev_was_blank = True
                else:
                    collapsed_lines.append(line)
                    prev_was_blank = False

            output_text = "\n".join(collapsed_lines)
            output_text += "\n"  # Final newline
            print(output_text, end="")

            if self.output_file:
                # Append collapsed lines to streamed lines for file output
                self._streamed_lines.extend(collapsed_lines)

        # Write to file if specified
        if self.output_file:
            try:
                with open(self.output_file, 'w', encoding='utf-8') as f:
                    # Combine streamed lines and remaining lines, remove color codes
                    import re
                    all_file_lines = self._streamed_lines + remaining_lines
                    if all_file_lines:
                        file_text = "\n".join(all_file_lines)
                        file_text += "\n"  # Final newline
                        file_text = re.sub(r'\033\[[0-9;]*m', '', file_text)
                        f.write(file_text)
            except IOError as e:
                print(
                    f"Error: Cannot write to output file {self.output_file}: {e}",
                    file=sys.stderr
                )

        # Clear all sections
        self.header_lines = []
        self.header_metadata = []
        self.working_verbose_lines = []
        self.working_verbose_metadata = []
        self.summary_lines = []
        self.summary_metadata = []
        self.warning_lines = []
        self.warning_metadata = []
        self.error_lines = []
        self.error_metadata = []
        self.final_message_lines = []
        self.final_message_metadata = []
        self._last_was_blank = {}
        self._header_printed = False
        self._streamed_lines = []

    def print_preview(self):
        """
        Print all current lines to stdout without clearing.

        Intended for showing output before interactive prompts.
        """
        all_lines = self._get_ordered_sections()
        if not all_lines:
            return
        output_text = "\n".join(all_lines)
        output_text += "\n"
        print(output_text, end="")

    def get_lines(self, filter_verbose=True):
        """
        Get all lines as a list in correct order (for custom processing).

        Args:
            filter_verbose: If True, filter based on verbose mode

        Returns:
            List of output lines in correct order
        """
        if filter_verbose:
            return self._get_ordered_sections()
        # If not filtering, combine all sections in order
        all_lines = []
        all_lines.extend(self.header_lines)
        all_lines.extend(self.working_verbose_lines)
        all_lines.extend(self.summary_lines)
        all_lines.extend(self.warning_lines)
        all_lines.extend(self.error_lines)
        all_lines.extend(self.final_message_lines)
        return all_lines

    def get_exit_code(self, no_fail=False):
        """
        Get the appropriate exit code based on errors found.

        Args:
            no_fail: If True, always return 0 (even if errors were found)

        Returns:
            0 if no errors found or no_fail is True, 1 if errors were found
        """
        if no_fail:
            return 0
        return 0 if not self._has_errors else 1

    def _clear_final_messages(self):
        """Clear final message section (used when switching between success/failure)."""
        self.final_message_lines = []
        self.final_message_metadata = []

    def clear(self):
        """Clear all accumulated lines from all sections."""
        self.header_lines = []
        self.header_metadata = []
        self.working_verbose_lines = []
        self.working_verbose_metadata = []
        self.summary_lines = []
        self.summary_metadata = []
        self.warning_lines = []
        self.warning_metadata = []
        self.error_lines = []
        self.error_metadata = []
        self.final_message_lines = []
        self.final_message_metadata = []
        self._last_was_blank = {}
        self._has_success_message = False
        self._has_failure_message = False
        self._errors_header_added = False
        self._warnings_header_added = False


def is_in_dot_directory(path: Path) -> bool:
    """
    Check if a path contains any directory starting with '.'.

    Args:
        path: Path object to check

    Returns:
        True if path contains any directory starting with '.' (except '.' itself), False otherwise
    """
    for part in path.parts:
        if part.startswith('.') and part != '.':
            return True
    return False


def find_markdown_files(
    target_paths: Optional[List[str]] = None,
    root_dir: Optional[Path] = None,
    default_dir: Optional[Path] = None,
    exclude_dirs: Optional[Set[str]] = None,
    verbose: bool = False,
    return_strings: bool = False
) -> List[Path]:
    """
    Find markdown files in the repository or target paths.

    Args:
        target_paths: Optional list of specific files or directories to check
        root_dir: Root directory to search from (when target_paths is None)
        default_dir: Default directory to search if target_paths is None and root_dir is None
        exclude_dirs: Set of directory names to exclude when scanning root_dir
        verbose: Whether to show detailed progress
        return_strings: If True, return list of strings instead of Path objects

    Returns:
        List of Path objects (or strings if return_strings=True) for markdown files found
    """
    md_files = []
    default_exclude_dirs = {
        'node_modules', 'vendor', 'tmp', '.git', '.venv', 'venv',
        '__pycache__', '.pytest_cache', 'dist', 'build',
        '.idea', '.vscode', '.cache'
    }
    if exclude_dirs is None:
        exclude_dirs = default_exclude_dirs

    if target_paths:
        for target_path in target_paths:
            target = Path(target_path)
            if not target.exists():
                if verbose:
                    print(
                        f"Warning: Target path does not exist: {target_path}",
                        file=sys.stderr
                    )
                continue

            if target.is_file():
                if target.suffix == '.md' and not is_in_dot_directory(target):
                    md_files.append(target)
                else:
                    if verbose:
                        print(
                            f"Warning: Target file is not a markdown file: {target_path}",
                            file=sys.stderr
                        )
            else:
                # Recursively find markdown files in target directory
                for md_file in target.rglob('*.md'):
                    if not is_in_dot_directory(md_file):
                        md_files.append(md_file)
    else:
        # Determine which directory to search
        search_dir = root_dir
        if search_dir is None:
            if default_dir is not None:
                search_dir = default_dir
            else:
                search_dir = Path('.')

        if not search_dir.exists():
            if verbose:
                print(f"Error: Search directory does not exist: {search_dir}", file=sys.stderr)
            return []

        # If default_dir is specified, only search that directory (non-recursive for glob)
        if default_dir is not None and root_dir is None:
            md_files = [
                f for f in sorted(search_dir.glob('*.md'))
                if not is_in_dot_directory(f)
            ]
        else:
            # Recursive search with exclusions
            for md_file in search_dir.rglob('*.md'):
                # Check if any excluded directory is in the path
                if any(excluded in md_file.parts for excluded in exclude_dirs):
                    continue
                # Also exclude dot directories
                if is_in_dot_directory(md_file):
                    continue
                md_files.append(md_file)

    if return_strings:
        return sorted([str(f) for f in md_files])
    return sorted(md_files)


def find_feature_files(
    target_paths: Optional[List[str]] = None,
    root_dir: Optional[Path] = None,
    default_dir: Optional[Path] = None,
    exclude_dirs: Optional[Set[str]] = None,
    verbose: bool = False,
    return_strings: bool = False
) -> List[Path]:
    """
    Find feature files (.feature) in the repository or target paths.

    Args:
        target_paths: Optional list of specific files or directories to check
        root_dir: Root directory to search from (when target_paths is None)
        default_dir: Default directory to search if target_paths is None and root_dir is None
        exclude_dirs: Set of directory names to exclude when scanning root_dir
        verbose: Whether to show detailed progress
        return_strings: If True, return list of strings instead of Path objects

    Returns:
        List of Path objects (or strings if return_strings=True) for feature files found
    """
    feature_files = []
    default_exclude_dirs = {
        'node_modules', 'vendor', 'tmp', '.git', '.venv', 'venv',
        '__pycache__', '.pytest_cache', 'dist', 'build',
        '.idea', '.vscode', '.cache'
    }
    if exclude_dirs is None:
        exclude_dirs = default_exclude_dirs

    if target_paths:
        for target_path in target_paths:
            target = Path(target_path)
            if not target.exists():
                if verbose:
                    print(
                        f"Warning: Target path does not exist: {target_path}",
                        file=sys.stderr
                    )
                continue

            if target.is_file():
                if target.suffix == '.feature' and not is_in_dot_directory(target):
                    feature_files.append(target)
                else:
                    if verbose:
                        print(
                            f"Warning: Target file is not a .feature file: {target_path}",
                            file=sys.stderr
                        )
            else:
                # Recursively find feature files in target directory
                for feature_file in target.rglob('*.feature'):
                    if not is_in_dot_directory(feature_file):
                        feature_files.append(feature_file)
    else:
        # Determine which directory to search
        search_dir = root_dir
        if search_dir is None:
            if default_dir is not None:
                search_dir = default_dir
            else:
                search_dir = Path('.')

        if not search_dir.exists():
            if verbose:
                print(f"Error: Search directory does not exist: {search_dir}", file=sys.stderr)
            return []

        # If default_dir is specified, search that directory recursively
        if default_dir is not None and root_dir is None:
            feature_files = [
                f for f in sorted(search_dir.rglob('*.feature'))
                if not is_in_dot_directory(f)
                and not any(excluded in f.parts for excluded in exclude_dirs)
            ]
        else:
            # Recursive search with exclusions
            for feature_file in search_dir.rglob('*.feature'):
                # Check if any excluded directory is in the path
                if any(excluded in feature_file.parts for excluded in exclude_dirs):
                    continue
                # Also exclude dot directories
                if is_in_dot_directory(feature_file):
                    continue
                feature_files.append(feature_file)

    if return_strings:
        return sorted([str(f) for f in feature_files])
    return sorted(feature_files)


def get_validation_exit_code(has_errors, no_fail=False):
    """
    Get the appropriate exit code for validation scripts.

    Args:
        has_errors: True if validation errors were found, False otherwise
        no_fail: If True, always return 0 (even if errors were found)

    Returns:
        0 if no errors found or no_fail is True, 1 if errors were found

    Note:
        This function is for scripts that track errors separately from OutputBuilder.
        Scripts using OutputBuilder should use output.get_exit_code(no_fail) instead.
    """
    if no_fail:
        return 0
    return 0 if not has_errors else 1


def get_workspace_root() -> Path:
    """
    Get the workspace root directory (parent of scripts directory).

    Returns:
        Path to workspace root
    """
    script_dir = Path(__file__).parent
    return script_dir.parent.parent


def import_module_with_fallback(module_name: str, script_dir: Path) -> types.ModuleType:
    """
    Import a module by name.

    Args:
        module_name: Name of module to import (e.g., '_validation_utils')
        script_dir: Directory containing the module file (unused)

    Returns:
        Imported module
    """
    return importlib.import_module(module_name)


def parse_paths(path_str: Optional[str]) -> Optional[List[str]]:
    """
    Parse comma-separated path string into list of paths.

    Args:
        path_str: Comma-separated string of paths, or None

    Returns:
        List of trimmed path strings, or None if path_str is None/empty
    """
    if not path_str:
        return None
    return [p.strip() for p in path_str.split(',') if p.strip()]


@dataclass(frozen=True)
class HeadingContext:
    """
    Context information about a markdown heading.

    Used to track heading information for code blocks and signatures.
    """
    heading_text: str  # The heading text (without # markers)
    heading_level: int  # Heading depth (1-6, where 1 is most general)
    heading_line: int  # Line number of the heading (1-indexed)
    file_path: Optional[str] = None  # Optional file path for context


@dataclass
class ProseSection:
    """
    Represents a prose-only section in a markdown document.

    This supports Overview blocks and prose subsections in index-style documents.
    """

    heading_str: str
    heading_num: Optional[str]
    heading_level: int
    heading_line: Optional[int]
    content: str
    parent_section: Optional["ProseSection"] = None
    child_sections: List["ProseSection"] = None
    has_code: bool = False
    code_blocks: List[Tuple[int, int, str]] = None
    file_path: Optional[str] = None
    lines: Optional[Tuple[int, int]] = None

    def __post_init__(self) -> None:
        if self.child_sections is None:
            self.child_sections = []
        if self.code_blocks is None:
            self.code_blocks = []
        if self.heading_num is not None:
            if not isinstance(self.heading_num, str):
                raise ValueError("heading_num must be a string or None")
            if not _RE_HEADING_NUM_PATTERN.match(self.heading_num):
                raise ValueError(
                    "heading_num must be a dotted number like '1', '2.4', or '3.5.6', got: %r"
                    % (self.heading_num,)
                )

    def path_label(self) -> str:
        parts: List[str] = []
        cur: Optional["ProseSection"] = self
        while cur is not None:
            parts.append(cur.heading_str)
            cur = cur.parent_section
        parts.reverse()
        return " > ".join(parts)


def extract_headings(content: str, skip_code_blocks: bool = True) -> List[Tuple[str, int, int]]:
    """
    Extract all headings from markdown content.

    Args:
        content: Markdown content as string
        skip_code_blocks: If True, skip headings inside code blocks

    Returns:
        List of tuples: (heading_text, heading_level, line_number)
        Lines are 1-indexed.
    """
    headings: List[Tuple[str, int, int]] = []
    lines = content.split('\n')
    in_code_block = False

    for i, line in enumerate(lines, 1):
        stripped_line = line.strip()

        if skip_code_blocks:
            # Check for code block boundaries
            if stripped_line.startswith('```'):
                in_code_block = not in_code_block
                continue

            # Skip lines inside code blocks
            if in_code_block:
                continue

        # Match markdown headings (# through ######)
        match = _RE_HEADING_PATTERN.match(stripped_line)
        if match:
            heading_level = len(match.group(1))
            heading_text = match.group(2).strip()
            headings.append((heading_text, heading_level, i))

    return headings


def extract_headings_from_file(
    file_path: Path, skip_code_blocks: bool = True, file_cache: Optional['FileContentCache'] = None
) -> List[Tuple[str, int, int]]:
    """
    Extract all headings from a markdown file.

    Args:
        file_path: Path to the markdown file
        skip_code_blocks: If True, skip headings inside code blocks
        file_cache: Optional FileContentCache instance to use for reading files

    Returns:
        List of tuples: (heading_text, heading_level, line_number)
        Lines are 1-indexed.
    """
    try:
        if file_cache:
            content = file_cache.get_content(file_path)
        else:
            with open(file_path, 'r', encoding='utf-8') as f:
                content = f.read()
        return extract_headings(content, skip_code_blocks=skip_code_blocks)
    except Exception as e:
        print(f"Error reading {file_path}: {e}", file=sys.stderr)
        return []


def extract_headings_with_anchors(
    file_path: Path, min_level: int = 1, max_level: int = 6,
    skip_code_blocks: bool = True, file_cache: Optional['FileContentCache'] = None
) -> Dict[str, Tuple[str, int, int]]:
    """
    Extract all headings from a markdown file and generate anchors.

    Args:
        file_path: Path to the markdown file
        min_level: Minimum heading level to include (1-6, default: 1)
        max_level: Maximum heading level to include (1-6, default: 6)
        skip_code_blocks: If True, skip headings inside code blocks
        file_cache: Optional FileContentCache instance to use for reading files

    Returns:
        Dictionary mapping anchor -> (heading_text, heading_level, line_number)
    """
    headings_dict = {}
    headings = extract_headings_from_file(
        file_path, skip_code_blocks=skip_code_blocks, file_cache=file_cache
    )
    for heading_text, heading_level, line_num in headings:
        if min_level <= heading_level <= max_level:
            anchor = generate_anchor_from_heading(heading_text, include_hash=False)
            headings_dict[anchor] = (heading_text, heading_level, line_num)
    return headings_dict


def extract_h2_plus_headings_with_sections(
    file_path: Path, skip_code_blocks: bool = True,
    file_cache: Optional['FileContentCache'] = None
) -> List[Tuple[int, str, int, str, Optional[str]]]:
    """
    Extract H2+ headings (## through ######) with anchors and section numbers.

    Args:
        file_path: Path to the markdown file
        skip_code_blocks: If True, skip headings inside code blocks
        file_cache: Optional FileContentCache instance to use for reading files

    Returns:
        List of tuples: (heading_level, heading_text, line_num, anchor, section_anchor)
        where:
        - heading_level is 2 for ##, 3 for ###, etc.
        - anchor is the plain anchor from heading text
        - section_anchor is the anchor with section number prefix (if section number exists)
    """
    headings_list = []
    headings = extract_headings_from_file(
        file_path, skip_code_blocks=skip_code_blocks, file_cache=file_cache
    )
    for heading_text, heading_level, line_num in headings:
        # Only include H2+ headings (level 2-6)
        if heading_level < 2:
            continue

        # Extract section number if present (e.g., "1.2.3 Heading" -> "1.2.3")
        section_match = re.match(r'^(\d+(?:\.\d+)*)\s+(.+)$', heading_text)
        section_anchor = None

        if section_match:
            # Heading has section number: "1.2.3 Heading Text"
            section_num = section_match.group(1)
            section_num_no_dots = section_num.replace('.', '')
            heading_text_without_section = section_match.group(2).strip()
            # Generate anchor from heading text without section number
            anchor = generate_anchor_from_heading(heading_text_without_section, include_hash=False)
            # Section anchor: section_num-anchor (e.g., "123-heading-text")
            section_anchor = f"{section_num_no_dots}-{anchor}"
        else:
            # Heading has no section number: just generate anchor from text
            anchor = generate_anchor_from_heading(heading_text, include_hash=False)

        headings_list.append((heading_level, heading_text, line_num, anchor, section_anchor))
    return headings_list


def extract_headings_with_section_numbers(
    file_path: Path, min_level: int = 2, max_level: int = 6,
    skip_code_blocks: bool = True, file_cache: Optional['FileContentCache'] = None
) -> Tuple[Set[str], Dict[str, Tuple[str, str]]]:
    """
    Parse markdown file to extract all heading anchors and section numbers.

    Args:
        file_path: Path to the markdown file
        min_level: Minimum heading level to include (1-6, default: 2 for H2+)
        max_level: Maximum heading level to include (1-6, default: 6)
        skip_code_blocks: If True, skip headings inside code blocks
        file_cache: Optional FileContentCache instance to use for reading files

    Returns:
        Tuple of (anchors set, sections dict where key is section_num and
        value is (heading_text, anchor))
    """
    anchors = set()
    sections = {}  # section_num -> (heading_text, anchor)

    if not file_path.exists():
        return anchors, sections

    headings = extract_headings_from_file(
        file_path, skip_code_blocks=skip_code_blocks, file_cache=file_cache
    )
    for heading_text, heading_level, line_num in headings:
        if min_level <= heading_level <= max_level:
            # Generate anchor from heading text (without '#' prefix)
            anchor = generate_anchor_from_heading(heading_text, include_hash=False)
            anchors.add(anchor)

            # Extract section number if present (e.g., "2.1 AddFile Package Method" -> "2.1")
            section_match = re.match(r'^(\d+(?:\.\d+)*)', heading_text)
            if section_match:
                section_num = section_match.group(1)
                sections[section_num] = (heading_text, anchor)

    return anchors, sections


def find_heading_before_line(
    content: str, line_num: int, prefer_deepest: bool = True
) -> Optional[HeadingContext]:
    """
    Find the heading context for a given line number in markdown content.

    Args:
        content: Markdown content as string
        line_num: Target line number (1-indexed)
        prefer_deepest: If True, return the most specific (deepest) heading.
                        If False, return the most recent heading.

    Returns:
        HeadingContext if a heading is found before the line, None otherwise.
    """
    lines = content.split('\n')

    if line_num < 1 or line_num > len(lines):
        return None

    # Find the most recent heading before this line

    if prefer_deepest:
        # Track heading stack to find the most specific heading
        heading_stack = []  # List of (level, text, line_num) tuples

        for i, line in enumerate(lines[:line_num], 1):
            match = _RE_HEADING_PATTERN.match(line.strip())
            if match:
                level = len(match.group(1))
                text = match.group(2).strip()
                # Remove headings at same or deeper level from stack
                heading_stack = [h for h in heading_stack if h[0] < level]
                # Add this heading
                heading_stack.append((level, text, i))

        # Get the most specific (deepest) heading
        if heading_stack:
            last_heading_level, last_heading, last_heading_line = heading_stack[-1]
            return HeadingContext(
                heading_text=last_heading,
                heading_level=last_heading_level,
                heading_line=last_heading_line
            )
    else:
        # Find the most recent heading (not necessarily deepest)
        for i in range(line_num - 1, -1, -1):
            if i < len(lines):
                match = _RE_HEADING_PATTERN.match(lines[i].strip())
                if match:
                    level = len(match.group(1))
                    text = match.group(2).strip()
                    return HeadingContext(
                        heading_text=text,
                        heading_level=level,
                        heading_line=i + 1
                    )

    return None


def find_heading_for_code_block(
    content: str, code_block_start_line: int
) -> Optional[str]:
    """
    Find the heading text that appears before a code block.

    This is a simpler version that just returns the heading text,
    useful for cases where only the text is needed.

    Args:
        content: Markdown content as string
        code_block_start_line: Line number where the code block starts (1-indexed)

    Returns:
        Heading text if found, None otherwise.
    """
    ctx = find_heading_before_line(content, code_block_start_line, prefer_deepest=False)
    return ctx.heading_text if ctx else None


def get_common_abbreviations() -> Set[str]:
    """
    Get comprehensive list of common abbreviations (case-insensitive matching).

    Returns:
        Set of abbreviations (all lowercase for case-insensitive matching)
    """
    return {
        # Titles
        'dr.', 'mr.', 'mrs.', 'ms.', 'prof.',
        # Academic degrees
        'ph.d.', 'm.d.', 'b.a.', 'm.a.', 'b.s.', 'm.s.',
        # Common abbreviations
        'etc.', 'i.e.', 'e.g.', 'vs.', 'a.m.', 'p.m.',
        # Business/location
        'inc.', 'ltd.', 'corp.', 'st.', 'ave.', 'blvd.',
    }


def contains_url(text: str) -> bool:
    """
    Check if text contains a URL.

    Detects:
    - http:// and https:// URLs
    - www. URLs (with word boundaries)
    - mailto: links

    Args:
        text: Text to check

    Returns:
        True if text contains a URL, False otherwise
    """
    url_patterns = [
        r'https?://',  # http:// or https://
        r'\bwww\.',    # www. with word boundary
        r'mailto:',    # mailto: links
    ]
    for pattern in url_patterns:
        if re.search(pattern, text, re.IGNORECASE):
            return True
    return False


def count_sentences(text: str) -> int:
    """
    Count sentences in text, handling edge cases.

    Splits on sentence-ending punctuation (., !, ?) followed by space/newline.
    Handles edge cases:
    - Abbreviations: using get_common_abbreviations() (normalize both to lowercase for comparison)
    - Decimals: \\d+\\.\\d+ pattern
    - URLs: using contains_url function
    - Ellipses: ... and Unicode ellipsis (…)
    - Hybrid approach: period + uppercase next char AND not in abbreviation list (case-insensitive)

    Args:
        text: Text to count sentences in

    Returns:
        Number of sentences (0 for empty/whitespace text, minimum 1 if text is non-empty)
    """
    if not text or not text.strip():
        return 0

    abbreviations = get_common_abbreviations()
    text_lower = text.lower()

    # Check if text contains URLs
    has_urls = contains_url(text)

    # Pattern for ellipses (not currently used but kept for potential future use)
    # ellipsis_pattern = re.compile(r'\.\.\.|…')

    # Split text into potential sentences using regex
    # Match sentence-ending punctuation followed by whitespace or end of string

    # Find all potential sentence endings
    matches = list(_RE_SENTENCE_END_PATTERN.finditer(text))
    if not matches:
        # No sentence-ending punctuation found, but text exists
        return 1 if text.strip() else 0

    sentences = []
    last_end = 0

    for match in matches:
        punct_pos = match.start()
        punct_end = match.end()
        # Adjust punct_end to include the whitespace if present
        if punct_end < len(text) and text[punct_end].isspace():
            # Skip whitespace
            while punct_end < len(text) and text[punct_end].isspace():
                punct_end += 1

        # Check if this period/exclamation/question mark is part of an ellipsis
        if punct_pos > 0 and punct_pos + 2 < len(text):
            if text[punct_pos - 1:punct_pos + 2] == '...' or text[punct_pos:punct_pos + 3] == '...':
                continue
        if punct_pos + 1 < len(text) and text[punct_pos:punct_pos + 2] == '..':
            continue

        # Check if this is part of a decimal number
        context_start = max(0, punct_pos - 10)
        context_end = min(len(text), punct_pos + 10)
        context = text[context_start:context_end]
        if _RE_DECIMAL_PATTERN.search(context):
            continue

        # Check if this is part of a URL
        if has_urls:
            url_context_start = max(0, punct_pos - 30)
            url_context_end = min(len(text), punct_pos + 30)
            url_context = text[url_context_start:url_context_end]
            if contains_url(url_context):
                continue

        # Check if this is an abbreviation
        # Look backwards to find the word before the punctuation
        word_start = punct_pos
        while word_start > 0 and (text[word_start - 1].isalnum() or text[word_start - 1] == '.'):
            word_start -= 1

        word_before = text_lower[word_start:punct_pos + 1]
        if word_before in abbreviations:
            continue

        # Check hybrid approach: if period and next char is uppercase, likely sentence end
        if text[punct_pos] == '.' and punct_end < len(text):
            # Find next non-whitespace character
            next_char_pos = punct_end
            while next_char_pos < len(text) and text[next_char_pos].isspace():
                next_char_pos += 1
            if next_char_pos < len(text):
                next_char = text[next_char_pos]
                # If uppercase and not an abbreviation, it's a sentence end
                if next_char.isupper() and word_before not in abbreviations:
                    # This is a sentence end
                    sentence = text[last_end:punct_end].strip()
                    if sentence:
                        sentences.append(sentence)
                    last_end = punct_end
                    continue

        # Regular sentence ending (followed by whitespace)
        sentence = text[last_end:punct_end].strip()
        if sentence:
            sentences.append(sentence)
        last_end = punct_end

    # Add remaining text as a sentence if any
    if last_end < len(text):
        remaining = text[last_end:].strip()
        if remaining:
            sentences.append(remaining)

    # Filter out empty sentences
    sentences = [s for s in sentences if s]

    return len(sentences) if sentences else (1 if text.strip() else 0)


def has_code_blocks(content: str, exclude_languages: Optional[Set[str]] = None) -> bool:
    """
    Check if content contains code blocks (any language, excluding specified).

    Extracts first word from language identifier by splitting on any non-alpha character.
    Examples: "go example" -> "go", "rust,no_run" -> "rust", "c++" -> "c"

    Args:
        content: Markdown content to check
        exclude_languages: Optional set of language identifiers to exclude
                         (e.g., {'text', 'markdown'})

    Returns:
        True if content contains code blocks (excluding specified languages)
    """
    lines = content.split('\n')
    in_code_block = False
    code_block_language = None

    for line in lines:
        stripped = line.strip()
        if stripped.startswith('```'):
            if in_code_block:
                # Closing code block
                in_code_block = False
                code_block_language = None
            else:
                # Opening code block
                in_code_block = True
                # Extract language identifier
                language_part = stripped[3:].strip()
                if language_part:
                    # Split on any non-alpha character to get first token
                    match = re.match(r'^([a-zA-Z]+)', language_part)
                    if match:
                        code_block_language = match.group(1).lower()
                    else:
                        code_block_language = None
                else:
                    code_block_language = None

                # Check if this language should be excluded
                if exclude_languages and code_block_language:
                    if code_block_language in exclude_languages:
                        # Skip this code block
                        continue

                # Found a code block that's not excluded
                return True

    return False


def build_heading_hierarchy(
    headings: List[Tuple[int, int, str]]  # (line_num, level, text)
) -> Dict[int, Optional[int]]:
    """
    Build parent-child relationship mapping for headings.

    Uses heading_stack approach similar to validate_heading_numbering.py.
    Each heading finds its most recent parent at the appropriate level.

    Args:
        headings: List of (line_num, level, text) tuples, sorted by line_num

    Returns:
        Dict mapping heading index (0-based) -> parent heading index (None if no parent).
        If H3+ appears before H2, it has no parent (None).
    """
    hierarchy = {}
    heading_stack = {}  # Maps level -> heading_index (current parent at that level)

    for idx, (line_num, level, text) in enumerate(headings):
        parent_index = None
        if level > 2:
            # H3 and beyond need a parent
            parent_level = level - 1
            if parent_level in heading_stack:
                parent_index = heading_stack[parent_level]

        hierarchy[idx] = parent_index

        # Update heading stack - set this heading as the current parent at its level
        heading_stack[level] = idx

        # Clear deeper levels when we move up in hierarchy
        levels_to_clear = [lvl for lvl in heading_stack.keys() if lvl > level]
        for lvl in levels_to_clear:
            del heading_stack[lvl]

    return hierarchy


def get_subheadings(
    heading_index: int,
    heading_level: int,
    all_headings: List[Tuple[int, int, str]],
    hierarchy: Dict[int, Optional[int]]
) -> List[int]:
    """
    Get all subheadings (all descendants at any level > heading_level) for a given heading.

    Args:
        heading_index: Index of the heading in all_headings list (0-based)
        heading_level: Level of the heading
        all_headings: List of (line_num, level, text) tuples
        hierarchy: Parent-child mapping from build_heading_hierarchy

    Returns:
        List of indices (0-based) for all subheadings (all descendants at any level > heading_level)
    """
    subheadings = []

    # Find all headings that are descendants of this heading
    # A heading is a descendant if it has this heading in its ancestor chain
    def is_descendant(child_idx: int) -> bool:
        current = child_idx
        while current is not None and current in hierarchy:
            parent = hierarchy[current]
            if parent == heading_index:
                return True
            current = parent
        return False

    for idx, (line_num, level, text) in enumerate(all_headings):
        if idx != heading_index and level > heading_level:
            if is_descendant(idx):
                subheadings.append(idx)

    return subheadings


def is_organizational_heading(
    content: str,
    heading_line: int,
    heading_level: int,
    all_headings: List[Tuple[int, int, str]],
    hierarchy: Dict[int, Optional[int]],
    max_prose_lines: int = 5
) -> dict:
    """
    Determine if a heading is purely organizational (grouping only).

    A heading is organizational if:
    - Has no code blocks (any language, except text/markdown)
    - Has max_prose_lines or fewer sentences
    - Only contains subheadings with no substantive content

    Args:
        content: Full markdown content
        heading_line: Line number of the heading (1-indexed)
        heading_level: Level of the heading (2-6)
        all_headings: List of (line_num, level, text) tuples
        hierarchy: Parent-child mapping from build_heading_hierarchy
        max_prose_lines: Maximum sentences before considered non-organizational

    Returns:
        Dict with:
        - is_organizational: bool - True if heading is organizational
        - is_empty: bool - True if heading has no content (0 sentences),
          False if it has minor informative content (1-5 sentences)
        - sentence_count: int - Number of sentences in the section
    """
    # Find the heading index
    heading_index = None
    for idx, (line_num, level, text) in enumerate(all_headings):
        if line_num == heading_line and level == heading_level:
            heading_index = idx
            break

    if heading_index is None:
        return {'is_organizational': False, 'is_empty': False, 'sentence_count': 0}

    # Find next heading (any level) to determine section boundaries
    next_heading_line = None
    for line_num, level, text in all_headings:
        if line_num > heading_line:
            next_heading_line = line_num
            break

    # Extract section content
    lines = content.split('\n')
    if next_heading_line:
        section_lines = lines[heading_line - 1:next_heading_line - 1]
    else:
        section_lines = lines[heading_line - 1:]

    section_content = '\n'.join(section_lines)

    # Extract prose (non-heading, non-code-block lines)
    prose_lines = []
    in_code_block = False
    for line in section_lines[1:]:  # Skip the heading line itself
        stripped = line.strip()
        if stripped.startswith('```'):
            in_code_block = not in_code_block
            continue
        if in_code_block:
            continue
        # Check if it's a heading
        if re.match(r'^#{1,6}\s+', stripped):
            continue
        if stripped:
            prose_lines.append(line)

    prose_text = '\n'.join(prose_lines)

    # Count sentences
    sentence_count = count_sentences(prose_text)

    # Check for code blocks (excluding text/markdown)
    if has_code_blocks(section_content, exclude_languages={'text', 'markdown'}):
        return {'is_organizational': False, 'is_empty': False, 'sentence_count': sentence_count}

    # Get subheadings
    subheadings = get_subheadings(heading_index, heading_level, all_headings, hierarchy)

    # Return organizational if: no code blocks AND sentence_count <= max_prose_lines AND
    # (sentence_count == 0 OR only subheadings)
    if sentence_count <= max_prose_lines:
        if sentence_count == 0 or (subheadings and len(prose_lines) <= max_prose_lines):
            return {
                'is_organizational': True,
                'is_empty': (sentence_count == 0),
                'sentence_count': sentence_count
            }

    return {'is_organizational': False, 'is_empty': False, 'sentence_count': sentence_count}


def generate_anchor_from_heading(heading: str, include_hash: bool = False) -> str:
    """
    Generate a GitHub-style markdown anchor from heading text.

    This function implements GitHub's markdown anchor generation algorithm:
    - Removes backticks but preserves their content (e.g., `` `code` `` -> `code`)
    - Converts to lowercase
    - Removes special characters except word characters, spaces, and hyphens
    - Collapses sequences of spaces and hyphens into a single hyphen
    - Strips leading and trailing hyphens

    Args:
        heading: The heading text (may contain markdown formatting like backticks)
        include_hash: If True, prefix the anchor with '#' (default: False)

    Returns:
        The generated anchor string (with '#' prefix if include_hash=True)

    Examples:
        >>> generate_anchor_from_heading("1.2.3 AddFile Package Method")
        '123-addfile-package-method'
        >>> generate_anchor_from_heading("File Management with `Package` type")
        'file-management-with-package-type'
        >>> generate_anchor_from_heading("Heading - With  Multiple   Spaces")
        'heading-with-multiple-spaces'
    """
    if not heading:
        return ""

    # Remove markdown formatting (backticks) but preserve their content
    # This matches GitHub's behavior: `` `code` `` becomes `code` in the anchor
    heading_clean = re.sub(r'`([^`]+)`', r'\1', heading)

    # Convert to lowercase
    heading_lower = heading_clean.lower()

    # Preserve " - " (space-hyphen-space) as "---" to match GitHub/markdownlint MD051
    _placeholder = 'TRPLDASH'
    heading_lower = heading_lower.replace(' - ', _placeholder)

    # Remove special characters except word characters, spaces, and hyphens
    anchor = re.sub(r'[^\w\s-]', '', heading_lower)

    # Collapse sequences of spaces and hyphens into a single hyphen
    anchor = re.sub(r'[-\s]+', '-', anchor)

    # Restore "---" for " - " so slug matches GitHub/markdownlint
    anchor = anchor.replace(_placeholder, '---')

    # Strip leading and trailing hyphens
    anchor = anchor.strip('-')

    # Add '#' prefix if requested
    if include_hash:
        return '#' + anchor if anchor else ""
    return anchor


def remove_backticks_keep_content(text: str) -> str:
    """
    Remove backticks from text but keep their contents.

    This removes the backtick characters but preserves the text that was
    enclosed in backticks. This is the standard behavior for both validation scripts.

    Args:
        text: Text that may contain backticks

    Returns:
        Text with backticks removed but content preserved

    Examples:
        "Heading with `code` example" => "Heading with code example"
        "`func()` and `var`" => "func() and var"
        "`code`" => "code"
        "No backticks here" => "No backticks here"
    """
    if not text:
        return text

    # Remove backticks but keep content
    # Pattern matches backtick, captures content, matches closing backtick
    result = re.sub(r'`([^`]*)`', r'\1', text)
    return result


def has_backticks(text: str) -> bool:
    """
    Check if text contains backticks.

    Args:
        text: Text to check for backticks

    Returns:
        True if text contains backticks, False otherwise

    Examples:
        "Heading with `code`" => True
        "Plain text heading" => False
        "" => False
        None => False
    """
    if not text:
        return False
    return '`' in text


def get_backticks_error_message() -> str:
    """
    Get the standard error message for backticks in headings.

    Returns:
        Standard error message string for backticks in headings
    """
    return ("Heading contains backticks. "
            "Headings should not contain backticks; use plain text instead.")


def is_safe_path(file_path: Path, repo_root: Path) -> bool:
    """
    Check if a path is safe (within repo and no traversal).

    Args:
        file_path: Path to check
        repo_root: Repository root directory

    Returns:
        True if path is safe (within repo root), False otherwise
    """
    try:
        # Resolve to absolute path and check it's within repo
        resolved = file_path.resolve()
        repo_resolved = repo_root.resolve()
        # Check that resolved path is within repo root
        return str(resolved).startswith(str(repo_resolved))
    except (OSError, ValueError):
        return False


def validate_file_name(filename: str) -> bool:
    """
    Validate that filename is safe (no path traversal, no separators).

    Args:
        filename: Filename to validate

    Returns:
        True if filename is safe, False otherwise
    """
    if not filename:
        return False
    # No path separators allowed
    if '/' in filename or '\\' in filename:
        return False
    # No parent directory references
    if '..' in filename:
        return False
    # No null bytes
    if '\x00' in filename:
        return False
    return True


def validate_spec_file_name(spec_file: str) -> bool:
    """
    Validate that spec file name is safe (no path traversal, no separators, .md extension).

    Args:
        spec_file: Spec file name to validate

    Returns:
        True if spec file name is safe, False otherwise
    """
    if not spec_file:
        return False
    # Must be a simple filename with .md extension
    # No path separators allowed
    if '/' in spec_file or '\\' in spec_file:
        return False
    # No parent directory references
    if '..' in spec_file:
        return False
    # Must end with .md
    if not spec_file.endswith('.md'):
        return False
    # Must be a valid filename (alphanumeric, underscore, hyphen, dot)
    if not re.match(r'^[a-zA-Z0-9_\-]+\.md$', spec_file):
        return False
    return True


def validate_anchor(anchor: str) -> bool:
    """
    Validate that anchor is safe (no path traversal, no separators).

    Args:
        anchor: Anchor string to validate

    Returns:
        True if anchor is safe, False otherwise
    """
    if not anchor:
        return True  # Empty anchor is OK
    # No path separators allowed
    if '/' in anchor or '\\' in anchor:
        return False
    # No parent directory references
    if '..' in anchor:
        return False
    # No null bytes
    if '\x00' in anchor:
        return False
    # Anchor should only contain alphanumeric, hyphens, underscores
    if not re.match(r'^[a-zA-Z0-9_\-]+$', anchor):
        return False
    return True


class FileContentCache:
    """
    Cache for file contents to avoid repeated reads.

    This class provides efficient caching of file contents to reduce I/O overhead
    when the same files are read multiple times during validation.
    """

    def __init__(self):
        """Initialize an empty cache."""
        self._cache: Dict[Path, str] = {}
        self._lines_cache: Dict[Path, List[str]] = {}

    def get_content(self, file_path: Path) -> str:
        """
        Get file content, using cache if available.

        Args:
            file_path: Path to the file to read

        Returns:
            File content as string

        Raises:
            IOError: If file cannot be read
        """
        if file_path not in self._cache:
            self._cache[file_path] = file_path.read_text(encoding='utf-8')
        return self._cache[file_path]

    def get_lines(self, file_path: Path) -> List[str]:
        """
        Get file content as list of lines, using cache if available.

        Args:
            file_path: Path to the file to read

        Returns:
            File content as list of lines (without newline characters)

        Raises:
            IOError: If file cannot be read
        """
        if file_path not in self._lines_cache:
            content = self.get_content(file_path)
            self._lines_cache[file_path] = content.split('\n')
        return self._lines_cache[file_path]

    def clear(self):
        """Clear all cached content."""
        self._cache.clear()
        self._lines_cache.clear()

    def has(self, file_path: Path) -> bool:
        """
        Check if file content is cached.

        Args:
            file_path: Path to check

        Returns:
            True if file is cached, False otherwise
        """
        return file_path in self._cache


class ValidationIssue:
    """
    Represents a validation issue found in markdown files.

    This is a shared class used across validation scripts for consistency.
    Issues are tracked as List[ValidationIssue] in validation functions.
    """

    def __init__(
        self,
        issue_type: str,
        file_path: Path,
        start_line: int,
        end_line: int,
        message: str,
        severity: str = "error",  # "error" or "warning"
        suggestion: Optional[str] = None,
        heading: Optional[str] = None,
        **kwargs
    ):
        """
        Create a ValidationIssue.

        Args:
            issue_type: Type of issue (e.g., 'missing_comment', 'heading_format')
            file_path: Path to the file (will be converted to string)
            start_line: Starting line number
            end_line: Ending line number
            message: Issue message
            severity: "error" or "warning" (default: "error")
            suggestion: Optional suggestion for fixing
            heading: Optional heading text
            **kwargs: Additional type-specific fields (e.g., def_name, def_kind, etc.)
        """
        self.issue_type = issue_type
        self.file = str(file_path)  # Convert Path to string
        self.start_line = start_line
        self.end_line = end_line
        self.message = message
        self.severity = severity.lower()  # Normalize to lowercase
        if self.severity not in ('error', 'warning'):
            raise ValueError(f"severity must be 'error' or 'warning', got '{severity}'")
        self.suggestion = suggestion
        self.heading = heading
        self.extra_fields = kwargs  # Store additional fields

    def to_dict(self) -> Dict:
        """Convert to dictionary for backward compatibility (JSON, reporting, etc.)."""
        result = {
            'type': self.issue_type,
            'file': self.file,
            'start_line': self.start_line,
            'end_line': self.end_line,
            'message': self.message,
            'severity': self.severity,
        }
        if self.suggestion:
            result['suggestion'] = self.suggestion
        if self.heading:
            result['heading'] = self.heading
        result.update(self.extra_fields)
        return result

    def format_message(self, no_color: bool = False) -> str:
        """Format issue message using format_issue_message utility."""
        return format_issue_message(
            self.severity,
            self.issue_type,
            self.file,
            self.start_line,
            self.message,
            self.suggestion,
            no_color
        )

    def matches(
        self,
        issue_type: Optional[str] = None,
        severity: Optional[str] = None
    ) -> bool:
        """
        Check if this issue matches the given filter criteria.

        Args:
            issue_type: Optional issue type to match (exact match)
            severity: Optional severity to match (exact match, case-insensitive)

        Returns:
            True if the issue matches all provided criteria, False otherwise.
            If no criteria are provided, returns True.
        """
        if issue_type is not None and self.issue_type != issue_type:
            return False
        if severity is not None and self.severity != severity.lower():
            return False
        return True

    def __repr__(self) -> str:
        """String representation for debugging."""
        return (
            f"ValidationIssue(type={self.issue_type!r}, file={self.file!r}, "
            f"line={self.start_line}, severity={self.severity!r})"
        )

    def __eq__(self, other) -> bool:
        """Equality comparison."""
        if not isinstance(other, ValidationIssue):
            return False
        return (
            self.issue_type == other.issue_type
            and self.file == other.file
            and self.start_line == other.start_line
            and self.end_line == other.end_line
            and self.message == other.message
            and self.severity == other.severity
        )
