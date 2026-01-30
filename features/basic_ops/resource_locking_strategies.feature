@domain:basic_ops @m2 @REQ-API_BASIC-212 @spec(api_basic_operations.md#3363-resource-locking)
Feature: Resource locking strategies

  @REQ-API_BASIC-212 @happy
  Scenario: Resource locking defines locking strategies for concurrent access
    Given concurrent access patterns to a package
    When resources are shared between operations
    Then resource locking strategies are defined
    And locks protect shared mutable state
    And lock ordering avoids deadlocks
    And locks are used consistently across operations that mutate state
    And locking strategies align with documented concurrency limitations

