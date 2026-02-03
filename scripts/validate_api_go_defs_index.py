#!/usr/bin/env python3
"""
Validate that all type, interface, function, and method definitions in tech specs
are documented in the API definitions index.

This script scans all markdown files in docs/tech_specs (excluding the index file itself)
and extracts Go definitions (types, interfaces, functions, methods)
from Go code blocks (```go), then verifies they are listed in the index file.
Constants and variables are intentionally excluded from the index.

IMPORTANT: This script only processes Go code blocks (```go). Other language code blocks
(```rust, ```zig, ```cpp, etc.) are ignored, as this script is specifically for Go API definitions.
"""

from __future__ import annotations

import argparse
import re
import os
import sys
from pathlib import Path
from typing import Dict, Optional

from lib import _validation_utils
from lib import _index_utils
from lib.go_defs_index import _go_defs_index_discovery as _discovery
from lib.go_defs_index import _go_defs_index_models as _models
from lib.go_defs_index import _go_defs_index_indexfile as _indexfile
from lib.go_defs_index import _go_defs_index_matching as _matching
from lib.go_defs_index import _go_defs_index_comparison as _comparison
from lib.go_defs_index import _go_defs_index_descriptions as _descriptions
from lib.go_defs_index import _go_defs_index_ordering as _ordering
from lib.go_defs_index import _go_defs_index_reporting as _reporting

ParsedIndex = _index_utils.ParsedIndex
discover_all_definitions_phase1 = _discovery.discover_all_definitions
DetectedDefinition = _models.DetectedDefinition
parse_index_indexfile = _indexfile.parse_index
compare_with_index_phase4 = _comparison.compare_with_index
check_entry_descriptions_phase5 = _descriptions.check_entry_descriptions
determine_ordering_phase6 = _ordering.determine_ordering
generate_report_phase7 = _reporting.generate_report

INDEX_FILENAME = "api_go_defs_index.md"

OutputBuilder = _validation_utils.OutputBuilder
parse_no_color_flag = _validation_utils.parse_no_color_flag
get_workspace_root = _validation_utils.get_workspace_root
ValidationIssue = _validation_utils.ValidationIssue
FileContentCache = _validation_utils.FileContentCache
DOCS_DIR = _validation_utils.DOCS_DIR
TECH_SPECS_DIR = _validation_utils.TECH_SPECS_DIR


_WHITESPACE_RE = re.compile(r"\s+")
_SENTENCE_SPLIT_RE = re.compile(r"(?<=[.!?])\s+")
_LABEL_RE = re.compile(
    (
        r"\b("
        r"Parameters|Returns|Validation|Behavior|Errors|Error|Important|Notes|Constraints|"
        r"Security|Usage|Examples|Example"
        r"):"
    ),
    flags=re.IGNORECASE,
)
_VALIDATION_INTRO_RE = re.compile(
    r"\b(performs?|performed)\s+(comprehensive\s+)?validation:\s*",
    flags=re.IGNORECASE,
)
_IMPERATIVE_RE = re.compile(
    r"\b(Ensure|Verify|Check|Validate|Reject|Detect|Require|Confirm)\b"
)


def _normalize_description_text(text: str) -> str:
    cleaned = _WHITESPACE_RE.sub(" ", (text or "").strip())
    # Prefer ASCII docs.
    cleaned = cleaned.replace("→", "=>")
    return cleaned


def _ensure_period(text: str) -> str:
    t = (text or "").strip()
    if not t:
        return ""
    if t[-1] not in ".!?":
        t += "."
    return t


def _split_sentences(text: str) -> list[str]:
    return [s.strip() for s in _SENTENCE_SPLIT_RE.split(text or "") if s.strip()]


def _looks_like_label(text: str) -> bool:
    # Heuristic: "Label: ..." within the first 40 chars.
    t = (text or "").strip()
    if ":" not in t:
        return False
    idx = t.find(":")
    return 2 <= idx <= 40


def _apply_validation_intro_normalization(text: str) -> str:
    # "X performs comprehensive validation: Ensure ..." =>
    # "X performs comprehensive validation.\nValidation: Ensure ..."
    return _VALIDATION_INTRO_RE.sub(r"\1 \2validation.\nValidation: ", text)


def _insert_newlines_before_labels(text: str) -> str:
    def _label_newline(m: re.Match) -> str:
        label = m.group(0)
        return label if not m.start() else "\n" + label

    return _LABEL_RE.sub(_label_newline, text)


