#!/usr/bin/env python3
"""
Shared utilities for parsing and processing Go code blocks in markdown files.

This module provides common functions for:
- Detecting and extracting Go code blocks from markdown
- Parsing Go function, method, and type signatures
- Normalizing Go signatures and type names
- Detecting example code (single lines and entire code blocks)
"""

import re
from dataclasses import dataclass
from typing import List, Optional, Tuple

from lib._validation_utils import find_heading_for_code_block


# Example detection markers
EXAMPLE_MARKERS = [
    'hypothetical', 'not the actual', 'this is not', 'not a real',
    'example only', 'example type', 'example interface', 'example struct',
    'example version', 'example pattern', 'illustration only',
    "not an actual", "shown for illustration"
]

EXAMPLE_NAME_PREFIXES = ('Example', 'Hypothetical', 'Mock', 'Test')

# Compiled regex patterns for performance
_RE_TYPE_NAME = re.compile(r'^\s*type\s+(\w+)')
_RE_FUNC_NAME = re.compile(r'^\s*func\s+(?:\([^)]+\)\s+)?(\w+)')
_RE_GO_COMMENT_SINGLE = re.compile(r'//.*$')
_RE_GO_COMMENT_SINGLE_MULTILINE = re.compile(r'//.*$', re.MULTILINE)
_RE_GO_COMMENT_MULTI = re.compile(r'/\*.*?\*/', flags=re.DOTALL)
_RE_GO_DOC_LINE = re.compile(r'^\s*//\s?(.*)$')
_RE_GO_BLOCK_COMMENT_START = re.compile(r'^\s*/\*\s?(.*)$')
_RE_GO_BLOCK_COMMENT_END = re.compile(r'^(.*)\*/\s*$')
_RE_INTERFACE_PATTERN = re.compile(r'^\s*(?:type\s+)?(\w+)(?:\s*\[[^\]]+\])?\s+interface\s*\{')
_RE_STRUCT_PATTERN = re.compile(r'^\s*type\s+(\w+)(?:\s*\[[^\]]+\])?\s+struct\s*\{')
_RE_ALIAS_PATTERN = re.compile(r'^\s*type\s+(\w+)(?:\s*\[[^\]]+\])?\s*=\s')
_RE_POINTER_PATTERN = re.compile(r'^\s*type\s+(\w+)(?:\s*\[[^\]]+\])?\s+\*')
_RE_SLICE_PATTERN = re.compile(r'^\s*type\s+(\w+)(?:\s*\[[^\]]+\])?\s+\[\]')
_RE_MAP_PATTERN = re.compile(r'^\s*type\s+(\w+)(?:\s*\[[^\]]+\])?\s+map\s*\[')
_RE_TYPE_PATTERN = re.compile(r'^\s*type\s+(\w+)(?:\s*\[[^\]]+\])?\s+\S')
_RE_AFTER_TYPE_PATTERN = re.compile(r'^\s*type\s+\w+(?:\s*\[[^\]]+\])?\s+(.+)')
_RE_REMOVE_BRACE = re.compile(r'\s*\{.*$')
_RE_METHOD_PATTERN = re.compile(
    r'^\s*func\s+(\([^)]+\))\s+(\w+)(?:\s*\[[^\]]+\])?\s*\(([^)]*)\)\s*(.*)$'
)
_RE_FUNC_PATTERN = re.compile(r'^\s*func\s+(\w+)(?:\s*\[[^\]]+\])?\s*\(([^)]*)\)\s*(.*)$')
_RE_RECEIVER_TYPE = re.compile(r'^\s*(?:\w+\s+)?(?:\*)?\s*(\w+(?:\[[^\]]+\])?)')
_RE_WHITESPACE = re.compile(r'\s+')
_RE_GENERIC_PARAMS = re.compile(r'\[[^\]]+\]')
_RE_PACKAGE_TYPE = re.compile(r'\b([a-z][a-z0-9_]*(?:\.[a-z][a-z0-9_]*)*)\.([A-Z][A-Za-z0-9_]*)\b')
_RE_METHOD_NORMALIZE = re.compile(r'func\s+(\([^)]+\))\s+(\w+)\s*\(([^)]*)\)\s*(.*)$')
_RE_FUNC_NORMALIZE = re.compile(r'func\s+(\w+)\s*\(([^)]*)\)\s*(.*)$')
_RE_FUNC_WITH_PARAMS = re.compile(r'func\s+(?:\([^)]+\)\s+)?(\w+)\s*\(([^)]*)\)\s*(.*)$')
_RE_RECEIVER_MATCH = re.compile(r'func\s+(\([^)]+\))\s+')
_RE_WHITESPACE_NORMALIZE = re.compile(r'\s+')
_RE_GENERICS_TAG = re.compile(r'\bgenerics\.(Tag|TagValueType|PathEntry)\b')
_RE_METADATA_TYPES = re.compile(
    r'\bmetadata\.(PackageMetadata|PackageInfo|FileEntry|PathMetadataEntry|ProcessingState)\b'
)
_RE_FILEFORMAT_TYPES = re.compile(r'\bfileformat\.(PackageHeader|FileIndex|IndexEntry)\b')
_RE_HEADER_TYPE = re.compile(r'\bHeader\b')
_RE_PKGERRORS_TYPES = re.compile(r'\bpkgerrors\.(ErrorType|PackageError)\b')
_RE_SIGNATURES_TYPES = re.compile(r'\bsignatures\.(Signature|SignatureInfo)\b')
_RE_FUNC_TYPE_DEF = re.compile(r'^\s*type\s+\w+(?:\s*\[[^\]]+\])?\s+func\s*\(')


