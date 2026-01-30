@domain:core @m2 @REQ-CORE-156 @spec(api_core.md#1273-writing-to-a-new-path) @spec(api_basic_operations.md#8-packagesettargetpath-method) @spec(api_writing.md#43-writing-signed-package-content-to-new-path)
Feature: Writing to a new path allows changing target path via SetTargetPath

  @REQ-CORE-156 @happy
  Scenario: SetTargetPath allows changing target path for writing
    Given a package opened for writing
    When SetTargetPath is called with a new path
    Then the target path is updated for subsequent writes
    And writing to a new path is supported
    And the behavior matches the writing specification
