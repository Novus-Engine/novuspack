@domain:basic_ops @m2 @REQ-API_BASIC-122 @spec(api_basic_operations.md#332-data-loading-strategies)
Feature: Data loading strategies

  @REQ-API_BASIC-122 @happy
  Scenario: Package supports eager and on-demand loading strategies
    Given a package opened from disk
    When data loading strategies are applied
    Then eager loading is defined for data required immediately after open
    And on-demand loading is defined for data loaded only when needed
    And the strategy avoids unnecessary I/O for unused data
    And loading strategy supports large packages efficiently
    And loading behavior is consistent across package operations

