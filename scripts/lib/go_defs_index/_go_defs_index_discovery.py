from __future__ import annotations

from pathlib import Path
import re
from typing import Dict, List, Optional

from lib._go_code_utils import (
    extract_go_doc_comment_above,
    find_go_code_blocks,
    is_example_code,
    normalize_generic_name,
    parse_go_def_signature,
)
from lib.go_defs_index._go_defs_index_models import DetectedDefinition
from lib._validation_utils import (
    FileContentCache,
    OutputBuilder,
    ValidationIssue,
    extract_headings,
    generate_anchor_from_heading,
    is_in_dot_directory,
    find_heading_before_line,
    find_heading_for_code_block,
)
from lib.go_defs_index._go_defs_index_headings import resolve_canonical_reference

INDEX_FILENAME = "api_go_defs_index.md"
_QUALIFIED_TYPE_RE = re.compile(
    r"\b([A-Za-z_][A-Za-z0-9_]*)\.([A-Za-z_][A-Za-z0-9_]*)\b"
)
_IDENT_RE = re.compile(r"[A-Za-z_][A-Za-z0-9_]*")


def _collect_type_names(type_str: str) -> List[str]:
    if not type_str:
        return []
    names: List[str] = []
    for match in _QUALIFIED_TYPE_RE.finditer(type_str):
        candidate = match.group(2)
        if candidate and candidate[0].isupper():
            names.append(candidate)
    for token in _IDENT_RE.findall(type_str):
        if token and token[0].isupper():
            names.append(token)
    seen = set()
    ordered = []
    for name in names:
        if name in seen:
            continue
        seen.add(name)
        ordered.append(name)
    return ordered


def _extract_param_type(param: str) -> str:
    text = param.strip()
    if not text:
        return ""
    if text.startswith("..."):
        text = text[3:].strip()
    parts = text.split()
    if len(parts) >= 2:
        return " ".join(parts[1:])
    return parts[0]


def _extract_signature_types(sig) -> tuple[List[str], List[str], List[str], List[str]]:
    input_types: List[str] = []
    output_types: List[str] = []
    referenced_types: List[str] = []
    referenced_methods: List[str] = []

    if sig.kind != "func":
        return input_types, output_types, referenced_types, referenced_methods

    params = sig.params or ""
    returns = sig.returns or ""

    for param in params.split(","):
        type_str = _extract_param_type(param)
        if not type_str:
            continue
        input_types.extend(_collect_type_names(type_str))

    returns_str = returns.strip()
    if returns_str.startswith("(") and returns_str.endswith(")"):
        returns_str = returns_str[1:-1]
    for item in returns_str.split(","):
        type_str = _extract_param_type(item)
        if not type_str:
            continue
        output_types.extend(_collect_type_names(type_str))

    signature_text = f"{sig.name}({params}) {returns}"
    method_pattern = r"\b([A-Za-z_][A-Za-z0-9_]*)\.([A-Za-z_][A-Za-z0-9_]*)\b"
    for match in re.finditer(method_pattern, signature_text):
        type_name = normalize_generic_name(match.group(1))
        method_name = match.group(2)
        if type_name and type_name[0].isupper():
            referenced_methods.append(f"{type_name}.{method_name}")
            referenced_types.append(type_name)

    referenced_types = list(dict.fromkeys(input_types + output_types + referenced_types))
    referenced_methods = list(dict.fromkeys(referenced_methods))
    return input_types, output_types, referenced_types, referenced_methods


def _build_duplicate_groups(
    definitions_by_name: Dict[str, List[DetectedDefinition]],
) -> Dict[str, List[DetectedDefinition]]:
    """Return groups of definitions with same name in multiple files."""
    duplicate_groups: Dict[str, List[DetectedDefinition]] = {}
    for name, defns in definitions_by_name.items():
        if len(defns) > 1:
            files = {d.file for d in defns}
            if len(files) > 1:
                duplicate_groups[name] = defns
    return duplicate_groups


