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
from typing import Callable, Dict, List, Optional, Set, Tuple

from lib._validation_utils import (
    OutputBuilder, find_markdown_files, parse_paths, get_workspace_root,
    get_validation_exit_code, HeadingContext, find_heading_before_line,
    ValidationIssue,
    DOCS_DIR, TECH_SPECS_DIR
)
from lib._go_code_utils import (
    is_example_signature_name,
    find_go_code_blocks, Signature,
    extract_interfaces_from_markdown,
    normalize_go_signature
)
from lib._validate_go_spec_signature_consistency_helpers import (
    extract_signatures_from_block as _extract_signatures_from_block,
    parse_cli_args,
    build_output,
    split_issues,
    emit_issues,
    emit_summary,
    emit_final_message,
)


def extract_signatures_from_markdown_file(file_path: Path, repo_root: Path) -> List[Signature]:
    """Extract all signatures from Go code blocks in a markdown file."""
    signatures = []

    try:
        try:
            relative_path = file_path.resolve().relative_to(repo_root.resolve())
        except ValueError:
            relative_path = file_path.resolve()
        content = file_path.read_text(encoding='utf-8')
        lines = content.split('\n')
        interface_signatures = extract_interfaces_from_markdown(
            content, file_path, start_line=1, parse_methods=True,
            skip_examples=True, lines=lines
        )
        signatures.extend(interface_signatures)
        go_blocks = find_go_code_blocks(content)
        for start_line, _end_line, code_content in go_blocks:
            block_lines = code_content.split('\n')
            signatures.extend(_extract_signatures_from_block(
                block_lines, start_line, relative_path, code_content, lines
            ))

    except (IOError, OSError) as e:
        # File read errors - log to stderr
        print(f"Warning: Could not read {file_path}: {e}", file=sys.stderr)
    except UnicodeDecodeError as e:
        # Encoding errors - log to stderr
        print(f"Warning: Could not decode {file_path} (encoding issue): {e}", file=sys.stderr)
    except (ValueError, KeyError, TypeError, AttributeError, RuntimeError) as e:
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
    except (TypeError, AttributeError, RuntimeError) as e:
        # Unexpected errors - log but don't fail
        print(f"Warning: Unexpected error reading {file_path}: {e}", file=sys.stderr)

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


def _canonical_score_heading(heading_ctx: Optional[HeadingContext]) -> Tuple[float, List[str]]:
    """Return (score_delta, reasons) for heading level."""
    if not heading_ctx:
        return (0.0, [])
    heading_score = max(0.1, 1.0 - (heading_ctx.heading_level - 1) * 0.2)
    return (heading_score * 0.4, [f"heading level {heading_ctx.heading_level}"])


def _canonical_score_body(sig: Signature) -> Tuple[float, List[str]]:
    """Return (score_delta, reasons) for has_body."""
    if sig.has_body:
        return (0.2, ["has body"])
    return (0.0, [])


def _canonical_score_first_occurrence(
    sig: Signature, all_sigs: List[Signature]
) -> Tuple[float, List[str]]:
    """Return (score_delta, reasons) for first occurrence in file."""
    file_path = sig.location.split(':', 1)[0]
    same_file_sigs = [s for s in all_sigs if s.location.startswith(file_path)]
    line_nums = []
    for s in same_file_sigs:
        try:
            line_nums.append(int(s.location.split(':', 1)[1]))
        except (ValueError, IndexError):
            pass
    if not line_nums:
        return (0.0, [])
    try:
        sig_line = int(sig.location.split(':', 1)[1])
        if sig_line == min(line_nums):
            return (0.15, ["first occurrence in file"])
    except (ValueError, IndexError):
        pass
    return (0.0, [])


def _canonical_score_file_name(file_path: str) -> Tuple[float, List[str]]:
    """Return (score_delta, reasons) for file name (core/basic vs advanced)."""
    file_name = file_path.lower()
    if 'core' in file_name or 'basic' in file_name:
        return (0.1, ["core/basic file"])
    if 'advanced' in file_name or 'extended' in file_name:
        return (-0.05, ["advanced/extended file"])
    return (0.0, [])


