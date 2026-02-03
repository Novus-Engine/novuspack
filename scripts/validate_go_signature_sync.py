#!/usr/bin/env python3
"""
Check Go type, method, and function signatures in implementation against tech specs.

This script:
1. Extracts all public signatures from the Go implementation
2. Extracts all signatures from tech specs markdown files
3. Compares them and reports:
   - Errors: Signatures that are out of sync (different parameters, return types, etc.)
   - Warnings: Missing signatures in implementation (not yet implemented)
   - Warnings: Additional signatures in implementation not in spec (probably helpers)

Usage:
    python3 scripts/validate_go_signature_sync.py [options]

Options:
    --verbose, -v           Show detailed progress information
    --specs-dir DIR         Directory containing tech specs (default: docs/tech_specs)
    --impl-dir DIR           Directory containing Go implementation (default: api/go)
    --output, -o FILE        Output file path for validation report (default: stdout)
    --no-color, --nocolor    Disable colored output
    --help, -h               Show this help message
"""

import argparse
import functools
import re
import shutil
import subprocess  # nosec B404
import sys
from pathlib import Path
from typing import Dict, List, Set, Tuple

from lib._validation_utils import (
    OutputBuilder, get_workspace_root, parse_no_color_flag,
    ValidationIssue, DOCS_DIR, TECH_SPECS_DIR
)
from lib._go_code_utils import (
    parse_go_def_signature,
    find_go_code_blocks, Signature,
    normalize_go_signature_with_params,
    extract_interfaces_from_go_file,
    extract_interfaces_from_markdown
)
from lib._validate_go_signature_sync_helpers import (
    emit_extra_in_impl_section,
    emit_mismatches,
    emit_missing_in_impl,
    emit_sync_final,
)


_EMPTY_INTERFACE_RE = re.compile(r'\binterface\s*\{\s*\}')
_ANY_TYPE_RE = re.compile(r'\bany\b')
_RE_INTERFACE_PATTERN = re.compile(r'^\s*type\s+\w+(?:\s*\[[^\]]+\])?\s+interface\s*\{')

# Compiled regex patterns for type re-exports (module level)
_RE_TYPE_REEXPORT_PATTERN = re.compile(r'^\s+(\w+)\s*=\s+[\w.]+\.(\w+)')
_RE_STANDALONE_REEXPORT_PATTERN = re.compile(r'^\s*type\s+(\w+)\s*=\s+[\w.]+\.(\w+)')
_RE_TYPE_BLOCK_PATTERN = re.compile(r'^\s*type\s*\($')
_RE_TYPE_BLOCK_END_PATTERN = re.compile(r'^\s*\)$')
_RE_GENERIC_REEXPORT_IN_BLOCK_PATTERN = re.compile(
    r'^\s+(\w+)\s*\[[^\]]+\]\s*=\s+[\w.]+\.(\w+)\s*\['
)
_RE_GENERIC_REEXPORT_STANDALONE_PATTERN = re.compile(
    r'^\s*type\s+(\w+)\s*\[[^\]]+\]\s*=\s+[\w.]+\.(\w+)\s*\['
)


def _go_list_public_package_dirs(impl_dir: Path) -> List[Path]:
    """
    Return directories of packages that are part of the Go toolchain's public surface.

    Uses `go list ./...` from impl_dir so that:
    - Directories Go ignores (e.g. starting with `_`) are excluded.
    - Only importable packages are considered.
    Excludes packages whose import path contains `/internal/` so that the
    validator only reports on the externally visible API surface.

    Returns:
        List of absolute Paths to package directories.
    """
    impl_dir_abs = impl_dir.resolve()
    go_path = shutil.which('go')
    if not go_path:
        print(
            f"Warning: 'go' not found in PATH (impl_dir={impl_dir_abs})",
            file=sys.stderr,
        )
        return []
    try:
        result = subprocess.run(  # nosec B603
            [go_path, 'list', '-f', '{{.ImportPath}} {{.Dir}}', './...'],
            cwd=str(impl_dir_abs),
            capture_output=True,
            text=True,
            check=False,
            timeout=60,
        )
    except (FileNotFoundError, subprocess.TimeoutExpired) as e:
        print(
            f"Warning: could not run 'go list ./...' in {impl_dir_abs}: {e}",
            file=sys.stderr,
        )
        return []
    if result.returncode:
        print(
            f"Warning: 'go list ./...' failed in {impl_dir_abs}: {result.stderr}",
            file=sys.stderr,
        )
        return []
    dirs: List[Path] = []
    for line in result.stdout.strip().splitlines():
        line = line.strip()
        if not line:
            continue
        parts = line.split(' ', 1)
        if len(parts) != 2:
            continue
        import_path, pkg_dir = parts[0], parts[1]
        if '/internal/' in import_path or import_path.endswith('/internal'):
            continue
        dirs.append(Path(pkg_dir).resolve())
    return dirs


