"""
Go code block and signature utilities for markdown.

Re-exports from _base and _rest for backward compatibility.
"""

from lib.go_markdown._base import (
    EXAMPLE_MARKERS,
    EXAMPLE_NAME_PREFIXES,
    Signature,
    determine_type_kind,
    extract_go_doc_comment_above,
    extract_receiver_type,
    find_go_code_blocks,
    is_example_code,
    is_example_signature_name,
    is_in_go_code_block,
    is_public_name,
    normalize_generic_name,
    parse_go_def_signature,
    remove_go_comments,
)
from lib.go_markdown._rest import (
    InterfaceParser,
    check_kind_word_after,
    count_go_definitions,
    extract_interfaces_from_go_file,
    extract_interfaces_from_markdown,
    find_definition_line_index,
    find_first_definition,
    is_continuation_line,
    is_definition_start_line,
    is_example_definition,
    is_signature_only_code_block,
    normalize_go_signature,
    normalize_go_signature_with_params,
)

__all__ = [
    "EXAMPLE_MARKERS",
    "EXAMPLE_NAME_PREFIXES",
    "Signature",
    "InterfaceParser",
    "check_kind_word_after",
    "count_go_definitions",
    "determine_type_kind",
    "extract_go_doc_comment_above",
    "extract_interfaces_from_go_file",
    "extract_interfaces_from_markdown",
    "extract_receiver_type",
    "find_go_code_blocks",
    "find_first_definition",
    "find_definition_line_index",
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