def _split_dash_items(text: str) -> list[str]:
    """
    Split embedded dash list text into items.

    Supports:
    - "- a - b - c"
    - " - a - b"
    - "- - a" (strips repeated leading '-' markers)
    """
    t = (text or "").strip().rstrip()
    if t.endswith(" -") or t.endswith("-"):
        t = t.rstrip(" -").rstrip()
    t = re.sub(r"\s+-\s+", "\n", t)
    t = re.sub(r"^(?:-\s+)+", "", t)

    items: list[str] = []
    for line in t.split("\n"):
        s = line.strip()
        if not s:
            continue
        s = re.sub(r"^(?:-\s+)+", "", s).strip()
        if s in ("-", "-.", "."):
            continue
        if s:
            items.append(s)
    return items


def _split_imperative_items(text: str) -> list[str]:
    """
    Split long validation/behavior text into items like:
    "Ensure ... Check ... Verify ..." => ["Ensure ...", "Check ...", "Verify ..."].
    """
    t = (text or "").strip()
    matches = list(_IMPERATIVE_RE.finditer(t))
    if len(matches) < 2:
        return []
    parts: list[str] = []
    for idx, m in enumerate(matches):
        start = m.start()
        end = matches[idx + 1].start() if idx + 1 < len(matches) else len(t)
        seg = t[start:end].strip()
        if seg:
            parts.append(seg)
    return parts


def _compact_context(entry_name: str, prefix: str) -> str:
    """
    Reduce a long prefix down to the nearest useful context label.

    Used when expanding embedded "- item" lists so we don't repeat the entire preceding
    paragraph for every bullet.
    """
    p = _normalize_description_text(prefix)
    if not p:
        return ""

    # If the prefix contains multiple sentences, keep only the last sentence.
    if "." in p:
        p = p.rsplit(".", 1)[-1].strip()

    # If the last sentence still starts with the entry name, drop it.
    if entry_name and p.lower().startswith(entry_name.lower()):
        p = p[len(entry_name):].strip(" :-")

    # If the context has a "Category: details" form, prefer just the category.
    if ":" in p:
        category = p.split(":", 1)[0].strip()
        if 3 <= len(category) <= 50:
            return category

    # If it's still too long, keep the tail end (more specific) rather than the head.
    if len(p) > 80:
        p = p[-80:].lstrip(" ,;-")

    return p.strip(" :-")


def _add_sentence_continuations(bullets: list[str]) -> list[str]:
    """
    Enforce one sentence per line:
    - First sentence stays on the bullet line.
    - Remaining sentences become continuation lines ("CONT: ...").
    """
    out_lines: list[str] = []
    max_bullets = 6
    max_total_lines = 12

    def _count_bullets(lines: list[str]) -> int:
        return len([x for x in lines if not x.startswith("CONT: ")])

    for b in [x for x in bullets if x]:
        if _count_bullets(out_lines) >= max_bullets:
            break
        sentences = _split_sentences(b)
        if not sentences:
            continue
        first_sentence = sentences[0].strip()
        if first_sentence in ("-", "-.", "."):
            continue
        out_lines.append(_ensure_period(first_sentence))
        for s in sentences[1:]:
            if len(out_lines) >= max_total_lines:
                break
            ss = s.strip()
            if ss in ("-", "-.", "."):
                continue
            out_lines.append("CONT: " + _ensure_period(ss))

    return out_lines


def _maybe_prefix_first_bullet(entry_name: str, raw_name: str, bullets: list[str]) -> None:
    if not bullets:
        return
    first = bullets[0].strip()
    first_l = first.lower()
    raw_l = (raw_name or "").strip().lower()
    name_l = (entry_name or "").strip().lower()
    is_structured_label = first_l.startswith(
        ("parameter:", "return:", "validation:", "behavior:", "error:", "errors:")
    )
    has_context_label = ":" in first and first.index(":") <= 40
    if is_structured_label or has_context_label:
        return
    if raw_l and first_l.startswith(raw_l):
        return
    if name_l and first_l.startswith(name_l):
        return
    if entry_name:
        bullets[0] = _ensure_period(f"{entry_name} {first}")


def _ensure_min_description_length(entry_name: str, cleaned: str, bullets: list[str]) -> list[str]:
    combined = " ".join([b for b in bullets if b]).strip()
    if len(combined) >= 20:
        return bullets

    fallback = cleaned
    if entry_name and not fallback.lower().startswith(entry_name.lower()):
        fallback = f"{entry_name} {fallback}"
    fallback = _ensure_period(fallback)
    if len(fallback) >= 20:
        return [fallback]
    return []


