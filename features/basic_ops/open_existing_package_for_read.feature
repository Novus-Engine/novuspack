@domain:basic_ops @m1 @REQ-API_BASIC-002 @spec(api_basic_operations.md#51-open-method)
Feature: Open existing package for read

  Scenario: Open a valid package
    Given a valid NovusPack file on disk
    When OpenPackage is invoked
    Then the package opens successfully

  Scenario: Invalid format returns structured error
    Given a corrupted or non NovusPack file
    When OpenPackage is invoked
    Then a structured error is returned

  @REQ-API_BASIC-016 @REQ-API_BASIC-017
  Scenario: OpenPackage validates path parameter
    Given an invalid path (empty or whitespace-only)
    When OpenPackage is called with invalid path
    Then structured validation error is returned
    And error indicates invalid path

  @REQ-API_BASIC-017 @REQ-API_BASIC-019
  Scenario: OpenPackage respects context cancellation
    Given a valid NovusPack file on disk
    And a cancelled context
    When OpenPackage is called
    Then structured context error is returned
    And error type is context cancellation

  @REQ-API_BASIC-017 @REQ-API_BASIC-019
  Scenario: OpenPackage respects context timeout
    Given a valid NovusPack file on disk
    And a context with timeout
    When OpenPackage is called
    And operation exceeds timeout
    Then structured context error is returned
    And error type is context timeout
