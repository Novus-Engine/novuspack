"""
Text extraction and keyword scoring helpers for Go defs index matching.
"""

from __future__ import annotations

import re
from typing import Dict, List, Optional, Set, Tuple
from lib.go_defs_index._go_defs_index_config import (
    KEYWORD_TO_SECTION_MAPPING,
    PRIORITY_PHRASES,
)
from lib.go_defs_index._go_defs_index_models import DetectedDefinition


def normalize_keyword(keyword: str) -> str:
    """
    Normalize keyword to lowercase while preserving compound words.

    Detects camelCase compounds (FileEntry -> fileentry) and all-caps acronyms
    (ACL -> acl, ML-KEM -> mlkem).
    """
    if not keyword:
        return ""

    known_compounds = {
        "fileentry",
        "pathmetadata",
        "fileentrytag",
        "pathmetadatatag",
        "bufferpool",
        "mlkem",
        "appid",
        "vendorid",
        "mimetype",
        "filetype",
        "errorcontext",
        "accesscontrol",
        "accesscontrollist",
    }

    normalized = keyword.lower()

    if normalized in known_compounds:
        return normalized

    if re.search(r"[a-z][A-Z]", keyword):
        spaced = re.sub(r"([a-z])([A-Z])", r"\\1 \\2", keyword)
        normalized = spaced.lower().replace(" ", "")
        return normalized

    if keyword.isupper() or (keyword.replace("-", "").isupper() and "-" in keyword):
        normalized = keyword.lower().replace("-", "")
        return normalized

    if "-" in normalized:
        normalized = normalized.replace("-", "")

    return normalized


def extract_heading_keywords(heading: str) -> List[str]:
    """
    Extract meaningful keywords from a heading.

    Removes numbers and common filler words.
    """
    if not heading:
        return []
    heading_clean = re.sub(r"^\\d+(?:\\.\\d+)*\\.?\\s*", "", heading)
    heading_clean = heading_clean.lower()
    words = re.findall(r"\\b[a-z]{3,}\\b", heading_clean)
    common_words = {
        "the",
        "and",
        "for",
        "with",
        "from",
        "that",
        "this",
        "are",
        "was",
        "were",
        "has",
        "have",
        "had",
        "but",
        "not",
        "all",
        "any",
        "can",
        "will",
        "may",
        "should",
        "must",
        "each",
        "such",
        "when",
        "where",
        "what",
        "which",
        "who",
        "why",
        "how",
        "method",
        "methods",
        "function",
        "functions",
        "type",
        "types",
        "definition",
        "definitions",
        "section",
        "subsection",
        "overview",
        "summary",
        "general",
        "core",
        "other",
    }
    keywords = [w for w in words if w not in common_words]
    return keywords[:5]


def extract_prose_keywords(section_content: str) -> List[str]:
    """
    Extract meaningful keywords from prose text in section (before code block).

    Args:
        section_content: Full section content (heading to next heading)

    Returns:
        List of meaningful keywords extracted from prose
    """
    if not section_content:
        return []
    lines = section_content.split("\\n")
    prose_lines = []
    for line in lines:
        if line.strip().startswith("```"):
            break
        if line.strip().startswith("#"):
            continue
        if not line.strip():
            continue
        prose_lines.append(line)
    prose_text = " ".join(prose_lines).lower()
    words = re.findall(r"\\b[a-z]{3,}\\b", prose_text)
    common_words = {
        "the",
        "and",
        "for",
        "with",
        "from",
        "that",
        "this",
        "are",
        "was",
        "were",
        "has",
        "have",
        "had",
        "but",
        "not",
        "all",
        "any",
        "can",
        "will",
        "may",
        "should",
        "must",
        "each",
        "such",
        "when",
        "where",
        "what",
        "which",
        "who",
        "why",
        "how",
    }
    keywords = [w for w in words if w not in common_words]
    return keywords[:10]


def _is_definition_line(line: str, definition: DetectedDefinition) -> bool:
    line_stripped = line.strip()
    if not line_stripped.startswith(("type ", "func ", "const ", "var ", "interface ")):
        return False
    return definition.name in line or definition.raw_name in line


def _find_definition_line(block_lines: List[str], definition: DetectedDefinition) -> Optional[int]:
    for i, line in enumerate(block_lines):
        if _is_definition_line(line, definition):
            return i
    for i, line in enumerate(block_lines):
        if line.strip().startswith("//"):
            continue
        if definition.name in line or definition.raw_name in line:
            return i
    return None


def _collect_comment_lines(block_lines: List[str], definition_line_idx: int) -> List[str]:
    comments: List[str] = []
    for i in range(definition_line_idx - 1, -1, -1):
        line = block_lines[i].strip()
        if not line:
            continue
        if line.startswith("//"):
            comment_text = line[2:].strip()
            if comment_text:
                comments.insert(0, comment_text)
            continue
        if "/*" in line and "*/" in line:
            comment_match = re.search(r"/\\*(.*?)\\*/", line)
            if comment_match:
                comments.insert(0, comment_match.group(1).strip())
            continue
        if line.startswith("/*"):
            comment_lines = [line[2:].strip()]
            for j in range(i + 1, definition_line_idx):
                next_line = block_lines[j].strip()
                if "*/" in next_line:
                    comment_lines.append(next_line.replace("*/", "").strip())
                    break
                comment_lines.append(next_line)
            comments.insert(0, " ".join(comment_lines))
            break
        break
    return comments


