@domain:streaming @m2 @REQ-STREAM-023 @spec(api_streaming.md#14-features)
Feature: Streaming Features and Capabilities

  @REQ-STREAM-023 @happy
  Scenario: Features define streaming capabilities
    Given a NovusPack package
    When streaming features are examined
    Then chunked reading is supported with configurable chunk sizes
    And buffer pool integration reuses buffers to reduce allocations
    And compression support handles compressed data transparently
    And encryption support handles encrypted data transparently
    And memory management provides configurable limits and pressure handling
    And performance monitoring provides built-in read speed and statistics

  @REQ-STREAM-023 @happy
  Scenario: Chunked reading provides optimal performance
    Given a NovusPack package
    And a FileStream with configured chunk size
    When chunked reading is performed
    Then chunks are read in configured sizes
    And chunk size is optimized for performance
    And memory usage is controlled through chunking

  @REQ-STREAM-023 @happy
  Scenario: Buffer pool integration reduces memory allocations
    Given a NovusPack package
    And a FileStream with buffer pool enabled
    When multiple chunks are read
    Then buffers are reused from pool
    And memory allocations are minimized
    And buffer pool prevents excessive allocations

  @REQ-STREAM-023 @happy
  Scenario: Compression and encryption are handled transparently
    Given a NovusPack package
    And a FileStream with compressed or encrypted data
    When streaming operations are performed
    Then compression is handled transparently
    And encryption is handled transparently
    And streaming interface remains consistent
