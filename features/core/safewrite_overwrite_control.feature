@domain:core @m2 @REQ-CORE-157 @spec(api_core.md#1274-safewrite-overwrite-control) @spec(api_writing.md#11-packagesafewrite-method)
Feature: SafeWrite overwrite control requires overwrite flag for existing files

  @REQ-CORE-157 @happy
  Scenario: SafeWrite requires overwrite flag for existing target files
    Given a package opened for writing
    And the target path already exists
    When SafeWrite is called without overwrite flag
    Then the write fails or returns an error
    And the overwrite flag must be set to overwrite existing files
    And the behavior matches the SafeWrite overwrite control specification
