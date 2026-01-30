@domain:basic_ops @m2 @REQ-API_BASIC-222 @spec(api_basic_operations.md#7216-implementation-body)
Feature: OpenPackageReadOnly implementation body

  @REQ-API_BASIC-222 @happy
  Scenario: OpenPackageReadOnly follows a wrapper-based implementation pattern
    Given OpenPackageReadOnly is called for an existing package file
    When it performs its implementation body
    Then it calls into shared open logic to load the package
    And it wraps the resulting package instance for read-only enforcement
    And it returns a Package interface that blocks writes
    And it propagates structured errors from the underlying open operation
    And it applies read-only enforcement consistently for all write operations

