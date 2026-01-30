#!/usr/bin/env python3
"""
Validate Go signature consistency within tech specs documentation.

This script checks for inconsistencies where the same function, method, or type
is defined multiple times in tech specs with different signatures.

It detects:
- Errors: Same signature defined with different parameters/return types
- Errors: Methods defined for interfaces that are NOT in the canonical interface definition
- Errors: Interface stubs with methods not in canonical definition
- Errors: Interface stub methods that differ from canonical definitions
- Errors: Duplicate identical signatures
- Warnings: Type/interface stubs (expected when fleshing out methods in separate sections)

Usage:
    python3 scripts/validate_go_spec_signature_consistency.py [options]

Options:
    --verbose, -v           Show detailed progress information
    --output, -o FILE       Write detailed report to FILE
    --path, -p PATHS        Check only the specified file(s) or
                            directory(ies) (recursive). Can be a single
                            path or comma-separated list of paths
    --nocolor, --no-color   Disable colored output
    --help, -h              Show this help message

Examples:
    # Basic validation
    python3 scripts/validate_go_spec_signature_consistency.py

    # Save report to file
    python3 scripts/validate_go_spec_signature_consistency.py \
        --output tmp/signature_consistency_report.txt

    # Verbose output
    python3 scripts/validate_go_spec_signature_consistency.py --verbose

    # Check specific file
    python3 scripts/validate_go_spec_signature_consistency.py \
        --path docs/tech_specs/api_basic_operations.md
"""

import re
import sys
from collections import defaultdict
from pathlib import Path
from typing import Dict, List, Optional, Tuple

scripts_dir = Path(__file__).parent
lib_dir = scripts_dir / "lib"

# Import shared utilities
for module_path in (str(scripts_dir), str(lib_dir)):
    if module_path not in sys.path:
        sys.path.insert(0, module_path)

from lib._validation_utils import (  # noqa: E402
    OutputBuilder, parse_no_color_flag,
    find_markdown_files, parse_paths, get_workspace_root,
    get_validation_exit_code, HeadingContext, find_heading_before_line,
    ValidationIssue,
    DOCS_DIR, TECH_SPECS_DIR
)
from lib._go_code_utils import (  # noqa: E402
    parse_go_def_signature,
    is_example_code, is_example_signature_name,
    find_go_code_blocks, Signature, is_public_name,
    extract_interfaces_from_markdown,
    InterfaceParser, normalize_go_signature
)

# Compiled regex patterns for performance (module level)
_RE_INTERFACE_PATTERN = re.compile(r'^\s*type\s+\w+(?:\s*\[[^\]]+\])?\s+interface\s*\{')
_RE_STRUCT_PATTERN = re.compile(r'^\s*type\s+(\w+)(?:\s*(\[[^\]]+\]))?\s+struct\s*\{')


def count_interface_methods(content: str, start_line: int, end_line: int) -> int:
    """
    Count methods in an interface definition within a specific line range.

    This is a specialized lightweight utility that uses the shared InterfaceParser
    class for brace depth tracking. It's simpler than the full extraction helpers
    since it only needs to count methods, not extract full signatures.

    Note: This function is currently not called but kept for potential future use.
    If needed, it could be refactored to use extract_interfaces_from_markdown(),
    but the current implementation is acceptable as it uses shared utilities.
    """
    lines = content.split('\n')
    method_count = 0
    interface_parser = InterfaceParser()

    for i in range(start_line - 1, min(end_line, len(lines))):
        line = lines[i]

        # Check for interface start using InterfaceParser
        interface_name = interface_parser.check_interface_start(line)
        if interface_name:
            continue

        if interface_parser.is_in_interface():
            still_in_interface = interface_parser.update_brace_depth(line)
            # Check for method signature if still in interface body
            if still_in_interface and interface_parser.brace_depth > 0:
                sig = parse_go_def_signature(line, location="")
                if sig and sig.kind in ('func', 'method'):
                    method_count += 1

            if not still_in_interface:
                break

    return method_count


def count_struct_fields(content: str, start_line: int, end_line: int) -> int:
    """Count fields in a struct definition."""
    lines = content.split('\n')
    field_count = 0
    in_struct = False
    brace_depth = 0

    for i in range(start_line - 1, min(end_line, len(lines))):
        line = lines[i]
        stripped = line.strip()

        # Check for struct start
        if re.match(r'^\s*type\s+\w+\s+struct\s*\{', line):
            in_struct = True
            brace_depth = stripped.count('{') - stripped.count('}')
            continue

        if in_struct:
            brace_depth += stripped.count('{') - stripped.count('}')
            # Check for field (not a comment, not empty, not a method)
            if brace_depth > 0 and stripped and not stripped.startswith('//'):
                # Simple heuristic: if it looks like a field (has identifier and type)
                if re.match(r'^\s*\w+\s+\w+', stripped):
                    # Make sure it's not a method
                    if not re.match(r'^\s*func\s+', stripped):
                        field_count += 1

            if brace_depth <= 0:
                break

    return field_count


