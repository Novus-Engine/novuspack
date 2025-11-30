@domain:compression @m2 @REQ-COMPR-121 @spec(api_package_compression.md#514-compresspackagestream-error-conditions)
Feature: CompressPackageStream Error Conditions

  @REQ-COMPR-121 @error
  Scenario: Security error when package is already signed
    Given a NovusPack package that is already signed
    When CompressPackageStream is called
    Then security error is returned
    And error indicates package is already signed
    And error indicates signed packages cannot be compressed

  @REQ-COMPR-121 @error
  Scenario: Validation error for invalid compression type
    Given a NovusPack package
    And invalid compression type is specified
    When CompressPackageStream is called
    Then validation error is returned
    And error indicates invalid compression type
    And error specifies valid compression types

  @REQ-COMPR-121 @error
  Scenario: Validation error for invalid stream configuration
    Given a NovusPack package
    And invalid StreamConfig is provided
    When CompressPackageStream is called
    Then validation error is returned
    And error indicates invalid stream configuration
    And error provides configuration details

  @REQ-COMPR-121 @error
  Scenario: I/O error when temporary file creation fails
    Given a CompressPackageStream operation requiring temporary files
    And temporary file creation fails
    When CompressPackageStream is called
    Then I/O error is returned
    And error indicates temporary file creation failure
    And error provides file system details

  @REQ-COMPR-121 @error
  Scenario: I/O error when disk space is insufficient
    Given a CompressPackageStream operation requiring disk space
    And insufficient disk space is available
    When CompressPackageStream is called
    Then I/O error is returned
    And error indicates insufficient disk space
    And error provides disk space details

  @REQ-COMPR-121 @error
  Scenario: Context error on cancellation
    Given a CompressPackageStream operation
    And a cancelled context
    When CompressPackageStream is called with cancelled context
    Then context error is returned
    And error type is context cancellation

  @REQ-COMPR-121 @error
  Scenario: Context error on timeout
    Given a CompressPackageStream operation
    And a context with expired timeout
    When CompressPackageStream is called
    Then context error is returned
    And error type is context timeout