def _emit_duplicate_errors(
    duplicate_groups: Dict[str, List[DetectedDefinition]],
    output: OutputBuilder,
    tech_specs_dir: Path,
    content_cache: FileContentCache,
) -> None:
    """Emit error lines for duplicate definition groups."""
    output.add_errors_header()
    output.add_line(
        f"Found {len(duplicate_groups)} duplicate definition(s) across multiple files:",
        section="error",
    )
    output.add_blank_line("error")
    for name, defns in sorted(duplicate_groups.items()):
        locations = []
        for defn in sorted(defns, key=lambda d: (d.file, d.code_block_start_line)):
            anchor = _anchor_for_definition(defn, tech_specs_dir, content_cache)
            locations.append(f"{defn.file}{anchor}:{defn.code_block_start_line}")
        locations_str = f"({','.join(locations)})"
        first_defn = defns[0]
        output.add_error_line(
            ValidationIssue.create(
                "Duplicate definition",
                Path(first_defn.file),
                first_defn.code_block_start_line,
                first_defn.code_block_start_line,
                message=f"Definition '{name}' found in multiple files {locations_str}",
                severity="error",
            ).format_message(no_color=output.no_color)
        )


def _anchor_for_definition(defn, tech_specs_dir: Path, content_cache) -> str:
    """Return anchor string for a definition (heading-based or #line-N)."""
    file_path = tech_specs_dir / defn.file
    if not file_path.exists():
        return f"#line-{defn.code_block_start_line}"
    try:
        file_content = content_cache.get_content(file_path)
        heading = find_heading_for_code_block(
            file_content, defn.code_block_start_line
        )
        if heading:
            return generate_anchor_from_heading(heading, include_hash=True)
    except (ValueError, AttributeError, KeyError, OSError, RuntimeError, TypeError):
        pass
    return f"#line-{defn.code_block_start_line}"


def _definition_from_signature(
    sig,
    block_lines: List[str],
    line_index: int,
    *,
    md_file: Path,
    start_line: int,
    code_content: str,
) -> DetectedDefinition:
    """Build one DetectedDefinition from a parsed signature (kind dispatch)."""
    def_comments = extract_go_doc_comment_above(block_lines, line_index)
    if sig.kind == "method" and sig.receiver:
        receiver_type = normalize_generic_name(sig.receiver)
        normalized_name = f"{receiver_type}.{sig.name}"
        return DetectedDefinition(
            name=normalized_name,
            kind="method",
            file=md_file.name,
            code_block_start_line=start_line,
            code_block_content=code_content,
            raw_name=sig.name,
            receiver_type=receiver_type,
            def_comments=def_comments or None,
        )
    if sig.kind == "func":
        normalized_name = normalize_generic_name(sig.name)
        (
            input_types,
            output_types,
            referenced_types,
            referenced_methods,
        ) = _extract_signature_types(sig)
        return DetectedDefinition(
            name=normalized_name,
            kind="func",
            file=md_file.name,
            code_block_start_line=start_line,
            code_block_content=code_content,
            raw_name=sig.name,
            input_types=input_types,
            output_types=output_types,
            referenced_types=referenced_types,
            referenced_methods=referenced_methods,
            def_comments=def_comments or None,
        )
    if sig.kind in ("interface", "struct"):
        normalized_name = normalize_generic_name(sig.name)
        return DetectedDefinition(
            name=normalized_name,
            kind=sig.kind,
            file=md_file.name,
            code_block_start_line=start_line,
            code_block_content=code_content,
            raw_name=sig.name,
            def_comments=def_comments or None,
        )
    normalized_name = normalize_generic_name(sig.name)
    return DetectedDefinition(
        name=normalized_name,
        kind="type",
        file=md_file.name,
        code_block_start_line=start_line,
        code_block_content=code_content,
        raw_name=sig.name,
        def_comments=def_comments or None,
    )


def _definitions_from_block(
    block_lines: List[str],
    code_content: str,
    start_line: int,
    *,
    md_file: Path,
    content: str,
    lines: List[str],
    headings: List[tuple],
    tech_specs_dir: Path,
    output: Optional[OutputBuilder],
    content_cache: FileContentCache,
) -> List[DetectedDefinition]:
    """Return all DetectedDefinitions from one code block (line defs only)."""
    heading_above = find_heading_before_line(content, start_line)
    heading_text = (
        heading_above.heading_text if heading_above and heading_above.heading_text else None
    )
    result: List[DetectedDefinition] = []
    for i, line in enumerate(block_lines):
        if is_example_code(
            code_content, start_line,
            lines=lines, heading_text=heading_text, check_single_line=i,
        ):
            continue
        sig = parse_go_def_signature(
            line, location=f"{md_file.name}:{start_line + i}",
        )
        if sig is None:
            continue
        defn = _definition_from_signature(
            sig, block_lines, i,
            md_file=md_file, start_line=start_line, code_content=code_content,
        )
        _populate_heading_context(
            defn, content, lines, headings,
            tech_specs_dir=tech_specs_dir,
            output=output,
            content_cache=content_cache,
        )
        result.append(defn)
    return result