def _clean_comment_text(comment_text: str) -> str:
    return_type_pattern = (
        r"\\breturns?\\s+(?:\\*?[A-Z][a-zA-Z0-9]*|error|bool|string|"
        r"int\\d*|uint\\d*|float\\d*)(?:\\s+if|\\s+when|\\s+on)?[\\s,.]*"
    )
    comment_text = re.sub(return_type_pattern, "", comment_text, flags=re.IGNORECASE)
    comment_text = re.sub(r"\\([^)]*\\*?[A-Z][a-zA-Z0-9]+[^)]*\\)", "", comment_text)
    comment_text = re.sub(r"\\s+", " ", comment_text)
    return comment_text.strip()


def extract_code_comments(definition: DetectedDefinition) -> str:
    """
    Extract comment text from code block above the definition.

    Looks for comments immediately before the definition line:
    - Single-line comments (//) on lines before the definition
    - Multi-line comments (/* */) ending before the definition
    - Strips function signatures (parameters and return types) from definition line

    Returns:
        Combined comment text, or empty string if none found
    """
    if not definition.code_block_content:
        return ""
    block_lines = definition.code_block_content.split("\\n")
    definition_line_idx = _find_definition_line(block_lines, definition)
    if definition_line_idx is None:
        return ""
    comments = _collect_comment_lines(block_lines, definition_line_idx)
    comment_text = " ".join(comments)
    return _clean_comment_text(comment_text)


def extract_keywords_from_comments(definition: DetectedDefinition) -> List[str]:
    """
    Extract domain-specific keywords from definition comments.

    Uses extract_code_comments() to get comment text, then normalizes and extracts
    domain-specific terms. Allows partial word matches for better coverage.
    """
    comment_text = extract_code_comments(definition)
    if not comment_text:
        return []
    comment_lower = comment_text.lower()
    keywords: List[str] = []
    seen: Set[str] = set()
    for phrase in PRIORITY_PHRASES:
        phrase_pattern_exact = r"\\b" + re.escape(phrase) + r"\\b"
        phrase_no_spaces = phrase.replace(" ", "")
        if re.search(phrase_pattern_exact, comment_lower) or phrase_no_spaces in comment_lower:
            normalized_phrase = phrase.replace(" ", "")
            if normalized_phrase not in seen:
                keywords.append(normalized_phrase)
                seen.add(normalized_phrase)
    words = re.findall(r"\\b[a-z]{3,}\\b", comment_lower)
    for word in words:
        normalized = normalize_keyword(word)
        if normalized and normalized not in seen:
            if normalized in KEYWORD_TO_SECTION_MAPPING:
                keywords.append(normalized)
                seen.add(normalized)
    for keyword in KEYWORD_TO_SECTION_MAPPING:
        if keyword in seen:
            continue
        keyword_no_spaces = keyword.replace(" ", "")
        if keyword_no_spaces in seen:
            continue
        if keyword_no_spaces in comment_lower:
            keywords.append(keyword_no_spaces)
            seen.add(keyword_no_spaces)
        elif " " in keyword:
            keyword_words = keyword.split()
            pattern_parts = []
            for i, word in enumerate(keyword_words):
                if i > 0:
                    pattern_parts.append(r"\\s+\\w{0,20}?\\s+")
                pattern_parts.append(r"\\b" + re.escape(word) + r"\\b")
            phrase_pattern = "".join(pattern_parts)
            if re.search(phrase_pattern, comment_lower):
                keywords.append(keyword_no_spaces)
                seen.add(keyword_no_spaces)
    return keywords


def get_section_kind(section_name: str) -> Optional[str]:
    """
    Determine what kind of definition a section accepts.
    """
    section_lower = section_name.lower()
    if "methods" in section_lower and "helper" not in section_lower:
        return "method"
    if "helper" in section_lower and "function" in section_lower:
        return "func"
    if (
        ("type" in section_lower or "types" in section_lower)
        and "method" not in section_lower
        and "function" not in section_lower
    ):
        return "type"
    if "interface" in section_lower and "method" not in section_lower:
        return "type"
    if "error" in section_lower and "type" in section_lower:
        return "type"
    if "generic" in section_lower and "type" in section_lower:
        return "type"
    return None