def _is_exported_receiver(sig: Signature) -> bool:
    """
    Return True if the signature is a function (no receiver) or a method on an exported type.

    Methods on unexported receiver types are not part of the public API surface.
    """
    if not sig.receiver:
        return True
    receiver_type = sig.receiver.strip('*').strip()
    if not receiver_type or '[' in receiver_type:
        base = receiver_type.split('[')[0].strip()
        return bool(base and base[0].isupper())
    return receiver_type[0].isupper()


def _is_public_surface_path(file_path: Path, impl_dir: Path) -> bool:
    """
    Return False if the file is under internal/ or a directory whose name starts with _.

    Matches Go visibility: internal packages and _-prefixed dirs (e.g. _bdd) are not public.
    """
    try:
        rel = file_path.resolve().relative_to(impl_dir.resolve())
    except ValueError:
        return True
    parts = rel.parts
    for part in parts:
        if part == 'internal' or part.startswith('_'):
            return False
    return True


def _gather_go_files(
    impl_dir: Path, public_dirs: List[Path], verbose: bool
) -> List[Path]:
    """
    Build list of non-test Go files from public package dirs or fallback to rglob.

    When go list succeeds, only files in public_dirs are used (Go toolchain view).
    When go list fails, fallback rglob still excludes internal/ and _-prefixed dirs
    so visibility matches Go semantics as much as possible without go list.
    """
    if public_dirs:
        go_files = []
        for pkg_dir in public_dirs:
            for f in pkg_dir.glob('*.go'):
                if not f.name.endswith('_test.go'):
                    go_files.append(f)
        if verbose:
            print(
                f"Scanning {len(go_files)} Go files in {len(public_dirs)} public package(s)..."
            )
        return go_files
    go_files = list(impl_dir.rglob('*.go'))
    go_files = [
        f for f in go_files
        if not f.name.endswith('_test.go') and _is_public_surface_path(f, impl_dir)
    ]
    if verbose:
        msg = (
            "Warning: 'go list' returned no public packages; "
            f"scanning {len(go_files)} Go files (excluding internal/ and _* dirs)."
        )
        print(msg)
    elif go_files:
        print(
            "Warning: 'go list ./...' failed or returned no packages; "
            "scanning Go files excluding internal/ and _* directories.",
            file=sys.stderr,
        )
    return go_files


def _process_one_impl_sig(
    sig: Signature,
    signatures: Dict[str, Signature],
    issues: List[ValidationIssue],
    verbose: bool,
) -> None:
    """Append issues for empty interface/any; add to signatures if public and exported."""
    if signature_has_empty_interface_input(sig):
        file_path, line_num = parse_location(sig.location)
        issues.append(ValidationIssue.create(
            "forbidden_empty_interface",
            file_path,
            line_num,
            line_num,
            message=(
                f"forbidden empty interface parameter type found in "
                f"implementation signature '{sig.normalized_key()}' at {sig.location}: "
                f"{sig.normalized_signature()}"
            ),
            severity='error',
            signature_key=sig.normalized_key(),
            location=sig.location
        ))
    if signature_uses_any_type(sig):
        file_path, line_num = parse_location(sig.location)
        issues.append(ValidationIssue.create(
            "discouraged_any_type",
            file_path,
            line_num,
            line_num,
            message=(
                f"discouraged any type usage found in "
                f"implementation signature '{sig.normalized_key()}' at {sig.location}: "
                f"{sig.normalized_signature()}"
            ),
            severity='warning',
            signature_key=sig.normalized_key(),
            location=sig.location
        ))
    if sig.is_public and _is_exported_receiver(sig):
        key = sig.normalized_key()
        if key in signatures:
            if verbose:
                print(f"    Warning: Duplicate signature {key} at {sig.location}")
        else:
            signatures[key] = sig


