from __future__ import annotations

from typing import Optional
from lib._index_utils import ParsedIndex
from lib._validation_utils import OutputBuilder


def determine_ordering(
    parsed_index: ParsedIndex,
    output: Optional[OutputBuilder] = None,
) -> None:
    """
    Phase 6: Determine correct ordering within sections.
    """
    _ = output
    parsed_index.sort_expected_entries()
