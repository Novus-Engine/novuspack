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
import sys
from pathlib import Path
from typing import List, Tuple, Optional, Dict
from dataclasses import dataclass

# Import heading utility from validation_utils to avoid duplication
lib_dir = Path(__file__).parent
scripts_dir = lib_dir.parent
if str(scripts_dir) not in sys.path:
    sys.path.insert(0, str(scripts_dir))

from lib._validation_utils import find_heading_for_code_block  # noqa: E402


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

    for i, line in enumerate(lines[:line_num], 1):
        if line.strip() == '```go':
            in_go_block = True
        elif line.strip() == '```' and in_go_block:
            in_go_block = False

    return in_go_block


def is_example_code(
    code: str,
    start_line: int,
    content: Optional[str] = None,
    lines: Optional[List[str]] = None,
    heading_text: Optional[str] = None,
    auto_find_heading: bool = False,
    check_prose_before_block: bool = True,
    check_single_line: Optional[int] = None,
    max_lines_to_check: int = 5
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
    # Derive content from lines if needed (for heading finding)
    if content is None and lines is not None:
        content = '\n'.join(lines)

    # Derive lines from content if needed (for prose checking)
    if lines is None:
        if content:
            lines = content.split('\n')
        else:
            # If we have code but no content/lines, create minimal lines list
            code_lines = code.split('\n')
            lines = [''] * (start_line - 1) + code_lines

    # Handle auto-finding heading (requires content)
    if auto_find_heading and content:
        heading_text = find_heading_for_code_block(content, start_line)

    if not lines:
        return False

    # Check heading for example indicators (highest priority) - done once per code block
    if heading_text:
        heading_lower = heading_text.lower()
        if any(marker in heading_lower for marker in EXAMPLE_MARKERS):
            return True
        # Also check for standalone "example" word in heading
        if "example" in heading_lower:
            return True

    # Check prose text immediately before the code block (between heading and code start)
    # Done once per code block
    if check_prose_before_block and start_line > 1:
        # Look back up to 10 lines before the code block for example indicators
        prose_start = max(0, start_line - 11)  # 0-indexed
        prose_end = start_line - 1  # 0-indexed, exclusive
        for j in range(prose_start, prose_end):
            if j < 0 or j >= len(lines):
                continue
            prose_line = lines[j]
            # Skip empty lines and markdown code block markers
            if not prose_line.strip() or prose_line.strip() in ('```', '```go'):
                continue
            # Skip markdown headings (they're already checked)
            if prose_line.strip().startswith('#'):
                continue
            prose_lower = prose_line.lower()
            if any(marker in prose_lower for marker in EXAMPLE_MARKERS):
                return True
            # Also check for standalone "example" word in prose,
            # but only if it's clearly an example marker
            # Skip if it's just "for example" or "example:" in normal prose
            if "example" in prose_lower:
                # Only flag if it's clearly an example marker, not just "for example" in prose
                example_phrases = [
                    "this is an example", "example:", "example code", "example only",
                    "example type", "example interface"
                ]
                if any(phrase in prose_lower for phrase in example_phrases):
                    return True
                # Skip common phrases that use "example" but aren't marking example code
                example_markers = ["example code", "example type", "example interface"]
                if "for example" in prose_lower and not any(
                    marker in prose_lower for marker in example_markers
                ):
                    continue

    # Determine what lines to check
    if check_single_line is not None:
        # Check single line within code block
        line_index = start_line - 1 + check_single_line  # Convert to 0-indexed
        if line_index < 0 or line_index >= len(lines):
            return False
        lines_to_check = [line_index]
    else:
        # Check multiple lines (code block)
        code_lines = code.split('\n')
        if not code_lines:
            return False

        # Check the first few lines of the code block for example markers
        lines_to_check = []
        for i, code_line in enumerate(code_lines[:max_lines_to_check]):
            if code_line.strip() and not code_line.strip().startswith('```'):
                line_idx = start_line - 1 + i  # Convert to 0-indexed
                if line_idx < len(lines):
                    lines_to_check.append(line_idx)

    # Check each line for example indicators
    for line_index in lines_to_check:
        if line_index < start_line - 1:
            continue

        # Check previous lines in this code block for example indicators
        for j in range(max(start_line - 1, line_index - 5), line_index):
            if j < 0 or j >= len(lines):
                continue
            prev_line = lines[j]
            prev_lower = prev_line.lower()
            if any(marker in prev_lower for marker in EXAMPLE_MARKERS):
                return True
            # Also check for standalone "example" word in code comments, but be more specific
            if "example" in prev_lower:
                # Only flag if it's clearly an example marker in comments
                comment_example_phrases = [
                    "// example", "// example:", "example code", "example only"
                ]
                if any(phrase in prev_lower for phrase in comment_example_phrases):
                    return True
                # Skip if it's just "for example" in a comment
                if "for example" in prev_lower:
                    continue

        # Check if type/function name indicates it's an example
        line = lines[line_index] if line_index < len(lines) else ''
        stripped = line.strip()

        # Check for type definitions
        type_name_match = _RE_TYPE_NAME.match(stripped)
        if type_name_match:
            type_name = type_name_match.group(1)
            if type_name.startswith(EXAMPLE_NAME_PREFIXES):
                return True

        # Check for function/method definitions
        func_name_match = _RE_FUNC_NAME.match(stripped)
        if func_name_match:
            func_name = func_name_match.group(1)
            if func_name.startswith(EXAMPLE_NAME_PREFIXES):
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

        # Multi-line block comment ending on this line: ... */
        end_match = _RE_GO_BLOCK_COMMENT_END.match(stripped)
        if end_match and "/*" not in stripped:
            block_parts: List[str] = []
            end_text = end_match.group(1).strip()
            if end_text and not _should_skip_doc_line(end_text):
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
                        if cleaned and not _should_skip_doc_line(cleaned):
                            block_parts.insert(0, cleaned)
                    break

                if stripped2:
                    # Strip leading "*" for common block comment style.
                    if stripped2.startswith("*"):
                        stripped2 = stripped2[1:].strip()
                    cleaned = stripped2.replace("*/", "").strip()
                    if cleaned and not _should_skip_doc_line(cleaned):
                        block_parts.insert(0, cleaned)
                i -= 1

            block_text = " ".join([p for p in block_parts if p]).strip()
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
    elif len(parts) == 1:
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
                if bracket_count == 0:
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
    receiver_type: Optional[str] = None,
    has_multiple_returns: bool = False,
    always_paren_returns: bool = False
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
        else:
            return f"func ({receiver_type}) {name}({normalized_params})"
    else:
        # Function without receiver
        if returns_str:
            return f"func {name}({normalized_params}) {returns_str}"
        else:
            return f"func {name}({normalized_params})"


def normalize_go_signature(sig_str: str) -> str:
    """
    Normalize a Go signature string for comparison.

    Removes comments, normalizes whitespace, standardizes package names,
    and extracts receiver types properly.

    Args:
        sig_str: Go signature string

    Returns:
        Normalized signature string
    """
    # Common preprocessing
    sig_str = _normalize_go_signature_preprocessing(sig_str, use_whitespace_normalize=False)

    # Normalize package names (general approach)
    sig_str = _normalize_package_names_general(sig_str)

    # Normalize parameter list (remove parameter names, keep types)
    def normalize_param_list(param_str: str) -> str:
        if not param_str.strip():
            return ""
        # Simple normalization: remove parameter names
        # Pattern: name Type -> Type
        params = []
        for param in param_str.split(','):
            param = param.strip()
            parts = param.split()
            if len(parts) >= 2:
                # Has name and type: keep type part
                params.append(' '.join(parts[1:]))
            else:
                params.append(param)
        return ", ".join(params)

    # Extract and normalize function signatures
    # Handle method with receiver - receiver can be in format (Type) or (var *Type)
    method_match = _RE_METHOD_NORMALIZE.match(sig_str)
    if method_match:
        receiver_str = method_match.group(1)
        name = method_match.group(2)
        params = method_match.group(3)
        returns = method_match.group(4).strip()

        # Extract receiver type
        receiver_type = _extract_receiver_type_safe(receiver_str)

        normalized_params = normalize_param_list(params)

        # Normalize return values
        normalized_returns, has_multiple_returns = _normalize_returns_simple(
            returns, normalize_param_list
        )

        # Format signature using shared helper
        return _format_normalized_signature(
            name=name,
            normalized_params=normalized_params,
            normalized_returns=normalized_returns,
            receiver_type=receiver_type,
            has_multiple_returns=has_multiple_returns,
            always_paren_returns=False
        )

    # Check for function without receiver
    func_match = _RE_FUNC_NORMALIZE.match(sig_str)
    if func_match:
        name = func_match.group(1)
        params = func_match.group(2)
        returns = func_match.group(3).strip()

        normalized_params = normalize_param_list(params)

        # Normalize return values
        normalized_returns, has_multiple_returns = _normalize_returns_simple(
            returns, normalize_param_list
        )

        # Format signature using shared helper
        return _format_normalized_signature(
            name=name,
            normalized_params=normalized_params,
            normalized_returns=normalized_returns,
            receiver_type=None,
            has_multiple_returns=has_multiple_returns,
            always_paren_returns=False
        )

    return sig_str


def normalize_go_signature_with_params(sig_str: str) -> str:
    """
    Normalize a Go signature string for comparison while preserving parameter names.

    This is a specialized version for sync validation that handles shorthand
    notation and keeps parameter names for exact matching. Use this when you need
    to compare signatures where parameter names must match exactly.

    The general-purpose `normalize_go_signature()` removes parameter names and
    is better suited for general signature normalization.

    Normalizes:
    - Extra whitespace
    - Comments
    - Generic type parameters (for comparison purposes)
    - Package name differences (generics.X vs X)

    Keeps:
    - Parameter names (must match exactly)
    - Return value names (must match exactly)
    """
    # Common preprocessing
    sig_str = _normalize_go_signature_preprocessing(sig_str, use_whitespace_normalize=True)

    # Normalize package names (specific approach for sync validation)
    sig_str = _normalize_package_names_specific(sig_str)

    # Remove parameter names, keep only types
    # Pattern: name Type -> Type
    # Handle: ctx context.Context, path string -> context.Context, string
    # Handle: offset, size int64 -> int64, int64

    def _is_parameter_name_only(param: str) -> bool:
        """Check if parameter looks like just a name (no type indicators)."""
        return not any(c in param for c in [' ', '.', '*', '[', ']', '(', ')'])

    def _can_split_normalized_param(normalized: str) -> bool:
        """Check if normalized parameter can be safely split by comma."""
        return (',' in normalized and
                not any(c in normalized for c in ['*[', '[]', 'map[']))

    def _process_param_token(param: str, normalize_single_param_func) -> List[Tuple[str, str]]:
        """Process a single parameter token, returning list of (tag, value) tuples."""
        if not param:
            return []

        if _is_parameter_name_only(param):
            # Just a name, might be part of shorthand - keep it for later processing
            return [('name', param)]

        normalized = normalize_single_param_func(param)
        if _can_split_normalized_param(normalized):
            return [('type', p.strip()) for p in normalized.split(',')]
        else:
            return [('type', normalized)]

    def _process_last_param_with_shorthand(
        param: str,
        params: List[Tuple[str, str]],
        normalize_single_param_func
    ) -> None:
        """Process the last parameter, handling shorthand notation."""
        if not param:
            return

        # Check if previous params ended with names (shorthand pattern)
        if params and params[-1][0] == 'name':
            # This is the type for the shorthand names
            type_part = normalize_single_param_func(param)
            # Replace all trailing 'name' entries with this type
            i = len(params) - 1
            while i >= 0 and params[i][0] == 'name':
                params[i] = ('type', type_part)
                i -= 1
        else:
            tokens = _process_param_token(param, normalize_single_param_func)
            params.extend(tokens)

    def _resolve_remaining_names(
        params: List[Tuple[str, str]]
    ) -> List[str]:
        """Resolve any remaining name entries to types, handling edge cases."""
        final_params = []
        i = 0
        while i < len(params):
            if params[i][0] == 'name':
                # Collect consecutive names
                names = [params[i][1]]
                i += 1
                while i < len(params) and params[i][0] == 'name':
                    names.append(params[i][1])
                    i += 1
                # If next is a type, use it; otherwise these are invalid
                if i < len(params) and params[i][0] == 'type':
                    type_part = params[i][1]
                    final_params.extend([type_part] * len(names))
                    i += 1
                else:
                    # Invalid - just use the names as-is (shouldn't happen)
                    final_params.extend(names)
            else:
                final_params.append(params[i][1])
                i += 1
        return final_params

    def normalize_param_list(param_str: str) -> str:
        if not param_str.strip():
            return ""
        # Split parameters by comma, but be careful with nested structures
        params = []
        current = ""
        paren_depth = 0
        bracket_depth = 0

        for char in param_str:
            if char == '(':
                paren_depth += 1
            elif char == ')':
                paren_depth -= 1
            elif char == '[':
                bracket_depth += 1
            elif char == ']':
                bracket_depth -= 1
            elif char == ',' and paren_depth == 0 and bracket_depth == 0:
                # Found a top-level comma separator
                param = current.strip()
                if param:
                    tokens = _process_param_token(param, normalize_single_param)
                    params.extend(tokens)
                current = ""
                continue
            current += char

        # Process last param
        if current.strip():
            _process_last_param_with_shorthand(
                current.strip(), params, normalize_single_param
            )

        # Handle any remaining name entries (shouldn't happen in valid Go, but handle gracefully)
        final_params = _resolve_remaining_names(params)
        return ", ".join(final_params)

    def _is_type_like(param: str) -> bool:
        """Check if parameter looks like a type (starts with type indicators)."""
        return param and (param.startswith('*') or param.startswith('[') or param[0].isupper())

    def _extract_type_from_shorthand(parts: List[str]) -> str:
        """Extract type from shorthand notation (e.g., 'offset, size int64')."""
        type_part = parts[-1]
        first_part = parts[0]
        name_list = [n.strip() for n in first_part.split(',')]
        # Return type repeated for each name
        return ", ".join([type_part] * len(name_list))

    def _extract_type_from_regular(parts: List[str]) -> str:
        """Extract type from regular notation (e.g., 'name Type' or 'name *package.Type')."""
        if len(parts) == 2:
            return parts[-1]
        else:
            # Multiple words: might be name *package.Type
            # Remove first word (the name)
            return " ".join(parts[1:])

    def normalize_single_param(param: str) -> str:
        """Normalize a single parameter, handling shorthand notation."""
        # Remove leading parameter names
        # Pattern: name Type or name1, name2 Type
        # Handle: offset, size int64 -> int64 (expand to int64, int64)

        parts = param.split()
        if len(parts) < 2:
            # Single identifier - might be just a type or just a name
            if _is_type_like(param):
                return param
            # Otherwise, it's probably just a name - return as-is (caller will handle)
            return param

        # Check if first part has commas (shorthand)
        first_part = parts[0]
        if ',' in first_part:
            # Shorthand: offset, size int64
            return _extract_type_from_shorthand(parts)
        else:
            # Regular: name Type - remove the name, keep the type
            return _extract_type_from_regular(parts)

    # Extract and normalize function signatures
    # Pattern: func Name(params) returns or func (r Receiver) Name(params) returns
    func_match = _RE_FUNC_WITH_PARAMS.match(sig_str)
    if func_match:
        name = func_match.group(1)
        params = func_match.group(2)
        returns = func_match.group(3).strip()

        normalized_params = normalize_param_list(params)
        # For returns, keep names and types - they must match exactly
        # Expand shorthand in returns too
        normalized_returns = normalize_param_list(returns) if returns else ""

        # Reconstruct using shared helper
        receiver_match = _RE_RECEIVER_MATCH.match(sig_str)
        if receiver_match:
            receiver = receiver_match.group(1)
            receiver_type = extract_receiver_type(receiver)
            return _format_normalized_signature(
                name=name,
                normalized_params=normalized_params,
                normalized_returns=normalized_returns,
                receiver_type=receiver_type,
                has_multiple_returns=False,  # Not used when always_paren_returns=True
                always_paren_returns=True
            )
        else:
            # For functions without receiver, always use parentheses for returns
            # (this matches the behavior expected by sync validation)
            return _format_normalized_signature(
                name=name,
                normalized_params=normalized_params,
                normalized_returns=normalized_returns,
                receiver_type=None,
                has_multiple_returns=False,  # Not used when always_paren_returns=True
                always_paren_returns=True
            )

    return sig_str


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
        elif self.kind in ('type', 'interface') and self.generic_params:
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
        elif self.kind == 'func':
            return f"func {self.name}({params}) ({returns})"
        elif self.kind == 'type':
            return f"type {self.name}"
        elif self.kind == 'interface':
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


class InterfaceParser:
    """
    Helper class for parsing Go interfaces with brace depth tracking.

    This handles the common pattern of tracking interface definitions
    and their methods across multiple scripts.
    """

    def __init__(self):
        self.in_interface = False
        self.current_interface: Optional[str] = None
        self.brace_depth = 0

    def reset(self):
        """Reset the parser state."""
        self.in_interface = False
        self.current_interface = None
        self.brace_depth = 0

    def check_interface_start(self, line: str) -> Optional[str]:
        """
        Check if a line starts an interface definition.

        Args:
            line: The line to check

        Returns:
            Interface name if this line starts an interface, None otherwise
        """
        # Pattern: type Name interface { or type Name[T] interface {
        interface_match = re.match(
            r'^\s*type\s+(\w+)(?:\s*(\[[^\]]+\]))?\s+interface\s*\{', line
        )
        if interface_match:
            self.in_interface = True
            self.current_interface = interface_match.group(1)
            stripped = line.strip()
            self.brace_depth = stripped.count('{') - stripped.count('}')
            return self.current_interface
        return None

    def update_brace_depth(self, line: str) -> bool:
        """
        Update brace depth for current interface.

        Args:
            line: The current line

        Returns:
            True if still inside interface, False if interface closed
        """
        if not self.in_interface:
            return False

        stripped = line.strip()
        self.brace_depth += stripped.count('{') - stripped.count('}')

        if self.brace_depth <= 0:
            self.in_interface = False
            self.current_interface = None
            return False

        return True

    def is_in_interface(self) -> bool:
        """Check if currently parsing an interface."""
        return self.in_interface

    def get_current_interface(self) -> Optional[str]:
        """Get the name of the current interface being parsed."""
        return self.current_interface


def extract_interfaces_from_go_file(
    file_path: Path,
    parse_methods: bool = True
) -> List[Signature]:
    """
    Extract all interfaces (and optionally their methods) from a Go source file.

    Args:
        file_path: Path to the Go source file
        parse_methods: If True, also extract interface methods as separate signatures

    Returns:
        List of Signature objects for interfaces (and their methods if parse_methods=True)
    """
    interfaces = []
    methods = []

    try:
        resolved_path = file_path.resolve()
        content = file_path.read_text(encoding='utf-8')
        lines = content.split('\n')

        interface_parser = InterfaceParser()

        for line_num, line in enumerate(lines, 1):
            stripped = line.strip()

            # Skip empty lines and comments
            if not stripped or stripped.startswith('//'):
                continue

            # Check for interface start using InterfaceParser
            interface_name = interface_parser.check_interface_start(line)
            if interface_name:
                is_public = is_public_name(interface_name) if interface_name else False
                interfaces.append(Signature(
                    name=interface_name,
                    kind='interface',
                    location=f"{resolved_path}:{line_num}",
                    is_public=is_public
                ))
                continue

            # Track interface brace depth using InterfaceParser
            if interface_parser.is_in_interface():
                # Check brace depth before updating to catch methods on closing line
                brace_depth_before = interface_parser.brace_depth
                current_interface = interface_parser.get_current_interface()
                still_in_interface = interface_parser.update_brace_depth(line)

                # Check for interface method if we're still in interface or on closing line
                if parse_methods and (
                    (still_in_interface and interface_parser.brace_depth > 0) or
                    (brace_depth_before > 0 and not still_in_interface and '{' not in stripped)
                ):
                    sig = parse_go_def_signature(line, location=f"{resolved_path}:{line_num}")
                    if sig and sig.kind in ('func', 'method'):
                        # Interface methods don't have receivers in the interface definition
                        # but we track them as methods of the interface type
                        methods.append(Signature(
                            name=sig.name,
                            kind='method',
                            receiver=current_interface,
                            params=sig.params,
                            returns=sig.returns,
                            location=f"{resolved_path}:{line_num}",
                            is_public=sig.is_public
                        ))

                if not still_in_interface:
                    # Interface closed, parser already reset
                    pass
                continue

    except Exception as e:
        print(f"Warning: Error reading {file_path}: {e}")

    # Return interfaces first, then methods
    return interfaces + methods


def extract_interfaces_from_markdown(
    content: str,
    file_path: Path,
    start_line: int = 1,
    parse_methods: bool = True,
    skip_examples: bool = True,
    lines: Optional[List[str]] = None
) -> List[Signature]:
    """
    Extract all interfaces (and optionally their methods) from Go code blocks in markdown.

    Args:
        content: Markdown content as string
        file_path: Path to the markdown file (for location strings)
        start_line: Starting line number for the content (default: 1)
        parse_methods: If True, also extract interface methods as separate signatures
        skip_examples: If True, skip interfaces that are marked as examples
        lines: Optional list of all lines in the file (for example detection)

    Returns:
        List of Signature objects for interfaces (and their methods if parse_methods=True)
    """
    interfaces = []
    methods = []

    try:
        resolved_path = file_path.resolve()
        go_blocks = find_go_code_blocks(content)

        for block_start_line, block_end_line, code_content in go_blocks:
            # Process each code block
            block_lines = code_content.split('\n')
            interface_parser = InterfaceParser()

            for i, line in enumerate(block_lines):
                # Calculate actual line number in file (1-indexed)
                line_num = block_start_line + i

                stripped = line.strip()

                # Skip empty lines and comments
                if not stripped or stripped.startswith('//'):
                    continue

                # Check if this is an example signature
                is_example = False
                if skip_examples:
                    is_example = is_example_code(
                        code_content, block_start_line,
                        lines=lines,
                        check_single_line=i
                    )

                # Check for interface start using InterfaceParser
                interface_name = interface_parser.check_interface_start(line)
                if interface_name:
                    # Skip if this is an example
                    if is_example:
                        continue

                    # Extract generic parameters from the line
                    generic_match = re.match(
                        r'^\s*type\s+\w+\s*(\[[^\]]+\])?\s+interface\s*\{', line
                    )
                    generic_params = generic_match.group(1) if generic_match else None
                    is_public = is_public_name(interface_name) if interface_name else False

                    # Check if this is a stub (interface without body or minimal body)
                    has_full_body = interface_parser.brace_depth > 0

                    # Count methods if it's a full interface
                    method_count = 0
                    if has_full_body and parse_methods:
                        # Look ahead to count methods within this code block
                        temp_brace_depth = interface_parser.brace_depth
                        for j in range(i + 1, len(block_lines)):
                            temp_line = block_lines[j]
                            temp_stripped = temp_line.strip()
                            if not temp_stripped or temp_stripped.startswith('//'):
                                continue
                            temp_brace_depth += temp_stripped.count('{') - temp_stripped.count('}')
                            temp_sig = parse_go_def_signature(temp_line, location="")
                            if (temp_sig and temp_sig.kind in ('func', 'method') and
                                    temp_brace_depth > 0):
                                method_count += 1
                            if temp_brace_depth <= 0:
                                break

                    interfaces.append(Signature(
                        name=interface_name,
                        kind='interface',
                        location=f"{resolved_path}:{line_num}",
                        is_public=is_public,
                        has_body=has_full_body,
                        method_count=method_count,
                        generic_params=generic_params
                    ))
                    continue

                # Track interface brace depth using InterfaceParser
                if interface_parser.is_in_interface():
                    # Check brace depth before updating to catch methods on closing line
                    brace_depth_before = interface_parser.brace_depth
                    current_interface = interface_parser.get_current_interface()
                    still_in_interface = interface_parser.update_brace_depth(line)

                    # Check for interface method if we're still in interface or on closing line
                    if parse_methods and (
                        (still_in_interface and interface_parser.brace_depth > 0) or
                        (brace_depth_before > 0 and not still_in_interface and '{' not in stripped)
                    ):
                        sig = parse_go_def_signature(line, location=f"{resolved_path}:{line_num}")
                        if sig and sig.kind in ('func', 'method'):
                            # Interface methods don't have receivers in the interface definition
                            # but we track them as methods of the interface type
                            methods.append(Signature(
                                name=sig.name,
                                kind='method',
                                receiver=current_interface,
                                params=sig.params,
                                returns=sig.returns,
                                location=f"{resolved_path}:{line_num}",
                                is_public=sig.is_public,
                                has_body=False  # Methods in interface are stubs
                            ))

                    if not still_in_interface:
                        # Interface closed, parser already reset
                        pass
                    continue

    except Exception as e:
        print(f"Warning: Error processing interfaces from {file_path}: {e}")

    # Return interfaces first, then methods
    return interfaces + methods


def find_definition_line_index(code: str, is_type: bool) -> Optional[int]:
    """
    Find the line index (0-indexed) of the first definition in code.

    Args:
        code: Go code block content
        is_type: True to find type definition, False to find function definition

    Returns:
        Line index (0-indexed) of the definition, or None if not found
    """
    lines = code.split('\n')
    for i, line in enumerate(lines):
        stripped = line.strip()
        # Skip comments and empty lines
        if not stripped or stripped.startswith('//'):
            continue

        sig = parse_go_def_signature(line)
        if sig:
            if is_type:
                # Check for type definition
                if sig.kind not in ('func', 'method'):
                    return i
            else:
                # Check for function/method definition
                if sig.kind in ('func', 'method'):
                    return i

    return None


def is_example_definition(
    code: str,
    start_line: int,
    lines: List[str],
    heading: Optional[str],
    is_type: bool
) -> bool:
    """
    Check if a definition (type or function) is example code.

    Uses the existing find_definition_line_index to find the signature,
    then checks if it's example code using the unified is_example_code utility.

    Args:
        code: Go code block content
        start_line: Line number where code block starts (1-indexed)
        lines: All lines of the file (for context checking)
        heading: Optional heading text (for example detection)
        is_type: True for type definitions, False for function/method definitions

    Returns:
        True if the definition is example code, False otherwise
    """
    def_line_idx = find_definition_line_index(code, is_type=is_type)
    if def_line_idx is None:
        return False

    return is_example_code(
        code, start_line,
        lines=lines,
        heading_text=heading,
        check_prose_before_block=True,
        check_single_line=def_line_idx  # 0-indexed line number within code block
    )


def count_go_definitions(
    code: str,
    filter_example: bool = False,
    lines: Optional[List[str]] = None,
    start_line: int = 1,
    heading_text: Optional[str] = None
) -> Dict[str, int]:
    """
    Count Go definitions in code using the unified parser.

    Args:
        code: Go code block content
        filter_example: If True, skip example code when counting
        lines: All lines of the file (required if filter_example=True)
        start_line: Line number where code block starts (required if filter_example=True)
        heading_text: Optional heading text (used for example detection if filter_example=True)

    Returns:
        Dictionary with counts:
        - 'func': count of functions
        - 'method': count of methods
        - 'type': count of all types (struct, interface, alias, etc., excluding function types)
        - 'func_type': count of function types (type Name func(...))
    """
    counts: Dict[str, int] = {
        'func': 0,
        'method': 0,
        'type': 0,
        'func_type': 0
    }

    code_lines = code.split('\n')
    for line_idx, line in enumerate(code_lines):
        # Filter example code if requested
        if filter_example and lines is not None:
            if is_example_code(
                code, start_line,
                lines=lines,
                heading_text=heading_text,
                check_prose_before_block=True,
                check_single_line=line_idx
            ):
                continue

        # Try to parse as definition
        sig = parse_go_def_signature(line)
        if sig:
            if sig.kind == 'func':
                counts['func'] += 1
            elif sig.kind == 'method':
                counts['method'] += 1
            else:
                # All other kinds are types (struct, interface, alias, etc.)
                counts['type'] += 1
        else:
            # Check for function types (parse_go_def_signature excludes these)
            stripped = line.strip()
            if stripped and not stripped.startswith('//'):
                if _RE_FUNC_TYPE_DEF.match(line):
                    counts['func_type'] += 1

    return counts


def is_definition_start_line(line: str) -> bool:
    """
    Return True if the line starts a type, func, method, or func type definition.

    Used to detect signature-only code blocks: only definition-start lines
    and continuation lines (braces, indented body) should be present.

    Args:
        line: Single line of Go code (may have leading/trailing whitespace).

    Returns:
        True if the line is a definition start (type, func, method, func type).
    """
    stripped = line.strip()
    if not stripped:
        return False
    if parse_go_def_signature(stripped):
        return True
    if _RE_FUNC_TYPE_DEF.match(stripped):
        return True
    return False


def is_continuation_line(line: str) -> bool:
    """
    Return True if the line is brace-only or indented (part of a definition body).

    Used together with is_definition_start_line to classify lines in a code block
    when detecting signature-only blocks (struct/interface bodies, func bodies).

    Args:
        line: Single line of Go code (may have leading/trailing whitespace).

    Returns:
        True if the line is only braces or starts with whitespace (continuation).
    """
    stripped = line.strip()
    if stripped in ('{', '}'):
        return True
    if len(line) > 0 and line[0].isspace():
        return True
    return False


def is_signature_only_code_block(code: str) -> bool:
    """
    Return True if the Go code block contains only definition signatures and bodies.

    After removing comments, every non-empty line must be either a definition-start
    line (type/func/method/func type) or a continuation line (brace-only or
    indented body). Comment lines are removed and not counted.

    Used by documentation audits to skip requirement coverage for sections that
    only list API signatures (e.g. struct definitions with fields, method stubs).

    Args:
        code: Go code block content (no markdown fences).

    Returns:
        True if the block has at least one definition and no other substantive lines.
    """
    cleaned = remove_go_comments(code, multiline=True)
    non_empty_lines = [line for line in cleaned.split('\n') if line.strip()]
    if not non_empty_lines:
        return False
    definition_start_count = 0
    for line in non_empty_lines:
        if is_definition_start_line(line):
            definition_start_count += 1
        elif not is_continuation_line(line):
            return False
    return definition_start_count >= 1


def find_first_definition(code: str, is_type: bool) -> Optional[Signature]:
    """
    Find and parse the first definition in code.

    This is a convenience function that combines finding the definition line
    and parsing it into a Signature object.

    Args:
        code: Go code block content
        is_type: True for type definitions, False for function/method definitions

    Returns:
        Signature object or None if no definition found
    """
    def_line_idx = find_definition_line_index(code, is_type=is_type)
    if def_line_idx is None:
        return None
    code_lines = code.split('\n')
    return parse_go_def_signature(code_lines[def_line_idx])


def check_kind_word_after(
    heading: str,
    search_term: str,
    kind_word: str,
    display_name: str,
    error_prefix: str,
    errors: List[str]
) -> None:
    """
    Check if kind word appears immediately after the search term in heading.

    Args:
        heading: The heading text (will be converted to lowercase internally)
        search_term: The search term to look for (e.g., "Package" or "FileEntry.GetState")
        kind_word: The expected kind word (e.g., "Method", "Function", "Struct")
        display_name: The name to display in error messages
        error_prefix: Prefix for error messages (e.g., "Method heading")
        errors: List to append error messages to

    Returns:
        None (modifies errors list in place)
    """
    heading_lower = heading.lower()
    search_term_lower = search_term.lower()
    kind_word_lower = kind_word.lower()

    if search_term_lower in heading_lower:
        # Find position of search term in heading (case-insensitive)
        search_pos = heading_lower.find(search_term_lower)
        if search_pos != -1:
            # Check if kind word appears immediately after (allowing for whitespace)
            after_search = heading_lower[search_pos + len(search_term_lower):].strip()
            if not after_search.startswith(kind_word_lower):
                errors.append(
                    f'{error_prefix} should include "{kind_word}" immediately after '
                    f'{display_name}'
                )
        else:
            # If search term is present but we couldn't find it case-insensitively,
            # check if kind word is anywhere in heading
            if kind_word_lower not in heading_lower:
                errors.append(
                    f'{error_prefix} should include "{kind_word}" immediately after '
                    f'{display_name}'
                )