def _bullets_from_unlabeled_chunk(entry_name: str, chunk: str) -> list[str]:
    chunk_stripped = (chunk or "").strip()
    if not chunk_stripped:
        return []

    # If this already looks like a dash-list in prose, split it into multiple bullets
    # rather than leaving embedded "- item" text.
    if chunk_stripped.startswith("- "):
        return [_ensure_period(i) for i in _split_dash_items(chunk_stripped)]

    if " - " in chunk:
        dash_pos = chunk.find(" - ")
        if dash_pos != -1 and ":" in chunk[:dash_pos]:
            # Use the last ":" before the first dash item as the split point.
            colon_pos = chunk.rfind(":", 0, dash_pos)
            prefix = chunk[:colon_pos].strip()
            rest = chunk[colon_pos + 1:].strip()

            prefix_bullets = [_ensure_period(p) for p in _split_sentences(prefix)] or [
                _ensure_period(prefix)
            ]
            items = _split_dash_items(rest)
            if not items:
                return prefix_bullets

            ctx = _compact_context(entry_name, prefix)
            item_bullets: list[str] = []
            for it in items:
                if _looks_like_label(it):
                    item_bullets.append(_ensure_period(it))
                elif ctx:
                    item_bullets.append(_ensure_period(f"{ctx}: {it}"))
                else:
                    item_bullets.append(_ensure_period(it))
            return prefix_bullets + item_bullets

        items = _split_dash_items(chunk)
        if len(items) >= 2:
            return [_ensure_period(i) for i in items]

    return [_ensure_period(p) for p in _split_sentences(chunk)] or [_ensure_period(chunk)]


def _bullets_from_labeled_chunk(label: str, rest: str) -> list[str]:
    label_cap = (label or "").capitalize()
    rest = (rest or "").strip()
    if not rest:
        return [_ensure_period(f"{label_cap}.")]

    items = _split_dash_items(rest)
    if items:
        prefix_map = {
            "Parameters": "Parameter",
            "Returns": "Return",
            "Errors": "Error",
            "Error": "Error",
        }
        item_prefix = prefix_map.get(label_cap, label_cap)
        return [_ensure_period(f"{item_prefix}: {it}") for it in items]

    if label_cap in ("Validation", "Behavior") and len(rest) >= 80:
        imperative_items = _split_imperative_items(rest)
        if imperative_items:
            return [_ensure_period(f"{label_cap}: {it}") for it in imperative_items]

    return [_ensure_period(f"{label_cap}: {rest}")]


def _derive_description_lines_from_comments(
    *,
    entry_name: str,
    raw_name: str,
    def_comments: str,
) -> list[str]:
    """
    Convert definition comments into index description bullet lines.

    Rules:
    - Use content derived from def_comments only (no generic placeholder text).
    - Keep one sentence per line when possible.
    - Ensure combined length >= 20 characters so description validation passes.
    """
    cleaned = _normalize_description_text(def_comments)
    if not cleaned:
        return []

    cleaned = _apply_validation_intro_normalization(cleaned)
    labeled = _insert_newlines_before_labels(cleaned)
    chunks = [c.strip() for c in labeled.split("\n") if c.strip()]

    bullets: list[str] = []
    for chunk in chunks if chunks else [cleaned]:
        m = _LABEL_RE.match(chunk)
        if not m:
            bullets.extend(_bullets_from_unlabeled_chunk(entry_name, chunk))
            continue
        label = m.group(1)
        rest = chunk.split(":", 1)[1].strip() if ":" in chunk else ""
        bullets.extend(_bullets_from_labeled_chunk(label, rest))

    _maybe_prefix_first_bullet(entry_name, raw_name, bullets)
    bullets = _ensure_min_description_length(entry_name, cleaned, bullets)
    return _add_sentence_continuations(bullets)


def _populate_missing_expected_descriptions(parsed_index: ParsedIndex) -> int:
    """
    Populate missing expected entry descriptions from def_comments for apply.

    Applies when:
    - Entry has no description_lines AND
      - it has def_comments (missing description), OR
      - it is newly added to the index AND has def_comments.
    """
    populated = 0
    for section_path in parsed_index.section_order:
        section = parsed_index.sections.get(section_path)
        if not section:
            continue
        for entry in section.expected_entries.values():
            if entry.description_lines:
                continue
            if not entry.def_comments:
                continue
            desc_lines = _derive_description_lines_from_comments(
                entry_name=entry.name,
                raw_name=entry.raw_name,
                def_comments=entry.def_comments,
            )
            if not desc_lines:
                continue
            entry.description_lines = desc_lines
            populated += 1
    return populated


