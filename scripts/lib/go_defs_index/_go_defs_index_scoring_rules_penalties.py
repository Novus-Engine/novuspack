"""
Penalty and special-case scoring rules.
"""

from __future__ import annotations

from typing import List, Optional, Tuple

from lib.go_defs_index._go_defs_index_scoring_domain import extract_error_domain
from lib.go_defs_index._go_defs_index_scoring_rules_core import (
    ScoringContext,
    _error_section_flags,
)


def score_hash_optional_types(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if ctx.definition.kind not in ("type", "struct"):
        return 0.0, []
    if "hash and optional data" not in ctx.section_lower:
        return 0.0, []
    hash_keywords = [
        "hashpurpose",
        "hashtype",
        "optionaldata",
        "processingstate",
        "tagvaluetype",
        "transformtype",
    ]
    if any(kw in ctx.name_lower for kw in hash_keywords):
        msg = (
            "Hash/Optional data type matches 'Hash and Optional Data Types' "
            "subsection: +30%"
        )
        return 0.30, [msg]
    return 0.0, []


def score_error_domain_match(ctx: ScoringContext) -> Tuple[float, List[str]]:
    is_error_section, is_error_definition, _ = _error_section_flags(
        ctx.section_lower,
        ctx.name_lower,
        ctx.definition.kind,
    )
    if not (is_error_section and is_error_definition):
        return 0.0, []
    error_domain = extract_error_domain(ctx.definition.name)
    if not error_domain:
        return 0.0, []
    if error_domain in ctx.section_lower:
        return 0.15, [f"Error domain '{error_domain}' matches subsection domain: +15%"]
    return 0.0, []


def _infer_penalty_domain(ctx: ScoringContext) -> Optional[str]:
    if "compression" in ctx.name_lower or "compress" in ctx.name_lower:
        return "compression"
    if "encryption" in ctx.name_lower or "encrypt" in ctx.name_lower:
        return "encryption"
    if "signature" in ctx.name_lower or "sign" in ctx.name_lower:
        return "signature"
    if "security" in ctx.name_lower:
        return "security"
    if "file" in ctx.name_lower and "handler" in ctx.name_lower:
        encrypt_keywords = ["encrypt", "aes", "chacha", "mlkem"]
        if any(kw in ctx.name_lower for kw in encrypt_keywords):
            return "encryption"
    return None


def _find_matching_type_section(ctx: ScoringContext, definition_domain: str) -> Optional[str]:
    for section in sorted(ctx.all_sections):
        section_lower = section.lower()
        if "type definition" in section_lower or "type definitions" in section_lower:
            if definition_domain in section_lower:
                return section
            if definition_domain == "hash" and (
                "hash" in section_lower and "optional" in section_lower
            ):
                return section
    return None


def score_type_operation_penalty(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if ctx.definition.kind not in ("type", "interface"):
        return 0.0, []
    is_operation_section = (
        "type definition" not in ctx.section_lower
        and "type definitions" not in ctx.section_lower
        and ctx.section_lower not in ["core interfaces", "generics"]
    )
    if not is_operation_section:
        return 0.0, []
    definition_domain = _infer_penalty_domain(ctx)
    if not definition_domain:
        return 0.0, []
    matching_type_section = _find_matching_type_section(ctx, definition_domain)
    if matching_type_section:
        msg = (
            "Type should be in Type Definitions section "
            f"'{matching_type_section}': -40%"
        )
        return -0.4, [msg]
    return 0.0, []


def _penalty_error_section(ctx: ScoringContext) -> Tuple[float, List[str], bool]:
    score = 0.0
    reasoning: List[str] = []
    kind_mismatch = False
    is_error_section, is_error_definition, is_error_helper = _error_section_flags(
        ctx.section_lower,
        ctx.name_lower,
        ctx.definition.kind,
    )
    if not (is_error_section and not is_error_helper):
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


def _penalty_method_in_type_section(ctx: ScoringContext) -> Tuple[float, List[str], bool]:
    if ctx.definition.kind != "method":
        return 0.0, [], False
    if "type definition" in ctx.section_lower or "type definitions" in ctx.section_lower:
        return -0.3, ["Kind mismatch: method in Type section: -30%"], True
    return 0.0, [], False


def _penalty_func_in_type_section(ctx: ScoringContext) -> Tuple[float, List[str], bool]:
    if ctx.definition.kind != "func":
        return 0.0, [], False
    if "type definition" in ctx.section_lower or "type definitions" in ctx.section_lower:
        return -0.5, ["Kind mismatch: function in Type section: -50%"], True
    return 0.0, [], False


def _penalty_type_in_method_section(ctx: ScoringContext) -> Tuple[float, List[str], bool]:
    if ctx.definition.kind not in ("type", "interface"):
        return 0.0, [], False
    if (
        ctx.definition.kind in ("type", "struct")
        and "method" in ctx.section_lower
        and "type" not in ctx.section_lower
    ):
        msg = f"Kind mismatch: {ctx.definition.kind} in Method section: -30%"
        return -0.3, [msg], True
    if ctx.definition.kind == "interface":
        if "method" in ctx.section_lower and "interface" not in ctx.section_lower:
            return -0.3, ["Kind mismatch: interface in Method section: -30%"], True
    return 0.0, [], False


def _apply_kind_mismatch_penalties(ctx: ScoringContext) -> Tuple[float, List[str], bool]:
    score = 0.0
    reasoning: List[str] = []
    kind_mismatch = False
    for penalty_func in (
        _penalty_error_section,
        _penalty_method_in_type_section,
        _penalty_func_in_type_section,
        _penalty_type_in_method_section,
    ):
        delta, delta_reasoning, mismatch = penalty_func(ctx)
        score += delta
        reasoning.extend(delta_reasoning)
        kind_mismatch = kind_mismatch or mismatch
    return score, reasoning, kind_mismatch


def _kind_section_map_matches(ctx: ScoringContext, kind_mismatch: bool) -> Tuple[float, List[str]]:
    if kind_mismatch:
        return 0.0, []
    kind_section_map = {
        "interface": [
            "Core Interfaces",
            "Type Definitions",
            "Metadata Types",
            "Generic Types",
            "Compression Types",
            "Encryption and Security Types",
            "Signature Types",
            "Streaming and Buffer Types",
            "Deduplication Types",
            "FileType System Types",
        ],
        "method": [
            "Methods",
            "File Management",
            "Package Writing",
            "Package Compression",
            "Package Metadata Methods",
            "Metadata Methods",
            "Basic Operations",
            "Security and Encryption Operations",
            "Digital Signatures",
            "Deduplication",
            "Streaming and Buffer Management",
        ],
        "type": ["Type Definitions", "Metadata Types"],
        "func": [
            "Basic Operations",
            "Metadata Helper Functions",
            "Package Metadata Methods",
            "Package Helper Functions",
            "File Management",
        ],
    }
    for kind_section in kind_section_map.get(ctx.definition.kind, []):
        if kind_section.lower() in ctx.section_lower:
            if (
                ctx.definition.kind in ("type", "struct", "interface")
                and "type definition" in ctx.section_lower
            ):
                msg = (
                    f"Kind '{ctx.definition.kind}' matches Type Definitions section: +20%"
                )
                return 0.20, [msg]
            return 0.15, [f"Kind '{ctx.definition.kind}' matches section type: +15%"]
    return 0.0, []


def score_kind_section_map(ctx: ScoringContext) -> Tuple[float, List[str]]:
    score = 0.0
    reasoning: List[str] = []
    delta, delta_reasoning, kind_mismatch = _apply_kind_mismatch_penalties(ctx)
    score += delta
    reasoning.extend(delta_reasoning)
    delta, delta_reasoning = _kind_section_map_matches(ctx, kind_mismatch)
    score += delta
    reasoning.extend(delta_reasoning)
    return score, reasoning


def score_general_heuristics(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if ctx.definition.kind != "func":
        return 0.0, []
    score = 0.0
    reasoning: List[str] = []
    if ctx.definition.file == "api_core.md" and "path" in ctx.name_lower:
        if any(token in ctx.name_lower for token in ["normalize", "todisplay", "validate"]):
            if "package helper function" in ctx.section_lower:
                score += 0.25
                reasoning.append(
                    "Path-related core functions prefer Package Helper Functions: +25%"
                )
            if "error helper" in ctx.section_lower:
                score -= 0.20
                reasoning.append("Path-related core functions are not error helpers: -20%")
    if "comment" in ctx.name_lower and "validate" in ctx.name_lower:
        if "metadata" in ctx.section_lower or "package metadata" in ctx.section_lower:
            score += 0.20
            reasoning.append("Comment validation prefers Metadata sections: +20%")
        if "encryption" in ctx.section_lower or "security" in ctx.section_lower:
            score -= 0.20
            reasoning.append("Comment validation is not encryption/security: -20%")
    return score, reasoning
