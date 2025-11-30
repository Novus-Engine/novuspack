@domain:security @m2 @REQ-SEC-031 @spec(security.md#411-comprehensive-validation)
Feature: Comprehensive Package Validation

  @REQ-SEC-031 @happy
  Scenario: Comprehensive validation validates package format
    Given an open NovusPack package
    And a valid context
    When Validate is called
    Then package format is validated
    And format validation is comprehensive
    And format errors are detected

  @REQ-SEC-031 @happy
  Scenario: Comprehensive validation validates package structure
    Given an open NovusPack package
    And a valid context
    When Validate is called
    Then package structure is validated
    And structure validation is comprehensive
    And structure errors are detected

  @REQ-SEC-031 @happy
  Scenario: Comprehensive validation validates package integrity
    Given an open NovusPack package
    And a valid context
    When ValidateIntegrity is called
    Then package integrity is validated using checksums
    And integrity validation is comprehensive
    And integrity errors are detected

  @REQ-SEC-031 @happy
  Scenario: Comprehensive validation validates all signatures
    Given an open NovusPack package
    And package with multiple signatures
    And a valid context
    When ValidateAllSignatures is called
    Then all signatures are validated
    And signature validation is comprehensive
    And signature errors are detected

  @REQ-SEC-031 @happy
  Scenario: Comprehensive validation provides security status
    Given an open NovusPack package
    And a valid context
    When GetSecurityStatus is called
    Then comprehensive security status information is returned
    And status includes signature information
    And status includes validation results
    And status includes checksum status
    And status includes security level
    And status includes error reporting

  @REQ-SEC-007 @REQ-SEC-011 @error
  Scenario: Comprehensive validation respects context cancellation
    Given an open NovusPack package
    And a cancelled context
    When validation operation is called
    Then operation stops
    And structured context error is returned
    And error type is context cancellation
