@domain:metadata @m2 @REQ-META-079 @spec(api_metadata.md#633-package-integrity)
Feature: Package Integrity

  @REQ-META-079 @happy
  Scenario: Package integrity ensures metadata-only package integrity
    Given a NovusPack package
    And a metadata-only package
    When package integrity is verified
    Then size validation requires enhanced validation for very small packages
    And structure validation ensures valid package structure without content
    And metadata consistency verification is performed

  @REQ-META-079 @happy
  Scenario: Size validation requires enhanced validation
    Given a NovusPack package
    And a metadata-only package
    When size validation is performed
    Then very small packages require enhanced validation
    And size validation detects suspiciously small packages
    And enhanced validation rules are applied

  @REQ-META-079 @happy
  Scenario: Structure validation ensures valid package structure
    Given a NovusPack package
    And a metadata-only package
    When structure validation is performed
    Then package structure is validated without content
    And structure must be valid even with no content files
    And structure validation ensures package integrity

  @REQ-META-079 @happy
  Scenario: Metadata consistency verification is performed
    Given a NovusPack package
    And a metadata-only package
    When metadata consistency is verified
    Then metadata files are verified for internal consistency
    And metadata relationships are validated
    And consistency violations are detected

  @REQ-META-079 @error
  Scenario: Package integrity validation detects violations
    Given a NovusPack package
    When integrity violations are detected
    Then violations are reported
    And appropriate errors are returned
    And integrity violations are clearly identified
