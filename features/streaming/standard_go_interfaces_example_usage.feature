@domain:streaming @m2 @REQ-STREAM-039 @spec(api_streaming.md#1536-example-usage)
Feature: Standard Go Interfaces Example Usage

  @REQ-STREAM-039 @happy
  Scenario: Standard Go interfaces example usage demonstrates interface compatibility
    Given a NovusPack package
    And an open FileStream
    When standard Go interfaces are used as in example
    Then Read can be used as io.Reader
    And ReadAt can be used as io.ReaderAt
    And interfaces enable compatibility with standard Go libraries
    And example demonstrates proper usage patterns

  @REQ-STREAM-039 @happy
  Scenario: Read interface usage follows io.Reader pattern
    Given a NovusPack package
    And an open FileStream
    When Read is used as io.Reader
    Then buffer is provided for reading
    And number of bytes read is returned
    And error is returned if read fails
    And usage follows standard Go io.Reader interface

  @REQ-STREAM-039 @happy
  Scenario: ReadAt interface usage follows io.ReaderAt pattern
    Given a NovusPack package
    And an open FileStream
    When ReadAt is used as io.ReaderAt
    Then buffer and offset are provided
    And number of bytes read from offset is returned
    And error is returned if read fails
    And usage follows standard Go io.ReaderAt interface

  @REQ-STREAM-039 @happy
  Scenario: Standard interfaces enable library compatibility
    Given a NovusPack package
    And an open FileStream
    When FileStream is passed to standard Go libraries
    Then libraries can use FileStream as io.Reader
    And libraries can use FileStream as io.ReaderAt
    And compatibility enables integration with Go ecosystem
    And standard interfaces provide interoperability

  @REQ-STREAM-039 @error
  Scenario: Standard interfaces handle errors correctly
    Given a NovusPack package
    And an open FileStream with error condition
    When standard interfaces are used during error
    Then appropriate error is returned from Read
    And appropriate error is returned from ReadAt
    And errors follow standard Go error handling patterns
