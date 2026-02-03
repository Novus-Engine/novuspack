"""
Title-case and capitalization helpers for heading numbering validation.

Pure functions used by validate_heading_numbering to apply Title Case rules
to heading text (preserving backticks, filenames, code-in-parens).
Case inside backticks is not checked and is preserved as-is.
"""

import re

_RE_SPLIT_WORDS = re.compile(r'\S+|\s+')
_RE_WHITESPACE_ONLY = re.compile(r'^\s+$')
_RE_FILENAME_PATTERN = re.compile(r'^[\w\-]+\.(\w+)$')
_RE_NON_WORD_CHARS = re.compile(r'[^\w]')
_RE_FIRST_LETTER = re.compile(r'[a-zA-Z]')

_COMMON_EXTENSIONS = [
    'go', 'md', 'txt', 'json', 'yaml', 'yml', 'xml', 'html', 'css', 'js',
    'ts', 'py', 'sh', 'bat', 'ps1', 'java', 'c', 'cpp', 'h', 'hpp',
    'rs', 'rb', 'php', 'sql', 'csv', 'tsv', 'log', 'conf', 'config',
    'ini', 'toml', 'lock', 'sum', 'mod', 'gitignore', 'editorconfig'
]

_LOWERCASE_WORDS = {
    'a', 'an', 'the',
    'and', 'but', 'or', 'nor', 'so', 'yet',
    'in', 'on', 'at', 'to', 'for', 'of', 'with', 'by', 'from', 'up', 'about',
    'into', 'onto', 'upon', 'over', 'under', 'above', 'below', 'across',
    'via', 'vs'
}

_CAPITALIZE_PREPOSITIONS = {
    'through', 'between', 'among', 'during', 'before', 'after', 'within',
    'without', 'against', 'along', 'around', 'behind', 'beside', 'beyond',
    'inside', 'outside', 'throughout', 'toward', 'towards', 'underneath'
}

_PROGRAMMING_KEYWORDS = {
    'return', 'if', 'else', 'for', 'while', 'do', 'switch', 'case',
    'break', 'continue', 'goto', 'throw', 'try', 'catch', 'finally',
    'new', 'delete', 'this', 'super', 'static', 'const', 'let', 'var',
    'function', 'class', 'interface', 'enum', 'type', 'import', 'export',
    'async', 'await', 'yield', 'def', 'lambda', 'pass', 'raise', 'except'
}

_PHRASAL_PARTICLES = {'up', 'down', 'out', 'off', 'in', 'on', 'over', 'away', 'back'}

_PHRASAL_VERB_BASES = {
    'clean', 'set', 'look', 'pick', 'give', 'make', 'break', 'build', 'call',
    'check', 'close', 'come', 'cut', 'fill', 'get', 'go', 'grow', 'hang',
    'hold', 'keep', 'line', 'live', 'move', 'open', 'pull', 'put', 'show',
    'sign', 'stand', 'start', 'take', 'turn', 'wake', 'warm', 'wrap', 'bring',
    'carry', 'catch', 'do', 'draw', 'drop', 'end', 'fall', 'find',
    'fix', 'follow', 'hand', 'head', 'help', 'join', 'jump', 'knock', 'lay',
    'leave', 'let', 'lie', 'lock', 'log', 'mix', 'pass', 'pay', 'point',
    'pop', 'push', 'read', 'run', 'send', 'shut', 'sit', 'slow', 'sort',
    'speed', 'split', 'spread', 'step', 'stick', 'stop', 'switch', 'talk',
    'tear', 'think', 'throw', 'tie', 'try', 'use', 'walk', 'wash',
    'watch', 'wear', 'wind', 'work', 'write'
}


def find_backtick_ranges(text):
    """Find all backtick-enclosed sections and their positions."""
    backtick_ranges = []
    i = 0
    while i < len(text):
        if text[i] == '`':
            start = i
            i += 1
            while i < len(text) and text[i] != '`':
                i += 1
            if i < len(text):
                end = i + 1
                backtick_ranges.append((start, end))
                i = end
            else:
                break
        else:
            i += 1
    return backtick_ranges


def is_in_backticks(pos, backtick_ranges):
    """Check if a character position is inside backticks."""
    for start, end in backtick_ranges:
        if start <= pos < end:
            return True
    return False


def should_preserve_part(part, part_start, part_end, backtick_ranges):
    """Check if a part should be preserved as-is (backticks, underscores, filenames)."""
    part_in_backticks = any(
        is_in_backticks(pos, backtick_ranges)
        for pos in range(part_start, part_end)
    )
    if part_in_backticks:
        return True
    if '_' in part:
        return True
    filename_match = _RE_FILENAME_PATTERN.match(part)
    if filename_match:
        extension = filename_match.group(1).lower()
        if extension in _COMMON_EXTENSIONS:
            return True
    return False


