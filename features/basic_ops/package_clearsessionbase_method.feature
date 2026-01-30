@domain:basic_ops @m2 @REQ-API_BASIC-146 @spec(api_basic_operations.md#193-packageclearsessionbase-method)
Feature: Package.ClearSessionBase method

  @REQ-API_BASIC-146 @happy
  Scenario: ClearSessionBase clears the package-level session base
    Given an open package
    And a session base path has been set for the package
    When Package.ClearSessionBase is called
    Then the session base path is cleared from package state
    And HasSessionBase returns false
    And GetSessionBase returns an empty string
    And subsequent path operations do not use the previous session base

