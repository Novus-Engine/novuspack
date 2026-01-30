@domain:streaming @m2 @REQ-STREAM-034 @spec(api_streaming.md#1531-purpose)
Feature: Standard Go interfaces purpose defines interface compatibility

  @REQ-STREAM-034 @happy
  Scenario: Standard Go interfaces purpose defines compatibility
    Given a FileStream that implements standard Go interfaces
    When Read or ReadAt is used
    Then the purpose defines interface compatibility with io.Reader and io.ReaderAt
    And the behavior matches the standard interfaces purpose specification
    And callers can use FileStream with standard Go I/O APIs
    And interface contracts are satisfied
