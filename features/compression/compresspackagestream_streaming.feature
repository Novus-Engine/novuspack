@domain:compression @m2 @REQ-COMPR-117 @spec(api_package_compression.md#51-compresspackagestream)
Feature: CompressPackageStream

  @REQ-COMPR-117 @happy
  Scenario: CompressPackageStream compresses large packages using streaming
    Given an open NovusPack package
    And a valid context
    And a compression type
    And a stream configuration
    When CompressPackageStream is called
    Then package content is compressed using streaming
    And streaming handles large packages efficiently
    And memory usage is controlled

  @REQ-COMPR-117 @happy
  Scenario: CompressPackageStream creates temporary files when needed
    Given an open NovusPack package
    And package size exceeds memory limits
    And a valid context
    And a stream configuration
    When CompressPackageStream is called
    Then temporary files are created as needed
    And temporary files are used for memory management
    And streaming avoids memory limitations

  @REQ-COMPR-117 @happy
  Scenario: CompressPackageStream compresses file entries, data, and index only
    Given an open NovusPack package
    And a valid context
    And a stream configuration
    When CompressPackageStream is called
    Then file entries are compressed
    And file data is compressed
    And package index is compressed
    And header remains uncompressed
    And comment remains uncompressed
    And signatures remain uncompressed

  @REQ-COMPR-117 @happy
  Scenario: CompressPackageStream uses adaptive processing based on configuration
    Given an open NovusPack package
    And a valid context
    And a stream configuration with adaptive settings
    When CompressPackageStream is called
    Then strategy automatically adjusts based on file size
    And strategy adapts to configuration
    And processing optimizes for conditions

  @REQ-COMPR-117 @happy
  Scenario: CompressPackageStream respects MaxMemoryUsage to prevent OOM
    Given an open NovusPack package
    And a valid context
    And a stream configuration with MaxMemoryUsage limit
    When CompressPackageStream is called
    Then memory usage stays within MaxMemoryUsage limit
    And out of memory errors are prevented
    And disk buffering is used if needed

  @REQ-COMPR-117 @happy
  Scenario: CompressPackageStream provides progress updates for long operations
    Given an open NovusPack package
    And a valid context
    And a stream configuration with ProgressCallback
    When CompressPackageStream is called
    Then progress updates are provided
    And callback receives bytes processed and total bytes
    And progress reporting supports user feedback

  @REQ-COMPR-117 @happy
  Scenario: CompressPackageStream uses parallel processing when enabled
    Given an open NovusPack package
    And a valid context
    And a stream configuration with parallel processing enabled
    When CompressPackageStream is called
    Then multiple CPU cores are used
    And parallel workers process chunks
    And compression performance is optimized

  @REQ-COMPR-117 @error
  Scenario: CompressPackageStream returns error if package is signed
    Given an open NovusPack package
    And package has signatures
    And a valid context
    And a stream configuration
    When CompressPackageStream is called
    Then security error is returned
    And error indicates package cannot be compressed when signed
    And error follows structured error format

  @REQ-COMPR-117 @error
  Scenario: CompressPackageStream returns error for invalid compression type
    Given an open NovusPack package
    And an invalid compression type
    And a valid context
    And a stream configuration
    When CompressPackageStream is called
    Then validation error is returned
    And error indicates invalid compression type
    And error follows structured error format
