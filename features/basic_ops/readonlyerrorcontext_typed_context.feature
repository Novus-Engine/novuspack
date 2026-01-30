@domain:basic_ops @m2 @REQ-API_BASIC-220 @spec(api_basic_operations.md#7214-readonlyerrorcontext-structure)
Feature: ReadOnlyErrorContext typed context

  @REQ-API_BASIC-220 @happy
  Scenario: ReadOnlyErrorContext provides typed context for read-only enforcement errors
    Given a write operation attempted on a read-only package
    When a structured error is returned
    Then ReadOnlyErrorContext provides typed context for the failure
    And context includes relevant operation details without leaking secrets
    And context supports programmatic inspection and logging
    And context fields are stable and documented
    And context is used consistently for read-only enforcement errors

