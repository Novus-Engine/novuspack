@domain:basic_ops @m2 @REQ-API_BASIC-024 @spec(api_basic_operations.md#211-constants)
Feature: Package format constants

  @REQ-API_BASIC-024 @happy
  Scenario: Package format constants are defined and accessible
    Given NovusPack package format constants
    When constants are accessed
    Then magic number constant is available
    And format version constant is available
    And header size constant is available
    And constants match package format specification

  @REQ-API_BASIC-024 @happy
  Scenario: Constants provide consistent package format values
    Given package format constants
    When constants are used in package operations
    Then magic number matches package header format
    And format version indicates supported version
    And header size matches specification
    And constants are read-only

  @REQ-API_BASIC-024 @happy
  Scenario: Constants can be used for validation
    Given package format constants
    When validating package format
    Then magic number constant can verify package identity
    And format version constant can verify compatibility
    And constants provide reference values for validation
