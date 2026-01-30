@domain:basic_ops @m2 @REQ-API_BASIC-189 @spec(api_basic_operations.md#2117-internal-package-purpose)
Feature: internal package purpose

  @REQ-API_BASIC-189 @happy
  Scenario: internal package purpose is helper functions not part of public API
    Given an internal helper package
    When the public API boundary is enforced
    Then internal helpers are not part of the public API contract
    And internal helpers can change without breaking consumers
    And internal helpers enable implementation reuse
    And internal helpers avoid cluttering public packages
    And internal package purpose aligns with Go internal package rules

