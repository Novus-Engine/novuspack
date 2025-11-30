@domain:basic_ops @m1 @REQ-API_BASIC-008 @spec(api_basic_operations.md#52-open-with-validation)
Feature: Open package with validation

  @happy
  Scenario: OpenWithValidation opens and validates package successfully
    Given a valid NovusPack package file
    When OpenWithValidation is called
    Then package is opened
    And full package validation is performed
    And package structure is validated
    And checksums are verified
    And signatures are validated if present
    And package is ready for operations

  @happy
  Scenario: OpenWithValidation validates package structure
    Given a NovusPack package file
    When OpenWithValidation is called
    Then package header format is validated
    And file entry structure is validated
    And data section integrity is validated
    And file index consistency is validated

  @happy
  Scenario: OpenWithValidation verifies checksums
    Given a NovusPack package file with calculated checksums
    When OpenWithValidation is called
    Then file checksums are verified
    And package CRC is verified if present
    And checksum mismatches are detected

  @happy
  Scenario: OpenWithValidation validates signatures if present
    Given a signed NovusPack package file
    When OpenWithValidation is called
    Then all signatures are validated
    And signature integrity is verified
    And signature metadata is validated

  @error
  Scenario: OpenWithValidation fails if validation fails
    Given a corrupted NovusPack package file
    When OpenWithValidation is called
    Then a structured validation error is returned
    And package is not opened
    And detailed error information is provided

  @error
  Scenario: OpenWithValidation fails if checksums don't match
    Given a NovusPack package file with corrupted data
    When OpenWithValidation is called
    Then a structured corruption error is returned
    And checksum mismatch is identified
    And package is not opened

  @error
  Scenario: OpenWithValidation fails if signatures are invalid
    Given a NovusPack package file with invalid signatures
    When OpenWithValidation is called
    Then a structured signature error is returned
    And signature validation failures are identified
    And package is not opened

  @error
  Scenario: OpenWithValidation respects context cancellation
    Given a NovusPack package file
    And a cancelled context
    When OpenWithValidation is called
    Then a structured context error is returned
    And validation is cancelled
