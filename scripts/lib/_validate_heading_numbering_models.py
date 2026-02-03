"""
Models and constants for heading numbering validation.
"""

import re

RE_HEADING_PATTERN = re.compile(r'^(#{1,})\s+(.+)$')
RE_NUMBERED_HEADING_PATTERN = re.compile(r'^([0-9]+(?:\.[0-9]+)*)\.?\s+(.+)$')

MAX_HEADING_NUMBER_SEGMENT = 20
MAX_ORGANIZATIONAL_PROSE_LINES = 5


class HeadingInfo:
    """Represents a heading with its metadata for sorting."""

    def __init__(self, file, line_num, level, numbers, *, heading_text, full_line,
                 parent=None, issue=None):
        self.file = file
        self.line_num = line_num
        self.level = level
        self.numbers = numbers
        self.heading_text = heading_text
        self.full_line = full_line
        self.parent = parent
        self.issue = issue
        self.original_number = '.'.join(map(str, numbers)) if numbers else ''
        self.corrected_number = None
        self.has_period = False
        self.corrected_capitalization = None

    def sort_key(self):
        """Return a sort key for proper numeric ordering."""
        return (tuple(self.numbers), self.level)