def is_in_code_parentheses(_part, parts, part_index, text):
    """Check if a word is inside parentheses with code-like content (backticks)."""
    text_up_to_here = ''.join(parts[:part_index + 1])
    if '(' not in text_up_to_here:
        return False
    last_open_paren = text_up_to_here.rfind('(')
    if last_open_paren < 0:
        return False
    text_before = ''.join(parts[:part_index + 1])
    open_count = text_before.count('(') - text_before.count(')')
    if open_count <= 0:
        return False
    text_from_paren = text[last_open_paren:]
    close_paren_pos = text_from_paren.find(')', 1)
    if close_paren_pos > 0:
        parens_content = text_from_paren[1:close_paren_pos]
    else:
        parens_content = text_from_paren[1:]
    return '`' in parens_content


def should_capitalize_word(
    word_clean, is_first, is_last, *, is_in_code_parens, previous_word_clean=None
):
    """
    Determine if a word should be capitalized based on title case rules.
    """
    if is_in_code_parens and word_clean in _PROGRAMMING_KEYWORDS:
        return False
    if is_first or is_last:
        return True
    if word_clean in _CAPITALIZE_PREPOSITIONS:
        return True
    if word_clean in _PHRASAL_PARTICLES and previous_word_clean:
        if previous_word_clean in _PHRASAL_VERB_BASES:
            return True
    if word_clean not in _LOWERCASE_WORDS:
        return True
    return False


def capitalize_word(part, should_cap):
    """Apply capitalization to a word part."""
    if not part:
        return part
    match = _RE_FIRST_LETTER.search(part)
    if not match:
        return part
    first_letter_idx = match.start()
    if should_cap:
        return (
            part[:first_letter_idx] +
            part[first_letter_idx].upper() +
            part[first_letter_idx + 1:]
        )
    has_internal_capitals = any(
        c.isupper() for c in part[first_letter_idx + 1:]
        if c.isalpha()
    )
    if has_internal_capitals:
        if part[first_letter_idx].isupper():
            return (
                part[:first_letter_idx] +
                part[first_letter_idx].lower() +
                part[first_letter_idx + 1:]
            )
        return part
    return (
        part[:first_letter_idx] +
        part[first_letter_idx].lower() +
        part[first_letter_idx + 1:].lower()
    )


def to_title_case(text):
    """
    Convert text to Title Case following standard rules.
    Preserves backticks, filenames, code-in-parens.
    """
    if not text:
        return text
    backtick_ranges = find_backtick_ranges(text)
    parts = _RE_SPLIT_WORDS.findall(text)
    result_parts = []
    word_indices = []
    char_pos = 0
    part_positions = []
    for i, part in enumerate(parts):
        if not _RE_WHITESPACE_ONLY.match(part):
            word_indices.append(i)
        part_positions.append((char_pos, char_pos + len(part)))
        char_pos += len(part)

    if not word_indices:
        return text

    for i, part in enumerate(parts):
        if _RE_WHITESPACE_ONLY.match(part):
            result_parts.append(part)
            continue
        part_start, part_end = part_positions[i]
        if should_preserve_part(part, part_start, part_end, backtick_ranges):
            result_parts.append(part)
            continue
        word_clean = _RE_NON_WORD_CHARS.sub('', part.lower())
        if not word_clean:
            result_parts.append(part)
            continue
        # Preserve single-letter words that are uppercase in the original
        # (e.g. "Option A", "Plan B") so we do not suggest "Option a".
        if len(word_clean) == 1:
            first_letter_match = _RE_FIRST_LETTER.search(part)
            if first_letter_match and part[first_letter_match.start()].isupper():
                result_parts.append(part)
                continue
        is_first = i == word_indices[0]
        is_last = i == word_indices[-1]
        previous_word_clean = None
        if not is_first:
            current_word_index = word_indices.index(i)
            if current_word_index > 0:
                prev_word_i = word_indices[current_word_index - 1]
                prev_part = parts[prev_word_i]
                previous_word_clean = _RE_NON_WORD_CHARS.sub('', prev_part.lower())
        is_in_code_parens = is_in_code_parentheses(part, parts, i, text)
        should_cap = should_capitalize_word(
            word_clean, is_first, is_last,
            is_in_code_parens=is_in_code_parens,
            previous_word_clean=previous_word_clean
        )
        result_parts.append(capitalize_word(part, should_cap))
    return ''.join(result_parts)
