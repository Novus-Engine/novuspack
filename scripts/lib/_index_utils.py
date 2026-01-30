from __future__ import annotations

import re
from dataclasses import dataclass, field
from typing import Dict, List, Literal, Optional, Set, TYPE_CHECKING

from lib._go_code_utils import normalize_generic_name
from lib._validation_utils import ProseSection, generate_anchor_from_heading

if TYPE_CHECKING:
    pass

SectionKind = Literal["type", "method", "func"]

_SECTION_NUMBER_RE = re.compile(r"^\d+(?:\.\d+)*$")
_INDEX_HEADING_RE = re.compile(r"^(#{1,6})\s+(.+)$")
_INDEX_SECTION_NUMBER_RE = re.compile(r"^(\d+(?:\.\d+)*)(?:\.)?\s+(.+)$")
_INDEX_ENTRY_RE = re.compile(r"^\s*-\s+\*\*`([^`]+)`\*\*")
_INDEX_LINK_RE = re.compile(r"\[([^\]]+)\]\(([^)#]+)(?:#([^)]+))?\)")
_TITLE_HEADING_RE = re.compile(r"^#\s+(.+)$")


def _split_words_lower(text: str) -> list[str]:
    # Keep this conservative and ASCII-focused.
    # We split on non-alphanumeric boundaries and drop empties.
    return [w.lower() for w in re.split(r"[^A-Za-z0-9]+", text) if w]


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

    def sort_key(self) -> str:
        if "." in self.name:
            return self.name.split(".", 1)[1].lower()
        return self.name.lower()


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

        def sort_key(name: str) -> tuple[str, int, str]:
            if not name:
                return ("", 1, "")
            first = name[0]
            return (name.lower(), 0 if first.isupper() else 1, name)

        sorted_names = sorted(self.expected_entries.keys(), key=sort_key)
        self.expected_entries = {name: self.expected_entries[name] for name in sorted_names}

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
                "section_number must be a dotted number like '1', '2.4', or '3.5.6', got: %r"
                % (value,)
            )

    @staticmethod
    def _validate_heading_level(value: int) -> None:
        if value not in (2, 3, 4):
            raise ValueError("heading_level must be 2, 3, or 4, got: %r" % (value,))

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
                "(parent=%r, child=%r)"
                % (parent_heading.heading_level, heading_level)
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


def _extract_entry_description(
    lines: list[str],
    entry_line_num: int,
) -> tuple[list[str], Optional[int], Optional[int]]:
    description_lines: List[str] = []
    description_start_line = None
    description_end_line = None

    i = entry_line_num
    while i < len(lines):
        next_line = lines[i]
        stripped = next_line.rstrip()

        if _INDEX_ENTRY_RE.match(stripped):
            break
        if re.match(r"^##+\s+", stripped):
            break
        if stripped and not (stripped.startswith("  ") or stripped.startswith("    ")):
            if not stripped.strip():
                i += 1
                continue
            break

        if stripped.startswith("  - ") or stripped.startswith("    - "):
            if description_start_line is None:
                description_start_line = i + 1
            bullet_text = stripped.split("- ", 1)[1].strip()
            if bullet_text:
                description_lines.append(bullet_text)
                description_end_line = i + 1
        elif (stripped.startswith("  ") or stripped.startswith("    ")) and description_lines:
            # Continuation of previous bullet line.
            description_lines[-1] += " " + stripped.lstrip().strip()
            description_end_line = i + 1

        i += 1

    return description_lines, description_start_line, description_end_line


