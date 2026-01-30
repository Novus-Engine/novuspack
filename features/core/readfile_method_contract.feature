@domain:core @m2 @REQ-CORE-084 @spec(api_core.md#115-readfile-method-contract)
Feature: ReadFile method contract defines read-only file reading interface

  @REQ-CORE-084 @happy
  Scenario: ReadFile provides a read-only file reading contract
    Given an opened package
    When ReadFile is invoked via the PackageReader interface
    Then the method contract defines read-only file access
    And callers receive file content without modifying package state
    And the behavior matches the ReadFile method specification