def find_go_code_blocks(content: str) -> List[Tuple[int, int, str]]:
    """
    Find all Go code blocks in markdown content.

    Args:
        content: Markdown content as string

    Returns:
        List of tuples: (start_line, end_line, code_content)
        Lines are 1-indexed.
    """
    go_blocks = []
    lines = content.split('\n')

    i = 0
    while i < len(lines):
        line = lines[i]

        # Check for Go code block start
        if line.strip() == '```go':
            start_line = i + 1  # 1-indexed for reporting
            code_lines = []
            i += 1

            # Collect code until closing ```
            while i < len(lines) and lines[i].strip() != '```':
                code_lines.append(lines[i])
                i += 1

            if i < len(lines):  # Found closing ```
                code_content = '\n'.join(code_lines)
                go_blocks.append((start_line, i + 1, code_content))

        i += 1

    return go_blocks


def is_in_go_code_block(content: str, line_num: int) -> bool:
    """
    Check if a given line number is inside a Go code block.

    Args:
        content: Markdown content as string
        line_num: Line number to check (1-indexed)

    Returns:
        True if the line is inside a ```go code block
    """
    lines = content.split('\n')
    in_go_block = False

    for _idx, line in enumerate(lines[:line_num], 1):
        if line.strip() == '```go':
            in_go_block = True
        elif line.strip() == '```' and in_go_block:
            in_go_block = False

    return in_go_block


def _heading_suggests_example(heading_text: Optional[str]) -> bool:
    """Return True if heading text suggests example code."""
    if not heading_text:
        return False
    heading_lower = heading_text.lower()
    if any(marker in heading_lower for marker in EXAMPLE_MARKERS):
        return True
    return "example" in heading_lower


def _prose_line_suggests_example(prose_lower: str) -> bool:
    """Return True if a prose line suggests example code."""
    if any(marker in prose_lower for marker in EXAMPLE_MARKERS):
        return True
    if "example" not in prose_lower:
        return False
    example_phrases = [
        "this is an example", "example:", "example code", "example only",
        "example type", "example interface"
    ]
    if any(phrase in prose_lower for phrase in example_phrases):
        return True
    example_markers = ["example code", "example type", "example interface"]
    if "for example" in prose_lower and not any(m in prose_lower for m in example_markers):
        return False
    return False


def _prose_before_block_suggests_example(
    lines: List[str],
    start_line: int,
    check_prose_before_block: bool,
) -> bool:
    """Return True if prose immediately before the code block suggests example."""
    if not check_prose_before_block or start_line <= 1:
        return False
    prose_start = max(0, start_line - 11)
    prose_end = start_line - 1
    for j in range(prose_start, prose_end):
        if j < 0 or j >= len(lines):
            continue
        line = lines[j]
        stripped = line.strip()
        if not stripped or stripped in ('```', '```go') or stripped.startswith('#'):
            continue
        if _prose_line_suggests_example(line.lower()):
            return True
    return False


def _comment_line_suggests_example(prev_lower: str) -> bool:
    """Return True if a code/comment line suggests example."""
    if any(marker in prev_lower for marker in EXAMPLE_MARKERS):
        return True
    if "example" not in prev_lower:
        return False
    comment_phrases = ["// example", "// example:", "example code", "example only"]
    if any(phrase in prev_lower for phrase in comment_phrases):
        return True
    return "for example" not in prev_lower