def _canonical_score_heading_keywords(
    sig: Signature, heading_ctx: Optional[HeadingContext]
) -> Tuple[float, List[str]]:
    """Return (score_delta, reasons) for heading keywords and name match."""
    if not heading_ctx:
        return (0.0, [])
    heading_lower = heading_ctx.heading_text.lower()
    sig_name_lower = sig.name.lower()
    score = 0.0
    reasons = []
    if sig_name_lower in heading_lower:
        score += 0.15
        reasons.append("signature name in heading")
    if 'definition' in heading_lower or 'definitions' in heading_lower:
        score += 0.1
        reasons.append("definition keyword in heading")
    if sig.kind in ('type', 'interface'):
        if any(kw in heading_lower for kw in ['type', 'struct', 'types', 'interfaces']):
            score += 0.1
            reasons.append("type-related keyword in heading")
    elif sig.kind == 'func':
        func_kw = ['function', 'functions', 'func', 'operation', 'operations']
        if any(kw in heading_lower for kw in func_kw):
            score += 0.1
            reasons.append("function-related keyword in heading")
    elif sig.kind == 'method' and sig.receiver:
        if any(kw in heading_lower for kw in ['method', 'methods']):
            score += 0.1
            reasons.append("method-related keyword in heading")
        if sig.receiver.lower() in heading_lower:
            score += 0.1
            reasons.append("receiver type in heading")
    return (score, reasons)


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
    reasons: List[str] = []
    s, r = _canonical_score_heading(heading_ctx)
    score += s
    reasons.extend(r)
    s, r = _canonical_score_body(sig)
    score += s
    reasons.extend(r)
    s, r = _canonical_score_first_occurrence(sig, all_sigs)
    score += s
    reasons.extend(r)
    s, r = _canonical_score_file_name(sig.location.split(':', 1)[0])
    score += s
    reasons.extend(r)
    s, r = _canonical_score_heading_keywords(sig, heading_ctx)
    score += s
    reasons.extend(r)
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


def _group_types_by_name(all_types: List[Signature]) -> Dict[str, List[Signature]]:
    """Group type/interface signatures by base name."""
    by_name: Dict[str, List[Signature]] = {}
    for sig in all_types:
        if sig.name not in by_name:
            by_name[sig.name] = []
        by_name[sig.name].append(sig)
    return by_name


def _append_duplicate_type_issue(
    name: str,
    sigs: List[Signature],
    get_first_location: Callable[[List[str]], Tuple[Path, int]],
    issues: List[ValidationIssue],
) -> None:
    """Append one duplicate_type issue for identical type/interface."""
    locations = [sig.location for sig in sigs]
    file_path, line_num = get_first_location(locations)
    message_parts = [f"Duplicate identical type/interface for '{name}':"]
    for sig in sigs:
        message_parts.append(f"  Location: {sig.location}")
    issues.append(ValidationIssue.create(
        "duplicate_type",
        file_path,
        line_num,
        line_num,
        message="\n".join(message_parts),
        severity='warning',
        type_name=name,
        locations=locations
    ))


def _append_type_conflict_issue(
    name: str,
    normalized_sigs: Dict[str, List[Signature]],
    get_first_location: Callable[[List[str]], Tuple[Path, int]],
    issues: List[ValidationIssue],
) -> None:
    """Append one conflicting_type_definitions issue."""
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
    issues.append(ValidationIssue.create(
        "conflicting_type_definitions",
        file_path,
        line_num,
        line_num,
        message="\n".join(message_parts),
        severity='error',
        type_name=name,
        locations=locations
    ))


def _build_normalized_type_sigs(sigs: List[Signature]) -> Dict[str, List[Signature]]:
    """Build dict of normalized signature string -> list of signatures (for types)."""
    normalized_sigs: Dict[str, List[Signature]] = {}
    for sig in sigs:
        sig_str = sig.name + (sig.generic_params or '')
        if sig_str not in normalized_sigs:
            normalized_sigs[sig_str] = []
        normalized_sigs[sig_str].append(sig)
    return normalized_sigs


def _append_stub_issues(
    stubs: List[Signature],
    canonical: Signature,
    name: str,
    parse_location: Callable[[str], Tuple[Path, int]],
    issues: List[ValidationIssue],
) -> None:
    """Append type_stub validation issues for stubs vs canonical."""
    for stub in stubs:
        if not stub.has_body or (
            canonical.has_body and canonical.method_count > stub.method_count
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
                'method_count' if canonical.kind == 'interface' else 'field_count'
            )
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
            issues.append(ValidationIssue.create(
                "type_stub",
                stub_file,
                stub_line,
                stub_line,
                message=message,
                severity='warning',
                type_name=name,
                canonical_location=canonical.location,
                stub_location=stub.location
            ))


