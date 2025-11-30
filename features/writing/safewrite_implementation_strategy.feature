@domain:writing @m2 @REQ-WRITE-013 @spec(api_writing.md#12-safewrite-implementation-strategy)
Feature: SafeWrite Implementation Strategy

  @REQ-WRITE-013 @happy
  Scenario: SafeWrite implementation creates temporary file in same directory
    Given an open NovusPack package
    When SafeWrite is called with the target path
    Then temporary file is created in same directory as target
    And temp file has unique name
    And temp file is used for writing operations

  @REQ-WRITE-013 @happy
  Scenario: SafeWrite implementation streams data for large files
    Given an open NovusPack package
    And the package is large (>100MB)
    When SafeWrite is called with the target path
    Then data is streamed from source package or temp files
    And streaming handles large content efficiently
    And memory usage is controlled

  @REQ-WRITE-013 @happy
  Scenario: SafeWrite implementation uses in-memory data for small files
    Given an open NovusPack package
    And the package is small (<10MB)
    When SafeWrite is called with the target path
    Then data is written from memory
    And in-memory writing is efficient
    And memory thresholds are respected

  @REQ-WRITE-013 @happy
  Scenario: SafeWrite implementation atomically renames temp file to target
    Given an open NovusPack package
    And package has been written to temp file
    When SafeWrite completes successfully
    Then temp file is atomically renamed to target path
    And original file is replaced atomically
    And no partial writes are possible

  @REQ-WRITE-013 @happy
  Scenario: SafeWrite implementation automatically cleans up temp file on failure
    Given an open NovusPack package
    And write operation encounters an error
    When SafeWrite fails
    Then temporary file is automatically cleaned up
    And rollback is performed
    And no temp files are left behind

  @REQ-WRITE-013 @happy
  Scenario: SafeWrite implementation accepts compressionType parameter
    Given an open NovusPack package
    When SafeWrite is called with compressionType parameter
    Then compressionType specifies compression behavior
    And compression type 0 indicates no compression
    And compression types 1-3 specify compression methods

  @REQ-WRITE-013 @happy
  Scenario: SafeWrite implementation compresses file entries, data, and index
    Given an open NovusPack package
    And compressionType parameter is non-zero
    When SafeWrite is called with the target path
    Then file entries are compressed
    And data section is compressed
    And index is compressed
    And header, comment, and signatures remain uncompressed

  @REQ-WRITE-013 @error
  Scenario: SafeWrite implementation handles streaming failures gracefully
    Given an open NovusPack package
    And streaming operation fails
    When SafeWrite encounters streaming error
    Then streaming failure is handled gracefully
    And cleanup is performed
    And error is returned with details