def discover_all_definitions(
    tech_specs_dir: Path,
    output: Optional[OutputBuilder] = None,
    index_filename: str = INDEX_FILENAME,
    file_cache: Optional[FileContentCache] = None,
) -> List[DetectedDefinition]:
    """
    Phase 1: Find and normalize all definitions from tech spec files.

    Scans all tech spec files, extracts Go definitions from code blocks,
    normalizes method names, filters example code, and checks for duplicates.
    Constants and variables are intentionally ignored.

    Args:
        tech_specs_dir: Directory containing tech spec markdown files
        index_filename: Index markdown filename to exclude
        output: Optional OutputBuilder for verbose output

    Returns:
        List of DetectedDefinition objects
    """
    definitions: List[DetectedDefinition] = []
    definitions_by_name: Dict[str, List[DetectedDefinition]] = {}
    content_cache = file_cache or FileContentCache()
    files_to_check = [
        f
        for f in tech_specs_dir.glob("*.md")
        if not is_in_dot_directory(f) and f.name != index_filename
    ]
    if output and output.verbose:
        msg = f"Scanning {len(files_to_check)} tech spec files for definitions..."
        output.add_verbose_line(msg)

    for md_file in files_to_check:
        try:
            content = content_cache.get_content(md_file)
            lines = content_cache.get_lines(md_file)
            headings = extract_headings(content, skip_code_blocks=True)
            go_blocks = find_go_code_blocks(content)
            for start_line, _end_line, code_content in go_blocks:
                block_lines = code_content.split("\n")
                block_defs = _definitions_from_block(
                    block_lines, code_content, start_line,
                    md_file=md_file, content=content, lines=lines, headings=headings,
                    tech_specs_dir=tech_specs_dir, output=output, content_cache=content_cache,
                )
                for defn in block_defs:
                    definitions.append(defn)
                    definitions_by_name.setdefault(defn.name, []).append(defn)

        except (OSError, UnicodeDecodeError, ValueError, RuntimeError) as e:
            if output:
                output.add_error_line(
                    ValidationIssue.create(
                        "Error reading file",
                        Path(md_file.name),
                        0,
                        0,
                        message=f"Could not read file: {e}",
                        severity="error",
                    ).format_message(no_color=output.no_color)
                )
            continue

    duplicate_groups = _build_duplicate_groups(definitions_by_name)
    if duplicate_groups and output:
        _emit_duplicate_errors(
            duplicate_groups, output, tech_specs_dir, content_cache
        )
    if output and output.verbose:
        output.add_verbose_line(f"Found {len(definitions)} total definitions")
        if duplicate_groups:
            total_duplicates = sum(len(defns) for defns in duplicate_groups.values())
            msg = (
                f"Found {len(duplicate_groups)} duplicate definition groups "
                f"({total_duplicates} total occurrences)"
            )
            output.add_verbose_line(msg)

    return definitions


def _populate_heading_context(
    definition: DetectedDefinition,
    content: str,
    lines: List[str],
    headings: List[tuple[str, int, int]],
    *,
    tech_specs_dir: Path,
    output: Optional[OutputBuilder],
    content_cache: FileContentCache,
) -> None:
    heading_context = find_heading_before_line(
        content, definition.code_block_start_line, prefer_deepest=False
    )
    if heading_context:
        definition.heading = heading_context.heading_text
        definition.heading_level = heading_context.heading_level
        definition.heading_line = heading_context.heading_line
    else:
        definition.heading = ""
        definition.heading_level = 0
        definition.heading_line = 0

    if definition.heading_line > 0:
        parent_heading = None
        parent_level = None
        for heading_text, level, line_num in headings:
            if line_num < definition.heading_line and level < definition.heading_level:
                if parent_level is None or level >= parent_level:
                    parent_heading = heading_text
                    parent_level = level
        definition.parent_heading = parent_heading
        definition.parent_heading_level = parent_level

        section_start = definition.heading_line - 1
        section_end = len(lines)
        for heading_text, level, line_num in headings:
            if line_num > definition.heading_line and level <= definition.heading_level:
                section_end = line_num - 1
                break
        definition.section_content = "\n".join(lines[section_start:section_end])

    canonical_file, canonical_heading, canonical_anchor = resolve_canonical_reference(
        definition,
        content,
        tech_specs_dir,
        output=output,
        file_cache=content_cache,
    )
    definition.canonical_file = canonical_file
    definition.canonical_heading = canonical_heading
    definition.canonical_anchor = canonical_anchor
