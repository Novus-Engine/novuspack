@domain:core @m2 @REQ-CORE-154 @spec(api_core.md#1271-fastwrite-path-restriction) @spec(api_writing.md#22-fastwrite-implementation-strategy)
Feature: FastWrite path restriction requires target path to match currently opened package path

  @REQ-CORE-154 @happy
  Scenario: FastWrite requires target path to match opened package path
    Given a package opened for writing
    When FastWrite is called
    Then the target path must match the currently opened package path
    And writing to a different path returns an error or is disallowed
    And the restriction is documented and enforced
