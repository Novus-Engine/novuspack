"""
Base scoring rules for Go defs index placement.
"""

from __future__ import annotations

import re
from dataclasses import dataclass
from typing import Dict, List, Optional, Set, Tuple

from lib._go_code_utils import normalize_generic_name
from lib.go_defs_index._go_defs_index_models import DetectedDefinition
from lib.go_defs_index._go_defs_index_scoring_domain import (
    extract_implementation_interface,
)
from lib.go_defs_index._go_defs_index_scoring_text import get_section_kind
from lib.go_defs_index._go_defs_index_shared import map_implementation_to_interface


@dataclass(frozen=True)
class ScoringContext:
    definition: DetectedDefinition
    section: str
    all_sections: Set[str]
    section_valid_types: Optional[Dict[str, Set[str]]]
    section_lower: str
    name_lower: str
    heading_lower: str
    content_lower: str
    detected_domain: Optional[str]


def _infer_function_types(defn: DetectedDefinition) -> Set[str]:
    if defn.kind != "func":
        return set()
    name = defn.name
    name_lower = name.lower()
    out: Set[str] = set()
    if name_lower.startswith("new") and len(name) > 3:
        candidate = name[3:]
        if candidate and candidate[0].isupper():
            out.add(normalize_generic_name(candidate).lower())
    for prefix in ("get", "set", "unmarshal", "marshal", "add", "remove", "has", "is"):
        if name_lower.startswith(prefix) and len(name) > len(prefix):
            rest = name[len(prefix):]
            if rest and rest[0].isupper():
                base = re.split(r"[^A-Za-z0-9]", rest, 1)[0]
                if base:
                    out.add(normalize_generic_name(base).lower())
    if name and name[0].isupper():
        base = re.split(r"[^A-Za-z0-9]", name, 1)[0]
        if base:
            out.add(normalize_generic_name(base).lower())
    return out


def _extract_primary_name_from_section(section_name: str) -> str:
    leaf = section_name.split(">")[-1].strip()
    leaf = re.sub(r"^\\d+(?:\\.\\d+)*\\.?\\s+", "", leaf)
    leaf_lower = leaf.lower()
    suffixes = [
        "interface",
        "type",
        "types",
        "structure",
        "struct",
        "definition",
        "definitions",
    ]
    for suffix in suffixes:
        if leaf_lower.endswith(" " + suffix):
            leaf_lower = leaf_lower[: -(len(suffix) + 1)].strip()
            break
    leaf_lower = re.sub(r"[^a-z0-9]", "", leaf_lower)
    return leaf_lower


def _is_core_package_type(
    name: str,
    kind: str,
    receiver_type: Optional[str] = None,
) -> bool:
    name_lower_check = name.lower()
    core_package_types = ["package", "packagereader", "packagewriter", "filepackage"]
    if name_lower_check in core_package_types:
        return True
    if kind == "method" and receiver_type:
        receiver_lower = receiver_type.lower()
        receiver_base = receiver_lower.split("[")[0] if "[" in receiver_lower else receiver_lower
        if receiver_base in core_package_types:
            return True
    return False


def _error_section_flags(
    section_lower: str,
    name_lower: str,
    definition_kind: str,
) -> Tuple[bool, bool, bool]:
    is_error_helper_functions = (
        "error" in section_lower and "helper" in section_lower and "function" in section_lower
    )
    is_error_section = (
        (
            "error types" in section_lower
            or ("error" in section_lower and ("type" in section_lower or "errors" in section_lower))
        )
        and not is_error_helper_functions
    )
    is_error_definition = (
        name_lower.startswith("err")
        or "error" in name_lower
        or definition_kind in ("type", "struct")
        and ("error" in name_lower or "err" in name_lower)
    )
    return is_error_section, is_error_definition, is_error_helper_functions


