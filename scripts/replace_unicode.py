#!/usr/bin/env python3
"""
Replace standard Unicode characters with ASCII in a Markdown file.

Replaces only in prose (not inside fenced code blocks or inline backticks).
Standard replacements:
  U+2192 (→)  ->  ->
  U+201C (")  ->  "
  U+201D (")  ->  "
  U+2019 (')  ->  '

Usage:
  make replace-unicode FILE="path/to/file.md" [DRY_RUN=1]
  python3 replace_unicode.py FILE [--dry-run]
  Default: overwrite FILE.  --dry-run / DRY_RUN=1: print to stdout, do not modify FILE.
"""

from __future__ import annotations

import argparse
import re
import sys

# Standard Unicode -> ASCII replacements (prose only)
REPLACEMENTS = (
    ("\u2192", "->"),   # RIGHTWARDS ARROW
    ("\u201c", "\""),   # LEFT DOUBLE QUOTATION MARK
    ("\u201d", "\""),   # RIGHT DOUBLE QUOTATION MARK
    ("\u2019", "'"),    # RIGHT SINGLE QUOTATION MARK (apostrophe)
)


def _replace_prose_line(line: str) -> str:
    """Apply replacements only outside inline code (backtick spans)."""
    result = []
    i = 0
    in_backticks = False
    while i < len(line):
        if line[i] == "`":
            result.append("`")
            in_backticks = not in_backticks
            i += 1
            continue
        if in_backticks:
            result.append(line[i])
            i += 1
            continue
        # Prose: apply first matching replacement
        replaced = False
        for old, new in REPLACEMENTS:
            if line[i:i + len(old)] == old:
                result.append(new)
                i += len(old)
                replaced = True
                break
        if not replaced:
            result.append(line[i])
            i += 1
    return "".join(result)


def process(lines: list[str]) -> list[str]:
    """Process lines: replace Unicode only in prose (skip fenced blocks and inline code)."""
    out = []
    fence = None  # current fence char ("```" or "~~~") or None
    fence_pattern = re.compile(r"^(\s*)(```|~~~)")

    for line in lines:
        # Track fenced code blocks
        match = fence_pattern.match(line)
        if match:
            _, delim = match.groups()
            if fence is None:
                fence = delim
                out.append(line)
                continue
            if delim == fence:
                fence = None
            out.append(line)
            continue

        if fence is not None:
            out.append(line)
            continue

        out.append(_replace_prose_line(line))

    return out


def main() -> int:
    parser = argparse.ArgumentParser(
        description="Replace standard Unicode with ASCII in a Markdown file (prose only)."
    )
    parser.add_argument("file", help="Path to the Markdown file")
    parser.add_argument(
        "--dry-run",
        action="store_true",
        help="Print result to stdout, do not modify the file",
    )
    args = parser.parse_args()

    try:
        with open(args.file, encoding="utf-8") as f:
            original = f.read()
    except OSError as e:
        print(f"replace_unicode: {e}", file=sys.stderr)
        return 1

    lines = original.splitlines(keepends=True)
    processed = process(lines)
    result = "".join(processed)

    if args.dry_run:
        print(result, end="")
        return 0

    try:
        with open(args.file, "w", encoding="utf-8") as f:
            f.write(result)
    except OSError as e:
        print(f"replace_unicode: {e}", file=sys.stderr)
        return 1

    return 0


if __name__ == "__main__":
    sys.exit(main())