@functools.lru_cache(maxsize=1)
def _cached_reexported_types(repo_root_str: str) -> Set[str]:
    """Return re-exported types from novuspack.go (cached by repo root)."""
    return extract_reexported_types_from_novuspack(Path(repo_root_str))


def signature_has_empty_interface_input(sig: Signature) -> bool:
    """
    Return True if a function/method signature has an empty interface (interface{})
    anywhere in its parameter list.
    """
    if sig.kind not in ('func', 'method'):
        return False
    if not sig.params:
        return False
    return bool(_EMPTY_INTERFACE_RE.search(sig.params))


def signature_uses_any_type(sig: Signature) -> bool:
    """
    Return True if a function/method signature uses `any` in parameters or returns.

    This intentionally checks for concrete `any` usage (e.g., `dest any`,
    `map[string]any`, `[]*Tag[any]`).
    It does not attempt to detect type parameter constraints, since those are not
    part of Signature.params / Signature.returns in this validator.
    """
    if sig.kind not in ('func', 'method'):
        return False
    if _ANY_TYPE_RE.search(sig.params or ""):
        return True
    if _ANY_TYPE_RE.search(sig.returns or ""):
        return True
    return False


def extract_signatures_from_go_file(file_path: Path) -> List[Signature]:
    """Extract all public signatures from a Go file."""
    signatures = []

    try:
        # Resolve path to absolute to avoid ../ in output
        resolved_path = file_path.resolve()
        content = file_path.read_text(encoding='utf-8')
        lines = content.split('\n')

        # Use shared helper to extract interfaces and their methods
        interface_signatures = extract_interfaces_from_go_file(file_path, parse_methods=True)
        signatures.extend(interface_signatures)

        # Extract other signatures (functions, methods, types) that aren't interfaces
        for line_num, line in enumerate(lines, 1):
            stripped = line.strip()

            # Skip empty lines and comments
            if not stripped or stripped.startswith('//'):
                continue

            # Skip interface definitions (already handled by helper)
            if _RE_INTERFACE_PATTERN.match(stripped):
                continue

            # Check for any Go definition (function, method, or type)
            sig = parse_go_def_signature(line, location=f"{resolved_path}:{line_num}")
            if sig:
                if sig.kind in ('func', 'method'):
                    signatures.append(Signature(
                        name=sig.name,
                        kind=sig.kind,
                        receiver=sig.receiver,
                        params=sig.params,
                        returns=sig.returns,
                        location=f"{resolved_path}:{line_num}",
                        is_public=sig.is_public
                    ))
                elif sig.kind != 'interface':  # Interfaces already handled
                    signatures.append(Signature(
                        name=sig.name,
                        kind=sig.kind,
                        location=f"{resolved_path}:{line_num}",
                        is_public=sig.is_public
                    ))

    except (MemoryError, RuntimeError, BufferError) as e:
        print(f"Warning: Error reading {file_path}: {e}")

    return signatures


