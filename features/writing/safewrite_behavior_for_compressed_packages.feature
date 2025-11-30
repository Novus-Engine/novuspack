@domain:writing @m2 @REQ-WRITE-034 @spec(api_writing.md#5211-behavior-for-compressed-packages)
Feature: SafeWrite Behavior for Compressed Packages

  @REQ-WRITE-034 @happy
  Scenario: SafeWrite with compressed packages requires decompression before writing
    Given an open NovusPack package
    And the package is compressed (compression type in header flags)
    When SafeWrite is called with the target path
    Then package is decompressed before writing operations
    And modifications can be applied to decompressed package
    And decompression is performed correctly

  @REQ-WRITE-034 @happy
  Scenario: SafeWrite with compressed packages preserves original compression settings
    Given an open NovusPack package
    And the package is compressed with Zstd (compression type = 1)
    When SafeWrite is called with the target path
    Then original compression settings are preserved in header
    And compression type is maintained in header flags
    And header retains compression information

  @REQ-WRITE-034 @happy
  Scenario: SafeWrite with compressed packages recompresses after modifications
    Given an open NovusPack package
    And the package is compressed
    And package modifications have been applied
    When SafeWrite is called with the target path
    Then package is recompressed with original compression type
    And written package maintains compression state
    And compression is applied to file entries, data, and index

  @REQ-WRITE-034 @happy
  Scenario: SafeWrite with compressed packages uses streaming for large packages
    Given an open NovusPack package
    And the package is compressed
    And the package is large (>100MB)
    When SafeWrite is called with the target path
    Then streaming is used for large compressed packages
    And memory usage is controlled during decompression/recompression
    And streaming handles large content efficiently

  @REQ-WRITE-034 @happy
  Scenario: SafeWrite with compressed packages preserves header uncompressed
    Given an open NovusPack package
    And the package is compressed
    When SafeWrite is called with the target path
    Then header remains uncompressed for direct access
    And header can be read without decompression
    And header flags are accessible

  @REQ-WRITE-034 @happy
  Scenario: SafeWrite with compressed packages preserves comment uncompressed
    Given an open NovusPack package
    And the package is compressed
    And the package has a comment
    When SafeWrite is called with the target path
    Then comment remains uncompressed for easy reading
    And comment can be read without decompression
    And comment access is direct

  @REQ-WRITE-034 @happy
  Scenario: SafeWrite with compressed packages preserves signatures uncompressed
    Given an open NovusPack package
    And the package is compressed
    And the package has signatures
    When SafeWrite is called with the target path
    Then signatures remain uncompressed for validation
    And signatures can be validated without decompression
    And signature access is direct

  @REQ-WRITE-034 @error
  Scenario: SafeWrite with compressed packages returns error when decompression fails
    Given an open NovusPack package
    And the package is compressed
    And compressed data is corrupted
    When SafeWrite is called with the target path
    Then DecompressionFailure error is returned
    And error indicates decompression failure
    And error follows structured error format

  @REQ-WRITE-034 @error
  Scenario: SafeWrite with compressed packages returns error when recompression fails
    Given an open NovusPack package
    And the package is compressed
    And modifications have been applied
    And recompression fails
    When SafeWrite is called with the target path
    Then CompressionFailure error is returned
    And error indicates compression failure
    And error follows structured error format