# is_example_signature_name now imported from _go_code_utils


def extract_signatures_from_markdown_file(file_path: Path, repo_root: Path) -> List[Signature]:
    """Extract all signatures from Go code blocks in a markdown file."""
    signatures = []

    try:
        # Get relative path from repo root
        try:
            relative_path = file_path.resolve().relative_to(repo_root.resolve())
        except ValueError:
            # If path is not under repo_root, use absolute path as fallback
            relative_path = file_path.resolve()
        content = file_path.read_text(encoding='utf-8')
        lines = content.split('\n')

        # Use shared helper to extract interfaces and their methods
        interface_signatures = extract_interfaces_from_markdown(
            content, file_path, start_line=1, parse_methods=True,
            skip_examples=True, lines=lines
        )
        signatures.extend(interface_signatures)

        # Extract other signatures (structs, functions, methods, types) that aren't interfaces
        go_blocks = find_go_code_blocks(content)

        for start_line, end_line, code_content in go_blocks:
            block_lines = code_content.split('\n')

            for i, line in enumerate(block_lines):
                line_num = start_line + i
                stripped = line.strip()

                # Skip empty lines and comments
                if not stripped or stripped.startswith('//'):
                    continue

                # Skip interface definitions (already handled by helper)
                if _RE_INTERFACE_PATTERN.match(stripped):
                    continue

                # Check if this is an example signature using the library function
                is_example = is_example_code(
                    code_content, start_line,
                    lines=lines,
                    check_single_line=i
                )

                # Check for struct start
                struct_match = _RE_STRUCT_PATTERN.match(stripped)
                if struct_match:
                    # Skip if this is an example
                    if is_example:
                        continue

                    name = struct_match.group(1)
                    generic_params = struct_match.group(2)  # e.g., "[T any]"
                    is_public = is_public_name(name) if name else False
                    brace_depth = stripped.count('{') - stripped.count('}')
                    has_full_body = brace_depth > 0

                    # Count fields if it's a full struct
                    field_count = 0
                    if has_full_body:
                        # Count fields within this code block
                        temp_brace_depth = brace_depth
                        for j in range(i + 1, len(block_lines)):
                            temp_line = block_lines[j]
                            temp_stripped = temp_line.strip()
                            if not temp_stripped or temp_stripped.startswith('//'):
                                continue
                            temp_brace_depth += temp_stripped.count('{') - temp_stripped.count('}')
                            if temp_brace_depth > 0 and temp_stripped:
                                # Simple heuristic: if it looks like a field
                                # (has identifier and type)
                                if re.match(r'^\s*\w+\s+\w+', temp_stripped):
                                    # Make sure it's not a method
                                    if not re.match(r'^\s*func\s+', temp_stripped):
                                        field_count += 1
                            if temp_brace_depth <= 0:
                                break

                    signatures.append(Signature(
                        name=name,
                        kind='type',
                        location=f"{relative_path}:{line_num}",
                        is_public=is_public,
                        has_body=has_full_body,
                        field_count=field_count,
                        generic_params=generic_params
                    ))
                    continue

                # Check for any Go definition (function, method, or type)
                sig = parse_go_def_signature(line, location=f"{relative_path}:{line_num}")
                if sig:
                    if sig.kind in ('func', 'method'):
                        # Standalone method/function definitions are full
                        signatures.append(Signature(
                            name=sig.name,
                            kind=sig.kind,
                            receiver=sig.receiver,
                            params=sig.params,
                            returns=sig.returns,
                            location=f"{relative_path}:{line_num}",
                            is_public=sig.is_public,
                            has_body=True
                        ))
                    else:
                        # Type definition (not interface/struct, already handled)
                        # Skip if this is an example
                        if is_example:
                            continue

                        if sig.kind != 'interface':  # Interfaces already handled
                            signatures.append(Signature(
                                name=sig.name,
                                kind=sig.kind,
                                location=f"{relative_path}:{line_num}",
                                is_public=sig.is_public,
                                has_body=False,  # Type aliases don't have bodies
                                generic_params=sig.generic_params
                            ))

    except (IOError, OSError) as e:
        # File read errors - log to stderr
        print(f"Warning: Could not read {file_path}: {e}", file=sys.stderr)
    except UnicodeDecodeError as e:
        # Encoding errors - log to stderr
        print(f"Warning: Could not decode {file_path} (encoding issue): {e}", file=sys.stderr)
    except Exception as e:
        # Unexpected errors - log to stderr
        print(f"Warning: Unexpected error reading {file_path}: {e}", file=sys.stderr)

    return signatures


