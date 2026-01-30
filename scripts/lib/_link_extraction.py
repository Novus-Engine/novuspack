#!/usr/bin/env python3
"""
Link extraction utilities for markdown validation.
"""

import re
import sys
from pathlib import Path


_RE_MARKDOWN_LINK = re.compile(r'\[([^\]]+)\]\(([^)]+)\)')
_RE_CODE_COMMENT = re.compile(r'^\s*(//|#|/\*)')
_RE_PATH_SEP = re.compile(r'[./]')
_RE_SPACE_COMMA = re.compile(r'[\s,]')


def is_in_inline_code(content, match_start, match_end):
    """
    Check if a match is within inline code (between backticks).

    This helps exclude function signatures like `func(param Type)`.
    """
    # Look backwards and forwards for backticks
    before = content[:match_start]
    after = content[match_end:]

    # Count backticks before and after the match on the same line
    line_start = before.rfind('\n') + 1
    line_end = after.find('\n')
    if line_end == -1:
        line_end = len(after)

    before_on_line = content[line_start:match_start]
    after_on_line = after[:line_end]

    # Count backticks before and after
    backticks_before = before_on_line.count('`')
    backticks_after = after_on_line.count('`')

    # If odd number of backticks before and after, we're inside inline code
    return (backticks_before % 2 == 1) and (backticks_after % 2 == 1)


def is_in_code_block_comment(line_content):
    """
    Check if a line is a comment in a code block.

    Looks for common comment patterns:
    - // (Go, C++, JavaScript)
    - # (Python, Shell, Ruby)
    - /* or */ (C-style block comments)
    """
    stripped = line_content.strip()
    return (
        stripped.startswith('//')
        or stripped.startswith('#')
        or stripped.startswith('/*')
        or '*/' in stripped
        or _RE_CODE_COMMENT.match(line_content)
    )


def has_strikethrough(line_content):
    """
    Check if a line contains any strikethrough text (~~text~~).

    Used to skip validation of links on lines with strikethrough,
    which typically indicate obsolete/deprecated requirements.
    """
    return '~~' in line_content


def extract_links(file_path, file_cache=None):
    """
    Extract all markdown links from a file, including those in code comments.

    Strategy:
    1. Find all markdown links in the file
    2. Exclude links that are in inline code (between backticks)
    3. Include links in code block comments
    4. Exclude links in code blocks that aren't in comments (function signatures)

    Args:
        file_path: Path to the file to extract links from
        file_cache: Optional FileContentCache instance to use for reading files

    Returns:
        List of tuples: (link_text, link_target, line_number)
    """
    links = []
    try:
        # Use cache if provided, otherwise read directly
        if file_cache:
            content = file_cache.get_content(Path(file_path))
        else:
            with open(file_path, 'r', encoding='utf-8') as f:
                content = f.read()

        lines = content.split('\n')

        # Track code blocks
        in_code_block = False

        for match in _RE_MARKDOWN_LINK.finditer(content):
            link_text = match.group(1)
            link_target = match.group(2)
            match_start = match.start()
            match_end = match.end()

            # Skip external links
            if link_target.startswith(('http://', 'https://', 'mailto:')):
                continue

            # Calculate line number
            line_num = content[:match_start].count('\n') + 1
            line_content = lines[line_num - 1] if line_num <= len(lines) else ""

            # Determine if we're in a code block at this line
            in_code_block = False
            for i in range(line_num - 1):
                if lines[i].strip().startswith('```'):
                    in_code_block = not in_code_block

            # Check if this match is in inline code
            if is_in_inline_code(content, match_start, match_end):
                continue

            # Skip links on lines with strikethrough (obsolete/deprecated requirements)
            if has_strikethrough(line_content):
                continue

            # If in code block, only include if it's in a comment
            if in_code_block:
                if not is_in_code_block_comment(line_content):
                    continue

            # Additional filter: skip obvious function parameters
            # (no path separator, contains spaces or special chars like 'error, key')
            if not _RE_PATH_SEP.search(link_target) and _RE_SPACE_COMMA.search(link_target):
                continue

            links.append((link_text, link_target, line_num))

    except (IOError, OSError) as e:
        # File read errors - log to stderr
        print(f"Error reading {file_path}: {e}", file=sys.stderr)
    except UnicodeDecodeError as e:
        # Encoding errors - log to stderr
        print(f"Error decoding {file_path} (encoding issue): {e}", file=sys.stderr)
    return links
