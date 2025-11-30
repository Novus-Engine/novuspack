@domain:streaming @m2 @REQ-STREAM-008 @REQ-STREAM-016 @spec(api_streaming.md#151-stream-information)
Feature: File Stream Information

  @REQ-STREAM-008 @REQ-STREAM-016 @happy
  Scenario: Stream information methods report stream state
    Given a NovusPack package
    And an open FileStream
    When stream information methods are called
    Then Size returns total stream size in bytes
    And Position returns current position in stream
    And IsClosed returns boolean indicating closure status
    And information methods are thread-safe

  @REQ-STREAM-008 @REQ-STREAM-016 @happy
  Scenario: Size returns total stream size
    Given a NovusPack package
    And an open FileStream
    When Size method is called
    Then total size of stream in bytes is returned
    And size reflects actual file size
    And size enables progress calculation

  @REQ-STREAM-008 @REQ-STREAM-016 @happy
  Scenario: Position returns current read position
    Given a NovusPack package
    And an open FileStream
    When Position method is called
    Then current read position in stream is returned
    And position reflects bytes read so far
    And position advances as data is read

  @REQ-STREAM-008 @REQ-STREAM-016 @happy
  Scenario: IsClosed returns stream closure status
    Given a NovusPack package
    And a FileStream
    When IsClosed method is called
    Then boolean indicating if stream is closed is returned
    And false is returned for open streams
    And true is returned for closed streams

  @REQ-STREAM-008 @REQ-STREAM-016 @happy
  Scenario: Stream information enables stream state queries
    Given a NovusPack package
    And an open FileStream
    When stream state is queried using information methods
    Then size and position enable progress tracking
    And closure status prevents operations on closed streams
    And information supports stream management

  @REQ-STREAM-008 @REQ-STREAM-016 @happy
  Scenario: Stream offset and position parameters are validated
    Given a NovusPack package
    And an open FileStream
    When position-related operations are performed
    Then offset parameters are validated as non-negative
    And position parameters are validated within file size
    And validation prevents invalid position access

  @REQ-STREAM-008 @REQ-STREAM-016 @error
  Scenario: Stream information methods handle errors correctly
    Given a NovusPack package
    And a FileStream with error condition
    When stream information methods are called
    Then appropriate error handling is performed
    And errors reflect stream state issues
    And error handling follows structured error format