def parse_entry_descriptions(
    index_content: str,
    entry_line_numbers: Dict[str, int],
) -> Dict[str, tuple[bool, Optional[str], Optional[int]]]:
    """
    Parse descriptive text (indented bullets) for each index entry.

    Args:
        index_content: Full content of the index file
        entry_line_numbers: Dict mapping normalized entry name -> line number

    Returns:
        Dict mapping normalized entry name ->
        (has_description, description_text, description_line_start)
    """
    lines = index_content.split("\n")
    descriptions: Dict[str, tuple[bool, Optional[str], Optional[int]]] = {}
    entry_pattern = r"^\s*-\s+\*\*`([^`]+)`\*\*"

    line_to_entry: Dict[int, str] = {
        line_num: name for name, line_num in entry_line_numbers.items()
    }

    for line_num, line in enumerate(lines, 1):
        if line_num not in line_to_entry:
            continue

        entry_name = line_to_entry[line_num]
        description_lines = []
        description_start_line = None

        i = line_num
        while i < len(lines):
            next_line = lines[i]
            stripped = next_line.rstrip()

            if re.match(entry_pattern, stripped):
                break
            if re.match(r"^##+\s+", stripped):
                break
            if stripped and not (stripped.startswith("  ") or stripped.startswith("    ")):
                if not stripped.strip():
                    i += 1
                    continue
                break

            if stripped.startswith("  - ") or stripped.startswith("    - "):
                if description_start_line is None:
                    description_start_line = i
                bullet_text = stripped.split("- ", 1)[1].strip()
                if bullet_text:
                    description_lines.append(bullet_text)
            elif (stripped.startswith("  ") or stripped.startswith("    ")) and description_lines:
                if description_lines:
                    description_lines[-1] += " " + stripped.lstrip().strip()

            i += 1

        if description_lines:
            description_text = " ".join(description_lines).strip()
            has_description = len(description_text) >= 20
            descriptions[entry_name] = (
                has_description,
                description_text if has_description else None,
                description_start_line,
            )
        else:
            descriptions[entry_name] = (False, None, None)

    return descriptions


def _parse_overview_section(
    lines: list[str],
    end_line: int,
) -> Optional[ProseSection]:
    headings: List[tuple[int, str, int]] = []
    for line_num, line in enumerate(lines[:end_line], 1):
        match = _INDEX_HEADING_RE.match(line)
        if not match:
            continue
        level = len(match.group(1))
        heading_text = match.group(2).strip()
        if level < 2:
            continue
        headings.append((level, heading_text, line_num))

    if not headings:
        return None

    root: Optional[ProseSection] = None
    stack: List[ProseSection] = []
    nodes_by_line: Dict[int, ProseSection] = {}

    def _split_heading(heading_text: str) -> tuple[str, Optional[str]]:
        num_match = _INDEX_SECTION_NUMBER_RE.match(heading_text)
        if not num_match:
            return heading_text.strip(), None
        return num_match.group(2).strip(), num_match.group(1).strip()

    for level, heading_text, line_num in headings:
        heading_str, heading_num = _split_heading(heading_text)
        node = ProseSection(
            heading_str=heading_str,
            heading_num=heading_num,
            heading_level=level,
            heading_line=line_num,
            content="",
        )
        while stack and stack[-1].heading_level >= level:
            stack.pop()
        if stack:
            parent = stack[-1]
            node.parent_section = parent
            parent.child_sections.append(node)
        elif root is None:
            root = node
        else:
            node.parent_section = root
            root.child_sections.append(node)
        stack.append(node)
        nodes_by_line[line_num] = node

    for idx, (_level, _heading_text, line_num) in enumerate(headings):
        node_line_start = line_num + 1
        node_line_end = end_line if idx + 1 >= len(headings) else headings[idx + 1][2] - 1
        node_lines = lines[node_line_start - 1:node_line_end]
        content = "\n".join(node_lines).strip()
        node = nodes_by_line.get(line_num)
        if node is None:
            continue
        node.content = content
        node.lines = (node_line_start, node_line_end)

        in_code = False
        code_start = None
        for offset, line in enumerate(node_lines, node_line_start):
            stripped = line.strip()
            if stripped.startswith("```"):
                if not in_code:
                    in_code = True
                    code_start = offset
                else:
                    in_code = False
                    code_type = stripped[3:].strip().split()[0] if stripped[3:].strip() else ""
                    node.code_blocks.append((code_start or offset, offset, code_type))
                    node.has_code = True

    return root


