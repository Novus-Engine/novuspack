@domain:core @m2 @REQ-CORE-056 @spec(api_core.md#127-allowed-target-paths-and-overwrite-behavior) @spec(api_writing.md#44-signed-package-writing-error-conditions)
Feature: Allowed target paths and overwrite behavior define write restrictions

  @REQ-CORE-056 @happy
  Scenario: Write operations enforce allowed target path and overwrite rules
    Given a package has a configured target path
    When a write operation is requested for a target path
    Then only allowed target paths are accepted
    And overwrite behavior follows the configured overwrite rules
    And invalid target paths produce a structured error