def extract_heading_context(file_path: Path, line_num: int) -> Optional[HeadingContext]:
    """
    Extract heading context for a given line number in a markdown file.

    Returns the most recent heading before the given line number.
    Uses shared utility from _go_code_utils.
    """
    try:
        content = file_path.read_text(encoding='utf-8')
        ctx = find_heading_before_line(content, line_num, prefer_deepest=True)
        if ctx:
            # Add file_path to the context
            return HeadingContext(
                heading_text=ctx.heading_text,
                heading_level=ctx.heading_level,
                heading_line=ctx.heading_line,
                file_path=str(file_path)
            )
    except (IOError, OSError, UnicodeDecodeError):
        # File read/encoding errors - return None silently
        pass
    except (ValueError, IndexError, KeyError):
        # Data structure errors - return None silently
        pass
    except Exception as e:
        # Unexpected errors - log but don't fail
        print(f"Warning: Unexpected error reading {file_path}: {e}", file=sys.stderr)
        pass

    return None


def get_signature_heading_context(sig: Signature, repo_root: Path) -> Optional[HeadingContext]:
    """Get heading context for a signature."""
    location_parts = sig.location.split(':', 1)
    if len(location_parts) != 2:
        return None

    file_path = repo_root / location_parts[0]
    try:
        line_num = int(location_parts[1])
        return extract_heading_context(file_path, line_num)
    except (ValueError, TypeError):
        return None


def score_canonical_signature(
    sig: Signature,
    heading_ctx: Optional[HeadingContext],
    all_sigs: List[Signature]
) -> Tuple[float, str]:
    """
    Score a signature to determine if it's likely canonical.

    Returns:
        Tuple of (score, reason) where higher score = more likely canonical.
        Score ranges from 0.0 to 1.0.
    """
    score = 0.0
    reasons = []

    # Prefer less deeply nested headings (lower level = higher score)
    if heading_ctx:
        # H1 = 1.0, H2 = 0.8, H3 = 0.6, H4 = 0.4, H5 = 0.2, H6 = 0.1
        heading_score = max(0.1, 1.0 - (heading_ctx.heading_level - 1) * 0.2)
        score += heading_score * 0.4  # 40% weight
        reasons.append(f"heading level {heading_ctx.heading_level}")

    # Prefer signatures with bodies (more complete definitions)
    if sig.has_body:
        score += 0.2  # 20% weight
        reasons.append("has body")

    # Prefer earlier line numbers (first occurrence)
    file_path = sig.location.split(':', 1)[0]
    same_file_sigs = [s for s in all_sigs if s.location.startswith(file_path)]
    if same_file_sigs:
        line_nums = []
        for s in same_file_sigs:
            try:
                line_nums.append(int(s.location.split(':', 1)[1]))
            except (ValueError, IndexError):
                pass
        if line_nums:
            try:
                sig_line = int(sig.location.split(':', 1)[1])
                if sig_line == min(line_nums):
                    score += 0.15  # 15% weight
                    reasons.append("first occurrence in file")
            except (ValueError, IndexError):
                pass

    # Prefer files with more general names (e.g., api_core.md over api_core_advanced.md)
    file_name = file_path.lower()
    if 'core' in file_name or 'basic' in file_name:
        score += 0.1  # 10% weight
        reasons.append("core/basic file")
    elif 'advanced' in file_name or 'extended' in file_name:
        score -= 0.05  # Penalty
        reasons.append("advanced/extended file")

    # Prefer signatures in sections with relevant keywords and signature name
    if heading_ctx:
        heading_lower = heading_ctx.heading_text.lower()
        sig_name_lower = sig.name.lower()

        # Check for signature name in heading
        if sig_name_lower in heading_lower:
            score += 0.15  # 15% weight
            reasons.append("signature name in heading")

        # "definition" or "definitions" applies to all signature types
        if 'definition' in heading_lower or 'definitions' in heading_lower:
            score += 0.1  # 10% weight
            reasons.append("definition keyword in heading")

        # For type/interface/struct definitions: look for type-related keywords
        if sig.kind in ('type', 'interface'):
            type_keywords = ['type', 'struct', 'types', 'interfaces']
            if any(keyword in heading_lower for keyword in type_keywords):
                score += 0.1  # 10% weight
                reasons.append("type-related keyword in heading")

        # For functions: look for function-related keywords
        elif sig.kind == 'func':
            func_keywords = ['function', 'functions', 'func', 'operation', 'operations']
            if any(keyword in heading_lower for keyword in func_keywords):
                score += 0.1  # 10% weight
                reasons.append("function-related keyword in heading")

        # For methods: look for method-related keywords and receiver type
        elif sig.kind == 'method' and sig.receiver:
            method_keywords = ['method', 'methods']
            if any(keyword in heading_lower for keyword in method_keywords):
                score += 0.1  # 10% weight
                reasons.append("method-related keyword in heading")

            # Also check if receiver type name is in heading
            receiver_lower = sig.receiver.lower()
            if receiver_lower in heading_lower:
                score += 0.1  # 10% weight
                reasons.append("receiver type in heading")

    return (min(1.0, max(0.0, score)), ", ".join(reasons) if reasons else "no specific indicators")


