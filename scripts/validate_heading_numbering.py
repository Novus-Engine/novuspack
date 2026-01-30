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
- Organizational headings with no content are flagged as errors
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

import re
import sys
from pathlib import Path
from collections import defaultdict
from typing import List

scripts_dir = Path(__file__).parent
lib_dir = scripts_dir / "lib"

# Import shared utilities
for module_path in (str(scripts_dir), str(lib_dir)):
    if module_path not in sys.path:
        sys.path.insert(0, module_path)

# Import shared utilities
from lib._validation_utils import (  # noqa: E402
    OutputBuilder, parse_no_color_flag, format_issue_message,
    get_workspace_root, parse_paths,
    build_heading_hierarchy, is_organizational_heading,
    remove_backticks_keep_content, has_backticks, get_backticks_error_message,
    ValidationIssue, find_markdown_files
)


# Constants for validation thresholds
MAX_HEADING_NUMBER_SEGMENT = 20  # Maximum value for any number segment in heading numbering
MAX_ORGANIZATIONAL_PROSE_LINES = 5  # Maximum lines of prose for organizational heading check

# Compiled regex patterns for performance (module level)
_RE_HEADING_PATTERN = re.compile(r'^(#{1,})\s+(.+)$')
_RE_NUMBERED_HEADING_PATTERN = re.compile(r'^([0-9]+(?:\.[0-9]+)*)\.?\s+(.+)$')
_RE_SPLIT_WORDS = re.compile(r'\S+|\s+')
_RE_WHITESPACE_ONLY = re.compile(r'^\s+$')
_RE_FILENAME_PATTERN = re.compile(r'^[\w\-]+\.(\w+)$')
_RE_NON_WORD_CHARS = re.compile(r'[^\w]')
_RE_NUMBERING_PREFIX = re.compile(r'^([0-9]+(?:\.[0-9]+)*)\.?\s+(.+)$')
_RE_FIRST_LETTER = re.compile(r'[a-zA-Z]')


class HeadingInfo:
    """Represents a heading with its metadata for sorting."""
    def __init__(self, file, line_num, level, numbers, heading_text, full_line,
                 parent=None, issue=None):
        self.file = file
        self.line_num = line_num
        self.level = level
        self.numbers = numbers  # List of integers
        self.heading_text = heading_text
        self.full_line = full_line
        self.parent = parent  # Reference to parent HeadingInfo (if any)
        self.issue = issue  # Reference to related HeadingIssue (if any)
        self.original_number = '.'.join(map(str, numbers))  # Original number as string
        self.corrected_number = None  # Will be set during correction calculation
        self.has_period = False  # For H2: whether period follows number in original
        self.corrected_capitalization = None  # Will be set during capitalization check

    def sort_key(self):
        """Return a sort key for proper numeric ordering."""
        # Create a tuple: (numbers as tuple, level)
        # Sort by numeric values first, then by level for same numbers
        # This ensures proper numeric sorting (e.g., [1, 10] comes after [1, 2])
        return (tuple(self.numbers), self.level)