def extract_signatures_from_markdown_file(file_path: Path) -> List[Signature]:
    """Extract all signatures from Go code blocks in a markdown file."""
    signatures = []

    try:
        # Resolve path to absolute to avoid ../ in output
        resolved_path = file_path.resolve()
        content = file_path.read_text(encoding='utf-8')
        lines = content.split('\n')

        # Use shared helper to extract interfaces and their methods
        interface_signatures = extract_interfaces_from_markdown(
            content, file_path, start_line=1, parse_methods=True,
            skip_examples=False, lines=lines
        )
        signatures.extend(interface_signatures)

        # Extract other signatures (functions, methods, types) that aren't interfaces
        go_blocks = find_go_code_blocks(content)

        for start_line, _end_line, code_content in go_blocks:
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

                # Check for any Go definition (function, method, or type)
                sig = parse_go_def_signature(line, location=f"{resolved_path}:{line_num}")
                if sig:
                    if sig.kind in ('func', 'method'):
                        signatures.append(Signature(
                            name=sig.name,
                            kind=sig.kind,
                            receiver=sig.receiver,
                            params=sig.params,
                            returns=sig.returns,
                            location=f"{resolved_path}:{line_num}",
                            is_public=sig.is_public
                        ))
                    elif sig.kind != 'interface':  # Interfaces already handled
                        signatures.append(Signature(
                            name=sig.name,
                            kind=sig.kind,
                            location=f"{resolved_path}:{line_num}",
                            is_public=sig.is_public
                        ))

    except (IOError, OSError) as e:
        # File read errors - log to stderr
        print(f"Warning: Could not read {file_path}: {e}", file=sys.stderr)
    except UnicodeDecodeError as e:
        # Encoding errors - log to stderr
        print(f"Warning: Could not decode {file_path} (encoding issue): {e}", file=sys.stderr)
    except (MemoryError, RuntimeError, BufferError) as e:
        # Unexpected errors - log to stderr
        print(f"Warning: Unexpected error reading {file_path}: {e}", file=sys.stderr)

    return signatures


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


def collect_go_signatures(
    impl_dir: Path, verbose: bool = False
) -> Tuple[Dict[str, Signature], List[ValidationIssue]]:
    """
    Collect public signatures from Go implementation files.

    Only scans packages that are part of the Go toolchain's importable surface
    (via `go list ./...`), excluding `internal/` packages, and only counts
    methods on exported receiver types as public.
    """
    signatures = {}
    issues: List[ValidationIssue] = []
    public_dirs = _go_list_public_package_dirs(impl_dir)
    go_files = _gather_go_files(impl_dir, public_dirs, verbose)
    if verbose and public_dirs and not go_files:
        print(f"  (No non-test .go files in {len(public_dirs)} package dirs)")
    for go_file in go_files:
        if verbose:
            print(f"  Reading {go_file.relative_to(impl_dir)}")
        for sig in extract_signatures_from_go_file(go_file):
            _process_one_impl_sig(sig, signatures, issues, verbose)
    return signatures, issues


def collect_spec_signatures(
    specs_dir: Path, verbose: bool = False
) -> Tuple[Dict[str, Signature], List[ValidationIssue]]:
    """Collect all signatures from tech specs markdown files."""
    signatures = {}
    issues: List[ValidationIssue] = []

    # Find all .md files
    md_files = list(specs_dir.glob('*.md'))

    if verbose:
        print(f"Scanning {len(md_files)} markdown files...")

    for md_file in md_files:
        if verbose:
            print(f"  Reading {md_file.name}")

        file_sigs = extract_signatures_from_markdown_file(md_file)
        for sig in file_sigs:
            if signature_has_empty_interface_input(sig):
                file_path, line_num = parse_location(sig.location)
                issues.append(ValidationIssue.create(
                    "forbidden_empty_interface",
                    file_path,
                    line_num,
                    line_num,
                    message=(
                        f"forbidden empty interface parameter type found in "
                        f"spec signature '{sig.normalized_key()}' at {sig.location}: "
                        f"{sig.normalized_signature()}"
                    ),
                    severity='error',
                    signature_key=sig.normalized_key(),
                    location=sig.location
                ))
            if signature_uses_any_type(sig):
                file_path, line_num = parse_location(sig.location)
                issues.append(ValidationIssue.create(
                    "discouraged_any_type",
                    file_path,
                    line_num,
                    line_num,
                    message=(
                        f"discouraged any type usage found in "
                        f"spec signature '{sig.normalized_key()}' at {sig.location}: "
                        f"{sig.normalized_signature()}"
                    ),
                    severity='warning',
                    signature_key=sig.normalized_key(),
                    location=sig.location
                ))
            if sig.is_public:  # Only track public signatures
                key = sig.normalized_key()
                if key in signatures:
                    # Duplicate - keep the first one, but note the conflict
                    if verbose:
                        print(f"    Warning: Duplicate signature {key} at {sig.location}")
                else:
                    signatures[key] = sig

    return signatures, issues


