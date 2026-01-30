@domain:core @m2 @REQ-CORE-151 @spec(api_core.md#1256-fastwrite-error-conditions) @spec(api_writing.md#25-fastwrite-error-handling)
Feature: FastWrite error conditions reference common writer error mapping table

  @REQ-CORE-151 @happy
  Scenario: FastWrite errors use the common writer error mapping table
    Given a package opened for writing
    When FastWrite encounters an error
    Then the error is mapped using the common writer error mapping table
    And the returned error is structured
    And error types follow the writer error mapping rules
