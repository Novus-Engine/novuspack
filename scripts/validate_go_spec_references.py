#!/usr/bin/env python3
"""
Validate Go file specification references against api_go_defs_index.md and spec files.

This script:
1. Parses all Go files for "Specification:" comments
2. Extracts file and section/anchor references
3. Validates that referenced files exist
4. Validates that section anchors exist in those files
5. Optionally checks against api_go_defs_index.md for method/type references

Usage:
    python3 scripts/validate_go_spec_references.py [options]

Options:
    --verbose, -v           Show detailed progress information
    --repo-root DIR         Repository root directory (default: parent of script directory)
    --check-index           Also validate against api_go_defs_index.md (not yet implemented)
    --output, -o FILE        Output file path for validation report (default: stdout)
    --help, -h               Show this help message
"""

import argparse
import re
import sys
from pathlib import Path
from typing import Dict, List, Optional, Set, Tuple

scripts_dir = Path(__file__).parent
lib_dir = scripts_dir / "lib"

# Import shared utilities
if str(scripts_dir) not in sys.path:
    sys.path.insert(0, str(scripts_dir))

from lib._validation_utils import (  # noqa: E402
    OutputBuilder, get_workspace_root, format_issue_message, parse_no_color_flag,
    ValidationIssue,
    is_safe_path, validate_spec_file_name, validate_anchor,
    extract_headings_with_section_numbers,
    FileContentCache, DOCS_DIR, TECH_SPECS_DIR,
    extract_h2_plus_headings_with_sections
)


class SpecReference:
    """Represents a specification reference from a Go file."""

    def __init__(self, file_path: Path, line_num: int, raw_ref: str):
        self.file_path = file_path
        self.line_num = line_num
        self.raw_ref = raw_ref
        self.spec_file: Optional[str] = None
        self.section: Optional[str] = None
        self.heading: Optional[str] = None
        self.is_valid_format = False
        self.function_name: Optional[str] = None  # Extracted from Go code context
        self.suggested_ref: Optional[str] = None  # Suggested correct reference
        self._parse()

    def _parse(self):
        """Parse the raw reference into components.

        Expected format: file_name.md: 4.5 Descriptive Heading
        Also handles: file_name.md#anchor or file_name.md Section X
        """
        # Remove relative path prefixes and validate no traversal remains
        ref = self.raw_ref.strip()
        # Remove any path traversal attempts
        ref = re.sub(r'^\.\./', '', ref)
        ref = re.sub(r'^\.\.\\', '', ref)
        # Check for any remaining traversal attempts
        if '..' in ref or '/' in ref or '\\' in ref:
            # If there's any path separator or traversal, reject it
            return

        # Check for the required format: file.md: section heading
        # Pattern: filename.md: section_number [.] heading_text
        # Optional period after section (e.g. "1. Core" or "1 Core") for markdown-style headings.
        pattern = r'^([a-zA-Z0-9_\-]+\.md):\s+(\d+(?:\.\d+)*)\.?\s+(.+)$'
        match = re.match(pattern, ref)

        if match:
            self.is_valid_format = True
            self.spec_file = match.group(1)
            self.section = match.group(2)
            self.heading = match.group(3).strip()
        else:
            # Try to parse anchor format: file.md#anchor
            anchor_pattern = r'^([a-zA-Z0-9_\-]+\.md)#([^#\s]+)$'
            anchor_match = re.match(anchor_pattern, ref)
            if anchor_match:
                self.spec_file = anchor_match.group(1)
                anchor = anchor_match.group(2)
                # Try to extract section number from anchor
                # Examples: "21-addfile" -> "2.1", "116-getmetadata" -> "1.1.6",
                #           "74-header" -> "7.4", "3-tag-management" -> "3"
                # Pattern: digits at start, optionally with dots
                section_match = re.match(r'^(\d+)(?:-|$)', anchor)
                if section_match:
                    digits_str = section_match.group(1)
                    # Convert "21" -> "2.1", "116" -> "1.1.6", "74" -> "7.4", "3" -> "3"
                    if len(digits_str) == 1:
                        self.section = digits_str
                    elif len(digits_str) == 2:
                        self.section = f"{digits_str[0]}.{digits_str[1]}"
                    elif len(digits_str) == 3:
                        self.section = f"{digits_str[0]}.{digits_str[1]}.{digits_str[2]}"
                    elif len(digits_str) == 4:
                        self.section = (
                            f"{digits_str[0]}.{digits_str[1]}."
                            f"{digits_str[2]}.{digits_str[3]}"
                        )
                return

            # Try "Section X" format
            section_pattern = r'^([a-zA-Z0-9_\-]+\.md)\s+Section\s+(\d+(?:\.\d+)*)(?:\s*-\s*(.+))?$'
            section_match = re.match(section_pattern, ref)
            if section_match:
                self.spec_file = section_match.group(1)
                self.section = section_match.group(2)
                if section_match.group(3):
                    self.heading = section_match.group(3).strip()
                return

            # Invalid format - try to extract what we can for error reporting
            if ':' in ref:
                parts = ref.split(':', 1)
                self.spec_file = parts[0].strip()
            elif ref.endswith('.md'):
                self.spec_file = ref.strip()
            else:
                # Try to extract filename if it looks like a file reference
                md_match = re.search(r'([a-zA-Z0-9_\-]+\.md)', ref)
                if md_match:
                    self.spec_file = md_match.group(1)

    def __repr__(self):
        if self.is_valid_format:
            return f"{self.spec_file}: {self.section} {self.heading}"
        return self.raw_ref


