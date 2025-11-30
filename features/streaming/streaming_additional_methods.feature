@domain:streaming @m2 @REQ-STREAM-024 @spec(api_streaming.md#15-additional-methods)
Feature: Streaming Additional Methods

  @REQ-STREAM-024 @happy
  Scenario: Additional methods provide extended streaming operations
    Given a NovusPack package
    And an open FileStream
    When additional streaming methods are used
    Then Size returns total size of stream
    And Position returns current read position
    And IsClosed checks if stream is closed
    And Progress returns detailed progress information
    And EstimatedTimeRemaining estimates completion time
    And Read implements io.Reader interface
    And ReadAt implements io.ReaderAt interface

  @REQ-STREAM-024 @happy
  Scenario: Stream information methods report stream state
    Given a NovusPack package
    And an open FileStream
    When stream information methods are called
    Then Size returns total stream size in bytes
    And Position returns current position in stream
    And IsClosed returns boolean indicating closure status
    And information methods are thread-safe

  @REQ-STREAM-024 @happy
  Scenario: Progress monitoring methods track stream progress
    Given a NovusPack package
    And an open FileStream
    When progress monitoring methods are called
    Then Progress returns bytesRead, totalBytes, readSpeed, and elapsed time
    And EstimatedTimeRemaining returns estimated completion time
    And progress information enables progress tracking

  @REQ-STREAM-024 @happy
  Scenario: Standard Go interfaces enable compatibility
    Given a NovusPack package
    And an open FileStream
    When standard Go interfaces are used
    Then Read implements io.Reader interface
    And ReadAt implements io.ReaderAt interface
    And interfaces enable compatibility with standard Go libraries
