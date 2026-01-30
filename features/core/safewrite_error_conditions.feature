@domain:core @m2 @REQ-CORE-143 @spec(api_core.md#1246-safewrite-error-conditions) @spec(api_writing.md#16-safewrite-error-handling)
Feature: SafeWrite error conditions reference common writer error mapping table

  @REQ-CORE-143 @happy
  Scenario: SafeWrite errors use the common writer error mapping table
    Given a package opened for writing
    When SafeWrite encounters an error
    Then the error is mapped using the common writer error mapping table
    And the returned error is structured
    And error types follow the writer error mapping rules
