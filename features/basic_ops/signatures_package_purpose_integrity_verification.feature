@domain:basic_ops @m2 @REQ-API_BASIC-184 @spec(api_basic_operations.md#21141-signatures-package-purpose)
Feature: signatures package purpose

  @REQ-API_BASIC-184 @happy
  Scenario: signatures package purpose is integrity verification of packages
    Given signed package support
    When signature structures are used
    Then the signatures package purpose is to support package integrity verification
    And signature structures enable detection of signed package state
    And signature workflows integrate with signing and validation operations
    And signature purpose aligns with immutability and secure signing requirements
    And signature behavior is consistent across package lifecycle operations

