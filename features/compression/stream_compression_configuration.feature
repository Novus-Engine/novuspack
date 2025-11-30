@domain:compression @m2 @REQ-COMPR-044 @spec(api_package_compression.md#11341-configuration)
Feature: Stream Compression Configuration

  @REQ-COMPR-044 @happy
  Scenario: StreamConfig is created with chunk size settings for streaming compression
    Given a compression operation for large packages
    When StreamConfig is created for streaming compression
    Then StreamConfig includes chunk size settings
    And chunk size is set to reasonable size such as 1MB
    And chunk size enables processing of large packages

  @REQ-COMPR-044 @happy
  Scenario: StreamConfig enables UseTempFiles for large packages
    Given a compression operation for large packages
    When StreamConfig is created with UseTempFiles enabled
    Then temporary files are used for packages exceeding memory limits
    And streaming avoids memory limitations
    And large packages can be processed

  @REQ-COMPR-044 @happy
  Scenario: StreamConfig configuration supports streaming compression workflow
    Given an open NovusPack package
    And a StreamConfig with streaming settings
    And a valid context
    And a compression type
    When CompressPackageStream is called with configuration
    Then package is compressed using streaming
    And streaming configuration is applied
    And workflow handles large packages efficiently

  @REQ-COMPR-044 @happy
  Scenario: StreamConfig configuration works with Write after streaming compression
    Given an open NovusPack package
    And package is compressed using CompressPackageStream
    And a valid context
    And an output file path
    When Write is called with CompressionNone
    Then compressed package is written to output file
    And no additional compression is applied
    And workflow completes successfully

  @REQ-COMPR-044 @happy
  Scenario: StreamConfig configuration enables processing packages exceeding RAM
    Given a large package that exceeds available RAM
    And a StreamConfig with appropriate settings
    When streaming compression is used
    Then package is processed without memory limitations
    And temporary files are used as needed
    And processing succeeds for large packages
