@domain:basic_ops @m2 @REQ-API_BASIC-210 @spec(api_basic_operations.md#3361-current-limitations)
Feature: Concurrency current limitations

  @REQ-API_BASIC-210 @happy
  Scenario: Current limitations define thread safety and concurrency constraints
    Given concurrent access to a package instance
    When concurrency constraints are evaluated
    Then current limitations define what is not safe concurrently
    And limitations identify operations that require exclusive access
    And limitations prevent incorrect assumptions by consumers
    And limitations inform locking and synchronization strategy
    And limitations align with documented thread safety behavior

