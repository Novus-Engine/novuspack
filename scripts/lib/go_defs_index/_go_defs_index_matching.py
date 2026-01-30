"""
Phase 3 matching / placement logic for validate_api_go_defs_index.

This module extracts the matching and scoring logic (determine_section_placement and
calculate_confidence_score) plus supporting helpers, in order to keep
validate_api_go_defs_index.py focused on orchestration.
"""

from __future__ import annotations

from dataclasses import dataclass
from typing import Dict, List, Optional, Set, Tuple, TYPE_CHECKING

from lib._go_code_utils import normalize_generic_name
from lib.go_defs_index._go_defs_index_models import DetectedDefinition, IndexEntry
from lib._index_utils import IndexSection
from lib.go_defs_index._go_defs_index_scoring import calculate_confidence_score
from lib.go_defs_index._go_defs_index_shared import map_implementation_to_interface
from lib._validation_utils import OutputBuilder

if TYPE_CHECKING:
    from lib._index_utils import ParsedIndex

# Keep consistent with validate_api_go_defs_index.py
CONFIDENCE_THRESHOLD = 0.75


def _expected_kind_for_definition(defn: DetectedDefinition) -> str:
    """Return section kind that must match for placement (method/func/type)."""
    if defn.kind == "method":
        return "method"
    if defn.kind == "func":
        return "func"
    if defn.kind == "type":
        return "type"
    return "type"


def _definition_sort_key(defn: DetectedDefinition) -> Tuple[int, str]:
    kind_order = {"type": 0, "method": 1, "func": 2}
    return (kind_order.get(defn.kind, 99), defn.name.lower())


def _populate_section_valid_types(index_sections: List[IndexSection]) -> None:
    """
    Set valid_types on each IndexSection from index entries and structure.

    Type sections: valid types = entry names that are types (no ".").
    Method sections: valid types inherited from parent type section.
    Function sections: valid_types left empty.
    """
    for sec in index_sections:
        if sec.kind == "type":
            sec.valid_types = {
                normalize_generic_name(e.name).lower()
                for e in sec.entries
                if "." not in e.name
            }
        else:
            sec.valid_types = set()

    for sec in index_sections:
        if sec.kind != "method":
            continue
        if sec.parent_heading and sec.parent_heading.kind == "type":
            sec.valid_types = set(sec.parent_heading.valid_types)


_METHOD_CATEGORY_RULES = {
    "Package": [
        ("Package Signature Methods", ["signature", "sign", "validatesignature"]),
        ("Package Lifecycle Methods", ["open", "close", "isopen", "isreadonly"]),
        ("Package File Management Methods", ["addfile", "removefile", "extract"]),
        (
            "Package Information and Queries Methods",
            ["getinfo", "getmetadata", "listfiles", "fileexists", "getfile"],
        ),
        ("Package Metadata Methods", ["comment", "appid", "vendorid", "identity"]),
        ("Package Compression Methods", ["compress", "decompress"]),
        (
            "Package Path and Configuration Methods",
            ["targetpath", "extractroot", "sessionbase"],
        ),
    ],
    "FileEntry": [
        (
            "FileEntry Data Management Methods",
            ["loaddata", "unloaddata", "setdata", "getdata"],
        ),
        (
            "FileEntry Transformation Methods",
            ["compress", "decompress", "encrypt", "decrypt", "transform"],
        ),
        ("FileEntry Tag Management Methods", ["tag", "tags"]),
        (
            "FileEntry Path and Metadata Methods",
            ["path", "metadata", "associate"],
        ),
        ("FileEntry Source Management Methods", ["source", "current", "original"]),
        ("FileEntry Marshaling Methods", ["marshal", "unmarshal"]),
        (
            "FileEntry Query Methods",
            ["is", "has", "get", "compressed", "encrypted"],
        ),
    ],
    "PackageReader": [
        ("PackageReader Read Operations", ["read", "readfile"]),
        ("PackageReader Query Operations", ["list", "getinfo", "getmetadata", "query"]),
    ],
    "PackageWriter": [
        ("PackageWriter Write Operations", ["write", "writedata", "writefile"]),
    ],
}

_METHOD_CATEGORY_DEFAULTS = {
    "Package": "Package Other Methods",
    "FileEntry": "FileEntry Other Methods",
    "PackageReader": "PackageReader Other Methods",
    "PackageWriter": "PackageWriter Other Methods",
}


def _categorize_by_keywords(
    method_lower: str,
    categories: List[Tuple[str, List[str]]],
    fallback: str,
) -> str:
    for category, keywords in categories:
        if any(word in method_lower for word in keywords):
            return category
    return fallback


