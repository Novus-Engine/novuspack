from __future__ import annotations

import re
from dataclasses import dataclass, field
from typing import Dict, List, Literal, Optional, Set, TYPE_CHECKING

from lib._validation_utils import ProseSection
from lib import _index_utils_parsing
from lib import _index_utils_rendering

if TYPE_CHECKING:
    pass

SectionKind = Literal["type", "method", "func"]

_SECTION_NUMBER_RE = re.compile(r"^\d+(?:\.\d+)*$")


def _split_words_lower(text: str) -> list[str]:
    # Keep this conservative and ASCII-focused.
    # We split on non-alphanumeric boundaries and drop empties.
    return [w.lower() for w in re.split(r"[^A-Za-z0-9]+", text) if w]


def _entry_sort_key(name: str) -> tuple[int, str, str]:
    if not name:
        return (1, "", "")
    first = name[0]
    return (0 if first.isupper() else 1, name.lower(), name)


@dataclass(slots=True)
class IndexEntry:
    """An entry in the index file."""

    name: str  # Normalized
    raw_name: str
    current_section: str
    link_text: str
    link_file: str
    link_anchor: Optional[str]
    line_number: int
    has_description: bool = False
    description_text: Optional[str] = None
    description_line_start: Optional[int] = None
    description_line_end: Optional[int] = None
    description_lines: list[str] = field(default_factory=list)
    kind: Optional[SectionKind] = None
    entry_status: Optional[str] = None
    def_comments: Optional[str] = None
    expected_link_file: Optional[str] = None
    expected_link_anchor: Optional[str] = None
    needs_link_update: bool = False
    source_file: Optional[str] = None
    source_line: Optional[int] = None
    confidence_score: Optional[float] = None
    confidence_reasoning: list[str] = field(default_factory=list)
    suggested_section: Optional[str] = None

    def link_target(self) -> str:
        if self.link_anchor:
            return f"{self.link_file}#{self.link_anchor}"
        return self.link_file

    def sort_key(self) -> tuple[int, str, str]:
        if "." in self.name:
            method_name = self.name.split(".", 1)[1]
            return _entry_sort_key(method_name)
        return _entry_sort_key(self.name)


