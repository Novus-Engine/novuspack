#!/usr/bin/env python3
"""
Shared utilities for parsing and processing Go code blocks in markdown files.

This module provides common functions for:
- Detecting and extracting Go code blocks from markdown
- Parsing Go function, method, and type signatures
- Normalizing Go signatures and type names
- Detecting example code (single lines and entire code blocks)

Facade: re-exports from lib.go_markdown for backward compatibility.
Implementation lives in scripts/lib/go_markdown/ (_base, _rest).
"""

from lib.go_markdown import (
    EXAMPLE_MARKERS,
    EXAMPLE_NAME_PREFIXES,
    InterfaceParser,
    Signature,
    check_kind_word_after,
    count_go_definitions,
    determine_type_kind,
    extract_go_doc_comment_above,
    extract_interfaces_from_go_file,
    extract_interfaces_from_markdown,
    extract_receiver_type,
    find_definition_line_index,
    find_first_definition,
    find_go_code_blocks,
    is_continuation_line,
    is_definition_start_line,
    is_example_code,
    is_example_definition,
    is_example_signature_name,
    is_in_go_code_block,
    is_public_name,
    is_signature_only_code_block,
    normalize_generic_name,
    normalize_go_signature,
    normalize_go_signature_with_params,
    parse_go_def_signature,
    remove_go_comments,
)

__all__ = [
    "EXAMPLE_MARKERS",
    "EXAMPLE_NAME_PREFIXES",
    "InterfaceParser",
    "Signature",
    "check_kind_word_after",
    "count_go_definitions",
    "determine_type_kind",
    "extract_go_doc_comment_above",
    "extract_interfaces_from_go_file",
    "extract_interfaces_from_markdown",
    "extract_receiver_type",
    "find_definition_line_index",
    "find_first_definition",
    "find_go_code_blocks",
    "is_continuation_line",
    "is_definition_start_line",
    "is_example_code",
    "is_example_definition",
    "is_example_signature_name",
    "is_in_go_code_block",
    "is_public_name",
    "is_signature_only_code_block",
    "normalize_generic_name",
    "normalize_go_signature",
    "normalize_go_signature_with_params",
    "parse_go_def_signature",
    "remove_go_comments",
]
