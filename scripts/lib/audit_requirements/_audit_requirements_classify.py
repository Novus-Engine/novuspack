"""
Section content and heading classification for requirements coverage audit.
"""

import re
from typing import List, Tuple, Optional

from lib._validation_utils import count_sentences, has_code_blocks, contains_url

from lib.audit_requirements._audit_requirements_config import (
    MIN_MEANINGFUL_SENTENCES,
    MIN_PROSE_LINES,
    ARCHITECTURAL_SCORE_THRESHOLD,
    FUNCTION_SIGNATURE_WEIGHT,
    TYPE_DEFINITION_WEIGHT,
    EXAMPLE_CODE_WEIGHT,
    FUNCTIONAL_KEYWORD_PATTERNS,
    ARCHITECTURAL_KEYWORD_PATTERNS,
    NON_ARCHITECTURAL_PHRASE_PATTERNS,
    IMPLEMENTATION_PATTERN,
    LINK_ONLY_PATTERN,
    LABEL_ONLY_PATTERN,
    URL_ONLY_PATTERN,
)


def extract_section_content(
    content: str,
    heading_line: int,
    all_headings: List[Tuple[int, int, str]],
    lines: Optional[List[str]] = None
) -> str:
    """
    Extract section content from heading to next heading (any level).
    """
    if lines is None:
        lines = content.split('\n')

    next_heading_line = None
    for line_num, _level, _text in all_headings:
        if line_num > heading_line:
            next_heading_line = line_num
            break

    if next_heading_line:
        section_lines = lines[heading_line - 1:next_heading_line - 1]
    else:
        section_lines = lines[heading_line - 1:]

    return '\n'.join(section_lines)


def calculate_architectural_score(text: str) -> int:
    """Calculate architectural confidence score for given text."""
    score = 0
    sorted_arch_patterns = sorted(
        ARCHITECTURAL_KEYWORD_PATTERNS, key=lambda x: x[1], reverse=True
    )
    for pattern, weight in sorted_arch_patterns:
        if pattern.search(text):
            score += weight
    for pattern, weight in NON_ARCHITECTURAL_PHRASE_PATTERNS:
        if pattern.search(text):
            score += weight
    return score


def _strip_list_marker(line: str) -> str:
    """Remove leading list markers from a line."""
    return re.sub(r'^\s*(?:[-*+]|\d+\.)\s+', '', line)


def _is_link_only_line(line: str) -> bool:
    """Check if a line is only a markdown link or URL."""
    stripped = line.strip()
    if not stripped:
        return False
    if LINK_ONLY_PATTERN.match(stripped):
        return True
    if URL_ONLY_PATTERN.match(stripped):
        return True
    if contains_url(stripped):
        cleaned = re.sub(r'(https?://\S+|www\.\S+|mailto:\S+)', '', stripped).strip()
        return not cleaned
    return False


def _is_label_only_line(line: str) -> bool:
    """Check if a line is only a label with no prose."""
    stripped = line.strip()
    if not stripped:
        return False
    if LABEL_ONLY_PATTERN.match(stripped) and not contains_url(stripped):
        return True
    return False


def _extract_prose_lines(section_content: str) -> List[str]:
    """Extract prose lines from a section, filtering link-only/label-only lines."""
    prose_lines = []
    in_code_block = False
    for line in section_content.split('\n'):
        stripped = line.strip()
        if stripped.startswith('```'):
            in_code_block = not in_code_block
            continue
        if in_code_block:
            continue
        if re.match(r'^#{1,6}\s+', stripped):
            continue
        if not stripped:
            continue
        cleaned = _strip_list_marker(line).strip()
        if _is_link_only_line(cleaned) or _is_label_only_line(cleaned):
            continue
        prose_lines.append(cleaned)
    return prose_lines


