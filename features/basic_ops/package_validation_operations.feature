@domain:basic_ops @m1 @REQ-API_BASIC-010 @spec(api_basic_operations.md#71-package-validation)
Feature: Package validation operations

  @happy
  Scenario: Validate performs comprehensive package validation
    Given an open NovusPack package
    When Validate is called
    Then package header format is validated
    And package version is validated
    And file entry structure is validated
    And data section integrity is validated
    And checksums are verified
    And signatures are validated if present

  @happy
  Scenario: Validate checks package header format
    Given an open NovusPack package
    When Validate is called
    Then header magic number is verified
    And header format version is verified
    And header structure is validated
    And header field consistency is checked

  @happy
  Scenario: Validate checks file entry structure
    Given an open NovusPack package
    When Validate is called
    Then all file entries are validated
    And file entry format is verified
    And file entry consistency is checked
    And path entries are validated

  @happy
  Scenario: Validate verifies data section integrity
    Given an open NovusPack package
    When Validate is called
    Then file data checksums are verified
    And data section structure is validated
    And file data consistency is checked

  @happy
  Scenario: Validate verifies checksums
    Given an open NovusPack package with calculated checksums
    When Validate is called
    Then file checksums are verified against data
    And package CRC is verified if present
    And checksum mismatches are detected

  @happy
  Scenario: Validate validates signatures if present
    Given a signed open NovusPack package
    When Validate is called
    Then all signatures are validated
    And signature integrity is verified
    And signature metadata is validated

  @error
  Scenario: Validate fails if package is not open
    Given a closed NovusPack package
    When Validate is called
    Then a structured validation error is returned

  @error
  Scenario: Validate fails if package format is invalid
    Given an open NovusPack package with invalid format
    When Validate is called
    Then a structured validation error is returned
    And detailed error information is provided

  @error
  Scenario: Validate fails if checksums don't match
    Given an open NovusPack package with corrupted data
    When Validate is called
    Then a structured corruption error is returned
    And checksum mismatches are identified

  @error
  Scenario: Validate respects context cancellation
    Given an open NovusPack package
    And a cancelled context
    When Validate is called
    Then a structured context error is returned
    And validation is cancelled