def categorize_method(method_name: str, receiver_type: str) -> str:
    """
    Categorize a method into a logical group for sub-subsection placement.

    Returns category name like "Package Lifecycle Methods",
    "FileEntry Data Management Methods", etc.
    The category name includes the type prefix.
    """
    method_lower = method_name.lower()

    if receiver_type == "PathMetadataEntry":
        # PathMetadataEntry methods are grouped under "Metadata Methods" in the index.
        return "Metadata Methods"

    categories = _METHOD_CATEGORY_RULES.get(receiver_type)
    if not categories:
        return "Other Methods"
    fallback = _METHOD_CATEGORY_DEFAULTS.get(receiver_type, "Other Methods")
    return _categorize_by_keywords(method_lower, categories, fallback)


def _section_candidates_by_kind(
    index_sections: Dict[str, IndexSection],
    section_order: List[str],
    kind: str,
) -> List[IndexSection]:
    candidates: List[IndexSection] = []
    for section_path in section_order:
        section = index_sections.get(section_path)
        if not section or section.kind != kind:
            continue
        candidates.append(section)
    return candidates


def _get_section_by_receiver(
    receiver_type: str,
    type_sections: List[IndexSection],
) -> Optional[IndexSection]:
    receiver = normalize_generic_name(receiver_type).lower()
    for section in type_sections:
        if receiver in section.expected_entries:
            return section
    return None


def _normalize_receiver_type(receiver_type: str) -> str:
    if not receiver_type:
        return receiver_type
    mapped = map_implementation_to_interface(receiver_type)
    return normalize_generic_name(mapped)


def _is_signature_package_method(defn: DetectedDefinition) -> bool:
    if defn.file == "api_signatures.md":
        return True
    name_lower = defn.name.lower()
    if "signature" in name_lower:
        return True
    method_name = name_lower.split(".", 1)[1] if "." in name_lower else name_lower
    return method_name.startswith("sign") or method_name.startswith("validate")


def _select_method_section_by_category(
    defn: DetectedDefinition,
    receiver_type: str,
    candidates: List[IndexSection],
) -> Optional[IndexSection]:
    method_name = defn.name.split(".", 1)[1] if "." in defn.name else defn.name
    category = categorize_method(method_name, receiver_type)
    for section in candidates:
        if section.heading_text == category:
            return section
    if receiver_type == "Package" and _is_signature_package_method(defn):
        for section in candidates:
            if section.heading_text == "Package Other Methods":
                return section
    return None


def _best_section_for_definition(
    definition: DetectedDefinition,
    candidate_sections: List[IndexSection],
    all_sections: Set[str],
    section_valid_types: Dict[str, Set[str]],
) -> Tuple[Optional[IndexSection], float, List[str]]:
    best_section = None
    best_score = 0.0
    best_reasoning: List[str] = []
    for section in candidate_sections:
        section_path = section.path_label()
        score, reasoning = calculate_confidence_score(
            definition, section_path, all_sections, section_valid_types
        )
        if score > best_score:
            best_score = score
            best_reasoning = reasoning
            best_section = section
    return best_section, best_score, best_reasoning


def _assign_definition_to_section(
    definition: DetectedDefinition,
    section: IndexSection,
) -> None:
    entry = definition.to_index_entry(section.path_label())
    entry.suggested_section = section.path_label()
    if definition.confidence_score:
        entry.confidence_score = definition.confidence_score
    if definition.confidence_reasoning:
        entry.confidence_reasoning = list(definition.confidence_reasoning)
    section.expected_entries[entry.name] = entry


@dataclass
class PlacementContext:
    parsed_index: "ParsedIndex"
    type_sections: List[IndexSection]
    method_sections: List[IndexSection]
    func_sections: List[IndexSection]
    all_sections: Set[str]
    section_valid_types: Dict[str, Set[str]]
    unsorted_types: Optional[IndexSection]
    unsorted_methods: Optional[IndexSection]
    unsorted_funcs: Optional[IndexSection]
    unresolved_types: Set[str]
    definitions_by_name: Dict[str, DetectedDefinition]


