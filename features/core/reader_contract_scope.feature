@domain:core @m2 @REQ-CORE-186 @spec(api_core.md#1111-reader-contract-scope)
Feature: Reader contract scope defines PackageReader method assumptions

  @REQ-CORE-186 @happy
  Scenario: Reader contract scope defines method assumptions
    Given an opened package and PackageReader interface
    When PackageReader methods are invoked
    Then the scope defines the assumptions for each method
    And callers can rely on the documented assumptions
    And the behavior matches the reader contract scope specification
