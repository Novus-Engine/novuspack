@domain:basic_ops @m2 @REQ-API_BASIC-148 @spec(api_basic_operations.md#965-session-base-purpose)
Feature: Session base purpose for absolute path conversion

  @REQ-API_BASIC-148 @happy
  Scenario: Session base exists to support package-level absolute path conversion relative to a base
    Given an open package
    And a session base path configured for the package session
    When a relative path is converted to an absolute filesystem path
    Then the session base path is used as the base for conversion
    And path conversion behavior is consistent across file operations
    And session base purpose is to centralize absolute path conversion logic
    And session base is package-scoped for the duration of the open session
    And session base can be cleared to disable base-relative conversion

