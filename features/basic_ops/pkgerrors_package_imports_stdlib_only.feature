@domain:basic_ops @m2 @REQ-API_BASIC-181 @spec(api_basic_operations.md#2113-error-handling-package-imports)
Feature: pkgerrors package imports

  @REQ-API_BASIC-181 @happy
  Scenario: pkgerrors package imports define standard library only dependencies
    Given the pkgerrors package
    When dependencies are evaluated
    Then imports are limited to standard library dependencies
    And dependency boundaries avoid coupling to higher-level API packages
    And the package remains reusable across all domains
    And import boundaries support stable builds without dependency cycles
    And dependency choices support portability and reliability

