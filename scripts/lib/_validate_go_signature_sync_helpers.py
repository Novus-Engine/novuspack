"""Helpers for Go signature sync validation: report emission."""

import re
import sys
from typing import Dict, List, Tuple

from lib._go_code_utils import Signature
from lib._validation_utils import OutputBuilder, format_issue_message


def _parse_location(location_str: str) -> Tuple[str, int]:
    """Parse 'path:line' into (path, line_num)."""
    if ':' in location_str:
        parts = location_str.rsplit(':', 1)
        try:
            return parts[0], int(parts[1])
        except (ValueError, IndexError):
            pass
    return location_str, 0


def emit_mismatches(
    output: OutputBuilder, mismatches: List, *, no_color: bool = False
) -> bool:
    """Emit mismatch errors; return True if any."""
    if not mismatches:
        return False
    output.add_errors_header()
    output.add_line(f"Found {len(mismatches)} signature mismatch(es):", section="error")
    output.add_blank_line("error")
    for key, impl_sig, spec_sig in sorted(mismatches):
        file_path, line_num = _parse_location(impl_sig.location)
        msg = format_issue_message(
            "error",
            "signature_mismatch",
            file_path,
            line_num=line_num if line_num else None,
            message=key,
            no_color=no_color,
        )
        output.add_error_line(msg)
        output.add_error_line(f"  Implementation: {impl_sig.normalized_signature()}")
        output.add_error_line(f"    Location: {impl_sig.location}")
        output.add_error_line(f"  Specification:  {spec_sig.normalized_signature()}")
        output.add_error_line(f"    Location: {spec_sig.location}")
    return True


def emit_missing_in_impl(
    output: OutputBuilder,
    missing_in_impl: List[str],
    spec_sigs: Dict[str, Signature],
    *,
    no_color: bool = False,
) -> bool:
    """Emit missing-in-impl warnings; return True if any."""
    if not missing_in_impl:
        return False
    output.add_warnings_header()
    output.add_line(
        f"Found {len(missing_in_impl)} signature(s) in specs but not in implementation:",
        section="warning",
    )
    output.add_blank_line("warning")
    for key in sorted(missing_in_impl):
        spec_sig = spec_sigs[key]
        file_path, line_num = _parse_location(spec_sig.location)
        msg = format_issue_message(
            "warning",
            "missing_in_implementation",
            file_path,
            line_num=line_num if line_num else None,
            message=key,
            no_color=no_color,
        )
        output.add_warning_line(msg)
        output.add_warning_line(f"    Signature: {spec_sig.normalized_signature()}")
        output.add_warning_line(f"    Location: {spec_sig.location}")
    return True


