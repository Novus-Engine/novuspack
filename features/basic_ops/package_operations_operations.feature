@domain:basic_ops @m2 @REQ-API_BASIC-051 @spec(api_basic_operations.md#7-package-operations)
Feature: Package Operations

  @REQ-API_BASIC-051 @happy
  Scenario: Package operations include Validate method
    Given an open NovusPack package
    And a valid context
    When Validate is called
    Then package format is validated
    And package structure is validated
    And package integrity is validated

  @REQ-API_BASIC-051 @happy
  Scenario: Package operations include Defragment method
    Given an open NovusPack package
    And a valid context
    When Defragment is called
    Then unused space is removed
    And file entries are reorganized
    And data sections are compacted
    And package structure is optimized

  @REQ-API_BASIC-051 @happy
  Scenario: Package operations include GetInfo method
    Given an open NovusPack package
    And a valid context
    When GetInfo is called
    Then comprehensive package information is returned
    And package information includes file count and sizes
    And package information includes VendorID and AppID

  @REQ-API_BASIC-051 @happy
  Scenario: Package operations include ReadHeader method
    Given a NovusPack package file
    And a reader for the package file
    When ReadHeader is called with reader
    Then package header is read
    And header information is returned

  @REQ-API_BASIC-051 @happy
  Scenario: Package operations support checking package state
    Given a NovusPack package
    When package state is checked
    Then IsOpen method indicates if package is open
    And IsReadOnly method indicates if package is read-only
    And GetPath method returns current package file path

  @REQ-API_BASIC-051 @error
  Scenario: Package operations validate package is open
    Given a NovusPack package that is not open
    And a valid context
    When package operation is attempted
    Then validation error is returned
    And error indicates package must be open

  @REQ-API_BASIC-051 @error
  Scenario: Package operations handle context cancellation
    Given an open NovusPack package
    And a cancelled context
    When package operation is attempted
    Then context cancellation error is returned
    And operation is cancelled
