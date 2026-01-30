"""
Compatibility exports for core scoring rules.
"""

from __future__ import annotations

from lib.go_defs_index._go_defs_index_scoring_rules_core_base import (
    ScoringContext,
    _error_section_flags,
    _is_core_package_type,
    score_current_section,
    score_exact_name_match,
    score_function_type_interaction,
    score_implementation_mapping,
    score_kind_blockers,
    score_kind_positive_match,
    score_strict_kind_matching,
)
from lib.go_defs_index._go_defs_index_scoring_rules_core_domain import (
    score_constructor_functions,
    score_domain_match,
    score_domain_type_subsection,
    score_keyword_comment_matching,
)

__all__ = [
    "ScoringContext",
    "_error_section_flags",
    "_is_core_package_type",
    "score_current_section",
    "score_exact_name_match",
    "score_function_type_interaction",
    "score_implementation_mapping",
    "score_kind_blockers",
    "score_kind_positive_match",
    "score_strict_kind_matching",
    "score_constructor_functions",
    "score_domain_match",
    "score_domain_type_subsection",
    "score_keyword_comment_matching",
]
