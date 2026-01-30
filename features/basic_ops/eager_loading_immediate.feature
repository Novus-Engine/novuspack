@domain:basic_ops @m2 @REQ-API_BASIC-194 @spec(api_basic_operations.md#3321-eager-loading-immediate)
Feature: Eager loading immediate

  @REQ-API_BASIC-194 @happy
  Scenario: Package open loads required data immediately
    Given a package opened from disk
    When eager loading is applied
    Then required data is loaded immediately on open
    And eager loading ensures required metadata is available after open
    And eager loading avoids additional I/O for required initial operations
    And eager loading is consistent with documented loading strategy
    And eager loading supports predictable API reads after open

