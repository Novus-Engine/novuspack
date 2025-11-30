@domain:core @m2 @REQ-CORE-042 @spec(api_core.md#91-general-metadata-operations)
Feature: Core Metadata Operations

  @REQ-CORE-042 @happy
  Scenario: General metadata operations provide metadata access
    Given an open NovusPack package
    And a valid context
    When general metadata operations are used
    Then metadata can be accessed
    And metadata can be read
    And metadata operations provide package information

  @REQ-CORE-042 @happy
  Scenario: General metadata operations support metadata retrieval
    Given an open NovusPack package
    And a valid context
    When metadata is retrieved
    Then metadata information is returned
    And metadata includes package details
    And metadata can be used for package inspection

  @REQ-CORE-042 @happy
  Scenario: General metadata operations provide comprehensive package information
    Given an open NovusPack package
    And a valid context
    When metadata is accessed
    Then package information is available
    And package details can be inspected
    And metadata includes package structure

  @REQ-CORE-042 @happy
  Scenario: General metadata operations integrate with context
    Given an open NovusPack package
    And a valid context
    When metadata operations are performed
    Then context cancellation is supported
    And context timeout handling is supported
    And operations respect context cancellation

  @REQ-CORE-042 @error
  Scenario: General metadata operations handle package not open errors
    Given a closed NovusPack package
    And a valid context
    When metadata operations are attempted
    Then a structured error is returned
    And error indicates package is not open
    And error follows structured error format

  @REQ-CORE-042 @happy
  Scenario: General metadata operations support package inspection
    Given an open NovusPack package
    And a valid context
    When package is inspected via metadata
    Then package structure can be examined
    And package contents can be analyzed
    And package state can be determined