def _append_method_inconsistency_if_stubs_differ(
    method_key: str,
    method_sigs: List[Signature],
    canonical: Signature,
    get_first_location: Callable[[List[str]], Tuple[Path, int]],
    issues: List[ValidationIssue],
) -> None:
    """If all stubs have different normalized signatures, append one issue."""
    normalized_methods: Dict[str, List[Signature]] = {}
    for method_sig in method_sigs:
        norm_sig = normalize_go_signature(method_sig.normalized_signature())
        if norm_sig not in normalized_methods:
            normalized_methods[norm_sig] = []
        normalized_methods[norm_sig].append(method_sig)
    if len(normalized_methods) <= 1:
        return
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
    issues.append(ValidationIssue.create(
        "method_signature_inconsistency",
        file_path,
        line_num,
        line_num,
        message="\n".join(message_parts),
        severity='error',
        method_name=method_key,
        interface_name=canonical.name,
        locations=locations
    ))


def _check_one_interface_method_consistency(
    method_key: str,
    method_sigs: List[Signature],
    canonical: Signature,
    *,
    get_first_location: Callable[[List[str]], Tuple[Path, int]],
    parse_location: Callable[[str], Tuple[Path, int]],
    issues: List[ValidationIssue],
) -> None:
    """Check one interface method's consistency (all stubs or canonical vs stubs)."""
    canonical_method = None
    stub_methods = []
    for method_sig in method_sigs:
        if method_sig.has_body:
            canonical_method = method_sig
        else:
            stub_methods.append(method_sig)
    if not canonical_method:
        _append_method_inconsistency_if_stubs_differ(
            method_key, method_sigs, canonical, get_first_location, issues
        )
        return
    for stub_method in stub_methods:
        canonical_norm = normalize_go_signature(canonical_method.normalized_signature())
        stub_norm = normalize_go_signature(stub_method.normalized_signature())
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
            issues.append(ValidationIssue.create(
                "interface_stub_method_differs",
                stub_file,
                stub_line,
                stub_line,
                message=message,
                severity='error',
                method_name=method_key,
                canonical_location=canonical_method.location,
                stub_location=stub_method.location
            ))


def _collect_canonical_interface_method_names(
    canonical: Signature, signatures: Dict[str, List[Signature]]
) -> set:
    """Return set of method names defined in canonical interface block (same file, ~100 lines)."""
    canonical_file, canonical_line_str = canonical.location.split(':', 1)
    try:
        canonical_line = int(canonical_line_str)
    except ValueError:
        canonical_line = 0
    names = set()
    for sig_list in signatures.values():
        for sig in sig_list:
            if not (sig.kind == 'method' and sig.receiver == canonical.name and not sig.has_body):
                continue
            sig_file, sig_line_str = sig.location.split(':', 1)
            try:
                sig_line = int(sig_line_str)
                line_ok = canonical_line <= sig_line <= canonical_line + 100
                if sig_file == canonical_file and line_ok:
                    names.add(sig.name)
            except ValueError:
                pass
    return names


def _build_interface_methods_map(
    canonical: Signature, signatures: Dict[str, List[Signature]]
) -> Dict[str, List[Signature]]:
    """Return map method_name -> list of Signature for methods on canonical.name."""
    out: Dict[str, List[Signature]] = {}
    for sig_list in signatures.values():
        for sig in sig_list:
            if sig.kind == 'method' and sig.receiver == canonical.name:
                out.setdefault(sig.name, []).append(sig)
    return out


def _check_interface_method_consistency(
    canonical: Signature,
    signatures: Dict[str, List[Signature]],
    get_first_location: Callable[[List[str]], Tuple[Path, int]],
    parse_location: Callable[[str], Tuple[Path, int]],
    issues: List[ValidationIssue],
) -> None:
    """Check interface method consistency (canonical vs stubs, method not in canonical)."""
    canonical_method_names = _collect_canonical_interface_method_names(canonical, signatures)
    interface_methods = _build_interface_methods_map(canonical, signatures)
    if canonical_method_names and canonical.has_body:
        for method_key, method_sigs in interface_methods.items():
            if method_key not in canonical_method_names:
                locations = [msig.location for msig in method_sigs]
                file_path, line_num = get_first_location(locations)
                message_parts = [
                    f"Method '{method_key}' is defined for interface '{canonical.name}' "
                    "but is not in the canonical interface definition:",
                    f"  Canonical interface: {canonical.location}",
                    f"  Methods in canonical: {sorted(canonical_method_names)}",
                ]
                for method_sig in method_sigs:
                    message_parts.append(f"  Method location: {method_sig.location}")
                issues.append(ValidationIssue.create(
                    "method_not_in_canonical_interface",
                    file_path, line_num, line_num,
                    message="\n".join(message_parts),
                    severity='error',
                    method_name=method_key,
                    interface_name=canonical.name,
                    canonical_location=canonical.location,
                    locations=locations,
                ))
    for method_key, method_sigs in interface_methods.items():
        if len(method_sigs) <= 1:
            continue
        _check_one_interface_method_consistency(
            method_key, method_sigs, canonical,
            get_first_location=get_first_location,
            parse_location=parse_location,
            issues=issues,
        )