def _name_suggests_example(stripped: str) -> bool:
    """Return True if type or function name suggests example."""
    type_match = _RE_TYPE_NAME.match(stripped)
    if type_match and type_match.group(1).startswith(EXAMPLE_NAME_PREFIXES):
        return True
    func_match = _RE_FUNC_NAME.match(stripped)
    return bool(func_match and func_match.group(1).startswith(EXAMPLE_NAME_PREFIXES))


def _code_line_suggests_example(
    lines: List[str],
    line_index: int,
    start_line: int,
) -> bool:
    """Return True if this code line or its preceding lines suggest example."""
    for j in range(max(start_line - 1, line_index - 5), line_index):
        if j < 0 or j >= len(lines):
            continue
        if _comment_line_suggests_example(lines[j].lower()):
            return True
    line = lines[line_index] if line_index < len(lines) else ''
    return _name_suggests_example(line.strip())


def _get_lines_to_check(
    code: str,
    start_line: int,
    lines: List[str],
    check_single_line: Optional[int],
    max_lines_to_check: int,
) -> Optional[List[int]]:
    """Return list of line indices to check for example markers, or None if invalid."""
    if check_single_line is not None:
        line_index = start_line - 1 + check_single_line
        if line_index < 0 or line_index >= len(lines):
            return None
        return [line_index]
    code_lines = code.split('\n')
    if not code_lines:
        return None
    result = []
    for i, code_line in enumerate(code_lines[:max_lines_to_check]):
        if code_line.strip() and not code_line.strip().startswith('```'):
            line_idx = start_line - 1 + i
            if line_idx < len(lines):
                result.append(line_idx)
    return result


def is_example_code(
    code: str,
    start_line: int,
    *,
    content: Optional[str] = None,
    lines: Optional[List[str]] = None,
    heading_text: Optional[str] = None,
    auto_find_heading: bool = False,
    check_prose_before_block: bool = True,
    check_single_line: Optional[int] = None,
    max_lines_to_check: int = 5,
) -> bool:
    """
    Check if Go code is example code.

    This unified function can check:
    - A single line within a code block
    - Multiple lines in a code block (default: first 5 lines)
    - An entire code block

    Looks for example markers in:
    - The heading above the code block (if provided or auto-found)
    - Prose text immediately before the code block (if check_prose_before_block is True)
    - Previous lines within the code block
    - The name of the type/function definition

    Args:
        code: The code block content (without ```go markers) OR a single line
        start_line: Line number where the code block starts (1-indexed)
        content: Full markdown content (preferred - used for heading finding and prose checking)
        lines: All lines of the file as a list
            (alternative to content; content will be derived if needed)
        heading_text: Optional heading text above the code block
        auto_find_heading: If True, automatically find heading from content
        check_prose_before_block: If True, check prose text between heading and code block
        check_single_line: If provided (0-indexed line number within code), check only that line
        max_lines_to_check: Maximum number of lines to check in code block (default: 5)

    Returns:
        True if the code appears to be example code
    """
    if content is None and lines is not None:
        content = '\n'.join(lines)
    if lines is None:
        lines = content.split('\n') if content else [''] * (start_line - 1) + code.split('\n')
    if auto_find_heading and content:
        heading_text = find_heading_for_code_block(content, start_line)
    if not lines:
        return False
    if _heading_suggests_example(heading_text):
        return True
    if _prose_before_block_suggests_example(lines, start_line, check_prose_before_block):
        return True
    lines_to_check = _get_lines_to_check(
        code, start_line, lines, check_single_line, max_lines_to_check
    )
    if not lines_to_check:
        return False
    for line_index in lines_to_check:
        if line_index >= start_line - 1 and _code_line_suggests_example(
            lines, line_index, start_line
        ):
            return True
    return False


def is_example_signature_name(name: str) -> bool:
    """
    Check if a signature name indicates it's an example.

    Args:
        name: Signature name to check

    Returns:
        True if the name starts with example prefixes
    """
    return name.startswith(EXAMPLE_NAME_PREFIXES)


