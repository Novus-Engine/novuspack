@skip @domain:streaming @m2 @spec(api_streaming.md#12-core-types)
Feature: Streaming Core Types and Structures

# This feature captures high-level expectations for core streaming types and structures.
# Detailed runnable scenarios live in the dedicated streaming feature files.

  @REQ-STREAM-019 @architecture
  Scenario: FileStream and StreamConfig provide core streaming structures
    Given a streaming implementation for large file access
    When the implementation exposes a FileStream
    Then the FileStream is configurable via a StreamConfig structure
    And configuration controls buffering and read behavior

  @REQ-STREAM-041 @architecture
  Scenario: BufferPool provides reusable buffers for streaming operations
    Given a series of streaming operations that allocate buffers
    When a BufferPool is used
    Then buffers are reused across operations
    And the pool exposes statistics and total size information