def score_strict_kind_matching(ctx: ScoringContext) -> Tuple[float, List[str], bool]:
    section_kind = get_section_kind(ctx.section)
    if section_kind:
        expected_kind: Optional[str] = None
        if ctx.definition.kind == "method":
            expected_kind = "method"
        elif ctx.definition.kind == "func":
            expected_kind = "func"
        elif ctx.definition.kind in ("type", "struct", "interface", "alias"):
            expected_kind = "type"

        if expected_kind and section_kind != expected_kind:
            return (
                -1.0,
                [
                    "STRICT kind mismatch (blocked)",
                    f"definition.kind={ctx.definition.kind}, section_kind={section_kind}",
                ],
                True,
            )

    if (
        ctx.definition.kind == "method"
        and ctx.definition.receiver_type
        and ctx.section_valid_types
        and ctx.section in ctx.section_valid_types
        and ctx.section_valid_types[ctx.section]
    ):
        mapped_receiver = map_implementation_to_interface(ctx.definition.receiver_type)
        receiver_normalized = normalize_generic_name(mapped_receiver).lower()
        if receiver_normalized not in ctx.section_valid_types[ctx.section]:
            return (
                -1.0,
                [
                    "Receiver type not allowed by section structure (blocked)",
                    f"receiver={mapped_receiver}, section={ctx.section}",
                ],
                True,
            )
        return (
            0.50,
            [f"Receiver type match (index structure): +50% ({mapped_receiver})"],
            False,
        )

    return 0.0, [], False


