@domain:basic_ops @m2 @REQ-API_BASIC-154 @spec(api_basic_operations.md#10312-packageerrorcontext-structure)
Feature: PackageErrorContext structure

  @REQ-API_BASIC-154 @happy
  Scenario: PackageErrorContext captures structured error details for package operations
    Given a structured error produced by a package operation
    When a PackageErrorContext is included
    Then the context captures operation-specific fields for package operations
    And the context supports debugging and correlation across layers
    And the context fields are stable and documented
    And the context can be serialized for logs or external reporting
    And the context is used consistently across package operations

