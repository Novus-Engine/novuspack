@domain:compression @m2 @REQ-COMPR-082 @spec(api_package_compression.md#134-memory-usage)
Feature: Compression Memory Usage

  @REQ-COMPR-082 @happy
  Scenario: Compression requires additional memory for compression buffers
    Given a compression operation
    When compression is performed
    Then additional memory is required for compression buffers
    And memory usage scales with package size
    And buffer memory is allocated during compression

  @REQ-COMPR-082 @happy
  Scenario: Large files use CompressPackageStream with memory limits
    Given a large package file
    And memory constraints exist
    When compression is performed
    Then CompressPackageStream is used with appropriate memory limits
    And advanced configuration manages memory usage
    And streaming avoids memory limitations

  @REQ-COMPR-082 @happy
  Scenario: Compression uses automatic fallback to disk buffering when memory limits exceeded
    Given a compression operation
    And memory limits are exceeded
    When compression is performed
    Then automatic fallback to disk buffering occurs
    And memory usage stays within limits
    And compression continues using disk

  @REQ-COMPR-082 @happy
  Scenario: Decompression requires memory for decompressed content
    Given a decompression operation
    When decompression is performed
    Then memory is required for decompressed content
    And memory needs may require temporary storage for large packages
    And streaming is used for memory-constrained environments

  @REQ-COMPR-082 @happy
  Scenario: Large files use chunked decompression with temp file management
    Given a large compressed package
    And memory constraints exist
    When decompression is performed
    Then chunked decompression is used
    And temporary file management handles memory constraints
    And decompression succeeds for large packages

  @REQ-COMPR-082 @happy
  Scenario: Memory limits enforce MaxMemoryUsage to prevent system OOM
    Given a compression or decompression operation
    And MaxMemoryUsage is configured
    When operation is performed
    Then MaxMemoryUsage limit is enforced
    And system out-of-memory errors are prevented
    And memory usage stays within configured limits