def _add_unresolved_entry(
    definition: DetectedDefinition,
    unsorted_section: Optional[IndexSection],
    suggested_section: Optional[str] = None,
    confidence_score: Optional[float] = None,
    confidence_reasoning: Optional[List[str]] = None,
) -> Optional[IndexEntry]:
    if not unsorted_section:
        return None
    entry = definition.to_index_entry(unsorted_section.path_label())
    entry.entry_status = "unresolved"
    if suggested_section:
        entry.suggested_section = suggested_section
    if confidence_score is not None:
        entry.confidence_score = confidence_score
    if confidence_reasoning:
        entry.confidence_reasoning = list(confidence_reasoning)
    unsorted_section.expected_entries[entry.name] = entry
    return entry


def _place_type_definitions(
    definitions: List[DetectedDefinition],
    context: PlacementContext,
) -> None:
    for defn in [d for d in definitions if d.kind == "type"]:
        best_section, best_score, best_reasoning = _best_section_for_definition(
            defn,
            context.type_sections,
            context.all_sections,
            context.section_valid_types,
        )
        defn.confidence_score = best_score
        defn.confidence_reasoning = best_reasoning
        if best_section and best_score >= CONFIDENCE_THRESHOLD:
            _assign_definition_to_section(defn, best_section)
            continue
        entry = _add_unresolved_entry(
            defn,
            context.unsorted_types,
            best_section.path_label() if best_section else None,
            best_score,
            best_reasoning,
        )
        if entry:
            context.unresolved_types.add(entry.name)


def _place_method_definition(
    defn: DetectedDefinition,
    context: PlacementContext,
) -> None:
    if not defn.receiver_type:
        _add_unresolved_entry(defn, context.unsorted_methods)
        return
    receiver = _normalize_receiver_type(defn.receiver_type)
    if receiver in context.unresolved_types:
        current_section = context.parsed_index.find_section_by_current_entry(defn.name)
        if current_section and current_section.kind == "method":
            entry = defn.to_index_entry(current_section.path_label())
            entry.suggested_section = current_section.path_label()
            current_section.expected_entries[entry.name] = entry
            return
        _add_unresolved_entry(defn, context.unsorted_methods)
        return
    receiver_section = _get_section_by_receiver(receiver, context.type_sections)
    if not receiver_section:
        _add_unresolved_entry(defn, context.unsorted_methods)
        return
    candidates = [
        child for child in receiver_section.children if child.kind == "method"
    ]
    if not candidates:
        _add_unresolved_entry(defn, context.unsorted_methods)
        return
    category_section = _select_method_section_by_category(defn, receiver, candidates)
    if category_section is not None:
        defn.confidence_score = 1.0
        defn.confidence_reasoning = ["Category match: structure-first placement"]
        _assign_definition_to_section(defn, category_section)
        return
    best_section, best_score, best_reasoning = _best_section_for_definition(
        defn,
        candidates,
        context.all_sections,
        context.section_valid_types,
    )
    defn.confidence_score = best_score
    defn.confidence_reasoning = best_reasoning
    if best_section and best_score >= CONFIDENCE_THRESHOLD:
        _assign_definition_to_section(defn, best_section)
        return
    _add_unresolved_entry(
        defn,
        context.unsorted_methods,
        best_section.path_label() if best_section else None,
        best_score,
        best_reasoning,
    )


def _place_method_definitions(
    definitions: List[DetectedDefinition],
    context: PlacementContext,
) -> None:
    for defn in [d for d in definitions if d.kind == "method"]:
        _place_method_definition(defn, context)


def _propagate_unresolved_types(
    context: PlacementContext,
) -> None:
    if not context.unsorted_types:
        return
    for section in context.method_sections:
        if not section.parent_heading or section.parent_heading.kind != "type":
            continue
        parent = section.parent_heading
        for entry in section.expected_entries.values():
            if "." not in entry.name:
                continue
            receiver = normalize_generic_name(entry.name.split(".", 1)[0])
            if receiver not in context.unresolved_types:
                continue
            if not context.parsed_index.find_current_entry(entry.name):
                continue
            if receiver in parent.current_entries or receiver in parent.expected_entries:
                continue
            type_entry = context.unsorted_types.expected_entries.pop(receiver, None)
            if not type_entry:
                defn = context.definitions_by_name.get(receiver)
                if not defn:
                    continue
                type_entry = defn.to_index_entry(parent.path_label())
            type_entry.current_section = parent.path_label()
            type_entry.suggested_section = parent.path_label()
            parent.expected_entries[receiver] = type_entry
            context.unresolved_types.discard(receiver)
            receiver_lower = normalize_generic_name(receiver).lower()
            parent.valid_types.add(receiver_lower)
            context.section_valid_types.setdefault(parent.path_label(), set()).add(
                receiver_lower
            )