def get_public_api_types() -> List[str]:
    """
    Get the list of public API types.

    This includes types re-exported from novuspack.go plus a fallback list.

    Returns:
        List of public API type names
    """
    repo_root = get_workspace_root()
    reexported = _cached_reexported_types(str(repo_root))
    public_types = list(reexported) if reexported else []

    # Add fallback types that might not be in novuspack.go but are still public API
    fallback_types = [
        'FileEntry', 'Package', 'PackageReader', 'PackageWriter',
        'ConfigBuilder', 'Tag', 'Option', 'PathMetadataEntry',
        'FileStream', 'BufferPool', 'ErrorType', 'PackageError',
        'FileInfo', 'PathInfo', 'AddFileOptions', 'ExtractPathOptions',
        'RemoveDirectoryOptions', 'ProcessingState', 'FileSource',
        'TransformPipeline', 'TransformStage', 'SignatureInfo',
        'SigningKey', 'CompressionStrategy', 'EncryptionStrategy',
        'PackageBuilder', 'PackageHeader', 'FileIndex', 'IndexEntry',
        'HashEntry', 'OptionalDataEntry', 'PackageComment', 'PackageInfo',
        'SecurityLevel', 'Signature', 'PathEntry', 'Result', 'Strategy',
        'Validator', 'ValidationRule', 'TagValueType', 'CreateOptions',
        'CompressionType', 'EncryptionType'
    ]

    # Merge and deduplicate (case-insensitive)
    type_set = set()
    for t in public_types + fallback_types:
        type_set.add(t)

    return list(type_set)


def extract_reexported_types_from_novuspack(root_dir: Path) -> Set[str]:
    """
    Extract all re-exported types from api/go/novuspack.go.

    This file re-exports types from subpackages, making them part of the public API.

    Args:
        root_dir: Root directory of the project (should contain api/go/)

    Returns:
        Set of type names that are re-exported (public API types)
    """
    reexported_types = set()
    novuspack_file = root_dir / "api" / "go" / "novuspack.go"

    if not novuspack_file.exists():
        return reexported_types

    try:
        content = novuspack_file.read_text(encoding='utf-8')
        lines = content.split('\n')

        # Pattern to match type re-exports:
        # type Name = package.Type
        # type ( Name1 = package.Type1 Name2 = package.Type2 )
        # Within type blocks, lines are indented: \tName = package.Type
        in_type_block = False

        for line in lines:
            # Check for start of type block
            if _RE_TYPE_BLOCK_PATTERN.match(line):
                in_type_block = True
                continue
            # Check for end of type block
            if in_type_block and _RE_TYPE_BLOCK_END_PATTERN.match(line):
                in_type_block = False
                continue

            # Match type re-exports within type blocks (indented lines)
            if in_type_block:
                match = _RE_TYPE_REEXPORT_PATTERN.match(line)
                if match:
                    exported_name = match.group(1)
                    original_name = match.group(2)
                    reexported_types.add(exported_name)
                    # If the exported name differs from original, add both
                    if exported_name != original_name:
                        reexported_types.add(original_name)
            else:
                # Match standalone type re-exports (not in blocks)
                match = _RE_STANDALONE_REEXPORT_PATTERN.match(line)
                if match:
                    exported_name = match.group(1)
                    original_name = match.group(2)
                    reexported_types.add(exported_name)
                    if exported_name != original_name:
                        reexported_types.add(original_name)

            # Also check for generic type re-exports:
            # type Option[T any] = generics.Option[T]
            # Within blocks: \tOption[T any] = generics.Option[T]
            if in_type_block:
                generic_match = _RE_GENERIC_REEXPORT_IN_BLOCK_PATTERN.match(line)
            else:
                generic_match = _RE_GENERIC_REEXPORT_STANDALONE_PATTERN.match(line)
            if generic_match:
                exported_name = generic_match.group(1)
                original_name = generic_match.group(2)
                reexported_types.add(exported_name)
                if exported_name != original_name:
                    reexported_types.add(original_name)

    except (IOError, OSError) as e:
        # File read errors - log to stderr
        print(f"Warning: Could not read novuspack.go for re-exports: {e}", file=sys.stderr)
    except UnicodeDecodeError as e:
        # Encoding errors - log to stderr
        print(f"Warning: Could not decode novuspack.go (encoding issue): {e}", file=sys.stderr)
    except (MemoryError, RuntimeError, BufferError) as e:
        # Unexpected errors - log to stderr
        print(
            f"Warning: Unexpected error parsing novuspack.go for re-exports: {e}",
            file=sys.stderr
        )

    return reexported_types


