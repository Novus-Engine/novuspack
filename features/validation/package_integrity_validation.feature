@domain:validation @m2 @REQ-VALID-003 @spec(file_validation.md#12-file-content-validation)
Feature: Package integrity validation

  @happy
  Scenario: Package integrity validation checks all components
    Given an open package
    When package integrity validation is performed
    Then header is validated
    And file entries are validated
    And file data is validated
    And checksums are verified
    And signatures are validated if present

  @error
  Scenario: Integrity validation detects corruption
    Given a corrupted package
    When integrity validation is performed
    Then structured corruption error is returned
    And corruption location is identified

  @error
  Scenario: Integrity validation detects checksum mismatches
    Given a package with checksum mismatch
    When integrity validation is performed
    Then structured corruption error is returned
    And checksum mismatch is identified
