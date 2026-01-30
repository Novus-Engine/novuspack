@domain:security @m2 @REQ-SEC-077 @spec(api_security.md#21-securityvalidationresult-struct)
Feature: SecurityValidationResult struct provides validation result structure

  @REQ-SEC-077 @happy
  Scenario: SecurityValidationResult provides validation result structure
    Given security validation operations
    When validation results are returned
    Then SecurityValidationResult struct provides the result structure as specified
    And the behavior matches the SecurityValidationResult struct specification
    And result fields are populated correctly
    And callers can inspect validation status and details
