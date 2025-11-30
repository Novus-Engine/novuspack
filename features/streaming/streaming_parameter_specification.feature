@domain:streaming @m2 @REQ-STREAM-035 @REQ-STREAM-037 @REQ-STREAM-049 @spec(api_streaming.md#1532-read-parameters)
Feature: Streaming Parameter Specification

  @REQ-STREAM-035 @happy
  Scenario: Read parameters define read operation interface
    Given a NovusPack package
    And an open FileStream
    When Read method is called
    Then p parameter defines buffer to read data into
    And parameter is of type []byte
    And parameter enables sequential reading from stream

  @REQ-STREAM-035 @happy
  Scenario: Read buffer parameter accepts byte slice
    Given a NovusPack package
    And an open FileStream
    When Read is called with buffer p
    Then buffer parameter accepts []byte slice
    And data is read into provided buffer
    And buffer size determines read amount

  @REQ-STREAM-035 @happy
  Scenario: Read parameters support sequential reading
    Given a NovusPack package
    And an open FileStream
    When Read is called multiple times
    Then buffer parameter enables sequential reading
    And each read advances stream position
    And parameters follow io.Reader interface

  @REQ-STREAM-037 @happy
  Scenario: ReadAt parameters define random access read interface
    Given a NovusPack package
    And an open FileStream
    When ReadAt method is called
    Then p parameter defines buffer to read data into
    And off parameter defines offset to read from
    And parameters enable random access reading

  @REQ-STREAM-037 @happy
  Scenario: ReadAt buffer and offset parameters
    Given a NovusPack package
    And an open FileStream
    When ReadAt is called with p and off
    Then p parameter accepts []byte slice for buffer
    And off parameter accepts int64 for offset
    And offset enables reading from specific position

  @REQ-STREAM-037 @happy
  Scenario: ReadAt parameters support random access
    Given a NovusPack package
    And an open FileStream
    When ReadAt is called with different offsets
    Then parameters enable reading from any position
    And current stream position is not affected
    And parameters follow io.ReaderAt interface

  @REQ-STREAM-049 @happy
  Scenario: SetMaxTotalSize parameters define size limit configuration
    Given a NovusPack package
    And a BufferPool
    When SetMaxTotalSize is called
    Then maxSize parameter specifies maximum total size in bytes
    And parameter defines memory limit for buffer pool
    And configuration enables dynamic memory management

  @REQ-STREAM-049 @happy
  Scenario: SetMaxTotalSize maxSize parameter
    Given a NovusPack package
    And a BufferPool
    When SetMaxTotalSize is called with maxSize
    Then maxSize parameter is of type int64
    And value is specified in bytes
    And limit constrains total memory usage

  @REQ-STREAM-049 @error
  Scenario: Parameters handle validation errors
    Given a NovusPack package
    And an open FileStream or BufferPool
    When methods are called with invalid parameters
    Then structured error is returned
    And error indicates invalid parameter
    And error follows structured error format