def _get_prose_metrics(section_content: str) -> dict:
    """Calculate prose metrics for a section."""
    prose_lines = _extract_prose_lines(section_content)
    prose_text = '\n'.join(prose_lines).strip()
    sentence_count = count_sentences(prose_text)
    prose_line_count = len([line for line in prose_lines if line.strip()])
    has_meaningful_prose = (
        sentence_count >= MIN_MEANINGFUL_SENTENCES and
        prose_line_count >= MIN_PROSE_LINES
    )
    return {
        'prose_text': prose_text,
        'sentence_count': sentence_count,
        'prose_line_count': prose_line_count,
        'has_meaningful_prose': has_meaningful_prose,
    }


def _is_signature_only_section(section_content: str, go_code_utils_module) -> bool:
    """Check if a section contains only signature-only Go blocks."""
    go_blocks = go_code_utils_module.find_go_code_blocks(section_content)
    if not go_blocks:
        return False
    if has_code_blocks(section_content, exclude_languages={'go'}):
        return False
    section_lines = section_content.split('\n')
    for start_line, _end_line, code_content in go_blocks:
        if go_code_utils_module.is_example_code(
            code_content,
            start_line,
            content=section_content,
            lines=section_lines,
            auto_find_heading=False,
            check_prose_before_block=True
        ):
            return False
        if not go_code_utils_module.is_signature_only_code_block(code_content):
            return False
    return True


def _is_example_only_section(section_content: str, go_code_utils_module) -> bool:
    """Check if a section only contains example Go code blocks."""
    go_blocks = go_code_utils_module.find_go_code_blocks(section_content)
    if not go_blocks:
        return False
    if has_code_blocks(section_content, exclude_languages={'go'}):
        return False
    section_lines = section_content.split('\n')
    for start_line, _end_line, code_content in go_blocks:
        if not go_code_utils_module.is_example_code(
            code_content,
            start_line,
            content=section_content,
            lines=section_lines,
            auto_find_heading=False,
            check_prose_before_block=True
        ):
            return False
    return True


def analyze_section_content(
    section_content: str,
    _functional_keywords: List[str],
    _architectural_keywords: List[str],
    *,
    func_weight: int,
    type_weight: int,
    example_weight: int,
    go_code_utils_module
) -> dict:
    """Analyze section content to determine classification with scoring."""
    result = {
        'has_function_signatures': False,
        'has_type_definitions': False,
        'has_example_code': False,
        'has_functional_keywords': False,
        'has_architectural_keywords': False,
        'has_implementation_keyword': False,
        'architectural_score': 0,
        'function_signature_count': 0,
        'type_definition_count': 0,
        'example_code_count': 0,
        'content_score': 0,
        'content_type': 'mixed',
        'sentence_count': 0,
        'prose_line_count': 0,
        'has_meaningful_prose': False,
    }

    go_blocks = go_code_utils_module.find_go_code_blocks(section_content)
    lines = section_content.split('\n')

    for start_line, _end_line, code_content in go_blocks:
        counts = go_code_utils_module.count_go_definitions(
            code_content,
            filter_example=True,
            lines=lines,
            start_line=start_line
        )
        result['function_signature_count'] += counts['func'] + counts['method']
        result['type_definition_count'] += counts['type']
        if counts['func'] + counts['method'] > 0:
            result['has_function_signatures'] = True
        if counts['type'] > 0:
            result['has_type_definitions'] = True

        code_lines = code_content.split('\n')
        for line_idx, _code_line in enumerate(code_lines):
            is_example = go_code_utils_module.is_example_code(
                code_content, start_line,
                lines=lines,
                check_single_line=line_idx
            )
            if is_example:
                result['example_code_count'] += 1
                result['has_example_code'] = True

    result['content_score'] = (
        result['function_signature_count'] * func_weight +
        result['type_definition_count'] * type_weight +
        result['example_code_count'] * example_weight
    )

    prose_metrics = _get_prose_metrics(section_content)
    prose_text = prose_metrics['prose_text']
    result['sentence_count'] = prose_metrics['sentence_count']
    result['prose_line_count'] = prose_metrics['prose_line_count']
    result['has_meaningful_prose'] = prose_metrics['has_meaningful_prose']

    for pattern in FUNCTIONAL_KEYWORD_PATTERNS:
        if pattern.search(prose_text):
            result['has_functional_keywords'] = True
            break

    prose_architectural_score = calculate_architectural_score(prose_text)
    result['architectural_score'] = prose_architectural_score
    result['has_architectural_keywords'] = (
        prose_architectural_score >= ARCHITECTURAL_SCORE_THRESHOLD
    )

    if IMPLEMENTATION_PATTERN.search(prose_text):
        result['has_implementation_keyword'] = True

    if (result['content_score'] > 0 or result['has_function_signatures'] or
            result['has_type_definitions']):
        result['content_type'] = 'functional'
    elif result['has_functional_keywords']:
        result['content_type'] = 'functional'
    elif result['has_architectural_keywords']:
        result['content_type'] = 'architectural'
    elif result['has_implementation_keyword']:
        result['content_type'] = 'implementation'
    elif result['has_example_code'] and result['content_score'] < 0:
        result['content_type'] = 'documentation'
    else:
        result['content_type'] = 'mixed'

    return result


