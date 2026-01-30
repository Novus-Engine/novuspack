@domain:core @m2 @REQ-CORE-126 @spec(api_core.md#1212-write-durability)
Feature: Write durability defines when changes are written to disk

  @REQ-CORE-126 @happy
  Scenario: Write durability is defined for each write method
    Given a package with pending changes
    When Write, SafeWrite, or FastWrite is invoked
    Then durability guarantees are as specified for each method
    And callers know when changes are persisted to disk
    And the behavior matches the write durability specification
