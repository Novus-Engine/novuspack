@domain:basic_ops @m2 @REQ-API_BASIC-155 @spec(api_basic_operations.md#10313-securityerrorcontext-structure)
Feature: SecurityErrorContext structure

  @REQ-API_BASIC-155 @happy
  Scenario: SecurityErrorContext captures structured error details for security operations
    Given a structured error produced by a security operation
    When a SecurityErrorContext is included
    Then the context captures operation-specific fields for security operations
    And the context supports debugging security failures without leaking secrets
    And the context fields are stable and documented
    And the context can be serialized for logs or external reporting
    And the context is used consistently across security-related APIs

