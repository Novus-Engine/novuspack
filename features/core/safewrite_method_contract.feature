@domain:core @m2 @REQ-CORE-137 @spec(api_core.md#124-safewrite-method-contract) @spec(api_writing.md#11-packagesafewrite-method)
Feature: SafeWrite method contract defines atomic write method interface

  @REQ-CORE-137 @happy
  Scenario: SafeWrite provides an atomic write method contract
    Given a package opened for writing
    When SafeWrite is invoked
    Then the method contract defines atomic write behavior
    And the interface matches the SafeWrite specification
    And callers receive atomic write guarantees