def parse_index(index_content: str) -> ParsedIndex:
    """
    Parse docs/tech_specs/api_go_defs_index.md into structured sections and entries.
    """
    sections: Dict[str, IndexSection] = {}
    section_order: List[str] = []
    section_path_lines: Dict[str, List[int]] = {}

    current_h2: Optional[IndexSection] = None
    current_h3: Optional[IndexSection] = None
    current_h4: Optional[IndexSection] = None

    title = ""
    lines = index_content.split("\n")
    for line in lines:
        title_match = _TITLE_HEADING_RE.match(line)
        if title_match:
            title = title_match.group(1).strip()
            break

    first_def_section_line = len(lines)
    for line_num, line in enumerate(lines, 1):
        heading_match = _INDEX_HEADING_RE.match(line)
        if not heading_match:
            continue
        hashes = heading_match.group(1)
        heading_level = len(hashes)
        if heading_level not in (2, 3, 4):
            continue
        raw_heading = heading_match.group(2).strip()
        n = _INDEX_SECTION_NUMBER_RE.match(raw_heading)
        if not n:
            continue
        section_number = n.group(1).strip()
        if section_number.split(".")[0] == "0":
            continue
        first_def_section_line = line_num
        break

    overview = _parse_overview_section(lines, first_def_section_line - 1)
    if overview:
        overview.file_path = "api_go_defs_index.md"

    for line_num, line in enumerate(lines, 1):
        heading_match = _INDEX_HEADING_RE.match(line)
        if heading_match:
            hashes = heading_match.group(1)
            heading_level = len(hashes)
            if heading_level not in (2, 3, 4):
                continue

            raw_heading = heading_match.group(2).strip()
            n = _INDEX_SECTION_NUMBER_RE.match(raw_heading)
            if not n:
                continue

            section_number = n.group(1).strip()
            if section_number.split(".")[0] == "0":
                continue
            heading_text = n.group(2).strip()

            if heading_level == 2:
                parent = None
                current_h2 = None
                current_h3 = None
                current_h4 = None
            elif heading_level == 3:
                parent = current_h2
                current_h3 = None
                current_h4 = None
            else:
                parent = current_h3
                current_h4 = None

            if parent is None and heading_level != 2:
                continue

            heading_kind = IndexSection.derive_heading_kind(heading_text)
            node = IndexSection(
                section_number=section_number,
                heading_level=heading_level,
                parent_heading=parent,
                heading_text=heading_text,
                kind=heading_kind,
            )
            if parent is not None:
                parent.add_child(node)

            section_path = node.path_label()
            section_path_lines.setdefault(section_path, []).append(line_num)
            if section_path in sections:
                continue

            sections[section_path] = node
            section_order.append(section_path)

            if heading_level == 2:
                current_h2 = node
            elif heading_level == 3:
                current_h3 = node
            else:
                current_h4 = node
            continue

        entry_match = _INDEX_ENTRY_RE.match(line)
        if not entry_match:
            continue

        raw_name = entry_match.group(1)
        name = normalize_generic_name(raw_name)

        section_node = current_h4 or current_h3 or current_h2
        if section_node is None:
            continue

        link_match = _INDEX_LINK_RE.search(line)
        link_text = ""
        link_file = ""
        link_anchor: Optional[str] = None
        if link_match:
            link_text = link_match.group(1).strip()
            link_file = link_match.group(2).strip()
            anchor = link_match.group(3)
            link_anchor = anchor.strip() if anchor else None

        entry = IndexEntry(
            name=name,
            raw_name=raw_name,
            current_section=section_node.path_label(),
            link_text=link_text,
            link_file=link_file,
            link_anchor=link_anchor,
            line_number=line_num,
            kind=section_node.kind,
        )
        section_node.add_entry(entry)
        section_node.current_entries[name] = entry

    duplicates = {path: lines for path, lines in section_path_lines.items() if len(lines) > 1}
    if duplicates:
        details = ", ".join(
            f"{path} (lines {', '.join(str(line_num) for line_num in dup_lines)})"
            for path, dup_lines in sorted(duplicates.items())
        )
        raise ValueError(f"Duplicate headings detected in index file: {details}")

    for section_path in section_order:
        section = sections[section_path]
        for entry in section.entries:
            desc_lines, desc_start, desc_end = _extract_entry_description(
                lines, entry.line_number
            )
            entry.description_lines = desc_lines
            entry.description_line_start = desc_start
            entry.description_line_end = desc_end
            if desc_lines:
                desc_text = " ".join(desc_lines).strip()
                entry.description_text = desc_text
                entry.has_description = len(desc_text) >= 20

    unsorted_types = IndexSection(
        section_number="0",
        heading_level=2,
        parent_heading=None,
        heading_text="Unsorted Types",
        kind="type",
    )
    unsorted_methods = IndexSection(
        section_number="0",
        heading_level=2,
        parent_heading=None,
        heading_text="Unsorted Methods",
        kind="method",
    )
    unsorted_funcs = IndexSection(
        section_number="0",
        heading_level=2,
        parent_heading=None,
        heading_text="Unsorted Functions",
        kind="func",
    )
    unsorted_paths = [
        unsorted_types.path_label(),
        unsorted_methods.path_label(),
        unsorted_funcs.path_label(),
    ]
    sections[unsorted_paths[0]] = unsorted_types
    sections[unsorted_paths[1]] = unsorted_methods
    sections[unsorted_paths[2]] = unsorted_funcs

    return ParsedIndex(
        sections=sections,
        section_order=section_order,
        overview=overview,
        unsorted_paths=unsorted_paths,
        title=title,
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
        lines: List[str] = []
        for section_path in self.section_order:
            section = self.sections.get(section_path)
            if section is None:
                continue
            lines.append(section_path)
            entries: Dict[str, IndexEntry] = dict[str, IndexEntry](section.expected_entries)
            for name, entry in section.current_entries.items():
                if name in entries:
                    continue
                if entry.entry_status in ("orphaned", "removed"):
                    entries[name] = entry
            for name in sorted(entries.keys()):
                entry = entries[name]
                marker = ""
                if entry.entry_status == "added":
                    marker = " [ADDED]"
                elif entry.entry_status == "moved":
                    marker = " [MOVED]"
                elif entry.entry_status == "unresolved":
                    marker = " [UNRESOLVED]"
                elif entry.entry_status == "orphaned":
                    marker = " [ORPHANED]"
                elif entry.entry_status == "removed":
                    marker = " [REMOVED]"
                lines.append(f"- {name}{marker}")
            lines.append("")
        return lines

    def to_markdown(self) -> str:
        lines: List[str] = []
        title = self.title.strip() if self.title else ""
        if title:
            lines.append(f"# {title}")
            lines.append("")

        toc_lines = self._render_toc()
        if toc_lines:
            lines.extend(toc_lines)
            lines.append("")

        if self.overview:
            lines.extend(self._render_prose_section(self.overview))

        for section_path in self.section_order:
            section = self.sections.get(section_path)
            if section is None:
                continue
            lines.extend(self._render_section_markdown(section))

        return "\n".join(lines).rstrip() + "\n"

    def _render_toc(self) -> List[str]:
        toc_lines: List[str] = []
        if self.overview:
            label = self._format_prose_heading(self.overview)
            if label:
                toc_lines.append(
                    f"- [{label}]({generate_anchor_from_heading(label, include_hash=True)})"
                )
        for section_path in self.section_order:
            section = self.sections.get(section_path)
            if section is None:
                continue
            label = section.heading_label()
            indent = max(section.heading_level - 2, 0) * 2
            prefix = " " * indent
            anchor = generate_anchor_from_heading(label, include_hash=True)
            toc_lines.append(f"{prefix}- [{label}]({anchor})")
        return toc_lines

    def _render_prose_section(self, section: ProseSection) -> List[str]:
        lines: List[str] = []
        heading_label = self._format_prose_heading(section)
        if heading_label:
            lines.append(f"{'#' * section.heading_level} {heading_label}")
            lines.append("")
        if section.content:
            lines.extend(section.content.splitlines())
            lines.append("")
        for child in section.child_sections:
            lines.extend(self._render_prose_section(child))
        return lines

    @staticmethod
    def _format_prose_heading(section: ProseSection) -> str:
        if section.heading_num:
            if section.heading_level == 2:
                return f"{section.heading_num}. {section.heading_str}"
            return f"{section.heading_num} {section.heading_str}"
        return section.heading_str

    def _render_section_markdown(self, section: IndexSection) -> List[str]:
        lines = [f"{'#' * section.heading_level} {section.heading_label()}"]
        lines.append("")
        for entry in section.expected_entries.values():
            current_entry = self.find_current_entry(entry.name)
            raw_name = entry.raw_name
            link_text = entry.link_text
            link_target = entry.link_target()
            if current_entry:
                raw_name = current_entry.raw_name or raw_name
                link_text = current_entry.link_text or link_text
                link_target = current_entry.link_target()
                if current_entry.needs_link_update:
                    link_target = entry.link_target()
            link_label = link_text or "Spec"
            lines.append(f"- **`{raw_name}`** - [{link_label}]({link_target})")
            if entry.description_lines:
                for desc_line in entry.description_lines:
                    if desc_line.startswith("CONT: "):
                        lines.append(f"    {desc_line[len('CONT: '):]}")
                    else:
                        lines.append(f"  - {desc_line}")
        lines.append("")
        return lines
