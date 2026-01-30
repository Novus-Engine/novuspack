@domain:core @m2 @REQ-CORE-088 @spec(api_core.md#1164-listfiles-behavior) @spec(api_core.md#packagereaderlistfiles-behavior)
Feature: ListFiles behavior defines sorting, stability, and mutation handling

  @REQ-CORE-088 @happy
  Scenario: ListFiles behavior is consistent for sorting and stability
    Given an opened package
    When ListFiles is called multiple times without package mutation
    Then results are sorted consistently
    And results are stable across calls
    And mutation handling follows the specified behavior
