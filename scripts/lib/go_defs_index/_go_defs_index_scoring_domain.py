"""
Domain detection helpers for Go defs index scoring.
"""

from __future__ import annotations

import re
from typing import List, Optional, Tuple

from lib._go_code_utils import normalize_generic_name
from lib.go_defs_index._go_defs_index_config import (
    DOMAIN_FILE_MAP,
    KEYWORD_TO_SECTION_MAPPING,
)
from lib.go_defs_index._go_defs_index_models import DetectedDefinition

_DOMAIN_SUFFIXES: Tuple[str, ...] = ("strategy", "builder", "validator")
_DOMAIN_SUFFIX_KEYWORDS = {
    "compression": ["compression", "compress", "decompress"],
    "encryption": ["encryption", "encrypt", "security"],
    "signature": ["signature", "sign"],
    "streaming": ["stream", "streaming", "buffer"],
    "deduplication": ["dedup", "deduplication"],
}
_DOMAIN_FALLBACK_KEYWORDS = {
    "metadata": ["metadata", "comment", "tag", "pathmetadata", "fileentrytag"],
    "compression": ["compression", "compress", "decompress"],
    "encryption": ["encryption", "encrypt", "decrypt", "aes", "chacha", "mlkem", "cipher"],
    "signature": ["signature", "sign"],
    "streaming": ["streaming", "stream", "buffer", "chunk"],
    "deduplication": ["deduplication", "dedup"],
    "package": ["package"],
    "concurrency": ["concurrency", "thread", "worker", "safety"],
    "extraction": ["extract", "extraction"],
    "creation": ["create", "creation"],
    "filetype": ["filetype"],
}
_SECTION_DOMAIN_RULES = [
    ("metadata", ["package metadata", "metadata"]),
    ("generic", ["generic"]),
    ("streaming", ["streaming", "buffer"]),
    ("compression", ["compression", "compress", "decompress"]),
    ("encryption", ["encryption", "encrypt", "security"]),
    ("signature", ["signature", "sign"]),
    ("deduplication", ["deduplication", "dedup"]),
    ("filetype", ["filetype"]),
    ("writing", ["packagewriter", "package writing"]),
    ("package", ["package"]),
]


def _name_has_any(name_lower: str, keywords: List[str]) -> bool:
    return any(keyword in name_lower for keyword in keywords)


def _infer_domain_from_section_pattern(pattern: str) -> Optional[str]:
    pattern_lower = pattern.lower().strip()
    for domain, keywords in _SECTION_DOMAIN_RULES:
        if any(keyword in pattern_lower for keyword in keywords):
            return domain
    return None


def _domain_from_signature(definition: DetectedDefinition, normalized_lower: str) -> Optional[str]:
    if "signature" in normalized_lower or "signing" in normalized_lower:
        return "signature"
    if definition.receiver_type:
        receiver_normalized = normalize_generic_name(definition.receiver_type).lower()
        if "signature" in receiver_normalized or "signing" in receiver_normalized:
            return "signature"
    if definition.heading and "signature" in definition.heading.lower():
        return "signature"
    return None


def _domain_from_generic(definition: DetectedDefinition, normalized_name: str) -> Optional[str]:
    if "[" in normalized_name:
        return "generic"
    if definition.receiver_type:
        receiver_normalized = normalize_generic_name(definition.receiver_type)
        if "[" in receiver_normalized:
            return "generic"
    return None


def _domain_from_keyword_mapping(
    definition: DetectedDefinition,
    name_lower: str,
) -> Optional[str]:
    receiver_lower = ""
    if definition.receiver_type:
        receiver_lower = normalize_generic_name(definition.receiver_type).lower()
    name_for_matching = name_lower.replace(".", "")
    for keyword, mappings in KEYWORD_TO_SECTION_MAPPING.items():
        kw = keyword.replace(" ", "")
        if not kw:
            continue
        if kw in name_for_matching or (receiver_lower and kw in receiver_lower):
            for section_pattern, _weight in mappings:
                inferred = _infer_domain_from_section_pattern(section_pattern)
                if inferred:
                    return inferred
    return None


def _domain_from_suffix(name_lower: str) -> Optional[str]:
    if not any(suffix in name_lower for suffix in _DOMAIN_SUFFIXES):
        return None
    for domain, keywords in _DOMAIN_SUFFIX_KEYWORDS.items():
        if _name_has_any(name_lower, keywords):
            return domain
    return None


def _domain_from_fallback(name_lower: str) -> Optional[str]:
    for domain, keywords in _DOMAIN_FALLBACK_KEYWORDS.items():
        if _name_has_any(name_lower, keywords):
            return domain
    return None


