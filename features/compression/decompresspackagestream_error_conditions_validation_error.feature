@domain:compression @m2 @REQ-COMPR-126 @spec(api_package_compression.md#524-decompresspackagestream-error-conditions)
Feature: DecompressPackageStream Error Conditions

  @REQ-COMPR-126 @error
  Scenario: Validation error when package is not compressed
    Given a NovusPack package that is not compressed
    When DecompressPackageStream is called
    Then validation error is returned
    And error indicates package is not compressed
    And error prevents invalid decompression operation

  @REQ-COMPR-126 @error
  Scenario: Validation error for invalid stream configuration
    Given a compressed NovusPack package
    And invalid StreamConfig is provided
    When DecompressPackageStream is called
    Then validation error is returned
    And error indicates invalid stream configuration
    And error provides configuration details

  @REQ-COMPR-126 @error
  Scenario: Compression error when decompression operation fails
    Given a compressed NovusPack package
    When decompression operation fails
    Then compression error is returned
    And error indicates decompression operation failure
    And error provides details about failure

  @REQ-COMPR-126 @error
  Scenario: Compression error for algorithm-specific decompression failures
    Given a compressed NovusPack package
    And decompression algorithm encounters specific failure
    When DecompressPackageStream is called
    Then compression error is returned
    And error indicates algorithm-specific failure
    And error provides algorithm context

  @REQ-COMPR-126 @error
  Scenario: I/O error when streaming operation fails
    Given a DecompressPackageStream operation
    And streaming operation fails
    When DecompressPackageStream is called
    Then I/O error is returned
    And error indicates streaming operation failure
    And error provides streaming details

  @REQ-COMPR-126 @error
  Scenario: I/O error when disk space is insufficient
    Given a DecompressPackageStream operation requiring disk space
    And insufficient disk space is available
    When DecompressPackageStream is called
    Then I/O error is returned
    And error indicates insufficient disk space
    And error provides disk space details

  @REQ-COMPR-126 @error
  Scenario: Context error on cancellation
    Given a DecompressPackageStream operation
    And a cancelled context
    When DecompressPackageStream is called with cancelled context
    Then context error is returned
    And error type is context cancellation

  @REQ-COMPR-126 @error
  Scenario: Context error on timeout
    Given a DecompressPackageStream operation
    And a context with expired timeout
    When DecompressPackageStream is called
    Then context error is returned
    And error type is context timeout
