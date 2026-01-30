@domain:core @m2 @REQ-CORE-147 @spec(api_core.md#1252-fastwrite-parameters) @spec(api_writing.md#21-packagefastwrite-method)
Feature: FastWrite parameters define context for cancellation and timeout

  @REQ-CORE-147 @happy
  Scenario: FastWrite accepts context for cancellation and timeout
    Given a package opened for writing
    And a context for cancellation
    When FastWrite is called with the context
    Then the context is accepted as the first parameter
    And cancellation and timeout are respected
    And parameters match the FastWrite method contract
