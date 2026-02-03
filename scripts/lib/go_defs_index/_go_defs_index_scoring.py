"""
Scoring logic for Go definitions index placement.
"""

from __future__ import annotations

from typing import Dict, List, Optional, Set, Tuple

from lib.go_defs_index._go_defs_index_models import DetectedDefinition
from lib.go_defs_index._go_defs_index_scoring_domain import detect_definition_domain
from lib.go_defs_index._go_defs_index_scoring_rules_core import (
    ScoringContext,
    score_constructor_functions,
    score_current_section,
    score_domain_match,
    score_domain_type_subsection,
    score_exact_name_match,
    score_function_type_interaction,
    score_implementation_mapping,
    score_keyword_comment_matching,
    score_kind_blockers,
    score_kind_positive_match,
    score_strict_kind_matching,
)
from lib.go_defs_index._go_defs_index_scoring_rules_core_domain import (
    score_type_name_patterns,
)
from lib.go_defs_index._go_defs_index_scoring_rules_methods import (
    score_file_entry_method_categories,
    score_method_name_preferences,
    score_method_patterns,
    score_method_type_classification,
)
from lib.go_defs_index._go_defs_index_scoring_rules_penalties import (
    score_error_helper_functions,
    score_error_methods,
    score_error_domain_match,
    score_error_context_types,
    score_error_context_domain_mismatch,
    score_general_heuristics,
    score_hash_optional_types,
    score_kind_section_map,
    score_metadata_tag_helpers,
    score_package_config_preference,
    score_readonly_package_interface,
    score_readonly_type_preference,
    score_streaming_helper_functions,
    score_streaming_helper_mismatch,
    score_type_operation_penalty,
)
from lib.go_defs_index._go_defs_index_scoring_rules_type_keywords import (
    score_create_options_preference,
    score_error_type_keywords,
    score_file_entry_type_keywords,
    score_file_info_preference,
    score_generic_helper_functions,
    score_generic_type_keywords,
    score_metadata_type_keywords,
    score_other_type_helper_functions,
    score_other_types_suffix,
    score_recovery_file_header_preference,
    score_signature_type_keywords,
    score_generic_core_type_preference,
    score_package_error_constructor,
    score_signature_comment_helpers,
    score_security_error_context_types,
    score_package_open_helpers,
    score_package_read_header_helpers,
    score_metadata_destpath_helpers,
    score_package_helper_overrides,
)
from lib.go_defs_index._go_defs_index_scoring_rules_sections import (
    score_camelcase_match,
    score_comment_domain_match,
    score_content_keyword_match,
    score_file_patterns,
    score_heading_match,
    score_parent_heading_match,
    score_prose_keyword_match,
    score_subsection_keyword_match,
)


def calculate_confidence_score(
    definition: DetectedDefinition,
    section: str,
    all_sections: Set[str],
    section_valid_types: Optional[Dict[str, Set[str]]] = None,
) -> Tuple[float, List[str]]:
    """
    Calculate confidence score (0.0 to 1.0) and reasoning.
    """
    score = 0.0
    reasoning: List[str] = ["Base score: 0%"]

    ctx = ScoringContext(
        definition=definition,
        section=section,
        all_sections=all_sections,
        section_valid_types=section_valid_types,
        section_lower=section.lower(),
        name_lower=definition.name.lower(),
        heading_lower=definition.heading.lower() if definition.heading else "",
        content_lower=definition.section_content.lower(),
        detected_domain=detect_definition_domain(definition, definition.name.lower()),
    )

    strict_score, strict_reasoning, blocked = score_strict_kind_matching(ctx)
    if blocked:
        return strict_score, strict_reasoning
    score += strict_score
    reasoning.extend(strict_reasoning)

    delta, delta_reasoning = score_function_type_interaction(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_exact_name_match(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_implementation_mapping(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning, kind_mismatch = score_kind_blockers(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_kind_positive_match(ctx, kind_mismatch)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_keyword_comment_matching(ctx, kind_mismatch)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_current_section(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_constructor_functions(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_domain_match(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_type_name_patterns(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_error_context_types(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_error_context_domain_mismatch(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_domain_type_subsection(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_method_patterns(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_file_patterns(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_heading_match(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_camelcase_match(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_parent_heading_match(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_comment_domain_match(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_prose_keyword_match(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_content_keyword_match(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_subsection_keyword_match(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_method_type_classification(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_method_name_preferences(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_file_entry_method_categories(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_hash_optional_types(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_error_domain_match(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_error_methods(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_metadata_tag_helpers(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_error_helper_functions(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_package_config_preference(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_other_types_suffix(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_generic_type_keywords(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_metadata_type_keywords(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_streaming_helper_functions(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_streaming_helper_mismatch(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_generic_helper_functions(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_create_options_preference(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_signature_type_keywords(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_error_type_keywords(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_file_entry_type_keywords(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_file_info_preference(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_recovery_file_header_preference(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_other_type_helper_functions(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_generic_core_type_preference(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_package_error_constructor(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_signature_comment_helpers(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_security_error_context_types(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_package_open_helpers(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_package_read_header_helpers(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_metadata_destpath_helpers(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_package_helper_overrides(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_readonly_type_preference(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_readonly_package_interface(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_type_operation_penalty(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_kind_section_map(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    delta, delta_reasoning = score_general_heuristics(ctx)
    score += delta
    reasoning.extend(delta_reasoning)

    score = max(0.0, score)
    return score, reasoning
