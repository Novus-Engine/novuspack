"""SpecValidator for validate_go_spec_references.py."""

import re
from pathlib import Path
from typing import Dict, List, Optional, Set, Tuple

from lib._validation_utils import (
    format_issue_message,
    ValidationIssue,
    is_safe_path,
    validate_spec_file_name,
    validate_anchor,
    extract_headings_with_section_numbers,
    FileContentCache,
    DOCS_DIR,
    TECH_SPECS_DIR,
)
from lib._validate_go_spec_references_models import SpecReference
from lib._validate_go_spec_references_section_finder import SectionFinder


class _SpecValidatorContext:
    """Public adapter for SectionFinder (avoids protected-access)."""

    def __init__(self, validator):
        self._validator = validator
        self.file_cache = validator.file_cache
        self.spec_sections = validator.spec_sections
        self.index_entries = validator.index_entries
        self.index_anchors = validator.index_anchors

    @property
    def output(self):
        return self._validator.output_for_context

    @property
    def verbose(self):
        return self._validator.verbose

    def get_spec_file_path(self, spec_file):
        return self._validator.get_spec_file_path(spec_file)

    def parse_markdown_anchors(self, file_path):
        return self._validator.parse_markdown_anchors(file_path)

    def is_section_0_or_cross_reference(self, section_num, heading_text=""):
        return self._validator.is_section_0_or_cross_reference(
            section_num, heading_text
        )

    def clean_heading(self, section_num, heading_text):
        return self._validator.clean_heading(section_num, heading_text)

    def format_section_number(self, section_num, heading_text):
        return self._validator.format_section_number(section_num, heading_text)

    def validate_anchor(self, anchor):
        return self._validator.validate_anchor(anchor)

    def validate_spec_file_name(self, spec_file):
        return self._validator.validate_spec_file_name(spec_file)

    def is_safe_path(self, file_path):
        return self._validator.is_safe_path(file_path)


