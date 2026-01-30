@domain:core @m2 @REQ-CORE-142 @spec(api_core.md#1245-safewrite-detailed-behavior) @spec(api_writing.md#12-safewrite-implementation-strategy)
Feature: SafeWrite detailed behavior references atomic package writing specification

  @REQ-CORE-142 @happy
  Scenario: SafeWrite detailed behavior follows atomic writing spec
    Given a package opened for writing
    When SafeWrite is called
    Then the detailed behavior follows the atomic package writing specification
    And temp file and rename semantics are as specified
    And callers can rely on the documented behavior
