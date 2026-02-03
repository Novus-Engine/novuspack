#!/usr/bin/env python3
"""
Validate Go file specification references against api_go_defs_index.md and spec files.

This script:
1. Parses all Go files for "Specification:" comments
2. Extracts file and section/anchor references
3. Validates that referenced files exist
4. Validates that section anchors exist in those files
5. Optionally checks against api_go_defs_index.md for method/type references

Usage:
    python3 scripts/validate_go_spec_references.py [options]

Options:
    --verbose, -v           Show detailed progress information
    --repo-root DIR         Repository root directory (default: parent of script directory)
    --check-index           Also validate against api_go_defs_index.md (not yet implemented)
    --output, -o FILE        Output file path for validation report (default: stdout)
    --help, -h               Show this help message
"""

import argparse
import sys
from pathlib import Path

from lib._validation_utils import (
    OutputBuilder, get_workspace_root, parse_no_color_flag,
)
from lib._validate_go_spec_references_validator import SpecValidator


def main():
    parser = argparse.ArgumentParser(
        description="Validate Go file specification references against spec files",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  python3 scripts/validate_go_spec_references.py
  python3 scripts/validate_go_spec_references.py --verbose
  python3 scripts/validate_go_spec_references.py --output report.txt
        """
    )
    parser.add_argument(
        "--repo-root",
        type=Path,
        default=get_workspace_root(),
        help="Repository root directory (default: parent of script directory)",
    )
    parser.add_argument(
        "--check-index",
        action="store_true",
        help="Also validate against api_go_defs_index.md (not yet implemented)",
    )
    parser.add_argument(
        "--verbose", "-v",
        action="store_true",
        help="Show all references being checked",
    )
    parser.add_argument(
        "--output", "-o",
        type=str,
        help="Output file path for validation report (default: stdout)",
    )
    parser.add_argument(
        "--no-fail",
        action="store_true",
        help="Exit with code 0 even if errors are found",
    )
    parser.add_argument(
        "--nocolor", "--no-color",
        action="store_true",
        help="Disable colored output",
    )

    args = parser.parse_args()

    no_color = args.nocolor or parse_no_color_flag(sys.argv)
    output = OutputBuilder(
        "Go Specification References Validation",
        "Validates Go file specification references against spec files",
        no_color=no_color,
        verbose=args.verbose,
        output_file=args.output
    )

    validator = SpecValidator(args.repo_root)

    if args.verbose:
        output.add_verbose_line(f"Repository root: {args.repo_root}")
        output.add_verbose_line(f"Docs directory: {validator.docs_dir}")
        output.add_verbose_line(f"API Go directory: {validator.api_go_dir}")
        output.add_verbose_line(f"Index file: {validator.index_file}")
        output.add_blank_line("working_verbose")

    error_count, error_messages = validator.validate_all(
        _check_index=args.check_index, verbose=args.verbose, output=output
    )

    if error_messages:
        output.add_errors_header()
        for error_msg in error_messages:
            output.add_error_line(error_msg)
        output.add_blank_line("error")
        output.add_line(
            f"Found {error_count} error(s)",
            section="error"
        )
        has_suggestions = any(" -> " in msg for msg in error_messages)
        if has_suggestions:
            output.add_blank_line("error")
            output.add_line(
                "Note: After applying these reference updates, "
                "verify that each updated reference points to the correct content.",
                section="error"
            )
        output.add_failure_message("Validation failed. Please fix the errors above.")
    elif output.has_warnings():
        output.add_warnings_only_message(
            verbose_hint="Run with --verbose to see the full warning details.",
        )
    else:
        output.add_success_message("All specification references are valid!")

    output.print()
    return output.get_exit_code(args.no_fail)


if __name__ == "__main__":
    sys.exit(main())