def find_canonical_signature(
    sigs: List[Signature],
    repo_root: Path
) -> Tuple[Optional[Signature], Optional[HeadingContext], float, str]:
    """
    Find the most likely canonical signature from a list of duplicates.

    Returns:
        Tuple of (canonical_sig, heading_context, confidence, reason)
        confidence: 0.0 to 1.0, where >= 0.7 is high confidence
    """
    if not sigs:
        return None, None, 0.0, ""

    if len(sigs) == 1:
        ctx = get_signature_heading_context(sigs[0], repo_root)
        return sigs[0], ctx, 1.0, "only one occurrence"

    # Score each signature
    scored = []
    for sig in sigs:
        ctx = get_signature_heading_context(sig, repo_root)
        score, reason = score_canonical_signature(sig, ctx, sigs)
        scored.append((sig, ctx, score, reason))

    # Sort by score (highest first)
    scored.sort(key=lambda x: x[2], reverse=True)

    best_sig, best_ctx, best_score, best_reason = scored[0]

    # Calculate confidence based on score gap
    # Heading level is already factored into the score, so we don't override here
    if len(scored) > 1:
        second_score = scored[1][2]
        score_gap = best_score - second_score
        # Confidence based on score and gap between best and second best
        confidence = min(1.0, best_score * (1.0 + score_gap))
    else:
        confidence = best_score

    return best_sig, best_ctx, confidence, best_reason


# normalize_go_signature now imported from _go_code_utils (enhanced version)


def find_canonical_definition(signatures: List[Signature]) -> Optional[Signature]:
    """
    Find the canonical definition from a list of signatures.

    For types/interfaces: prefer the one with most fields/methods.
    For methods/functions: prefer the one with has_body=True.
    """
    if not signatures:
        return None

    # For methods/functions, prefer full definitions
    if signatures[0].kind in ('method', 'func'):
        for sig in signatures:
            if sig.has_body:
                return sig
        # If no full definition, return first
        return signatures[0]

    # For types/interfaces, prefer the one with most fields/methods
    if signatures[0].kind in ('type', 'interface'):
        canonical = signatures[0]
        max_count = (
            canonical.method_count
            if canonical.kind == 'interface'
            else canonical.field_count
        )

        for sig in signatures[1:]:
            count = sig.method_count if sig.kind == 'interface' else sig.field_count
            if count > max_count:
                canonical = sig
                max_count = count
            elif count == max_count and sig.has_body and not canonical.has_body:
                # Prefer one with body if counts are equal
                canonical = sig

        return canonical

    return signatures[0]