def remove_go_comments(text: str, multiline: bool = False) -> str:
    """
    Remove Go comments from text (single or multi-line).

    Args:
        text: Go code text
        multiline: If True, handles multi-line strings and block comments.
                   If False, also strips whitespace (for single-line usage).

    Returns:
        Text with comments removed (and stripped if multiline=False)
    """
    if multiline:
        text = _RE_GO_COMMENT_SINGLE_MULTILINE.sub('', text)
        text = _RE_GO_COMMENT_MULTI.sub('', text)
    else:
        text = _RE_GO_COMMENT_SINGLE.sub('', text)
        # For single-line usage, strip whitespace (matches original behavior)
        text = text.strip()
    return text


def _collect_block_comment_lines(code_lines, i, end_match, should_skip):
    """Collect multi-line block comment parts upward; return (block_text, new_i)."""
    block_parts: List[str] = []
    end_text = end_match.group(1).strip()
    if end_text and not should_skip(end_text):
        block_parts.insert(0, end_text)
    i -= 1
    while i >= 0:
        raw2 = code_lines[i].rstrip("\n")
        stripped2 = raw2.strip()
        start_match = _RE_GO_BLOCK_COMMENT_START.match(stripped2)
        if start_match:
            start_text = start_match.group(1).strip()
            if start_text and start_text != "*/":
                cleaned = start_text.replace("*/", "").strip()
                if cleaned and not should_skip(cleaned):
                    block_parts.insert(0, cleaned)
            break
        if stripped2:
            if stripped2.startswith("*"):
                stripped2 = stripped2[1:].strip()
            cleaned = stripped2.replace("*/", "").strip()
            if cleaned and not should_skip(cleaned):
                block_parts.insert(0, cleaned)
        i -= 1
    block_text = " ".join([p for p in block_parts if p]).strip()
    return (block_text, i)


def extract_go_doc_comment_above(
    code_lines: List[str],
    definition_line_index: int,
) -> str:
    """
    Extract doc comment text immediately above a definition line.

    This is an additive helper intended for future scoring improvements.

    Args:
        code_lines: List of Go code lines (no markdown fences).
        definition_line_index: 0-based index of the definition line within code_lines.

    Returns:
        Normalized doc comment text, or empty string if none.
    """
    if not code_lines:
        return ""
    if definition_line_index <= 0 or definition_line_index > len(code_lines) - 1:
        return ""

    # Walk upward collecting contiguous comment lines / blocks.
    collected: List[str] = []
    i = definition_line_index - 1

    def _should_skip_doc_line(text: str) -> bool:
        # Skip TODO/FIXME lines, but keep other doc comment content.
        t = (text or "").strip()
        if not t:
            return True
        upper = t.upper()
        return upper.startswith("TODO:") or upper.startswith("FIXME:")

    while i >= 0:
        raw = code_lines[i].rstrip("\n")
        stripped = raw.strip()

        if not stripped:
            # Allow blank lines between comment lines but stop if we already started
            # collecting and then hit a blank line (doc comments must be adjacent).
            if collected:
                break
            i -= 1
            continue

        # Single-line doc comment: // ...
        m = _RE_GO_DOC_LINE.match(raw)
        if m:
            text = m.group(1).strip()
            if text and not _should_skip_doc_line(text):
                collected.insert(0, text)
            i -= 1
            continue

        # Inline block comment: /* ... */
        if "/*" in stripped and "*/" in stripped:
            inner = _RE_GO_COMMENT_MULTI.sub(lambda mm: mm.group(0)[2:-2], stripped)
            inner = inner.strip()
            if inner and not _should_skip_doc_line(inner):
                collected.insert(0, inner)
            i -= 1
            continue

        end_match = _RE_GO_BLOCK_COMMENT_END.match(stripped)
        if end_match and "/*" not in stripped:
            block_text, i = _collect_block_comment_lines(
                code_lines, i, end_match, _should_skip_doc_line
            )
            if block_text:
                collected.insert(0, block_text)
            i -= 1
            continue

        # Not a comment line; stop.
        break

    # Normalize whitespace.
    out = " ".join(collected).strip()
    out = _RE_WHITESPACE_NORMALIZE.sub(" ", out)
    return out


