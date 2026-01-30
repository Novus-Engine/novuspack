@domain:core @m2 @REQ-CORE-050 @spec(api_core.md#114-listfiles-method-contract) @spec(api_core.md#116-listfiles-method-contract)
Feature: ListFiles returns results sorted by PrimaryPath alphabetically

  @REQ-CORE-050 @happy
  Scenario: ListFiles output is sorted by PrimaryPath
    Given an opened package containing multiple files
    When ListFiles is called
    Then returned FileInfo entries are sorted by PrimaryPath alphabetically
    And sorting uses PrimaryPath as the primary sort key
    And callers receive a deterministic ordering for the same package state
