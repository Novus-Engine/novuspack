@domain:writing @m2 @REQ-WRITE-015 @REQ-WRITE-033 @spec(api_writing.md#14-safewrite-performance-characteristics)
Feature: SafeWrite Performance Characteristics

  @REQ-WRITE-015 @happy
  Scenario: SafeWrite performance characteristics show slower speed than FastWrite
    Given an open NovusPack package
    When SafeWrite performance is examined
    Then SafeWrite is slower than FastWrite for updates
    And complete rewrite is required
    And all data must be written

  @REQ-WRITE-015 @happy
  Scenario: SafeWrite performance characteristics show intelligent memory usage
    Given an open NovusPack package
    When SafeWrite performance is examined
    Then memory usage is intelligent
    And streaming is used for large files (>100MB)
    And in-memory writing is used for small files (<10MB)
    And memory thresholds are respected

  @REQ-WRITE-015 @happy
  Scenario: SafeWrite performance characteristics show higher disk I/O
    Given an open NovusPack package
    When SafeWrite performance is examined
    Then disk I/O is higher than FastWrite
    And temporary file is created
    And final file is written
    And complete data is written

  @REQ-WRITE-015 @happy
  Scenario: SafeWrite performance characteristics show maximum safety level
    Given an open NovusPack package
    When SafeWrite performance is examined
    Then safety level is maximum
    And guaranteed atomicity is provided
    And rollback capability exists

  @REQ-WRITE-015 @happy
  Scenario: SafeWrite performance characteristics show excellent scalability
    Given an open NovusPack package
    And the package is large (>1GB)
    When SafeWrite performance is examined
    Then scalability is excellent
    And packages of any size are handled
    And streaming supports large files

  @REQ-WRITE-015 @happy
  Scenario: SafeWrite performance characteristics show full rollback capability
    Given an open NovusPack package
    When SafeWrite encounters an error
    Then full rollback is performed
    And temporary file is cleaned up
    And no partial writes occur

  @REQ-WRITE-033 @happy
  Scenario: SafeWrite with compressed packages handles decompression before write
    Given an open NovusPack package
    And the package is compressed
    When SafeWrite is called with the target path
    Then package is decompressed before write operations
    And modifications are applied to decompressed package
    And package is ready for recompression

  @REQ-WRITE-033 @happy
  Scenario: SafeWrite with compressed packages preserves original compression settings
    Given an open NovusPack package
    And the package is compressed with Zstd (compression type = 1)
    When SafeWrite is called with the target path
    Then original compression settings are preserved in header
    And compression type is maintained
    And header flags retain compression information

  @REQ-WRITE-033 @happy
  Scenario: SafeWrite with compressed packages recompresses after modifications
    Given an open NovusPack package
    And the package is compressed
    And package modifications have been applied
    When SafeWrite is called with the target path
    Then package is recompressed with original compression type
    And written package maintains compression state
    And compression is applied correctly

  @REQ-WRITE-033 @happy
  Scenario: SafeWrite with compressed packages uses streaming for large packages
    Given an open NovusPack package
    And the package is compressed
    And the package is large (>100MB)
    When SafeWrite is called with the target path
    Then streaming is used for large compressed packages
    And memory usage is controlled
    And decompression/recompression uses streaming

  @REQ-WRITE-033 @happy
  Scenario: SafeWrite with compressed packages preserves header uncompressed
    Given an open NovusPack package
    And the package is compressed
    When SafeWrite is called with the target path
    Then header remains uncompressed for direct access
    And header can be read without decompression
    And header flags are accessible

  @REQ-WRITE-033 @happy
  Scenario: SafeWrite with compressed packages preserves signatures uncompressed
    Given an open NovusPack package
    And the package is compressed
    And the package has signatures
    When SafeWrite is called with the target path
    Then signatures remain uncompressed for validation
    And signatures can be validated without decompression
    And signature access is direct

  @REQ-WRITE-033 @error
  Scenario: SafeWrite with compressed packages returns error when decompression fails
    Given an open NovusPack package
    And the package is compressed
    And compressed data is corrupted
    When SafeWrite is called with the target path
    Then DecompressionFailure error is returned
    And error indicates decompression failure
    And error follows structured error format
