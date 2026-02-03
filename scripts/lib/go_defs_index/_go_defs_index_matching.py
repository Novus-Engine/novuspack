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
from lib.go_defs_index._go_defs_index_matching_helpers import (
    adjust_related_section_for_function,
)
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
    if defn.kind in ("method", "func"):
        return (1 if defn.kind == "method" else 2, defn.name.lower())
    return (0, defn.name.lower())


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
        if sec.kind not in ("method", "func"):
            continue
        if sec.parent_heading and sec.parent_heading.kind == "type":
            sec.valid_types = set(sec.parent_heading.valid_types)


_METHOD_CATEGORY_RULES = {
    "Package": [
        ("Package Comment Methods", ["comment"]),
        ("Package Identity Methods", ["appid", "vendorid", "identity", "packageidentity"]),
        (
            "Package Special File Methods",
            [
                "indexfile",
                "manifestfile",
                "metadatafile",
                "signaturefile",
                "specialfile",
                "specialmetadata",
            ],
        ),
        (
            "Package Path Metadata Methods",
            [
                "pathmetadata",
                "directorymetadata",
                "filepathassociation",
                "pathconflicts",
                "pathstats",
                "pathtree",
                "pathfiles",
                "filesinpath",
                "destpath",
                "targetexists",
            ],
        ),
        ("Package Symlink Methods", ["symlink"]),
        ("Package Metadata-Only Methods", ["metadataonly"]),
        ("Package Info Methods", ["packageinfo"]),
        (
            "Package Metadata Validation Methods",
            [
                "validatemetadataonly",
                "validatepathmetadata",
                "validatespecial",
                "validatesymlink",
                "validatepathwithin",
            ],
        ),
        ("Package Metadata Internal Methods", ["load", "save", "updatefilepathassociations"]),
        (
            "Package File Encryption Methods",
            ["encryptfile", "decryptfile", "validatefileencryption", "fileencryptioninfo"],
        ),
        ("Package Write Methods", ["safewrite", "fastwrite", "write"]),
        ("Package Signature Management Methods", ["signature", "sign", "validatesignature"]),
        (
            "Package Lifecycle Methods",
            ["open", "close", "validate", "integrity", "defragment"],
        ),
        (
            "Package File Management Methods",
            ["addfile", "removefile", "extract", "updatefile", "addfilepath", "removefilepath"],
        ),
        ("Package Compression Methods", ["compress", "compressed", "compression", "decompress"]),
        (
            "Package Information and Queries Methods",
            [
                "getinfo",
                "getmetadata",
                "listfiles",
                "fileexists",
                "getfile",
                "getpath",
                "getpathstats",
                "getpathmetadata",
                "isopen",
                "isreadonly",
                "securitystatus",
                "multipath",
                "list",
                "find",
                "has",
            ],
        ),
        (
            "Package Path and Configuration Methods",
            ["targetpath", "extractroot", "sessionbase"],
        ),
    ],
    "FileEntry": [
        ("FileEntry Data Methods", ["getdata", "setdata", "loaddata", "unloaddata", "data"]),
        ("FileEntry Temp File Methods", ["tempfile", "temp"]),
        ("FileEntry Serialization Methods", ["marshal", "writedata", "writemeta", "writeto"]),
        ("FileEntry Path Methods", ["path", "symlink", "associate", "resolve"]),
        (
            "FileEntry Transformation Methods",
            [
                "compress",
                "decompress",
                "encrypt",
                "decrypt",
                "transform",
                "process",
                "pipeline",
                "set",
                "unset",
                "current",
                "original",
                "processingstate",
                "validate",
                "cleanup",
                "resume",
                "execute",
                "copy",
            ],
        ),
        ("FileEntry Query Methods", ["get", "has", "is"]),
    ],
    "Tag": [
        ("Tag Methods", ["get", "set"]),
    ],
}

_METHOD_CATEGORY_DEFAULTS = {
    "Package": "Package Other Methods",
    "FileEntry": "FileEntry Query Methods",
    "Tag": "Tag Methods",
}

