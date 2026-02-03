"""
Method-specific scoring rules.
"""

from __future__ import annotations

from typing import List, Tuple

from lib.go_defs_index._go_defs_index_scoring_rules_core import ScoringContext


def score_method_patterns(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if ctx.definition.kind != "method" or "." not in ctx.definition.name:
        return 0.0, []
    method_name = ctx.definition.name.split(".", 1)[1].lower()
    operation_patterns = {
        "get": ["query", "information", "info", "queries"],
        "has": ["query", "information", "info", "queries"],
        "is": ["query", "information", "info", "queries", "validation", "validate"],
        "list": ["query", "information", "info", "queries"],
        "find": ["query", "information", "info", "queries"],
        "exists": ["query", "information", "info", "queries"],
        "set": ["configuration", "config", "management", "state"],
        "add": ["operations", "management", "basic"],
        "remove": ["operations", "management", "basic"],
        "create": ["lifecycle", "creation", "operations"],
        "new": ["lifecycle", "creation", "operations"],
        "open": ["lifecycle"],
        "close": ["lifecycle"],
        "write": ["write", "writing", "operations"],
        "read": ["streaming", "stream", "operations"],
        "validate": ["validation", "validate", "verify"],
        "verify": ["validation", "validate", "verify"],
        "compress": ["compression", "compress"],
        "decompress": ["compression", "compress"],
        "encrypt": ["encryption", "encrypt"],
        "decrypt": ["encryption", "encrypt"],
    }
    for pattern, subsection_keywords in operation_patterns.items():
        pattern_matches = method_name.startswith(pattern) or pattern in method_name
        if not pattern_matches:
            continue
        for kw in subsection_keywords:
            if kw not in ctx.section_lower:
                continue
            if pattern in ["get", "has", "is", "list", "find", "exists"]:
                if (
                    "information" in ctx.section_lower
                    or "queries" in ctx.section_lower
                    or "query" in ctx.section_lower
                ):
                    msg = (
                        f"Method pattern '{pattern}*' matches query/info "
                        f"subsection: +25%"
                    )
                    return 0.25, [msg]
            elif pattern in ["add", "remove", "set"]:
                msg = (
                    f"Method pattern '{pattern}*' matches operation subsection: +15%"
                )
                return 0.15, [msg]
            msg = f"Method pattern '{pattern}*' matches subsection: +15%"
            return 0.15, [msg]
        break
    return 0.0, []


def _method_flags(method_name: str) -> Tuple[bool, bool]:
    method_lower = method_name.lower()
    is_getter = (
        method_name.startswith("Get")
        or method_name.startswith("Is")
        or method_name.startswith("Has")
        or method_lower.startswith("get")
        or method_lower.startswith("is")
        or method_lower.startswith("has")
    )
    is_transformation = (
        method_name.startswith("Add")
        or method_name.startswith("Update")
        or method_name.startswith("Set")
        or method_name.startswith("Remove")
        or method_name.startswith("Delete")
        or method_name.startswith("Modify")
        or method_lower.startswith("add")
        or method_lower.startswith("update")
        or method_lower.startswith("set")
        or method_lower.startswith("remove")
        or method_lower.startswith("delete")
        or method_lower.startswith("modify")
    )
    return is_getter, is_transformation


def _score_getter_section(ctx: ScoringContext, is_getter: bool) -> Tuple[float, List[str]]:
    if not is_getter:
        return 0.0, []
    if "query methods" in ctx.section_lower:
        msg = "Getter method (Get/Is/Has) matches Query Methods: +25%"
        return 0.25, [msg]
    if "data methods" in ctx.section_lower:
        msg = "Getter method (Get/Is/Has) matches Data Methods: +15%"
        return 0.15, [msg]
    if "transformation methods" in ctx.section_lower:
        msg = "Getter method (Get/Is/Has) does not match Transformation Methods: -25%"
        return -0.25, [msg]
    return 0.0, []


def _score_transformation_section(
    ctx: ScoringContext,
    is_transformation: bool,
) -> Tuple[float, List[str]]:
    if not is_transformation:
        return 0.0, []
    if "transformation methods" in ctx.section_lower:
        msg = (
            "Transformation method (Add/Update/Set/Remove) matches "
            "Transformation Methods: +25%"
        )
        return 0.25, [msg]
    if "data methods" in ctx.section_lower or "query methods" in ctx.section_lower:
        msg = (
            "Transformation method (Add/Update/Set/Remove) does not match "
            "Query/Data Methods: -25%"
        )
        return -0.25, [msg]
    return 0.0, []


def score_method_type_classification(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if ctx.definition.kind != "method" or "." not in ctx.definition.name:
        return 0.0, []
    method_name = ctx.definition.name.split(".", 1)[1]
    is_getter, is_transformation = _method_flags(method_name)
    score, reasoning = _score_getter_section(ctx, is_getter)
    if score or reasoning:
        return score, reasoning
    return _score_transformation_section(ctx, is_transformation)


def score_method_name_preferences(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if ctx.definition.kind != "method" or "." not in ctx.definition.name:
        return 0.0, []
    score = 0.0
    reasoning: List[str] = []
    method_name = ctx.definition.name.split(".", 1)[1].lower()
    if "metadata" in method_name:
        if "metadata" in ctx.section_lower and "compression" not in ctx.section_lower:
            score += 0.30
            reasoning.append(
                "Method name contains 'metadata' matches Metadata section: +30%"
            )
        elif "compression" in ctx.section_lower and "metadata" not in ctx.section_lower:
            score -= 0.30
            reasoning.append(
                "Method name contains 'metadata' does not match Compression section: -30%"
            )
    if "signaturefile" in method_name or "signature file" in method_name:
        if "metadata" in ctx.section_lower:
            score += 0.30
            reasoning.append(
                "Signature file method (special metadata) matches Metadata section: +30%"
            )
        elif "file management" in ctx.section_lower and "metadata" not in ctx.section_lower:
            score -= 0.30
            reasoning.append(
                "Signature file method (special metadata) does not match "
                "File Management section: -30%"
            )
    if "compression" in ctx.section_lower:
        if method_name.startswith(("get", "is", "has")):
            compression_keywords = [
                "compress",
                "decompress",
                "compressiontype",
                "compressionratio",
            ]
            is_operation = any(kw in method_name for kw in compression_keywords)
            if not is_operation:
                if "information" in ctx.section_lower or "queries" in ctx.section_lower:
                    score += 0.20
                    reasoning.append(
                        "Getter method about compression info prefers Information/Queries: +20%"
                    )
                elif "method" in ctx.section_lower:
                    score -= 0.20
                    reasoning.append(
                        "Getter method about compression info does not match "
                        "Compression Methods: -20%"
                    )
    return score, reasoning


def score_file_entry_method_categories(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if ctx.definition.kind != "method" or "." not in ctx.definition.name:
        return 0.0, []
    receiver = ctx.definition.receiver_type or ctx.definition.name.split(".", 1)[0]
    if receiver.lower() != "fileentry":
        return 0.0, []
    method_lower = ctx.definition.name.split(".", 1)[1].lower()
    category_rules = [
        (
            "query methods",
            ["get", "has", "is"],
        ),
        (
            "data methods",
            ["getdata", "setdata", "loaddata", "unloaddata", "data"],
        ),
        (
            "temp file methods",
            ["tempfile", "temp"],
        ),
        (
            "serialization methods",
            ["marshal", "writedata", "writemeta", "writeto"],
        ),
        (
            "path methods",
            ["path", "symlink", "associate", "resolve"],
        ),
        (
            "transformation methods",
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
    ]
    for section_keyword, tokens in category_rules:
        if section_keyword in ctx.section_lower and any(token in method_lower for token in tokens):
            return 0.40, [f"FileEntry {section_keyword} method matches section: +40%"]
    return 0.0, []
