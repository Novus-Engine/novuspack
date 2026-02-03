#!/usr/bin/env python3
"""Tests for Go defs index scoring helpers."""

from __future__ import annotations

import unittest

from ._go_defs_index_models import DetectedDefinition
from ._go_defs_index_scoring_domain import _domain_from_file_pattern, detect_definition_domain
from ._go_defs_index_scoring_rules_core_base import (
    ScoringContext,
    score_exact_name_match,
)
from ._go_defs_index_scoring_rules_core_domain import score_type_name_patterns
from ._go_defs_index_scoring_rules_penalties import score_error_context_types
from ._go_defs_index_scoring_rules_sections import score_file_patterns


def _make_definition(name: str, kind: str, file_name: str) -> DetectedDefinition:
    return DetectedDefinition(
        name=name,
        kind=kind,
        file=file_name,
        code_block_start_line=1,
        code_block_content="",
        raw_name=name,
    )


def _make_context(definition: DetectedDefinition, section: str) -> ScoringContext:
    return ScoringContext(
        definition=definition,
        section=section,
        all_sections=set(),
        section_valid_types=None,
        section_lower=section.lower(),
        name_lower=definition.name.lower(),
        heading_lower=definition.heading.lower() if definition.heading else "",
        content_lower=definition.section_content.lower(),
        detected_domain=None,
    )


class ScoringDomainTests(unittest.TestCase):
    """Tests for domain detection and file pattern helpers."""

    def test_domain_from_file_pattern(self) -> None:
        cases = {
            "api_file_mgmt_error.md": "error",
            "api_file_mgmt_compression.md": "compression",
            "api_file_mgmt_queries.md": "fileentry",
            "api_basic_operations.md": "package",
            "package_file_format.md": "package",
            "file_type_system.md": "filetype",
        }
        for file_name, expected in cases.items():
            with self.subTest(file_name=file_name):
                self.assertEqual(_domain_from_file_pattern(file_name), expected)

    def test_detect_definition_domain_fallback(self) -> None:
        definition = _make_definition("WorkerPool", "type", "api_misc.md")
        detected = detect_definition_domain(definition, definition.name.lower())
        self.assertEqual(detected, "concurrency")

    def test_score_file_patterns(self) -> None:
        definition = _make_definition("FileEntry", "type", "api_file_mgmt_removal.md")
        ctx = _make_context(definition, "4. FileEntry Types")
        score, reasoning = score_file_patterns(ctx)
        self.assertGreater(score, 0.0)
        self.assertTrue(any("File pattern match" in r for r in reasoning))

    def test_score_type_name_patterns(self) -> None:
        definition = _make_definition("CompressionConfig", "type", "api_package_compression.md")
        ctx = _make_context(definition, "6. Compression Types")
        score, reasoning = score_type_name_patterns(ctx)
        self.assertGreater(score, 0.0)
        self.assertTrue(any("Type pattern" in r for r in reasoning))

    def test_score_error_context_types(self) -> None:
        definition = _make_definition("CompressionErrorContext", "type", "api_core.md")
        ctx = _make_context(definition, "13. Error Types")
        score, reasoning = score_error_context_types(ctx)
        self.assertGreater(score, 0.0)
        self.assertTrue(any("ErrorContext type" in r for r in reasoning))

    def test_score_exact_name_match_interface_types(self) -> None:
        definition = _make_definition("Package", "type", "api_core.md")
        ctx = _make_context(definition, "1. Package Interface Types")
        score, reasoning = score_exact_name_match(ctx)
        self.assertGreater(score, 0.0)
        self.assertTrue(any("Exact type name match" in r for r in reasoning))


if __name__ == "__main__":
    unittest.main()