def _is_public_api_method(sig: Signature) -> bool:
    """Return True if sig is a method on a known public API type."""
    if not sig.receiver:
        return False
    receiver = sig.receiver.strip('*')
    if not receiver or not receiver[0].isupper():
        return False
    public_api_types = get_public_api_types()
    receiver_lower = receiver.lower()
    for api_type in public_api_types:
        if receiver_lower == api_type.lower():
            return True
    if '[' in receiver:
        base_type = receiver.split('[')[0].strip()
        return bool(base_type and base_type[0].isupper())
    return False


def _helper_score_from_path(
    file_name: str, parent_dir: str, location_lower: str, reasons: List[str]
) -> int:
    """Add path-based helper score and reasons; return score delta."""
    score = 0
    if 'helper' in file_name:
        score += 3
        reasons.append("filename contains 'helper'")
    if 'internal' in file_name:
        score += 3
        reasons.append("filename contains 'internal'")
    if 'test' in file_name and not file_name.endswith('_test.go'):
        score += 2
        reasons.append("filename contains 'test'")
    if 'helper' in parent_dir:
        score += 3
        reasons.append("parent directory contains 'helper'")
    if 'internal' in parent_dir:
        score += 3
        reasons.append("parent directory contains 'internal'")
    if '/internal/' in location_lower:
        score += 3
        reasons.append("in internal package path")
    if '/testhelpers/' in location_lower:
        score += 3
        reasons.append("in testhelpers package")
    if '/testutil/' in location_lower:
        score += 3
        reasons.append("in testutil package")
    if '/_bdd/' in location_lower:
        score += 2
        reasons.append("in BDD test package")
    return score


def _helper_score_from_name_receiver(sig: Signature, reasons: List[str]) -> int:
    """Add name/receiver-based helper score and reasons; return score delta."""
    helper_keywords = ['helper', 'internal', 'util', 'test', 'mock', 'stub']
    score = 0
    name_lower = sig.name.lower()
    for keyword in helper_keywords:
        if keyword in name_lower:
            score += 2
            reasons.append(f"name contains '{keyword}'")
    if sig.receiver:
        receiver_lower = sig.receiver.lower()
        for keyword in helper_keywords:
            if keyword in receiver_lower:
                score += 2
                reasons.append(f"receiver contains '{keyword}'")
        if receiver_lower.startswith('test') or receiver_lower.endswith('test'):
            score += 2
            reasons.append("receiver is test-related type")
    return score


def is_high_confidence_helper(sig: Signature) -> Tuple[bool, List[str]]:
    """
    Determine if a signature is a high-confidence helper function.

    This function is conservative: methods on public API types require
    very strong evidence to be considered helpers.

    Returns:
        - (is_helper, reasons): Tuple of boolean and list of reason strings
    """
    reasons: List[str] = []
    location_parts = sig.location.split(':')
    file_path = Path(location_parts[0])
    file_name = file_path.name.lower()
    parent_dir = file_path.parent.name.lower() if file_path.parent != file_path else ""
    if file_name.endswith('_test.go'):
        reasons.append("in test file (*_test.go)")
        return True, reasons

    is_public_api_method = _is_public_api_method(sig)
    helper_threshold = 6 if is_public_api_method else 3
    if is_public_api_method:
        reasons.append("method on public API type (requires stronger evidence)")

    score = _helper_score_from_path(
        file_name, parent_dir, sig.location.lower(), reasons
    )
    score += _helper_score_from_name_receiver(sig, reasons)
    return (score >= helper_threshold, reasons)


