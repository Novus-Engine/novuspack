@domain:compression @m2 @REQ-COMPR-043 @spec(api_package_compression.md#1134-option-4-stream-compression-for-large-packages)
Feature: Stream Compression for Large Packages

  @REQ-COMPR-043 @happy
  Scenario: Stream compression uses streaming compression for large packages
    Given an open NovusPack package
    And package is large and may exceed available memory
    And a valid context
    And a StreamConfig with streaming settings
    And a compression type
    When stream compression workflow is followed
    Then CompressPackageStream is used for compression
    And streaming avoids memory limitations
    And large packages are handled efficiently

  @REQ-COMPR-043 @happy
  Scenario: Stream compression uses temporary files for memory management
    Given an open NovusPack package
    And package exceeds available memory
    And a valid context
    And a StreamConfig with UseTempFiles enabled
    When CompressPackageStream is called
    Then temporary files are used for large packages
    And memory limits are avoided
    And compression succeeds for large packages
