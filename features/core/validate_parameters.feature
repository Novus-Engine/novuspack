@domain:core @m2 @REQ-CORE-118 @spec(api_core.md#1192-validate-parameters)
Feature: Validate parameters define context for cancellation and timeout

  @REQ-CORE-118 @happy
  Scenario: Validate accepts context for cancellation and timeout
    Given an opened package
    And a context that can be cancelled
    When Validate is called with the context
    Then validation can be cancelled via the context
    And timeouts are respected
