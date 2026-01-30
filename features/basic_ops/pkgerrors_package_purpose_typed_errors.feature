@domain:basic_ops @m2 @REQ-API_BASIC-179 @spec(api_basic_operations.md#2112-error-handling-package-purpose)
Feature: pkgerrors package purpose

  @REQ-API_BASIC-179 @happy
  Scenario: pkgerrors package purpose is typed errors with consistent context
    Given error handling across the API
    When error design is evaluated
    Then pkgerrors package purpose is to provide typed structured errors
    And errors include consistent context for debugging
    And errors support programmatic inspection and handling
    And error design avoids leaking sensitive information
    And error handling patterns are consistent across packages