def _process_method_group(
    key: str,
    sig_list: List[Signature],
    get_first_location: Callable[[List[str]], Tuple[Path, int]],
    issues: List[ValidationIssue],
) -> None:
    """Check one method/function group for inconsistency or duplicate."""
    methods = [s for s in sig_list if s.kind in ('method', 'func')]
    if not methods:
        return
    normalized_sigs: Dict[str, List[Signature]] = {}
    for sig in methods:
        norm_sig = normalize_go_signature(sig.normalized_signature())
        if norm_sig not in normalized_sigs:
            normalized_sigs[norm_sig] = []
        normalized_sigs[norm_sig].append(sig)
    if len(normalized_sigs) > 1:
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
        issues.append(ValidationIssue.create(
            "signature_inconsistency",
            file_path,
            line_num,
            line_num,
            message="\n".join(message_parts),
            severity='error',
            signature_key=key,
            locations=locations
        ))
    elif len(methods) > 1:
        locations = [sig.location for sig in methods]
        file_path, line_num = get_first_location(locations)
        message_parts = [f"Duplicate identical signature for '{key}':"]
        for sig in methods:
            message_parts.append(f"  Location: {sig.location}")
        issues.append(ValidationIssue.create(
            "duplicate_signature",
            file_path,
            line_num,
            line_num,
            message="\n".join(message_parts),
            severity='error',
            signature_key=key,
            locations=locations
        ))


