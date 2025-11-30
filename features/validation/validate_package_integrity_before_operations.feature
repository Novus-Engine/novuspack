@domain:validation @m2 @spec(file_validation.md#12-file-content-validation)
Feature: Validate Package Integrity Before Operations

  @REQ-VALID-002 @error
  Scenario: Integrity validation gates risky operations
    Given a NovusPack package
    And a package with integrity issues
    When risky operation is attempted
    Then operation is blocked by integrity validation
    And integrity validation prevents unsafe operations
    And error indicates integrity issues

  @REQ-VALID-002 @happy
  Scenario: Package integrity validation validates package state
    Given a NovusPack package
    And an open NovusPack package
    When package integrity validation is performed
    Then package structure is validated
    And package data integrity is validated
    And package metadata integrity is validated
    And integrity validation confirms package validity

  @REQ-VALID-002 @error
  Scenario: Package integrity validation detects integrity issues
    Given a NovusPack package
    And a package with integrity problems
    When package integrity validation is performed
    Then integrity issues are detected
    And integrity issues are reported
    And integrity validation prevents unsafe operations

  @REQ-VALID-002 @error
  Scenario: Integrity validation blocks operations on invalid packages
    Given a NovusPack package
    And a package with integrity issues
    When write operation is attempted
    Then operation is blocked
    And error indicates integrity validation failure
    And error provides details about integrity issues
