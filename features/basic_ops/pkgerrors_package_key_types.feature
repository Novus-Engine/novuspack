@domain:basic_ops @m2 @REQ-API_BASIC-180 @spec(api_basic_operations.md#21121-error-handling-package-key-types)
Feature: pkgerrors package key types

  @REQ-API_BASIC-180 @happy
  Scenario: pkgerrors defines key types for structured error modeling
    Given structured errors produced by the API
    When error types are inspected
    Then ErrorType is defined as a key type
    And PackageError is defined as a key type
    And ValidationErrorContext is defined as a key type
    And key types support structured error creation and handling
    And key types are used consistently across packages

