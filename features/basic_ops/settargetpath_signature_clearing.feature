@domain:basic_ops @m2 @REQ-API_BASIC-099 @spec(api_basic_operations.md#442-settargetpath-behavior)
Feature: SetTargetPath Signature Clearing

  @REQ-API_BASIC-099 @happy
  Scenario: SetTargetPath clears signatures when path differs
    Given a signed package with current path
    When SetTargetPath is called with a different path
    Then all signature information is cleared from memory
    And package becomes unsigned
    And new path is set for writing

  @REQ-API_BASIC-099 @happy
  Scenario: SetTargetPath clears signatures for signed package
    Given a package is signed
    And package has signature metadata
    When SetTargetPath is called with new path different from current
    Then signature data is removed
    And signature flags are cleared
    And package is marked as unsigned

  @REQ-API_BASIC-099 @happy
  Scenario: Signature clearing creates new unsigned package
    Given a signed package at original location
    When SetTargetPath changes path to new location
    Then writing creates new unsigned package at new location
    And original signed package remains unchanged
    And new package requires re-signing if signatures needed

  @REQ-API_BASIC-100 @happy
  Scenario: SetTargetPath preserves signatures when path equals current
    Given a signed package with current path
    When SetTargetPath is called with same path as current
    Then signature information is preserved
    And package remains signed
    And no signature clearing occurs

  @REQ-API_BASIC-100 @happy
  Scenario: SetTargetPath preserves signatures for identical path
    Given a signed package at "/path/to/package.nvpk"
    When SetTargetPath is called with "/path/to/package.nvpk"
    Then all signatures remain intact
    And signature metadata is not modified
    And package signed status unchanged

  @REQ-API_BASIC-099 @happy
  Scenario: SetTargetPath on unsigned package does not trigger clearing
    Given an unsigned package
    When SetTargetPath is called with different path
    Then no signature clearing occurs
    And package remains unsigned
    And target path is updated normally

  @REQ-API_BASIC-099 @happy
  Scenario: Signature clearing is required for immutability
    Given a signed package with cryptographic signatures
    When SetTargetPath changes location
    Then signature clearing preserves immutability principle
    And signed packages remain immutable at original location
    And new location creates new unsigned instance
