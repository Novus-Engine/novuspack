@skip @domain:streaming @m2 @spec(api_streaming.md#11-purpose)
Feature: Streaming Definitions

# This feature captures high-level streaming expectations from the streaming API specification.
# Detailed runnable scenarios live in the dedicated streaming feature files.

  @REQ-STREAM-018 @documentation
  Scenario: FileStream enables streaming access for large files
    Given a file that is too large to load entirely into memory
    When a FileStream is created for that file
    Then the caller can read the file contents in chunks
    And the stream maintains the current position and size information

  @REQ-STREAM-013 @constraint
  Scenario: Streaming operations respect context cancellation
    Given a streaming operation that is reading file content
    When the provided context is cancelled
    Then the operation stops as soon as practical
    And resources associated with the stream are released