class SpecValidator:
    """Validates specification references."""

    def __init__(self, repo_root: Path):
        self.repo_root = repo_root
        self.docs_dir = repo_root / DOCS_DIR / TECH_SPECS_DIR
        self.api_go_dir = repo_root / "api" / "go"
        self.index_file = self.docs_dir / "api_go_defs_index.md"

        # File content cache to avoid repeated reads
        self.file_cache = FileContentCache()

        # Cache of parsed spec files (file -> set of anchors)
        self.spec_anchors: Dict[str, Set[str]] = {}
        # file -> {section_num: (heading_text, anchor)}
        self.spec_sections: Dict[
            str, Dict[str, Tuple[str, str]]
        ] = {}

        # Load index file
        self.index_entries: Dict[str, str] = {}  # method/type -> spec_file
        self.index_link_texts: Dict[str, str] = {}  # method/type -> link text for context
        self.index_anchors: Dict[str, str] = {}  # method/type -> anchor (e.g., "11-hashtype-type")
        self._load_index()

    def _is_section_0_or_cross_reference(self, section_num: str, heading_text: str = "") -> bool:
        """Check if a section is section 0 or a cross-reference section (not source of truth)."""
        # Section 0 or sections starting with "0." are not source of truth
        if section_num == "0" or section_num.startswith("0."):
            return True
        # Check if heading contains cross-reference keywords
        if heading_text:
            heading_lower = heading_text.lower()
            if "cross-reference" in heading_lower or "cross-references" in heading_lower:
                return True
            if "overview" in heading_lower and section_num.startswith("0."):
                return True
        return False

    def _clean_heading(self, section_num: str, heading_text: str) -> str:
        """Remove section number from heading if present, handling edge cases."""
        # Special case: if section is "0" and heading starts with "0. ", return just the text after
        if section_num == "0":
            if heading_text.startswith("0. "):
                return heading_text[3:]
            elif heading_text.startswith("0 "):
                return heading_text[2:]

        # Remove section number prefix (e.g., "2.1 AddFile" -> "AddFile")
        # Match the exact section number at the start
        section_pattern = re.escape(section_num) + r'(?:\.\s+|\s+)'
        heading_clean = re.sub(r'^' + section_pattern, '', heading_text)

        # If that didn't work, try generic pattern
        if heading_clean == heading_text:
            heading_clean = re.sub(r'^\d+(?:\.\d+)*\s+', '', heading_text)

        return heading_clean

    def _format_section_number(
        self, section_num: str, heading_text: str
    ) -> str:
        """Format section number with period if original heading had one.

        Args:
            section_num: The section number (e.g., "11", "2.1", "8.1.6")
            heading_text: The full heading text (e.g., "11. HashType Type")

        Returns:
            Section number formatted with period if original had one (e.g., "11." or "2.1")
        """
        # Check if the original heading has a period after the section number
        if heading_text.startswith(section_num + '.'):
            return section_num + '.'
        return section_num

    def _load_index(self):
        """Load api_go_defs_index.md and extract method/type -> spec_file mappings with anchors."""
        if not self.index_file.exists():
            warning_msg = format_issue_message(
                "warning",
                "Index file not found",
                str(self.index_file),
                None,
                "skipping index validation",
                None,
                False
            )
            print(warning_msg)
            return

        # Verify index file is within repo
        if not self._is_safe_path(self.index_file):
            warning_msg = format_issue_message(
                "warning",
                "Index file path unsafe",
                str(self.index_file),
                None,
                "skipping index validation",
                None,
                False
            )
            print(warning_msg)
            return

        content = self.file_cache.get_content(self.index_file)

        # Pattern: **`Package.AddFile`** - [File Management API - AddFile]
        #          (api_file_mgmt_addition.md#anchor)
        pattern = r'\*\*`([^`]+)`\*\*\s*-\s*\[([^\]]+)\]\(([^)]+)\)'
        for match in re.finditer(pattern, content):
            method_type = match.group(1)
            link_text = match.group(2)
            # e.g., "api_file_mgmt_addition.md#21-addfile" or
            #      "api_file_mgmt_addition.md"
            link_target = match.group(3)

            # Extract spec file and anchor with security validation
            if '#' in link_target:
                spec_file, anchor = link_target.split('#', 1)
                # Validate anchor is safe
                if not self._validate_anchor(anchor):
                    continue  # Skip this entry if anchor is unsafe
                self.index_anchors[method_type] = anchor
            else:
                spec_file = link_target
                self.index_anchors[method_type] = None

            # Validate spec file name is safe
            if not self._validate_spec_file_name(spec_file):
                continue  # Skip this entry if filename is unsafe

            self.index_entries[method_type] = spec_file
            self.index_link_texts[method_type] = link_text

    def _parse_markdown_anchors(
        self, file_path: Path
    ) -> Tuple[Set[str], Dict[str, Tuple[str, str]]]:
        """Parse markdown file to extract all heading anchors and section numbers.

        Returns:
            Tuple of (anchors set, sections dict where key is section_num and
            value is (heading_text, anchor))
        """
        # Use shared utility function with cache
        return extract_headings_with_section_numbers(
            file_path, min_level=2, max_level=6, file_cache=self.file_cache
        )

    def _is_safe_path(self, file_path: Path) -> bool:
        """Check if a path is safe (within repo and no traversal)."""
        return is_safe_path(file_path, self.repo_root)

    def _validate_spec_file_name(self, spec_file: str) -> bool:
        """Validate that spec file name is safe (no path traversal, no separators)."""
        return validate_spec_file_name(spec_file)

    def _validate_anchor(self, anchor: str) -> bool:
        """Validate that anchor is safe (no path traversal, no separators)."""
        return validate_anchor(anchor)

    def _get_spec_file_path(self, spec_file: str) -> Optional[Path]:
        """Get the full path to a spec file with security validation."""
        # Validate file name is safe
        if not self._validate_spec_file_name(spec_file):
            return None

        # Construct path within docs directory
        file_path = self.docs_dir / spec_file

        # Verify the resolved path is within repo
        if not self._is_safe_path(file_path):
            return None

        return file_path

    def _validate_reference(self, ref: SpecReference) -> List[ValidationIssue]:
        """Validate a single reference. Returns list of ValidationIssue objects (empty if valid)."""
        errors: List[ValidationIssue] = []

        # Try to find correct reference from index if we have a function name
        if ref.function_name:
            correct_ref = self._find_correct_reference_from_index(ref.function_name)
            if correct_ref:
                spec_file, section_num, heading = correct_ref
                # Filter out section 0 and cross-reference sections
                if not self._is_section_0_or_cross_reference(section_num, heading):
                    # heading is already cleaned, section_num is already formatted
                    ref.suggested_ref = f"{spec_file}: {section_num} {heading}"

        # First check format
        if not ref.is_valid_format:
            # If we have a parsed section from anchor, try to validate it
            if ref.spec_file and ref.section:
                spec_path = self._get_spec_file_path(ref.spec_file)
                if spec_path and spec_path.exists():
                    # Load sections if not cached
                    if ref.spec_file not in self.spec_sections:
                        _, sections = self._parse_markdown_anchors(spec_path)
                        self.spec_sections[ref.spec_file] = sections

                    sections = self.spec_sections[ref.spec_file]
                    if ref.section in sections:
                        # Found the section! But check if it's section 0 or cross-reference
                        heading_text, _ = sections[ref.section]
                        if self._is_section_0_or_cross_reference(ref.section, heading_text):
                            # Don't suggest section 0 or cross-reference sections
                            # Try to find a better section
                            similar = [
                                s for s in sections.keys()
                                if not self._is_section_0_or_cross_reference(
                                    s, sections[s][0]
                                ) and (
                                    ref.section.startswith(s) or
                                    s.startswith(ref.section)
                                )
                            ]
                            if similar:
                                actual_heading, _ = sections[similar[0]]
                                heading_clean = self._clean_heading(similar[0], actual_heading)
                                section_formatted = self._format_section_number(
                                    similar[0], actual_heading
                                )
                                ref.suggested_ref = (
                                    f"{ref.spec_file}: {section_formatted} {heading_clean}"
                                )
                        else:
                            # Valid section - generate correct format
                            heading_clean = self._clean_heading(ref.section, heading_text)
                            section_formatted = self._format_section_number(
                                ref.section, heading_text
                            )
                            ref.suggested_ref = (
                                f"{ref.spec_file}: {section_formatted} {heading_clean}"
                            )
                    else:
                        # Section not found, try to find similar (skip section 0 and cross-refs)
                        similar = [
                            s for s in sections.keys()
                            if not self._is_section_0_or_cross_reference(
                                s, sections[s][0]
                            ) and (
                                ref.section.startswith(s) or
                                s.startswith(ref.section)
                            )
                        ]
                        if similar:
                            actual_heading, _ = sections[similar[0]]
                            heading_clean = self._clean_heading(similar[0], actual_heading)
                            section_formatted = self._format_section_number(
                                similar[0], actual_heading
                            )
                            ref.suggested_ref = (
                                f"{ref.spec_file}: {section_formatted} {heading_clean}"
                            )

            # Build error message
            message_parts = []
            if not ref.suggested_ref:
                message_parts.append(
                    "Invalid format. Expected: 'file_name.md: "
                    "section_number heading_text'"
                )
                message_parts.append(f"Got: '{ref.raw_ref}'")
                message_parts.append(
                    "Example: 'api_file_mgmt_addition.md: "
                    "2.1 AddFile Package Method'"
                )
            else:
                message_parts.append(f"Invalid reference: '{ref.raw_ref}'")

            # Check if spec file exists
            if ref.spec_file:
                spec_path = self._get_spec_file_path(ref.spec_file)
                if not spec_path or not spec_path.exists():
                    message_parts.append(f"Spec file not found: {ref.spec_file}")

            if message_parts:
                errors.append(ValidationIssue(
                    "invalid_spec_ref_format",
                    ref.file_path,
                    ref.line_num,
                    ref.line_num,
                    " ".join(message_parts),
                    severity='error',
                    suggestion=ref.suggested_ref,
                    raw_ref=ref.raw_ref,
                    spec_file=ref.spec_file
                ))
            return errors

        if not ref.spec_file:
            errors.append(ValidationIssue(
                "missing_spec_file",
                ref.file_path,
                ref.line_num,
                ref.line_num,
                "No spec file specified in reference",
                severity='error',
                raw_ref=ref.raw_ref
            ))
            return errors

        spec_path = self._get_spec_file_path(ref.spec_file)
        if not spec_path or not spec_path.exists():
            errors.append(ValidationIssue(
                "spec_file_not_found",
                ref.file_path,
                ref.line_num,
                ref.line_num,
                f"Spec file not found: {ref.spec_file}",
                severity='error',
                raw_ref=ref.raw_ref,
                spec_file=ref.spec_file
            ))
            return errors

        # Load anchors and sections for this file if not cached
        if ref.spec_file not in self.spec_anchors:
            anchors, sections = self._parse_markdown_anchors(spec_path)
            self.spec_anchors[ref.spec_file] = anchors
            self.spec_sections[ref.spec_file] = sections

        # Validate section exists
        sections = self.spec_sections[ref.spec_file]
        if ref.section not in sections:
            # Check for similar section numbers (skip section 0 and cross-refs)
            similar = [
                s for s in sections.keys()
                if not self._is_section_0_or_cross_reference(s, sections[s][0]) and (
                    ref.section.startswith(s) or s.startswith(ref.section)
                )
            ]
            if similar:
                actual_heading, _ = sections[similar[0]]
                # Remove section number from heading if present
                heading_clean = self._clean_heading(similar[0], actual_heading)
                if not ref.suggested_ref:
                    ref.suggested_ref = (
                        f"{ref.spec_file}: {similar[0]} {heading_clean}"
                    )
                message = (
                    f"Section '{ref.section}' not found. "
                    f"Did you mean: '{similar[0]} {actual_heading}'?"
                )
            else:
                message = f"Section '{ref.section}' not found in {ref.spec_file}"
                # Show available sections (excluding section 0 and cross-refs)
                available_sections = []
                if sections:
                    available = [
                        (num, heading) for num, (heading, _) in sorted(sections.items())
                        if not self._is_section_0_or_cross_reference(num, heading)
                    ][:5]
                    available_sections = [f"{num} {heading}" for num, heading in available]
                    if available_sections:
                        message += f" Available sections: {', '.join(available_sections)}..."

            errors.append(ValidationIssue(
                "section_not_found",
                ref.file_path,
                ref.line_num,
                ref.line_num,
                message,
                severity='error',
                suggestion=ref.suggested_ref,
                raw_ref=ref.raw_ref,
                spec_file=ref.spec_file,
                section=ref.section
            ))
            return errors

        # Validate heading matches (or is close to) the actual heading
        actual_heading, actual_anchor = sections[ref.section]

        # Normalize both headings for comparison (lowercase, remove extra spaces)
        normalized_ref = re.sub(r'\s+', ' ', ref.heading.lower().strip())
        normalized_actual = re.sub(r'\s+', ' ', actual_heading.lower().strip())

        # Check if they match (allowing for some flexibility)
        if normalized_ref != normalized_actual:
            # Check if the reference heading is a substring of the actual heading
            # (e.g., "AddFile" matches "2.1 AddFile Package Method")
            if normalized_ref not in normalized_actual:
                # Check if actual heading contains the reference heading (reverse)
                if normalized_actual not in normalized_ref:
                    # Remove section number from heading if present
                    heading_clean = self._clean_heading(ref.section, actual_heading)
                    if not ref.suggested_ref:
                        ref.suggested_ref = f"{ref.spec_file}: {ref.section} {heading_clean}"
                    errors.append(
                        f"  Heading mismatch for section {ref.section}"
                    )
                    errors.append(f"    Expected: '{actual_heading}'")
                    errors.append(f"    Got: '{ref.heading}'")
                    errors.append(
                        f"    Correct format: '{ref.spec_file}: {ref.section} {actual_heading}'"
                    )

        return errors

    def _extract_function_or_type_name(self, go_file: Path, line_num: int) -> Optional[str]:
        """Extract function or type name from Go code context around the specification comment."""
        try:
            lines = self.file_cache.get_lines(go_file)

            # Look backwards from the comment line to find function/type definition
            # Check up to 30 lines before the comment
            start_line = max(0, line_num - 30)
            context = ''.join(lines[start_line:line_num])

            # Pattern 1: func (receiver) MethodName(...) - most common
            # Try to get receiver type too for better matching
            match = re.search(r'func\s+\([^)]*(\w+)[^)]*\)\s+([A-Z][a-zA-Z0-9_]+)\s*\(', context)
            if match:
                receiver = match.group(1)
                method = match.group(2)
                # Try full qualified name first (e.g., "Package.AddFile")
                if receiver and receiver[0].isupper():
                    return f"{receiver}.{method}"
                return method

            # Pattern 2: func FunctionName(...)
            match = re.search(r'func\s+([A-Z][a-zA-Z0-9_]+)\s*\(', context)
            if match:
                return match.group(1)

            # Pattern 3: type TypeName
            match = re.search(r'type\s+([A-Z][a-zA-Z0-9_]+)\s+(?:struct|interface|\[|$)', context)
            if match:
                return match.group(1)

            # Pattern 4: const ConstName
            match = re.search(r'const\s+([A-Z][a-zA-Z0-9_]+)\s*=', context)
            if match:
                return match.group(1)

            # Pattern 5: var VarName
            match = re.search(r'var\s+([A-Z][a-zA-Z0-9_]+)\s*=', context)
            if match:
                return match.group(1)

            return None
        except (IOError, OSError, UnicodeDecodeError):
            # File read errors - return None silently as this is a helper function
            return None
        except Exception as e:
            # Unexpected errors - log but don't fail
            if self.verbose:
                print(f"Warning: Unexpected error extracting function name: {e}", file=sys.stderr)
            return None

    def _parse_anchor_to_section_and_heading(
        self, anchor: str, spec_file: str
    ) -> Optional[Tuple[str, str]]:
        """Parse an anchor back to section number and heading by looking it up in the spec file.

        Since anchors are generated by removing dots from section numbers, we can't reliably
        reverse them (e.g., "11" could be section "11" or "1.1"). Instead, we look up the
        actual heading in the spec file that matches this anchor.

        Args:
            anchor: The anchor string (e.g., "11-hashtype-type")
            spec_file: The spec file name (e.g., "api_file_mgmt_file_entry.md")

        Returns:
            Tuple of (section_num, heading_text) or None if not found
            Never returns section 0 or cross-reference sections.
        """
        if not anchor:
            return None

        # Validate anchor is safe
        if not self._validate_anchor(anchor):
            return None

        # Validate spec file name is safe
        if not self._validate_spec_file_name(spec_file):
            return None

        spec_path = self._get_spec_file_path(spec_file)
        if not spec_path or not spec_path.exists():
            return None

        # Double-check path is safe after resolution
        if not self._is_safe_path(spec_path):
            return None

        # Load sections if not cached
        if spec_file not in self.spec_sections:
            _, sections = self._parse_markdown_anchors(spec_path)
            self.spec_sections[spec_file] = sections

        sections = self.spec_sections[spec_file]

        # Find the section that has this anchor, but skip section 0 and cross-references
        for section_num, (heading_text, heading_anchor) in sections.items():
            if heading_anchor == anchor:
                # Filter out section 0 and cross-reference sections
                if not self._is_section_0_or_cross_reference(section_num, heading_text):
                    return (section_num, heading_text)

        # If not found in cached sections, try generating the anchor from all headings
        # and match it using shared utility function
        try:
            # Use shared utility to extract headings with section anchors
            headings = extract_h2_plus_headings_with_sections(
                spec_path, file_cache=self.file_cache
            )

            for heading_level, heading_text, line_num, plain_anchor, section_anchor in headings:
                # Check if the section_anchor matches (this is the format: "123-heading-text")
                if section_anchor == anchor:
                    # Extract section number from heading text
                    section_match = re.match(r'^(\d+(?:\.\d+)*)', heading_text)
                    if section_match:
                        section_num = section_match.group(1)
                        # Filter out section 0 and cross-reference sections
                        if not self._is_section_0_or_cross_reference(section_num, heading_text):
                            return (section_num, heading_text)

                # Also check plain anchor as fallback
                if plain_anchor == anchor:
                    # Extract section number from heading text
                    section_match = re.match(r'^(\d+(?:\.\d+)*)', heading_text)
                    if section_match:
                        section_num = section_match.group(1)
                        # Filter out section 0 and cross-reference sections
                        if not self._is_section_0_or_cross_reference(section_num, heading_text):
                            return (section_num, heading_text)
        except (IOError, OSError, UnicodeDecodeError):
            # File read errors - return None silently
            pass
        except Exception as e:
            # Unexpected errors - log but don't fail
            if self.verbose:
                print(f"Warning: Unexpected error parsing anchor: {e}", file=sys.stderr)
            pass

        return None

    def _try_exact_match(self, function_name: str) -> Optional[Tuple[str, str, str]]:
        """Try to find an exact match for the function name in the index."""
        if function_name not in self.index_entries:
            return None

        spec_file = self.index_entries[function_name]
        anchor = self.index_anchors.get(function_name)
        if anchor:
            parsed = self._parse_anchor_to_section_and_heading(anchor, spec_file)
            if parsed:
                section_num, heading = parsed
                # Filter out section 0 and cross-reference sections
                if not self._is_section_0_or_cross_reference(section_num, heading):
                    heading_clean = self._clean_heading(section_num, heading)
                    section_formatted = self._format_section_number(section_num, heading)
                    return (spec_file, section_formatted, heading_clean)

        # Fallback to old method if no anchor
        fallback = self._find_section_for_spec_file(spec_file, function_name)
        if fallback:
            spec_file_fb, section_num, heading_text = fallback
            # Filter out section 0 and cross-reference sections
            if not self._is_section_0_or_cross_reference(section_num, heading_text):
                heading_clean = self._clean_heading(section_num, heading_text)
                section_formatted = self._format_section_number(section_num, heading_text)
                return (spec_file_fb, section_formatted, heading_clean)
        return None

    def _calculate_match_score(self, function_base: str, index_name: str) -> int:
        """Calculate match score for a function name against an index entry.

        Returns:
            Score: 100 = exact base match, 50 = receiver match, 0 = no match
        """
        index_base = index_name.split('.')[-1] if '.' in index_name else index_name

        if function_base == index_base:
            # Exact base name match (e.g., "AddFile" == "AddFile")
            return 100
        elif index_name.endswith('.' + function_base):
            # Receiver match (e.g., "Package.AddFile" ends with ".AddFile")
            return 50
        return 0

    def _create_candidate_from_anchor(
        self, anchor: str, spec_file: str, score: int
    ) -> Optional[Tuple[int, str, str, str]]:
        """Create a candidate tuple from an anchor match."""
        parsed = self._parse_anchor_to_section_and_heading(anchor, spec_file)
        if not parsed:
            return None

        section_num, heading = parsed
        # Filter out section 0 and cross-reference sections
        if self._is_section_0_or_cross_reference(section_num, heading):
            return None

        heading_clean = self._clean_heading(section_num, heading)
        section_formatted = self._format_section_number(section_num, heading)
        return (score, spec_file, section_formatted, heading_clean)

    def _create_candidate_from_fallback(
        self, fallback: Tuple[str, str, str], score: int
    ) -> Optional[Tuple[int, str, str, str]]:
        """Create a candidate tuple from a fallback match."""
        spec_file_fb, section_num, heading_text = fallback
        # Filter out section 0 and cross-reference sections
        if self._is_section_0_or_cross_reference(section_num, heading_text):
            return None

        heading_clean = self._clean_heading(section_num, heading_text)
        section_formatted = self._format_section_number(section_num, heading_text)
        return (score, spec_file_fb, section_formatted, heading_clean)

    def _find_partial_matches(
        self, function_base: str
    ) -> List[Tuple[int, str, str, str]]:
        """Find partial matches for a function name in the index.

        Returns:
            List of candidate tuples: (score, spec_file, section_formatted, heading_clean)
        """
        candidates = []

        for index_name, spec_file in self.index_entries.items():
            score = self._calculate_match_score(function_base, index_name)
            if score == 0:
                continue

            anchor = self.index_anchors.get(index_name)
            if anchor:
                candidate = self._create_candidate_from_anchor(anchor, spec_file, score)
                if candidate:
                    candidates.append(candidate)
            else:
                # Fallback to old method if no anchor
                fallback = self._find_section_for_spec_file(spec_file, index_name)
                if fallback:
                    candidate = self._create_candidate_from_fallback(fallback, score)
                    if candidate:
                        candidates.append(candidate)

        return candidates

    def _find_correct_reference_from_index(
        self, function_name: str
    ) -> Optional[Tuple[str, str, str]]:
        """Find correct reference from index. Returns (spec_file, section_num, heading) or None.

        Now uses anchors from the index file directly, assuming they are correct.
        """
        if not function_name:
            return None

        # Try exact match first
        exact_match = self._try_exact_match(function_name)
        if exact_match:
            return exact_match

        # Try partial matches (e.g., "AddFile" matches "Package.AddFile")
        # Extract base name (last component after dot)
        function_base = function_name.split('.')[-1] if '.' in function_name else function_name

        candidates = self._find_partial_matches(function_base)

        # Return the best match (highest score, first in case of ties)
        if candidates:
            # Sort by score descending, then return first
            candidates.sort(key=lambda x: x[0], reverse=True)
            _, spec_file, section_formatted, heading_clean = candidates[0]
            return (spec_file, section_formatted, heading_clean)

        return None

    def _find_section_for_spec_file(
        self, spec_file: str, context: str
    ) -> Optional[Tuple[str, str, str]]:
        """Find the section number and heading for a spec file based on context."""
        spec_path = self._get_spec_file_path(spec_file)
        if not spec_path or not spec_path.exists():
            return None

        # Load sections if not cached
        if spec_file not in self.spec_sections:
            _, sections = self._parse_markdown_anchors(spec_path)
            self.spec_sections[spec_file] = sections

        sections = self.spec_sections[spec_file]

        # Extract function name from context (e.g., "Package.AddFile" -> "AddFile")
        function_name = context.split('.')[-1] if '.' in context else context

        # Try to find a section that matches the function/type name
        # Look for sections that contain the function name or keywords from it
        context_lower = context.lower()
        function_name_lower = function_name.lower()
        context_words = set(re.findall(r'[a-zA-Z]+', context_lower))

        best_match = None
        best_score = 0

        for section_num, (heading_text, _) in sections.items():
            # Skip section 0 and cross-reference sections - they are not source of truth
            if self._is_section_0_or_cross_reference(section_num, heading_text):
                continue

            heading_lower = heading_text.lower()
            heading_words = set(re.findall(r'[a-zA-Z]+', heading_lower))

            # Score based on:
            # 1. Exact function name match in heading (highest priority)
            # 2. Word overlap
            score = 0
            if function_name_lower in heading_lower:
                score += 100  # Strong match
            if function_name in heading_text:
                score += 50   # Case-sensitive match

            # Word overlap score
            overlap = len(context_words & heading_words)
            score += overlap * 10

            if score > best_score:
                best_score = score
                best_match = (spec_file, section_num, heading_text)

        # If we found a good match (not section 0), return it
        if best_match and best_score >= 10:  # Require at least some match
            return best_match

        # Otherwise, try to find a section with the function name in the content
        # by reading the spec file and searching for the function name
        try:
            lines = self.file_cache.get_lines(spec_path)

            # Search for function name in headings or near headings
            for i, line in enumerate(lines):
                if re.match(r'^#{2,6}\s+', line):  # This is a heading
                    # Check if function name appears in nearby content
                    nearby = ' '.join(lines[max(0, i):min(len(lines), i + 10)])
                    if function_name in nearby or function_name_lower in nearby.lower():
                        # Extract section number from this heading
                        heading_match = re.match(r'^#{2,6}\s+(.+)', line)
                        if heading_match:
                            heading_text = heading_match.group(1).strip()
                            section_match = re.match(r'^(\d+(?:\.\d+)*)', heading_text)
                            if section_match:
                                section_num = section_match.group(1)
                                # Filter out section 0 and cross-reference sections
                                if not self._is_section_0_or_cross_reference(
                                        section_num, heading_text):
                                    return (spec_file, section_num, heading_text)
        except (IOError, OSError, UnicodeDecodeError):
            # File read errors - continue to fallback logic
            pass
        except Exception as e:
            # Unexpected errors - log but don't fail
            if self.verbose:
                print(f"Warning: Unexpected error finding section: {e}", file=sys.stderr)
            pass

        # Try to find first non-zero, non-cross-reference section
        if sections:
            for section_num, (heading_text, _) in sorted(sections.items()):
                if not self._is_section_0_or_cross_reference(section_num, heading_text):
                    return (spec_file, section_num, heading_text)

        # Never return section 0 or cross-reference sections
        return None

    def find_spec_references(self, go_file: Path) -> List[SpecReference]:
        """Find all specification references in a Go file."""
        references = []

        try:
            lines = self.file_cache.get_lines(go_file)
            for line_num, line in enumerate(lines, 1):
                # Match "// Specification: ..." comments
                match = re.search(r'//\s*Specification:\s*(.+)', line)
                if match:
                    ref_text = match.group(1).strip()
                    # Handle multi-line references (some may have continuation)
                    if ref_text:
                        ref = SpecReference(go_file, line_num, ref_text)
                        # Try to extract function/type name for index lookup
                        ref.function_name = (
                            self._extract_function_or_type_name(
                                go_file, line_num
                            )
                        )
                        references.append(ref)
        except (IOError, OSError) as e:
            # File read errors - create ValidationIssue
            error = ValidationIssue(
                "file_read_error",
                go_file,
                0,
                0,
                f"Could not read file: {e}",
                severity='error'
            )
            self.issues.append(error)
        except UnicodeDecodeError as e:
            # Encoding errors - create ValidationIssue
            error = ValidationIssue(
                "file_encoding_error",
                go_file,
                0,
                0,
                f"Could not decode file (encoding issue): {e}",
                severity='error'
            )
            self.issues.append(error)
        except Exception as e:
            # Unexpected errors - create ValidationIssue
            error = ValidationIssue(
                "unexpected_error",
                go_file,
                0,
                0,
                f"Unexpected error reading file: {e}",
                severity='error'
            )
            self.issues.append(error)

        return references

    def validate_all(
        self, check_index: bool = False, verbose: bool = False, output=None
    ) -> Tuple[int, List[str]]:
        """Validate all references in Go files. Returns (error_count, error_messages)."""
        error_count = 0
        error_messages: List[str] = []
        all_issues: List[ValidationIssue] = []

        # Find all Go files
        go_files = list(self.api_go_dir.rglob("*.go"))

        if output:
            output.add_verbose_line(
                f"Scanning {len(go_files)} Go files for specification references..."
            )
        else:
            print(f"Scanning {len(go_files)} Go files for specification references...")

        for go_file in go_files:
            references = self.find_spec_references(go_file)
            if not references:
                continue

            if verbose:
                rel_path = go_file.relative_to(self.repo_root)
                if output:
                    output.add_verbose_line(
                        f"  Checking {rel_path} ({len(references)} reference(s))"
                    )
                else:
                    print(f"  Checking {rel_path} ({len(references)} reference(s))")

            for ref in references:
                if verbose:
                    if output:
                        output.add_verbose_line(f"    Validating: {ref.raw_ref}")
                    else:
                        print(f"    Validating: {ref.raw_ref}")
                errors = self._validate_reference(ref)
                if errors:
                    error_count += len(errors)
                    all_issues.extend(errors)
                    # Format error messages for output
                    for error in errors:
                        if isinstance(error, ValidationIssue):
                            error_messages.append(error.format_message(no_color=False))
                        else:
                            error_messages.append(str(error))

        return error_count, error_messages