class SpecValidator:
    """Validates specification references."""

    def __init__(self, repo_root: Path):
        self.repo_root = repo_root
        self.docs_dir = repo_root / DOCS_DIR / TECH_SPECS_DIR
        self.api_go_dir = repo_root / "api" / "go"
        self.index_file = self.docs_dir / "api_go_defs_index.md"
        self.verbose = False
        self.issues: List[ValidationIssue] = []
        self._output = None  # Set in validate_all; used for warnings

        # File content cache to avoid repeated reads
        self.file_cache = FileContentCache()

        # Cache of parsed spec files (file -> set of anchors)
        self.spec_anchors: Dict[str, Set[str]] = {}
        # file -> {section_num: (heading_text, anchor)}
        self.spec_sections: Dict[
            str, Dict[str, Tuple[str, str]]
        ] = {}

        # Index file (loaded on first validate_all when output is available)
        self.index_entries: Dict[str, str] = {}  # method/type -> spec_file
        self.index_link_texts: Dict[str, str] = {}  # method/type -> link text for context
        self.index_anchors: Dict[str, str] = {}  # method/type -> anchor (e.g., "11-hashtype-type")
        self._index_loaded = False
        self._section_finder = SectionFinder(_SpecValidatorContext(self))

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
            if heading_text.startswith("0 "):
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
        self, section_num: str, _heading_text: str
    ) -> str:
        """Format section number for reference strings."""
        return section_num

    def _ensure_index_loaded(self, output=None):
        """Load index once; emit warnings via output if provided."""
        if self._index_loaded:
            return
        self._index_loaded = True
        self._load_index(output)

    def _load_index(self, output=None):
        """Load api_go_defs_index.md and extract method/type -> spec_file mappings with anchors."""
        if not self.index_file.exists():
            warning_msg = format_issue_message(
                "warning",
                "Index file not found",
                str(self.index_file),
                message="skipping index validation",
                no_color=output.no_color if output else False
            )
            if output:
                output.add_warning_line(warning_msg)
            else:
                print(warning_msg)  # noqa: T201
            return

        # Verify index file is within repo
        if not self._is_safe_path(self.index_file):
            warning_msg = format_issue_message(
                "warning",
                "Index file path unsafe",
                str(self.index_file),
                message="skipping index validation",
                no_color=output.no_color if output else False
            )
            if output:
                output.add_warning_line(warning_msg)
            else:
                print(warning_msg)  # noqa: T201
            return

        content = self.file_cache.get_content(self.index_file)

        # Pattern: **`Package.AddFile`** - [File Management API - AddFile]
        pattern = r'\*\*`([^`]+)`\*\*\s*-\s*\[([^\]]+)\]\(([^)]+)\)'
        for match in re.finditer(pattern, content):
            method_type = match.group(1)
            link_text = match.group(2)
            link_target = match.group(3)

            if '#' in link_target:
                spec_file, anchor = link_target.split('#', 1)
                if not self._validate_anchor(anchor):
                    continue
                self.index_anchors[method_type] = anchor
            else:
                spec_file = link_target
                self.index_anchors[method_type] = None

            if not self._validate_spec_file_name(spec_file):
                continue

            self.index_entries[method_type] = spec_file
            self.index_link_texts[method_type] = link_text

    def _parse_markdown_anchors(
        self, file_path: Path
    ) -> Tuple[Set[str], Dict[str, Tuple[str, str]]]:
        """Parse markdown file to extract all heading anchors and section numbers."""
        return extract_headings_with_section_numbers(
            file_path, min_level=2, max_level=6, file_cache=self.file_cache
        )

    def _is_safe_path(self, file_path: Path) -> bool:
        """Check if a path is safe (within repo and no traversal)."""
        return is_safe_path(file_path, self.repo_root)

    def _validate_spec_file_name(self, spec_file: str) -> bool:
        """Validate that spec file name is safe."""
        return validate_spec_file_name(spec_file)

    def _validate_anchor(self, anchor: str) -> bool:
        """Validate that anchor is safe."""
        return validate_anchor(anchor)

    # Public API for _SpecValidatorContext / SectionFinder
    @property
    def output_for_context(self):
        return self._output

    def get_spec_file_path(self, spec_file: str) -> Optional[Path]:
        return self._get_spec_file_path(spec_file)

    def parse_markdown_anchors(
        self, file_path: Path
    ) -> Tuple[Set[str], Dict[str, Tuple[str, str]]]:
        return self._parse_markdown_anchors(file_path)

    def is_section_0_or_cross_reference(
        self, section_num: str, heading_text: str = ""
    ) -> bool:
        return self._is_section_0_or_cross_reference(section_num, heading_text)

    def clean_heading(self, section_num: str, heading_text: str) -> str:
        return self._clean_heading(section_num, heading_text)

    def format_section_number(self, section_num: str, heading_text: str) -> str:
        return self._format_section_number(section_num, heading_text)

    def validate_anchor(self, anchor: str) -> bool:
        return self._validate_anchor(anchor)

    def validate_spec_file_name(self, spec_file: str) -> bool:
        return self._validate_spec_file_name(spec_file)

    def is_safe_path(self, file_path: Path) -> bool:
        return self._is_safe_path(file_path)

    def _get_spec_file_path(self, spec_file: str) -> Optional[Path]:
        """Get the full path to a spec file with security validation."""
        if not self._validate_spec_file_name(spec_file):
            return None
        file_path = self.docs_dir / spec_file
        if not self._is_safe_path(file_path):
            return None
        return file_path

    def _compute_suggested_ref_from_sections(
        self, ref: SpecReference, sections: Dict[str, Tuple[str, str]]
    ) -> Optional[str]:
        """Compute suggested_ref from cached sections (section in or similar)."""
        if ref.section in sections:
            heading_text, _ = sections[ref.section]
            if self._is_section_0_or_cross_reference(ref.section, heading_text):
                similar = [
                    s for s in sections.keys()
                    if not self._is_section_0_or_cross_reference(s, sections[s][0])
                    and (ref.section.startswith(s) or s.startswith(ref.section))
                ]
                if not similar:
                    return None
                section_key = similar[0]
            else:
                section_key = ref.section
        else:
            similar = [
                s for s in sections.keys()
                if not self._is_section_0_or_cross_reference(s, sections[s][0])
                and (ref.section.startswith(s) or s.startswith(ref.section))
            ]
            if not similar:
                return None
            section_key = similar[0]
        actual_heading, _ = sections[section_key]
        heading_clean = self._clean_heading(section_key, actual_heading)
        section_formatted = self._format_section_number(section_key, actual_heading)
        return f"{ref.spec_file}: {section_formatted} {heading_clean}"

    def _try_suggest_ref_from_index(self, ref: SpecReference) -> None:
        """If ref has function_name, try to set ref.suggested_ref from index."""
        if not ref.function_name:
            return
        correct_ref = self._section_finder.find_correct_reference_from_index(
            ref.function_name
        )
        if not correct_ref:
            return
        spec_file, section_num, heading = correct_ref
        if not self._is_section_0_or_cross_reference(section_num, heading):
            ref.suggested_ref = f"{spec_file}: {section_num} {heading}"

    def _ensure_spec_sections_loaded(self, ref: SpecReference, spec_path: Path) -> None:
        """Load spec_anchors and spec_sections for ref.spec_file if not cached."""
        if ref.spec_file not in self.spec_anchors:
            anchors, sections = self._parse_markdown_anchors(spec_path)
            self.spec_anchors[ref.spec_file] = anchors
            self.spec_sections[ref.spec_file] = sections

    def _build_invalid_format_issues(self, ref: SpecReference) -> List[ValidationIssue]:
        if ref.spec_file and ref.section:
            spec_path = self._get_spec_file_path(ref.spec_file)
            if spec_path and spec_path.exists():
                self._ensure_spec_sections_loaded(ref, spec_path)
                ref.suggested_ref = self._compute_suggested_ref_from_sections(
                    ref, self.spec_sections[ref.spec_file]
                )
        message_parts = []
        if not ref.suggested_ref:
            message_parts.extend([
                "Invalid format. Expected: 'file_name.md: section_number heading_text'",
                f"Got: '{ref.raw_ref}'",
                "Example: 'api_file_mgmt_addition.md: 2.1 AddFile Package Method'",
            ])
        else:
            message_parts.append(f"Invalid reference: '{ref.raw_ref}'")
        if ref.spec_file:
            spec_path = self._get_spec_file_path(ref.spec_file)
            if not spec_path or not spec_path.exists():
                message_parts.append(f"Spec file not found: {ref.spec_file}")
        if not message_parts:
            return []
        return [
            ValidationIssue.create(
                "invalid_spec_ref_format",
                ref.file_path,
                ref.line_num,
                ref.line_num,
                message=" ".join(message_parts),
                severity="error",
                suggestion=ref.suggested_ref,
                raw_ref=ref.raw_ref,
                spec_file=ref.spec_file,
            )
        ]

    def _build_missing_spec_file_issue(self, ref: SpecReference) -> List[ValidationIssue]:
        if ref.spec_file:
            return []
        return [
            ValidationIssue.create(
                "missing_spec_file",
                ref.file_path,
                ref.line_num,
                ref.line_num,
                message="No spec file specified in reference",
                severity="error",
                raw_ref=ref.raw_ref,
            )
        ]

    def _resolve_spec_path_errors(
        self,
        ref: SpecReference,
    ) -> Tuple[Optional[Path], List[ValidationIssue]]:
        spec_path = self._get_spec_file_path(ref.spec_file)
        if spec_path and spec_path.exists():
            return spec_path, []
        return None, [
            ValidationIssue.create(
                "spec_file_not_found",
                ref.file_path,
                ref.line_num,
                ref.line_num,
                message=f"Spec file not found: {ref.spec_file}",
                severity="error",
                raw_ref=ref.raw_ref,
                spec_file=ref.spec_file,
            )
        ]

    def _build_section_not_found_issues(
        self,
        ref: SpecReference,
        sections: Dict[str, Tuple[str, str]],
    ) -> List[ValidationIssue]:
        if ref.section in sections:
            return []
        similar = [
            s for s in sections.keys()
            if not self._is_section_0_or_cross_reference(s, sections[s][0])
            and (ref.section.startswith(s) or s.startswith(ref.section))
        ]
        if similar:
            actual_heading, _ = sections[similar[0]]
            heading_clean = self._clean_heading(similar[0], actual_heading)
            if not ref.suggested_ref:
                ref.suggested_ref = f"{ref.spec_file}: {similar[0]} {heading_clean}"
            message = (
                f"Section '{ref.section}' not found. "
                f"Did you mean: '{similar[0]} {heading_clean}'?"
            )
        else:
            message = f"Section '{ref.section}' not found in {ref.spec_file}"
            if sections:
                available = [
                    (num, self._clean_heading(num, heading))
                    for num, (heading, _) in sorted(sections.items())
                    if not self._is_section_0_or_cross_reference(num, heading)
                ][:5]
                if available:
                    message += (
                        " Available sections: "
                        f"{', '.join(f'{n} {h}' for n, h in available)}..."
                    )
        return [
            ValidationIssue.create(
                "section_not_found",
                ref.file_path,
                ref.line_num,
                ref.line_num,
                message=message,
                severity="error",
                suggestion=ref.suggested_ref,
                raw_ref=ref.raw_ref,
                spec_file=ref.spec_file,
                section=ref.section,
            )
        ]

    def _build_heading_mismatch_issues(
        self,
        ref: SpecReference,
        actual_heading: str,
    ) -> List[ValidationIssue]:
        heading_clean = self._clean_heading(ref.section, actual_heading)
        ref_heading_clean = self._clean_heading(ref.section, ref.heading)
        normalized_ref = re.sub(r'\s+', ' ', ref_heading_clean.lower().strip())
        normalized_actual = re.sub(r'\s+', ' ', heading_clean.lower().strip())
        if (
            normalized_ref == normalized_actual
            or normalized_ref in normalized_actual
            or normalized_actual in normalized_ref
        ):
            return []
        if not ref.suggested_ref:
            ref.suggested_ref = f"{ref.spec_file}: {ref.section} {heading_clean}"
        return [
            ValidationIssue.create(
                "heading_mismatch",
                ref.file_path,
                ref.line_num,
                ref.line_num,
                message=(
                    f"Heading mismatch for section {ref.section}. "
                    f"Expected: '{heading_clean}'. Got: '{ref.heading}'. "
                    f"Correct format: '{ref.spec_file}: {ref.section} {heading_clean}'"
                ),
                severity="error",
                suggestion=ref.suggested_ref,
                raw_ref=ref.raw_ref,
                spec_file=ref.spec_file,
                section=ref.section,
            )
        ]

    def _validate_reference(self, ref: SpecReference) -> List[ValidationIssue]:
        """Validate a single reference. Returns list of ValidationIssue objects (empty if valid)."""
        self._try_suggest_ref_from_index(ref)

        if not ref.is_valid_format:
            return self._build_invalid_format_issues(ref)

        missing_spec_file = self._build_missing_spec_file_issue(ref)
        if missing_spec_file:
            return missing_spec_file

        spec_path, spec_path_errors = self._resolve_spec_path_errors(ref)
        if spec_path_errors:
            return spec_path_errors

        self._ensure_spec_sections_loaded(ref, spec_path)
        sections = self.spec_sections[ref.spec_file]
        section_errors = self._build_section_not_found_issues(ref, sections)
        if section_errors:
            return section_errors

        actual_heading, _ = sections[ref.section]
        return self._build_heading_mismatch_issues(ref, actual_heading)

    def find_spec_references(self, go_file: Path) -> List[SpecReference]:
        """Find all specification references in a Go file."""
        references = []
        try:
            lines = self.file_cache.get_lines(go_file)
            for line_num, line in enumerate(lines, 1):
                match = re.search(r'//\s*Specification:\s*(.+)', line)
                if match:
                    ref_text = match.group(1).strip()
                    if ref_text:
                        ref = SpecReference(go_file, line_num, ref_text)
                        ref.function_name = self._section_finder.extract_function_or_type_name(
                            go_file, line_num
                        )
                        references.append(ref)
        except (IOError, OSError) as e:
            error = ValidationIssue.create(
                "file_read_error", go_file, 0, 0,
                message=f"Could not read file: {e}", severity='error'
            )
            self.issues.append(error)
        except UnicodeDecodeError as e:
            error = ValidationIssue.create(
                "file_encoding_error", go_file, 0, 0,
                message=f"Could not decode file (encoding issue): {e}", severity='error'
            )
            self.issues.append(error)
        except (ValueError, KeyError, TypeError, AttributeError, RuntimeError) as e:
            error = ValidationIssue.create(
                "unexpected_error", go_file, 0, 0,
                message=f"Unexpected error reading file: {e}", severity='error'
            )
            self.issues.append(error)
        return references

    def validate_all(
        self, _check_index: bool = False, verbose: bool = False, output=None
    ) -> Tuple[int, List[str]]:
        """Validate all references in Go files. Returns (error_count, error_messages)."""
        all_issues: List[ValidationIssue] = []
        self.verbose = verbose
        self.issues = all_issues
        self._output = output

        if output:
            self._ensure_index_loaded(output)

        go_files = list(self.api_go_dir.rglob("*.go"))
        if output:
            output.add_verbose_line(
                f"Scanning {len(go_files)} Go files for specification references..."
            )

        for go_file in go_files:
            references = self.find_spec_references(go_file)
            if not references:
                continue
            if verbose and output:
                rel_path = go_file.relative_to(self.repo_root)
                output.add_verbose_line(
                    f"  Checking {rel_path} ({len(references)} reference(s))"
                )
            for ref in references:
                if verbose and output:
                    output.add_verbose_line(f"    Validating: {ref.raw_ref}")
                errors = self._validate_reference(ref)
                if errors:
                    all_issues.extend(errors)

        error_messages = [
            issue.format_message(no_color=False) if isinstance(issue, ValidationIssue)
            else str(issue)
            for issue in all_issues
        ]
        return len(all_issues), error_messages