def _apply_index_updates(
    index_file: Path,
    parsed_index: ParsedIndex,
) -> bool:
    # --apply is intentionally interactive-only. In some environments (e.g. when stdin is
    # redirected by a runner), input() can block without showing the prompt. To avoid a
    # "silent hang", read confirmation from /dev/tty when possible.
    has_tty_stdin = sys.stdin.isatty()
    tty_path = "/dev/tty"
    has_dev_tty = os.path.exists(tty_path)

    # Description fixes are intentionally applied by --apply when they can be derived
    # from def_comments, even if there are no structural/link changes.
    has_description_fix_candidates = False
    for section_path in parsed_index.section_order:
        section = parsed_index.sections.get(section_path)
        if not section:
            continue
        for name, current_entry in section.current_entries.items():
            if current_entry.has_description:
                continue
            expected_entry = section.expected_entries.get(name)
            if not expected_entry or not expected_entry.def_comments:
                continue
            has_description_fix_candidates = True
            break
        if has_description_fix_candidates:
            break

    pending_changes = any(
        [
            parsed_index.get_added_entries(),
            parsed_index.get_moved_entries(),
            parsed_index.get_removed_entries(),
            parsed_index.get_orphans(),
            parsed_index.get_link_update_entries(),
            parsed_index.get_reordered_entries(),
            has_description_fix_candidates,
        ]
    )

    if not pending_changes:
        print("No high-confidence updates to apply.")
        return False

    print(
        "Apply will overwrite the entire index file, including overview and table of contents."
    )
    # NOTE: This is intentionally set up to prevent automation from using it; DO NOT CHANGE THIS.
    print("Type 'yes' to confirm.")
    sys.stdout.flush()
    try:
        if has_tty_stdin:
            confirm = input().lower().strip()
        else:
            # Use the controlling terminal explicitly when stdin isn't interactive.
            if not has_dev_tty:
                print("ERROR: --apply requires an interactive terminal (stdin must be a TTY)")
                sys.exit(1)
            try:
                with open(tty_path, "r", encoding="utf-8") as tty_in:
                    confirm = tty_in.readline().lower().strip()
            except OSError as e:
                print(
                    "ERROR: --apply requires an interactive terminal "
                    f"(failed to read {tty_path}: {e})"
                )
                sys.exit(1)
    except (EOFError, KeyboardInterrupt):
        print("Apply canceled.")
        return False
    if confirm != "yes":
        print("Apply canceled.")
        return False

    parsed_index.sync_expected_descriptions()
    _populate_missing_expected_descriptions(parsed_index)
    updated_content = parsed_index.to_markdown()
    index_file.write_text(updated_content, encoding="utf-8")
    print("Index file updated.")
    return True


def _fatal_validation_issue(
    *,
    output: OutputBuilder,
    no_fail: bool,
    issue: ValidationIssue,
) -> None:
    output.add_error_line(issue.format_message(no_color=output.no_color))
    output.add_failure_message("Validation failed. Please fix the errors above.")
    output.print()
    sys.exit(output.get_exit_code(no_fail))


def _read_index_file_or_exit(
    *,
    file_cache: FileContentCache,
    index_file: Path,
    output: OutputBuilder,
    no_fail: bool,
) -> str:
    try:
        return file_cache.get_content(index_file)
    except (IOError, OSError) as e:
        _fatal_validation_issue(
            output=output,
            no_fail=no_fail,
            issue=ValidationIssue.create(
                "index_file_read_error",
                index_file,
                0,
                0,
                message=f"Could not read index file: {e}",
                severity="error",
            ),
        )
    except UnicodeDecodeError as e:
        _fatal_validation_issue(
            output=output,
            no_fail=no_fail,
            issue=ValidationIssue.create(
                "index_file_decode_error",
                index_file,
                0,
                0,
                message=f"Could not decode index file (encoding issue): {e}",
                severity="error",
            ),
        )
    except (ValueError, KeyError, TypeError, RuntimeError, MemoryError, BufferError) as e:
        _fatal_validation_issue(
            output=output,
            no_fail=no_fail,
            issue=ValidationIssue.create(
                "index_file_unexpected_error",
                index_file,
                0,
                0,
                message=f"Unexpected error reading index file: {e}",
                severity="error",
            ),
        )
    return ""


