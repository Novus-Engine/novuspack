"""
Section/heading-related scoring rules.
"""

from __future__ import annotations

import re
from typing import List, Set, Tuple

from lib._go_code_utils import normalize_generic_name
from lib.go_defs_index._go_defs_index_scoring_domain import (
    extract_domain_keywords_from_subsection,
)
from lib.go_defs_index._go_defs_index_scoring_rules_core import (
    ScoringContext,
    _is_core_package_type,
)
from lib.go_defs_index._go_defs_index_scoring_text import (
    extract_code_comments,
    extract_heading_keywords,
    extract_prose_keywords,
)
from lib.go_defs_index._go_defs_index_shared import map_implementation_to_interface


_FILE_SECTION_PATTERNS = [
    ("file_mgmt_error", ["Error Types"]),
    ("file_mgmt_compression", ["Compression Types", "FileEntry Types"]),
    (
        "file_mgmt",
        [
            "FileEntry Types",
            "FileEntry Query Methods",
            "FileEntry Data Methods",
            "FileEntry Temp File Methods",
            "FileEntry Serialization Methods",
            "FileEntry Path Methods",
            "FileEntry Transformation Methods",
            "FileEntry Helper Functions",
            "Tag Methods",
        ],
    ),
    (
        "core",
        [
            "Package Interface Types",
            "Package Lifecycle Methods",
            "Package File Management Methods",
            "Package Information and Queries Methods",
            "Package Comment Methods",
            "Package Identity Methods",
            "Package Special File Methods",
            "Package Path Metadata Methods",
            "Package Symlink Methods",
            "Package Metadata-Only Methods",
            "Package Info Methods",
            "Package Metadata Validation Methods",
            "Package Metadata Internal Methods",
            "Package Compression Methods",
            "Package Path and Configuration Methods",
            "Package File Encryption Methods",
            "Package Signature Management Methods",
            "Package Write Methods",
            "Package Other Methods",
            "Package Helper Functions",
            "Error Types",
        ],
    ),
    (
        "basic_operation",
        [
            "Package Interface Types",
            "Package Lifecycle Methods",
            "Package File Management Methods",
            "Package Information and Queries Methods",
            "Package Helper Functions",
        ],
    ),
    (
        "basic_operations",
        [
            "Package Interface Types",
            "Package Lifecycle Methods",
            "Package File Management Methods",
            "Package Information and Queries Methods",
            "Package Helper Functions",
        ],
    ),
    (
        "file_format",
        [
            "Package Interface Types",
            "Package Information and Queries Methods",
            "Package File Management Methods",
        ],
    ),
    ("compression", ["Compression Types", "Compression Methods", "Compression Helper Functions"]),
    (
        "streaming",
        [
            "Streaming and Buffer Types",
            "Streaming and Buffer Methods",
            "Streaming and Buffer Helper Functions",
        ],
    ),
    (
        "security",
        [
            "Encryption and Security Types",
            "Encryption and Security Methods",
            "Encryption and Security Helper Functions",
        ],
    ),
    (
        "encryption",
        [
            "Encryption and Security Types",
            "Encryption and Security Methods",
            "Encryption and Security Helper Functions",
        ],
    ),
    ("signature", ["Signature Types", "Signature Methods", "Signature Helper Functions"]),
    (
        "metadata",
        [
            "Package Metadata Types",
            "Package Comment Methods",
            "Package Identity Methods",
            "Package Special File Methods",
            "Package Path Metadata Methods",
            "Package Symlink Methods",
            "Package Metadata-Only Methods",
            "Package Info Methods",
            "Package Metadata Validation Methods",
            "Package Metadata Internal Methods",
            "Package Metadata Type Methods",
            "Package Metadata Helper Functions",
            "Package Interface Types",
        ],
    ),
    (
        "deduplication",
        [
            "Package File Management Methods",
            "Package Information and Queries Methods",
        ],
    ),
    ("generic", ["Generic Types", "Generic Methods", "Generic Helper Functions"]),
    (
        "writing",
        [
            "Package Write Methods",
            "Package Helper Functions",
        ],
    ),
    (
        "file_type",
        ["FileType System Types", "FileType System Methods", "FileType System Helper Functions"],
    ),
]


def _file_tokens(file_name: str) -> List[str]:
    normalized = file_name.lower().replace(".md", "").replace("-", "_")
    return [token for token in normalized.split("_") if token]