def compare_signatures(
    impl_sigs: Dict[str, Signature],
    spec_sigs: Dict[str, Signature]
) -> Tuple[List[Tuple[str, Signature, Signature]], List[str], List[str]]:
    """
    Compare implementation and spec signatures.

    Returns:
        - List of (name, impl_sig, spec_sig) tuples for mismatched signatures
        - List of names missing in implementation
        - List of names in implementation but not in spec
    """
    mismatches = []
    missing_in_impl = []
    extra_in_impl = []

    # Find mismatches and missing in implementation
    for key, spec_sig in spec_sigs.items():
        if key not in impl_sigs:
            missing_in_impl.append(key)
        else:
            impl_sig = impl_sigs[key]
            # Compare normalized signatures
            impl_norm = normalize_go_signature_with_params(impl_sig.normalized_signature())
            spec_norm = normalize_go_signature_with_params(spec_sig.normalized_signature())

            if impl_norm != spec_norm:
                mismatches.append((key, impl_sig, spec_sig))

    # Find extra in implementation
    for key in impl_sigs:
        if key not in spec_sigs:
            extra_in_impl.append(key)

    return mismatches, missing_in_impl, extra_in_impl


def _is_public_api_receiver(receiver: str, public_api_types: List[str]) -> bool:
    """Return True if receiver is a known public API type (or generic of one)."""
    if not receiver or not receiver[0].isupper():
        return False
    receiver_lower = receiver.lower()
    for api_type in public_api_types:
        if receiver_lower == api_type.lower():
            return True
    if '[' in receiver:
        base_type = receiver.split('[')[0].strip()
        if base_type and base_type[0].isupper():
            base_lower = base_type.lower()
            for api_type in public_api_types:
                if base_lower == api_type.lower():
                    return True
    return False


def _classify_extra_in_impl(
    extra_in_impl: List[str],
    impl_sigs: Dict[str, Signature],
    public_api_types: List[str],
) -> Tuple[
    List[Tuple[str, Signature, List[str]]],
    List[str],
    List[Tuple[str, Signature]],
]:
    """Classify extra implementation keys into helpers, low-confidence, or API errors."""
    high_confidence_helpers: List[Tuple[str, Signature, List[str]]] = []
    low_confidence_extra: List[str] = []
    errors_public_api_missing: List[Tuple[str, Signature]] = []
    for key in extra_in_impl:
        impl_sig = impl_sigs[key]
        is_helper, reasons = is_high_confidence_helper(impl_sig)
        if is_helper:
            high_confidence_helpers.append((key, impl_sig, reasons))
        elif impl_sig.receiver and _is_public_api_receiver(
            impl_sig.receiver.strip('*'), public_api_types
        ):
            errors_public_api_missing.append((key, impl_sig))
        else:
            low_confidence_extra.append(key)
    return high_confidence_helpers, low_confidence_extra, errors_public_api_missing


def main():
    """Entry point: parse args, check paths, run validation."""
    parser = argparse.ArgumentParser(description='Check Go signatures against tech specs')
    parser.add_argument('--verbose', '-v', action='store_true', help='Show detailed progress')
    parser.add_argument(
        '--specs-dir', type=str, default=f'{DOCS_DIR}/{TECH_SPECS_DIR}',
        help=f'Directory containing tech specs (default: {DOCS_DIR}/{TECH_SPECS_DIR})',
    )
    parser.add_argument(
        '--impl-dir', type=str, default='api/go',
        help='Directory containing Go implementation (default: api/go)',
    )
    parser.add_argument('--output', '-o', type=str, help='Output file path')
    parser.add_argument(
        '--no-color',
        '--nocolor',
        action='store_true',
        help='Disable colored output',
    )
    parser.add_argument('--no-fail', action='store_true', help='Exit 0 even if errors found')
    args = parser.parse_args()

    repo_root = get_workspace_root()
    specs_dir = repo_root / args.specs_dir
    impl_dir = repo_root / args.impl_dir
    if not specs_dir.exists():
        print(f"Error: Specs directory not found: {specs_dir}", file=sys.stderr)
        sys.exit(1)
    if not impl_dir.exists():
        print(f"Error: Implementation directory not found: {impl_dir}", file=sys.stderr)
        sys.exit(1)

    no_color = getattr(args, 'nocolor', False) or parse_no_color_flag(sys.argv)
    output = OutputBuilder(
        "Go Signature Sync Validation",
        "Checks Go signatures in implementation against tech specs",
        no_color=no_color,
        verbose=args.verbose,
        output_file=args.output,
    )
    _main_run(args, specs_dir, impl_dir, output, no_color)