def _parse_index_or_exit(
    *,
    index_content: str,
    index_file: Path,
    output: OutputBuilder,
    no_fail: bool,
) -> ParsedIndex:
    try:
        return parse_index_indexfile(index_content)
    except ValueError as e:
        _fatal_validation_issue(
            output=output,
            no_fail=no_fail,
            issue=ValidationIssue.create(
                "duplicate_headings",
                index_file,
                0,
                0,
                message=f"{e}",
                severity="error",
            ),
        )
    return parse_index_indexfile(index_content)


def _add_index_summary(
    *,
    output: OutputBuilder,
    definitions_count: int,
    parsed_index: ParsedIndex,
    description_errors: int,
    zero_confidence_counts: Optional[Dict[str, int]] = None,
) -> int:
    added_entries = len(parsed_index.get_added_entries())
    moved_entries = len(parsed_index.get_moved_entries())
    orphaned_entries = len(parsed_index.get_orphans())
    removed_entries = len(parsed_index.get_removed_entries())
    link_updates = len(parsed_index.get_link_update_entries())
    unresolved_entries = len(parsed_index.get_unresolved_entries())

    total_issues = (
        added_entries
        + orphaned_entries
        + moved_entries
        + removed_entries
        + link_updates
        + unresolved_entries
        + description_errors
    )

    if total_issues > 0:
        summary_items = []
        if added_entries:
            summary_items.append(("  missing definition(s):", added_entries))
        if orphaned_entries:
            summary_items.append(("  orphaned entry/entries:", orphaned_entries))
        if moved_entries:
            summary_items.append(("  section mismatch(es):", moved_entries))
        if removed_entries:
            summary_items.append(("  removed entry/entries:", removed_entries))
        if link_updates:
            summary_items.append(("  link update(s):", link_updates))
        if unresolved_entries:
            summary_items.append(("  unresolved definition(s):", unresolved_entries))
        if description_errors > 0:
            summary_items.append(("  description error(s):", description_errors))

        if summary_items:
            output.add_summary_header()
            summary_items.insert(0, ("Total issues:", total_issues))
            if zero_confidence_counts:
                summary_items.extend(_format_zero_confidence_summary(zero_confidence_counts))
            output.add_summary_section(summary_items)
            output.add_failure_message("Validation failed. Please fix the errors above.")
        return total_issues

    output.add_summary_header()
    summary_items = [
        ("Definitions checked:", definitions_count),
        ("All definitions indexed:", definitions_count),
    ]
    if zero_confidence_counts:
        summary_items.extend(_format_zero_confidence_summary(zero_confidence_counts))
    output.add_summary_section(summary_items)
    return 0


def _format_zero_confidence_summary(
    zero_confidence_counts: Dict[str, int],
) -> list[tuple[str, int]]:
    summary_items: list[tuple[str, int]] = []
    for label, key in (
        ("Zero-confidence types:", "type"),
        ("Zero-confidence functions:", "func"),
        ("Zero-confidence methods:", "method"),
        ("Zero-confidence total:", "total"),
    ):
        count = zero_confidence_counts.get(key, 0)
        if count:
            summary_items.append((label, count))
    return summary_items


def _count_zero_confidence(definitions: list[DetectedDefinition]) -> Dict[str, int]:
    counts = {"type": 0, "func": 0, "method": 0, "total": 0}
    for definition in definitions:
        score = definition.confidence_score
        if score is None or score > 0.0:
            continue
        if definition.kind in counts:
            counts[definition.kind] += 1
        counts["total"] += 1
    return counts


def _build_arg_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        description=(
            "Validate that all Go definitions in tech specs are in the "
            "Go API definitions index"
        )
    )
    parser.add_argument(
        "--verbose",
        "-v",
        action="store_true",
        help="Show verbose output",
    )
    parser.add_argument(
        "--index-file",
        type=str,
        default=f"{DOCS_DIR}/{TECH_SPECS_DIR}/{INDEX_FILENAME}",
        help="Path to the API definitions index file",
    )
    parser.add_argument(
        "--output",
        "-o",
        type=str,
        metavar="FILE",
        help="Write detailed output to FILE",
    )
    parser.add_argument(
        "--nocolor",
        "--no-color",
        action="store_true",
        help="Disable colored output",
    )
    parser.add_argument(
        "--no-fail",
        action="store_true",
        help="Exit with code 0 even if errors are found",
    )
    parser.add_argument(
        "--apply",
        action="store_true",
        help="Apply high-confidence index updates (interactive confirmation required)",
    )
    return parser


