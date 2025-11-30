@domain:streaming @m2 @REQ-STREAM-001 @spec(api_streaming.md#1-file-streaming-interface)
Feature: Stream file contents in chunks

  @happy
  Scenario: Sequential chunk reads
    Given a package with a large file
    When I read the file via a stream
    Then I should receive sequential chunks in order

  @happy
  Scenario: File stream read handles large files efficiently
    Given a package with very large file
    When file is read via stream
    Then chunks are read sequentially
    And memory usage is controlled
    And streaming is efficient

  @happy
  Scenario: Streaming decompresses compressed files transparently
    Given a package with compressed file
    When file is read via stream
    Then decompression occurs during streaming
    And decompressed chunks are returned
    And transparency is maintained

  @happy
  Scenario: Streaming decrypts encrypted files transparently
    Given a package with encrypted file
    When file is read via stream with correct key
    Then decryption occurs during streaming
    And decrypted chunks are returned
    And transparency is maintained

  @error
  Scenario: Streaming fails for non-existent files
    Given a package
    When stream read is attempted for non-existent file
    Then structured validation error is returned

  @error
  Scenario: Streaming respects context cancellation
    Given a package with large file
    And a cancelled context
    When stream read is called
    Then structured context error is returned
    And stream is closed properly
