@domain:basic_ops @m2 @REQ-API_BASIC-113 @spec(api_basic_operations.md#22-import-patterns)
Feature: Import patterns

  @REQ-API_BASIC-113 @happy
  Scenario: Import patterns define recommended and direct subpackage import usage
    Given a consumer writing Go code against the API
    When choosing how to import NovusPack
    Then a recommended import pattern is defined for common use cases
    And direct subpackage imports are defined for package-specific functionality
    And import patterns minimize consumer boilerplate
    And import patterns avoid circular dependency hazards
    And import guidance is consistent with the package organization

