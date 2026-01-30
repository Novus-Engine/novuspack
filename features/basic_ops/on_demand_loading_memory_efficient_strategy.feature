@domain:basic_ops @m2 @REQ-API_BASIC-202 @spec(api_basic_operations.md#3341-on-demand-loading)
Feature: On-demand loading memory efficiency

  @REQ-API_BASIC-202 @happy
  Scenario: On-demand loading provides a memory-efficient loading strategy
    Given a large package with many file entries and metadata
    When on-demand loading is used
    Then memory footprint is reduced compared to fully loaded state
    And data is loaded only when required by operations
    And on-demand loading avoids unnecessary allocations
    And the strategy supports scalable package usage patterns
    And behavior aligns with documented on-demand loading design

