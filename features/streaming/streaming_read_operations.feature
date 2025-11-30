@domain:streaming @m2 @REQ-STREAM-036 @spec(api_streaming.md#1533-read-returns)
Feature: Streaming Read Operations

  @REQ-STREAM-036 @happy
  Scenario: Read returns define read operation results
    Given a NovusPack package
    And an open FileStream
    And a buffer for reading
    When Read method is called
    Then number of bytes read is returned
    And error is returned if read fails
    And Read implements io.Reader interface
    And Read follows standard Go reader semantics

  @REQ-STREAM-036 @happy
  Scenario: Read performs sequential reading
    Given a NovusPack package
    And an open FileStream
    When Read is called with buffer
    Then data is read into buffer
    And bytes read count reflects actual bytes read
    And stream position advances by bytes read
    And Read continues from current position

  @REQ-STREAM-036 @happy
  Scenario: Read handles partial reads correctly
    Given a NovusPack package
    And an open FileStream
    When Read is called with buffer larger than available data
    Then Read returns bytes available
    And error is nil if some data was read
    And Read follows io.Reader contract

  @REQ-STREAM-036 @error
  Scenario: Read handles errors correctly
    Given a NovusPack package
    And an open FileStream
    When Read encounters error
    Then error is returned
    And bytes read count reflects data read before error
    And error follows structured error format