def check_signature_consistency(
    signatures: Dict[str, List[Signature]],
    verbose: bool = False,
    repo_root: Optional[Path] = None
) -> List[ValidationIssue]:
    """
    Check for signature inconsistencies.

    Returns:
        - List of ValidationIssue objects (with severity='error' or 'warning')
    """
    issues: List[ValidationIssue] = []

    def parse_location(location_str: str) -> Tuple[Path, int]:
        """Parse location string 'file:line' into Path and line number."""
        if ':' in location_str:
            file_str, line_str = location_str.split(':', 1)
            try:
                line_num = int(line_str)
            except ValueError:
                line_num = 1
            return Path(file_str), line_num
        return Path(location_str), 1

    def get_first_location(location_strs: List[str]) -> Tuple[Path, int]:
        """Get the first location from a list of location strings."""
        if location_strs:
            return parse_location(location_strs[0])
        return Path("unknown"), 1

    # Collect all types/interfaces first to check for conflicts across groups
    all_types = []
    for sig_list in signatures.values():
        for sig in sig_list:
            if sig.kind in ('type', 'interface'):
                all_types.append(sig)

    # Check types/interfaces for conflicts (grouped by base name)
    if all_types:
        # Group by base name
        by_name = {}
        for sig in all_types:
            if sig.name not in by_name:
                by_name[sig.name] = []
            by_name[sig.name].append(sig)

        # Check each name group for conflicts
        for name, sigs in by_name.items():
            if len(sigs) <= 1:
                continue

            # Create normalized signatures including generics
            normalized_sigs = {}
            for sig in sigs:
                # Include generics in signature
                sig_str = sig.name
                if sig.generic_params:
                    sig_str += sig.generic_params

                if sig_str not in normalized_sigs:
                    normalized_sigs[sig_str] = []
                normalized_sigs[sig_str].append(sig)

            # If different normalized signatures, it's a conflict
            if len(normalized_sigs) > 1:
                locations = [
                    sig.location
                    for sig_list in normalized_sigs.values()
                    for sig in sig_list
                ]
                file_path, line_num = get_first_location(locations)
                message_parts = [f"Conflicting type definitions for '{name}':"]
                for norm_sig, sig_list in normalized_sigs.items():
                    message_parts.append(f"  {norm_sig}:")
                    for sig in sig_list:
                        message_parts.append(f"    Location: {sig.location}")
                issues.append(ValidationIssue(
                    "conflicting_type_definitions",
                    file_path,
                    line_num,
                    line_num,
                    "\n".join(message_parts),
                    severity='error',
                    type_name=name,
                    locations=locations
                ))
                continue  # Skip stub detection for this group

            # If same normalized signature, check for duplicates or stubs
            canonical = find_canonical_definition(sigs)
            if not canonical:
                # No canonical found, but we have multiple identical signatures
                # This is a duplicate - will be handled by the duplicate_type check below
                continue

            stubs = [s for s in sigs if s != canonical]

            # If there are no stubs (all are identical and canonical), it's still a duplicate
            if not stubs and len(sigs) > 1:
                locations = [sig.location for sig in sigs]
                file_path, line_num = get_first_location(locations)
                message_parts = [f"Duplicate identical type/interface for '{name}':"]
                for sig in sigs:
                    message_parts.append(f"  Location: {sig.location}")
                issues.append(ValidationIssue(
                    "duplicate_type",
                    file_path,
                    line_num,
                    line_num,
                    "\n".join(message_parts),
                    severity='warning',
                    type_name=name,
                    locations=locations
                ))

            # Warn about stubs
            for stub in stubs:
                if not stub.has_body or (
                    canonical.has_body and
                    canonical.method_count > stub.method_count
                ):
                    canonical_count = (
                        canonical.method_count
                        if canonical.kind == 'interface'
                        else canonical.field_count
                    )
                    stub_count = (
                        stub.method_count
                        if stub.kind == 'interface'
                        else stub.field_count
                    )
                    count_type = (
                        'method_count'
                        if canonical.kind == 'interface'
                        else 'field_count'
                    )
                    canonical_file, canonical_line = parse_location(canonical.location)
                    stub_file, stub_line = parse_location(stub.location)
                    message = (
                        f"Type/interface stub detected for '{name}':"
                        f"  Canonical: {canonical.location} "
                        f"(has_body={canonical.has_body}, "
                        f"{count_type}={canonical_count})\n"
                        f"  Stub: {stub.location} "
                        f"(has_body={stub.has_body}, "
                        f"{count_type}={stub_count})"
                    )
                    issues.append(ValidationIssue(
                        "type_stub",
                        stub_file,
                        stub_line,
                        stub_line,
                        message,
                        severity='warning',
                        type_name=name,
                        canonical_location=canonical.location,
                        stub_location=stub.location
                    ))

            # For interfaces, check method signatures in stubs vs canonical
            if canonical.kind == 'interface':
                # First, collect methods that are in the canonical interface definition
                # (methods inside the interface body have has_body=False and are in the same file
                #  and within a reasonable line range of the canonical definition)
                canonical_method_names = set()
                canonical_file, canonical_line_str = canonical.location.split(':', 1)
                try:
                    canonical_line = int(canonical_line_str)
                except ValueError:
                    canonical_line = 0

                # Methods in the canonical interface should be in the same file,
                # have has_body=False (inside interface body), and be within ~100 lines
                # of the interface definition (reasonable range for interface body)
                for sig_list in signatures.values():
                    for sig in sig_list:
                        if (sig.kind == 'method' and
                                sig.receiver == canonical.name and
                                not sig.has_body):
                            sig_file, sig_line_str = sig.location.split(':', 1)
                            try:
                                sig_line = int(sig_line_str)
                                # Check if method is in same file and within reasonable range
                                if (sig_file == canonical_file and
                                        canonical_line <= sig_line <= canonical_line + 100):
                                    canonical_method_names.add(sig.name)
                            except ValueError:
                                pass

                # Get all methods for this interface from all signatures
                interface_methods = {}
                for sig_list in signatures.values():
                    for sig in sig_list:
                        if sig.kind == 'method' and sig.receiver == canonical.name:
                            method_key = sig.name
                            if method_key not in interface_methods:
                                interface_methods[method_key] = []
                            interface_methods[method_key].append(sig)

                # Check if any methods are defined that are NOT in the canonical definition
                if canonical_method_names and canonical.has_body:
                    for method_key, method_sigs in interface_methods.items():
                        if method_key not in canonical_method_names:
                            # Method is defined but not in canonical interface definition
                            locations = [msig.location for msig in method_sigs]
                            file_path, line_num = get_first_location(locations)
                            message_parts = [
                                f"Method '{method_key}' is defined for interface "
                                f"'{canonical.name}' but is not in the canonical "
                                f"interface definition:",
                                f"  Canonical interface: {canonical.location}",
                                f"  Methods in canonical: {sorted(canonical_method_names)}"
                            ]
                            for method_sig in method_sigs:
                                message_parts.append(f"  Method location: {method_sig.location}")
                            issues.append(ValidationIssue(
                                "method_not_in_canonical_interface",
                                file_path,
                                line_num,
                                line_num,
                                "\n".join(message_parts),
                                severity='error',
                                method_name=method_key,
                                interface_name=canonical.name,
                                canonical_location=canonical.location,
                                locations=locations
                            ))

                # Check each method for consistency
                for method_key, method_sigs in interface_methods.items():
                    if len(method_sigs) <= 1:
                        continue

                    # Find canonical method (prefer has_body=True)
                    canonical_method = None
                    stub_methods = []
                    for method_sig in method_sigs:
                        if method_sig.has_body:
                            canonical_method = method_sig
                        else:
                            stub_methods.append(method_sig)

                    if not canonical_method:
                        # All are stubs, compare them
                        normalized_methods = {}
                        for method_sig in method_sigs:
                            norm_sig = normalize_go_signature(method_sig.normalized_signature())
                            if norm_sig not in normalized_methods:
                                normalized_methods[norm_sig] = []
                            normalized_methods[norm_sig].append(method_sig)

                        if len(normalized_methods) > 1:
                            locations = [
                                msig.location
                                for msigs in normalized_methods.values()
                                for msig in msigs
                            ]
                            file_path, line_num = get_first_location(locations)
                            message_parts = [
                                f"Method signature inconsistency for "
                                f"'{method_key}' in interface '{canonical.name}':"
                            ]
                            for norm_sig, msigs in normalized_methods.items():
                                message_parts.append(f"  Signature: {norm_sig}")
                                for msig in msigs:
                                    message_parts.append(f"    Location: {msig.location}")
                            issues.append(ValidationIssue(
                                "method_signature_inconsistency",
                                file_path,
                                line_num,
                                line_num,
                                "\n".join(message_parts),
                                severity='error',
                                method_name=method_key,
                                interface_name=canonical.name,
                                locations=locations
                            ))
                    else:
                        # Compare stubs against canonical
                        for stub_method in stub_methods:
                            canonical_norm = normalize_go_signature(
                                canonical_method.normalized_signature()
                            )
                            stub_norm = normalize_go_signature(
                                stub_method.normalized_signature()
                            )

                            if canonical_norm != stub_norm:
                                stub_file, stub_line = parse_location(stub_method.location)
                                message = (
                                    f"Interface stub method "
                                    f"'{method_key}' differs from canonical:"
                                    f"  Canonical: {canonical_norm}\n"
                                    f"    Location: {canonical_method.location}\n"
                                    f"  Stub: {stub_norm}\n"
                                    f"    Location: {stub_method.location}"
                                )
                                issues.append(ValidationIssue(
                                    "interface_stub_method_differs",
                                    stub_file,
                                    stub_line,
                                    stub_line,
                                    message,
                                    severity='error',
                                    method_name=method_key,
                                    canonical_location=canonical_method.location,
                                    stub_location=stub_method.location
                                ))

        # Check methods/functions (grouped by normalized_key)
        for key, sig_list in signatures.items():
            if len(sig_list) <= 1:
                continue

            # Group by kind - skip types/interfaces as they're handled above
            methods = [s for s in sig_list if s.kind in ('method', 'func')]

            # Check methods/functions
            if methods:
                # Normalize signatures for comparison
                normalized_sigs = {}
                for sig in methods:
                    norm_sig = normalize_go_signature(sig.normalized_signature())
                    if norm_sig not in normalized_sigs:
                        normalized_sigs[norm_sig] = []
                    normalized_sigs[norm_sig].append(sig)

                # Check for different signatures
                if len(normalized_sigs) > 1:
                    # Different signatures found - ERROR
                    locations = [
                        sig.location
                        for sigs in normalized_sigs.values()
                        for sig in sigs
                    ]
                    file_path, line_num = get_first_location(locations)
                    message_parts = [f"Signature inconsistency for '{key}':"]
                    for norm_sig, sigs in normalized_sigs.items():
                        message_parts.append(f"  Signature: {norm_sig}")
                        for sig in sigs:
                            message_parts.append(f"    Location: {sig.location}")
                    issues.append(ValidationIssue(
                        "signature_inconsistency",
                        file_path,
                        line_num,
                        line_num,
                        "\n".join(message_parts),
                        severity='error',
                        signature_key=key,
                        locations=locations
                    ))
                elif len(methods) > 1:
                    # Same signature, multiple locations - ERROR
                    # This will be handled in main() with canonical detection
                    locations = [sig.location for sig in methods]
                    file_path, line_num = get_first_location(locations)
                    message_parts = [f"Duplicate identical signature for '{key}':"]
                    for sig in methods:
                        message_parts.append(f"  Location: {sig.location}")
                    issues.append(ValidationIssue(
                        "duplicate_signature",
                        file_path,
                        line_num,
                        line_num,
                        "\n".join(message_parts),
                        severity='error',
                        signature_key=key,
                        locations=locations
                    ))

    return issues


