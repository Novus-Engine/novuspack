@domain:basic_ops @m2 @REQ-API_BASIC-044 @spec(api_basic_operations.md#76-session-base-management) @spec(api_basic_operations.md#96-session-base-management)
Feature: HasSessionBase reports whether session base is set

  @REQ-API_BASIC-044 @happy
  Scenario: HasSessionBase is true only when a session base path is currently set
    Given an open package
    And no session base path has been set
    When HasSessionBase is called
    Then it returns false
    When SetSessionBase is called with a valid base path
    And HasSessionBase is called again
    Then it returns true
    When ClearSessionBase is called
    Then HasSessionBase returns false

