"""
Heading numbering validation helpers.

Public surface for C0302 module splitting (validate_heading_numbering).
"""

from lib.heading_numbering._checks import (
    check_duplicate_headings,
    check_excessive_numbering,
    check_h2_period_consistency,
    check_heading_capitalization,
    check_organizational_headings,
    check_single_word_headings,
)

__all__ = [
    "check_duplicate_headings",
    "check_excessive_numbering",
    "check_h2_period_consistency",
    "check_heading_capitalization",
    "check_organizational_headings",
    "check_single_word_headings",
]