def detect_definition_domain(definition: DetectedDefinition, name_lower: str) -> Optional[str]:
    normalized_name = normalize_generic_name(definition.name)
    normalized_lower = normalized_name.lower()
    signature_domain = _domain_from_signature(definition, normalized_lower)
    if signature_domain:
        return signature_domain
    if definition.file in DOMAIN_FILE_MAP:
        return DOMAIN_FILE_MAP[definition.file]
    generic_domain = _domain_from_generic(definition, normalized_name)
    if generic_domain:
        return generic_domain
    keyword_domain = _domain_from_keyword_mapping(definition, name_lower)
    if keyword_domain:
        return keyword_domain
    suffix_domain = _domain_from_suffix(name_lower)
    if suffix_domain:
        return suffix_domain
    return _domain_from_fallback(name_lower)


def extract_implementation_interface(definition: DetectedDefinition) -> Optional[str]:
    """
    Try to detect an interface name from section prose (implementation hints).
    """
    if not definition.section_content:
        return None

    struct_name = definition.name
    section_content = definition.section_content
    implementation_context = re.finditer(
        r"implements\\s+([A-Z][A-Za-z0-9]+)",
        section_content,
        re.IGNORECASE,
    )
    for match in implementation_context:
        interface_name = match.group(1)
        if interface_name != struct_name and interface_name in section_content:
            start_pos = max(0, match.start() - 200)
            end_pos = min(len(section_content), match.end() + 200)
            context = section_content[start_pos:end_pos]
            if struct_name.lower() in context.lower():
                return interface_name

    return None


def _extract_subsection_keywords(section: str) -> List[str]:
    """
    Extract keywords from subsection names for matching.
    """
    section_lower = section.lower()

    is_symlink_section = "symlink" in section_lower or "link" in section_lower
    is_path_metadata_section = "path metadata" in section_lower or "pathmetadata" in section_lower

    section_text = re.sub(r"^\\d+\\.\\d*\\s*", "", section_lower)
    section_text = re.sub(r"\\s+types?\\s*$", "", section_text)
    section_text = re.sub(r"\\s+errors?\\s*$", "", section_text)
    section_text = re.sub(r"\\s+methods?\\s*$", "", section_text)
    section_text = re.sub(r"\\s+operations?\\s*$", "", section_text)

    words = re.split(r"[\\s\\-]+", section_text)
    priority_keywords: List[str] = []
    generic_keywords: List[str] = []

    ambiguous_keywords = {
        "file",
        "path",
        "package",
        "type",
        "error",
        "basic",
        "aware",
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

    specific_domain_keywords = {
        "encryption",
        "encrypt",
        "decrypt",
        "compression",
        "compress",
        "decompress",
        "signature",
        "sign",
        "validation",
        "validate",
        "query",
        "queries",
        "information",
        "info",
        "lifecycle",
        "pattern",
        "streaming",
        "stream",
        "buffer",
        "deduplication",
        "writing",
        "write",
        "security",
        "mlkem",
        "ml-kem",
        "symlink",
        "link",
        "conversion",
        "convert",
        "concurrent",
        "status",
        "hash",
        "optional",
        "data",
        "tag",
        "value",
        "transform",
        "processing",
        "state",
        "comment",
        "pathmetadata",
        "path-metadata",
    }

    for word in words:
        word = word.strip()
        if len(word) < 3:
            continue
        if word in ["and", "the", "for", "with", "from", "that", "this", "to", "in", "on", "at"]:
            continue
        if re.match(r"^\\d+\\.\\d*$", word):
            continue

        if is_symlink_section:
            if word in ["symlink", "link", "conversion", "convert"]:
                priority_keywords.append(word)
            continue
        if is_path_metadata_section:
            if word in ["pathmetadata", "path-metadata", "path"]:
                priority_keywords.append(word)
            continue

        if word in specific_domain_keywords:
            priority_keywords.append(word)
        elif word in ambiguous_keywords:
            continue
        elif len(word) >= 4 and word not in ambiguous_keywords:
            if word in ["fileentry", "fileentry", "compression", "encryption", "signature"]:
                priority_keywords.append(word)

    return priority_keywords + generic_keywords


def extract_domain_keywords_from_subsection(section: str) -> List[str]:
    """
    Extract domain keywords from subsection names.
    """
    return _extract_subsection_keywords(section)


def extract_error_domain(error_name: str) -> Optional[str]:
    """
    Extract domain from error name.
    """
    name_lower = error_name.lower()
    name_clean = re.sub(r"^err", "", name_lower)
    name_clean = re.sub(r"error$", "", name_clean)

    domain_patterns = [
        r"file",
        r"compression",
        r"encryption",
        r"encrypt",
        r"signature",
        r"sign",
        r"package",
        r"security",
        r"validation",
        r"validate",
        r"streaming",
        r"stream",
        r"buffer",
        r"writing",
        r"write",
        r"basic",
        r"operation",
        r"management",
        r"metadata",
        r"path",
    ]

    for pattern in domain_patterns:
        if re.search(pattern, name_clean):
            return pattern

    return None
