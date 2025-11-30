@domain:streaming @m2 @REQ-STREAM-006 @spec(api_streaming.md#1-file-streaming-interface)
Feature: File stream operations

  @happy
  Scenario: NewFileStream creates file stream for large files
    Given an open package with large file
    When NewFileStream is called with file path
    Then FileStream is created
    And stream is ready for reading
    And stream is not closed

  @REQ-STREAM-007 @happy
  Scenario: ReadChunk reads chunk of data from stream
    Given a FileStream for large file
    When ReadChunk is called with chunk size
    Then chunk of data is returned
    And chunk size matches request or end of file
    And stream position advances

  @REQ-STREAM-007 @happy
  Scenario: Seek seeks to specific position in stream
    Given a FileStream
    When Seek is called with offset and whence
    Then stream position is set
    And position is accessible via Position method
    And subsequent reads start from new position

  @REQ-STREAM-007 @happy
  Scenario: Close closes file stream
    Given an open FileStream
    When Close is called
    Then stream is closed
    And IsClosed returns true
    And resources are released

  @REQ-STREAM-007 @happy
  Scenario: GetStats gets streaming statistics
    Given a FileStream that has been used
    When GetStats is called
    Then streaming statistics are returned
    And bytes read are included
    And read operations count is included
    And performance metrics are included

  @REQ-STREAM-008 @happy
  Scenario: Size returns total size of stream
    Given a FileStream
    When Size is called
    Then total size of stream is returned
    And size matches file size

  @REQ-STREAM-008 @happy
  Scenario: Position returns current position in stream
    Given a FileStream that has been read
    When Position is called
    Then current position in stream is returned
    And position reflects read operations

  @REQ-STREAM-008 @happy
  Scenario: IsClosed checks if stream is closed
    Given a FileStream
    When IsClosed is called
    Then false is returned if stream is open
    And true is returned if stream is closed

  @REQ-STREAM-009 @happy
  Scenario: Progress returns detailed progress information
    Given a FileStream in use
    When Progress is called
    Then detailed progress information is returned
    And bytes read are included
    And percentage complete is included
    And elapsed time is included

  @REQ-STREAM-009 @happy
  Scenario: EstimatedTimeRemaining estimates time remaining
    Given a FileStream in use
    When EstimatedTimeRemaining is called
    Then estimated time remaining is returned
    And estimate is based on current progress
    And estimate is reasonable

  @REQ-STREAM-010 @happy
  Scenario: FileStream implements standard Go interfaces
    Given a FileStream instance
    When interface compliance is checked
    Then FileStream implements io.Reader interface
    And FileStream implements io.ReaderAt interface
    And standard Go functions work with FileStream
    And interface methods behave correctly

  @error
  Scenario: ReadChunk fails if stream is closed
    Given a closed FileStream
    When ReadChunk is called
    Then structured validation error is returned

  @error
  Scenario: Seek fails with invalid offset
    Given a FileStream
    When Seek is called with invalid offset
    Then structured validation error is returned

  @REQ-STREAM-013 @REQ-STREAM-014 @error
  Scenario: NewFileStream validates file path parameter
    Given an open package
    When NewFileStream is called with empty path
    Then structured validation error is returned
    And error indicates invalid path

  @REQ-STREAM-013 @REQ-STREAM-016 @error
  Scenario: Seek validates offset parameter
    Given a FileStream
    When Seek is called with negative offset
    Then structured validation error is returned
    And error indicates invalid offset

  @REQ-STREAM-013 @REQ-STREAM-017 @error
  Scenario: File stream operations respect context cancellation
    Given a FileStream in use
    And a cancelled context
    When file stream operation is called
    Then structured context error is returned
    And error type is context cancellation
    And stream is closed
    And resources are released
