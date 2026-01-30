@domain:core @m2 @REQ-CORE-140 @spec(api_core.md#1243-safewrite-returns) @spec(api_writing.md#11-packagesafewrite-method)
Feature: SafeWrite returns define error return for write failures

  @REQ-CORE-140 @happy
  Scenario: SafeWrite returns an error on write failures
    Given a package opened for writing
    When SafeWrite is called and a failure occurs
    Then an error is returned
    And the error is structured
    And the return contract matches the SafeWrite specification
