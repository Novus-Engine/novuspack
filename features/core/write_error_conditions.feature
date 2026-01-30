@domain:core @m2 @REQ-CORE-135 @spec(api_core.md#1236-write-error-conditions) @spec(api_writing.md#57-error-handling)
Feature: Write error conditions reference common writer error mapping table

  @REQ-CORE-135 @happy
  Scenario: Write errors use the common writer error mapping table
    Given a package opened for writing
    When Write encounters an error
    Then the error is mapped using the common writer error mapping table
    And the returned error is structured
    And error types follow the writer error mapping rules