def extract_section_level(section: str, kind: str) -> str:
    """
    Extract appropriate section level based on definition kind.
    """
    if " > " in section:
        parts = section.split(" > ", 1)
        h2_part = parts[0].strip()
        h3_part = parts[1].strip() if len(parts) > 1 else ""
        section_kind = get_section_kind(section)
        if kind in ("type", "struct", "interface"):
            return h2_part
        if kind in ("method", "func"):
            if h3_part:
                return h3_part
            return h2_part
        if section_kind == "type" and kind in ("type", "struct", "interface"):
            return h2_part
        if section_kind == "method" and kind == "method":
            return h3_part if h3_part else h2_part
        if section_kind == "func" and kind == "func":
            return h3_part if h3_part else h2_part
        if kind in ("type", "struct", "interface"):
            return h2_part
        return h3_part if h3_part else h2_part
    return section


def _match_section_pattern(section_lower: str, pattern: str) -> Tuple[bool, bool]:
    pattern_lower = pattern.lower()
    if pattern_lower in section_lower:
        return True, False
    pattern_words = [w for w in pattern_lower.split() if len(w) >= 4]
    if pattern_words and all(word in section_lower for word in pattern_words):
        return True, True
    return False, False


def _apply_keyword_pattern_score(
    section_lower: str,
    section_pattern: str,
    weight_type: str,
    weight_map: Dict[str, float],
    base_weight: float,
) -> Tuple[float, bool]:
    matched, partial = _match_section_pattern(section_lower, section_pattern)
    if not matched:
        return 0.0, False
    weight = weight_map.get(weight_type, base_weight)
    if partial:
        weight *= 0.8
    return weight, partial


def _score_priority_phrases(
    keywords: List[str],
    section_lower: str,
    keyword_mapping: Dict[str, List[Tuple[str, str]]],
    weight_map: Dict[str, float],
    matched_words: Set[str],
) -> Tuple[float, List[str]]:
    total_score = 0.0
    reasoning: List[str] = []
    for phrase in PRIORITY_PHRASES:
        normalized_phrase = phrase.replace(" ", "")
        if phrase not in keywords and normalized_phrase not in keywords:
            continue
        if phrase not in keyword_mapping:
            continue
        for section_pattern, weight_type in keyword_mapping[phrase]:
            weight, partial = _apply_keyword_pattern_score(
                section_lower,
                section_pattern,
                weight_type,
                weight_map,
                0.15,
            )
            if weight <= 0.0:
                continue
            total_score += weight
            if partial:
                reasoning.append(
                    f"Priority phrase '{phrase}' partially matches "
                    f"section: +{int(weight * 100)}%"
                )
            else:
                reasoning.append(
                    f"Priority phrase '{phrase}' matches section: +{int(weight * 100)}%"
                )
            matched_words.update(word.lower() for word in phrase.split())
            break
    return total_score, reasoning


def _score_keywords(
    keywords: List[str],
    section_lower: str,
    keyword_mapping: Dict[str, List[Tuple[str, str]]],
    weight_map: Dict[str, float],
    matched_words: Set[str],
) -> Tuple[float, List[str]]:
    total_score = 0.0
    reasoning: List[str] = []
    normalized_priority = {p.replace(" ", "") for p in PRIORITY_PHRASES}
    for keyword in keywords:
        if keyword in matched_words:
            continue
        if keyword in PRIORITY_PHRASES or keyword.replace(" ", "") in normalized_priority:
            continue
        if keyword not in keyword_mapping:
            continue
        for section_pattern, weight_type in keyword_mapping[keyword]:
            weight, partial = _apply_keyword_pattern_score(
                section_lower,
                section_pattern,
                weight_type,
                weight_map,
                0.05,
            )
            if weight <= 0.0:
                continue
            total_score += weight
            if partial:
                reasoning.append(
                    f"Keyword '{keyword}' partially matches section: +{int(weight * 100)}%"
                )
            else:
                reasoning.append(
                    f"Keyword '{keyword}' matches section: +{int(weight * 100)}%"
                )
            break
    return total_score, reasoning


def match_keywords_to_section(
    keywords: List[str],
    section_level: str,
    keyword_mapping: Dict[str, List[Tuple[str, str]]],
) -> Tuple[float, List[str]]:
    """
    Match extracted keywords against section level using keyword mapping.
    """
    if not keywords or not section_level:
        return (0.0, [])
    section_lower = section_level.lower()
    total_score = 0.0
    reasoning: List[str] = []
    matched_words: Set[str] = set()
    weight_map = {"strong": 0.15, "medium": 0.10, "weak": 0.05}
    phrase_score, phrase_reasoning = _score_priority_phrases(
        keywords,
        section_lower,
        keyword_mapping,
        weight_map,
        matched_words,
    )
    total_score += phrase_score
    reasoning.extend(phrase_reasoning)
    keyword_score, keyword_reasoning = _score_keywords(
        keywords,
        section_lower,
        keyword_mapping,
        weight_map,
        matched_words,
    )
    total_score += keyword_score
    reasoning.extend(keyword_reasoning)
    if total_score > 0.35:
        total_score = 0.35
        if len(reasoning) > 1:
            reasoning.append("(capped at 35% total)")
    return (total_score, reasoning)
