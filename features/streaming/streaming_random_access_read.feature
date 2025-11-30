@domain:streaming @m2 @REQ-STREAM-038 @spec(api_streaming.md#1535-readat-returns)
Feature: Streaming Random Access Read

  @REQ-STREAM-038 @happy
  Scenario: ReadAt returns define random access read results
    Given a NovusPack package
    And an open FileStream
    And a buffer for reading
    And an offset position
    When ReadAt method is called
    Then number of bytes read is returned
    And error is returned if read fails
    And ReadAt implements io.ReaderAt interface
    And ReadAt reads from specified offset

  @REQ-STREAM-038 @happy
  Scenario: ReadAt performs random access reading
    Given a NovusPack package
    And an open FileStream
    When ReadAt is called with offset
    Then data is read from specified offset
    And bytes read count reflects actual bytes read
    And stream position is not changed by ReadAt
    And ReadAt enables random access

  @REQ-STREAM-038 @happy
  Scenario: ReadAt enables concurrent access
    Given a NovusPack package
    And an open FileStream
    When ReadAt is called concurrently from different offsets
    Then each ReadAt reads from its specified offset
    And concurrent reads do not interfere
    And ReadAt is thread-safe

  @REQ-STREAM-038 @error
  Scenario: ReadAt handles invalid offset
    Given a NovusPack package
    And an open FileStream
    And an invalid offset
    When ReadAt is called with invalid offset
    Then error is returned
    And error indicates invalid offset
    And error follows structured error format
