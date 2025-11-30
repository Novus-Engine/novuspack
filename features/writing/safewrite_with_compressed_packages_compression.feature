@domain:writing @m2 @REQ-WRITE-033 @spec(api_writing.md#521-safewrite-with-compressed-packages)
Feature: SafeWrite with Compressed Packages

  @REQ-WRITE-033 @happy
  Scenario: SafeWrite with compressed packages handles compression
    Given a NovusPack package
    And a compressed package
    When SafeWrite is called with compressed package
    Then SafeWrite handles compression correctly
    And decompression is required before writing
    And recompression occurs after writing
    And compression settings are preserved

  @REQ-WRITE-033 @happy
  Scenario: Decompression required before writing
    Given a NovusPack package
    And a compressed package
    When SafeWrite is performed
    Then package must be decompressed before writing
    And decompression enables write operations
    And decompression process is transparent to user

  @REQ-WRITE-033 @happy
  Scenario: Compression preservation maintains original settings
    Given a NovusPack package
    And a compressed package with original compression settings
    When SafeWrite is performed
    Then original compression settings are preserved in header
    And compression type is maintained
    And preservation ensures package consistency

  @REQ-WRITE-033 @happy
  Scenario: Recompression occurs after writing
    Given a NovusPack package
    And a compressed package
    When SafeWrite is performed
    Then package is recompressed after writing if originally compressed
    And recompression maintains compression state
    And recompression uses original compression settings

  @REQ-WRITE-033 @happy
  Scenario: Memory management uses streaming for large packages
    Given a NovusPack package
    And a large compressed package
    When SafeWrite is performed
    Then streaming is used for large compressed packages
    And memory management prevents excessive memory usage
    And streaming enables efficient handling of large packages

  @REQ-WRITE-033 @happy
  Scenario: Header comment and signature access remain uncompressed
    Given a NovusPack package
    And a compressed package
    When SafeWrite is performed
    Then header remains uncompressed for direct access
    And comment remains uncompressed for easy reading
    And signatures remain uncompressed for validation
    And uncompressed access supports package operations

  @REQ-WRITE-033 @error
  Scenario: SafeWrite with compressed packages handles errors correctly
    Given a NovusPack package
    And a compressed package with error conditions
    When SafeWrite encounters errors
    Then structured error is returned
    And error indicates compression-related failure
    And error follows structured error format
