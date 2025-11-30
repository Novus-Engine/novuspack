@domain:streaming @m2 @REQ-STREAM-010 @spec(api_streaming.md#153-standard-go-interfaces)
Feature: Standard Go Interfaces

  @REQ-STREAM-010 @happy
  Scenario: FileStream implements standard Go interfaces for compatibility
    Given a NovusPack package
    And an open FileStream
    When standard Go interfaces are examined
    Then Read implements io.Reader interface
    And ReadAt implements io.ReaderAt interface
    And interfaces enable compatibility with standard Go libraries

  @REQ-STREAM-010 @happy
  Scenario: Read implements io.Reader interface
    Given a NovusPack package
    And an open FileStream
    When Read is used as io.Reader
    Then Read accepts buffer parameter p []byte
    And Read returns number of bytes read and error
    And implementation follows standard Go io.Reader contract

  @REQ-STREAM-010 @happy
  Scenario: ReadAt implements io.ReaderAt interface
    Given a NovusPack package
    And an open FileStream
    When ReadAt is used as io.ReaderAt
    Then ReadAt accepts buffer parameter p []byte and offset off int64
    And ReadAt returns number of bytes read and error
    And implementation follows standard Go io.ReaderAt contract

  @REQ-STREAM-010 @happy
  Scenario: Standard interfaces enable library compatibility
    Given a NovusPack package
    And an open FileStream
    When FileStream is passed to standard Go libraries
    Then libraries can use FileStream as io.Reader
    And libraries can use FileStream as io.ReaderAt
    And compatibility enables integration with Go ecosystem

  @REQ-STREAM-010 @happy
  Scenario: Sequential reading with io.Reader
    Given a NovusPack package
    And an open FileStream
    When Read is called sequentially
    Then data is read from current stream position
    And position advances after each read
    And sequential reading follows io.Reader semantics

  @REQ-STREAM-010 @happy
  Scenario: Random access reading with io.ReaderAt
    Given a NovusPack package
    And an open FileStream
    When ReadAt is called with different offsets
    Then data is read from specified offset
    And current stream position is not affected
    And random access enables efficient data retrieval

  @REQ-STREAM-010 @error
  Scenario: Standard interfaces handle errors correctly
    Given a NovusPack package
    And an open FileStream with error condition
    When standard interfaces are used during error
    Then appropriate error is returned from Read
    And appropriate error is returned from ReadAt
    And errors follow standard Go error handling patterns