def _main_run(
    args,
    specs_dir: Path,
    impl_dir: Path,
    output: OutputBuilder,
    no_color: bool,
) -> None:
    """Run validation after paths and output are set (collect, compare, emit)."""
    if args.verbose:
        output.add_verbose_line("Collecting signatures from Go implementation...")
    impl_sigs, impl_issues = collect_go_signatures(impl_dir, args.verbose)
    if args.verbose:
        output.add_verbose_line(f"Found {len(impl_sigs)} public signatures in implementation")
        output.add_blank_line("working_verbose")
        output.add_verbose_line("Collecting signatures from tech specs...")
    spec_sigs, spec_issues = collect_spec_signatures(specs_dir, args.verbose)
    if args.verbose:
        output.add_verbose_line(f"Found {len(spec_sigs)} public signatures in specs")
        output.add_blank_line("working_verbose")
        output.add_verbose_line("Comparing signatures...")

    all_issues = impl_issues + spec_issues
    empty_interface_errors = [
        i for i in all_issues
        if i.matches(issue_type='forbidden_empty_interface', severity='error')
    ]
    any_type_warnings = [
        i for i in all_issues
        if i.matches(issue_type='discouraged_any_type', severity='warning')
    ]
    if empty_interface_errors:
        for error in empty_interface_errors:
            output.add_error_line(error.format_message(no_color=no_color))
        output.print()
        sys.exit(1)
    if any_type_warnings:
        output.add_warnings_header()
        output.add_line(
            f"Found {len(any_type_warnings)} signature(s) using discouraged any type.",
            section="warning",
        )
        if args.verbose:
            for w in any_type_warnings:
                output.add_warning_line(w.format_message(no_color=no_color))
        else:
            for w in any_type_warnings[:25]:
                output.add_warning_line(w.format_message(no_color=no_color))
            if len(any_type_warnings) > 25:
                output.add_warning_line(
                    f"{len(any_type_warnings) - 25} additional warning(s) suppressed. "
                    "Re-run with --verbose to see all."
                )

    mismatches, missing_in_impl, extra_in_impl = compare_signatures(impl_sigs, spec_sigs)
    has_errors = emit_mismatches(output, mismatches, no_color=no_color)
    has_warnings = emit_missing_in_impl(
        output, missing_in_impl, spec_sigs, no_color=no_color
    )
    public_api_types = get_public_api_types()
    high_confidence_helpers, low_confidence_extra, errors_public_api_missing = (
        _classify_extra_in_impl(extra_in_impl, impl_sigs, public_api_types)
    )
    he_extra, hw_extra = emit_extra_in_impl_section(
        output, args,
        extra_in_impl=extra_in_impl,
        impl_sigs=impl_sigs,
        high_confidence_helpers=high_confidence_helpers,
        low_confidence_extra=low_confidence_extra,
        errors_public_api_missing=errors_public_api_missing,
        no_color=no_color,
    )
    has_errors = has_errors or he_extra
    has_warnings = has_warnings or hw_extra
    low_confidence_count = sum(
        1 for key in extra_in_impl if not is_high_confidence_helper(impl_sigs[key])[0]
    )
    high_confidence_count = len(extra_in_impl) - low_confidence_count
    emit_sync_final(
        output, args,
        has_errors=has_errors,
        has_warnings=has_warnings,
        impl_sigs=impl_sigs,
        spec_sigs=spec_sigs,
        mismatches=mismatches,
        missing_in_impl=missing_in_impl,
        extra_in_impl=extra_in_impl,
        low_confidence_count=low_confidence_count,
        high_confidence_count=high_confidence_count,
    )


if __name__ == '__main__':
    main()
