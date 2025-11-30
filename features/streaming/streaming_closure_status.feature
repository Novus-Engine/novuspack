@domain:streaming @m2 @REQ-STREAM-028 @spec(api_streaming.md#1514-isclosed-returns)
Feature: Streaming Closure Status

  @REQ-STREAM-028 @happy
  Scenario: IsClosed returns stream closure status
    Given a NovusPack package
    And an open FileStream
    When IsClosed method is called
    Then boolean indicating closure status is returned
    And IsClosed returns false for open stream
    And IsClosed returns true for closed stream

  @REQ-STREAM-028 @happy
  Scenario: IsClosed reflects stream state changes
    Given a NovusPack package
    And an open FileStream
    When stream is closed
    Then IsClosed returns true
    And closure status is persistent
    And IsClosed enables state checking

  @REQ-STREAM-028 @happy
  Scenario: IsClosed enables safe stream operations
    Given a NovusPack package
    And a FileStream
    When IsClosed is checked before operations
    Then operations can validate stream state
    And operations can prevent errors on closed streams
    And IsClosed enables defensive programming

  @REQ-STREAM-028 @error
  Scenario: IsClosed handles stream state correctly
    Given a NovusPack package
    When stream state changes
    Then IsClosed reflects current state accurately
    And IsClosed is thread-safe
