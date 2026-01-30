@domain:basic_ops @m2 @REQ-API_BASIC-156 @spec(api_basic_operations.md#10314-ioerrorcontext-structure)
Feature: IOErrorContext structure

  @REQ-API_BASIC-156 @happy
  Scenario: IOErrorContext captures structured error details for I/O operations
    Given a structured error produced by an I/O operation
    When an IOErrorContext is included
    Then the context captures I/O-specific fields such as paths and operation types
    And the context supports debugging I/O failures consistently
    And the context fields are stable and documented
    And the context can be serialized for logs or external reporting
    And the context is used consistently across I/O-heavy APIs

