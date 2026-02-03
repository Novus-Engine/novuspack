"""Helpers for validate_links.py: heading extraction, text normalization, anchor suggestion."""

import re
from pathlib import Path
from typing import List

from lib._validation_utils import (
    extract_headings_with_anchors,
    validate_anchor,
)

# Compiled regex patterns for performance (module level)
_RE_MARKDOWN_FORMAT = re.compile(r'[*_`]')
_RE_SPECIAL_CHARS = re.compile(r'[^a-zA-Z0-9 .-]')
_RE_SPLIT_WORDS = re.compile(r'[\s.\-]+')
_RE_CAMEL_CASE = re.compile(r'([a-z])([A-Z])')
_RE_NUMBERING_PREFIX = re.compile(r'^([0-9]+(?:\.[0-9]+)*)\.?\s+(.+)$')


def extract_headings(file_path, file_cache=None):
    """
    Extract all headings from a markdown file and generate anchors.

    Args:
        file_path: Path to the file
        file_cache: Optional FileContentCache instance to use for reading files

    Returns:
        dict: Mapping of anchor -> (heading_text, heading_level, line_number)
    """
    return extract_headings_with_anchors(Path(file_path), file_cache=file_cache)


def normalize_text_for_matching(text: str) -> str:
    """
    Normalize text for matching by removing markdown formatting,
    converting to lowercase, and removing common words.

    Args:
        text: Text to normalize

    Returns:
        Normalized text string
    """
    text = _RE_MARKDOWN_FORMAT.sub('', text)
    text = text.lower()
    text = _RE_SPECIAL_CHARS.sub('', text)
    return text.strip()


def extract_words(text: str) -> List[str]:
    """
    Extract words from text, handling various separators.

    Handles:
    - Spaces: "Add File" -> ["add", "file"]
    - Dots: "Package.AddFile" -> ["package", "add", "file"]
    - Hyphens: "add-file" -> ["add", "file"]
    - CamelCase: "AddFile" -> ["add", "file"]
    - Mixed: "Package.AddFile" -> ["package", "add", "file"]

    Args:
        text: Text to extract words from

    Returns:
        List of normalized words (all lowercase)
    """
    text = _RE_MARKDOWN_FORMAT.sub('', text)
    words = _RE_SPLIT_WORDS.split(text)
    all_words = []
    for word in words:
        if not word:
            continue
        camel_split = _RE_CAMEL_CASE.sub(r'\1 \2', word)
        camel_words = camel_split.split()
        all_words.extend([w.lower() for w in camel_words if w])
    stop_words = {'the', 'a', 'an', 'and', 'or', 'of', 'in', 'on', 'at', 'to', 'for', 'with', 'by'}
    return [w for w in all_words if w and w not in stop_words]


def strip_numbering_prefix(text: str) -> str:
    """
    Strip numbering prefix from heading text (e.g., "1.2.3 Add File" -> "Add File").

    Args:
        text: Heading text that may contain numbering

    Returns:
        Text with numbering prefix removed
    """
    match = _RE_NUMBERING_PREFIX.match(text)
    if match:
        return match.group(2).strip()
    return text


def calculate_word_match_score(link_words: List[str], heading_words: List[str]) -> float:
    """
    Calculate word matching score between link text and heading text.

    Args:
        link_words: List of words from link text
        heading_words: List of words from heading text

    Returns:
        Score from 0-100 based on word matching
    """
    if not link_words or not heading_words:
        return 0.0

    link_set = set(link_words)
    heading_set = set(heading_words)

    if link_set == heading_set:
        return 100.0

    if link_set.issubset(heading_set):
        return 90.0

    matching_words = link_set.intersection(heading_set)
    if matching_words:
        match_ratio = len(matching_words) / len(link_set)
        return 60.0 + (match_ratio * 30.0)

    partial_matches = 0
    for link_word in link_words:
        for heading_word in heading_words:
            if link_word in heading_word or heading_word in link_word:
                partial_matches += 1
                break

    if partial_matches > 0:
        partial_ratio = partial_matches / len(link_words)
        return 20.0 + (partial_ratio * 40.0)

    return 0.0


def suggest_anchor(link_text, broken_anchor, target_file, heading_cache, verbose=False):
    """
    Suggest correct anchor based on weighted heuristics.

    Args:
        link_text: Text from the markdown link
        broken_anchor: The broken anchor that was not found
        target_file: Path to the target file
        heading_cache: Dictionary mapping file paths to heading dictionaries
                      (anchor -> (heading_text, heading_level, line_number))
        verbose: If True, return detailed scoring information

    Returns:
        Tuple of (suggested_anchor, confidence_score) or None if no good match found.
        If verbose=True, returns (suggested_anchor, confidence_score, score_details)
    """
    headings_dict = heading_cache.get(str(target_file), {})
    if not headings_dict:
        return None

    normalized_link_text = normalize_text_for_matching(link_text)
    link_words = extract_words(link_text)
    broken_anchor_words = extract_words(broken_anchor.replace('-', ' '))

    best_match = None
    best_score = 0.0
    best_details = {}

    for anchor, (heading_text, heading_level, line_num) in headings_dict.items():
        if not validate_anchor(anchor):
            continue

        heading_text_no_numbering = strip_numbering_prefix(heading_text)
        normalized_heading = normalize_text_for_matching(heading_text_no_numbering)
        heading_words = extract_words(heading_text_no_numbering)

        scores = {}

        word_score = calculate_word_match_score(link_words, heading_words)
        scores['word_match'] = word_score
        weighted_word = word_score * 0.4

        anchor_score = 0.0
        if anchor == broken_anchor:
            anchor_score = 100.0
        elif broken_anchor in anchor:
            anchor_score = 70.0
        elif anchor in broken_anchor:
            anchor_score = 50.0
        else:
            anchor_words = extract_words(anchor.replace('-', ' '))
            anchor_word_score = calculate_word_match_score(broken_anchor_words, anchor_words)
            anchor_score = anchor_word_score * 0.5
        scores['anchor_similarity'] = anchor_score
        weighted_anchor = anchor_score * 0.3

        context_score = 0.0
        if heading_level == 2:
            context_score += 30.0
        elif heading_level == 3:
            context_score += 20.0
        elif heading_level == 4:
            context_score += 10.0
        if line_num < 100:
            context_score += 5.0
        scores['context'] = context_score
        weighted_context = context_score * 0.2

        norm_score = 0.0
        if normalized_link_text == normalized_heading:
            norm_score = 100.0
        elif (normalized_link_text in normalized_heading or
              normalized_heading in normalized_link_text):
            norm_score = 60.0
        scores['normalization'] = norm_score
        weighted_norm = norm_score * 0.1

        total_score = weighted_word + weighted_anchor + weighted_context + weighted_norm

        if total_score > best_score:
            best_score = total_score
            best_match = anchor
            best_details = {
                'heading_text': heading_text,
                'heading_level': heading_level,
                'line_num': line_num,
                'scores': scores,
                'total_score': total_score
            }

    if best_match and best_score >= 70.0:
        if verbose:
            return (best_match, best_score, best_details)
        return (best_match, best_score)

    return None
