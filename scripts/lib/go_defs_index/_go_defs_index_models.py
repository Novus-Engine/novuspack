from __future__ import annotations

from dataclasses import dataclass, field
from typing import List, Optional, TYPE_CHECKING

from lib._index_utils import IndexEntry

if TYPE_CHECKING:
    pass


@dataclass
class DetectedDefinition:
    """A Go definition detected in tech spec files."""

    # Phase 1 fields
    name: str  # Normalized: "Package.Close", "Package"
    kind: str  # "type", "method", "func"
    file: str  # Source file: "api_core.md"
    code_block_start_line: int
    code_block_content: str
    raw_name: str  # Original name before normalization (for methods)
    receiver_type: Optional[str] = None  # For methods only
    input_types: List[str] = field(default_factory=list)
    output_types: List[str] = field(default_factory=list)
    referenced_types: List[str] = field(default_factory=list)
    referenced_methods: List[str] = field(default_factory=list)
    def_comments: Optional[str] = None

    # Phase 2 fields (added later)
    heading: str = ""
    heading_level: int = 0
    heading_line: int = 0
    parent_heading: Optional[str] = None
    parent_heading_level: Optional[int] = None
    section_content: str = ""
    current_section: Optional[str] = None  # Current section in index (if exists)
    canonical_file: str = ""
    canonical_heading: str = ""
    canonical_anchor: str = ""

    # Phase 3 fields (added later)
    suggested_section: Optional[str] = None
    confidence_score: float = 0.0  # 0.0 to 1.0
    confidence_reasoning: List[str] = field(default_factory=list)
    is_resolved: bool = False

    def __post_init__(self) -> None:
        if self.kind in ("method", "func"):
            return
        if self.kind in ("type", "struct", "interface", "alias"):
            self.kind = "type"
            return
        self.kind = "type"

    def to_index_entry(self, current_section: str) -> IndexEntry:
        link_text = self.canonical_heading or self.heading or self.raw_name
        def_comments = None
        if self.def_comments:
            def_comments = " ".join(self.def_comments.splitlines()).strip()
            if not def_comments:
                def_comments = None
        # IndexEntry.link_target() adds the '#' separator itself, so link_anchor must
        # NOT include a leading '#'. Canonical anchors are stored with '#' because
        # they come from markdown links / generate_anchor_from_heading(include_hash=True).
        link_anchor = (self.canonical_anchor or "").lstrip("#").strip() or None
        return IndexEntry(
            name=self.name,
            raw_name=self.raw_name,
            current_section=current_section,
            link_text=link_text,
            link_file=self.canonical_file,
            link_anchor=link_anchor,
            line_number=self.code_block_start_line,
            kind=self.kind,
            def_comments=def_comments,
            source_file=self.file,
            source_line=self.code_block_start_line,
            confidence_score=self.confidence_score,
            confidence_reasoning=list(self.confidence_reasoning),
        )
