from __future__ import annotations

from typing import List

from lib._validation_utils import ProseSection, generate_anchor_from_heading


def format_prose_heading(section: ProseSection) -> str:
    if section.heading_num:
        if section.heading_level == 2:
            return f"{section.heading_num}. {section.heading_str}"
        return f"{section.heading_num} {section.heading_str}"
    return section.heading_str


def render_toc(parsed_index) -> List[str]:
    toc_lines: List[str] = []
    if parsed_index.overview:
        label = format_prose_heading(parsed_index.overview)
        if label:
            toc_lines.append(
                f"- [{label}]({generate_anchor_from_heading(label, include_hash=True)})"
            )
    for section_path in parsed_index.section_order:
        section = parsed_index.sections.get(section_path)
        if section is None:
            continue
        label = section.heading_label()
        indent = max(section.heading_level - 2, 0) * 2
        prefix = " " * indent
        anchor = generate_anchor_from_heading(label, include_hash=True)
        toc_lines.append(f"{prefix}- [{label}]({anchor})")
    return toc_lines


def render_prose_section(parsed_index, section: ProseSection) -> List[str]:
    lines: List[str] = []
    heading_label = format_prose_heading(section)
    if heading_label:
        lines.append(f"{'#' * section.heading_level} {heading_label}")
        lines.append("")
    if section.content:
        lines.extend(section.content.splitlines())
        lines.append("")
    for child in section.child_sections:
        lines.extend(render_prose_section(parsed_index, child))
    return lines


def render_section_markdown(parsed_index, section) -> List[str]:
    lines = [f"{'#' * section.heading_level} {section.heading_label()}"]
    lines.append("")
    for entry in section.expected_entries.values():
        current_entry = parsed_index.find_current_entry(entry.name)
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


def index_to_markdown(parsed_index) -> str:
    lines: List[str] = []
    title = parsed_index.title.strip() if parsed_index.title else ""
    if title:
        lines.append(f"# {title}")
        lines.append("")

    toc_lines = render_toc(parsed_index)
    if toc_lines:
        lines.extend(toc_lines)
        lines.append("")

    if parsed_index.overview:
        lines.extend(render_prose_section(parsed_index, parsed_index.overview))

    for section_path in parsed_index.section_order:
        section = parsed_index.sections.get(section_path)
        if section is None:
            continue
        lines.extend(render_section_markdown(parsed_index, section))

    return "\n".join(lines).rstrip() + "\n"


def render_full_tree(parsed_index) -> List[str]:
    lines: List[str] = []
    for section_path in parsed_index.section_order:
        section = parsed_index.sections.get(section_path)
        if section is None:
            continue
        lines.append(section_path)
        entries = dict[str, object](section.expected_entries)
        for name, entry in section.current_entries.items():
            if name in entries:
                continue
            if entry.entry_status in ("orphaned", "removed"):
                entries[name] = entry
        for entry in sorted(entries.values(), key=lambda item: item.sort_key()):
            marker = ""
            if entry.entry_status == "added":
                marker = " [ADDED]"
            elif entry.entry_status == "moved":
                marker = " [MOVED]"
            elif entry.entry_status == "reordered":
                marker = " [REORDERED]"
            elif entry.entry_status == "unresolved":
                marker = " [UNRESOLVED]"
            elif entry.entry_status == "orphaned":
                marker = " [ORPHANED]"
            elif entry.entry_status == "removed":
                marker = " [REMOVED]"
            lines.append(f"- {entry.name}{marker}")
        lines.append("")
    return lines