def score_function_type_interaction(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if (
        ctx.definition.kind == "func"
        and ctx.section_valid_types
        and ctx.section in ctx.section_valid_types
        and ctx.section_valid_types[ctx.section]
    ):
        func_types = _infer_function_types(ctx.definition)
        section_types = ctx.section_valid_types[ctx.section]
        if func_types and section_types & func_types:
            return 0.60, ["Function interacts with type in section: +60%"]
    return 0.0, []


def score_exact_name_match(ctx: ScoringContext) -> Tuple[float, List[str]]:
    score = 0.0
    reasoning: List[str] = []
    section_leaf_lower = ctx.section.split(">")[-1].strip().lower()
    section_primary = _extract_primary_name_from_section(ctx.section)

    if ctx.definition.kind == "interface" and "interface" in section_leaf_lower:
        mapped_name = map_implementation_to_interface(ctx.definition.name)
        mapped_norm = re.sub(
            r"[^a-z0-9]",
            "",
            normalize_generic_name(mapped_name).lower(),
        )
        if mapped_norm and section_primary and mapped_norm == section_primary:
            score += 0.60
            reasoning.append(f"Exact interface name match ({mapped_name}): +60%")

    if ctx.definition.kind in ("type", "struct") and (
        " type" in section_leaf_lower
        or " types" in section_leaf_lower
        or " structure" in section_leaf_lower
        or " struct" in section_leaf_lower
    ):
        mapped_name = map_implementation_to_interface(ctx.definition.name)
        mapped_norm = re.sub(
            r"[^a-z0-9]",
            "",
            normalize_generic_name(mapped_name).lower(),
        )
        if mapped_norm and section_primary and mapped_norm == section_primary:
            score += 0.60
            reasoning.append(f"Exact type name match ({mapped_name}): +60%")

    return score, reasoning


def score_implementation_mapping(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if ctx.definition.kind not in ("type", "struct"):
        return 0.0, []
    score = 0.0
    reasoning: List[str] = []
    extracted_interface = extract_implementation_interface(ctx.definition)
    if extracted_interface:
        extracted_lower = extracted_interface.lower()
        interface_pattern = r"\\b" + re.escape(extracted_lower) + r"\\s+interface"
        if re.search(interface_pattern, ctx.section_lower):
            score += 0.40
            msg = (
                f"Implementation reference: {ctx.definition.name} implements "
                f"{extracted_interface} (from section content): +40%"
            )
            reasoning.append(msg)

    mapped_type = map_implementation_to_interface(ctx.definition.name)
    if mapped_type != ctx.definition.name:
        mapped_lower = mapped_type.lower()
        interface_pattern = r"\\b" + re.escape(mapped_lower) + r"\\s+interface"
        if re.search(interface_pattern, ctx.section_lower):
            if not extracted_interface or extracted_interface != mapped_type:
                score += 0.30
                msg = (
                    f"Implementation type ({ctx.definition.name}) maps to interface "
                    f"({mapped_type}): +30%"
                )
                reasoning.append(msg)

    return score, reasoning


def _score_error_section_penalties(ctx: ScoringContext) -> Tuple[float, List[str], bool]:
    score = 0.0
    reasoning: List[str] = []
    kind_mismatch = False
    is_error_section, is_error_definition, _ = _error_section_flags(
        ctx.section_lower,
        ctx.name_lower,
        ctx.definition.kind,
    )
    if not is_error_section:
        return 0.0, [], False
    if ctx.definition.kind in ("method", "func"):
        score -= 0.5
        reasoning.append(f"Non-error {ctx.definition.kind} in Error Types section: -50%")
        kind_mismatch = True
    elif ctx.definition.kind in ("type", "struct") and not is_error_definition:
        score -= 0.5
        reasoning.append(f"Non-error {ctx.definition.kind} in Error Types section: -50%")
        kind_mismatch = True
    elif ctx.definition.kind == "interface" and not is_error_definition:
        score -= 0.5
        reasoning.append("Non-error interface in Error Types section: -50%")
        kind_mismatch = True
    return score, reasoning, kind_mismatch


def _score_method_blockers(ctx: ScoringContext) -> Tuple[float, List[str], bool]:
    if "helper" in ctx.section_lower and "function" in ctx.section_lower:
        return (
            -999.0,
            ["CRITICAL: method in Helper Functions section (must be function): BLOCKED"],
            True,
        )
    if ("type" in ctx.section_lower or "types" in ctx.section_lower) and "method" not in (
        ctx.section_lower
    ):
        return (
            -999.0,
            ["CRITICAL: method in Type section (must be type): BLOCKED"],
            True,
        )
    return 0.0, [], False


def _score_func_blockers(ctx: ScoringContext) -> Tuple[float, List[str], bool]:
    if ("type" in ctx.section_lower or "types" in ctx.section_lower) and "function" not in (
        ctx.section_lower
    ):
        return (
            -999.0,
            ["CRITICAL: function in Type section (must be type): BLOCKED"],
            True,
        )
    return 0.0, [], False


def _score_type_blockers(ctx: ScoringContext) -> Tuple[float, List[str], bool]:
    if "helper" in ctx.section_lower and "function" in ctx.section_lower:
        msg = (
            f"CRITICAL: {ctx.definition.kind} in Helper Functions section "
            f"(must be function): BLOCKED"
        )
        return -999.0, [msg], True
    if "methods" in ctx.section_lower and "type" not in ctx.section_lower:
        msg = (
            f"CRITICAL: {ctx.definition.kind} in Methods section "
            f"(must be method): BLOCKED"
        )
        return -999.0, [msg], True
    if ctx.definition.kind == "interface":
        if (
            ("type definition" in ctx.section_lower or "type definitions" in ctx.section_lower)
            and "interface" not in ctx.section_lower
        ):
            return -0.3, ["Kind mismatch: interface in Type section: -30%"], True
        if "method" in ctx.section_lower and "interface" not in ctx.section_lower:
            return -0.3, ["Kind mismatch: interface in Method section: -30%"], True
    return 0.0, [], False


def _score_kind_blockers_by_kind(ctx: ScoringContext) -> Tuple[float, List[str], bool]:
    if ctx.definition.kind == "method":
        return _score_method_blockers(ctx)
    if ctx.definition.kind == "func":
        return _score_func_blockers(ctx)
    if ctx.definition.kind in ("type", "struct", "interface"):
        return _score_type_blockers(ctx)
    return 0.0, [], False


def score_kind_blockers(ctx: ScoringContext) -> Tuple[float, List[str], bool]:
    score = 0.0
    reasoning: List[str] = []
    kind_mismatch = False
    delta, delta_reasoning, error_mismatch = _score_error_section_penalties(ctx)
    score += delta
    reasoning.extend(delta_reasoning)
    kind_mismatch = kind_mismatch or error_mismatch
    delta, delta_reasoning, blocker_mismatch = _score_kind_blockers_by_kind(ctx)
    score += delta
    reasoning.extend(delta_reasoning)
    kind_mismatch = kind_mismatch or blocker_mismatch
    return score, reasoning, kind_mismatch


def score_kind_positive_match(ctx: ScoringContext, kind_mismatch: bool) -> Tuple[float, List[str]]:
    if kind_mismatch:
        return 0.0, []
    if ctx.definition.kind == "method" and "methods" in ctx.section_lower:
        return 0.20, ["Kind 'method' matches Methods section: +20%"]
    if (
        ctx.definition.kind == "func"
        and "helper" in ctx.section_lower
        and "function" in ctx.section_lower
    ):
        return 0.20, ["Kind 'func' matches Helper Functions section: +20%"]
    if (
        ctx.definition.kind in ("type", "struct", "interface")
        and ("type" in ctx.section_lower or "types" in ctx.section_lower)
    ):
        return 0.20, [f"Kind '{ctx.definition.kind}' matches Type section: +20%"]
    return 0.0, []


def score_current_section(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if ctx.definition.current_section and ctx.definition.current_section == ctx.section:
        return 0.10, ["Current section match: +10%"]
    return 0.0, []
