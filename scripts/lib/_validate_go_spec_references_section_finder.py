"""Section and index lookup for SpecValidator (validate_go_spec_references)."""

import re
from pathlib import Path
from typing import List, Optional, Tuple

from lib._validation_utils import extract_h2_plus_headings_with_sections


class SectionFinder:
    """Finds sections and correct references from index; used by SpecValidator."""

    def __init__(self, ctx):
        """ctx: object with file_cache, spec_sections, index_entries, index_anchors,
        get_spec_file_path, parse_markdown_anchors, is_section_0_or_cross_reference,
        clean_heading, format_section_number, validate_anchor, validate_spec_file_name,
        is_safe_path, output, verbose.
        """
        self._v = ctx

    def extract_function_or_type_name(self, go_file: Path, line_num: int) -> Optional[str]:
        """Extract function or type name from Go code context around the specification comment."""
        try:
            lines = self._v.file_cache.get_lines(go_file)
            start_line = max(0, line_num - 30)
            context = ''.join(lines[start_line:line_num])

            match = re.search(r'func\s+\([^)]*(\w+)[^)]*\)\s+([A-Z][a-zA-Z0-9_]+)\s*\(', context)
            if match:
                receiver, method = match.group(1), match.group(2)
                if receiver and receiver[0].isupper():
                    return f"{receiver}.{method}"
                return method
            match = re.search(r'func\s+([A-Z][a-zA-Z0-9_]+)\s*\(', context)
            if match:
                return match.group(1)
            match = re.search(r'type\s+([A-Z][a-zA-Z0-9_]+)\s+(?:struct|interface|\[|$)', context)
            if match:
                return match.group(1)
            match = re.search(r'const\s+([A-Z][a-zA-Z0-9_]+)\s*=', context)
            if match:
                return match.group(1)
            match = re.search(r'var\s+([A-Z][a-zA-Z0-9_]+)\s*=', context)
            if match:
                return match.group(1)
            return None
        except (IOError, OSError, UnicodeDecodeError):
            return None
        except (ValueError, KeyError, TypeError, AttributeError, RuntimeError) as e:
            if self._v.output and self._v.verbose:
                self._v.output.add_warning_line(
                    f"Unexpected error extracting function name: {e}"
                )
            return None

    def find_correct_reference_from_index(
        self, function_name: str
    ) -> Optional[Tuple[str, str, str]]:
        """Find correct reference from index. Returns (spec_file, section_num, heading) or None."""
        if not function_name:
            return None
        exact = self._try_exact_match(function_name)
        if exact:
            return exact
        function_base = function_name.split('.')[-1] if '.' in function_name else function_name
        candidates = self._find_partial_matches(function_base)
        if candidates:
            candidates.sort(key=lambda x: x[0], reverse=True)
            _, spec_file, section_formatted, heading_clean = candidates[0]
            return (spec_file, section_formatted, heading_clean)
        return None

    def _try_exact_match(self, function_name: str) -> Optional[Tuple[str, str, str]]:
        if function_name not in self._v.index_entries:
            return None
        spec_file = self._v.index_entries[function_name]
        anchor = self._v.index_anchors.get(function_name)
        if anchor:
            parsed = self._parse_anchor_to_section_and_heading(anchor, spec_file)
            if parsed:
                section_num, heading = parsed
                if not self._v.is_section_0_or_cross_reference(section_num, heading):
                    heading_clean = self._v.clean_heading(section_num, heading)
                    section_formatted = self._v.format_section_number(section_num, heading)
                    return (spec_file, section_formatted, heading_clean)
        fallback = self._find_section_for_spec_file(spec_file, function_name)
        if fallback:
            spec_file_fb, section_num, heading_text = fallback
            if not self._v.is_section_0_or_cross_reference(section_num, heading_text):
                heading_clean = self._v.clean_heading(section_num, heading_text)
                section_formatted = self._v.format_section_number(section_num, heading_text)
                return (spec_file_fb, section_formatted, heading_clean)
        return None

    def _calculate_match_score(self, function_base: str, index_name: str) -> int:
        index_base = index_name.split('.')[-1] if '.' in index_name else index_name
        if function_base == index_base:
            return 100
        if index_name.endswith('.' + function_base):
            return 50
        return 0

    def _create_candidate_from_anchor(
        self, anchor: str, spec_file: str, score: int
    ) -> Optional[Tuple[int, str, str, str]]:
        parsed = self._parse_anchor_to_section_and_heading(anchor, spec_file)
        if not parsed:
            return None
        section_num, heading = parsed
        if self._v.is_section_0_or_cross_reference(section_num, heading):
            return None
        heading_clean = self._v.clean_heading(section_num, heading)
        section_formatted = self._v.format_section_number(section_num, heading)
        return (score, spec_file, section_formatted, heading_clean)

    def _create_candidate_from_fallback(
        self, fallback: Tuple[str, str, str], score: int
    ) -> Optional[Tuple[int, str, str, str]]:
        spec_file_fb, section_num, heading_text = fallback
        if self._v.is_section_0_or_cross_reference(section_num, heading_text):
            return None
        heading_clean = self._v.clean_heading(section_num, heading_text)
        section_formatted = self._v.format_section_number(section_num, heading_text)
        return (score, spec_file_fb, section_formatted, heading_clean)

    def _find_partial_matches(
        self, function_base: str
    ) -> List[Tuple[int, str, str, str]]:
        candidates = []
        for index_name, spec_file in self._v.index_entries.items():
            score = self._calculate_match_score(function_base, index_name)
            if not score:
                continue
            anchor = self._v.index_anchors.get(index_name)
            if anchor:
                candidate = self._create_candidate_from_anchor(anchor, spec_file, score)
                if candidate:
                    candidates.append(candidate)
            else:
                fallback = self._find_section_for_spec_file(spec_file, index_name)
                if fallback:
                    candidate = self._create_candidate_from_fallback(fallback, score)
                    if candidate:
                        candidates.append(candidate)
        return candidates

    def _parse_anchor_to_section_and_heading(
        self, anchor: str, spec_file: str
    ) -> Optional[Tuple[str, str]]:
        if (not anchor or not self._v.validate_anchor(anchor)
                or not self._v.validate_spec_file_name(spec_file)):
            return None
        spec_path = self._v.get_spec_file_path(spec_file)
        if not spec_path or not spec_path.exists() or not self._v.is_safe_path(spec_path):
            return None
        if spec_file not in self._v.spec_sections:
            _, sections = self._v.parse_markdown_anchors(spec_path)
            self._v.spec_sections[spec_file] = sections
        result = self._find_section_for_anchor_in_cached(
            self._v.spec_sections[spec_file], anchor
        )
        if result is not None:
            return result
        return self._find_section_for_anchor_in_headings(spec_path, anchor)

    def _find_section_for_anchor_in_cached(
        self, sections: dict, anchor: str
    ) -> Optional[Tuple[str, str]]:
        for section_num, (heading_text, heading_anchor) in sections.items():
            if (heading_anchor == anchor
                    and not self._v.is_section_0_or_cross_reference(
                        section_num, heading_text
                    )):
                return (section_num, heading_text)
        return None

    def _find_section_for_anchor_in_headings(
        self, spec_path: Path, anchor: str
    ) -> Optional[Tuple[str, str]]:
        try:
            headings = extract_h2_plus_headings_with_sections(
                spec_path, file_cache=self._v.file_cache
            )
            for _hl, heading_text, _ln, plain_anchor, section_anchor in headings:
                if anchor not in (section_anchor, plain_anchor):
                    continue
                section_match = re.match(r'^(\d+(?:\.\d+)*)', heading_text)
                if section_match and not self._v.is_section_0_or_cross_reference(
                    section_match.group(1), heading_text
                ):
                    return (section_match.group(1), heading_text)
        except (IOError, OSError, UnicodeDecodeError):
            pass
        except (ValueError, KeyError, TypeError, AttributeError, RuntimeError) as e:
            if self._v.output and self._v.verbose:
                self._v.output.add_warning_line(
                    f"Unexpected error parsing anchor: {e}"
                )
        return None

    def _find_section_for_spec_file(
        self, spec_file: str, context: str
    ) -> Optional[Tuple[str, str, str]]:
        spec_path = self._v.get_spec_file_path(spec_file)
        if not spec_path or not spec_path.exists():
            return None
        if spec_file not in self._v.spec_sections:
            _, sections = self._v.parse_markdown_anchors(spec_path)
            self._v.spec_sections[spec_file] = sections
        sections = self._v.spec_sections[spec_file]
        function_name = context.split('.')[-1] if '.' in context else context
        context_lower = context.lower()
        function_name_lower = function_name.lower()
        context_words = set(re.findall(r'[a-zA-Z]+', context_lower))

        best_match = None
        best_score = 0
        for section_num, (heading_text, _) in sections.items():
            if self._v.is_section_0_or_cross_reference(section_num, heading_text):
                continue
            heading_lower = heading_text.lower()
            heading_words = set(re.findall(r'[a-zA-Z]+', heading_lower))
            score = 0
            if function_name_lower in heading_lower:
                score += 100
            if function_name in heading_text:
                score += 50
            overlap = len(context_words & heading_words)
            score += overlap * 10
            if score > best_score:
                best_score = score
                best_match = (spec_file, section_num, heading_text)

        if best_match and best_score >= 10:
            return best_match

        try:
            lines = self._v.file_cache.get_lines(spec_path)
            content_match = self._find_section_in_content_by_function(
                lines, spec_file, function_name, function_name_lower
            )
            if content_match:
                return content_match
        except (IOError, OSError, UnicodeDecodeError):
            pass
        except (ValueError, KeyError, TypeError, AttributeError, RuntimeError) as e:
            if self._v.output and self._v.verbose:
                self._v.output.add_warning_line(
                    f"Unexpected error finding section: {e}"
                )

        if sections:
            for section_num, (heading_text, _) in sorted(sections.items()):
                if not self._v.is_section_0_or_cross_reference(
                    section_num, heading_text
                ):
                    return (spec_file, section_num, heading_text)
        return None

    def _find_section_in_content_by_function(
        self,
        lines: List[str],
        spec_file: str,
        function_name: str,
        function_name_lower: str,
    ) -> Optional[Tuple[str, str, str]]:
        for i, line in enumerate(lines):
            if not re.match(r'^#{2,6}\s+', line):
                continue
            nearby = ' '.join(lines[max(0, i):min(len(lines), i + 10)])
            if (function_name not in nearby
                    and function_name_lower not in nearby.lower()):
                continue
            heading_match = re.match(r'^#{2,6}\s+(.+)', line)
            if not heading_match:
                continue
            heading_text = heading_match.group(1).strip()
            section_match = re.match(r'^(\d+(?:\.\d+)*)', heading_text)
            if not section_match:
                continue
            section_num = section_match.group(1)
            if self._v.is_section_0_or_cross_reference(section_num, heading_text):
                continue
            return (spec_file, section_num, heading_text)
        return None
