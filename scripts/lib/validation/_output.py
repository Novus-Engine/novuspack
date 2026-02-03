"""Output and issue formatting for validation scripts."""

import os
import sys
import re
from pathlib import Path
from typing import Optional, Dict

COLOR_GREEN = "32"
COLOR_RED = "31"
COLOR_YELLOW = "33"
COLOR_RESET = "0"

# Standard separator width
SEPARATOR_WIDTH = 80


def _section_has_header(lines, separator: str, title: str, max_look: int = 3) -> bool:
    """Return True if lines start with separator followed by title within max_look."""
    for i in range(min(max_look, len(lines))):
        if lines[i] == separator and i + 1 < len(lines) and lines[i + 1].strip() == title:
            return True
    return False


def _ensure_section_header(section_lines: list, title: str, separator: str) -> list:
    """Prepend separator and title header to section_lines if not already present."""
    if _section_has_header(section_lines, separator, title):
        return section_lines
    return [separator, title, separator] + section_lines


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
    severity, issue_type, file_path,
    *,
    line_num=None,
    message=None,
    suggestion=None,
    no_color=False,
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

    def __init__(self, title, description, no_color=False, *, output_file=None, verbose=False):
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
        self._has_warnings_only_message = False  # Warnings-only final message (no errors)
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
        # Clear other final messages if present (mutually exclusive)
        if self._has_failure_message or self._has_warnings_only_message:
            self._clear_final_messages()
        self._has_success_message = True
        self._has_failure_message = False
        self._has_warnings_only_message = False

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
        # Clear other final messages if present (mutually exclusive)
        if self._has_success_message or self._has_warnings_only_message:
            self._clear_final_messages()
        self._has_success_message = False
        self._has_failure_message = True
        self._has_warnings_only_message = False

        self._add_blank_to_section("final_message")
        # Automatically prepend ❌ if not already present
        if not message.startswith("❌ "):
            message = f"❌ {message}"
        colored_msg = colorize(message, COLOR_RED, self.no_color)
        self._add_to_section("final_message", colored_msg)
        self._add_blank_to_section("final_message")

    def add_warnings_only_message(
        self,
        message="Warnings detected. Review the warnings above.",
        verbose_hint=None,
    ):
        """
        Add warnings-only final message (no errors).

        Use when validation passed but there are warnings. Adds: 1 blank before,
        message with ⚠️ prefix (yellow), optional verbose hint line, 1 blank after.

        Args:
            message: Main message (⚠️ prepended if not present).
            verbose_hint: If set and not verbose, add a second line (e.g. run --verbose).
        """
        if self._has_success_message or self._has_failure_message:
            self._clear_final_messages()
        self._has_success_message = False
        self._has_failure_message = False
        self._has_warnings_only_message = True

        self._add_blank_to_section("final_message")
        if not message.startswith("⚠️ "):
            message = f"⚠️ {message}"
        colored_msg = colorize(message, COLOR_YELLOW, self.no_color)
        self._add_to_section("final_message", colored_msg)
        if verbose_hint and not self.verbose:
            self._add_to_section("final_message", verbose_hint)
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
        for line, (_line_type, verbose_only) in zip(section_lines, section_metadata):
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
            separator = "=" * SEPARATOR_WIDTH
            all_lines.extend(_ensure_section_header(warnings, "Warnings", separator))

        # 5. Errors (with header if any errors exist, but skip if only headers/separators)
        errors = self._filter_section(self.error_lines, self.error_metadata)
        non_header_errors = [
            line for line in errors
            if line.strip() and line.strip() != "Errors" and
            not (line == "=" * SEPARATOR_WIDTH)
        ]
        if non_header_errors:
            separator = "=" * SEPARATOR_WIDTH
            all_lines.extend(_ensure_section_header(errors, "Errors", separator))

        # 6. Final messages
        all_lines.extend(
            self._filter_section(self.final_message_lines, self.final_message_metadata)
        )

        return all_lines

    def _get_remaining_lines_to_print(self, all_lines: list) -> list:
        """
        Return the slice of all_lines that should be printed.

        If verbose and header already printed, returns only summary onward.
        Otherwise returns all_lines.
        """
        if not (self.verbose and self._header_printed):
            return all_lines
        separator = "=" * SEPARATOR_WIDTH
        filtered_header = self._filter_section(self.header_lines, self.header_metadata)
        filtered_verbose = self._filter_section(
            self.working_verbose_lines, self.working_verbose_metadata
        )
        skip_count = len(filtered_header) + len(filtered_verbose)
        summary_start = None
        for i, line in enumerate(all_lines):
            if (line == separator and i + 1 < len(all_lines) and
                    all_lines[i + 1].strip() == "Summary"):
                summary_start = i
                break
        if summary_start is not None:
            return all_lines[summary_start:]
        if 0 < skip_count < len(all_lines):
            return all_lines[skip_count:]
        if skip_count > 0:
            return []
        for i, line in enumerate(all_lines):
            if line == separator and i < 3:
                continue
            has_colon_with_digit = (
                ':' in line and
                any(c.isdigit() for c in line.split(':', 1)[-1].strip())
            )
            if has_colon_with_digit or (line == separator and i > 2):
                return all_lines[i:]
        return []

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
        remaining_lines = self._get_remaining_lines_to_print(all_lines)

        if remaining_lines:
            # Collapse consecutive blank lines (max 2 consecutive)
            collapsed_lines = []
            prev_was_blank = False
            for line in remaining_lines:
                is_blank = (not line or not line.strip())
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

    def has_warnings(self) -> bool:
        """
        Return True if warnings were recorded.
        """
        return self._has_warnings

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
        self._has_warnings_only_message = False
        self._errors_header_added = False
        self._warnings_header_added = False


class ValidationIssue:
    """
    Represents a validation issue found in markdown files.

    This is a shared class used across validation scripts for consistency.
    Issues are tracked as List[ValidationIssue] in validation functions.
    Use ValidationIssue.create(...) for R0917-friendly construction (≤5 positional).
    """

    @classmethod
    def create(
        cls,
        issue_type: str,
        file_path: Path,
        start_line: int,
        end_line: int,
        *,
        message: str,
        **kwargs
    ) -> "ValidationIssue":
        """Create a ValidationIssue (avoids too-many-positional-arguments)."""
        return cls(
            issue_type, file_path, start_line, end_line,
            message=message, **kwargs
        )

    def __init__(
        self,
        issue_type: str,
        file_path: Path,
        start_line: int,
        end_line: int,
        *,
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
            line_num=self.start_line,
            message=self.message,
            suggestion=self.suggestion,
            no_color=no_color,
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
