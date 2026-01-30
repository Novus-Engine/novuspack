@domain:basic_ops @m2 @REQ-API_BASIC-109 @spec(api_basic_operations.md#76-session-base-management) @spec(api_basic_operations.md#96-session-base-management)
Feature: ClearSessionBase clears the session base path

  @REQ-API_BASIC-109 @happy
  Scenario: ClearSessionBase removes the current session base path from the package state
    Given an open package
    And a session base path has been set
    When ClearSessionBase is called
    Then the session base path is removed from in-memory package state
    And HasSessionBase reports that no session base is set
    And GetSessionBase returns an empty string
    And subsequent path operations do not use a previously configured session base
    And the package behaves as if no session base was ever configured