def main():
    """Main entry point."""
    # Show help if requested
    if '--help' in sys.argv or '-h' in sys.argv:
        print(__doc__)
        return 0

    # Parse command line arguments
    verbose = '--verbose' in sys.argv or '-v' in sys.argv
    no_color = parse_no_color_flag(sys.argv)
    no_fail = '--no-fail' in sys.argv
    output_file = None
    target_paths_str = None

    for i, arg in enumerate(sys.argv):
        if arg in ('--output', '-o') and i + 1 < len(sys.argv):
            output_file = sys.argv[i + 1]
        elif arg in ('--path', '-p') and i + 1 < len(sys.argv):
            target_paths_str = sys.argv[i + 1]

    # Parse comma-separated paths
    target_paths = parse_paths(target_paths_str)

    # Create output builder (header streams immediately if verbose)
    output = OutputBuilder(
        "Go Signature Consistency",
        "Validates signature consistency within tech specs",
        no_color=no_color,
        verbose=verbose,
        output_file=output_file
    )

    # Find repository root
    repo_root = get_workspace_root()

    # Find markdown files to audit
    default_dir = Path(DOCS_DIR) / TECH_SPECS_DIR
    md_files = find_markdown_files(
        target_paths=target_paths,
        default_dir=default_dir,
        verbose=verbose
    )

    if not md_files:
        output.add_error_line('No markdown files found')
        output.print()
        return 1

    if verbose:
        output.add_verbose_line(f'Found {len(md_files)} markdown file(s) to audit')
        output.add_blank_line("working_verbose")

    # Collect all signatures
    all_signatures = []
    for md_file in md_files:
        if verbose:
            output.add_verbose_line(f'Extracting signatures from {md_file.name}...')
        file_sigs = extract_signatures_from_markdown_file(md_file, repo_root)
        all_signatures.extend(file_sigs)

    if verbose:
        output.add_verbose_line(f'Found {len(all_signatures)} total signatures')
        output.add_blank_line("working_verbose")

    # Group signatures by normalized key, filtering out examples
    signatures_by_key = defaultdict(list)
    for sig in all_signatures:
        # Skip example signatures
        if is_example_signature_name(sig.name):
            continue
        key = sig.normalized_key()
        signatures_by_key[key].append(sig)

    # Check for inconsistencies
    if verbose:
        output.add_verbose_line('Checking for signature inconsistencies...')
        output.add_blank_line("working_verbose")

    issues = check_signature_consistency(
        signatures_by_key, verbose, repo_root=repo_root
    )

    # Check if signature counts don't match (indicates duplicates that weren't caught)
    # Find all duplicate keys and generate specific warnings for any that weren't reported
    duplicate_count = len(all_signatures) - len(signatures_by_key)
    if duplicate_count > 0:
        # Track which keys were already reported in issues
        reported_keys = set()
        import re
        for issue in issues:
            # issue is a ValidationIssue
            if isinstance(issue, ValidationIssue):
                # Try to extract key from extra_fields
                key = (
                    issue.extra_fields.get('signature_key') or
                    issue.extra_fields.get('type_name') or
                    issue.extra_fields.get('method_name')
                )
                if key:
                    reported_keys.add(key)
                else:
                    # Fallback: try to extract from message
                    match = re.search(r"for '([^']+)'|'([^']+)'", issue.message)
                    if match:
                        reported_keys.add(match.group(1) or match.group(2))

        # Find which keys have duplicates that weren't already reported
        for key, sig_list in signatures_by_key.items():
            # Filter out any example signatures that might have slipped through
            sig_list = [sig for sig in sig_list if not is_example_signature_name(sig.name)]
            if len(sig_list) < 2:
                continue  # Not enough non-example signatures to be a duplicate
            if key not in reported_keys:
                # Find canonical signature
                canonical_sig, canonical_ctx, confidence, reason = find_canonical_signature(
                    sig_list, repo_root
                )

                # Build error message with canonical suggestion
                locations = [sig.location for sig in sig_list]
                # Get first location from list
                if locations:
                    location_str = locations[0]
                    if ':' in location_str:
                        file_str, line_str = location_str.split(':', 1)
                        try:
                            line_num = int(line_str)
                        except ValueError:
                            line_num = 1
                        file_path = Path(file_str)
                    else:
                        file_path = Path(location_str)
                        line_num = 1
                else:
                    file_path = Path("unknown")
                    line_num = 1
                # Build message parts in a single pass through sig_list
                if canonical_sig and confidence >= 0.7:
                    # High confidence - suggest canonical
                    message_parts = [
                        f"Duplicate identical signature for '{key}':",
                        f"  Suggested canonical (confidence: {confidence:.0%}): "
                        f"{canonical_sig.location} (reason: {reason})",
                        "  Other locations:"
                    ]
                    # Single loop: collect non-canonical locations
                    other_locations = [
                        f"    {sig.location}" for sig in sig_list if sig != canonical_sig
                    ]
                    message_parts.extend(other_locations)
                    suggestion = f"Use canonical: {canonical_sig.location}"
                else:
                    # Low confidence - show all locations without suggestion
                    message_parts = [f"Duplicate identical signature for '{key}':"]
                    # Single loop: collect all locations
                    all_locations = [f"  Location: {sig.location}" for sig in sig_list]
                    message_parts.extend(all_locations)
                    if canonical_sig:
                        message_parts.append(
                            f"  Note: {canonical_sig.location} may be canonical "
                            f"(low confidence: {confidence:.0%}, reason: {reason})"
                        )
                    suggestion = None
                issues.append(ValidationIssue(
                    "duplicate_signature",
                    file_path,
                    line_num,
                    line_num,
                    "\n".join(message_parts),
                    severity='error',
                    suggestion=suggestion,
                    signature_key=key,
                    locations=locations,
                    canonical_location=canonical_sig.location if canonical_sig else None
                ))

    # Report results (headers will be added automatically by add_error_line/add_warning_line)
    # Filter issues by severity in a single loop
    errors = []
    warnings = []
    for issue in issues:
        if issue.matches(severity='error'):
            errors.append(issue)
        if issue.matches(severity='warning'):
            warnings.append(issue)

    for error in errors:
        output.add_error_line(error.format_message(no_color=no_color))

    for warning in warnings:
        output.add_warning_line(warning.format_message(no_color=no_color))

    # Always add summary section
    summary_items = [
        ("Signatures checked:", len(all_signatures)),
        ("Unique definitions:", len(signatures_by_key)),
    ]
    if errors:
        summary_items.append(("Errors found:", len(errors)))
    if warnings:
        summary_items.append(("Warnings found:", len(warnings)))
    output.add_summary_header()
    output.add_summary_section(summary_items)

    # Final message: success if no errors, but mention warnings if present
    if not errors:
        if warnings:
            output.add_success_message(
                'All signatures are consistent! (Some warnings were found - see above)'
            )
        else:
            output.add_success_message('All signatures are consistent!')
    else:
        output.add_failure_message("Validation failed. Please fix the errors above.")

    output.print()
    # Only treat actual errors as failures, not warnings
    has_errors = bool(errors)
    return get_validation_exit_code(has_errors, no_fail)


if __name__ == '__main__':
    sys.exit(main())
