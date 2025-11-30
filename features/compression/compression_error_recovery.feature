@domain:compression @m2 @REQ-COMPR-060 @spec(api_package_compression.md#122-error-recovery)
Feature: Compression Error Recovery

  @REQ-COMPR-060 @happy
  Scenario: Compression failure recovery preserves package state
    Given a package compression operation
    When compression operation fails
    Then package remains in original state
    And no partial compression state exists
    And package can retry with different compression type

  @REQ-COMPR-060 @happy
  Scenario: Compression failure allows retry with different compression type
    Given a compression operation that failed
    When retry is attempted with different compression type
    Then package state is preserved from original failure
    And new compression attempt can proceed
    And package remains usable

  @REQ-COMPR-060 @happy
  Scenario: Decompression failure recovery preserves compressed data
    Given a package decompression operation
    When decompression operation fails
    Then package remains compressed
    And original compressed data is preserved
    And recovery or backup can be attempted

  @REQ-COMPR-060 @happy
  Scenario: Decompression failure allows recovery attempts
    Given a decompression operation that failed
    When recovery is attempted
    Then original compressed data is still available
    And backup can be used if available
    And package data integrity is maintained

  @REQ-COMPR-060 @error
  Scenario: Error recovery handles compression errors gracefully
    Given compression operations that may fail
    When errors occur during compression
    Then error recovery mechanisms activate
    And package state is protected
    And error details are available for diagnosis