@dataclass(slots=True)
class IndexSection:
    """
    Represents a single indexed heading section in docs/tech_specs/api_go_defs_index.md.

    This is intentionally small and strict so callers can rely on normalized data.
    """

    section_number: str
    heading_level: int
    parent_heading: Optional[IndexSection]
    heading_text: str
    kind: SectionKind
    keywords_strong: list[str] = field(default_factory=list)
    keywords_med: list[str] = field(default_factory=list)
    keywords_weak: list[str] = field(default_factory=list)
    related_types: list[str] = field(default_factory=list)
    entries: list[IndexEntry] = field(default_factory=list)
    children: list[IndexSection] = field(default_factory=list)
    # Types valid for this section (populated by matching: type sections = entry names;
    # method sections = receiver from heading). Used for function-type association weight.
    valid_types: Set[str] = field(default_factory=set)
    # Current index entries for this section.
    current_entries: Dict[str, IndexEntry] = field(default_factory=dict)
    # Expected entries for this section.
    expected_entries: Dict[str, IndexEntry] = field(default_factory=dict)

    # Derived / normalized.
    heading_words: list[str] = field(init=False)

    def __post_init__(self) -> None:
        self.section_number = self.section_number.strip()
        # Normalize "1." => "1" so callers can compare consistently.
        if self.section_number.endswith("."):
            self.section_number = self.section_number[:-1]

        self._validate_section_number(self.section_number)
        self._validate_heading_level(self.heading_level)
        self._validate_parent_heading(self.heading_level, self.parent_heading)

        self.heading_text = self.heading_text.strip()
        if not self.heading_text:
            raise ValueError("heading_text must be non-empty")

        self.heading_words = _split_words_lower(self.heading_text)

        self.keywords_strong = self._normalize_lower_list(self.keywords_strong)
        self.keywords_med = self._normalize_lower_list(self.keywords_med)
        self.keywords_weak = self._normalize_lower_list(self.keywords_weak)
        self.related_types = self._normalize_lower_list(self.related_types)

    def heading_label(self) -> str:
        if self.heading_level == 2:
            return f"{self.section_number}. {self.heading_text}"
        return f"{self.section_number} {self.heading_text}"

    def path_label(self) -> str:
        parts: list[str] = []
        cur: Optional[IndexSection] = self
        while cur is not None:
            parts.append(cur.heading_label())
            cur = cur.parent_heading
        parts.reverse()
        return " > ".join(parts)

    def add_entry(self, entry: IndexEntry) -> None:
        self.entries.append(entry)

    def extend_entries(self, entries: list[IndexEntry]) -> None:
        self.entries.extend(entries)

    def add_child(self, child: IndexSection) -> None:
        self.children.append(child)

    def sort_entries(self) -> None:
        self.entries.sort(key=lambda entry: entry.sort_key())

    def sort_expected_entries(self) -> None:
        if not self.expected_entries:
            return

        sorted_items = sorted(
            self.expected_entries.items(),
            key=lambda item: item[1].sort_key(),
        )
        self.expected_entries = dict(sorted_items)

    def iter_sections(self) -> list[IndexSection]:
        ordered = [self]
        for child in self.children:
            ordered.extend(child.iter_sections())
        return ordered

    def to_markdown(self) -> list[str]:
        lines = [f"{'#' * self.heading_level} {self.heading_label()}"]
        lines.append("")
        for entry in self.entries:
            link_text = entry.link_text or "Spec"
            lines.append(f"- **`{entry.raw_name}`** - [{link_text}]({entry.link_target()})")
            if entry.description_lines:
                for desc_line in entry.description_lines:
                    if desc_line.startswith("CONT: "):
                        lines.append(f"    {desc_line[len('CONT: '):]}")
                    else:
                        lines.append(f"  - {desc_line}")
        lines.append("")
        return lines

    @staticmethod
    def derive_heading_kind(heading_text: str) -> SectionKind:
        """
        Derive heading kind from the heading text using basic matching.

        Rules:
        - If the heading contains the word "types" (plural), it is a "type" heading.
        - Else if the heading contains the word "functions" (plural), it is a "func" heading.
        - Otherwise it is a "method" heading.
        """
        words = _split_words_lower(heading_text)
        if "types" in words:
            return "type"
        if "functions" in words:
            return "func"
        return "method"

    @staticmethod
    def _validate_section_number(value: str) -> None:
        if not _SECTION_NUMBER_RE.match(value):
            raise ValueError(
                "section_number must be a dotted number like '1', '2.4', or '3.5.6', got: "
                f"{value!r}"
            )

    @staticmethod
    def _validate_heading_level(value: int) -> None:
        if value not in (2, 3, 4):
            raise ValueError(f"heading_level must be 2, 3, or 4, got: {value!r}")

    @staticmethod
    def _validate_parent_heading(
        heading_level: int, parent_heading: Optional[IndexSection]
    ) -> None:
        if heading_level == 2:
            if parent_heading is not None:
                raise ValueError("H2 headings must have parent_heading=None")
            return

        if parent_heading is None:
            raise ValueError("H3/H4 headings must have a non-None parent_heading")

        if parent_heading.heading_level >= heading_level:
            raise ValueError(
                "parent_heading.heading_level must be less than heading_level "
                f"(parent={parent_heading.heading_level!r}, child={heading_level!r})"
            )

    @staticmethod
    def _normalize_lower_list(values: list[str]) -> list[str]:
        out: list[str] = []
        for v in values:
            s = v.strip().lower()
            if not s:
                continue
            out.append(s)
        return out


def parse_entry_descriptions(
    index_content: str,
    entry_line_numbers: Dict[str, int],
) -> Dict[str, tuple[bool, Optional[str], Optional[int]]]:
    return _index_utils_parsing.parse_entry_descriptions(index_content, entry_line_numbers)


def parse_index(index_content: str) -> "ParsedIndex":
    return _index_utils_parsing.parse_index(
        index_content,
        index_section_cls=IndexSection,
        index_entry_cls=IndexEntry,
        parsed_index_cls=ParsedIndex,
    )


