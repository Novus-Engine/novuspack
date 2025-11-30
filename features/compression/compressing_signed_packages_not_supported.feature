@domain:compression @m2 @REQ-COMPR-028 @spec(api_package_compression.md#1021-not-supported)
Feature: Compressing Signed Packages Not Supported

  @REQ-COMPR-028 @error
  Scenario: Signed packages cannot be compressed
    Given an open NovusPack package
    And package has signatures
    And a valid context
    And a compression type
    When compression operation is attempted
    Then compression is not supported
    And error is returned
    And error indicates signed packages cannot be compressed

  @REQ-COMPR-028 @error
  Scenario: Compression of signed packages returns error
    Given an open NovusPack package
    And package has signatures
    And a valid context
    When compression operation is attempted
    Then security error is returned
    And error indicates operation not supported
    And error follows structured error format

  @REQ-COMPR-028 @happy
  Scenario: Workflow for compressing previously signed packages
    Given an open NovusPack package
    And package has signatures
    When proper workflow is followed
    Then first step is to clear signatures
    And second step is to compress package
    And third step is to re-sign if needed
    And workflow handles signed packages correctly
