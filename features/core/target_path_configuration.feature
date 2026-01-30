@domain:core @m2 @REQ-CORE-127 @spec(api_core.md#1213-target-path-configuration) @spec(api_core.md#writing-target-path-configuration)
Feature: Target path configuration defines how package target path is configured

  @REQ-CORE-127 @happy
  Scenario: Target path is configured for write operations
    Given a package opened for writing
    When the target path is configured
    Then the configuration follows the target path specification
    And SetTargetPath or equivalent allows changing the target
    And write operations use the configured target path
