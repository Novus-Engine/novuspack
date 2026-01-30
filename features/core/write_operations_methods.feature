@domain:core @m2 @REQ-CORE-125 @spec(api_core.md#1211-write-operations)
Feature: Write operations define Write, SafeWrite, and FastWrite methods

  @REQ-CORE-125 @happy
  Scenario: PackageWriter exposes Write, SafeWrite, and FastWrite
    Given a package opened for writing
    When the PackageWriter interface is used
    Then Write, SafeWrite, and FastWrite methods are available
    And each method has a defined purpose and durability guarantee
    And the behavior matches the write operations specification
