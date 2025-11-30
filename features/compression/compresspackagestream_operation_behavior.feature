@domain:compression @m2 @REQ-COMPR-120 @spec(api_package_compression.md#513-compresspackagestream-behavior)
Feature: CompressPackageStream Operation Behavior

  @REQ-COMPR-120 @happy
  Scenario: CompressPackageStream uses streaming for large package content
    Given an open NovusPack package
    And package is large
    And a valid context
    And a compression type
    And a StreamConfig
    When CompressPackageStream is called
    Then streaming is used for compression
    And large packages are handled efficiently
    And memory limitations are avoided

  @REQ-COMPR-120 @happy
  Scenario: CompressPackageStream creates temporary files when needed
    Given an open NovusPack package
    And package exceeds available memory
    And a valid context
    And a compression type
    And a StreamConfig with temp file support
    When CompressPackageStream is called
    Then temporary files are created when needed
    And memory management uses disk storage
    And compression continues despite memory constraints

  @REQ-COMPR-120 @happy
  Scenario: CompressPackageStream compresses file entries, data, and index only
    Given an open NovusPack package
    And a valid context
    And a compression type
    And a StreamConfig
    When CompressPackageStream is called
    Then file entries are compressed
    And file data is compressed
    And package index is compressed
    And header, comment, and signatures remain uncompressed

  @REQ-COMPR-120 @happy
  Scenario: CompressPackageStream uses chunked processing for large files
    Given an open NovusPack package
    And package is large
    And a valid context
    And a compression type
    And a StreamConfig with chunk size settings
    When CompressPackageStream is called
    Then chunked processing is used
    And chunks are processed sequentially or in parallel
    And memory usage is controlled per chunk

  @REQ-COMPR-120 @error
  Scenario: CompressPackageStream returns error if package is signed
    Given an open NovusPack package
    And package has signatures
    And a valid context
    And a compression type
    And a StreamConfig
    When CompressPackageStream is called
    Then security error is returned
    And error indicates package is signed
    And error follows structured error format
