@domain:metadata @m2 @REQ-META-082 @spec(api_metadata.md#641-metadata-only-package-validation)
Feature: Metadata-Only Package Validation

  @REQ-META-082 @happy
  Scenario: Metadata-only package validation performs comprehensive validation
    Given a NovusPack package
    And a metadata-only package
    When ValidateMetadataOnlyPackage is called
    Then FileCount is checked and must be 0
    And HasSpecialMetadataFiles is checked and must be true
    And all special metadata files are validated
    And malicious metadata patterns are checked
    And signature scope is verified to include all metadata
    And metadata consistency is verified
    And package structure is validated

  @REQ-META-082 @happy
  Scenario: Metadata-only package validation checks file count
    Given a NovusPack package
    And a metadata-only package
    When validation checks file count
    Then FileCount must be 0
    And validation fails if regular files exist

  @REQ-META-082 @happy
  Scenario: Metadata-only package validation validates special files
    Given a NovusPack package
    And a metadata-only package
    When validation validates special files
    Then all special metadata files are validated
    And file formats are verified
    And file content is validated

  @REQ-META-082 @happy
  Scenario: Metadata-only package validation checks for malicious patterns
    Given a NovusPack package
    And a metadata-only package
    When validation checks for malicious patterns
    Then metadata is scanned for injection patterns
    And malicious patterns are detected
    And validation fails if malicious patterns are found

  @REQ-META-082 @error
  Scenario: Metadata-only package validation detects violations
    Given a NovusPack package
    When validation detects violations
    Then appropriate errors are returned
    And violations are clearly reported
    And errors follow structured error format
