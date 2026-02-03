"""
Domain-focused scoring rules for Go defs index placement.
"""

from __future__ import annotations

from typing import Dict, List, Tuple

from lib.go_defs_index._go_defs_index_config import KEYWORD_TO_SECTION_MAPPING
from lib.go_defs_index._go_defs_index_scoring_domain import (
    extract_domain_keywords_from_subsection,
)
from lib.go_defs_index._go_defs_index_scoring_rules_core_base import (
    ScoringContext,
    _is_core_package_type,
)
from lib.go_defs_index._go_defs_index_scoring_text import (
    extract_keywords_from_comments,
    extract_section_level,
    match_keywords_to_section,
)


def _section_contains_any(section_lower: str, keywords: List[str]) -> bool:
    return any(keyword in section_lower for keyword in keywords)


def score_keyword_comment_matching(
    ctx: ScoringContext,
    kind_mismatch: bool,
) -> Tuple[float, List[str]]:
    if kind_mismatch:
        return 0.0, []
    section_level = extract_section_level(ctx.section, ctx.definition.kind)
    keywords = extract_keywords_from_comments(ctx.definition)
    if not keywords:
        return 0.0, []
    keyword_score, keyword_reasoning = match_keywords_to_section(
        keywords,
        section_level,
        KEYWORD_TO_SECTION_MAPPING,
    )
    score = keyword_score
    reasoning = list(keyword_reasoning)

    name_lower = ctx.definition.name.lower()
    is_error_function_name = (
        (name_lower.startswith("as") and "error" in name_lower)
        or (name_lower.startswith("get") and "error" in name_lower)
        or (name_lower.startswith("add") and "error" in name_lower and "context" in name_lower)
        or ("error" in name_lower and "map" in name_lower)
    )
    if (
        is_error_function_name
        and "error" in ctx.section_lower
        and "helper" in ctx.section_lower
        and "function" in ctx.section_lower
    ):
        score += 0.25
        reasoning.append("Error keyword match => Error Helper Functions: +25%")

    return score, reasoning


