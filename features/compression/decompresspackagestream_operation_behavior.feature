@domain:compression @m2 @REQ-COMPR-125 @spec(api_package_compression.md#523-decompresspackagestream-behavior)
Feature: DecompressPackageStream Operation Behavior

  @REQ-COMPR-125 @happy
  Scenario: DecompressPackageStream uses streaming for large package content
    Given an open NovusPack package
    And package is compressed and large
    And a valid context
    And a StreamConfig
    When DecompressPackageStream is called
    Then streaming is used for decompression
    And large packages are handled efficiently
    And memory limitations are avoided

  @REQ-COMPR-125 @happy
  Scenario: DecompressPackageStream manages memory efficiently for large packages
    Given an open NovusPack package
    And package is compressed and large
    And a valid context
    And a StreamConfig with memory settings
    When DecompressPackageStream is called
    Then memory is managed efficiently
    And streaming prevents excessive memory usage
    And decompression succeeds for large packages

  @REQ-COMPR-125 @happy
  Scenario: DecompressPackageStream decompresses all compressed content
    Given an open NovusPack package
    And package is compressed
    And a valid context
    And a StreamConfig
    When DecompressPackageStream is called
    Then all compressed content is decompressed
    And file entries are decompressed
    And file data is decompressed
    And package index is decompressed

  @REQ-COMPR-125 @happy
  Scenario: DecompressPackageStream uses chunked processing for large files
    Given an open NovusPack package
    And package is compressed and large
    And a valid context
    And a StreamConfig with chunk size settings
    When DecompressPackageStream is called
    Then chunked processing is used
    And chunks are processed sequentially or in parallel
    And memory usage is controlled per chunk

  @REQ-COMPR-125 @happy
  Scenario: DecompressPackageStream updates package compression state
    Given an open NovusPack package
    And package is compressed
    And a valid context
    And a StreamConfig
    When DecompressPackageStream is called
    Then package compression state is updated
    And IsCompressed returns false after decompression
    And package state reflects decompressed status
