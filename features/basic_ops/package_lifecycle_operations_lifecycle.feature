@domain:basic_ops @m2 @REQ-API_BASIC-025 @spec(api_basic_operations.md#3-package-lifecycle-operations)
Feature: Package Lifecycle Operations

  @REQ-API_BASIC-025 @happy
  Scenario: Package lifecycle follows create pattern
    Given a valid context
    When package lifecycle is initiated
    Then first step is Create to create new package
    And package is prepared for operations

  @REQ-API_BASIC-025 @happy
  Scenario: Package lifecycle follows open pattern
    Given an existing package file
    And a valid context
    When package lifecycle is initiated
    Then Open step loads existing package
    And package is ready for operations

  @REQ-API_BASIC-025 @happy
  Scenario: Package lifecycle allows operations after create or open
    Given a package that is created or opened
    When package is in valid state
    Then Operations step allows various operations
    And files can be added
    And metadata can be modified
    And package operations can be performed

  @REQ-API_BASIC-025 @happy
  Scenario: Package lifecycle completes with close
    Given a package that has completed operations
    When package lifecycle is completed
    Then Close step releases resources
    And package file handle is closed
    And memory resources are freed
    And package state is cleared

  @REQ-API_BASIC-025 @happy
  Scenario: Package lifecycle pattern is simple and clear
    Given the NovusPack system
    When lifecycle pattern is examined
    Then pattern consists of four main steps
    And Create, Open, Operations, and Close are distinct phases
    And pattern supports both new and existing packages

  @REQ-API_BASIC-025 @error
  Scenario: Package lifecycle validates state transitions
    Given a package in invalid state
    When lifecycle operation is attempted
    Then validation error is returned
    And error indicates invalid state transition
    And operation is not performed