def _find_related_section(
    defn: DetectedDefinition,
    type_sections: List[IndexSection],
) -> Optional[IndexSection]:
    for method_name in defn.referenced_methods:
        receiver = method_name.split(".", 1)[0]
        receiver = _normalize_receiver_type(receiver)
        related_section = _get_section_by_receiver(receiver, type_sections)
        if related_section:
            return related_section
    for type_name in defn.referenced_types:
        receiver = _normalize_receiver_type(type_name)
        related_section = _get_section_by_receiver(receiver, type_sections)
        if related_section:
            return related_section
    return None


def _place_function_definitions(
    definitions: List[DetectedDefinition],
    context: PlacementContext,
) -> None:
    for defn in [d for d in definitions if d.kind == "func"]:
        related_section = _find_related_section(defn, context.type_sections)
        candidates = context.func_sections
        if related_section:
            related_candidates = [
                child for child in related_section.children if child.kind == "func"
            ]
            if related_candidates:
                candidates = related_candidates
        best_section, best_score, best_reasoning = _best_section_for_definition(
            defn,
            candidates,
            context.all_sections,
            context.section_valid_types,
        )
        defn.confidence_score = best_score
        defn.confidence_reasoning = best_reasoning
        if best_section and best_score >= CONFIDENCE_THRESHOLD:
            _assign_definition_to_section(defn, best_section)
            continue
        _add_unresolved_entry(
            defn,
            context.unsorted_funcs,
            best_section.path_label() if best_section else None,
            best_score,
            best_reasoning,
        )


def _emit_verbose_placements(
    definitions: List[DetectedDefinition],
    parsed_index: ParsedIndex,
    output: OutputBuilder,
) -> None:
    for defn in definitions:
        section = parsed_index.find_section_by_expected_entry(defn.name)
        entry = section.expected_entries.get(defn.name) if section else None
        if entry:
            reasoning_str = ", ".join(entry.confidence_reasoning)
            target_section = entry.suggested_section or entry.current_section
            if entry.confidence_score is None:
                score_str = "N/A"
            else:
                score_str = f"{int(entry.confidence_score * 100)}%"
            output.add_verbose_line(
                f"{defn.name} -> {target_section}: {score_str} ({reasoning_str})"
            )
        else:
            output.add_verbose_line(
                f"{defn.name} -> (no section): 0% (no valid matches)"
            )


def determine_section_placement(
    definitions: List[DetectedDefinition],
    parsed_index: ParsedIndex,
    output: Optional[OutputBuilder] = None,
) -> List[DetectedDefinition]:
    """
    Phase 3: Determine which index section each definition belongs to.

    Populates expected_entries on IndexSection objects.
    """
    for section in parsed_index.sections.values():
        section.expected_entries = {}

    # Ensure processing order: types, then methods, then funcs.
    definitions.sort(key=_definition_sort_key)

    all_sections = set(parsed_index.section_order)
    _populate_section_valid_types(list(parsed_index.sections.values()))
    section_valid_types = {
        section.path_label(): section.valid_types for section in parsed_index.sections.values()
    }

    type_sections = _section_candidates_by_kind(
        parsed_index.sections, parsed_index.section_order, "type"
    )
    method_sections = _section_candidates_by_kind(
        parsed_index.sections, parsed_index.section_order, "method"
    )
    func_sections = _section_candidates_by_kind(
        parsed_index.sections, parsed_index.section_order, "func"
    )

    unsorted_types = parsed_index.sections.get(parsed_index.unsorted_paths[0])
    unsorted_methods = parsed_index.sections.get(parsed_index.unsorted_paths[1])
    unsorted_funcs = parsed_index.sections.get(parsed_index.unsorted_paths[2])

    definitions_by_name = {defn.name: defn for defn in definitions}
    unresolved_types: Set[str] = set()
    context = PlacementContext(
        parsed_index=parsed_index,
        type_sections=type_sections,
        method_sections=method_sections,
        func_sections=func_sections,
        all_sections=all_sections,
        section_valid_types=section_valid_types,
        unsorted_types=unsorted_types,
        unsorted_methods=unsorted_methods,
        unsorted_funcs=unsorted_funcs,
        unresolved_types=unresolved_types,
        definitions_by_name=definitions_by_name,
    )

    _place_type_definitions(definitions, context)
    _place_method_definitions(definitions, context)
    _propagate_unresolved_types(context)
    _place_function_definitions(definitions, context)
    if output and output.verbose:
        _emit_verbose_placements(definitions, parsed_index, output)

    return definitions


# End of module.