def determine_type_kind(line: str) -> Optional[str]:
    """
    Determine the kind of a Go type definition from a line.

    This function extracts the kind ('interface', 'struct', 'alias', or 'type') from a Go type
    definition line. It checks interfaces first, then structs, then type aliases, then other types.

    Args:
        line: Line of Go code

    Returns:
        'interface', 'struct', 'alias', 'pointer', 'slice', 'map', 'type',
        or None if not a type definition

    Examples:
        - "type Package interface {" -> 'interface'
        - "type FileEntry struct {" -> 'struct'
        - "type ProcessingState uint8" -> 'type'
        - "type Option[T] struct {" -> 'struct'
        - "type Name = SomeType" -> 'alias' (type alias)
        - "type Name[T] = SomeType[T]" -> 'alias' (generic type alias)
        - "type Name *SomeType" -> 'pointer' (pointer type)
        - "type Name []SomeType" -> 'slice' (slice type)
        - "type Name map[K]V" -> 'map' (map type)
        - "type Name SomeType" -> 'type' (regular type definition)
    """
    line_clean = remove_go_comments(line)

    # Check for interface definitions FIRST (before type definitions)
    # This ensures interfaces are correctly classified, not as types
    # Pattern: type Name interface { or Name interface {
    interface_match = _RE_INTERFACE_PATTERN.match(line_clean)
    if interface_match:
        return 'interface'

    # Check for struct definitions (distinct from other types)
    # Pattern: type Name struct { or type Name[T] struct {
    struct_match = _RE_STRUCT_PATTERN.match(line_clean)
    if struct_match:
        return 'struct'

    # Check for type aliases: type Name = Type or type Name[T] = Type
    # This must be checked before other type definitions
    alias_match = _RE_ALIAS_PATTERN.match(line_clean)
    if alias_match:
        return 'alias'

    # Check for pointer types: type Name *SomeType or type Name[T] *SomeType
    pointer_match = _RE_POINTER_PATTERN.match(line_clean)
    if pointer_match:
        return 'pointer'

    # Check for slice types: type Name []SomeType or type Name[T] []SomeType
    slice_match = _RE_SLICE_PATTERN.match(line_clean)
    if slice_match:
        return 'slice'

    # Check for map types: type Name map[K]V or type Name[T] map[K]V
    map_match = _RE_MAP_PATTERN.match(line_clean)
    if map_match:
        return 'map'

    # Check for other type definitions
    # (custom types, etc. - excludes structs, interfaces, aliases, pointers, slices, maps)
    # Pattern: type Name SomeType or type Name[T] SomeType
    # This handles regular type definitions (may or may not have generics)
    type_match = _RE_TYPE_PATTERN.match(line_clean)
    if type_match:
        # Make sure it's not already matched by struct/interface/alias/pointer/slice/map patterns
        # and it's not a function type
        # Check that it doesn't start with pointer, slice, or map patterns
        if ('struct' not in line_clean and 'interface' not in line_clean
                and '=' not in line_clean and 'func(' not in line_clean):
            # Check if it's a pointer, slice, or map (already handled above)
            # by checking if the pattern after type name matches those
            after_type_match = _RE_AFTER_TYPE_PATTERN.match(line_clean)
            if after_type_match:
                after_type = after_type_match.group(1).strip()
                # If it doesn't start with *, [], or map[, it's a regular type
                if (not after_type.startswith('*')
                        and not after_type.startswith('[]')
                        and not after_type.startswith('map[')):
                    return 'type'

    return None


