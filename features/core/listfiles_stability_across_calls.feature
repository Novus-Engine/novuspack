@domain:core @m2 @REQ-CORE-051 @spec(api_core.md#114-listfiles-method-contract) @spec(api_core.md#116-listfiles-method-contract)
Feature: ListFiles results are stable across calls when package state unchanged

  @REQ-CORE-051 @happy
  Scenario: ListFiles is stable across repeated calls
    Given an opened package
    And the package state does not change between calls
    When ListFiles is called multiple times
    Then each call returns the same ordered results
    And the ordering is stable for an unchanged package state