class HeadingValidator:
    """Validates heading numbering in markdown files."""

    # Pattern to match markdown headings: ##+ followed by optional space and number
    HEADING_PATTERN = _RE_HEADING_PATTERN

    # Pattern to match numbered headings: starts with number(s) followed by
    # optional period and space. Uses [0-9]+\. pattern to match each number
    # segment explicitly. Handles both "1 Title" and "1. Title" formats
    NUMBERED_HEADING_PATTERN = _RE_NUMBERED_HEADING_PATTERN

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

    def check_and_record_backticks(self, filepath, line_num, heading_text, full_line, heading_info):
        """
        Check for backticks in heading text and record error if found.

        Args:
            filepath: Path to the file
            line_num: Line number of the heading
            heading_text: The heading text to check (without number prefix)
            full_line: The full heading line
            heading_info: HeadingInfo object to associate with the error
        """
        if has_backticks(heading_text):
            error = ValidationIssue(
                "heading_formatting",
                Path(filepath),
                line_num,
                line_num,
                get_backticks_error_message(),
                severity='error',
                heading=full_line,
                heading_info=heading_info
            )
            self.issues.append(error)
            heading_info.issue = error
            if self.first_error_line[filepath] is None:
                self.first_error_line[filepath] = error.start_line

    def build_corrected_full_line(self, heading_info):
        """
        Build the corrected full heading line with numbering and capitalization fixes.
        Removes backticks from the heading text in suggestions.

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

        # Determine the corrected heading text (capitalization)
        corrected_text = heading_info.heading_text
        if heading_info.corrected_capitalization:
            corrected_text = heading_info.corrected_capitalization

        # Remove backticks from the heading text for suggestions (keep content)
        corrected_text = remove_backticks_keep_content(corrected_text)

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
        else:
            return f"{prefix} {corrected_number} {corrected_text}"

    def _find_backtick_ranges(self, text):
        """Find all backtick-enclosed sections and their positions."""
        backtick_ranges = []
        i = 0
        while i < len(text):
            if text[i] == '`':
                start = i
                i += 1
                # Find the closing backtick
                while i < len(text) and text[i] != '`':
                    i += 1
                if i < len(text):  # Found closing backtick
                    end = i + 1
                    backtick_ranges.append((start, end))
                    i = end
                else:
                    # Unclosed backtick, treat as regular text
                    break
            else:
                i += 1
        return backtick_ranges

    def _is_in_backticks(self, pos, backtick_ranges):
        """Check if a character position is inside backticks."""
        for start, end in backtick_ranges:
            if start <= pos < end:
                return True
        return False

    def _should_preserve_part(self, part, part_start, part_end, backtick_ranges):
        """Check if a part should be preserved as-is (backticks, underscores, filenames)."""
        # Check if part is inside backticks
        part_in_backticks = any(
            self._is_in_backticks(pos, backtick_ranges)
            for pos in range(part_start, part_end)
        )
        if part_in_backticks:
            return True

        # Check if word contains underscores (e.g., file_name.go, some_function)
        if '_' in part:
            return True

        # Check if word looks like a filename
        common_extensions = [
            'go', 'md', 'txt', 'json', 'yaml', 'yml', 'xml', 'html', 'css', 'js',
            'ts', 'py', 'sh', 'bat', 'ps1', 'java', 'c', 'cpp', 'h', 'hpp',
            'rs', 'rb', 'php', 'sql', 'csv', 'tsv', 'log', 'conf', 'config',
            'ini', 'toml', 'lock', 'sum', 'mod', 'gitignore', 'editorconfig'
        ]
        filename_match = _RE_FILENAME_PATTERN.match(part)
        if filename_match:
            extension = filename_match.group(1).lower()
            if extension in common_extensions:
                return True

        return False

    def _is_in_code_parentheses(self, part, parts, part_index, text):
        """Check if a word is inside parentheses with code-like content (backticks)."""
        text_up_to_here = ''.join(parts[:part_index + 1])
        if '(' not in text_up_to_here:
            return False

        # Find the most recent unclosed parenthesis
        last_open_paren = text_up_to_here.rfind('(')
        if last_open_paren < 0:
            return False

        # Count parentheses to see if we're inside an unclosed one
        text_before = ''.join(parts[:part_index + 1])
        open_count = text_before.count('(') - text_before.count(')')
        if open_count <= 0:
            return False

        # We're inside parentheses - check if they contain backticks
        text_from_paren = text[last_open_paren:]
        close_paren_pos = text_from_paren.find(')', 1)
        if close_paren_pos > 0:
            parens_content = text_from_paren[1:close_paren_pos]
        else:
            parens_content = text_from_paren[1:]

        return '`' in parens_content

    def _should_capitalize_word(
        self, word_clean, is_first, is_last, is_in_code_parens, previous_word_clean=None
    ):
        """
        Determine if a word should be capitalized based on title case rules.

        Args:
            word_clean: The word to check (lowercase, cleaned)
            is_first: True if this is the first word
            is_last: True if this is the last word
            is_in_code_parens: True if word is in code-like parentheses context
            previous_word_clean: The previous word (lowercase, cleaned) for phrasal verb detection
        """
        lowercase_words = {
            'a', 'an', 'the',  # Articles
            'and', 'but', 'or', 'nor', 'for', 'so', 'yet',  # Coordinating conjunctions
            'in', 'on', 'at', 'to', 'for', 'of', 'with', 'by', 'from', 'up', 'about',
            'into', 'onto', 'upon', 'over', 'under', 'above', 'below', 'across',
            'via'  # Short prepositions
        }
        capitalize_prepositions = {
            'through', 'between', 'among', 'during', 'before', 'after', 'within',
            'without', 'against', 'along', 'around', 'behind', 'beside', 'beyond',
            'inside', 'outside', 'throughout', 'toward', 'towards', 'underneath'
        }
        programming_keywords = {
            'return', 'if', 'else', 'for', 'while', 'do', 'switch', 'case',
            'break', 'continue', 'goto', 'throw', 'try', 'catch', 'finally',
            'new', 'delete', 'this', 'super', 'static', 'const', 'let', 'var',
            'function', 'class', 'interface', 'enum', 'type', 'import', 'export',
            'async', 'await', 'yield', 'def', 'lambda', 'pass', 'raise', 'except'
        }

        # Phrasal verb particles that should be capitalized when part of a phrasal verb
        phrasal_particles = {'up', 'down', 'out', 'off', 'in', 'on', 'over', 'away', 'back'}

        # Common phrasal verb bases (verbs that commonly form phrasal verbs)
        phrasal_verb_bases = {
            'clean', 'set', 'look', 'pick', 'give', 'make', 'break', 'build', 'call',
            'check', 'close', 'come', 'cut', 'fill', 'get', 'go', 'grow', 'hang',
            'hold', 'keep', 'line', 'live', 'move', 'open', 'pull', 'put', 'show',
            'sign', 'stand', 'start', 'take', 'turn', 'wake', 'warm', 'wrap', 'bring',
            'carry', 'catch', 'come', 'do', 'draw', 'drop', 'end', 'fall', 'find',
            'fix', 'follow', 'hand', 'head', 'help', 'join', 'jump', 'knock', 'lay',
            'leave', 'let', 'lie', 'lock', 'log', 'mix', 'pass', 'pay', 'point',
            'pop', 'push', 'read', 'run', 'send', 'shut', 'sit', 'slow', 'sort',
            'speed', 'split', 'spread', 'step', 'stick', 'stop', 'switch', 'talk',
            'tear', 'think', 'throw', 'tie', 'try', 'turn', 'use', 'walk', 'wash',
            'watch', 'wear', 'wind', 'work', 'write'
        }

        # If word is in code-like parentheses context, keep programming keywords lowercase
        if is_in_code_parens and word_clean in programming_keywords:
            return False
        if is_first or is_last:
            return True
        if word_clean in capitalize_prepositions:
            return True

        # Check if this is a phrasal verb particle (like "up" in "Clean Up")
        if word_clean in phrasal_particles and previous_word_clean:
            if previous_word_clean in phrasal_verb_bases:
                # This is part of a phrasal verb, so capitalize it
                return True

        if word_clean not in lowercase_words:
            return True
        return False

    def _capitalize_word(self, part, should_capitalize):
        """Apply capitalization to a word part."""
        if not part:
            return part

        # Find first letter using regex
        match = _RE_FIRST_LETTER.search(part)
        if not match:
            return part

        first_letter_idx = match.start()

        if should_capitalize:
            # Capitalize first letter, preserve rest
            return (
                part[:first_letter_idx] +
                part[first_letter_idx].upper() +
                part[first_letter_idx + 1:]
            )
        else:
            # Lowercase the word, but preserve existing capitalization
            # that suggests proper nouns
            has_internal_capitals = any(
                c.isupper() for c in part[first_letter_idx + 1:]
                if c.isalpha()
            )

            if has_internal_capitals:
                # Preserve existing capitalization (likely proper noun/acronym)
                # Only lowercase the first letter if it's uppercase
                if part[first_letter_idx].isupper():
                    return (
                        part[:first_letter_idx] +
                        part[first_letter_idx].lower() +
                        part[first_letter_idx + 1:]
                    )
                else:
                    return part
            else:
                # No internal capitals, lowercase normally
                return (
                    part[:first_letter_idx] +
                    part[first_letter_idx].lower() +
                    part[first_letter_idx + 1:].lower()
                )

    def to_title_case(self, text):
        """
        Convert text to Title Case following standard rules.

        Title Case rules:
        - Capitalize first and last word
        - Capitalize all major words (nouns, verbs, adjectives, adverbs)
        - Lowercase articles (a, an, the)
        - Lowercase coordinating conjunctions (and, but, or, nor, for, so, yet)
        - Lowercase short prepositions (in, on, at, to, for, of, with, by, etc.)
          unless they're the first or last word
        - Capitalize prepositions of 4+ letters (through, between, etc.)
        - Preserve existing capitalization of proper nouns and acronyms
        - Preserve words inside backticks (code) exactly as-is
        """
        if not text:
            return text

        # Find all backtick-enclosed sections
        backtick_ranges = self._find_backtick_ranges(text)

        # Split text into words and separators, preserving structure
        parts = _RE_SPLIT_WORDS.findall(text)
        result_parts = []
        word_indices = []  # Track which parts are words (not whitespace)

        # First pass: identify words and track their positions in original text
        char_pos = 0
        part_positions = []  # Track character positions for each part
        for i, part in enumerate(parts):
            if not _RE_WHITESPACE_ONLY.match(part):  # Not just whitespace
                word_indices.append(i)
            part_positions.append((char_pos, char_pos + len(part)))
            char_pos += len(part)

        if not word_indices:
            return text

        # Process each part
        for i, part in enumerate(parts):
            if _RE_WHITESPACE_ONLY.match(part):
                # Preserve whitespace as-is
                result_parts.append(part)
                continue

            part_start, part_end = part_positions[i]

            # Check if part should be preserved as-is
            if self._should_preserve_part(part, part_start, part_end, backtick_ranges):
                result_parts.append(part)
                continue

            # This is a word
            word_clean = _RE_NON_WORD_CHARS.sub('', part.lower())
            if not word_clean:
                result_parts.append(part)
                continue

            # Check if it's the first or last word
            is_first = i == word_indices[0]
            is_last = i == word_indices[-1]

            # Get previous word for phrasal verb detection
            previous_word_clean = None
            if not is_first:
                # Find the previous word index
                current_word_index = word_indices.index(i)
                if current_word_index > 0:
                    prev_word_i = word_indices[current_word_index - 1]
                    prev_part = parts[prev_word_i]
                    previous_word_clean = _RE_NON_WORD_CHARS.sub('', prev_part.lower())

            # Check if this word is inside parentheses with code-like content
            is_in_code_parens = self._is_in_code_parentheses(part, parts, i, text)

            # Determine if word should be capitalized
            should_capitalize = self._should_capitalize_word(
                word_clean, is_first, is_last, is_in_code_parens, previous_word_clean
            )

            # Apply capitalization
            result_parts.append(self._capitalize_word(part, should_capitalize))

        return ''.join(result_parts)

    def check_capitalization(self, heading_text):
        """
        Check if heading text follows Title Case and return corrected version.

        Returns:
            tuple: (is_correct, corrected_text)
        """
        if not heading_text:
            return True, heading_text

        corrected = self.to_title_case(heading_text)
        is_correct = (heading_text == corrected)

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
            error = ValidationIssue(
                "file_read_error",
                Path(filepath),
                0,
                0,
                f"Could not read file: {e}",
                severity='error'
            )
            self.issues.append(error)
            self.log(f"  Error reading file: {e}")
            return None
        except UnicodeDecodeError as e:
            # Encoding errors - create ValidationIssue
            error = ValidationIssue(
                "file_encoding_error",
                Path(filepath),
                0,
                0,
                f"Could not decode file (encoding issue): {e}",
                severity='error'
            )
            self.issues.append(error)
            self.log(f"  Error decoding file: {e}")
            return None
        except Exception as e:
            # Unexpected errors - create ValidationIssue
            error = ValidationIssue(
                "unexpected_error",
                Path(filepath),
                0,
                0,
                f"Unexpected error reading file: {e}",
                severity='error'
            )
            self.issues.append(error)
            self.log(f"  Unexpected error reading file: {e}")
            return None

    def _process_heading_line(
        self, filepath, line_num, line, stripped_line, heading_stack, in_code_block
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
            warning = ValidationIssue(
                "heading_leading_whitespace",
                Path(filepath),
                line_num,
                line_num,
                msg,
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
            warning = ValidationIssue(
                "heading_leading_whitespace",
                Path(filepath),
                line_num,
                line_num,
                msg,
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
        numbers, title = self.parse_heading_number(heading_text)

        # If this heading is not numbered, create HeadingInfo for it
        if numbers is None:
            return self._create_unnumbered_heading(
                filepath, line_num, level, heading_text, line, heading_stack
            ), in_code_block, True

        # Process numbered heading
        return self._create_numbered_heading(
            filepath, line_num, level, heading_text, line, numbers, heading_stack
        ), in_code_block, False

    def _create_unnumbered_heading(
        self, filepath, line_num, level, heading_text, line, heading_stack
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
            filepath, line_num, level, [], heading_text, line.strip(),
            parent=parent_heading, issue=None
        )
        heading_info.original_number = "MISSING"
        heading_info.has_period = False

        # Check for backticks in heading text
        self.check_and_record_backticks(
            filepath, line_num, heading_text, line.strip(), heading_info
        )

        return heading_info

    def _check_h1_heading(self, filepath, line_num, heading_text, full_line):
        """Track H1 headings and validate numbering/duplicates."""
        self.h1_counts[filepath] += 1
        if self.h1_first_line[filepath] is None:
            self.h1_first_line[filepath] = line_num
        else:
            msg = "More than one H1 heading found. Only the first H1 heading is valid."
            error = ValidationIssue(
                "heading_multiple_h1",
                Path(filepath),
                line_num,
                line_num,
                msg,
                severity='error',
                heading=full_line
            )
            self.issues.append(error)
            if self.first_error_line[filepath] is None:
                self.first_error_line[filepath] = line_num

        numbers, _ = self.parse_heading_number(heading_text)
        if numbers is not None:
            msg = "H1 heading should not be numbered."
            warning = ValidationIssue(
                "heading_h1_numbering",
                Path(filepath),
                line_num,
                line_num,
                msg,
                severity='warning',
                heading=full_line
            )
            self.issues.append(warning)

    def _create_numbered_heading(
        self, filepath, line_num, level, heading_text, line, numbers, heading_stack
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
            filepath, line_num, level, numbers, title, line.strip(),
            parent=parent_heading, issue=None
        )
        # Set original_number from the actual string extracted from the file
        heading_info.original_number = original_number_str
        heading_info.has_period = has_period  # Track period for H2 headings

        # Check for backticks in heading text (after number removed)
        self.check_and_record_backticks(
            filepath, line_num, title, line.strip(), heading_info
        )

        return heading_info

    def _validate_heading_structure(self, filepath, headings, unnumbered_headings):
        """
        Validate heading structure after parsing.

        Returns:
            List of headings (may be modified with issues)
        """
        # Check if we have any headings
        if not headings:
            return []

        # Filter H2 headings once and reuse
        h2_headings = [h for h in headings if h.level == 2]
        if not h2_headings:
            # No H2 headings - skip numbering validation but return headings for other checks
            return headings

        # Check if the first H2 heading is numbered
        first_h2 = min(h2_headings, key=lambda h: h.line_num)
        first_h2_is_numbered = (
            first_h2.numbers and
            len(first_h2.numbers) > 0 and
            first_h2.original_number != "MISSING"
        )

        if first_h2_is_numbered and unnumbered_headings:
            # First H2 is numbered, so unnumbered headings are errors
            for line_num, level, heading_text, full_line, heading_info in unnumbered_headings:
                msg = (f"H{level} heading is missing numbering. "
                       "This document uses numbered headings, so all headings must be numbered.")
                error = ValidationIssue(
                    "heading_missing_numbering",
                    Path(filepath),
                    line_num,
                    line_num,
                    msg,
                    severity='error',
                    heading=full_line,
                    heading_info=heading_info
                )
                self.issues.append(error)
                heading_info.issue = error
                if self.first_error_line[filepath] is None:
                    self.first_error_line[filepath] = line_num

        # Check if the first H2 heading is numbered
        # Only perform numbering validation if the first H2 is numbered
        is_first_h2_numbered = (
            first_h2.numbers and
            len(first_h2.numbers) > 0 and
            first_h2.original_number != "MISSING"
        )

        if not is_first_h2_numbered:
            # First H2 is not numbered - skip numbering validation
            # but still return headings for other checks (capitalization, duplicates, etc.)
            return headings

        # First H2 is numbered - validate it must be "0. Title" or "1. Title"
        if first_h2.numbers[0] not in [0, 1]:
            error = ValidationIssue(
                "heading_first_h2_numbering",
                Path(filepath),
                first_h2.line_num,
                first_h2.line_num,
                (f"First H2 heading must be numbered '0' or '1', "
                 f"got '{first_h2.numbers[0]}'. "
                 "Please run a markdown linter to fix basic heading order, "
                 "then re-run this script."),
                severity='error',
                heading=first_h2.full_line,
                heading_info=first_h2
            )
            self.issues.append(error)
            first_h2.issue = error
            return headings  # Continue with validation despite this error

        # Verify that only H2 headings have parent = None
        for heading in headings:
            if heading.level == 2 and heading.parent is not None:
                # This shouldn't happen, but log it if it does
                self.log(f"  Warning: H2 heading at line {heading.line_num} has a parent")
            elif heading.level > 2 and heading.parent is None:
                # This is an error - H3+ should have a parent
                error = ValidationIssue(
                    "heading_no_parent",
                    Path(filepath),
                    heading.line_num,
                    heading.line_num,
                    (f"H{heading.level} heading has no parent. "
                     "Please run a markdown linter to fix basic heading order, "
                     "then re-run this script."),
                    severity='error',
                    heading=heading.full_line,
                    heading_info=heading
                )
                self.issues.append(error)
                heading.issue = error

        return headings

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
                filepath, line_num, line, stripped_line, heading_stack, in_code_block
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
            levels_to_clear = [lvl for lvl in heading_stack.keys() if lvl > heading_info.level]
            for lvl in levels_to_clear:
                del heading_stack[lvl]

        # Validate heading structure
        return self._validate_heading_structure(filepath, headings, unnumbered_headings)

    def calculate_corrected_numbers(self, filepath, headings):
        """
        Calculate corrected numbers for ALL headings level by level.
        This happens BEFORE validation - we determine what the numbers SHOULD be,
        then compare with original numbers to find errors.
        """
        if not headings:
            return

        # Filter numbered headings and numbered H2 headings in a single pass
        numbered_headings = []
        h2_headings = []
        for h in headings:
            if h.numbers and h.original_number != "MISSING":
                numbered_headings.append(h)
                if h.level == 2:
                    h2_headings.append(h)

        if not numbered_headings:
            # No numbered headings - skip numbering calculation
            return

        # Process H2 headings first (sorted by line number)
        h2_headings.sort(key=lambda h: h.line_num)

        # Determine starting number from first H2 heading
        # If first H2 is "0", start from 0; otherwise start from 1
        start_number = 0
        if h2_headings:
            first_h2_numbers = h2_headings[0].numbers
            if first_h2_numbers and first_h2_numbers[0] == 0:
                start_number = 0
            else:
                start_number = 1

        h2_sequence = start_number - 1  # Will be incremented to start_number
        for heading in h2_headings:
            h2_sequence += 1
            # corrected_number is just the number (no period)
            heading.corrected_number = str(h2_sequence)

        # Process H3+ headings level by level
        # Only process numbered headings
        for level in range(3, 7):  # H3 through H6
            level_headings = [
                h for h in headings
                if h.level == level and h.numbers and h.original_number != "MISSING"
            ]
            if not level_headings:
                continue

            level_headings.sort(key=lambda h: h.line_num)

            # Track sequence for each parent (by parent's corrected_number)
            parent_sequences = {}  # Maps parent corrected_number -> current sequence

            for heading in level_headings:
                # Get parent's corrected number
                if heading.parent and heading.parent.corrected_number:
                    parent_corrected = heading.parent.corrected_number

                    # Initialize sequence for this parent if needed
                    if parent_corrected not in parent_sequences:
                        parent_sequences[parent_corrected] = 0

                    # Increment sequence for this parent
                    parent_sequences[parent_corrected] += 1
                    sequence_num = parent_sequences[parent_corrected]

                    # Build corrected number: parent.corrected_number.sequence_num
                    heading.corrected_number = f"{parent_corrected}.{sequence_num}"
                else:
                    # No parent or parent doesn't have corrected number yet
                    # This shouldn't happen if structure is correct, but handle gracefully
                    heading.corrected_number = heading.original_number

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

            # Expected depth based on heading level
            # H2 => 1 number, H3 => 2 numbers, H4 => 3 numbers, etc.
            expected_depth = level - 1
            actual_depth = len(numbers)

            # Check depth matches heading level
            if actual_depth != expected_depth:
                error = ValidationIssue(
                    "heading_depth_mismatch",
                    Path(filepath),
                    heading.line_num,
                    heading.line_num,
                    f"H{level} heading has {actual_depth} number(s), "
                    f"expected {expected_depth}",
                    severity='error',
                    heading=heading.full_line,
                    heading_info=heading
                )
                self.issues.append(error)
                heading.issue = error
                if self.first_error_line[filepath] is None:
                    self.first_error_line[filepath] = heading.line_num
                continue

            # Check parent relationship
            if level > 2:  # H3 and beyond need to match parent
                if heading.parent is None:
                    error = ValidationIssue(
                        "heading_no_parent_in_validation",
                        Path(filepath),
                        heading.line_num,
                        heading.line_num,
                        f"H{level} heading has no parent. "
                        "Please run a markdown linter to fix basic heading order.",
                        severity='error',
                        heading=heading.full_line,
                        heading_info=heading
                    )
                    self.issues.append(error)
                    heading.issue = error
                    if self.first_error_line[filepath] is None:
                        self.first_error_line[filepath] = heading.line_num
                    continue

            # Check if original_number differs from corrected_number
            if heading.original_number != heading.corrected_number:
                # Find previous heading at same level with same parent for better error message
                prev_heading = None
                if level > 2 and heading.parent:
                    for h in headings_by_line:
                        if h.line_num < heading.line_num and h.level == level:
                            if h.parent == heading.parent:
                                prev_heading = h
                                break

                # Build error message
                if prev_heading:
                    msg = (f"Non-sequential numbering: got '{heading.original_number}', "
                           f"expected '{heading.corrected_number}' "
                           f"(previous was '{prev_heading.original_number}')")
                else:
                    msg = (f"Non-sequential numbering: got '{heading.original_number}', "
                           f"expected '{heading.corrected_number}'")

                error = ValidationIssue(
                    "heading_non_sequential",
                    Path(filepath),
                    heading.line_num,
                    heading.line_num,
                    msg,
                    severity='error',
                    heading=heading.full_line,
                    heading_info=heading
                )
                self.issues.append(error)
                heading.issue = error
                if self.first_error_line[filepath] is None:
                    self.first_error_line[filepath] = heading.line_num

    def check_excessive_numbering(self, filepath, headings):
        """
        Check for H3+ headings where the depth-specific number in the
        corrected numbering exceeds 20. This check is performed AFTER
        corrected heading numbering is calculated.

        For H3 headings (e.g., "1.2"), checks the second number (index 1).
        For H4 headings (e.g., "1.2.3"), checks the third number (index 2).
        For H5 headings (e.g., "1.2.3.4"), checks the fourth number (index 3).
        And so on.
        """
        if not headings:
            return

        # Check only H3+ headings (level >= 3)
        h3_plus_headings = [h for h in headings if h.level >= 3]
        if not h3_plus_headings:
            return

        for heading in h3_plus_headings:
            # Skip if corrected_number is not set
            if not heading.corrected_number:
                continue

            # Parse corrected_number string (e.g., "1.2.3" or "3.25")
            try:
                number_segments = [int(n) for n in heading.corrected_number.split('.')]
            except (ValueError, AttributeError):
                # If parsing fails, skip this heading
                continue

            # For heading level L, check the segment at index (L - 2)
            # H3 (level 3) has 2 numbers, check index 1 (second number)
            # H4 (level 4) has 3 numbers, check index 2 (third number)
            # H5 (level 5) has 4 numbers, check index 3 (fourth number)
            segment_index = heading.level - 2

            # Ensure we have enough segments
            if segment_index >= len(number_segments):
                continue

            # Check if the depth-specific segment exceeds maximum
            segment = number_segments[segment_index]
            if segment > MAX_HEADING_NUMBER_SEGMENT:
                msg = (
                    f"H{heading.level} heading has numbering "
                    f"'{heading.corrected_number}' where number {segment} "
                    f"(at depth {heading.level - 1}) exceeds "
                    f"{MAX_HEADING_NUMBER_SEGMENT}. "
                    "Consider restructuring the document to reduce nesting depth."
                )
                warning = ValidationIssue(
                    "heading_excessive_numbering",
                    Path(filepath),
                    heading.line_num,
                    heading.line_num,
                    msg,
                    severity='warning',
                    heading=heading.full_line,
                    heading_info=heading
                )
                self.issues.append(warning)

    def check_single_word_headings(self, filepath, headings):
        """
        Check for H4+ headings where the title (after numbering) is a single word.
        This is a warning as single-word headings may be too vague or unclear.
        """
        if not headings:
            return

        # Check only H4+ headings (level >= 4)
        h4_plus_headings = [h for h in headings if h.level >= 4]
        if not h4_plus_headings:
            return

        for heading in h4_plus_headings:
            # heading_text contains the title after the number has been removed
            if not heading.heading_text:
                continue

            # Strip whitespace and check if it's a single word (no spaces)
            title = heading.heading_text.strip()
            if title and ' ' not in title:
                msg = (f"H{heading.level} heading has a single-word title '{title}'. "
                       "Consider using a more descriptive multi-word heading.")
                warning = ValidationIssue(
                    "heading_single_word",
                    Path(filepath),
                    heading.line_num,
                    heading.line_num,
                    msg,
                    severity='warning',
                    heading=heading.full_line,
                    heading_info=heading
                )
                self.issues.append(warning)

    def check_duplicate_headings(self, filepath, headings):
        """
        Check for duplicate headings (excluding numbering) across all levels.
        All occurrences after the first are flagged as errors.
        """
        if not headings:
            return

        # Group headings by their title text (case-insensitive, normalized)
        # heading_text contains the title after the number has been removed
        heading_groups = defaultdict(list)
        for heading in headings:
            if not heading.heading_text:
                continue
            # Normalize: strip whitespace and convert to lowercase for comparison
            normalized_title = heading.heading_text.strip().lower()
            if normalized_title:
                heading_groups[normalized_title].append(heading)

        # Find duplicates and flag all but the first occurrence as errors
        for normalized_title, heading_list in heading_groups.items():
            if len(heading_list) > 1:
                # Sort by line number to ensure first occurrence is the earliest
                heading_list.sort(key=lambda h: h.line_num)

                # Flag all subsequent occurrences as errors
                for duplicate_heading in heading_list[1:]:
                    # Find all other occurrences for the error message
                    other_locations = [
                        f"line {h.line_num}" for h in heading_list
                        if h.line_num != duplicate_heading.line_num
                    ]
                    other_locations_str = ", ".join(other_locations)

                    msg = (f"Duplicate heading title '{duplicate_heading.heading_text}' "
                           f"(also appears at {other_locations_str}). "
                           "Each heading should have a unique title.")
                    error = ValidationIssue(
                        "heading_duplicate",
                        Path(filepath),
                        duplicate_heading.line_num,
                        duplicate_heading.line_num,
                        msg,
                        severity='error',
                        heading=duplicate_heading.full_line,
                        heading_info=duplicate_heading
                    )
                    self.issues.append(error)
                    if duplicate_heading.issue is None:
                        duplicate_heading.issue = error
                    if self.first_error_line[filepath] is None:
                        self.first_error_line[filepath] = duplicate_heading.line_num

    def is_go_code_related_heading(self, heading_text):
        """
        Check if a heading is related to Go code elements.

        Returns True if the heading appears to reference a Go code element
        (struct, function, method, interface, type) based on patterns that
        match what the Go code blocks validator expects.

        This allows headings to use actual Go identifiers (like "readOnlyPackage")
        instead of Title Case, which is required by the Go code blocks validator.
        """
        if not heading_text:
            return False

        # Patterns that indicate Go code elements:
        # 1. camelCase identifiers (starts with lowercase, has uppercase in middle)
        #    Examples: "readOnlyPackage", "filePackage", "newFileEntry"
        camel_case_pattern = r'\b[a-z][a-zA-Z]*[A-Z][a-zA-Z]*\b'
        if re.search(camel_case_pattern, heading_text):
            return True

        # 2. Method patterns: "Type.MethodName" (camelCase.TypeMethod)
        method_pattern = r'\b[a-z][a-zA-Z]*\.[A-Z][a-zA-Z]*\b'
        if re.search(method_pattern, heading_text):
            return True

        # 3. Headings that end with Go kind words followed by actual identifiers
        #    Examples: "readOnlyPackage Struct", "filePackage Struct", "newFileEntry Function"
        go_kind_words = ['Struct', 'Function', 'Method', 'Interface', 'Type']
        for kind_word in go_kind_words:
            # Check if heading ends with "KindWord" and has camelCase before it
            pattern = rf'\b[a-z][a-zA-Z]*\s+{kind_word}\b'
            if re.search(pattern, heading_text):
                return True
            # Also check for "Type.MethodName Method" pattern
            method_kind_pattern = rf'\b[a-z][a-zA-Z]*\.[A-Z][a-zA-Z]*\s+{kind_word}\b'
            if re.search(method_kind_pattern, heading_text):
                return True

        return False

    def check_heading_capitalization(self, filepath, headings):
        """
        Check if headings follow Title Case capitalization.
        Adds warnings for headings with incorrect capitalization.

        Skips capitalization checks for headings that reference Go code elements,
        as those must use actual Go identifiers (not Title Case) to satisfy
        the Go code blocks validator.
        """
        if not headings:
            return

        for heading in headings:
            if not heading.heading_text:
                continue

            # Skip capitalization check for Go code-related headings
            # These headings must use actual Go identifiers (e.g., "readOnlyPackage Struct")
            # instead of Title Case (e.g., "ReadOnlyPackage Struct") to satisfy
            # the Go code blocks validator requirements
            if self.is_go_code_related_heading(heading.heading_text):
                continue

            is_correct, corrected = self.check_capitalization(heading.heading_text)
            if not is_correct:
                heading.corrected_capitalization = corrected
                msg = (f"Incorrect capitalization: got '{heading.heading_text}', "
                       f"expected '{corrected}'")
                warning = ValidationIssue(
                    "heading_capitalization",
                    Path(filepath),
                    heading.line_num,
                    heading.line_num,
                    msg,
                    severity='warning',
                    heading=heading.full_line,
                    heading_info=heading
                )
                self.issues.append(warning)

    def check_organizational_headings(self, filepath, headings):
        """
        Check for organizational headings with no content.
        Uses shared utility function is_organizational_heading.
        Errors on headings that are purely organizational with no content.
        """
        if not headings:
            return

        # Read file content
        try:
            with open(filepath, 'r', encoding='utf-8') as f:
                content = f.read()
        except (IOError, OSError) as e:
            # File read errors - log and return
            self.log(f"  Error reading file for organizational check: {e}")
            return
        except UnicodeDecodeError as e:
            # Encoding errors - log and return
            self.log(f"  Error decoding file for organizational check (encoding issue): {e}")
            return
        except Exception as e:
            # Unexpected errors - log and return
            self.log(f"  Unexpected error reading file for organizational check: {e}")
            return

        # Build heading hierarchy using shared utility
        # Convert HeadingInfo objects to (line_num, level, text) tuples
        headings_for_hierarchy = [
            (h.line_num, h.level, h.heading_text)
            for h in headings
        ]
        # Sort by line number
        headings_for_hierarchy.sort(key=lambda x: x[0])
        hierarchy = build_heading_hierarchy(headings_for_hierarchy)

        # Check each heading
        for heading in headings:
            if heading.issue:  # Skip headings that already have errors
                continue

            # Check if heading is organizational
            try:
                result = is_organizational_heading(
                    content,
                    heading.line_num,
                    heading.level,
                    headings_for_hierarchy,
                    hierarchy,
                    max_prose_lines=MAX_ORGANIZATIONAL_PROSE_LINES
                )

                # Error on organizational headings with no content
                if result.get('is_organizational') and result.get('is_empty'):
                    msg = ("Organizational heading with no content. "
                           "Headings should have substantive content or be removed.")
                    error = ValidationIssue(
                        "organizational_heading",
                        Path(filepath),
                        heading.line_num,
                        heading.line_num,
                        msg,
                        severity='error',
                        heading=heading.full_line,
                        heading_info=heading
                    )
                    self.issues.append(error)
                    heading.issue = error
                    if self.first_error_line[filepath] is None:
                        self.first_error_line[filepath] = heading.line_num
            except (ValueError, IndexError, KeyError) as e:
                # Data structure errors - log but don't fail
                self.log(f"  Error checking organizational heading at line {heading.line_num}: {e}")
            except Exception as e:
                # Unexpected errors - log but don't fail
                self.log(
                    f"  Unexpected error checking organizational heading at line "
                    f"{heading.line_num}: {e}"
                )

    def check_h2_period_consistency(self, filepath, headings):
        """
        Check if H2 headings have consistent period usage.
        If first H2 has period, all should have period. If not, none should.
        Only called when there are no errors.
        """
        h2_headings = [h for h in headings if h.level == 2]
        if not h2_headings:
            return

        h2_headings.sort(key=lambda h: h.line_num)
        first_h2 = h2_headings[0]
        expected_has_period = first_h2.has_period

        for heading in h2_headings[1:]:  # Skip first one
            if heading.has_period != expected_has_period:
                expected_str = "with period" if expected_has_period else "without period"
                actual_str = "with period" if heading.has_period else "without period"
                msg = (f"H2 heading period inconsistency: first H2 is {expected_str}, "
                       f"but this heading is {actual_str}. "
                       f"All H2 headings should match the first one.")
                warning = ValidationIssue(
                    "heading_period_inconsistency",
                    Path(filepath),
                    heading.line_num,
                    heading.line_num,
                    msg,
                    severity='warning',
                    heading=heading.full_line,
                    heading_info=heading
                )
                self.issues.append(warning)

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
        # Filter issues by severity in a single loop
        errors = []
        warnings = []
        for issue in self.issues:
            if issue.matches(severity='error'):
                errors.append(issue)
            if issue.matches(severity='warning'):
                warnings.append(issue)

        # Show errors first
        if errors:
            output_builder.add_errors_header()
            output_builder.add_line(
                f"Found {len(errors)} error(s):",
                section="error"
            )
            output_builder.add_blank_line("error")

            # Group errors by issue type (for counting)
            errors_by_type = defaultdict(list)
            for error in errors:
                errors_by_type[error.issue_type].append(error)

            # Group errors by (file, line_num) to combine multiple errors on same line
            errors_by_line = defaultdict(list)
            for error in errors:
                key = (error.file, error.start_line)
                errors_by_line[key].append(error)

            # Display organizational heading errors first
            org_errors = errors_by_type.get("organizational_heading", [])
            if org_errors:
                output_builder.add_line(
                    f"Organizational Heading Errors ({len(org_errors)}):",
                    section="error"
                )
                output_builder.add_blank_line("error")

                # Get organizational errors grouped by line
                org_errors_by_line = {
                    k: v for k, v in errors_by_line.items()
                    if any(
                        isinstance(e, ValidationIssue) and
                        e.issue_type == "organizational_heading"
                        for e in v
                    )
                }

                # Sort by file and line number
                sorted_org_lines = sorted(org_errors_by_line.keys(), key=lambda k: (k[0], k[1]))

                for file, line_num in sorted_org_lines:
                    line_errors = [
                        e for e in org_errors_by_line[(file, line_num)]
                        if isinstance(e, ValidationIssue) and
                        e.issue_type == "organizational_heading"
                    ]
                    if not line_errors:
                        continue

                    rel_file = self.rel_path(file)
                    # Combine all error messages for this line
                    messages = [
                        e.message if isinstance(e, ValidationIssue)
                        else e.get('message', '')
                        for e in line_errors
                    ]
                    combined_message = "; ".join(messages)

                    # Get suggestion from first error with heading_info
                    suggestion = None
                    for error in line_errors:
                        heading_info = error.extra_fields.get('heading_info')
                        if heading_info:
                            suggestion = self.build_corrected_full_line(heading_info)
                            break

                    error_msg = format_issue_message(
                        "error",
                        "Organizational heading",
                        rel_file,
                        line_num,
                        combined_message,
                        suggestion,
                        self.no_color
                    )
                    output_builder.add_error_line(error_msg)

                output_builder.add_blank_line("error")

            # Display heading formatting errors (e.g., backticks)
            formatting_errors = errors_by_type.get("heading_formatting", [])
            if formatting_errors:
                if org_errors:
                    output_builder.add_separator(section="error")
                    output_builder.add_blank_line("error")

                output_builder.add_line(
                    f"Heading Formatting Errors ({len(formatting_errors)}):",
                    section="error"
                )
                output_builder.add_blank_line("error")

                # Get formatting errors grouped by line
                formatting_errors_by_line = {
                    k: v for k, v in errors_by_line.items()
                    if any(
                        isinstance(e, ValidationIssue) and
                        e.issue_type == "heading_formatting"
                        for e in v
                    )
                }

                # Sort by file and line number
                sorted_formatting_lines = sorted(
                    formatting_errors_by_line.keys(), key=lambda k: (k[0], k[1])
                )

                for file, line_num in sorted_formatting_lines:
                    line_errors = [
                        e for e in formatting_errors_by_line[(file, line_num)]
                        if isinstance(e, ValidationIssue) and
                        e.issue_type == "heading_formatting"
                    ]
                    if not line_errors:
                        continue

                    rel_file = self.rel_path(file)
                    # Combine all error messages for this line
                    messages = [
                        e.message if isinstance(e, ValidationIssue)
                        else e.get('message', '')
                        for e in line_errors
                    ]
                    combined_message = "; ".join(messages)

                    # Get suggestion from first error with heading_info
                    suggestion = None
                    for error in line_errors:
                        heading_info = error.extra_fields.get('heading_info')
                        if heading_info:
                            suggestion = self.build_corrected_full_line(heading_info)
                            break

                    error_msg = format_issue_message(
                        "error",
                        "Heading formatting",
                        rel_file,
                        line_num,
                        combined_message,
                        suggestion,
                        self.no_color
                    )
                    output_builder.add_error_line(error_msg)

                output_builder.add_blank_line("error")

            # Display heading numbering errors
            # Get all numbering-related errors
            # (everything except organizational and formatting)
            numbering_errors = [
                e for e in errors
                if isinstance(e, ValidationIssue) and
                e.issue_type not in (
                    "organizational_heading", "heading_formatting"
                )
            ]
            if numbering_errors:
                if org_errors or formatting_errors:
                    output_builder.add_separator(section="error")
                    output_builder.add_blank_line("error")

                output_builder.add_line(
                    f"Heading Numbering Errors ({len(numbering_errors)}):",
                    section="error"
                )
                output_builder.add_blank_line("error")

                # Get numbering errors grouped by line
                numbering_errors_by_line = {
                    k: v for k, v in errors_by_line.items()
                    if any(
                        isinstance(e, ValidationIssue) and
                        e.issue_type not in (
                            "organizational_heading", "heading_formatting"
                        )
                        for e in v
                    )
                }

                # Sort by file and line number
                sorted_numbering_lines = sorted(
                    numbering_errors_by_line.keys(), key=lambda k: (k[0], k[1])
                )

                for file, line_num in sorted_numbering_lines:
                    # All errors in numbering_errors_by_line are already numbering errors
                    # (filtered to exclude organizational_heading and heading_formatting)
                    line_errors = numbering_errors_by_line[(file, line_num)]
                    if not line_errors:
                        continue

                    rel_file = self.rel_path(file)
                    # Combine all error messages for this line
                    messages = [
                        e.message if isinstance(e, ValidationIssue)
                        else e.get('message', '')
                        for e in line_errors
                    ]
                    combined_message = "; ".join(messages)

                    # Get suggestion from first error with heading_info
                    suggestion = None
                    for error in line_errors:
                        heading_info = error.extra_fields.get('heading_info')
                        if heading_info:
                            suggestion = self.build_corrected_full_line(heading_info)
                            break

                    error_msg = format_issue_message(
                        "error",
                        "Heading numbering",
                        rel_file,
                        line_num,
                        combined_message,
                        suggestion,
                        self.no_color
                    )
                    output_builder.add_error_line(error_msg)

                output_builder.add_blank_line("error")

            # Show sorted headings from first error for each file
            # (only for numbering errors, not organizational)
            if numbering_errors:
                for filepath in sorted(self.headings_from_first_error.keys()):
                    errored_headings = self.headings_from_first_error[filepath]
                    if not errored_headings:
                        continue

                    # Filter out headings that only have duplicate errors (no numbering errors)
                    # Only include headings that have numbering errors (original != corrected)
                    # Also include unnumbered headings (original_number == "MISSING")
                    headings_with_numbering_errors = []
                    for heading in errored_headings:
                        # Check if this heading has a numbering error
                        # Include if original_number is "MISSING" (unnumbered heading)
                        if heading.original_number == "MISSING" and heading.corrected_number:
                            headings_with_numbering_errors.append(heading)
                        elif heading.original_number and heading.corrected_number:
                            current_for_comparison = heading.original_number.rstrip('.')
                            correct_for_comparison = heading.corrected_number.rstrip('.')
                            has_numbering_error = (current_for_comparison != correct_for_comparison)
                            if has_numbering_error:
                                headings_with_numbering_errors.append(heading)

                    # Skip this section if there are no numbering errors (only duplicate errors)
                    if not headings_with_numbering_errors:
                        continue

                    first_error_line = self.first_error_line[filepath]
                    rel_file = self.rel_path(filepath)

                    # CRITICAL: This format must be preserved for
                    # apply_heading_corrections.py parsing
                    output_builder.add_separator(section="error")
                    output_builder.add_line(
                        f"Sorted headings from first error (line {first_error_line}) "
                        f"in {rel_file}:",
                        section="error"
                    )
                    output_builder.add_separator(section="error")
                    output_builder.add_blank_line("error")
                    output_builder.add_line(
                        "The following headings should be in this order "
                        "(sorted by numeric values):",
                        section="error"
                    )
                    output_builder.add_blank_line("error")
                    output_builder.add_line(
                        "Format: Line X: [CURRENT] -> [CORRECT] Title",
                        section="error"
                    )
                    output_builder.add_blank_line("error")

                    # Sort headings with numbering errors for display: first by line number
                    # (document order), then by numeric values
                    sorted_headings = sorted(
                        headings_with_numbering_errors,
                        key=lambda h: (h.line_num, h.sort_key())
                    )

                    # Display in sorted order with correct numbering
                    # Find max line number for alignment
                    max_line_num = (
                        max(h.line_num for h in sorted_headings)
                        if sorted_headings else 0
                    )
                    line_num_width = len(str(max_line_num))

                    # Determine period pattern from first H2 heading for display
                    h2_headings_in_output = [h for h in sorted_headings if h.level == 2]
                    display_period = False
                    if h2_headings_in_output:
                        first_h2 = min(h2_headings_in_output, key=lambda h: h.line_num)
                        display_period = first_h2.has_period

                    for heading in sorted_headings:
                        current_number_str = heading.original_number
                        # corrected_number should always be set after
                        # calculate_corrected_numbers
                        # If it's not set, that's a bug, but use original_number
                        # as fallback for safety
                        if heading.corrected_number is None:
                            correct_number_str = current_number_str
                        else:
                            correct_number_str = heading.corrected_number

                        # For unnumbered headings, display "MISSING" as current
                        if current_number_str == "MISSING":
                            current_display = "MISSING"
                        # For display, add period to H2 headings if first H2 has period
                        elif heading.level == 2 and display_period:
                            current_display = f"{current_number_str}."
                        else:
                            current_display = current_number_str

                        # For comparison, strip periods from both
                        current_for_comparison = current_number_str.rstrip('.')
                        correct_for_comparison = correct_number_str.rstrip('.')

                        # Determine if numbering needs to change
                        needs_change = (current_for_comparison != correct_for_comparison)

                        # Check if this is a duplicate heading error (not a numbering error)
                        is_duplicate_error = (
                            heading.issue and
                            isinstance(heading.issue, ValidationIssue) and
                            heading.issue.matches(issue_type="heading_duplicate")
                        )

                        # For display, add period to H2 headings if first H2 has period
                        if heading.level == 2 and display_period:
                            correct_display = f"{correct_number_str}."
                        else:
                            correct_display = correct_number_str

                        # Get heading text with capitalization correction if available
                        heading_text_display = heading.heading_text
                        if heading.corrected_capitalization:
                            heading_text_display = heading.corrected_capitalization

                        # Use different format for duplicate errors when
                        # numbering is correct
                        # This format won't match the correction pattern in
                        # apply_heading_corrections.py
                        if is_duplicate_error and not needs_change:
                            output_builder.add_error_line(
                                f"Line {heading.line_num:{line_num_width}d}: "
                                f"{'#' * heading.level} [{current_display}] (DUPLICATE) "
                                f"{heading_text_display}"
                            )
                        else:
                            # Standard correction format for numbering errors
                            # Include capitalization correction in heading text if available
                            output_builder.add_error_line(
                                f"Line {heading.line_num:{line_num_width}d}: "
                                f"{'#' * heading.level} [{current_display}] -> "
                                f"[{correct_display}] {heading_text_display}"
                            )

                    output_builder.add_blank_line("error")

        # Show warnings
        if warnings:
            output_builder.add_warnings_header()
            output_builder.add_line(
                f"Found {len(warnings)} warning(s):",
                section="warning"
            )
            output_builder.add_blank_line("warning")

            # Group warnings by (file, line_num) to combine multiple warnings on same line
            warnings_by_line = defaultdict(list)
            for warning in warnings:
                # warning is a ValidationIssue
                if isinstance(warning, ValidationIssue):
                    file_key = warning.file
                    line_key = warning.start_line
                else:
                    file_key = warning.get('file', '')
                    line_key = warning.get('line_num', 0)
                key = (file_key, line_key)
                warnings_by_line[key].append(warning)

            # Sort by file and line number for consistent output
            sorted_warning_lines = sorted(warnings_by_line.keys(), key=lambda k: (k[0], k[1]))

            for file, line_num in sorted_warning_lines:
                line_warnings = warnings_by_line[(file, line_num)]
                rel_file = self.rel_path(file)

                # Combine all warning messages for this line
                messages = []
                for warning in line_warnings:
                    # warning is a ValidationIssue
                    if isinstance(warning, ValidationIssue):
                        message = warning.message
                    else:
                        message = warning.get('message', '')
                    # Extract expected value from message if it contains "expected"
                    # Remove the "expected" part from message, keeping "got 'X'"
                    if "expected" in message.lower():
                        message = re.sub(
                            r",\s*expected\s+['\"][^'\"]+['\"]", "", message
                        )
                    messages.append(message)

                combined_message = "; ".join(messages)

                # Get suggestion from first warning with heading_info
                suggestion = None
                for warning in line_warnings:
                    if isinstance(warning, ValidationIssue):
                        heading_info = warning.extra_fields.get('heading_info')
                    else:
                        heading_info = getattr(warning, 'heading_info', None)
                    if heading_info:
                        suggestion = self.build_corrected_full_line(heading_info)
                        break

                warning_msg = format_issue_message(
                    "warning",
                    "Heading numbering",
                    rel_file,
                    line_num,
                    combined_message,
                    suggestion,
                    self.no_color
                )
                output_builder.add_warning_line(warning_msg)

        # Overall status
        if not errors and not warnings:
            # Add summary section with statistics
            files_checked = len(self.all_headings)
            total_headings = sum(len(headings) for headings in self.all_headings.values())
            summary_items = [
                ("Files checked:", files_checked),
                ("Headings checked:", total_headings),
            ]
            output_builder.add_summary_header()
            output_builder.add_summary_section(summary_items)
            output_builder.add_success_message("All heading numbering is valid!")
        elif not errors:
            output_builder.add_line(
                "No heading numbering errors found (only warnings).",
                section="final_message"
            )
        else:
            output_builder.add_failure_message("Validation failed. Please fix the errors above.")


def show_help():
    """Show help message."""
    print(__doc__)


def main():
    """Main entry point."""
    import argparse

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