def parse_go_def_signature(line: str, location: str = "") -> Optional[Signature]:
    """
    Parse a Go definition signature from a line (function, method, or type).

    Args:
        line: Line of Go code
        location: Optional location string (file path and line number)

    Returns:
        Signature object or None if no definition found
        - For functions/methods: kind='func' or 'method', includes params and returns
        - For types: kind='type', 'interface', 'struct', etc., includes generic_params
    """
    line_clean = remove_go_comments(line)

    # Try to parse as function/method first
    # Remove opening brace if present
    line_no_brace = _RE_REMOVE_BRACE.sub('', line_clean).strip()

    # Method: func (r *Receiver) Name(params) returns
    method_match = _RE_METHOD_PATTERN.match(line_no_brace)
    if method_match:
        receiver_str = method_match.group(1)
        name = method_match.group(2)
        params = method_match.group(3)
        returns = method_match.group(4).strip()
        receiver_type = extract_receiver_type(receiver_str)
        return Signature(
            name=name,
            kind='method',
            receiver=receiver_type,
            params=params,
            returns=returns,
            location=location,
            is_public=is_public_name(name)
        )

    # Function: func Name(params) returns
    func_match = _RE_FUNC_PATTERN.match(line_no_brace)
    if func_match:
        name = func_match.group(1)
        params = func_match.group(2)
        returns = func_match.group(3).strip()
        return Signature(
            name=name,
            kind='func',
            params=params,
            returns=returns,
            location=location,
            is_public=is_public_name(name)
        )

    # Try to parse as type definition
    kind = determine_type_kind(line_clean)
    if kind is not None:
        # Special-case: interfaces may be written as:
        # - type Name interface { ... }
        # - Name interface { ... }
        #
        # determine_type_kind() supports both forms, but the generic type match below
        # only matches "type Name ...", so handle interface explicitly.
        if kind == 'interface':
            interface_match = re.match(
                r'^\s*(?:type\s+)?(\w+)(?:\s*(\[[^\]]+\]))?\s+interface\s*\{',
                line_clean,
            )
            if interface_match:
                name = interface_match.group(1)
                generic_params = interface_match.group(2)  # e.g., "[T any]"
                return Signature(
                    name=name,
                    kind='interface',
                    generic_params=generic_params,
                    location=location,
                    is_public=is_public_name(name),
                )

        # Extract name and generic parameters
        type_match = re.match(
            r'^\s*type\s+(\w+)(?:\s*(\[[^\]]+\]))?\s+',
            line_clean
        )
        if type_match:
            name = type_match.group(1)
            generic_params = type_match.group(2)  # e.g., "[T any]"
            return Signature(
                name=name,
                kind=kind,  # 'type', 'interface', 'struct', 'alias', etc.
                generic_params=generic_params,
                location=location,
                is_public=is_public_name(name)
            )

    return None


def extract_receiver_type(receiver_str: str, normalize_generics: bool = False) -> str:
    """
    Extract the type name from a receiver string.

    Args:
        receiver_str: Receiver string like "(r *Receiver)" or "(o *Option[T])" or just "Package"
        normalize_generics: If True, remove generic parameters from the type name

    Returns:
        Type name (e.g., "Receiver" or "Option")
    """
    # If already just a type name (starts with capital, no parentheses), return as-is
    if receiver_str and receiver_str[0].isupper() and '(' not in receiver_str:
        if normalize_generics:
            return normalize_generic_name(receiver_str)
        return receiver_str

    # Remove parentheses if present
    receiver_clean = receiver_str.strip('()').strip()

    # Pattern: variableName *TypeName or variableName TypeName
    # Also handle: *TypeName (no variable name)
    match = _RE_RECEIVER_TYPE.match(receiver_clean)
    if match:
        type_name = match.group(1)
        if normalize_generics:
            return normalize_generic_name(type_name)
        return type_name

    # Fallback: split by spaces and take last part
    parts = receiver_clean.split()
    if len(parts) >= 2:
        # Has variable name, last part is type
        type_name = parts[-1]
        if normalize_generics:
            return normalize_generic_name(type_name)
        return type_name
    if len(parts) == 1:
        # Single word - could be type name or pointer
        single_word = parts[0]
        # Remove leading * if present
        if single_word.startswith('*'):
            type_name = single_word[1:]
        else:
            type_name = single_word
        if normalize_generics:
            return normalize_generic_name(type_name)
        return type_name

    # If single word and starts with capital, it's likely already a type name
    if receiver_clean and receiver_clean[0].isupper():
        # Remove leading * if present
        if receiver_clean.startswith('*'):
            type_name = receiver_clean[1:]
        else:
            type_name = receiver_clean
        if normalize_generics:
            return normalize_generic_name(type_name)
        return type_name

    return receiver_clean


def normalize_generic_name(name: str) -> str:
    """
    Normalize a generic type name by removing generic parameters.

    Args:
        name: Type name that may include generics (e.g., "Option[T]", "BufferPool[T any]")

    Returns:
        Base type name without generics (e.g., "Option", "BufferPool")

    Examples:
        - "Option[T]" -> "Option"
        - "BufferPool[T]" -> "BufferPool"
        - "ConfigBuilder[T]" -> "ConfigBuilder"
        - "Option" -> "Option" (no change if no generics)
        - "Container[Option[T]]" -> "Container" (handles nested generics)
        - "Type[]" -> "Type" (handles empty brackets)
    """
    # Remove generic parameters like [T], [T any], [T, U], [], etc.
    # Handle nested brackets by repeatedly removing the rightmost bracket pair
    # This ensures we remove innermost brackets first
    result = name
    while True:
        # Find the rightmost [ that has a matching ]
        # We'll work backwards to find balanced brackets
        last_open = result.rfind('[')
        if last_open == -1:
            break  # No more brackets

        # Find the matching closing bracket
        bracket_count = 0
        found_close = False
        for i in range(last_open, len(result)):
            if result[i] == '[':
                bracket_count += 1
            elif result[i] == ']':
                bracket_count -= 1
                if not bracket_count:
                    # Found matching bracket, remove this bracket pair
                    result = result[:last_open] + result[i + 1:]
                    found_close = True
                    break

        if not found_close:
            # Unmatched bracket, just remove the [
            result = result[:last_open] + result[last_open + 1:]

    return result


