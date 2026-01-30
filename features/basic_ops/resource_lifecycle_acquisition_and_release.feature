@domain:basic_ops @m2 @REQ-API_BASIC-127 @spec(api_basic_operations.md#337-resource-lifecycle)
Feature: Resource lifecycle

  @REQ-API_BASIC-127 @happy
  Scenario: Resource lifecycle defines acquisition and release patterns
    Given a package opened for reading or writing
    When resources are acquired during operations
    Then resource acquisition patterns are defined
    And resources are released during close and cleanup
    And partial failure paths still release resources safely
    And resource lifecycle is consistent across package states
    And resource lifecycle supports predictable long-running processes

