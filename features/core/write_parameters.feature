@domain:core @m2 @REQ-CORE-131 @spec(api_core.md#1232-write-parameters) @spec(api_writing.md#533-packagewrite-method)
Feature: Write parameters define context for cancellation and timeout

  @REQ-CORE-131 @happy
  Scenario: Write accepts context for cancellation and timeout
    Given a package opened for writing
    And a context for cancellation
    When Write is called with the context
    Then the context is accepted as the first parameter
    And cancellation and timeout are respected
    And the parameters match the Write method contract