def _normalize_go_signature_preprocessing(
    sig_str: str, use_whitespace_normalize: bool = False
) -> str:
    """Common preprocessing for signature normalization.

    Args:
        sig_str: Go signature string
        use_whitespace_normalize: If True, use _RE_WHITESPACE_NORMALIZE,
            else _RE_WHITESPACE

    Returns:
        Preprocessed signature string
    """
    # Remove comments
    sig_str = remove_go_comments(sig_str, multiline=True)

    # Normalize whitespace
    if use_whitespace_normalize:
        sig_str = _RE_WHITESPACE_NORMALIZE.sub(' ', sig_str)
    else:
        sig_str = _RE_WHITESPACE.sub(' ', sig_str)
    sig_str = sig_str.strip()

    # Remove generic type parameters
    sig_str = _RE_GENERIC_PARAMS.sub('', sig_str)

    return sig_str


def _normalize_package_names_general(sig_str: str) -> str:
    """Normalize package-qualified type names to short names (general approach).

    Pattern: package.Type -> Type
    Only normalizes internal NovusPack packages, not standard library packages.
    """
    standard_lib_packages = {
        'context', 'errors', 'fmt', 'io', 'os', 'strings', 'bytes', 'time',
        'sync', 'reflect', 'encoding', 'encoding/json', 'encoding/binary',
        'crypto', 'net', 'path', 'path/filepath', 'syscall', 'unicode',
        'math', 'sort', 'strconv', 'bufio', 'compress', 'archive', 'hash'
    }

    def replace_package_type(match):
        full_match = match.group(0)
        package_part = match.group(1)
        type_name = match.group(2)

        # Check if this is a standard library package
        base_package = (
            package_part.split('.')[0] if '.' in package_part else package_part
        )
        if base_package in standard_lib_packages:
            return full_match  # Keep standard library types as-is

        # For internal packages, return just the type name
        return type_name

    return _RE_PACKAGE_TYPE.sub(replace_package_type, sig_str)


def _normalize_package_names_specific(sig_str: str) -> str:
    """Normalize package names using specific regex substitutions.

    For sync validation. Handles re-exported types: generics.X -> X,
    metadata.X -> X, etc.
    """
    sig_str = _RE_GENERICS_TAG.sub(r'\1', sig_str)
    sig_str = _RE_METADATA_TYPES.sub(r'\1', sig_str)
    sig_str = _RE_FILEFORMAT_TYPES.sub(r'\1', sig_str)
    sig_str = _RE_HEADER_TYPE.sub('PackageHeader', sig_str)
    sig_str = _RE_PKGERRORS_TYPES.sub(r'\1', sig_str)
    sig_str = _RE_SIGNATURES_TYPES.sub(r'\1', sig_str)
    return sig_str


def _normalize_returns_simple(
    returns: str, normalize_param_list_func
) -> Tuple[str, bool]:
    """Normalize return values - simple approach (removes names).

    Returns:
        Tuple of (normalized_returns, has_multiple_returns)
    """
    normalized_returns = ""
    has_multiple_returns = False
    if returns:
        returns_stripped = returns.strip()
        if returns_stripped.startswith('(') and returns_stripped.endswith(')'):
            returns_content = returns_stripped[1:-1].strip()
            normalized_returns = normalize_param_list_func(returns_content)
            has_multiple_returns = True
        elif ',' in returns_stripped:
            normalized_returns = normalize_param_list_func(returns_stripped)
            has_multiple_returns = True
        else:
            normalized_returns = normalize_param_list_func(returns_stripped)
            has_multiple_returns = False
    return normalized_returns, has_multiple_returns