def main() -> None:
    args = _build_arg_parser().parse_args()

    # Create output builder early so fatal argument/path errors are consistent.
    no_color = args.nocolor or parse_no_color_flag(sys.argv)
    output = OutputBuilder(
        "Go API Definitions Index",
        "Validates all definitions are in the index",
        no_color=no_color,
        verbose=args.verbose,
        output_file=args.output,
    )

    repo_root = get_workspace_root()
    tech_specs_dir = repo_root / DOCS_DIR / TECH_SPECS_DIR
    index_file = repo_root / args.index_file

    if not tech_specs_dir.exists():
        _fatal_validation_issue(
            output=output,
            no_fail=args.no_fail,
            issue=ValidationIssue.create(
                "tech_specs_dir_not_found",
                tech_specs_dir,
                0,
                0,
                message=f"Tech specs directory not found: {tech_specs_dir}",
                severity="error",
            ),
        )

    if not index_file.exists():
        _fatal_validation_issue(
            output=output,
            no_fail=args.no_fail,
            issue=ValidationIssue.create(
                "index_file_not_found",
                index_file,
                0,
                0,
                message=f"Index file not found: {index_file}",
                severity="error",
            ),
        )

    file_cache = FileContentCache()
    index_content = _read_index_file_or_exit(
        file_cache=file_cache,
        index_file=index_file,
        output=output,
        no_fail=args.no_fail,
    )

    # Phase 1: Discovery - Find all definitions
    definitions = discover_all_definitions_phase1(
        tech_specs_dir,
        output,
        file_cache=file_cache,
    )

    # Phase 2: Heading Resolution is handled during discovery.

    # Parse index once and share across phases.
    parsed_index = _parse_index_or_exit(
        index_content=index_content,
        index_file=index_file,
        output=output,
        no_fail=args.no_fail,
    )

    if output.verbose:
        total_sections = len(parsed_index.section_order)
        total_current_entries = 0
        for section_path in parsed_index.section_order:
            section = parsed_index.sections.get(section_path)
            if not section:
                continue
            total_current_entries += len(section.current_entries)
        output.add_verbose_line(
            f"Parsed index: {total_sections} sections, {total_current_entries} current entries"
        )

    # Phase 3: Section Placement - Determine index section
    _matching.determine_section_placement(
        definitions,
        parsed_index,
        output=output,
    )

    # Phase 4: Index Comparison - Compare with current index
    compare_with_index_phase4(
        parsed_index,
        output,
    )

    # Phase 5: Entry Description Validation
    try:
        description_errors = check_entry_descriptions_phase5(
            parsed_index,
            index_file,
            output,
        )
    except (ValueError, KeyError, TypeError, AttributeError, RuntimeError) as e:
        if output:
            output.add_error_line(
                ValidationIssue.create(
                    "error_parsing_descriptions",
                    index_file,
                    0,
                    0,
                    message=f"Could not parse index file for description checking: {e}",
                    severity="error",
                ).format_message(no_color=output.no_color)
            )
        description_errors = 0

    # Phase 6: Sorting - Determine correct ordering
    determine_ordering_phase6(
        parsed_index,
        output,
    )

    # Phase 7: Reporting - Generate output
    index_file_name = Path(args.index_file).name
    generate_report_phase7(
        parsed_index,
        output,
        index_file_name,
    )

    zero_confidence_counts = _count_zero_confidence(definitions)
    total_issues = _add_index_summary(
        output=output,
        definitions_count=len(definitions),
        parsed_index=parsed_index,
        description_errors=description_errors,
        zero_confidence_counts=zero_confidence_counts,
    )

    if output.verbose:
        output.add_blank_line("final_message")
        output.add_line("Expected index (full tree):", section="final_message")
        output.add_blank_line("final_message")
        for line in parsed_index.render_full_tree():
            output.add_line(line, section="final_message")
        output.add_blank_line("final_message")

    if not total_issues:
        if output.has_warnings():
            output.add_warnings_only_message(
                verbose_hint="Run with --verbose to see the full warning details.",
            )
        else:
            output.add_success_message(
                "No errors or suggestions found. All definitions are correctly indexed."
            )

    final_exit_code = output.get_exit_code(args.no_fail)
    output.print()
    if args.apply:
        _apply_index_updates(index_file, parsed_index)
    sys.exit(final_exit_code)


if __name__ == "__main__":
    main()
