from __future__ import annotations

import re
from typing import Dict, List, Optional, TYPE_CHECKING, Type

from lib._go_code_utils import normalize_generic_name
from lib._validation_utils import ProseSection

if TYPE_CHECKING:
    from lib._index_utils import IndexEntry, IndexSection, ParsedIndex

_INDEX_HEADING_RE = re.compile(r"^(#{1,6})\s+(.+)$")
_INDEX_SECTION_NUMBER_RE = re.compile(r"^(\d+(?:\.\d+)*)(?:\.)?\s+(.+)$")
_INDEX_ENTRY_RE = re.compile(r"^\s*-\s+\*\*`([^`]+)`\*\*")
_INDEX_LINK_RE = re.compile(r"\[([^\]]+)\]\(([^)#]+)(?:#([^)]+))?\)")
_TITLE_HEADING_RE = re.compile(r"^#\s+(.+)$")


def _is_entry_or_heading_line(stripped: str, entry_pattern: str) -> bool:
    if re.match(entry_pattern, stripped):
        return True
    return bool(re.match(r"^##+\s+", stripped))


def _is_description_line(stripped: str) -> bool:
    return stripped.startswith("  ") or stripped.startswith("    ")


def _collect_description_lines(
    lines: list[str],
    entry_line_num: int,
    entry_pattern: str,
) -> tuple[list[str], Optional[int], Optional[int]]:
    description_lines: List[str] = []
    description_start_line = None
    description_end_line = None

    i = entry_line_num
    while i < len(lines):
        stripped = lines[i].rstrip()

        if _is_entry_or_heading_line(stripped, entry_pattern):
            break
        if stripped and not _is_description_line(stripped):
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
        elif _is_description_line(stripped) and description_lines:
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
    """
    lines = index_content.split("\n")
    descriptions: Dict[str, tuple[bool, Optional[str], Optional[int]]] = {}
    entry_pattern = r"^\s*-\s+\*\*`([^`]+)`\*\*"

    line_to_entry: Dict[int, str] = {
        line_num: name for name, line_num in entry_line_numbers.items()
    }

    for line_num, _line in enumerate(lines, 1):
        if line_num not in line_to_entry:
            continue

        entry_name = line_to_entry[line_num]
        (
            description_lines,
            description_start_line,
            _description_end_line,
        ) = _collect_description_lines(lines, line_num, entry_pattern)

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
            return heading_text, None
        return num_match.group(2).strip(), num_match.group(1).strip()

    for level, heading_text, line_num in headings:
        heading_str, heading_num = _split_heading(heading_text)
        node = ProseSection(
            heading_str=heading_str,
            heading_num=heading_num,
            heading_level=level,
            heading_line=line_num,
            content="",
            parent_section=None,
            child_sections=[],
            has_code=False,
            code_blocks=[],
        )
        while stack and stack[-1].heading_level >= level:
            stack.pop()
        if stack:
            node.parent_section = stack[-1]
            stack[-1].child_sections.append(node)
        else:
            root = node
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


def _extract_title(lines: list[str]) -> str:
    for line in lines:
        title_match = _TITLE_HEADING_RE.match(line)
        if title_match:
            return title_match.group(1).strip()
    return ""


def _find_first_def_section_line(lines: list[str]) -> int:
    for line_num, line in enumerate(lines, 1):
        heading = _parse_heading_line(line)
        if not heading:
            continue
        section_number, _heading_text, heading_level = heading
        if heading_level not in (2, 3, 4):
            continue
        if section_number.split(".")[0] == "0":
            continue
        return line_num
    return len(lines)


def _parse_heading_line(line: str) -> Optional[tuple[str, str, int]]:
    heading_match = _INDEX_HEADING_RE.match(line)
    if not heading_match:
        return None
    hashes = heading_match.group(1)
    heading_level = len(hashes)
    if heading_level not in (2, 3, 4):
        return None
    raw_heading = heading_match.group(2).strip()
    n = _INDEX_SECTION_NUMBER_RE.match(raw_heading)
    if not n:
        return None
    section_number = n.group(1).strip()
    heading_text = n.group(2).strip()
    return section_number, heading_text, heading_level


def _update_heading_context(
    heading_level: int,
    node: "IndexSection",
    current_h2: Optional["IndexSection"],
    current_h3: Optional["IndexSection"],
) -> tuple[Optional["IndexSection"], Optional["IndexSection"], Optional["IndexSection"]]:
    if heading_level == 2:
        return node, None, None
    if heading_level == 3:
        return current_h2, node, None
    return current_h2, current_h3, node


def _parse_sections(
    lines: list[str],
    index_section_cls: Type["IndexSection"],
) -> tuple[
    Dict[str, "IndexSection"],
    List[str],
    Dict[str, List[int]],
    Dict[int, "IndexSection"],
]:
    sections: Dict[str, "IndexSection"] = {}
    section_order: List[str] = []
    section_path_lines: Dict[str, List[int]] = {}
    nodes_by_line: Dict[int, IndexSection] = {}

    current_h2: Optional[IndexSection] = None
    current_h3: Optional[IndexSection] = None

    for line_num, line in enumerate(lines, 1):
        heading = _parse_heading_line(line)
        if not heading:
            continue
        section_number, heading_text, heading_level = heading
        if section_number.split(".")[0] == "0":
            continue
        parent = None
        if heading_level == 3:
            parent = current_h2
        elif heading_level == 4:
            parent = current_h3
        if parent is None and heading_level != 2:
            continue

        heading_kind = index_section_cls.derive_heading_kind(heading_text)
        node = index_section_cls(
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
        nodes_by_line[line_num] = node
        if section_path in sections:
            continue

        sections[section_path] = node
        section_order.append(section_path)

        current_h2, current_h3, _current_h4 = _update_heading_context(
            heading_level,
            node,
            current_h2,
            current_h3,
        )

    return sections, section_order, section_path_lines, nodes_by_line


def _parse_entries(
    lines: list[str],
    nodes_by_line: Dict[int, "IndexSection"],
    index_entry_cls: Type["IndexEntry"],
) -> None:
    current_h2: Optional["IndexSection"] = None
    current_h3: Optional["IndexSection"] = None
    current_h4: Optional["IndexSection"] = None

    for line_num, line in enumerate(lines, 1):
        if line_num in nodes_by_line:
            node = nodes_by_line[line_num]
            current_h2, current_h3, current_h4 = _update_heading_context(
                node.heading_level,
                node,
                current_h2,
                current_h3,
            )
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

        entry = index_entry_cls(
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


def _validate_unique_headings(section_path_lines: Dict[str, List[int]]) -> None:
    duplicates = {path: lines for path, lines in section_path_lines.items() if len(lines) > 1}
    if not duplicates:
        return
    details = ", ".join(
        f"{path} (lines {', '.join(str(line_num) for line_num in dup_lines)})"
        for path, dup_lines in sorted(duplicates.items())
    )
    raise ValueError(f"Duplicate headings detected in index file: {details}")


def _populate_entry_description_fields(
    lines: list[str],
    section_order: List[str],
    sections: Dict[str, "IndexSection"],
) -> None:
    for section_path in section_order:
        section = sections[section_path]
        for entry in section.entries:
            desc_lines, desc_start, desc_end = _collect_description_lines(
                lines,
                entry.line_number,
                r"^\s*-\s+\*\*`([^`]+)`\*\*",
            )
            entry.description_lines = desc_lines
            entry.description_line_start = desc_start
            entry.description_line_end = desc_end
            if desc_lines:
                desc_text = " ".join(desc_lines).strip()
                entry.description_text = desc_text
                entry.has_description = len(desc_text) >= 20


def _add_unsorted_sections(
    sections: Dict[str, "IndexSection"],
    index_section_cls: Type["IndexSection"],
) -> list[str]:
    unsorted_types = index_section_cls(
        section_number="0",
        heading_level=2,
        parent_heading=None,
        heading_text="Unsorted Types",
        kind="type",
    )
    unsorted_methods = index_section_cls(
        section_number="0",
        heading_level=2,
        parent_heading=None,
        heading_text="Unsorted Methods",
        kind="method",
    )
    unsorted_funcs = index_section_cls(
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
    return unsorted_paths


def parse_index(
    index_content: str,
    *,
    index_section_cls: Type["IndexSection"],
    index_entry_cls: Type["IndexEntry"],
    parsed_index_cls: Type["ParsedIndex"],
) -> "ParsedIndex":
    """
    Parse docs/tech_specs/api_go_defs_index.md into structured sections and entries.
    """
    lines = index_content.split("\n")
    title = _extract_title(lines)
    first_def_section_line = _find_first_def_section_line(lines)
    overview = _parse_overview_section(lines, first_def_section_line - 1)
    if overview:
        overview.file_path = "api_go_defs_index.md"

    sections, section_order, section_path_lines, nodes_by_line = _parse_sections(
        lines,
        index_section_cls,
    )
    _parse_entries(lines, nodes_by_line, index_entry_cls)
    _validate_unique_headings(section_path_lines)
    _populate_entry_description_fields(lines, section_order, sections)
    unsorted_paths = _add_unsorted_sections(sections, index_section_cls)

    return parsed_index_cls(
        sections=sections,
        section_order=section_order,
        overview=overview,
        unsorted_paths=unsorted_paths,
        title=title,
    )