def _extract_receiver_type_safe(receiver_str: str) -> str:
    """Extract receiver type safely, handling both (Type) and (var *Type) formats."""
    receiver_clean = receiver_str.strip('()').strip()
    # Check if it's already just a type name (single word, starts with capital)
    if (len(receiver_clean.split()) == 1 and receiver_clean and
            receiver_clean[0].isupper()):
        return receiver_clean

    # Has variable name, extract type
    receiver_type = extract_receiver_type(receiver_str, normalize_generics=False)
    # Fallback: if extraction failed, try to get last word
    if not receiver_type or receiver_type == receiver_str.strip('()'):
        parts = receiver_clean.split()
        if len(parts) >= 2:
            receiver_type = parts[-1]  # Last part is the type
        else:
            receiver_type = receiver_clean
    return receiver_type


def _format_normalized_signature(
    name: str,
    normalized_params: str,
    normalized_returns: str,
    *,
    receiver_type: Optional[str] = None,
    has_multiple_returns: bool = False,
    always_paren_returns: bool = False,
) -> str:
    """
    Format a normalized Go signature string from its components.

    Args:
        name: Function/method name
        normalized_params: Normalized parameter list (types only)
        normalized_returns: Normalized return values
        receiver_type: Receiver type (if method, None for function)
        has_multiple_returns: Whether there are multiple return values
            (for simple normalization)
        always_paren_returns: If True, always use parentheses for returns
            (for param-preserving normalization)

    Returns:
        Formatted signature string
    """
    # Format return values
    if normalized_returns:
        # Determine if we should use parentheses
        use_parens = always_paren_returns or has_multiple_returns
        if use_parens:
            returns_str = f"({normalized_returns})"
        else:
            returns_str = normalized_returns
    else:
        returns_str = ""

    # Format the signature
    if receiver_type:
        # Method with receiver
        if returns_str:
            return f"func ({receiver_type}) {name}({normalized_params}) {returns_str}"
        return f"func ({receiver_type}) {name}({normalized_params})"
    # Function without receiver
    if returns_str:
        return f"func {name}({normalized_params}) {returns_str}"
    return f"func {name}({normalized_params})"


@dataclass(frozen=True)
class Signature:
    """
    Represents a Go function, method, or type signature.

    This is a shared dataclass used across multiple validation scripts.
    Optional fields allow different scripts to track additional information
    as needed.
    """
    name: str
    kind: str  # 'func', 'method', 'type', 'interface'
    receiver: Optional[str] = None  # For methods: the receiver type
    params: str = ""  # Parameter list as string
    returns: str = ""  # Return types as string
    location: str = ""  # File path and line number
    is_public: bool = True  # Whether it's exported (starts with capital)
    # Optional fields for scripts that need more detail
    has_body: bool = False  # Whether this is a full definition with body
    method_count: int = 0  # For interfaces: number of methods in body
    field_count: int = 0  # For structs: number of fields in body
    generic_params: Optional[str] = None  # Generic parameters like "[T any]" or None

    def normalized_key(self) -> str:
        """Generate a normalized key for comparison."""
        if self.kind == 'method' and self.receiver:
            return f"{self.receiver}.{self.name}"
        if self.kind in ('type', 'interface') and self.generic_params:
            # Include generics in key to distinguish SigningKey from SigningKey[T]
            return f"{self.name}{self.generic_params}"
        return self.name

    def normalized_signature(self) -> str:
        """Generate a normalized signature string for comparison."""
        # Normalize whitespace and remove comments
        params = _RE_WHITESPACE_NORMALIZE.sub(' ', self.params.strip())
        returns = _RE_WHITESPACE_NORMALIZE.sub(' ', self.returns.strip())

        if self.kind == 'method':
            return f"func ({self.receiver}) {self.name}({params}) ({returns})"
        if self.kind == 'func':
            return f"func {self.name}({params}) ({returns})"
        if self.kind == 'type':
            return f"type {self.name}"
        if self.kind == 'interface':
            return f"type {self.name} interface"
        return f"{self.kind} {self.name}"

    def normalized_type_name(self) -> str:
        """
        Get normalized type name (without generics for display purposes).

        For types with generics, returns just the base name.
        For other types, returns the name as-is.
        """
        if self.generic_params:
            return normalize_generic_name(self.name)
        return self.name

    def is_method(self) -> bool:
        """Check if this is a method (has receiver)."""
        return self.kind == 'method' and self.receiver is not None


def is_public_name(name: str) -> bool:
    """
    Check if a name is public (exported) in Go.

    In Go, exported identifiers start with an uppercase letter.

    Args:
        name: The name to check

    Returns:
        True if the name is public (starts with uppercase letter)
    """
    return bool(name and name[0].isupper())
