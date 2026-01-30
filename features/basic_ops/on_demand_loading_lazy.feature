@domain:basic_ops @m2 @REQ-API_BASIC-195 @spec(api_basic_operations.md#3322-on-demand-loading-lazy)
Feature: On-demand loading lazy

  @REQ-API_BASIC-195 @happy
  Scenario: Package loads data only when needed
    Given a package opened from disk
    When on-demand loading is applied
    Then non-required data is not loaded during open
    And data is loaded only when an operation requires it
    And on-demand loading reduces unnecessary I/O
    And on-demand loading supports large package usage patterns
    And on-demand loading is consistent with documented loading strategy