def score_constructor_functions(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if ctx.definition.kind != "func" or not ctx.definition.name.startswith("New"):
        return 0.0, []
    score = 0.0
    reasoning: List[str] = []
    constructed_type = ctx.definition.name[3:]
    constructed_lower = constructed_type.lower()

    if constructed_lower in ctx.section_lower:
        score += 0.25
        msg = (
            f"Constructor function 'New{constructed_type}' matches constructed type "
            f"in section: +25%"
        )
        reasoning.append(msg)
    if "helper" in ctx.section_lower and "function" in ctx.section_lower:
        if "constructor" in ctx.section_lower or "package" in ctx.section_lower:
            score += 0.15
            reasoning.append("Constructor function matches constructor/helper section: +15%")
    if "package" in constructed_lower and "package" in ctx.section_lower:
        score += 0.20
        reasoning.append("Package-related constructor matches Package section: +20%")
    return score, reasoning


def _score_domain_metadata(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if not _section_contains_any(ctx.section_lower, ["metadata", "comment", "pathmetadata"]):
        return 0.0, []
    return 0.30, ["Domain match: metadata-related => Metadata section: +30%"]


def _score_domain_compression(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if not _section_contains_any(ctx.section_lower, ["compression", "compress"]):
        return 0.0, []
    return 0.30, ["Domain match: compression-related => Compression Types: +30%"]


def _score_domain_generic(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if "generic" not in ctx.section_lower:
        return 0.0, []
    return 0.30, ["Domain match: generic-related => Generic Types: +30%"]


def _score_domain_package(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if "package" not in ctx.section_lower:
        return 0.0, []
    if _is_core_package_type(
        ctx.definition.name,
        ctx.definition.kind,
        ctx.definition.receiver_type,
    ):
        return 0.20, ["Domain match: package-related => Package section: +20%"]
    return 0.10, ["Domain match: package-related => Package section: +10% (weak)"]


def _score_domain_creation(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if "package" not in ctx.section_lower:
        return 0.0, []
    return 0.20, ["Domain match: creation-related => Package section: +20%"]


def _score_domain_extraction(ctx: ScoringContext) -> Tuple[float, List[str]]:
    keywords = ["extract", "fileentry"]
    if not _section_contains_any(ctx.section_lower, keywords):
        return 0.0, []
    msg = "Domain match: extraction-related => Extraction/FileEntry section: +20%"
    return 0.20, [msg]


def _score_domain_concurrency(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if "generic" not in ctx.section_lower:
        return 0.0, []
    return 0.20, ["Domain match: concurrency-related => Generic Types section: +20%"]


def _score_domain_encryption(ctx: ScoringContext) -> Tuple[float, List[str]]:
    keywords = ["encryption", "security", "encrypt"]
    if not _section_contains_any(ctx.section_lower, keywords):
        return 0.0, []
    msg = "Domain match: encryption-related => Encryption and Security: +30%"
    return 0.30, [msg]


def _score_domain_signature(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if not _section_contains_any(ctx.section_lower, ["signature", "sign"]):
        return 0.0, []
    return 0.30, ["Domain match: signature-related => Signature Types: +30%"]


def _score_domain_streaming(ctx: ScoringContext) -> Tuple[float, List[str]]:
    keywords = ["streaming", "stream", "buffer"]
    if not _section_contains_any(ctx.section_lower, keywords):
        return 0.0, []
    return 0.30, ["Domain match: streaming-related => Streaming and Buffer: +30%"]


def _score_domain_deduplication(ctx: ScoringContext) -> Tuple[float, List[str]]:
    keywords = ["file management", "information and queries", "package file management"]
    if not _section_contains_any(ctx.section_lower, keywords):
        return 0.0, []
    return 0.20, [
        "Domain match: deduplication-related => Package file management/queries: +20%"
    ]


def _score_domain_filetype(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if "filetype" not in ctx.section_lower:
        return 0.0, []
    return 0.30, ["Domain match: filetype-related => FileType System Types: +30%"]


def _score_domain_writing(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if not _section_contains_any(ctx.section_lower, ["package write methods", "writing"]):
        return 0.0, []
    return 0.20, ["Domain match: writing-related => Package Write Methods: +20%"]


def score_domain_match(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if ctx.detected_domain is None:
        return 0.0, []
    scorers = {
        "metadata": _score_domain_metadata,
        "compression": _score_domain_compression,
        "generic": _score_domain_generic,
        "package": _score_domain_package,
        "creation": _score_domain_creation,
        "extraction": _score_domain_extraction,
        "concurrency": _score_domain_concurrency,
        "encryption": _score_domain_encryption,
        "signature": _score_domain_signature,
        "streaming": _score_domain_streaming,
        "deduplication": _score_domain_deduplication,
        "filetype": _score_domain_filetype,
        "writing": _score_domain_writing,
    }
    scorer = scorers.get(ctx.detected_domain)
    if not scorer:
        return 0.0, []
    return scorer(ctx)


def score_type_name_patterns(ctx: ScoringContext) -> Tuple[float, List[str]]:
    """Score based on type name patterns like *Config, *Builder, *Strategy."""
    if ctx.definition.kind not in ("type", "struct", "interface"):
        return 0.0, []

    name_lower = ctx.name_lower
    patterns = {
        "config": (["compression", "encryption", "streaming", "signature", "package"], 0.15),
        "builder": (["compression", "encryption", "config", "signature", "streaming"], 0.10),
        "strategy": (["compression", "encryption", "signature", "streaming"], 0.15),
        "validator": (["compression", "encryption", "validation", "signature"], 0.10),
        "handler": (["encryption", "file"], 0.10),
        "pool": (["buffer", "compression", "resource", "streaming", "worker"], 0.10),
        "errorcontext": (["error"], 0.20),
        "options": (["file", "package", "compression", "extraction"], 0.10),
        "info": (["compression", "file", "package", "signature"], 0.10),
    }

    for suffix, (domain_keywords, bonus) in patterns.items():
        if not name_lower.endswith(suffix):
            continue
        for keyword in domain_keywords:
            if keyword in ctx.section_lower:
                reason = (
                    f"Type pattern '*{suffix}' matches section domain: +{int(bonus * 100)}%"
                )
                return bonus, [reason]
    return 0.0, []


def score_domain_type_subsection(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if ctx.definition.kind not in ("type", "struct"):
        return 0.0, []
    if "type" not in ctx.section_lower:
        return 0.0, []
    if "definition" not in ctx.section_lower and "definitions" not in ctx.section_lower:
        return 0.0, []
    domain_keywords = extract_domain_keywords_from_subsection(ctx.section)
    if not domain_keywords:
        return 0.0, []
    domain_priority_keywords: Dict[str, List[str]] = {
        "metadata": ["metadata", "comment", "tag", "pathmetadata", "fileentrytag"],
        "compression": ["compression", "compress", "decompress"],
        "encryption": ["encryption", "encrypt", "decrypt", "aes", "chacha", "mlkem", "cipher"],
        "security": ["security", "validation", "validate", "verify"],
        "signature": ["signature", "sign"],
        "streaming": ["streaming", "stream", "buffer", "chunk"],
        "deduplication": ["deduplication", "dedup"],
        "package": ["package"],
        "concurrency": ["concurrency", "thread", "worker", "safety"],
        "extraction": ["extract", "extraction"],
        "creation": ["create", "creation"],
        "generic": ["generic"],
        "filetype": ["filetype"],
        "writing": ["write", "writing"],
    }
    for keyword in domain_keywords:
        if keyword in ctx.name_lower:
            if ctx.detected_domain and keyword in domain_priority_keywords.get(
                ctx.detected_domain,
                [],
            ):
                msg = (
                    f"Priority domain keyword '{keyword}' in type name "
                    f"matches subsection: +30%"
                )
                return 0.30, [msg]
            msg = f"Domain keyword '{keyword}' in type name matches subsection: +15%"
            return 0.15, [msg]
    return 0.0, []