def _check_exclusion_patterns(heading_text: str, exclusion_patterns: List[str]):
    """Check if heading matches any exclusion patterns."""
    heading_text_lower = heading_text.lower()
    for pattern in exclusion_patterns:
        if pattern.lower() in heading_text_lower:
            return {
                'classification': 'excluded',
                'needs_requirement': False,
                'severity_if_missing': 'none',
                'reason': 'matches exclusion pattern',
            }
    return None


def _analyze_heading_keywords(heading_text: str) -> dict:
    """Analyze heading text for functional and implementation keywords."""
    heading_has_functional = False
    for pattern in FUNCTIONAL_KEYWORD_PATTERNS:
        if pattern.search(heading_text):
            heading_has_functional = True
            break

    heading_has_implementation = bool(IMPLEMENTATION_PATTERN.search(heading_text))
    heading_architectural_score = calculate_architectural_score(heading_text)

    return {
        'has_functional': heading_has_functional,
        'has_implementation': heading_has_implementation,
        'architectural_score': heading_architectural_score,
    }


def _classify_by_definitions(analysis: dict):
    """Classify heading based on function/type definitions."""
    has_definitions = (
        analysis['has_function_signatures'] or
        analysis['has_type_definitions']
    )
    if (has_definitions or analysis['content_score'] > 0 or
            (not analysis['content_score'] and has_definitions)):
        reason = (
            'contains function signatures' if analysis['has_function_signatures']
            else 'contains type definitions'
        )
        return {
            'classification': 'functional',
            'needs_requirement': True,
            'severity_if_missing': 'error',
            'reason': reason,
        }
    return None


def _heading_looks_like_type_definition(heading_text: str) -> bool:
    """Return True if heading title suggests a single type/struct/interface definition."""
    text_lower = heading_text.lower()
    without_number = re.sub(r'^\d+(?:\.\d+)*\s+', '', text_lower).strip()
    if without_number.endswith(' struct'):
        return True
    if without_number.endswith(' interface'):
        return True
    return False


def _classify_signature_only(
    section_content: str, analysis: dict, go_code_utils_module, heading_text: str = ''
):
    """Classify sections that are signature-only or lack meaningful prose."""
    if not analysis['has_meaningful_prose']:
        if _is_example_only_section(section_content, go_code_utils_module):
            return {
                'classification': 'example_only',
                'needs_requirement': False,
                'severity_if_missing': 'none',
                'reason': 'example-only section',
            }
        if _is_signature_only_section(section_content, go_code_utils_module):
            return {
                'classification': 'signature_only',
                'needs_requirement': False,
                'severity_if_missing': 'none',
                'reason': 'signature-only section',
            }
        if _heading_looks_like_type_definition(heading_text):
            go_blocks = go_code_utils_module.find_go_code_blocks(section_content)
            if go_blocks and not has_code_blocks(section_content, exclude_languages={'go'}):
                return {
                    'classification': 'signature_only',
                    'needs_requirement': False,
                    'severity_if_missing': 'none',
                    'reason': 'type-definition section (struct/interface)',
                }
        return {
            'classification': 'non_prose',
            'needs_requirement': False,
            'severity_if_missing': 'none',
            'reason': 'insufficient meaningful prose',
        }
    return None


