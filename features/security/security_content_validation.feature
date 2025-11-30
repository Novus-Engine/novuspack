@domain:security @m2 @REQ-SEC-034 @spec(security.md#421-content-validation)
Feature: Security Content Validation

  @REQ-SEC-034 @happy
  Scenario: Content validation supports empty files
    Given an open NovusPack package
    And an empty file
    When content validation is performed
    Then empty files are supported
    And empty files are valid
    And validation passes for empty files

  @REQ-SEC-034 @happy
  Scenario: Content validation rejects nil data
    Given an open NovusPack package
    And file with nil data
    When content validation is performed
    Then nil data is prohibited
    And nil data is rejected
    And validation error is returned

  @REQ-SEC-034 @happy
  Scenario: Content validation validates content integrity before addition
    Given an open NovusPack package
    And file content
    When file is added to package
    Then content integrity is validated before addition
    And validation ensures content correctness
    And invalid content is rejected

  @REQ-SEC-034 @happy
  Scenario: Content validation validates and normalizes file paths
    Given an open NovusPack package
    And file with path
    When content validation is performed
    Then path validation is performed
    And paths are normalized
    And invalid paths are rejected

  @REQ-SEC-034 @happy
  Scenario: Content validation validates file content format
    Given an open NovusPack package
    And file with content
    When content validation is performed
    Then file content format is validated
    And format validation ensures correctness
    And invalid formats are rejected

  @REQ-SEC-034 @error
  Scenario: Content validation fails with corrupted content
    Given an open NovusPack package
    And file with corrupted content
    When content validation is performed
    Then structured validation error is returned
    And error indicates content corruption