def _tokens_contain_sequence(tokens: List[str], pattern_tokens: List[str]) -> bool:
    if not tokens or not pattern_tokens:
        return False
    start_index = 0
    for token in pattern_tokens:
        try:
            found_index = tokens.index(token, start_index)
        except ValueError:
            return False
        start_index = found_index + 1
    return True


def score_file_patterns(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if not ctx.definition.file:
        return 0.0, []
    section_leaf = ctx.section_lower.split(">")[-1].strip()
    file_tokens = _file_tokens(ctx.definition.file)
    for pattern, section_keywords in _FILE_SECTION_PATTERNS:
        pattern_tokens = _file_tokens(pattern)
        if not _tokens_contain_sequence(file_tokens, pattern_tokens):
            continue
        for section_keyword in section_keywords:
            keyword_lower = section_keyword.lower()
            if keyword_lower in section_leaf or section_leaf in keyword_lower:
                return 0.15, [f"File pattern match ({ctx.definition.file}): +15%"]
    return 0.0, []


def score_heading_match(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if not ctx.heading_lower:
        return 0.0, []
    score = 0.0
    reasoning: List[str] = []
    section_text = re.sub(r"^\\d+\\.\\s*", "", ctx.section_lower)
    heading_keywords = extract_heading_keywords(ctx.definition.heading)

    domain_mismatch = False
    if (
        "signature" in ctx.heading_lower
        and "signature" not in ctx.section_lower
        and "sign" not in ctx.section_lower
    ):
        if ctx.definition.kind == "method" and "file" in ctx.name_lower:
            domain_mismatch = True
    elif "encryption" in ctx.heading_lower or "encrypt" in ctx.heading_lower:
        if "encryption" not in ctx.section_lower and "encrypt" not in ctx.section_lower:
            if "file" in ctx.name_lower and "file" not in ctx.section_lower:
                domain_mismatch = True

    if domain_mismatch:
        return -0.2, ["Heading domain mismatch (heading mentions different domain): -20%"]

    if ctx.heading_lower in section_text or section_text in ctx.heading_lower:
        score += 0.20
        reasoning.append(f"Heading text match with '{ctx.section}': +20%")
    elif heading_keywords:
        matched_keywords = [kw for kw in heading_keywords if kw in ctx.section_lower]
        if matched_keywords:
            score += 0.15
            msg = f"Heading keyword match ({', '.join(matched_keywords)} in section): +15%"
            reasoning.append(msg)

    return score, reasoning


def score_camelcase_match(ctx: ScoringContext) -> Tuple[float, List[str]]:
    def split_camel_case_words(text: str) -> List[str]:
        if not text:
            return []
        normalized = normalize_generic_name(text)
        normalized = re.sub(r"[^a-zA-Z0-9]", "", normalized)
        parts = re.findall(r"[A-Z]+[a-z0-9]*|[a-z0-9]+", normalized)
        return [p.lower() for p in parts if p]

    section_leaf = ctx.section.split(">")[-1].strip()
    section_leaf = re.sub(r"^\\d+(?:\\.\\d+)*\\.?\\s+", "", section_leaf)
    section_leaf_tokens = set(re.findall(r"\\b[a-z0-9]+\\b", section_leaf.lower()))

    camel_sources: List[str] = []
    if ctx.definition.kind == "method":
        if ctx.definition.receiver_type:
            camel_sources.append(map_implementation_to_interface(ctx.definition.receiver_type))
        else:
            camel_sources.append(ctx.definition.name.split(".", 1)[0])
    else:
        camel_sources.append(map_implementation_to_interface(ctx.definition.name))

    camel_words: List[str] = []
    for src in camel_sources:
        camel_words.extend(split_camel_case_words(src))

    camel_words = [w for w in camel_words if len(w) >= 4]
    matched_words = sorted({w for w in camel_words if w in section_leaf_tokens})
    if not matched_words:
        return 0.0, []
    camel_score = min(0.15 * len(matched_words), 0.30)
    msg = (
        f"camelCase word match ({', '.join(matched_words)} in section): "
        f"+{int(camel_score * 100)}%"
    )
    return camel_score, [msg]


def score_parent_heading_match(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if not ctx.definition.parent_heading:
        return 0.0, []
    parent_lower = ctx.definition.parent_heading.lower()
    section_text = re.sub(r"^\\d+\\.\\s*", "", ctx.section_lower)

    if parent_lower in section_text or section_text in parent_lower:
        return 0.15, [f"Parent heading match with '{ctx.section}': +15%"]
    if "pathmetadata" in parent_lower and "metadata" in section_text:
        msg = "Parent heading related term match ('PathMetadata' -> 'Metadata'): +15%"
        return 0.15, [msg]
    if "metadata" in parent_lower and "metadata" in section_text:
        return 0.15, ["Parent heading metadata term match: +15%"]
    return 0.0, []


def score_comment_domain_match(ctx: ScoringContext) -> Tuple[float, List[str]]:
    comment_text = extract_code_comments(ctx.definition)
    if not comment_text:
        return 0.0, []
    comment_lower = comment_text.lower()
    domain_priority_keywords = {
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
    score = 0.0
    reasoning: List[str] = []
    domains = [
        domain
        for domain, keywords in domain_priority_keywords.items()
        if any(kw in comment_lower for kw in keywords)
    ]
    for domain in domains:
        if domain == "compression" and (
            "compression" in ctx.section_lower or "compress" in ctx.section_lower
        ):
            score += 0.15
            reasoning.append("Code comment domain match (compression): +15%")
        elif domain == "encryption" and (
            "encryption" in ctx.section_lower or "security" in ctx.section_lower
        ):
            score += 0.15
            reasoning.append("Code comment domain match (encryption): +15%")
        elif domain == "package" and "package" in ctx.section_lower:
            if _is_core_package_type(
                ctx.definition.name,
                ctx.definition.kind,
                ctx.definition.receiver_type,
            ):
                score += 0.15
                reasoning.append("Code comment domain match (package): +15%")
            else:
                score += 0.10
                reasoning.append("Code comment domain match (package): +10% (weak)")

    comment_words = re.findall(r"\\b[a-z]{4,}\\b", comment_lower)
    for word in comment_words[:5]:
        excluded_words = ["type", "types", "method", "methods", "function", "functions"]
        if word in ctx.section_lower and word not in excluded_words:
            score += 0.10
            reasoning.append(f"Code comment keyword '{word}' matches section: +10%")
            break

    return score, reasoning


def score_prose_keyword_match(ctx: ScoringContext) -> Tuple[float, List[str]]:
    if not ctx.definition.section_content:
        return 0.0, []
    prose_keywords = extract_prose_keywords(ctx.definition.section_content)
    if not prose_keywords:
        return 0.0, []
    matched = []
    for keyword in prose_keywords[:5]:
        excluded_words = [
            "type",
            "types",
            "method",
            "methods",
            "function",
            "functions",
            "section",
        ]
        if keyword in ctx.section_lower and keyword not in excluded_words:
            matched.append(keyword)
    if not matched:
        return 0.0, []
    msg = f"Prose keyword match ({', '.join(matched[:3])} in section): +15%"
    return 0.15, [msg]


def score_content_keyword_match(ctx: ScoringContext) -> Tuple[float, List[str]]:
    keywords = {
        "compression": ["Package Compression", "Compression"],
        "streaming": ["Streaming and Buffer Management"],
        "security": ["Security and Encryption Operations"],
        "encryption": ["Security and Encryption Operations"],
        "signature": ["Digital Signatures"],
        "deduplication": ["Deduplication"],
        "writing": ["Package Writing"],
    }
    keyword_matches = 0
    score = 0.0
    reasoning: List[str] = []
    for keyword, sections in keywords.items():
        if keyword not in ctx.name_lower:
            continue
        for kw_section in sections:
            if kw_section.lower() in ctx.section_lower:
                keyword_matches += 1
                if keyword_matches <= 2:
                    score += 0.05
                    reasoning.append(
                        f"Content keyword '{keyword}' in definition name: +5%"
                    )
                break
        if keyword_matches >= 2:
            break
    return score, reasoning


def _match_subsection_in_name(
    ctx: ScoringContext,
    subsection_keywords: List[str],
    matched_keywords: Set[str],
    method_name_part: str,
) -> Tuple[float, List[str], int]:
    score = 0.0
    reasoning: List[str] = []
    matches = 0
    ambiguous_keywords = {
        "file",
        "path",
        "package",
        "type",
        "error",
        "basic",
        "management",
        "metadata",
        "operations",
        "operation",
        "methods",
        "method",
        "types",
        "definitions",
        "definition",
    }
    for keyword in subsection_keywords:
        if keyword not in ctx.name_lower or keyword in matched_keywords:
            continue
        is_ambiguous = keyword in ambiguous_keywords
        if method_name_part and keyword in method_name_part:
            if is_ambiguous:
                continue
            score += 0.30
            reasoning.append(f"Subsection keyword '{keyword}' in method name: +30%")
        else:
            if is_ambiguous:
                continue
            score += 0.20
            reasoning.append(f"Subsection keyword '{keyword}' in definition name: +20%")
        matched_keywords.add(keyword)
        matches += 1
        if matches >= 2:
            break
    return score, reasoning, matches


def _match_subsection_in_heading(
    ctx: ScoringContext,
    subsection_keywords: List[str],
    matched_keywords: Set[str],
) -> Tuple[float, List[str], int]:
    score = 0.0
    reasoning: List[str] = []
    matches = 0
    ambiguous_keywords = {
        "file",
        "path",
        "package",
        "type",
        "error",
        "basic",
        "management",
        "metadata",
        "operations",
        "operation",
        "methods",
        "method",
        "types",
        "definitions",
        "definition",
    }
    domain_specific_keywords = {
        "signature",
        "encryption",
        "encrypt",
        "compression",
        "compress",
        "deduplication",
        "streaming",
        "stream",
        "security",
        "validation",
    }
    if not ctx.heading_lower:
        return 0.0, [], 0
    for keyword in subsection_keywords:
        if keyword in ambiguous_keywords or keyword in matched_keywords:
            continue
        if keyword not in ctx.heading_lower:
            continue
        if keyword in domain_specific_keywords and keyword not in ctx.name_lower:
            continue
        score += 0.10
        reasoning.append(f"Subsection keyword '{keyword}' in heading: +10%")
        matched_keywords.add(keyword)
        matches += 1
        if matches >= 1:
            break
    return score, reasoning, matches


def _match_subsection_in_content(
    ctx: ScoringContext,
    subsection_keywords: List[str],
    matched_keywords: Set[str],
) -> Tuple[float, List[str], int]:
    score = 0.0
    reasoning: List[str] = []
    matches = 0
    ambiguous_keywords = {
        "file",
        "path",
        "package",
        "type",
        "error",
        "basic",
        "management",
        "metadata",
        "operations",
        "operation",
        "methods",
        "method",
        "types",
        "definitions",
        "definition",
    }
    domain_specific_keywords = {
        "signature",
        "encryption",
        "encrypt",
        "compression",
        "compress",
        "deduplication",
        "streaming",
        "stream",
        "security",
        "validation",
    }
    for keyword in subsection_keywords:
        if keyword in ambiguous_keywords or keyword in matched_keywords:
            continue
        if keyword not in ctx.content_lower:
            continue
        if keyword in domain_specific_keywords and keyword not in ctx.name_lower:
            continue
        score += 0.05
        reasoning.append(f"Subsection keyword '{keyword}' in content: +5%")
        matched_keywords.add(keyword)
        matches += 1
        if matches >= 1:
            break
    return score, reasoning, matches


def score_subsection_keyword_match(ctx: ScoringContext) -> Tuple[float, List[str]]:
    subsection_keywords = extract_domain_keywords_from_subsection(ctx.section)
    if not subsection_keywords:
        return 0.0, []
    score = 0.0
    reasoning: List[str] = []
    matched_keywords: Set[str] = set()
    method_name_part = ""
    if ctx.definition.kind == "method" and "." in ctx.definition.name:
        method_name_part = ctx.definition.name.split(".", 1)[1].lower()
    delta, delta_reasoning, matches = _match_subsection_in_name(
        ctx,
        subsection_keywords,
        matched_keywords,
        method_name_part,
    )
    score += delta
    reasoning.extend(delta_reasoning)
    if matches < 1:
        delta, delta_reasoning, matches = _match_subsection_in_heading(
            ctx,
            subsection_keywords,
            matched_keywords,
        )
        score += delta
        reasoning.extend(delta_reasoning)
    if not matches:
        delta, delta_reasoning, _matches = _match_subsection_in_content(
            ctx,
            subsection_keywords,
            matched_keywords,
        )
        score += delta
        reasoning.extend(delta_reasoning)
    return score, reasoning
