@domain:basic_ops @m2 @REQ-API_BASIC-045 @spec(api_basic_operations.md#76-session-base-management) @spec(api_basic_operations.md#96-session-base-management)
Feature: Session base path is runtime-only

  @REQ-API_BASIC-045 @happy
  Scenario: Session base path is not persisted to the package file
    Given an open package
    And a session base path has been set at runtime
    When the package is written and closed
    And the package is re-opened from disk
    Then the session base path is not restored from the package file
    And GetSessionBase returns an empty string unless set again at runtime
    And the session base must be explicitly configured for each package session
    And persistence of package data does not include session base state

