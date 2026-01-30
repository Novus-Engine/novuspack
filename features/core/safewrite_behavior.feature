@domain:core @m2 @REQ-CORE-141 @spec(api_core.md#1244-safewrite-behavior) @spec(api_writing.md#12-safewrite-implementation-strategy)
Feature: SafeWrite behavior defines atomic write process

  @REQ-CORE-141 @happy
  Scenario: SafeWrite performs atomic write process
    Given a package opened for writing
    When SafeWrite is called
    Then the write process is atomic
    And the behavior matches the SafeWrite specification
    And existing file is replaced only on success