def check_signature_consistency(
    signatures: Dict[str, List[Signature]],
    _verbose: bool = False,
    _repo_root: Optional[Path] = None
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

    if all_types:
        by_name = _group_types_by_name(all_types)
        for name, sigs in by_name.items():
            if len(sigs) <= 1:
                continue
            normalized_sigs = _build_normalized_type_sigs(sigs)
            if len(normalized_sigs) > 1:
                _append_type_conflict_issue(
                    name, normalized_sigs, get_first_location, issues
                )
                continue

            # If same normalized signature, check for duplicates or stubs
            canonical = find_canonical_definition(sigs)
            if not canonical:
                # No canonical found, but we have multiple identical signatures
                # This is a duplicate - will be handled by the duplicate_type check below
                continue

            stubs = [s for s in sigs if s != canonical]

            if not stubs and len(sigs) > 1:
                _append_duplicate_type_issue(name, sigs, get_first_location, issues)
            _append_stub_issues(stubs, canonical, name, parse_location, issues)

            if canonical.kind == 'interface':
                _check_interface_method_consistency(
                    canonical, signatures, get_first_location, parse_location, issues
                )

        for key, sig_list in signatures.items():
            if len(sig_list) <= 1:
                continue
            _process_method_group(key, sig_list, get_first_location, issues)

    return issues


def _reported_keys_from_issues(issues: List[ValidationIssue]) -> Set[str]:
    """Extract signature/type/method keys already reported in issues."""
    reported = set()
    for issue in issues:
        if not isinstance(issue, ValidationIssue):
            continue
        key = (
            issue.extra_fields.get('signature_key') or
            issue.extra_fields.get('type_name') or
            issue.extra_fields.get('method_name')
        )
        if key:
            reported.add(key)
        else:
            match = re.search(r"for '([^']+)'|'([^']+)'", issue.message)
            if match:
                reported.add(match.group(1) or match.group(2))
    return reported


def _append_unreported_duplicate_issue(
    key: str,
    sig_list: List[Signature],
    repo_root: Path,
    issues: List[ValidationIssue],
) -> None:
    """Append one duplicate_signature issue for key/sig_list if unreported."""
    canonical_sig, _ctx, confidence, reason = find_canonical_signature(
        sig_list, repo_root
    )
    locations = [sig.location for sig in sig_list]
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
    if canonical_sig and confidence >= 0.7:
        message_parts = [
            f"Duplicate identical signature for '{key}':",
            f"  Suggested canonical (confidence: {confidence:.0%}): "
            f"{canonical_sig.location} (reason: {reason})",
            "  Other locations:"
        ]
        other_locations = [
            f"    {sig.location}" for sig in sig_list if sig != canonical_sig
        ]
        message_parts.extend(other_locations)
        suggestion = f"Use canonical: {canonical_sig.location}"
    else:
        message_parts = [f"Duplicate identical signature for '{key}':"]
        all_locations = [f"  Location: {sig.location}" for sig in sig_list]
        message_parts.extend(all_locations)
        if canonical_sig:
            message_parts.append(
                f"  Note: {canonical_sig.location} may be canonical "
                f"(low confidence: {confidence:.0%}, reason: {reason})"
            )
        suggestion = None
    issues.append(ValidationIssue.create(
        "duplicate_signature",
        file_path,
        line_num,
        line_num,
        message="\n".join(message_parts),
        severity='error',
        suggestion=suggestion,
        signature_key=key,
        locations=locations,
        canonical_location=canonical_sig.location if canonical_sig else None
    ))


def _collect_signatures(
    md_files: List[Path],
    repo_root: Path,
    output: OutputBuilder,
    verbose: bool,
) -> List[Signature]:
    all_signatures: List[Signature] = []
    for md_file in md_files:
        if verbose:
            output.add_verbose_line(f'Extracting signatures from {md_file.name}...')
        file_sigs = extract_signatures_from_markdown_file(md_file, repo_root)
        all_signatures.extend(file_sigs)
    return all_signatures


def _group_signatures(
    all_signatures: List[Signature],
) -> Dict[str, List[Signature]]:
    signatures_by_key: Dict[str, List[Signature]] = defaultdict(list)
    for sig in all_signatures:
        if is_example_signature_name(sig.name):
            continue
        key = sig.normalized_key()
        signatures_by_key[key].append(sig)
    return signatures_by_key


def _collect_unreported_duplicates(
    all_signatures: List[Signature],
    signatures_by_key: Dict[str, List[Signature]],
    repo_root: Path,
    issues: List[ValidationIssue],
) -> None:
    duplicate_count = len(all_signatures) - len(signatures_by_key)
    if duplicate_count <= 0:
        return
    reported_keys = _reported_keys_from_issues(issues)
    for key, sig_list in signatures_by_key.items():
        sig_list = [sig for sig in sig_list if not is_example_signature_name(sig.name)]
        if len(sig_list) < 2 or key in reported_keys:
            continue
        _append_unreported_duplicate_issue(key, sig_list, repo_root, issues)


def main():
    """Main entry point."""
    # Show help if requested
    if '--help' in sys.argv or '-h' in sys.argv:
        print(__doc__)
        return 0

    verbose, no_color, no_fail, output_file, target_paths_str = parse_cli_args(
        sys.argv
    )

    # Parse comma-separated paths
    target_paths = parse_paths(target_paths_str)

    # Create output builder (header streams immediately if verbose)
    output = build_output(verbose, no_color, output_file)

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
    all_signatures = _collect_signatures(md_files, repo_root, output, verbose)

    if verbose:
        output.add_verbose_line(f'Found {len(all_signatures)} total signatures')
        output.add_blank_line("working_verbose")

    # Group signatures by normalized key, filtering out examples
    signatures_by_key = _group_signatures(all_signatures)

    # Check for inconsistencies
    if verbose:
        output.add_verbose_line('Checking for signature inconsistencies...')
        output.add_blank_line("working_verbose")

    issues = check_signature_consistency(
        signatures_by_key, verbose, _repo_root=repo_root
    )

    # Check if signature counts don't match (indicates duplicates that weren't caught)
    _collect_unreported_duplicates(
        all_signatures, signatures_by_key, repo_root, issues
    )

    # Report results (headers will be added automatically by add_error_line/add_warning_line)
    # Filter issues by severity in a single loop
    errors, warnings = split_issues(issues)

    emit_issues(output, errors, warnings, no_color)

    # Always add summary section
    emit_summary(output, all_signatures, signatures_by_key, errors, warnings)

    # Final message: success if no errors, but mention warnings if present
    emit_final_message(output, errors, warnings)

    output.print()
    # Only treat actual errors as failures, not warnings
    has_errors = bool(errors)
    return get_validation_exit_code(has_errors, no_fail)


if __name__ == '__main__':
    sys.exit(main())
