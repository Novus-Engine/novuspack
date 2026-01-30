@domain:basic_ops @m2 @REQ-API_BASIC-178 @spec(api_basic_operations.md#21111-error-handling-package-pkgerrors)
Feature: pkgerrors package provides structured error handling

  @REQ-API_BASIC-178 @happy
  Scenario: pkgerrors package provides structured error handling primitives
    Given structured error handling for the API
    When errors are created and returned
    Then pkgerrors package provides structured error handling primitives
    And errors have a consistent typed shape across packages
    And error metadata is available for diagnosis and remediation
    And error construction supports wrapping with context
    And structured errors are used instead of legacy sentinel errors

