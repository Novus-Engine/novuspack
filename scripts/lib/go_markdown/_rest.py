"""Go code block discovery, signatures, interfaces, definitions (part 2)."""

import re
from pathlib import Path
from typing import Dict, List, Optional, Tuple

from lib.go_markdown._base import (
    Signature,
    _RE_FUNC_NORMALIZE,
    _RE_FUNC_TYPE_DEF,
    _RE_FUNC_WITH_PARAMS,
    _RE_METHOD_NORMALIZE,
    _RE_RECEIVER_MATCH,
    _extract_receiver_type_safe,
    _format_normalized_signature,
    _normalize_go_signature_preprocessing,
    _normalize_package_names_general,
    _normalize_package_names_specific,
    _normalize_returns_simple,
    extract_receiver_type,
    find_go_code_blocks,
    is_example_code,
    is_public_name,
    parse_go_def_signature,
    remove_go_comments,
)


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
            elif char == ',' and not paren_depth and not bracket_depth:
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

    except (OSError, UnicodeDecodeError, ValueError, KeyError) as e:
        print(f"Warning: Error reading {file_path}: {e}")

    # Return interfaces first, then methods
    return interfaces + methods


def _maybe_append_interface_method(
    line, line_num, resolved_path, interface_parser, *,
    parse_methods, methods, stripped
):
    """If line is an interface method in body, append Signature to methods."""
    brace_depth_before = interface_parser.brace_depth
    current_interface = interface_parser.get_current_interface()
    still_in_interface = interface_parser.update_brace_depth(line)
    if not parse_methods:
        return
    if not (
        (still_in_interface and interface_parser.brace_depth > 0) or
        (brace_depth_before > 0 and not still_in_interface and '{' not in stripped)
    ):
        return
    sig = parse_go_def_signature(line, location=f"{resolved_path}:{line_num}")
    if sig and sig.kind in ('func', 'method'):
        methods.append(Signature(
            name=sig.name,
            kind='method',
            receiver=current_interface,
            params=sig.params,
            returns=sig.returns,
            location=f"{resolved_path}:{line_num}",
            is_public=sig.is_public,
            has_body=False
        ))


def _count_interface_methods(block_lines, start_index, interface_parser) -> int:
    """Count method signatures in interface body from block_lines starting at start_index."""
    method_count = 0
    temp_brace_depth = interface_parser.brace_depth
    for j in range(start_index + 1, len(block_lines)):
        temp_line = block_lines[j]
        temp_stripped = temp_line.strip()
        if not temp_stripped or temp_stripped.startswith('//'):
            continue
        temp_brace_depth += temp_stripped.count('{') - temp_stripped.count('}')
        temp_sig = parse_go_def_signature(temp_line, location="")
        if (temp_sig and temp_sig.kind in ('func', 'method') and temp_brace_depth > 0):
            method_count += 1
        if temp_brace_depth <= 0:
            break
    return method_count


def extract_interfaces_from_markdown(
    content: str,
    file_path: Path,
    *,
    start_line: int = 1,
    parse_methods: bool = True,
    skip_examples: bool = True,
    lines: Optional[List[str]] = None,
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
    _ = start_line  # reserved for future use
    interfaces = []
    methods = []

    try:
        resolved_path = file_path.resolve()
        go_blocks = find_go_code_blocks(content)

        for block_start_line, _block_end_line, code_content in go_blocks:
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

                    method_count = (
                        _count_interface_methods(block_lines, i, interface_parser)
                        if (has_full_body and parse_methods)
                        else 0
                    )
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

                if interface_parser.is_in_interface():
                    _maybe_append_interface_method(
                        line, line_num, resolved_path, interface_parser,
                        parse_methods=parse_methods, methods=methods, stripped=stripped
                    )
                    continue

    except (OSError, UnicodeDecodeError, ValueError, KeyError) as e:
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
    *,
    display_name: str,
    error_prefix: str,
    errors: List[str],
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
            # Check if kind word appears immediately after (allow optional ` and whitespace)
            after = heading_lower[search_pos + len(search_term_lower):].strip()
            after_search = after.lstrip('`').strip()
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