def _classify_by_keywords(
    heading_keywords: dict, analysis: dict, total_architectural_score: float
):
    """Classify heading based on keywords."""
    if heading_keywords['has_functional'] or analysis['has_functional_keywords']:
        return {
            'classification': 'functional',
            'needs_requirement': True,
            'severity_if_missing': 'error',
            'reason': 'has functional keywords',
        }

    if heading_keywords['has_implementation'] or analysis['has_implementation_keyword']:
        return {
            'classification': 'implementation',
            'needs_requirement': False,
            'severity_if_missing': 'none',
            'reason': 'implementation detail',
        }

    is_architectural = total_architectural_score >= ARCHITECTURAL_SCORE_THRESHOLD
    if is_architectural:
        return {
            'classification': 'architectural',
            'needs_requirement': False,
            'severity_if_missing': 'warning',
            'reason': f'has architectural keywords (score: {total_architectural_score})',
        }

    return None


def _classify_by_content_type(analysis: dict):
    """Classify heading based on content type."""
    if ((analysis['has_example_code'] and analysis['content_score'] < 0) or
            analysis['content_type'] == 'documentation'):
        return {
            'classification': 'documentation',
            'needs_requirement': False,
            'severity_if_missing': 'none',
            'reason': 'documentation only',
        }
    return None


def _classify_by_hierarchy(heading_level: int, analysis: dict) -> dict:
    """Classify heading based on hierarchy level."""
    if heading_level == 2:
        return {
            'classification': 'functional',
            'needs_requirement': True,
            'severity_if_missing': 'error',
            'reason': 'H2 heading requires coverage',
        }
    if heading_level >= 4:
        if analysis['has_function_signatures'] or analysis['has_type_definitions']:
            reason = (
                'contains function signatures' if analysis['has_function_signatures']
                else 'contains type definitions'
            )
            return {
                'classification': 'functional',
                'needs_requirement': True,
                'severity_if_missing': 'error',
                'reason': reason,
            }
        return {
            'classification': 'implementation',
            'needs_requirement': False,
            'severity_if_missing': 'none',
            'reason': 'implementation detail',
        }
    return {
        'classification': 'functional',
        'needs_requirement': True,
        'severity_if_missing': 'error',
        'reason': 'H3 heading requires coverage',
    }


def classify_heading(
    heading_text: str,
    heading_level: int,
    section_content: str,
    exclusion_patterns: List[str],
    *,
    _heading_line: int,
    functional_keywords: List[str],
    _max_org_lines: int,
    go_code_utils_module
) -> dict:
    """
    Classify a heading to determine requirement coverage needs.
    """
    excluded = _check_exclusion_patterns(heading_text, exclusion_patterns)
    if excluded:
        return excluded

    analysis = analyze_section_content(
        section_content,
        functional_keywords,
        [],
        func_weight=FUNCTION_SIGNATURE_WEIGHT,
        type_weight=TYPE_DEFINITION_WEIGHT,
        example_weight=EXAMPLE_CODE_WEIGHT,
        go_code_utils_module=go_code_utils_module
    )

    heading_keywords = _analyze_heading_keywords(heading_text)

    total_architectural_score = (
        (heading_keywords['architectural_score'] * 2) +
        analysis.get('architectural_score', 0)
    )

    classification = _classify_signature_only(
        section_content, analysis, go_code_utils_module, heading_text
    )
    if classification:
        return classification

    classification = _classify_by_definitions(analysis)
    if classification:
        return classification

    classification = _classify_by_keywords(
        heading_keywords, analysis, total_architectural_score
    )
    if classification:
        return classification

    classification = _classify_by_content_type(analysis)
    if classification:
        return classification

    return _classify_by_hierarchy(heading_level, analysis)
