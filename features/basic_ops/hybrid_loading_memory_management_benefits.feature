@domain:basic_ops @m2 @REQ-API_BASIC-197 @spec(api_basic_operations.md#3324-memory-management-benefits)
Feature: Memory management benefits of hybrid loading

  @REQ-API_BASIC-197 @happy
  Scenario: Hybrid loading strategy provides memory management benefits
    Given large packages with substantial metadata
    When a hybrid loading strategy is used
    Then memory usage is reduced compared to fully eager loading
    And required data remains available for common operations
    And less frequently used data is loaded only when needed
    And hybrid loading balances performance and memory efficiency
    And memory management benefits align with documented loading design