@dataclass(slots=True)
class ParsedIndex:
    """
    Parsed representation of docs/tech_specs/api_go_defs_index.md,
    representing the entire index file as a tree of sections and entries.
    Tracks both the current and expected state of the index file,
    with the expected state being the desired state of the index file.
    """

    sections: dict[str, IndexSection]
    section_order: list[str]
    overview: Optional[ProseSection]
    unsorted_paths: list[str]
    title: str

    def find_section_by_current_entry(self, name: str) -> Optional[IndexSection]:
        for section_path in self.section_order:
            section = self.sections.get(section_path)
            if section and name in section.current_entries:
                return section
        return None

    def find_section_by_expected_entry(self, name: str) -> Optional[IndexSection]:
        for section_path in self.section_order:
            section = self.sections.get(section_path)
            if section and name in section.expected_entries:
                return section
        return None

    def find_current_entry(self, name: str) -> Optional[IndexEntry]:
        section = self.find_section_by_current_entry(name)
        if not section:
            return None
        return section.current_entries.get(name)

    def get_orphans(self) -> List[IndexEntry]:
        orphans: List[IndexEntry] = []
        for section_path in self.section_order:
            section = self.sections.get(section_path)
            if not section:
                continue
            for entry in section.current_entries.values():
                if entry.entry_status == "orphaned":
                    orphans.append(entry)
        return orphans

    def get_added_entries(self) -> List[IndexEntry]:
        added: List[IndexEntry] = []
        for section_path in self.section_order:
            section = self.sections.get(section_path)
            if not section:
                continue
            for entry in section.expected_entries.values():
                if entry.entry_status == "added":
                    added.append(entry)
        return added

    def get_moved_entries(self) -> List[IndexEntry]:
        moved: List[IndexEntry] = []
        for section_path in self.section_order:
            section = self.sections.get(section_path)
            if not section:
                continue
            for entry in section.expected_entries.values():
                if entry.entry_status == "moved":
                    moved.append(entry)
        return moved

    def get_removed_entries(self) -> List[IndexEntry]:
        removed: List[IndexEntry] = []
        for section_path in self.section_order:
            section = self.sections.get(section_path)
            if not section:
                continue
            for entry in section.current_entries.values():
                if entry.entry_status == "removed":
                    removed.append(entry)
        return removed

    def get_link_update_entries(self) -> List[IndexEntry]:
        updates: List[IndexEntry] = []
        for section_path in self.section_order:
            section = self.sections.get(section_path)
            if not section:
                continue
            for entry in section.current_entries.values():
                if entry.needs_link_update:
                    updates.append(entry)
        return updates

    def get_reordered_entries(self) -> List[IndexEntry]:
        reordered: List[IndexEntry] = []
        for section_path in self.section_order:
            section = self.sections.get(section_path)
            if not section:
                continue
            for entry in section.current_entries.values():
                if entry.entry_status == "reordered":
                    reordered.append(entry)
            for entry in section.expected_entries.values():
                if entry.entry_status == "reordered":
                    reordered.append(entry)
        return reordered

    def get_unresolved_entries(self) -> List[IndexEntry]:
        unresolved: List[IndexEntry] = []
        for section_path in self.unsorted_paths:
            section = self.sections.get(section_path)
            if not section:
                continue
            for entry in section.expected_entries.values():
                if entry.entry_status == "unresolved":
                    unresolved.append(entry)
        return unresolved

    def sort_expected_entries(self) -> None:
        for section_path in self.section_order:
            section = self.sections.get(section_path)
            if not section:
                continue
            section.sort_expected_entries()

    def sync_expected_descriptions(self) -> None:
        description_map: Dict[str, List[str]] = {}
        for section_path in self.section_order:
            section = self.sections.get(section_path)
            if not section:
                continue
            for entry in section.current_entries.values():
                if entry.description_lines:
                    description_map[entry.name] = entry.description_lines
        for section_path in self.section_order:
            section = self.sections.get(section_path)
            if not section:
                continue
            for entry in section.expected_entries.values():
                if entry.description_lines:
                    continue
                if entry.name in description_map:
                    entry.description_lines = list(description_map[entry.name])

    def render_full_tree(self) -> List[str]:
        return _index_utils_rendering.render_full_tree(self)

    def to_markdown(self) -> str:
        return _index_utils_rendering.index_to_markdown(self)

    def _render_toc(self) -> List[str]:
        return _index_utils_rendering.render_toc(self)

    def _render_prose_section(self, section: ProseSection) -> List[str]:
        return _index_utils_rendering.render_prose_section(self, section)

    @staticmethod
    def _format_prose_heading(section: ProseSection) -> str:
        return _index_utils_rendering.format_prose_heading(section)

    def _render_section_markdown(self, section: IndexSection) -> List[str]:
        return _index_utils_rendering.render_section_markdown(self, section)