def emit_extra_in_impl_section(
    output: OutputBuilder,
    args,
    *,
    extra_in_impl: List[str],
    impl_sigs: Dict[str, Signature],
    high_confidence_helpers: List[Tuple[str, Signature, List[str]]],
    low_confidence_extra: List[str],
    errors_public_api_missing: List[Tuple[str, Signature]],
    no_color: bool = False,
) -> Tuple[bool, bool]:
    """Emit extra-in-impl section. Returns (has_errors, has_warnings)."""
    if not extra_in_impl:
        return (False, False)
    has_errors = False
    has_warnings = False
    if high_confidence_helpers:
        if args.verbose:
            output.add_verbose_line(
                f"Suppressed {len(high_confidence_helpers)} high-confidence helper function(s)"
            )
            output.add_blank_line("working_verbose")
            for key, impl_sig, reasons in sorted(high_confidence_helpers):
                output.add_verbose_line(f"  {key}")
                output.add_verbose_line(f"    Signature: {impl_sig.normalized_signature()}")
                output.add_verbose_line(f"    Location: {impl_sig.location}")
                output.add_verbose_line(f"    Reasons: {', '.join(reasons)}")

    if errors_public_api_missing:
        has_errors = True
        output.add_errors_header()
        output.add_line(
            f"Found {len(errors_public_api_missing)} public API method(s) "
            "in implementation but not in specs:",
            section="error",
        )
        output.add_blank_line("error")
        output.add_line(
            "These are public methods on public API types and MUST be documented in tech specs.",
            section="error",
        )
        output.add_blank_line("error")
        for key, impl_sig in sorted(errors_public_api_missing):
            file_path, line_num = _parse_location(impl_sig.location)
            msg = format_issue_message(
                "error",
                "public_api_not_in_spec",
                file_path,
                line_num=line_num if line_num else None,
                message=key,
                no_color=no_color,
            )
            output.add_error_line(msg)
            output.add_error_line(f"  Implementation: {impl_sig.normalized_signature()}")
            output.add_error_line(f"    Location: {impl_sig.location}")

    if low_confidence_extra:
        has_warnings = True
        output.add_warnings_header()
        output.add_line(
            f"Found {len(low_confidence_extra)} signature(s) in implementation but not in specs:",
            section="warning",
        )
        if high_confidence_helpers:
            if args.verbose:
                msg = (
                    f"(Suppressed {len(high_confidence_helpers)} high-confidence helper "
                    "function(s) - see above)"
                )
            else:
                msg = (
                    f"(Suppressed {len(high_confidence_helpers)} high-confidence helper "
                    "function(s) - use --verbose to see them)"
                )
            output.add_line(msg, section="warning")
        if errors_public_api_missing:
            output.add_line(
                f"(Also found {len(errors_public_api_missing)} public API method(s) "
                "missing from specs - see errors above)",
                section="warning",
            )
        output.add_blank_line("warning")
        output.add_line("(These may be helper functions, but should be checked)", section="warning")
        output.add_blank_line("warning")
        for key in sorted(low_confidence_extra):
            impl_sig = impl_sigs[key]
            file_path, line_num = _parse_location(impl_sig.location)
            msg = format_issue_message(
                "warning",
                "extra_in_implementation",
                file_path,
                line_num=line_num if line_num else None,
                message=key,
                no_color=no_color,
            )
            output.add_warning_line(msg)
            output.add_warning_line(f"    Signature: {impl_sig.normalized_signature()}")
            output.add_warning_line(f"    Location: {impl_sig.location}")
    elif high_confidence_helpers and not args.verbose:
        output.add_verbose_line(
            f"Suppressed {len(high_confidence_helpers)} high-confidence helper function(s)"
        )
        output.add_verbose_line("(Use --verbose to see the list of suppressed helpers)")
    return (has_errors, has_warnings)


def emit_sync_final(
    output: OutputBuilder,
    args,
    *,
    has_errors: bool,
    has_warnings: bool,
    impl_sigs: Dict[str, Signature],
    spec_sigs: Dict[str, Signature],
    mismatches: List,
    missing_in_impl: List[str],
    extra_in_impl: List[str],
    low_confidence_count: int,
    high_confidence_count: int,
) -> None:
    """Emit success or failure summary and exit."""
    if not has_errors and not has_warnings:
        output.add_success_message("All signatures are in sync!")
        if args.verbose:
            output.add_verbose_line(f"  - {len(impl_sigs)} signatures in implementation")
            output.add_verbose_line(f"  - {len(spec_sigs)} signatures in specs")
        output.print()
        sys.exit(0)
    summary_parts = []
    if mismatches:
        summary_parts.append(f"{len(mismatches)} mismatch(es)")
    if missing_in_impl:
        summary_parts.append(f"{len(missing_in_impl)} missing in implementation")
    if extra_in_impl:
        if low_confidence_count > 0:
            summary_parts.append(f"{low_confidence_count} extra in implementation")
        if high_confidence_count > 0:
            summary_parts.append(f"{high_confidence_count} helper(s) suppressed")
    summary_items = []
    for part in summary_parts:
        match = re.search(r'(\d+)\s+(.+)', part)
        if match:
            count = int(match.group(1))
            label = match.group(2).strip()
            label = label[0].upper() + label[1:] if label else label
            summary_items.append((f"{label}:", count))
    if summary_items:
        output.add_summary_section(summary_items)
    if has_errors:
        output.add_failure_message("Validation failed. Please fix the errors above.")
    else:
        output.add_warnings_only_message(
            verbose_hint="Run with --verbose to see the full warning details.",
        )
    output.print()
    sys.exit(output.get_exit_code(args.no_fail))
