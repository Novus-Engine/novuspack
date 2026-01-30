@domain:basic_ops @m2 @REQ-API_BASIC-223 @spec(api_basic_operations.md#7217-readonlyerror-helper)
Feature: readOnlyError helper

  @REQ-API_BASIC-223 @happy
  Scenario: readOnlyError helper creates structured errors for read-only enforcement
    Given a write operation attempted on a read-only package
    When the enforcement layer creates an error
    Then the readOnlyError helper constructs a structured error
    And the error includes ReadOnlyErrorContext for debugging
    And the error uses a stable error type for programmatic inspection
    And the error message is actionable for consumers
    And the helper is used consistently across read-only enforcement paths

