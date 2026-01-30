@domain:core @m2 @REQ-CORE-049 @REQ-CORE-190 @REQ-CORE-191 @REQ-CORE-192 @REQ-CORE-193 @spec(api_core.md#2-package-path-semantics) @spec(api_core.md#214-unicode-normalization) @spec(api_core.md#215-path-length-limits) @spec(api_core.md#216-path-normalization-on-storage) @spec(api_core.md#23-path-display-and-extraction)
Feature: Package path semantics define package-internal path rules

  @REQ-CORE-049 @REQ-CORE-190 @REQ-CORE-191 @REQ-CORE-192 @REQ-CORE-193 @happy
  Scenario: Package-internal paths follow defined normalization and validation rules
    Given a package path provided by a caller
    When the path is normalized for storage
    Then unicode normalization rules are applied
    And path length limits are enforced
    And normalized paths follow the package path semantics
    And display and extraction use the defined display rules
