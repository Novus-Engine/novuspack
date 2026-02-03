"""
Helpers for Go defs index matching adjustments.
"""

from __future__ import annotations

from typing import Optional

from lib.go_defs_index._go_defs_index_models import DetectedDefinition
from lib._index_utils import IndexSection


def adjust_related_section_for_function(
    defn: DetectedDefinition,
    related_section: Optional[IndexSection],
) -> Optional[IndexSection]:
    if not related_section:
        return related_section
    name_lower = defn.name.lower()
    if "fileentrytag" in name_lower:
        return None
    if defn.name == "NewPackageError":
        return None
    if "packagewithoptions" in name_lower:
        return None
    if name_lower in {"readheaderfrompath", "setdestpath"}:
        return None
    return related_section
