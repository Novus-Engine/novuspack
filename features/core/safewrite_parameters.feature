@domain:core @m2 @REQ-CORE-139 @spec(api_core.md#1242-safewrite-parameters) @spec(api_writing.md#11-packagesafewrite-method)
Feature: SafeWrite parameters define context and overwrite flag

  @REQ-CORE-139 @happy
  Scenario: SafeWrite accepts context and overwrite flag
    Given a package opened for writing
    And a context and overwrite flag
    When SafeWrite is called with context and overwrite
    Then the context is accepted for cancellation
    And the overwrite flag controls overwrite behavior
    And parameters match the SafeWrite contract
