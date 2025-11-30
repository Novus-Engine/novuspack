@domain:security @m2 @REQ-SEC-039 @spec(security.md#512-access-control)
Feature: Access Control Mechanisms

  @REQ-SEC-039 @happy
  Scenario: Access control provides per-file encryption selection
    Given an open NovusPack package
    And a file entry
    When access control is configured
    Then per-file encryption type selection is available
    And encryption can be selected per file
    And encryption selection provides access control

  @REQ-SEC-039 @happy
  Scenario: Access control provides per-file compression selection
    Given an open NovusPack package
    And a file entry
    When access control is configured
    Then per-file compression algorithm selection is available
    And compression can be selected per file
    And compression selection provides access control

  @REQ-SEC-039 @happy
  Scenario: Access control provides file-specific security flags
    Given an open NovusPack package
    And a file entry
    When access control is configured
    Then file-specific security options are available
    And security restrictions can be set per file
    And security flags control file access

  @REQ-SEC-039 @happy
  Scenario: Access control provides metadata protection
    Given an open NovusPack package
    And sensitive metadata
    When access control is configured
    Then secure storage of sensitive metadata is provided
    And metadata protection controls access
    And protected metadata is secured

  @REQ-SEC-007 @error
  Scenario: Access control operations respect context cancellation
    Given an open NovusPack package
    And a cancelled context
    When access control operation is called
    Then structured context error is returned
    And error type is context cancellation
