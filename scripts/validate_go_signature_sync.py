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
    --help, -h               Show this help message
"""

import argparse
import re
import sys
from pathlib import Path
from typing import Dict, List, Optional, Set, Tuple

scripts_dir = Path(__file__).parent
lib_dir = scripts_dir / "lib"

# Import shared utilities
for module_path in (str(scripts_dir), str(lib_dir)):
    if module_path not in sys.path:
        sys.path.insert(0, module_path)

from lib._validation_utils import (  # noqa: E402
    OutputBuilder, get_workspace_root, parse_no_color_flag,
    ValidationIssue, DOCS_DIR, TECH_SPECS_DIR
)
from lib._go_code_utils import (  # noqa: E402
    parse_go_def_signature,
    find_go_code_blocks, Signature,
    normalize_go_signature_with_params,
    extract_interfaces_from_go_file,
    extract_interfaces_from_markdown
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

# Module-level cache for re-exported types from novuspack.go
_REEXPORTED_TYPES: Optional[Set[str]] = None


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

    except Exception as e:
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
    except Exception as e:
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
    """Collect all public signatures from Go implementation files."""
    signatures = {}
    issues: List[ValidationIssue] = []

    # Find all .go files, excluding test files
    go_files = list(impl_dir.rglob('*.go'))
    go_files = [f for f in go_files if not f.name.endswith('_test.go')]

    if verbose:
        print(f"Scanning {len(go_files)} Go files...")

    for go_file in go_files:
        if verbose:
            print(f"  Reading {go_file.relative_to(impl_dir)}")

        file_sigs = extract_signatures_from_go_file(go_file)
        for sig in file_sigs:
            if signature_has_empty_interface_input(sig):
                file_path, line_num = parse_location(sig.location)
                issues.append(ValidationIssue(
                    "forbidden_empty_interface",
                    file_path,
                    line_num,
                    line_num,
                    f"forbidden empty interface parameter type found in "
                    f"implementation signature '{sig.normalized_key()}' at {sig.location}: "
                    f"{sig.normalized_signature()}",
                    severity='error',
                    signature_key=sig.normalized_key(),
                    location=sig.location
                ))
            if signature_uses_any_type(sig):
                file_path, line_num = parse_location(sig.location)
                issues.append(ValidationIssue(
                    "discouraged_any_type",
                    file_path,
                    line_num,
                    line_num,
                    f"discouraged any type usage found in "
                    f"implementation signature '{sig.normalized_key()}' at {sig.location}: "
                    f"{sig.normalized_signature()}",
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
                issues.append(ValidationIssue(
                    "forbidden_empty_interface",
                    file_path,
                    line_num,
                    line_num,
                    f"forbidden empty interface parameter type found in "
                    f"spec signature '{sig.normalized_key()}' at {sig.location}: "
                    f"{sig.normalized_signature()}",
                    severity='error',
                    signature_key=sig.normalized_key(),
                    location=sig.location
                ))
            if signature_uses_any_type(sig):
                file_path, line_num = parse_location(sig.location)
                issues.append(ValidationIssue(
                    "discouraged_any_type",
                    file_path,
                    line_num,
                    line_num,
                    f"discouraged any type usage found in "
                    f"spec signature '{sig.normalized_key()}' at {sig.location}: "
                    f"{sig.normalized_signature()}",
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
    global _REEXPORTED_TYPES

    # If we haven't loaded re-exported types yet, try to load them
    if _REEXPORTED_TYPES is None:
        repo_root = get_workspace_root()
        _REEXPORTED_TYPES = extract_reexported_types_from_novuspack(repo_root)

    # Start with re-exported types
    public_types = list(_REEXPORTED_TYPES) if _REEXPORTED_TYPES else []

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
    except Exception as e:
        # Unexpected errors - log to stderr
        print(
            f"Warning: Unexpected error parsing novuspack.go for re-exports: {e}",
            file=sys.stderr
        )

    return reexported_types


def is_high_confidence_helper(sig: Signature) -> Tuple[bool, List[str]]:
    """
    Determine if a signature is a high-confidence helper function.

    This function is conservative: methods on public API types require
    very strong evidence to be considered helpers.

    Returns:
        - (is_helper, reasons): Tuple of boolean and list of reason strings
    """
    reasons = []
    score = 0

    # Parse file path from location (format: "path/to/file.go:123")
    location_parts = sig.location.split(':')
    file_path_str = location_parts[0]
    file_path = Path(file_path_str)
    file_name = file_path.name.lower()
    parent_dir = file_path.parent.name.lower() if file_path.parent != file_path else ""

    # Check if file is a test file (*_test.go)
    # This is the strongest indicator - test files are always helpers
    if file_name.endswith('_test.go'):
        reasons.append("in test file (*_test.go)")
        # Test files are always helpers, return early
        return True, reasons

    # Check if this is a method on a public API type
    # Public API types are those that start with capital letter and are
    # likely part of the documented API surface
    is_public_api_method = False
    if sig.receiver:
        receiver = sig.receiver.strip('*')  # Remove pointer indicator
        # Check if receiver is a public type (starts with capital)
        if receiver and receiver[0].isupper():
            # Public API types are determined by what's re-exported in novuspack.go
            # This ensures we only flag methods on truly public API types as errors
            public_api_types = get_public_api_types()
            # Check if receiver matches a known API type (case-insensitive)
            receiver_lower = receiver.lower()
            for api_type in public_api_types:
                if receiver_lower == api_type.lower():
                    is_public_api_method = True
                    break
            # Also check if receiver looks like a generic type (contains [)
            if '[' in receiver:
                base_type = receiver.split('[')[0].strip()
                if base_type and base_type[0].isupper():
                    is_public_api_method = True

    # If this is a method on a public API type, require much stronger evidence
    # to be considered a helper
    if is_public_api_method:
        # Methods on public API types need score >= 6 to be considered helpers
        # (vs normal threshold of 3)
        helper_threshold = 6
        reasons.append("method on public API type (requires stronger evidence)")
    else:
        # Normal threshold for non-API methods
        helper_threshold = 3

    # Check filename patterns
    if 'helper' in file_name:
        score += 3
        reasons.append("filename contains 'helper'")
    if 'internal' in file_name:
        score += 3
        reasons.append("filename contains 'internal'")
    if 'test' in file_name and not file_name.endswith('_test.go'):
        # Test in filename but not a test file pattern
        score += 2
        reasons.append("filename contains 'test'")

    # Check parent directory patterns
    if 'helper' in parent_dir:
        score += 3
        reasons.append("parent directory contains 'helper'")
    if 'internal' in parent_dir:
        score += 3
        reasons.append("parent directory contains 'internal'")

    # Check package name in full file path
    location_lower = sig.location.lower()
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

    # Check function/method name patterns
    name_lower = sig.name.lower()
    helper_keywords = ['helper', 'internal', 'util', 'test', 'mock', 'stub']
    for keyword in helper_keywords:
        if keyword in name_lower:
            score += 2
            reasons.append(f"name contains '{keyword}'")

    # Check receiver type (for methods)
    if sig.receiver:
        receiver_lower = sig.receiver.lower()
        for keyword in helper_keywords:
            if keyword in receiver_lower:
                score += 2
                reasons.append(f"receiver contains '{keyword}'")
        # Check for test-related receiver types
        if receiver_lower.startswith('test') or receiver_lower.endswith('test'):
            score += 2
            reasons.append("receiver is test-related type")

    # Use threshold based on whether this is a public API method
    is_helper = score >= helper_threshold

    return is_helper, reasons


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


def main():
    parser = argparse.ArgumentParser(
        description='Check Go signatures against tech specs'
    )
    parser.add_argument(
        '--verbose', '-v',
        action='store_true',
        help='Show detailed progress information'
    )
    parser.add_argument(
        '--specs-dir',
        type=str,
        default=f'{DOCS_DIR}/{TECH_SPECS_DIR}',
        help=f'Directory containing tech specs (default: {DOCS_DIR}/{TECH_SPECS_DIR})'
    )
    parser.add_argument(
        '--impl-dir',
        type=str,
        default='api/go',
        help='Directory containing Go implementation (default: api/go)'
    )
    parser.add_argument(
        '--output', '-o',
        type=str,
        help='Output file path for validation report (default: stdout)'
    )
    parser.add_argument(
        '--no-fail',
        action='store_true',
        help='Exit with code 0 even if errors are found'
    )

    args = parser.parse_args()

    # Determine paths
    repo_root = get_workspace_root()
    specs_dir = repo_root / args.specs_dir
    impl_dir = repo_root / args.impl_dir

    if not specs_dir.exists():
        print(f"Error: Specs directory not found: {specs_dir}", file=sys.stderr)
        sys.exit(1)

    if not impl_dir.exists():
        print(f"Error: Implementation directory not found: {impl_dir}", file=sys.stderr)
        sys.exit(1)

    # Create output builder (header streams immediately if verbose)
    no_color = getattr(args, 'nocolor', False) or parse_no_color_flag(sys.argv)
    output = OutputBuilder(
        "Go Signature Sync Validation",
        "Checks Go signatures in implementation against tech specs",
        no_color=no_color,
        verbose=args.verbose,
        output_file=args.output
    )

    # Collect signatures
    if args.verbose:
        output.add_verbose_line("Collecting signatures from Go implementation...")
    impl_sigs, impl_issues = collect_go_signatures(impl_dir, args.verbose)

    if args.verbose:
        output.add_verbose_line(f"Found {len(impl_sigs)} public signatures in implementation")
        output.add_blank_line("working_verbose")
        output.add_verbose_line("Collecting signatures from tech specs...")
    spec_sigs, spec_issues = collect_spec_signatures(
        specs_dir, args.verbose
    )

    if args.verbose:
        output.add_verbose_line(f"Found {len(spec_sigs)} public signatures in specs")
        output.add_blank_line("working_verbose")
        output.add_verbose_line("Comparing signatures...")

    # Combine all issues and filter in a single loop
    all_issues = impl_issues + spec_issues
    empty_interface_errors = []
    any_type_warnings = []
    for issue in all_issues:
        if issue.matches(issue_type='forbidden_empty_interface', severity='error'):
            empty_interface_errors.append(issue)
        if issue.matches(issue_type='discouraged_any_type', severity='warning'):
            any_type_warnings.append(issue)

    if empty_interface_errors:
        for error in empty_interface_errors:
            output.add_error_line(error.format_message(no_color=no_color))
        output.print()
        sys.exit(1)
    if any_type_warnings:
        output.add_warnings_header()
        output.add_line(
            f"Found {len(any_type_warnings)} signature(s) using discouraged any type.",
            section="warning"
        )
        if args.verbose:
            for warning in any_type_warnings:
                output.add_warning_line(warning.format_message(no_color=no_color))
        else:
            max_show = 25
            for warning in any_type_warnings[:max_show]:
                output.add_warning_line(warning.format_message(no_color=no_color))
            suppressed = len(any_type_warnings) - max_show
            if suppressed > 0:
                output.add_warning_line(
                    f"{suppressed} additional warning(s) suppressed. "
                    "Re-run with --verbose to see all."
                )

    # Compare
    mismatches, missing_in_impl, extra_in_impl = compare_signatures(impl_sigs, spec_sigs)

    # Report results
    has_errors = False
    has_warnings = False

    # Errors: Mismatched signatures
    if mismatches:
        has_errors = True
        output.add_errors_header()
        output.add_line(
            f"Found {len(mismatches)} signature mismatch(es):",
            section="error"
        )
        output.add_blank_line("error")

        for key, impl_sig, spec_sig in sorted(mismatches):
            output.add_error_line(f"Signature: {key}")
            output.add_error_line(f"  Implementation: {impl_sig.normalized_signature()}")
            output.add_error_line(f"    Location: {impl_sig.location}")
            output.add_error_line(f"  Specification:  {spec_sig.normalized_signature()}")
            output.add_error_line(f"    Location: {spec_sig.location}")

    # Warnings: Missing in implementation
    if missing_in_impl:
        has_warnings = True
        output.add_warnings_header()
        output.add_line(
            f"Found {len(missing_in_impl)} signature(s) in specs "
            f"but not in implementation:",
            section="warning"
        )
        output.add_blank_line("warning")

        for key in sorted(missing_in_impl):
            spec_sig = spec_sigs[key]
            output.add_warning_line(f"  {key}")
            output.add_warning_line(f"    Signature: {spec_sig.normalized_signature()}")
            output.add_warning_line(f"    Location: {spec_sig.location}")

    # Errors/Warnings: Extra in implementation (filter out high-confidence helpers)
    if extra_in_impl:
        high_confidence_helpers = []
        low_confidence_extra = []
        errors_public_api_missing = []

        for key in extra_in_impl:
            impl_sig = impl_sigs[key]
            is_helper, reasons = is_high_confidence_helper(impl_sig)
            if is_helper:
                high_confidence_helpers.append((key, impl_sig, reasons))
            else:
                # Check if this is a public method on a public API type
                # These should be errors, not warnings
                is_public_api_error = False
                if impl_sig.receiver:
                    receiver = impl_sig.receiver.strip('*')
                    if receiver and receiver[0].isupper():
                        # Known public API types
                        public_api_types = [
                            'FileEntry', 'Package', 'PackageReader', 'PackageWriter',
                            'ConfigBuilder', 'Tag', 'Option', 'PathMetadataEntry',
                            'FileStream', 'BufferPool', 'ErrorType', 'PackageError',
                            'FileInfo', 'PathInfo', 'AddFileOptions', 'ExtractPathOptions',
                            'RemoveDirectoryOptions', 'ProcessingState', 'FileSource',
                            'TransformPipeline', 'TransformStage', 'SignatureInfo',
                            'SigningKey', 'CompressionStrategy', 'EncryptionStrategy'
                        ]
                        receiver_lower = receiver.lower()
                        for api_type in public_api_types:
                            if receiver_lower == api_type.lower():
                                is_public_api_error = True
                                break
                        # Also check generic types
                        if '[' in receiver:
                            base_type = receiver.split('[')[0].strip()
                            if base_type and base_type[0].isupper():
                                base_lower = base_type.lower()
                                for api_type in public_api_types:
                                    if base_lower == api_type.lower():
                                        is_public_api_error = True
                                        break

                if is_public_api_error:
                    errors_public_api_missing.append((key, impl_sig))
                else:
                    low_confidence_extra.append(key)

        # Report suppressed helpers
        if high_confidence_helpers:
            if args.verbose:
                output.add_verbose_line(
                    f"Suppressed {len(high_confidence_helpers)} high-confidence "
                    f"helper function(s)"
                )
                output.add_blank_line("working_verbose")
                for key, impl_sig, reasons in sorted(high_confidence_helpers):
                    output.add_verbose_line(f"  {key}")
                    output.add_verbose_line(f"    Signature: {impl_sig.normalized_signature()}")
                    output.add_verbose_line(f"    Location: {impl_sig.location}")
                    output.add_verbose_line(f"    Reasons: {', '.join(reasons)}")

        # Report errors: Public API methods missing from specs
        if errors_public_api_missing:
            has_errors = True
            output.add_errors_header()
            output.add_line(
                f"Found {len(errors_public_api_missing)} public API method(s) "
                f"in implementation but not in specs:",
                section="error"
            )
            output.add_blank_line("error")
            output.add_line(
                "These are public methods on public API types and MUST be documented "
                "in tech specs.",
                section="error"
            )
            output.add_blank_line("error")

            for key, impl_sig in sorted(errors_public_api_missing):
                output.add_error_line(f"Signature: {key}")
                output.add_error_line(f"  Implementation: {impl_sig.normalized_signature()}")
                output.add_error_line(f"    Location: {impl_sig.location}")

        # Report low-confidence extra functions
        if low_confidence_extra:
            has_warnings = True
            output.add_warnings_header()
            output.add_line(
                f"Found {len(low_confidence_extra)} signature(s) in "
                f"implementation but not in specs:",
                section="warning"
            )
            if high_confidence_helpers:
                if args.verbose:
                    output.add_line(
                        f"(Suppressed {len(high_confidence_helpers)} high-confidence "
                        f"helper function(s) - see above)",
                        section="warning"
                    )
                else:
                    output.add_line(
                        f"(Suppressed {len(high_confidence_helpers)} high-confidence "
                        f"helper function(s) - use --verbose to see them)",
                        section="warning"
                    )
            if errors_public_api_missing:
                output.add_line(
                    f"(Also found {len(errors_public_api_missing)} public API method(s) "
                    f"missing from specs - see errors above)",
                    section="warning"
                )
            output.add_blank_line("warning")
            output.add_line(
                "(These may be helper functions, but should be checked)",
                section="warning"
            )
            output.add_blank_line("warning")

            for key in sorted(low_confidence_extra):
                impl_sig = impl_sigs[key]
                output.add_warning_line(f"  {key}")
                output.add_warning_line(f"    Signature: {impl_sig.normalized_signature()}")
                output.add_warning_line(f"    Location: {impl_sig.location}")
        elif high_confidence_helpers and not args.verbose:
            # Only helpers, no low-confidence extras, and not verbose
            output.add_verbose_line(
                f"Suppressed {len(high_confidence_helpers)} high-confidence "
                f"helper function(s)"
            )
            output.add_verbose_line("(Use --verbose to see the list of suppressed helpers)")

    # Summary
    if not has_errors and not has_warnings:
        output.add_success_message("All signatures are in sync!")
        if args.verbose:
            output.add_verbose_line(f"  - {len(impl_sigs)} signatures in implementation")
            output.add_verbose_line(f"  - {len(spec_sigs)} signatures in specs")
        output.print()
        sys.exit(0)
    else:
        summary_parts = []
        if mismatches:
            summary_parts.append(f"{len(mismatches)} mismatch(es)")
        if missing_in_impl:
            summary_parts.append(f"{len(missing_in_impl)} missing in implementation")
        if extra_in_impl:
            # Count only low-confidence extras in summary
            low_confidence_count = sum(
                1 for key in extra_in_impl
                if not is_high_confidence_helper(impl_sigs[key])[0]
            )
            high_confidence_count = len(extra_in_impl) - low_confidence_count
            if low_confidence_count > 0:
                summary_parts.append(
                    f"{low_confidence_count} extra in implementation"
                )
            if high_confidence_count > 0:
                summary_parts.append(
                    f"{high_confidence_count} helper(s) suppressed"
                )

        # Convert summary_parts to (label, value) format
        summary_items = []
        for part in summary_parts:
            # Extract count and label from part (e.g., "5 mismatch(es)" -> ("Mismatch(es):", 5))
            import re as re_module
            match = re_module.search(r'(\d+)\s+(.+)', part)
            if match:
                count = int(match.group(1))
                label = match.group(2).strip()
                # Capitalize first letter
                label = label[0].upper() + label[1:] if label else label
                summary_items.append((f"{label}:", count))

        if summary_items:
            output.add_summary_section(summary_items)

        output.add_failure_message("Validation failed. Please fix the errors above.")
        output.print()

        # Exit with error code if there are actual errors (mismatches)
        exit_code = output.get_exit_code(args.no_fail)
        sys.exit(exit_code)


if __name__ == '__main__':
    main()
