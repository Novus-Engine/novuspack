@domain:basic_ops @m2 @REQ-API_BASIC-108 @spec(api_basic_operations.md#76-session-base-management) @spec(api_basic_operations.md#96-session-base-management)
Feature: GetSessionBase returns the session base path

  @REQ-API_BASIC-108 @happy
  Scenario: GetSessionBase returns the current session base path or an empty string when not set
    Given an open package
    And no session base path has been set
    When GetSessionBase is called
    Then an empty string is returned
    When SetSessionBase is called with a valid base path
    And GetSessionBase is called again
    Then the configured base path is returned
    And the returned base path represents the package's current session base setting

