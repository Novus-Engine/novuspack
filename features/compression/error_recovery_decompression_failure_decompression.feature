@domain:compression @m2 @REQ-COMPR-063 @spec(api_package_compression.md#1222-error-recovery-decompression-failure)
Feature: Error Recovery Decompression Failure

  @REQ-COMPR-063 @happy
  Scenario: Package remains compressed after decompression failure
    Given an open NovusPack package
    And package is compressed
    And a decompression operation fails
    When decompression failure occurs
    Then package remains compressed
    And package compression state is unchanged
    And package is still compressed

  @REQ-COMPR-063 @happy
  Scenario: Original compressed data is preserved after decompression failure
    Given an open NovusPack package
    And package is compressed
    And package has original compressed data
    And a decompression operation fails
    When decompression failure occurs
    Then original compressed data is preserved
    And compressed data is intact
    And no data loss occurs

  @REQ-COMPR-063 @happy
  Scenario: Recovery or backup can be attempted after decompression failure
    Given an open NovusPack package
    And package is compressed
    And a decompression operation fails
    When recovery is attempted
    Then recovery can be attempted
    And backup can be used if available
    And package state allows recovery attempts

  @REQ-COMPR-063 @happy
  Scenario: Error recovery ensures compressed data integrity after decompression failure
    Given an open NovusPack package
    And package is compressed
    And a decompression operation fails
    When decompression failure occurs
    Then compressed data integrity is maintained
    And compressed data remains usable
    And no corruption occurs in compressed data