_PACKAGE_METHOD_OVERRIDE_EXACT = {
    "addpathtoexistingentry": "Package File Management Methods",
    "isopen": "Package Information and Queries Methods",
    "loadpathmetadata": "Package Metadata Internal Methods",
    "loadspecialmetadatafiles": "Package Metadata Internal Methods",
    "loadsymlinkmetadatafile": "Package Symlink Methods",
    "savepathmetadatafile": "Package Metadata Internal Methods",
    "savesymlinkmetadatafile": "Package Metadata Internal Methods",
    "updatefilemetadata": "Package Path Metadata Methods",
    "validatemetadataonlyintegrity": "Package Metadata Validation Methods",
    "validatemetadataonlypackage": "Package Metadata Validation Methods",
    "validatepathmetadata": "Package Metadata Validation Methods",
    "validatespecialfiles": "Package Metadata Validation Methods",
    "validatesymlinkpaths": "Package Metadata Validation Methods",
}

_PACKAGE_METHOD_OVERRIDE_CONTAINS = [
    ("validatefileencryption", "Package File Encryption Methods"),
    ("encryptioninfo", "Package File Encryption Methods"),
    ("compressioninfo", "Package Compression Methods"),
    ("listcompressedfiles", "Package Compression Methods"),
    ("bytag", "Package Information and Queries Methods"),
    ("metadataindex", "Package Compression Methods"),
    ("multipath", "Package Information and Queries Methods"),
    ("filepathassociations", "Package Metadata Internal Methods"),
    ("sessionbase", "Package Path and Configuration Methods"),
    ("targetpath", "Package Path and Configuration Methods"),
]

_PACKAGE_METHOD_OVERRIDE_PREFIXES = [
    ("updatefile", "Package File Management Methods"),
]


def _categorize_by_keywords(
    method_lower: str,
    categories: List[Tuple[str, List[str]]],
    fallback: str,
) -> str:
    for category, keywords in categories:
        if any(word in method_lower for word in keywords):
            return category
    return fallback


def _package_category_overrides(
    method_lower: str,
) -> Optional[str]:
    exact_override = _PACKAGE_METHOD_OVERRIDE_EXACT.get(method_lower)
    if exact_override:
        return exact_override
    for token, category in _PACKAGE_METHOD_OVERRIDE_CONTAINS:
        if token in method_lower:
            return category
    for prefix, category in _PACKAGE_METHOD_OVERRIDE_PREFIXES:
        if method_lower.startswith(prefix) and "pattern" not in method_lower:
            return category
    return None


def _package_category_from_file(
    defn: DetectedDefinition,
    method_lower: str,
) -> Optional[str]:
    if not defn.file:
        return None
    file_name = defn.file
    if file_name in (
        "api_file_mgmt_addition.md",
        "api_file_mgmt_removal.md",
        "api_file_mgmt_extraction.md",
    ):
        return "Package File Management Methods"
    if file_name == "api_file_mgmt_queries.md":
        if "compressed" in method_lower:
            return "Package Compression Methods"
        return "Package Information and Queries Methods"
    if file_name == "api_deduplication.md":
        return "Package Information and Queries Methods"
    if file_name == "api_signatures.md":
        return "Package Signature Management Methods"
    if file_name == "api_package_compression.md":
        return "Package Compression Methods"
    if file_name == "api_security.md" and "signature" in method_lower:
        return "Package Signature Management Methods"
    return None