def main():
    parser = argparse.ArgumentParser(
        description="Validate Go file specification references against spec files",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  python3 scripts/validate_go_spec_references.py
  python3 scripts/validate_go_spec_references.py --verbose
  python3 scripts/validate_go_spec_references.py --output report.txt
        """
    )
    parser.add_argument(
        "--repo-root",
        type=Path,
        default=get_workspace_root(),
        help="Repository root directory (default: parent of script directory)",
    )
    parser.add_argument(
        "--check-index",
        action="store_true",
        help="Also validate against api_go_defs_index.md (not yet implemented)",
    )
    parser.add_argument(
        "--verbose", "-v",
        action="store_true",
        help="Show all references being checked",
    )
    parser.add_argument(
        "--output", "-o",
        type=str,
        help="Output file path for validation report (default: stdout)",
    )
    parser.add_argument(
        "--no-fail",
        action="store_true",
        help="Exit with code 0 even if errors are found",
    )
    parser.add_argument(
        "--nocolor", "--no-color",
        action="store_true",
        help="Disable colored output",
    )

    args = parser.parse_args()

    # Create output builder (header streams immediately if verbose)
    no_color = args.nocolor or parse_no_color_flag(sys.argv)
    output = OutputBuilder(
        "Go Specification References Validation",
        "Validates Go file specification references against spec files",
        no_color=no_color,
        verbose=args.verbose,
        output_file=args.output
    )

    validator = SpecValidator(args.repo_root)

    if args.verbose:
        output.add_verbose_line(f"Repository root: {args.repo_root}")
        output.add_verbose_line(f"Docs directory: {validator.docs_dir}")
        output.add_verbose_line(f"API Go directory: {validator.api_go_dir}")
        output.add_verbose_line(f"Index file: {validator.index_file}")
        output.add_blank_line("working_verbose")

    error_count, error_messages = validator.validate_all(
        check_index=args.check_index, verbose=args.verbose, output=output
    )

    if error_messages:
        output.add_errors_header()
        for error_msg in error_messages:
            output.add_error_line(error_msg)
        output.add_blank_line("error")
        output.add_line(
            f"Found {error_count} error(s)",
            section="error"
        )
        # Check if any messages contain suggestions (-> format)
        has_suggestions = any(" -> " in msg for msg in error_messages)
        if has_suggestions:
            output.add_blank_line("error")
            output.add_line(
                "Note: After applying these reference updates, "
                "verify that each updated reference points to the correct content.",
                section="error"
            )
        output.add_failure_message("Validation failed. Please fix the errors above.")
        output.print()
        return output.get_exit_code(args.no_fail)

    output.add_success_message("All specification references are valid!")
    output.print()
    return 0


if __name__ == "__main__":
    sys.exit(main())
