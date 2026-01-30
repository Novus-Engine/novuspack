"""
Parse validation output from validate_api_go_defs_index for the apply script.

Provides dataclasses and parse_validation_output() to extract missing definitions,
move entries, link updates, description updates, and orphaned entry names.
"""

from __future__ import annotations

import re
from dataclasses import dataclass
from typing import Dict, List, Optional, Tuple

from lib._go_code_utils import normalize_generic_name

MISSING_DEFS_HEADER = re.compile(
    r"^Found\s+\d+\s+high-confidence\s+sorted\s+definition\(s\)\s+not in index:\s*$",
    re.MULTILINE,
)
MISSING_DESC_RE = re.compile(
    r"^ERROR: (Missing description|Description too short): .*?:\d+: "
    r"Entry `([^`]+)`.* -> (.+)$"
)
WRONG_SECTION_RE = re.compile(
    r"^ERROR: Wrong section: .*?:\d+: `([^`]+)` in '([^']+)' -> "
    r"Move to '([^']+)'"
)
LINK_UPDATE_RE = re.compile(
    r"^ERROR: Incorrect link: .*?:\d+: `([^`]+)`: ([^ ]+) -> Update to: (\S+)$"
)
ORPHANED_ENTRY_RE = re.compile(
    r"^ERROR: Orphaned entry: .*?:\d+: `([^`]+)` not found"
)


@dataclass(frozen=True)
class MissingDefinition:
    """A missing definition suggested by the validator."""

    name: str
    kind: str
    file: str
    line: int
    suggested_section: str
    canonical_file: str
    canonical_anchor: Optional[str]


@dataclass(frozen=True)
class MoveEntry:
    """Move an entry from one section to another."""

    name: str
    from_section: str
    to_section: str


@dataclass(frozen=True)
class LinkUpdate:
    """Update link file and optional anchor for an entry."""

    name: str
    link_file: str
    link_anchor: Optional[str]


@dataclass(frozen=True)
class DescriptionUpdate:
    """Add or replace description for an entry."""

    name: str
    description: str


def _split_link_target(target: str) -> Tuple[str, Optional[str]]:
    if "#" in target:
        file_name, anchor = target.split("#", 1)
        return file_name, anchor
    return target, None


def _extract_missing_defs_block(content: str) -> Optional[str]:
    preferred = MISSING_DEFS_HEADER.search(content)
    if not preferred:
        return None
    tail = content[preferred.end():]
    end_match = re.search(r"(?m)^(Found\s+\d+|ERROR:|=+)", tail)
    end_idx = end_match.start() if end_match else len(tail)
    return tail[:end_idx]


def _build_missing_def(fields: Dict[str, str]) -> Optional[MissingDefinition]:
    name = fields.get("name")
    kind = fields.get("kind")
    file_line = fields.get("file_line", "")
    section = fields.get("section", "")
    canonical = fields.get("canonical", "")
    if not name or not kind or not section or not canonical:
        return None
    if " (confidence:" in section:
        section = section.split(" (confidence:", 1)[0].strip()
    file_part = ""
    line_num = 0
    if file_line:
        file_part, _, line_part = file_line.partition(":")
        try:
            line_num = int(line_part)
        except ValueError:
            return None
    if canonical.strip() == "(add to parent type section)":
        canonical_file = ""
        canonical_anchor = None
    else:
        canonical_file, canonical_anchor = _split_link_target(canonical)
    return MissingDefinition(
        name=normalize_generic_name(name),
        kind=kind,
        file=file_part,
        line=line_num,
        suggested_section=section,
        canonical_file=canonical_file,
        canonical_anchor=canonical_anchor,
    )


def _parse_missing_definitions(content: str) -> List[MissingDefinition]:
    defs_section = _extract_missing_defs_block(content)
    if not defs_section:
        return []

    defs: List[MissingDefinition] = []
    current: Dict[str, str] = {}

    for raw_line in defs_section.splitlines():
        line = raw_line.rstrip()
        if not line.strip():
            continue
        if line.startswith("  ") and not line.lstrip().startswith("-"):
            if current:
                built = _build_missing_def(current)
                if built:
                    defs.append(built)
                current = {}
            current["name"] = line.strip()
            continue
        if line.strip().startswith("- Kind:"):
            current["kind"] = line.split(":", 1)[1].strip()
        elif line.strip().startswith("- File:"):
            current["file_line"] = line.split(":", 1)[1].strip()
        elif line.strip().startswith("- Suggested section:"):
            current["section"] = line.split(":", 1)[1].strip()
        elif line.strip().startswith("- Canonical location:"):
            current["canonical"] = line.split(":", 1)[1].strip()

    if current:
        built = _build_missing_def(current)
        if built:
            defs.append(built)

    return defs


def _parse_move_entries(content: str) -> List[MoveEntry]:
    moves: List[MoveEntry] = []
    for line in content.splitlines():
        match = WRONG_SECTION_RE.match(line.strip())
        if not match:
            continue
        name, from_section, to_section = match.groups()
        moves.append(
            MoveEntry(
                name=normalize_generic_name(name),
                from_section=from_section,
                to_section=to_section,
            )
        )
    return moves


def _parse_link_updates(content: str) -> List[LinkUpdate]:
    updates: List[LinkUpdate] = []
    for line in content.splitlines():
        match = LINK_UPDATE_RE.match(line.strip())
        if not match:
            continue
        name, _current, suggested = match.groups()
        link_file, link_anchor = _split_link_target(suggested)
        updates.append(
            LinkUpdate(
                name=normalize_generic_name(name),
                link_file=link_file,
                link_anchor=link_anchor,
            )
        )
    return updates


def _parse_description_updates(content: str) -> List[DescriptionUpdate]:
    updates: List[DescriptionUpdate] = []
    for line in content.splitlines():
        match = MISSING_DESC_RE.match(line.strip())
        if not match:
            continue
        _issue_type, name, suggestion = match.groups()
        marker = "Review the definition comments for potential summary:"
        if marker not in suggestion:
            continue
        description = suggestion.split(marker, 1)[1].strip()
        if len(description) < 20:
            continue
        updates.append(
            DescriptionUpdate(
                name=normalize_generic_name(name),
                description=description,
            )
        )
    return updates


def _parse_orphaned_entries(content: str) -> List[str]:
    names: List[str] = []
    for line in content.splitlines():
        match = ORPHANED_ENTRY_RE.match(line.strip())
        if not match:
            continue
        name = match.group(1)
        names.append(normalize_generic_name(name))
    return names


def parse_validation_output(content: str) -> Tuple[
    List[MissingDefinition],
    List[MoveEntry],
    List[LinkUpdate],
    List[DescriptionUpdate],
    List[str],
]:
    """Parse validator output into missing defs, moves, link/description updates, orphaned names."""
    return (
        _parse_missing_definitions(content),
        _parse_move_entries(content),
        _parse_link_updates(content),
        _parse_description_updates(content),
        _parse_orphaned_entries(content),
    )
