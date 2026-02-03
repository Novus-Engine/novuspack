"""Markdown and heading utilities for validation scripts."""

import re
import sys
from pathlib import Path
from typing import Optional, List, Set, Tuple, Dict
from dataclasses import dataclass

from lib.validation._fs import FileContentCache

# Compiled regex patterns for performance (module level)
_RE_HEADING_PATTERN = re.compile(r'^(#{1,6})\s+(.+)$')
_RE_DECIMAL_PATTERN = re.compile(r'\d+\.\d+')
_RE_SENTENCE_END_PATTERN = re.compile(r'[.!?]+(?=\s+|$)')
_RE_HEADING_NUM_PATTERN = re.compile(r'^\d+(?:\.\d+)*$')


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
                    f"heading_num must be a dotted number like '1', '2.4', or '3.5.6', "
                    f"got: {self.heading_num!r}"
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
    except (OSError, UnicodeDecodeError, ValueError) as e:
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
    for heading_text, heading_level, _line_num in headings:
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


def _is_ellipsis_at(text: str, punct_pos: int) -> bool:
    """Return True if punct_pos is part of an ellipsis."""
    if punct_pos > 0 and punct_pos + 2 < len(text):
        if text[punct_pos - 1:punct_pos + 2] == '...' or text[punct_pos:punct_pos + 3] == '...':
            return True
    return punct_pos + 1 < len(text) and text[punct_pos:punct_pos + 2] == '..'


def _is_decimal_near(text: str, punct_pos: int) -> bool:
    """Return True if punct_pos is inside a decimal number."""
    context = text[max(0, punct_pos - 10):min(len(text), punct_pos + 10)]
    return bool(_RE_DECIMAL_PATTERN.search(context))


def _is_url_near(text: str, punct_pos: int, has_urls: bool) -> bool:
    """Return True if punct_pos is inside a URL."""
    if not has_urls:
        return False
    url_context = text[
        max(0, punct_pos - 30):min(len(text), punct_pos + 30)
    ]
    return contains_url(url_context)


def _is_abbreviation_at(
    text: str, text_lower: str, punct_pos: int, _abbreviations: set
) -> str:
    """Return the word before punct_pos (for abbreviation check). Empty if not found."""
    word_start = punct_pos
    while word_start > 0 and (text[word_start - 1].isalnum() or text[word_start - 1] == '.'):
        word_start -= 1
    return text_lower[word_start:punct_pos + 1]


def _should_skip_sentence_boundary(
    text: str,
    text_lower: str,
    punct_pos: int,
    abbreviations: set,
    has_urls: bool,
) -> bool:
    """Return True if this punctuation position is not a real sentence end."""
    if _is_ellipsis_at(text, punct_pos):
        return True
    if _is_decimal_near(text, punct_pos):
        return True
    if _is_url_near(text, punct_pos, has_urls):
        return True
    word_before = _is_abbreviation_at(text, text_lower, punct_pos, abbreviations)
    return word_before in abbreviations


def _try_append_hybrid_sentence(
    text: str,
    punct_pos: int,
    *,
    punct_end: int,
    last_end: int,
    word_before: str,
    abbreviations: set,
    sentences: list,
) -> bool:
    """If period+uppercase sentence end, append and return True. Else return False."""
    if text[punct_pos] != '.' or punct_end >= len(text):
        return False
    next_char_pos = punct_end
    while next_char_pos < len(text) and text[next_char_pos].isspace():
        next_char_pos += 1
    if (next_char_pos >= len(text) or
            not text[next_char_pos].isupper() or
            word_before in abbreviations):
        return False
    sentence = text[last_end:punct_end].strip()
    if sentence:
        sentences.append(sentence)
    return True


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
    has_urls = contains_url(text)
    matches = list(_RE_SENTENCE_END_PATTERN.finditer(text))
    if not matches:
        return 1 if text.strip() else 0

    sentences = []
    last_end = 0

    for match in matches:
        punct_pos = match.start()
        punct_end = match.end()
        if punct_end < len(text) and text[punct_end].isspace():
            while punct_end < len(text) and text[punct_end].isspace():
                punct_end += 1

        if _should_skip_sentence_boundary(
            text, text_lower, punct_pos, abbreviations, has_urls
        ):
            continue

        word_before = _is_abbreviation_at(text, text_lower, punct_pos, abbreviations)
        if _try_append_hybrid_sentence(
            text, punct_pos,
            punct_end=punct_end, last_end=last_end, word_before=word_before,
            abbreviations=abbreviations, sentences=sentences,
        ):
            last_end = punct_end
            continue

        sentence = text[last_end:punct_end].strip()
        if sentence:
            sentences.append(sentence)
        last_end = punct_end

    if last_end < len(text):
        remaining = text[last_end:].strip()
        if remaining:
            sentences.append(remaining)
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

    for idx, (_line_num, level, _text) in enumerate(headings):
        parent_index = None
        if level > 2:
            # H3 and beyond need a parent
            parent_level = level - 1
            parent_index = heading_stack.get(parent_level)

        hierarchy[idx] = parent_index

        # Update heading stack - set this heading as the current parent at its level
        heading_stack[level] = idx

        # Clear deeper levels when we move up in hierarchy
        levels_to_clear = [lvl for lvl in heading_stack if lvl > level]
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

    for idx, (_line_num, level, _text) in enumerate(all_headings):
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
    *,
    max_prose_lines: int = 5,
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
    for idx, (line_num, level, _text) in enumerate(all_headings):
        if line_num == heading_line and level == heading_level:
            heading_index = idx
            break

    if heading_index is None:
        return {'is_organizational': False, 'is_empty': False, 'sentence_count': 0}

    # Find next heading (any level) to determine section boundaries
    next_heading_line = None
    for line_num, _level, _text in all_headings:
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
        if not sentence_count or (subheadings and len(prose_lines) <= max_prose_lines):
            return {
                'is_organizational': True,
                'is_empty': (not sentence_count),
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
    if '/' in anchor or '\\' in anchor:
        return False
    if '..' in anchor:
        return False
    if '\x00' in anchor:
        return False
    if not re.match(r'^[a-zA-Z0-9_\-]+$', anchor):
        return False
    return True
