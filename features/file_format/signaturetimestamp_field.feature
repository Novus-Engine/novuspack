@domain:file_format @m2 @v2 @REQ-FILEFMT-068 @spec(package_file_format.md#723-signaturetimestamp-field)
Feature: SignatureTimestamp Field

  @REQ-FILEFMT-068 @happy
  Scenario: SignatureTimestamp field stores signature timestamp
    Given a signature block
    When SignatureTimestamp field is examined
    Then SignatureTimestamp is a 32-bit unsigned integer
    And SignatureTimestamp stores timestamp when signature was created
    And timestamp indicates signature creation time

  @REQ-FILEFMT-068 @happy
  Scenario: SignatureTimestamp uses Unix nanoseconds format
    Given a signature block
    When SignatureTimestamp is examined
    Then timestamp is Unix timestamp in nanoseconds format
    And timestamp represents nanoseconds since Unix epoch
    And timestamp format enables precise time tracking

  @REQ-FILEFMT-068 @happy
  Scenario: SignatureTimestamp supports valid range
    Given a signature block
    When SignatureTimestamp is examined
    Then timestamp range is 0 to 4294967295
    And timestamp represents valid Unix nanoseconds
    And timestamp is within supported range

  @REQ-FILEFMT-068 @happy
  Scenario: SignatureTimestamp records signature creation time
    Given a signature block
    And signature is created at known time
    When SignatureTimestamp is examined
    Then timestamp matches signature creation time
    And timestamp enables signature time verification
    And timestamp supports temporal validation

  @REQ-FILEFMT-068 @error
  Scenario: SignatureTimestamp with invalid range is rejected
    Given a signature block
    And SignatureTimestamp exceeds valid range
    When signature is validated
    Then validation may flag invalid timestamp
    And timestamp range violations are detected
