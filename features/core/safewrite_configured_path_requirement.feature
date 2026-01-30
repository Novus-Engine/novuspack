@domain:core @m2 @REQ-CORE-155 @spec(api_core.md#1272-safewrite-configured-path-requirement) @spec(api_writing.md#12-safewrite-implementation-strategy)
Feature: SafeWrite configured path requirement writes to Package's configured target path

  @REQ-CORE-155 @happy
  Scenario: SafeWrite writes to the Package configured target path
    Given a package opened for writing with a configured target path
    When SafeWrite is called
    Then the write goes to the Package's configured target path
    And the configured path requirement is enforced
    And the behavior matches the SafeWrite specification
