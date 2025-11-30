@domain:compression @m2 @REQ-COMPR-061 @spec(api_package_compression.md#1221-error-recovery-compression-failure)
Feature: Error Recovery Compression Failure

  @REQ-COMPR-061 @happy
  Scenario: Package remains in original state after compression failure
    Given an open NovusPack package
    And package is uncompressed
    And a compression operation fails
    When compression failure occurs
    Then package remains in original uncompressed state
    And package state is unchanged
    And no partial compression state exists

  @REQ-COMPR-061 @happy
  Scenario: No partial compression state after compression failure
    Given an open NovusPack package
    And package is uncompressed
    And a compression operation fails
    When compression failure occurs
    Then no partial compression state is left
    And package is fully uncompressed
    And package integrity is maintained

  @REQ-COMPR-061 @happy
  Scenario: Compression can be retried with different compression type after failure
    Given an open NovusPack package
    And package is uncompressed
    And compression operation fails with one compression type
    When retry is attempted with different compression type
    Then retry is possible
    And different compression type can be used
    And package state allows retry

  @REQ-COMPR-061 @happy
  Scenario: Error recovery ensures package consistency after compression failure
    Given an open NovusPack package
    And package is uncompressed
    And a compression operation fails
    When compression failure occurs
    Then package consistency is maintained
    And package can continue to be used
    And no corruption occurs
