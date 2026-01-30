@domain:basic_ops @m2 @REQ-API_BASIC-147 @spec(api_basic_operations.md#964-package-hassessionbase-method)
Feature: Package.HasSessionBase method

  @REQ-API_BASIC-147 @happy
  Scenario: Package.HasSessionBase returns true only when a session base is set
    Given an open package
    And no session base path has been configured
    When Package.HasSessionBase is called
    Then it returns false
    When SetSessionBase is called with a valid base path
    And Package.HasSessionBase is called again
    Then it returns true
    When ClearSessionBase is called
    Then Package.HasSessionBase returns false