def categorize_method(defn: DetectedDefinition, receiver_type: str) -> str:
    """
    Categorize a method into a logical group for sub-subsection placement.

    Returns category name like "Package Lifecycle Methods",
    "FileEntry Data Management Methods", etc.
    The category name includes the type prefix.
    """
    method_name = defn.name.split(".", 1)[1] if "." in defn.name else defn.name
    method_lower = method_name.lower()

    if receiver_type == "PathMetadataEntry":
        # PathMetadataEntry methods are grouped under Package Path Metadata Methods.
        return "Package Path Metadata Methods"

    if receiver_type == "Tag":
        return "Tag Methods"

    if receiver_type == "Package":
        override_category = _package_category_overrides(method_lower)
        if override_category:
            return override_category
        file_category = _package_category_from_file(defn, method_lower)
        if file_category:
            return file_category

    if receiver_type == "FileEntry":
        if method_lower in {"getdata", "setdata", "loaddata", "unloaddata"}:
            return "FileEntry Data Methods"
        if "tempfile" in method_lower:
            return "FileEntry Temp File Methods"
        if method_lower.startswith(("marshal", "writedata", "writemeta", "writeto")):
            return "FileEntry Serialization Methods"
        if any(
            token in method_lower
            for token in ("pathmetadata", "symlink", "path", "associate", "resolve")
        ):
            return "FileEntry Path Methods"
        if method_lower.startswith(("get", "has", "is")):
            return "FileEntry Query Methods"
        if method_lower.startswith(
            (
                "set",
                "compress",
                "decompress",
                "encrypt",
                "decrypt",
                "process",
                "transform",
                "resume",
                "execute",
                "cleanup",
                "validate",
                "copy",
                "unset",
            )
        ) or any(
            token in method_lower
            for token in ("pipeline", "current", "original", "processingstate")
        ):
            return "FileEntry Transformation Methods"

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
    receiver = normalize_generic_name(receiver_type)
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
    category = categorize_method(defn, receiver_type)
    for section in candidates:
        if section.heading_text == category:
            return section
    if receiver_type not in _METHOD_CATEGORY_RULES:
        for section in candidates:
            if section.heading_text.endswith("Other Type Methods"):
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
    for defn in [d for d in definitions if d.kind not in ("method", "func")]:
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
    receiver_section = _get_section_by_receiver(receiver, context.type_sections)
    if not receiver_section:
        if receiver in context.unresolved_types:
            current_section = context.parsed_index.find_section_by_current_entry(defn.name)
            if current_section and current_section.kind == "method":
                score, reasoning = calculate_confidence_score(
                    defn,
                    current_section.path_label(),
                    context.all_sections,
                    context.section_valid_types,
                )
                defn.confidence_score = score
                defn.confidence_reasoning = reasoning
                entry = defn.to_index_entry(current_section.path_label())
                entry.suggested_section = current_section.path_label()
                current_section.expected_entries[entry.name] = entry
                return
        _add_unresolved_entry(defn, context.unsorted_methods)
        return
    candidates: List[IndexSection] = []
    for child in receiver_section.children:
        if child.kind == "method":
            candidates.append(child)
        for grandchild in child.children:
            if grandchild.kind == "method":
                candidates.append(grandchild)
    if not candidates:
        _add_unresolved_entry(defn, context.unsorted_methods)
        return
    best_section, best_score, best_reasoning = _best_section_for_definition(
        defn,
        candidates,
        context.all_sections,
        context.section_valid_types,
    )
    category_section = _select_method_section_by_category(defn, receiver, candidates)
    if category_section is not None:
        defn.confidence_score = best_score
        defn.confidence_reasoning = list(best_reasoning)
        defn.confidence_reasoning.append(
            "Category match: structure-first placement (placement override)"
        )
        _assign_definition_to_section(defn, category_section)
        return
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
        related_section = adjust_related_section_for_function(defn, related_section)
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
        if not best_section:
            current_section = context.parsed_index.find_section_by_current_entry(defn.name)
            if current_section and current_section.kind == "func":
                best_score, best_reasoning = calculate_confidence_score(
                    defn,
                    current_section.path_label(),
                    context.all_sections,
                    context.section_valid_types,
                )
                best_section = current_section
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
        if not entry:
            for unsorted_path in parsed_index.unsorted_paths:
                unsorted_section = parsed_index.sections.get(unsorted_path)
                if not unsorted_section:
                    continue
                if defn.name in unsorted_section.expected_entries:
                    section = unsorted_section
                    entry = unsorted_section.expected_entries.get(defn.name)
                    break
        if entry:
            reasoning_str = ", ".join(entry.confidence_reasoning)
            target_section = (
                entry.suggested_section
                or entry.current_section
                or (section.path_label() if section else None)
            )
            if entry.confidence_score is None:
                score_str = "N/A"
            else:
                score_str = f"{int(entry.confidence_score * 100)}%"
            output.add_verbose_line(
                f"{defn.name} -> {target_section}: {score_str} ({reasoning_str})"
            )
            if (
                entry.confidence_score is not None
                and entry.confidence_score < CONFIDENCE_THRESHOLD
            ):
                output.add_warning_line(
                    (
                        f"Low-confidence placement: {defn.name} -> {target_section} "
                        f"({score_str})"
                    ),
                    verbose_only=True,
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

    for defn in definitions:
        current_section = parsed_index.find_section_by_current_entry(defn.name)
        defn.current_section = current_section.path_label() if current_section else None

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
