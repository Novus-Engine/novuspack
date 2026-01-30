@domain:core @m2 @REQ-CORE-149 @spec(api_core.md#1254-fastwrite-method-behavior) @spec(api_writing.md#22-fastwrite-implementation-strategy)
Feature: FastWrite method behavior defines in-place update process

  @REQ-CORE-149 @happy
  Scenario: FastWrite performs in-place update process
    Given a package opened for writing at an existing path
    When FastWrite is called
    Then the update process is in place
    And the behavior matches the FastWrite specification
    And the existing file is updated without temp file
